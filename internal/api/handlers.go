package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/z2-cli/internal/auth"
	"github.com/z2-cli/internal/cache"
	"github.com/z2-cli/internal/chart"
	"github.com/z2-cli/internal/service"
	"github.com/z2-cli/internal/stats"
)

// frontendOrigin extracts the origin (scheme://host) that the user's browser
// is on, so post-auth redirects land on the frontend (which may be the Vite
// dev server on a different port). The Referer-based origin is only trusted
// if it shares the same hostname as the request (allowing different ports
// for local development).
func frontendOrigin(r *http.Request) string {
	requestHost := stripPort(r.Host)

	if ref := r.Header.Get("Referer"); ref != "" {
		if u, err := url.Parse(ref); err == nil && u.Host != "" {
			if stripPort(u.Host) == requestHost {
				return u.Scheme + "://" + u.Host
			}
		}
	}
	scheme := r.Header.Get("X-Forwarded-Proto")
	if scheme == "" {
		scheme = "http"
	}
	return scheme + "://" + r.Host
}

// stripPort returns the hostname portion of a host:port string.
func stripPort(host string) string {
	if i := strings.LastIndex(host, ":"); i != -1 {
		return host[:i]
	}
	return host
}

func handleAuthLogin(w http.ResponseWriter, r *http.Request) {
	config, err := auth.LoadConfig()
	if err != nil {
		http.Redirect(w, r, "/settings?auth_error="+url.QueryEscape("Strava credentials not configured"), http.StatusTemporaryRedirect)
		return
	}

	// Build redirect URI from the incoming request
	scheme := r.Header.Get("X-Forwarded-Proto")
	if scheme == "" {
		scheme = "http"
	}
	redirectURI := scheme + "://" + r.Host + "/api/auth/callback"

	// Generate state and sign it with client secret for CSRF protection
	state := auth.GenerateState()
	sig := auth.SignState(state, config.ClientSecret)

	http.SetCookie(w, &http.Cookie{
		Name:     "z2_oauth_state",
		Value:    state + "." + sig,
		Path:     "/api/auth/callback",
		MaxAge:   300,
		HttpOnly: true,
		Secure:   scheme == "https",
		SameSite: http.SameSiteLaxMode,
	})

	// Remember the frontend origin so the callback can redirect back to it
	// (important when the Vite dev server is on a different port).
	http.SetCookie(w, &http.Cookie{
		Name:     "z2_auth_origin",
		Value:    frontendOrigin(r),
		Path:     "/api/auth/callback",
		MaxAge:   300,
		HttpOnly: true,
		Secure:   scheme == "https",
		SameSite: http.SameSiteLaxMode,
	})

	authURL := auth.BuildAuthorizeURL(config.ClientID, redirectURI, state)
	http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
}

