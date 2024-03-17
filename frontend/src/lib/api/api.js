/**
 *
 * @param {string} dateFrom
 * @param {string} dateTo
 *
 */
export async function getData(dateFrom, dateTo, wrStore) {
	const response = await fetch(
		`http://localhost:8080/api/dashboard/v1/data?dateFrom=${dateFrom}&dateTo=${dateTo}`
	);
	const currData = await response.json();
	wrStore.set(JSON.stringify(currData));
}
