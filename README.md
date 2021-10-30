# Message Store

Implementation of an event message store backend.

## Capacity Requirements

1. Throughput: process 10 mi obj/second - POST
2. Query Load - 100 reqs/second

## Schema

Objects will have the following fields:

1. id: 128-bit UUID string
2. Attributes: Array consisting of the following nested fields:

  1. key: string of upto 256 chars
  2. value: string upto 256 chars

3. timestamp: unix-epoch timestamp in microseconds

## APIs

We will expose an REST API for CRUD operations.

### Retrieve Object

```text
GET /mydata/<id>
```

Return an object with the provided uuid.

Response codes:

- 200 for non-empty response
- 400 invalid request params
- 404 no record found

### Search Objects

```text
GET /mydata/<start_timestamp>/<end_timestamp>/<key>/<value>
```

Return a list of object UUID satisfying the criteria:

1. key/value is matching
2. timestamp is in the range (start_timestamp, end_timestamp)

Response codes:

- 200 for non-empty response
- 400 invalid request params
- 404 for no records found

**Response limit is capped at 500.**

### Create objects

```text
POST /mydata
```

This will persist/update a record in the store.

If UUID does not exist, create new record in store. If a record with id exists already, a new document will be created with a different UUID.

Request payload:

```json
{
"id": "uuid",
"attributes":[
    {
        "key1":"variable length string upto 256-chars",
        "value1":"variable length string upto 256-chars"
    },
     {
        "key2":"variable length string upto 256-chars",
        "value2":"variable length string upto 256-chars"
    }
],
"timestamp":"unix-epoch timestamp"
}
```

## Stack

1. GoLang
2. Gin/Gonic
3. GORM
4. Sqlite/Postgres

## Configuration Params

1. Number of worker threads
2. Backend conn string
