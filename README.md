# Asgard: Guard AI Server

Handles all Guard AI traffic between the guard-ai-client and guard-ai-worker (Hiemdall).
Serves several endpoints which interacts with a PostgresDB instance with 3 main records defined in models.
Asgard bridges worker ingested data and broadcasts notifications when needed.

## Setup

- `go mod tidy`
- Create a `.env` file populated by `.sample.env` 
- `go run main.go`

# Endpoints

## Events

All routes related to events

* GET `/events/near/:uuid`
    - Get's all events near a user identified by the URL param `uuid`
    - Response:
     ```
    events: [Event]
     ```

## User

* POST `/user/ping`
    - Updates the user's location
    - Request:
     ```
     id: UUID,
     location: "Lat,Lon"
     ```

* POST `/user`
    - Creates a new user to track their location to send notifications
    - The user id must be present in the `auth` table managed by Supabase
    - Request:
    ```
    Id: UUID
    location: "Lat,Lon"
    ```

## Worker

* POST `/worker/record`
    - Ingests all new `Logs` and `Events` and sends push notifications to users when appropriate
    - Requires a bearer token only distributed to private workers
    - Request:
     ```
     Logs: [Log]
     Events: [Event]
     ```

# Models

## Log

```go
type Log struct {
	Id        uuid.UUID          `json:"id"`
	Region    string             `json:"region"`
	Utterance string             `json:"utterance"`
	CreatedAt pgtype.Timestamptz `json:"created_at"`
}
```

## Event

```go
type Event struct {
	Id        uuid.UUID          `json:"id"`
	Level     string             `json:"level"`
	Location  Point              `json:"location"`
	Category  string             `json:"category"`
	LogId     uuid.UUID          `json:"log_id"`
	CreatedAt pgtype.Timestamptz `json:"created_at"`
}
```

## User

```go
type User struct {
	Id        uuid.UUID          `json:"id"`
	LastPing  pgtype.Timestamptz `json:"last_ping"`
	CreatedAt pgtype.Timestamptz `json:"created_at"`
	Location  Point              `json:"location"`
}
```
