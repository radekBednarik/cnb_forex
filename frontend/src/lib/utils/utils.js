export function getDateString(toSubtract = 0) {
	const currDate = new Date();

	if (toSubtract !== 0) {
		currDate.setDate(currDate.getDate() + toSubtract);
	}
	const year = currDate.getFullYear();
	const month = String(currDate.getMonth() + 1).padStart(2, '0');
	const day = String(currDate.getDate()).padStart(2, '0');

	return `${year}-${month}-${day}`;
}
