# Welcome to the Country Information Service!
This is a service to fetch and display various information about a country.

## Deployment
The service can either be hosted locally or online on services such as [Render](https://render.com/).

## Setup & Installation

### Prerequisites

- **Go:** Version 1.23 or higher.

### 1. Clone the Repository

```bash
git clone https://github.com/iamk3v/assignment-1-cloud.git
```
### 2. Install Dependencies
```bash
go mod tidy
```

### 3. Run the Application
```bash
go run main.go
```

## API
This API has three resource root paths:
```
/countryinfo/v1/info/
/countryinfo/v1/population/
/countryinfo/v1/status/
```
For info regarding how to use the endpoint, the following syntax applies:
- `{:value}` - indicates mandatory input parameters
- `{value}` - indicates optional input parameters
- `{?:param}` - indicates mandatory HTTP parameters
- `{?param}` - Indicates optional HTTP parameters

## Info Endpoint
The info endpoint returns general information for a given country, based on their [2-letter ISO 3166-2 country codes](https://en.wikipedia.org/wiki/ISO_3166-2)

This endpoint can be invoked with:
```
Method: GET
Path: info/{:two_letter_country_code}
```
To limit the amount of cities, the {?limit} parameter can be passed:
```
Method: GET
Path: info/{:two_letter_country_code}{?limit={limit}}
```
### Example Requests
```
GET info/no
GET info/no/?limit=10
```
### Example Response
Request: `GET info/no/?limit=10`

Response: 
```json
{
  "Name": "Norway",
  "Continents": ["Europe"],
  "Population": 5379475,
  "Languages": {"nno":"Norwegian Nynorsk","nob":"Norwegian Bokmål","smi":"Sami"},
  "Borders": ["FIN","SWE","RUS"],
  "Flag": "https://flagcdn.com/w320/no.png",
  "Capital": ["Oslo"],
  "Cities": ["Abelvaer","Adalsbruk","Adland","Agotnes","Agskardet","Aker","Akkarfjord","Akrehamn","Al","Alen"]
}
```

## Population Endpoint
The population endpoint returns population levels for individual years for a given country, 
as well as the mean value of those, based on their [2-letter ISO 3166-2 country codes](https://en.wikipedia.org/wiki/ISO_3166-2).

This endpoint can be invoked with:
```
Method: GET
Path: population/{:two_letter_country_code}
```
To limit the amount of years, the {?limit} parameter can be passed:
```
Method: GET
Path: population/{:two_letter_country_code}{?limit={:startYear-endYear}}
```
### Example Requests
```
GET population/no
GET population/no?limit=2010-2015
```
### Example Response
Request: `GET population/no?limit=2010-2015`

Response:
```json
{
  "Mean": 5044395,
  "Values": [
    {
      "year": 2010,
      "value": 4889252
    },
    {
      "year": 2011,
      "value": 4953088
    },
    {
      "year": 2012,
      "value": 5018573
    },
    {
      "year": 2013,
      "value": 5079623
    },
    {
      "year": 2014,
      "value": 5137232
    },
    {
      "year": 2015,
      "value": 5188607
    }
  ]
}
```
## Status Endpoint
The status endpoint returns the availability of the individual services this service depends on.

This endpoint can be invoked with:
```
Method: GET
Path: status/
```
### Example Request
```
GET status/
```

### Example Response
```json
{
  "countriesnowapi": "200 - https://http.cat/200",
  "restcountriesapi": "200 - https://http.cat/200",
  "uptime": "00h:22m:56s",
  "version": "v1"
}
```
