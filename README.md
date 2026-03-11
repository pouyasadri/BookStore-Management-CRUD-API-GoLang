# BookStore Management CRUD API

A comprehensive RESTful API for managing a bookstore inventory with authors, books, and user authentication. Built with Go, GORM v2, and MySQL.

## Features

✨ **Core Features**
- ✅ Complete CRUD operations for books and authors
- ✅ Advanced search and filtering
- ✅ Pagination support
- ✅ Relational data modeling (Books & Authors)
- ✅ User authentication with JWT tokens
- ✅ Role-based access control via middleware

🛠️ **Technical Excellence**
- ✅ Production-grade error handling
- ✅ Comprehensive middleware (CORS, logging, panic recovery)
- ✅ Auto-generated Swagger/OpenAPI documentation
- ✅ Environment-based configuration (.env support)
- ✅ Clean code architecture (controllers, models, routes pattern)
- ✅ Unit tests with example coverage

🐳 **DevOps Ready**
- ✅ Multi-stage Dockerfile for optimized production builds
- ✅ Docker Compose for local development
- ✅ GitHub Actions CI/CD pipeline
- ✅ Makefile for common tasks

## Quick Start

### Prerequisites

- Go 1.25 or higher
- MySQL 8.0 or higher
- Docker & Docker Compose (optional, for containerized setup)
- Make (optional, for command shortcuts)

### Local Development Setup

1. **Clone the repository**
   ```bash
   git clone https://github.com/pouyasadri/BookStore-Management-CRUD-API-GoLang.git
   cd BookStore-Management-CRUD-API-GoLang
   ```

2. **Set up environment variables**
   ```bash
   cp .env.example .env
   # Edit .env with your database credentials
   ```

3. **Install dependencies**
   ```bash
   go mod download
   ```

4. **Run the application**
   ```bash
   make run
   # or without make:
   go run ./cmd/main
   ```

The API will be available at `http://localhost:8080`

### Docker Compose Setup (Recommended)

```bash
# Start all services (MySQL + API)
make docker-up

# View logs
make docker-logs

# Stop services
make docker-down
```

This automatically sets up MySQL and the API with proper configuration.

## API Documentation

### Interactive Swagger UI

Once the server is running, visit: **http://localhost:8080/swagger/**

### Health Check
```
GET /health
```

## Authentication

### Register New User
```bash
POST /auth/register
Content-Type: application/json

{
  "username": "john_doe",
  "email": "john@example.com",
  "password": "securepassword123"
}
```

Response:
```json
{
  "code": 201,
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "username": "john_doe",
      "email": "john@example.com"
    }
  }
}
```

### Login
```bash
POST /auth/login
Content-Type: application/json

{
  "username": "john_doe",
  "password": "securepassword123"
}
```

### Using JWT Token

All protected endpoints require the `Authorization` header:

```bash
Authorization: Bearer <your_jwt_token>
```

## API Endpoints

### Books (Protected Routes)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/book/` | Get all books with optional filtering and pagination |
| POST | `/book/` | Create a new book |
| GET | `/book/{bookId}` | Get a specific book |
| PUT | `/book/{bookId}` | Update a book |
| DELETE | `/book/{bookId}` | Delete a book |

