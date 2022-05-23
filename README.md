# PORTS

Service that maintains ports data.

## Overview

The logical structure of the service is reflected in a number of packages:

- service: core package that maintains the business logic, updating and serving port data.
- internal/parser: JSON parsing logic that provides the Feed interface that is used by the service to process bulk update of port details.
- store: embedded key-value store (LMST) that is biased towards fast writes.
- transport/rest: http service that interfaces the core service providing endpoints for client interaction.

The project is organised around the core service, which is 'surrounded' by various 'adapters' that depend on the domain constructs defined by the service. Store, and transport are such adapters. This allows the adapters to be replaced by a different mechanism, for example the transport with a new adapter that works from the file system instead of processing http requests.

## Running the service

### Directly (compiled from source)

Assuming that go >= 1.15 is present and you're in the ./cmd directory:

- Download dependencies

```sh
go get ./...
```

- Build the executable

```sh
go build -o ports
```

- Run the service

```sh
./ports -db_path=<path-for-database-files>
```

The service listens on port 8000 by default, that can be overridden by the `addr` parameter.

### Docker

From the project root directory:

- Build the container

```sh
docker build . -t ports
```

- Run the service in the container (exposing port 8000)

```sh
docker run --rm --name ports -p:8000:8000 ports
```

## HTTP Endpoints

### POST /ports

Allows a client to upload a json structure as defined in the original task description. The JSON input is parsed as a stream. Existing entries are overwritten by the provided update.

Example (where ports.json contains the updates):

```sh
curl -v 'localhost:8000/ports' -d@ports.json
```

### GET /ports

Allows a client to retrieve port information. Provides an array of port objects ordered by port UN/LOCODE. Partial or paginated results are retrieved by providing `from_id` and `limit` parameters. If `from_id` is not provided then the first N entry is returned where N is defined by the `limit` parameter. If limit is not provided, then the maximum number of returned entries will be set to 1000.

```curl
curl -v 'localhost:8000/ports?from_id=ZWBUQ?limit=2'
```

```json
[
  {
    "ID": "ZWHRE",
    "Details": {
      "name": "Harare",
      "city": "Harare",
      "country": "Zimbabwe",
      "coordinates": [
        31.03351,
        -17.8251657
      ],
      "province": "Harare",
      "timezone": "Africa/Harare",
      "unlocs": [
        "ZWHRE"
      ]
    }
  },
  {
    "ID": "ZWUTA",
    "Details": {
      "name": "Mutare",
      "city": "Mutare",
      "country": "Zimbabwe",
      "coordinates": [
        32.650351,
        -18.9757714
      ],
      "province": "Manicaland",
      "timezone": "Africa/Harare",
      "unlocs": [
        "ZWUTA"
      ]
    }
  }
]

```
