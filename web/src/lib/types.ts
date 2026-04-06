export interface Activity {
	id: number;
	name: string;
	type: string;
	sport_type: string;
	start_date_local: string;
	distance: number;
	moving_time: number;
	elapsed_time: number;
	average_heartrate: number;
	max_heartrate: number;
	has_heartrate: boolean;
}

export interface Summary {
	count: number;
	avg_ef: number;
	avg_hr: number;
	avg_pace: number;
	total_km: number;
}

export interface RunsResponse {
	current_runs: Activity[] | null;
	prior_runs: Activity[] | null;
	current: Summary;
	prior: Summary;
	zone2_hr: number;
	weeks_back: number;
	year?: number;
	ef_trend: number;
}

export interface ChartDataResponse {
	dates: string[] | null;
	ef: (number | null)[] | null;
	pace: (number | null)[] | null;
	pace_mi: (number | null)[] | null;
	distance: (number | null)[] | null;
	distance_mi: (number | null)[] | null;
	hr: (number | null)[] | null;
}

export interface LeaderboardResponse {
	runs: Activity[] | null;
	total_count: number;
	page: number;
	page_size: number;
}

export interface ConfigResponse {
	zone2_hr: number;
}

export interface AuthStatusResponse {
	authenticated: boolean;
	message?: string;
	expires_at?: number;
}
