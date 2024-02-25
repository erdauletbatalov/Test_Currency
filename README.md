# Currency Exchange Service API Documentation
## Introduction
The Currency Exchange Service API provides endpoints to retrieve and save currency exchange rates. It interacts with a local Microsoft SQL Server database to store and retrieve currency data. This document outlines the available endpoints and their usage.

## Base URL

http://localhost:PORT/
## Endpoints
### 1. Retrieve Currency Exchange Rates
GET /currency/{date}/{code}
This endpoint retrieves currency exchange rates for a specific date and currency code.

date: The date for which currency exchange rates are requested (format: "YYYY-MM-DD").
code: (Optional) The currency code for which exchange rates are requested. If not provided, exchange rates for all currencies on the specified date are returned.

Request
```bash
GET /currency/2024-02-25/USD
```

Response
```json
[
  {
    "title": "United States Dollar",
    "code": "USD",
    "value": 1.00,
    "a_date": "2024-02-25"
  }
]
```

### 2. Save Currency Exchange Rates

POST /currency/save/{date}

This endpoint saves currency exchange rates fetched from a public API into the local database.

- date: The date for which currency exchange rates are to be fetched and saved (format: "YYYY-MM-DD").
Request
```bash
POST /currency/save/2024-02-25
```
Response
```json
{
  "success": true
}
```
## Service Operation
The Currency Exchange Service operates as follows:

1) When a request is made to save currency exchange rates for a specific date (POST /currency/save/{date}), the service retrieves data from a public API provided by the National Bank.
2) The retrieved data is asynchronously saved to the local database in a separate Goroutine.
Upon successful saving of the data, a success response is returned to the client without waiting for the database operation to complete.
3) If there is an error during the database operation, it is logged, and an error response is returned to the client.
This documentation provides an overview of the available endpoints, their usage, and the operation of the Currency Exchange Service.