const KM_TO_MILE = 1.60934;

export function formatDistance(meters: number): string {
	const mi = meters / 1000 / KM_TO_MILE;
	return `${mi.toFixed(2)} mi`;
}

export function formatTotalDistance(km: number): string {
	const mi = km / KM_TO_MILE;
	return `${mi.toFixed(1)} mi`;
}

export function formatPace(secondsPerKm: number): string {
	if (secondsPerKm <= 0) return '-';
	const secondsPerMi = secondsPerKm * KM_TO_MILE;
	const mins = Math.floor(secondsPerMi / 60);
	const secs = Math.round(secondsPerMi % 60);
	return `${mins}:${secs.toString().padStart(2, '0')} /mi`;
}

export function formatDuration(seconds: number): string {
	const h = Math.floor(seconds / 3600);
	const m = Math.floor((seconds % 3600) / 60);
	const s = seconds % 60;
	if (h > 0) return `${h}h ${m}m`;
	return `${m}m ${s}s`;
}

export function formatEF(ef: number): string {
	return ef.toFixed(4);
}

export function formatHR(hr: number): string {
	return Math.round(hr).toString();
}

export function formatDate(isoString: string): string {
	const d = new Date(isoString);
	return d.toLocaleDateString('en-GB', { day: 'numeric', month: 'short', year: 'numeric' });
}

export function efficiencyFactor(distance: number, movingTime: number, avgHR: number): number {
	if (movingTime === 0 || avgHR === 0) return 0;
	const speed = distance / movingTime;
	return speed / avgHR;
}
