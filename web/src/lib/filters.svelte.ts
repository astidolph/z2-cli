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

// Individual reactive properties — ensures fine-grained Svelte 5 tracking
let weeks = $state(defaults.weeks);
let year = $state(defaults.year);
let day = $state(defaults.day);
let minDistance = $state(defaults.minDistance);
let maxDistance = $state(defaults.maxDistance);
let maxHR = $state(defaults.maxHR);
let showAll = $state(defaults.showAll);

/**
 * Returns a reactive object backed by individual $state variables.
 * Reading any property inside a $derived or $effect will track it.
 */
export function getFilters(): GlobalFilters {
	return {
		get weeks() { return weeks; },
		get year() { return year; },
		get day() { return day; },
		get minDistance() { return minDistance; },
		get maxDistance() { return maxDistance; },
		get maxHR() { return maxHR; },
		get showAll() { return showAll; },
	} as GlobalFilters;
}

export function setFilters(patch: Partial<GlobalFilters>) {
	// Enforce mutual exclusivity: setting weeks clears year and vice versa
	if (patch.weeks !== undefined && patch.weeks > 0) {
		patch.year = 0;
	} else if (patch.year !== undefined && patch.year > 0) {
		patch.weeks = 0;
	}

	if (patch.weeks !== undefined) weeks = patch.weeks;
	if (patch.year !== undefined) year = patch.year;
	if (patch.day !== undefined) day = patch.day;
	if (patch.minDistance !== undefined) minDistance = patch.minDistance;
	if (patch.maxDistance !== undefined) maxDistance = patch.maxDistance;
	if (patch.maxHR !== undefined) maxHR = patch.maxHR;
	if (patch.showAll !== undefined) showAll = patch.showAll;
}

export function resetFilters() {
	weeks = defaults.weeks;
	year = defaults.year;
	day = defaults.day;
	minDistance = defaults.minDistance;
	maxDistance = defaults.maxDistance;
	maxHR = defaults.maxHR;
	showAll = defaults.showAll;
}
