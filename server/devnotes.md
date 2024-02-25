# devnotes

## dashboard endpoint

### endpoint

`/api/dashboard/v1/data`

### query params

- dateFrom

- dateTo

- date

### response body structure

```json
{
  "data": [
    {
      "2024-02-25": [
        {
          "currName": "dolar",
          "currSymbol": "USD",
          "value": "20.01"
        }
      ]
    }
  ]
}
```

### limiting

- response body should contain max of 30 days

  - pagination needed for longer time spans

    - will lead to repeated calls to retrieve all needed data

## browser

- store data in session storage or cache on the server?
