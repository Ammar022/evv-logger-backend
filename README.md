# Mini EVV Logger – Caregiver Shift Tracker

A comprehensive Electronic Visit Verification (EVV) system for tracking caregiver visits and managing care schedules.

## Features

- **Schedule Management**: View and manage caregiver schedules
- **Visit Tracking**: Log start and end times with geolocation
- **Task Management**: Track care activities and their completion status
- **Dashboard Statistics**: Overview of schedules, visits, and completion rates
- **Real-time Geolocation**: Browser-based location tracking for visit verification

## Tech Stack

- **Backend**: Go (Golang) with Gorilla Mux
- **Database**: PostgreSQL
- **API Documentation**: Swagger/OpenAPI
- **Containerization**: Docker & Docker Compose

## Quick Start

### Prerequisites

- Go 1.24 or higher
- PostgreSQL 15 or higher
- Docker & Docker Compose (optional)

### Local Development

1. **Clone the repository**
   ```bash
   git clone https://github.com/Ammar022/evv-logger-backend
   cd evv-logger-backend
   ```

2. **Set up environment variables**
   ```bash
   cp .env.example .env
   # Edit .env with your database credentials
   ```

3. **Install dependencies**
   ```bash
   go mod tidy
   ```

4. **Start PostgreSQL** (using Docker)
   ```bash
   docker-compose up -d db
   ```

5. **Run the application**
   ```bash
   go run cmd/api/main.go
   ```

6. **Access the API**
   - API Server: http://localhost:8080
   - Swagger Documentation: http://localhost:8080/swagger/

### Using Docker

```bash
# Build and run with Docker Compose
docker-compose up --build

# Or run just the database
docker-compose up -d db
```

## API Usage Examples

### Start a Visit
```bash
curl -X POST http://localhost:8080/schedules/1/start \
  -H "Content-Type: application/json" \
  -d '{
    "timestamp": "2023-05-15T09:00:00Z",
    "location": {
      "latitude": 40.7128,
      "longitude": -74.0060
    }
  }'
```

### Update a Task
```bash
curl -X POST http://localhost:8080/tasks/1/update \
  -H "Content-Type: application/json" \
  -d '{
    "status": "completed",
    "reason": "Task completed successfully"
  }'
```

### Project Structure
```
internal/
  api/          # HTTP handlers and routes
  middleware/   # CORS and other middleware
  models/       # Data structures and DTOs
  repo/         # Database operations
cmd/api/       # Main application entry
docs/          # Swagger documentation
```

### Testing
Run all tests with coverage:
```bash
go test -cover ./...
```

View Swagger docs after starting server:
```bash
open http://localhost:8080/swagger/
