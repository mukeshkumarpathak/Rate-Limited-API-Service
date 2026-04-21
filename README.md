# Rate-Limited-API-Service
Rate-Limited API Service Implementation

## Steps to Run the Project

1. Ensure Go is installed on your system (version 1.16 or later recommended).
2. Clone or navigate to the project directory.
3. Run `go mod tidy` to download dependencies (though this project has no external dependencies).
4. Run `go run main.go` to start the server.
5. The server will start on port 8080.

## API Endpoints

- `POST /request`: Accepts JSON payload `{"user_id": "string", "payload": "string"}`. Returns "Request accepted" if within rate limit, or 429 error if exceeded.
- `GET /stats`: Returns JSON with per-user request stats: `{"stats": {"user_id": count, ...}}`.

## Design Decisions

- **In-Memory Storage**: All data is stored in memory using Go maps and slices, as specified. No database is used for simplicity.
- **Rate Limiting**: Implemented using a sliding window approach with timestamps. For each user, we maintain a list of request timestamps and clean out entries older than 1 minute. This allows up to 5 requests per minute per user.
- **Thread Safety**: Used `sync.Mutex` to protect shared data structures from concurrent access by HTTP handlers.
- **Simple HTTP Server**: Used Go's standard `net/http` package for the server, keeping dependencies minimal.
- **JSON Handling**: Used `encoding/json` for request/response serialization.

## What You Would Improve with More Time

- **Persistence**: Add database support (e.g., Redis or PostgreSQL) for data persistence across restarts.
- **Configuration**: Make rate limits and port configurable via environment variables or config file.
- **Logging**: Add proper logging with levels (info, error) and structured logs.
- **Metrics**: Integrate metrics collection (e.g., Prometheus) for monitoring request rates and errors.
- **Testing**: Add comprehensive unit tests and integration tests.
- **Rate Limiting Algorithms**: Support different algorithms (fixed window, token bucket) and configurable limits per user or endpoint.
- **Authentication**: Add user authentication and authorization.
- **Error Handling**: More detailed error responses with error codes and messages.
- **Graceful Shutdown**: Handle server shutdown gracefully to finish processing ongoing requests.
- **Docker**: Containerize the application for easier deployment.
