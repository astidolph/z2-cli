package strava

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const baseURL = "https://www.strava.com/api/v3"

type Client struct {
	accessToken string
	httpClient  *http.Client
}

func NewClient(accessToken string) *Client {
	return &Client{
		accessToken: accessToken,
		httpClient:  &http.Client{Timeout: 30 * time.Second},
	}
}

func (c *Client) get(path string) ([]byte, error) {
	req, err := http.NewRequest("GET", baseURL+path, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.accessToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	return body, nil
}

type Activity struct {
	ID               int64   `json:"id"`
	Name             string  `json:"name"`
	Type             string  `json:"type"`
	SportType        string  `json:"sport_type"`
	StartDateLocal   string  `json:"start_date_local"`
	Distance         float64 `json:"distance"`
	MovingTime       int     `json:"moving_time"`
	ElapsedTime      int     `json:"elapsed_time"`
	AverageHeartrate float64 `json:"average_heartrate"`
	MaxHeartrate     float64 `json:"max_heartrate"`
	HasHeartrate     bool    `json:"has_heartrate"`
}

func (a *Activity) StartTime() (time.Time, error) {
	return time.Parse("2006-01-02T15:04:05Z", a.StartDateLocal)
}

func (c *Client) GetActivities(after, before time.Time, page, perPage int) ([]Activity, error) {
	path := fmt.Sprintf("/athlete/activities?after=%d&before=%d&page=%d&per_page=%d",
		after.Unix(), before.Unix(), page, perPage)

	body, err := c.get(path)
	if err != nil {
		return nil, err
	}

	var activities []Activity
	if err := json.Unmarshal(body, &activities); err != nil {
		return nil, fmt.Errorf("could not parse activities: %w", err)
	}

	return activities, nil
}

func (c *Client) GetAllRunsSince(since time.Time) ([]Activity, error) {
	var allRuns []Activity
	page := 1
	now := time.Now()

	for {
		activities, err := c.GetActivities(since, now, page, 100)
		if err != nil {
			return nil, err
		}
		if len(activities) == 0 {
			break
		}

		for _, a := range activities {
			if a.Type == "Run" || a.SportType == "Run" {
				allRuns = append(allRuns, a)
			}
		}

		page++
	}

	return allRuns, nil
}
