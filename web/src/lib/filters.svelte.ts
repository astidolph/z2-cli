/**
 * Global filter store shared across all pages.
 *
 * Time range: either `weeks` or `year` is active, never both.
 * Setting one clears the other.
 */

export interface GlobalFilters {
	weeks: number;
	year: number; // 0 = not set
	day: string; // '' = all days
	minDistance: number; // km, 0 = no minimum
	maxDistance: number; // km, 0 = no maximum
	maxHR: number; // 0 = no limit
	showAll: boolean;
}

const defaults: GlobalFilters = {
	weeks: 12,
	year: 0,
	day: '',
	minDistance: 0,
	maxDistance: 0,
	maxHR: 0,
	showAll: false
};

let filters = $state<GlobalFilters>({ ...defaults });

export function getFilters(): GlobalFilters {
	return filters;
}

export function setFilters(patch: Partial<GlobalFilters>) {
	// Enforce mutual exclusivity: setting weeks clears year and vice versa
	if (patch.weeks !== undefined && patch.weeks > 0) {
		patch.year = 0;
	} else if (patch.year !== undefined && patch.year > 0) {
		patch.weeks = 0;
	}
	Object.assign(filters, patch);
}

export function resetFilters() {
	Object.assign(filters, { ...defaults });
}
