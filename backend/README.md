# Gemini API Wrapper - Refactored Backend

A clean, well-organized backend implementation for wrapping Google Gemini API with enhanced structure and beginner-friendly organization.

## Project Structure

```
backend/
├── cmd/
│   └── api/
│       └── main.go          # Application entry point
├── internal/                # Private application code
│   ├── handlers/            # HTTP handlers for each endpoint
│   ├── services/            # Business logic (Gemini API integration)
│   ├── models/              # Data structures and request/response models
│   ├── middleware/          # Custom middleware (CORS, logging)
├── go.mod                   # Go module definition
├── go.sum                   # Dependency checksums
├── Makefile                 # Build and development commands
├── README.md                # Documentation
├── .env.example             # Example environment file
├── .gitignore               # Git ignore patterns
├── embedding-request.json   # Example embedding request
├── request.json             # Example chat completion request
└── main_test.go             # Basic tests
```

## Features

- **Clean Architecture**: Well-organized codebase following Go best practices
- **Beginner-Friendly**: Clear structure and documentation for new developers
- **Secure**: API keys are stored in environment variables
- **Scalable**: Easy to extend with new endpoints and features
- **Well-Documented**: Comprehensive documentation and examples
- **Production Ready**: Includes CORS support, structured logging, and health checks

## Getting Started

### Prerequisites

- Go 1.19 or higher
- A Google Gemini API key

### Installation

1. Clone the repository:
   ```bash
   git clone <repository-url>
   cd backend
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

3. Set up environment variables:
   Copy the example environment file and add your Gemini API key:
   ```bash
   cp .env.example .env
   ```
   
   Edit the `.env` file and add your API key:
   ```
   GEMINI_API_KEY=your_actual_api_key_here
   PORT=8080
   ```

### Building the Application

You can build the application using either the Makefile or Go commands:

**Using Makefile:**
```bash
make build
```

**Using Go commands:**
```bash
go build -o api cmd/api/main.go
```

### Running the Application

**Using Makefile:**
```bash
make run
```

**Using Go commands:**
```bash
go run cmd/api/main.go
```

**Using the compiled binary:**
```bash
./api
```

## API Endpoints

### Health Check
- **Endpoint**: `GET /health`
- **Description**: Check if the server is running
- **Response**: 
  ```json
  {
    "status": "ok",
    "message": "Gemini API Wrapper is running"
  }
  ```

### Chat Completions
- **Endpoint**: `POST /v1/chat/completions`
- **Description**: Generate text completions using Gemini models
- **Request Body**:
  ```json
  {
    "model": "gemini-1.5-flash-latest",
    "messages": [
      {
        "role": "user",
        "content": "Hello, how are you?"
      }
    ],
    "stream": true,
    "temperature": 0.8
  }
  ```
- **Response**: Raw response from Google Gemini API

### Embeddings
- **Endpoint**: `POST /v1/embeddings`
- **Description**: Generate embeddings for text input
- **Request Body**:
  ```json
  {
    "model": "models/gemini-embedding-001",
    "content": {
      "parts": [
        {
          "text": "Hello, how are you?"
        }
      ]
    }
  }
  ```
- **Response**: Raw embedding vector from Google Gemini API

## Example Requests

### Chat Completions (Non-Streaming)
```bash
curl -X POST http://localhost:8080/v1/chat/completions \
  -H "Content-Type: application/json" \
  -d '{
    "model": "gemini-1.5-flash-latest",
    "messages": [
      {
        "role": "user",
        "content": "What is Go programming language?"
      }
    ],
    "stream": false
  }'
```

### Chat Completions (Streaming)
```bash
curl -X POST http://localhost:8080/v1/chat/completions \
  -H "Content-Type: application/json" \
  -d '{
    "model": "gemini-1.5-flash-latest",
    "messages": [
      {
        "role": "user",
        "content": "Tell me a joke."
      }
    ],
    "stream": true
  }'
```

### Embeddings
```bash
curl -X POST http://localhost:8080/v1/embeddings \
  -H "Content-Type: application/json" \
  -d '{
    "model": "models/gemini-embedding-001",
    "content": {
      "parts": [
        {
          "text": "Hello, world!"
        }
      ]
    }
  }'
```

## Makefile Commands

The project includes a Makefile with useful commands:

- `make build` - Build the application
- `make run` - Run the application
- `make deps` - Install dependencies
- `make clean` - Clean build artifacts
- `make test` - Run tests
- `make coverage` - Run tests with coverage
- `make fmt` - Format code
- `make vet` - Vet code
- `make help` - Show all available commands

## Project Structure Details

### cmd/api/main.go
The main entry point of the application that initializes the server, loads environment variables, and sets up routes.

### internal/handlers/
Contains HTTP handlers for each endpoint:
- `chat_handler.go` - Handles chat completions requests
- `embedding_handler.go` - Handles embeddings requests

### internal/services/
Contains business logic for interacting with the Gemini API:
- `gemini_service.go` - Service for calling Gemini API endpoints

### internal/models/
Contains data structures used throughout the application:
- `models.go` - Request/response models for both chat and embeddings

### internal/middleware/
Contains custom middleware functions:
- `middleware.go` - CORS and logging middleware

## Development

### Adding New Endpoints

1. Create a new handler in `internal/handlers/`
2. Add the route in `cmd/api/main.go`
3. Create corresponding models in `internal/models/` if needed
4. Update the README with documentation for the new endpoint

### Testing

Run tests with:
```bash
make test
```

Or with coverage:
```bash
make coverage
```

## Environment Variables

- `GEMINI_API_KEY` - Your Google Gemini API key (required)
- `PORT` - Port to run the server on (default: 8080)
- `ENV` - Environment mode (development, production)

## Dependencies

- [Gin](https://github.com/gin-gonic/gin) - HTTP web framework
- [godotenv](https://github.com/joho/godotenv) - Environment variable loading

## Contributing

1. Fork the repository
2. Create a new branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a pull request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgments

- Google for providing the Gemini API
- The Gin framework for the excellent HTTP handling