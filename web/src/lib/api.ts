import type { RunsResponse, ChartDataResponse, ConfigResponse, AuthStatusResponse, LeaderboardResponse } from './types';
import type { GlobalFilters } from './filters.svelte';

const BASE = '/api';

function handleUnauthorized(res: Response) {
	if (res.status === 401 && window.location.pathname !== '/settings') {
		window.location.href = '/settings';
	}
}

async function get<T>(path: string): Promise<T> {
	const res = await fetch(`${BASE}${path}`);
	if (!res.ok) {
		handleUnauthorized(res);
		const body = await res.json().catch(() => ({ error: res.statusText }));
		throw new Error(body.error || res.statusText);
	}
	return res.json();
}

async function put<T>(path: string, body: unknown): Promise<T> {
	const res = await fetch(`${BASE}${path}`, {
		method: 'PUT',
		headers: { 'Content-Type': 'application/json' },
		body: JSON.stringify(body)
	});
	if (!res.ok) {
		handleUnauthorized(res);
		const data = await res.json().catch(() => ({ error: res.statusText }));
		throw new Error(data.error || res.statusText);
	}
	return res.json();
}

async function post<T>(path: string): Promise<T> {
	const res = await fetch(`${BASE}${path}`, { method: 'POST' });
	if (!res.ok) {
		handleUnauthorized(res);
		const data = await res.json().catch(() => ({ error: res.statusText }));
		throw new Error(data.error || res.statusText);
	}
	return res.json();
}

export interface RunsParams {
	weeks?: number;
	day?: string;
	minDistance?: number;
	all?: boolean;
	sort?: string;
	asc?: boolean;
	refresh?: boolean;
}

function buildQuery(params: RunsParams): string {
	const q = new URLSearchParams();
	if (params.weeks) q.set('weeks', String(params.weeks));
	if (params.day) q.set('day', params.day);
	if (params.minDistance) q.set('minDistance', String(params.minDistance));
	if (params.all) q.set('all', 'true');
	if (params.sort) q.set('sort', params.sort);
	if (params.asc) q.set('asc', 'true');
	if (params.refresh) q.set('refresh', 'true');
	const str = q.toString();
	return str ? `?${str}` : '';
}

export interface LeaderboardParams {
	page?: number;
	year?: number;
	minDistance?: number;
	maxDistance?: number;
	maxHR?: number;
}

function buildLeaderboardQuery(params: LeaderboardParams): string {
	const q = new URLSearchParams();
	if (params.page) q.set('page', String(params.page));
	if (params.year) q.set('year', String(params.year));
	if (params.minDistance) q.set('minDistance', String(params.minDistance));
	if (params.maxDistance) q.set('maxDistance', String(params.maxDistance));
	if (params.maxHR) q.set('maxHR', String(params.maxHR));
	const str = q.toString();
	return str ? `?${str}` : '';
}

const KM_TO_METERS = 1000;

/** Convert global filters to RunsParams for /api/runs and /api/chart-data */
export function filtersToRunsParams(f: GlobalFilters, extra?: Pick<RunsParams, 'sort' | 'asc'>): RunsParams {
	const p: RunsParams = {};
	if (f.weeks > 0) p.weeks = f.weeks;
	if (f.day) p.day = f.day;
	if (f.minDistance > 0) p.minDistance = f.minDistance;
	if (f.showAll) p.all = true;
	if (extra?.sort) p.sort = extra.sort;
	if (extra?.asc) p.asc = extra.asc;
	return p;
}

/** Convert global filters to LeaderboardParams for /api/leaderboard */
export function filtersToLeaderboardParams(f: GlobalFilters, extra?: { page?: number }): LeaderboardParams {
	const p: LeaderboardParams = {};
	if (f.year > 0) p.year = f.year;
	if (f.minDistance > 0) p.minDistance = f.minDistance * KM_TO_METERS;
	if (f.maxDistance > 0) p.maxDistance = f.maxDistance * KM_TO_METERS;
	if (f.maxHR > 0) p.maxHR = f.maxHR;
	if (extra?.page) p.page = extra.page;
	return p;
}

export const api = {
	getRuns: (params: RunsParams = {}) => get<RunsResponse>(`/runs${buildQuery(params)}`),
	getChartData: (params: RunsParams = {}) => get<ChartDataResponse>(`/chart-data${buildQuery(params)}`),
	getConfig: () => get<ConfigResponse>('/config'),
	putConfig: (body: { zone2_hr?: number; age?: number }) => put<ConfigResponse>('/config', body),
	getAuthStatus: () => get<AuthStatusResponse>('/auth/status'),
	refresh: () => post<{ status: string }>('/refresh'),
	getLeaderboard: (params: LeaderboardParams = {}) => get<LeaderboardResponse>(`/leaderboard${buildLeaderboardQuery(params)}`),
	refreshLeaderboard: () => post<{ status: string }>('/leaderboard/refresh')
};