**Query Parameters for GET /book/**
- `author` - Filter by author name (partial match)
- `publication` - Filter by publication (partial match)
- `page` - Page number (default: 1)
- `limit` - Results per page (default: 10, max: 100)

**Example:**
```bash
GET /book/?author=tolkien&publication=harper&page=1&limit=20
Authorization: Bearer <token>
```

### Authors (Protected Routes)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/author/` | Get all authors with pagination |
| POST | `/author/` | Create a new author |
| GET | `/author/{authorId}` | Get a specific author with books |
| PUT | `/author/{authorId}` | Update an author |
| DELETE | `/author/{authorId}` | Delete an author |

**Query Parameters for GET /author/**
- `name` - Filter by author name (partial match)
- `page` - Page number (default: 1)
- `limit` - Results per page (default: 10, max: 100)

## Request/Response Examples

### Create a Book
```bash
POST /book/
Authorization: Bearer <token>
Content-Type: application/json

{
  "name": "The Hobbit",
  "author": "J.R.R. Tolkien",
  "publication": "Harper & Brothers"
}
```

Response (201):
```json
{
  "code": 201,
  "data": {
    "id": 1,
    "name": "The Hobbit",
    "author": "J.R.R. Tolkien",
    "publication": "Harper & Brothers",
    "createdAt": 1710171600,
    "updatedAt": 1710171600
  }
}
```

### Error Response
```json
{
  "code": 400,
  "message": "Validation failed",
  "details": "Book name is required"
}
```

## Environment Variables

```env
# Database Configuration
DB_HOST=127.0.0.1          # MySQL host
DB_PORT=3306               # MySQL port
DB_USER=root               # MySQL username
DB_PASSWORD=               # MySQL password (empty for no password)
DB_NAME=bookstore          # Database name

# Server Configuration
SERVER_PORT=8080           # API port

# Security
JWT_SECRET=your_secret_key_here_change_in_production
```

## Project Structure

```
BookStore-Management-CRUD-API-GoLang/
├── cmd/main/
│   ├── main.go             # Application entry point
│   └── main                # Compiled binary
├── pkg/
│   ├── config/
│   │   └── app.go          # Database connection & config
│   ├── controllers/
│   │   ├── book-controller.go
│   │   ├── author-controller.go
│   │   └── auth-controller.go
│   ├── middleware/
│   │   ├── auth.go         # JWT authentication middleware
│   │   ├── cors.go         # CORS headers middleware
│   │   ├── logging.go      # Request logging middleware
│   │   └── recovery.go     # Panic recovery middleware
│   ├── models/
│   │   ├── book.go
│   │   ├── author.go
│   │   ├── user.go
│   │   └── models_test.go
│   ├── routes/
│   │   └── bookstore-routes.go
│   └── utils/
│       ├── response.go     # Response helpers
│       ├── util.go         # Request body parsing
│       ├── jwt.go          # JWT token generation/validation
│       └── jwt_test.go
├── docs/                   # Auto-generated Swagger docs
├── .env                    # Environment variables (git-ignored)
├── .env.example            # Example environment variables
├── .gitignore              # Git ignore rules
├── Dockerfile              # Multi-stage Docker build
├── docker-compose.yml      # Local development setup
├── Makefile                # Development tasks
├── go.mod                  # Go module definition
├── go.sum                  # Dependency checksums
└── README.md               # This file
```

## Development Commands

Using Make (recommended):

```bash
make build              # Build the application
make run                # Run locally
make test               # Run unit tests
make swagger            # Regenerate Swagger docs
make fmt                # Format code
make lint               # Run linter
make docker-build       # Build Docker image
make docker-up          # Start Docker containers
make docker-down        # Stop Docker containers
make clean              # Clean build artifacts
make help               # Show all available commands
```

Without Make:

```bash
# Build
go build -o cmd/main/main ./cmd/main

# Run
go run ./cmd/main

# Test
go test -v ./...

# Format
go fmt ./...

# Lint
go vet ./...

# Generate Swagger docs
go run github.com/swaggo/swag/cmd/swag@latest init -g cmd/main/main.go
```

## Testing

Run all tests:
```bash
make test
```

Run with coverage:
```bash
make test-verbose
```

## Docker Usage

### Build Image
```bash
docker build -t bookstore-api:latest .
```

### Run Container
```bash
docker run -p 8080:8080 \
  -e DB_HOST=mysql \
  -e DB_USER=root \
  -e DB_PASSWORD=root \
  -e DB_NAME=bookstore \
  bookstore-api:latest
```

### Docker Compose
```bash
# Start services
docker-compose up -d

# View logs
docker-compose logs -f

# Stop services
docker-compose down
```

## CI/CD Pipeline

GitHub Actions workflow runs automatically on push and PR:
- ✅ Go linter (vet)
- ✅ Unit tests with race detection
- ✅ Code coverage reporting
- ✅ Docker build verification

See `.github/workflows/ci.yml` for details.

## Middleware Stack

1. **RecoveryMiddleware** - Catches panics and returns 500 errors gracefully
2. **CORSMiddleware** - Enables CORS for cross-origin requests
3. **LoggingMiddleware** - Logs all HTTP requests with method, path, status, and duration
4. **JWTMiddleware** - (Protected routes only) Validates JWT tokens

## Error Handling

All errors follow a consistent JSON format:

```json
{
  "code": 400,
  "message": "Validation failed",
  "details": "Additional error information"
}
```

HTTP Status Codes:
- `200` - Success
- `201` - Created
- `400` - Bad Request
- `401` - Unauthorized
- `404` - Not Found
- `409` - Conflict
- `500` - Internal Server Error

## Security Considerations

⚠️ **Important for Production:**

1. **Change JWT Secret**
   ```env
   JWT_SECRET=your_very_secure_random_key_with_sufficient_entropy
   ```

2. **Use Strong Database Passwords**
   ```env
   DB_PASSWORD=strong_secure_password_not_empty
   ```

3. **Enable HTTPS** - Update docker-compose and routing for HTTPS

4. **Hide Sensitive Data** - Never commit `.env` files (they're in .gitignore)

5. **Rate Limiting** - Consider adding rate limiting middleware for production

6. **Input Validation** - Always validate and sanitize user input

## Performance Notes

- Pagination is enforced (max 100 items per page)
- Database queries use proper indexing
- GORM v2 provides optimized query execution
- Connection pooling is configured by default
- Middleware chain is optimized for minimal overhead

## Troubleshooting

### Database Connection Error
```
Failed to connect to database: dial tcp 127.0.0.1:3306: connect: connection refused
```

**Solution:** Make sure MySQL is running and .env variables are correct
```bash
# Using Docker
make docker-up

# Or check MySQL is running
mysql -u root -p
```

### Port Already in Use
```
listen tcp :8080: bind: address already in use
```

**Solution:** Change the port in .env
```env
SERVER_PORT=9090
```

### JWT Token Errors
Ensure the `Authorization` header is properly formatted:
```bash
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit changes (`git commit -m 'Add amazing feature'`)
4. Push to branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the Apache License 2.0 - see the LICENSE file for details.

## Contact & Support

For issues, questions, or suggestions, please open an issue on GitHub.

---

**Made with ❤️ by [Pouya Sadri](https://github.com/pouyasadri)**

Last updated: March 2026
