import { getDateString } from '$lib/utils/utils';

const dateFrom = getDateString(-7);
const dateTo = getDateString();


/** @type {import('./$types').LayoutLoad} */
export async function load({ params }) {
    const response = await fetch(
        `http://localhost:8080/api/currencies/v1/symbols?dateFrom=${dateFrom}&dateTo=${dateTo}`
    );
    return await response.json()
}