func handleAuthCallback(w http.ResponseWriter, r *http.Request) {
	// Read the frontend origin saved during login so redirects land on the
	// correct host (e.g. Vite dev server on :5173 vs Go server on :8080).
	origin := ""
	if c, err := r.Cookie("z2_auth_origin"); err == nil {
		origin = c.Value
	}
	settingsRedirect := func(query string) string {
		return origin + "/settings?" + query
	}

	// Check if Strava returned an error
	if stravaErr := r.URL.Query().Get("error"); stravaErr != "" {
		http.Redirect(w, r, settingsRedirect("auth_error="+url.QueryEscape("Strava authorization denied")), http.StatusTemporaryRedirect)
		return
	}

	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")
	if code == "" || state == "" {
		http.Redirect(w, r, settingsRedirect("auth_error="+url.QueryEscape("Missing authorization code or state")), http.StatusTemporaryRedirect)
		return
	}

	// Validate state from cookie
	cookie, err := r.Cookie("z2_oauth_state")
	if err != nil {
		http.Redirect(w, r, settingsRedirect("auth_error="+url.QueryEscape("Missing state cookie — please try again")), http.StatusTemporaryRedirect)
		return
	}

	parts := strings.SplitN(cookie.Value, ".", 2)
	if len(parts) != 2 || parts[0] != state {
		http.Redirect(w, r, settingsRedirect("auth_error="+url.QueryEscape("Invalid state — please try again")), http.StatusTemporaryRedirect)
		return
	}

	config, err := auth.LoadConfig()
	if err != nil {
		http.Redirect(w, r, settingsRedirect("auth_error="+url.QueryEscape("Could not load config")), http.StatusTemporaryRedirect)
		return
	}

	if !auth.ValidateSignedState(parts[0], parts[1], config.ClientSecret) {
		http.Redirect(w, r, settingsRedirect("auth_error="+url.QueryEscape("Invalid state signature — please try again")), http.StatusTemporaryRedirect)
		return
	}

	// Clear the state and origin cookies
	http.SetCookie(w, &http.Cookie{
		Name:     "z2_oauth_state",
		Value:    "",
		Path:     "/api/auth/callback",
		MaxAge:   -1,
		HttpOnly: true,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "z2_auth_origin",
		Value:    "",
		Path:     "/api/auth/callback",
		MaxAge:   -1,
		HttpOnly: true,
	})

	token, err := auth.ExchangeCode(config.ClientID, config.ClientSecret, code)
	if err != nil {
		log.Printf("Token exchange failed: %v", err)
		http.Redirect(w, r, settingsRedirect("auth_error="+url.QueryEscape("Token exchange failed — please try again")), http.StatusTemporaryRedirect)
		return
	}

	if err := auth.SaveToken(token); err != nil {
		http.Redirect(w, r, settingsRedirect("auth_error="+url.QueryEscape("Could not save token")), http.StatusTemporaryRedirect)
		return
	}

	setSessionCookie(w, r)
	http.Redirect(w, r, settingsRedirect("auth_success=true"), http.StatusTemporaryRedirect)
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func handleAuthStatus(w http.ResponseWriter, r *http.Request) {
	token, err := auth.LoadToken()
	if err != nil {
		writeJSON(w, http.StatusOK, map[string]any{
			"authenticated": false,
			"message":       "Not authenticated — run 'z2-cli auth' to connect Strava",
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"authenticated": !token.IsExpired(),
		"expires_at":    token.ExpiresAt,
	})
}

func handleGetConfig(w http.ResponseWriter, r *http.Request) {
	config, err := auth.LoadConfig()
	if err != nil {
		log.Printf("Failed to load config: %v", err)
		writeError(w, http.StatusInternalServerError, "could not load config")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"zone2_hr": config.Zone2HR,
	})
}

func handlePutConfig(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1024)
	var body struct {
		Zone2HR *int `json:"zone2_hr"`
		Age     *int `json:"age"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	config, err := auth.LoadConfig()
	if err != nil {
		log.Printf("Failed to load config: %v", err)
		writeError(w, http.StatusInternalServerError, "could not load config")
		return
	}

	switch {
	case body.Zone2HR != nil:
		if *body.Zone2HR <= 0 {
			writeError(w, http.StatusBadRequest, "zone2_hr must be positive")
			return
		}
		config.Zone2HR = *body.Zone2HR
	case body.Age != nil:
		if *body.Age <= 0 {
			writeError(w, http.StatusBadRequest, "age must be positive")
			return
		}
		config.Zone2HR = 180 - *body.Age
	default:
		writeError(w, http.StatusBadRequest, "provide zone2_hr or age")
		return
	}

	if err := auth.SaveConfig(config); err != nil {
		log.Printf("Failed to save config: %v", err)
		writeError(w, http.StatusInternalServerError, "could not save config")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"zone2_hr": config.Zone2HR,
	})
}

func handleGetRuns(w http.ResponseWriter, r *http.Request) {
	query, err := parseRunsQuery(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	result, err := service.FetchRuns(query)
	if err != nil {
		log.Printf("Failed to fetch runs: %v", err)
		writeError(w, http.StatusInternalServerError, "could not fetch runs")
		return
	}

	// Add EF and trend to the response
	type runsResponse struct {
		*service.RunsResult
		EFTrend float64 `json:"ef_trend"`
	}

	writeJSON(w, http.StatusOK, runsResponse{
		RunsResult: result,
		EFTrend:    stats.TrendPercent(result.Current, result.Prior),
	})
}

type chartDataResponse struct {
	Dates      []string   `json:"dates"`
	EF         []*float64 `json:"ef"`
	Pace       []*float64 `json:"pace"`
	PaceMi     []*float64 `json:"pace_mi"`
	Distance   []*float64 `json:"distance"`
	DistanceMi []*float64 `json:"distance_mi"`
	HR         []*float64 `json:"hr"`
}

func handleGetChartData(w http.ResponseWriter, r *http.Request) {
	query, err := parseRunsQuery(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	result, err := service.FetchRuns(query)
	if err != nil {
		log.Printf("Failed to fetch chart data: %v", err)
		writeError(w, http.StatusInternalServerError, "could not fetch chart data")
		return
	}

	data := chart.BuildChartData(result.CurrentRuns)
	resp := chartDataResponse{
		Dates:      data.Dates,
		EF:         lineDataToFloats(data.EF),
		Pace:       lineDataToFloats(data.Pace),
		PaceMi:     lineDataToFloats(data.PaceMi),
		Distance:   lineDataToFloats(data.Distance),
		DistanceMi: lineDataToFloats(data.DistanceMi),
		HR:         lineDataToFloats(data.HR),
	}

	writeJSON(w, http.StatusOK, resp)
}

func handleRefresh(w http.ResponseWriter, r *http.Request) {
	if err := cache.Invalidate(); err != nil {
		log.Printf("Failed to invalidate cache: %v", err)
		writeError(w, http.StatusInternalServerError, "could not clear cache")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "cache cleared"})
}

// parseRunsQuery extracts RunsQuery fields from URL query parameters.
func parseRunsQuery(r *http.Request) (service.RunsQuery, error) {
	q := r.URL.Query()
	query := service.RunsQuery{
		WeeksBack: 12,
		SortBy:    "date",
	}

	if v := q.Get("weeks"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil || n <= 0 {
			return query, fmt.Errorf("invalid weeks parameter")
		}
		query.WeeksBack = n
	}
	if v := q.Get("day"); v != "" {
		query.Day = v
	}
	if v := q.Get("minDistance"); v != "" {
		f, err := strconv.ParseFloat(v, 64)
		if err != nil || f < 0 {
			return query, fmt.Errorf("invalid minDistance parameter")
		}
		query.MinDistance = f
	}
	if v := q.Get("all"); v == "true" {
		query.ShowAll = true
	}
	if v := q.Get("sort"); v != "" {
		switch v {
		case "date", "distance", "time", "hr", "pace", "ef":
			query.SortBy = v
		default:
			return query, fmt.Errorf("invalid sort parameter")
		}
	}
	if v := q.Get("asc"); v == "true" {
		query.Ascending = true
	}
	if v := q.Get("refresh"); v == "true" {
		query.ForceRefresh = true
	}

	return query, nil
}

func handleGetLeaderboard(w http.ResponseWriter, r *http.Request) {
	page := 1
	if v := r.URL.Query().Get("page"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil || n <= 0 {
			writeError(w, http.StatusBadRequest, "invalid page parameter")
			return
		}
		page = n
	}

	result, err := service.FetchLeaderboard(page)
	if err != nil {
		log.Printf("Failed to fetch leaderboard: %v", err)
		writeError(w, http.StatusInternalServerError, "could not fetch leaderboard")
		return
	}

	writeJSON(w, http.StatusOK, result)
}

func handleRefreshLeaderboard(w http.ResponseWriter, r *http.Request) {
	if err := service.RefreshLeaderboard(); err != nil {
		log.Printf("Failed to refresh leaderboard: %v", err)
		writeError(w, http.StatusInternalServerError, "could not refresh leaderboard data")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "leaderboard refreshed"})
}

// lineDataToFloats converts go-echarts LineData values to plain float64 pointers.
// nil Values are preserved as nil so they serialize to JSON null (Chart.js skips nulls).
func lineDataToFloats(data []opts.LineData) []*float64 {
	out := make([]*float64, len(data))
	for i, d := range data {
		if d.Value == nil {
			continue
		}
		var f float64
		switch v := d.Value.(type) {
		case string:
			parsed, err := strconv.ParseFloat(v, 64)
			if err != nil {
				continue
			}
			f = parsed
		case float64:
			f = v
		case int:
			f = float64(v)
		default:
			continue
		}
		out[i] = &f
	}
	return out
}
