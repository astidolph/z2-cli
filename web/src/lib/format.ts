const KM_TO_MILE = 1.60934;

export function formatDistance(meters: number): string {
	const km = meters / 1000;
	const mi = km / KM_TO_MILE;
	return `${mi.toFixed(2)} mi (${km.toFixed(2)} km)`;
}

export function formatTotalDistance(km: number): string {
	const mi = km / KM_TO_MILE;
	return `${mi.toFixed(1)} mi (${km.toFixed(1)} km)`;
}

export function formatPace(secondsPerKm: number): string {
	if (secondsPerKm <= 0) return '-';
	const secondsPerMi = secondsPerKm * KM_TO_MILE;
	const minsKm = Math.floor(secondsPerKm / 60);
	const secsKm = Math.round(secondsPerKm % 60);
	const minsMi = Math.floor(secondsPerMi / 60);
	const secsMi = Math.round(secondsPerMi % 60);
	return `${minsMi}:${secsMi.toString().padStart(2, '0')} /mi (${minsKm}:${secsKm.toString().padStart(2, '0')} /km)`;
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
