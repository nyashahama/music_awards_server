# 🎵 Music Awards Voting System

A robust, production-ready voting platform for music awards built with Go, PostgreSQL, and modern backend architecture patterns.

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](CONTRIBUTING.md)

## 📋 Table of Contents

- [Features](#features)
- [Tech Stack](#tech-stack)
- [Architecture](#architecture)
- [Getting Started](#getting-started)
- [API Documentation](#api-documentation)
- [Database Schema](#database-schema)
- [Development](#development)
- [Testing](#testing)
- [Deployment](#deployment)
- [Contributing](#contributing)
- [License](#license)

## ✨ Features

### User Management
- 🔐 Secure user registration and authentication with JWT
- 👤 User profile management
- 🔑 Password reset functionality
- 📧 Email notifications (welcome emails, login alerts)
- 👑 Role-based access control (User/Admin)

### Voting System
- 🗳️ Vote for nominees in different categories
- ♻️ Change votes before voting period closes
- 📊 Vote tracking and management
- 🎯 One vote per category per user
- 🔢 Available votes quota system (3 votes per user)

### Category Management
- 📁 Create and manage award categories
- 📝 Category descriptions
- 🏆 Active categories with votes tracking

### Nominee Management
- 🎤 Add artists/nominees to categories
- 🖼️ Nominee profiles with images
- 🎵 Sample works (JSONB format)
- 🏷️ Multi-category support for nominees

## 🛠️ Tech Stack

### Backend
- **Language:** Go 1.21+
- **Web Framework:** Gin
- **Database:** PostgreSQL 15+
- **ORM:** GORM
- **Authentication:** JWT (golang-jwt/jwt)
- **Password Hashing:** bcrypt

### Architecture Patterns
- Clean Architecture (Handlers → Services → Repositories → Models)
- Repository Pattern
- Dependency Injection
- Interface-based design

## 🏗️ Architecture

```
music-awards/
├── cmd/
│   └── api/
│       └── main.go                 # Application entry point
├── internal/
│   ├── models/                     # Database models
│   │   ├── user.go
│   │   ├── category.go
│   │   ├── nominee.go
│   │   ├── vote.go
│   │   └── nominee_category.go
│   ├── repositories/               # Data access layer
│   │   ├── user_repository.go
│   │   ├── category_repository.go
│   │   ├── nominee_repository.go
│   │   ├── vote_repository.go
│   │   └── nominee_category_repository.go
│   ├── services/                   # Business logic layer
│   │   ├── user_service.go
│   │   ├── category_service.go
│   │   ├── nominee_service.go
│   │   ├── vote_service.go
│   │   └── nominee_category_service.go
│   ├── handlers/                   # HTTP handlers (Controllers)
│   │   ├── user.go
│   │   ├── category.go
│   │   ├── nominee.go
│   │   ├── vote.go
│   │   └── nominee_category.go
│   ├── dtos/                       # Data Transfer Objects
│   │   ├── user_dto.go
│   │   ├── category_dto.go
│   │   ├── nominee_dto.go
│   │   └── vote_dto.go
│   ├── middleware/                 # HTTP middleware
│   │   ├── auth.go
│   │   └── admin.go
│   ├── security/                   # Security utilities
│   │   └── jwt.go
│   └── validation/                 # Input validation
│       └── validators.go
├── migrations/                     # Database migrations
├── config/                         # Configuration files
├── .env.example                    # Environment variables template
├── go.mod
├── go.sum
└── README.md
```

## 🚀 Getting Started

### Prerequisites

- Go 1.21 or higher
- PostgreSQL 15 or higher
- Git

### Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/nyashahama/music_awards_server.git
   cd music-awards
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Set up PostgreSQL database**
   ```bash
   createdb music_awards
   ```

4. **Configure environment variables**
   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

5. **Run database migrations**
   ```bash
   # Using your migration tool
   # or let GORM auto-migrate on first run
   ```

6. **Run the application**
   ```bash
   go run cmd/api/main.go
   ```

The API will be available at `http://localhost:8080`

### Environment Variables

Create a `.env` file in the root directory:

```env
# Server Configuration
SERVER_PORT=8080
SERVER_ENVIRONMENT=development

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=music_awards
DB_SSL_MODE=disable

# JWT Configuration
JWT_SECRET=your_super_secret_jwt_key_change_this_in_production
JWT_EXPIRATION=24h

# Email Configuration (optional)
EMAIL_FROM=noreply@musicawards.com
EMAIL_HOST=smtp.gmail.com
EMAIL_PORT=587
EMAIL_USERNAME=your_email@gmail.com
EMAIL_PASSWORD=your_app_password

# CORS
ALLOWED_ORIGINS=http://localhost:3000
```

## 📚 API Documentation

### Authentication Endpoints

#### Register User
```http
POST /auth/register
Content-Type: application/json

{
  "first_name": "John",
  "last_name": "Doe",
  "email": "john@example.com",
  "password": "SecurePass123!",
  "location": "New York, USA"
}
```

#### Login
```http
POST /auth/login
Content-Type: application/json

{
  "email": "john@example.com",
  "password": "SecurePass123!"
}

Response:
{
  "token": "eyJhbGciOiJIUzI1NiIs..."
}
```

#### Forgot Password
```http
POST /auth/forgot-password
Content-Type: application/json

{
  "email": "john@example.com"
}
```

#### Reset Password
```http
POST /auth/reset-password
Content-Type: application/json

{
  "token": "reset_token_here",
  "new_password": "NewSecurePass123!"
}
```

### User Endpoints

All user endpoints require authentication (JWT token in Authorization header).

```http
Authorization: Bearer {your_jwt_token}
```

#### Get User Profile
```http
GET /users/:id
```

#### Update User Profile
```http
PUT /users/:id
Content-Type: application/json

{
  "first_name": "Jane",
  "last_name": "Smith",
  "location": "Los Angeles, USA"
}
```

#### List All Users (Admin only)
```http
GET /users
```

#### Promote User to Admin (Admin only)
```http
POST /users/:id/promote
```

### Category Endpoints

#### List All Categories (Public)
```http
GET /categories
```

#### Get Category Details (Public)
```http
GET /categories/:categoryId
```

#### Create Category (Admin only)
```http
POST /categories
Content-Type: application/json
Authorization: Bearer {admin_token}

{
  "name": "Best New Artist",
  "description": "Recognizing exceptional emerging talent"
}
```

#### Update Category (Admin only)
```http
PUT /categories/:categoryId
Content-Type: application/json
Authorization: Bearer {admin_token}

{
  "name": "Best New Artist 2024",
  "description": "Updated description"
}
```

#### Delete Category (Admin only)
```http
DELETE /categories/:categoryId
Authorization: Bearer {admin_token}
```

### Nominee Endpoints

#### List All Nominees (Public)
```http
GET /nominees
```

#### Get Nominee Details (Public)
```http
GET /nominees/:id
```

#### Create Nominee (Admin only)
```http
POST /nominees
Content-Type: application/json
Authorization: Bearer {admin_token}

{
  "name": "Artist Name",
  "description": "Talented musician from...",
  "image_url": "https://example.com/artist.jpg",
  "sample_works": ["Song 1", "Song 2", "Album 3"],
  "category_ids": ["uuid1", "uuid2"]
}
```

#### Update Nominee (Admin only)
```http
PUT /nominees/:id
Content-Type: application/json
Authorization: Bearer {admin_token}

{
  "name": "Updated Artist Name",
  "description": "Updated bio"
}
```

### Vote Endpoints

All vote endpoints require authentication.

#### Cast Vote
```http
POST /votes
Content-Type: application/json
Authorization: Bearer {token}

{
  "category_id": "category-uuid",
  "nominee_id": "nominee-uuid"
}
```

#### Get User's Votes
```http
GET /votes
Authorization: Bearer {token}
```

#### Get Available Votes
```http
GET /votes/available
Authorization: Bearer {token}

Response:
{
  "available_votes": 3
}
```

#### Change Vote
```http
PUT /votes/:id
Content-Type: application/json
Authorization: Bearer {token}

{
  "nominee_id": "new-nominee-uuid"
}
```

#### Delete Vote
```http
DELETE /votes/:id
Authorization: Bearer {token}
```

#### Get Category Votes (Admin only)
```http
GET /votes/category/:category_id
Authorization: Bearer {admin_token}
```

#### Get All Votes (Admin only)
```http
GET /votes/all
Authorization: Bearer {admin_token}
```

## 🗄️ Database Schema

### Users Table
```sql
CREATE TABLE users (
    user_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL,
    location VARCHAR(255),
    available_votes INTEGER NOT NULL DEFAULT 5,
    reset_token VARCHAR(255),
    reset_token_expires_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Categories Table
```sql
CREATE TABLE categories (
    category_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) UNIQUE NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Nominees Table
```sql
CREATE TABLE nominees (
    nominee_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    sample_works JSONB,
    image_url VARCHAR(500),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Votes Table
```sql
CREATE TABLE votes (
    vote_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    category_id UUID NOT NULL REFERENCES categories(category_id) ON DELETE CASCADE,
    nominee_id UUID NOT NULL REFERENCES nominees(nominee_id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, category_id)
);
```

### Nominee Categories (Join Table)
```sql
CREATE TABLE nominee_categories (
    nominee_id UUID REFERENCES nominees(nominee_id) ON DELETE CASCADE,
    category_id UUID REFERENCES categories(category_id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (nominee_id, category_id)
);
```

## 💻 Development

### Project Structure

The project follows **Clean Architecture** principles:

- **Models**: Database entities and business domain models
- **Repositories**: Data access layer with database operations
- **Services**: Business logic and domain rules
- **Handlers**: HTTP request handlers and response formatting
- **DTOs**: Data transfer objects for API requests/responses
- **Middleware**: Request interceptors (auth, logging, etc.)

### Code Style

This project follows Go best practices:
- Use `gofmt` for code formatting
- Follow effective Go guidelines
- Use meaningful variable names
- Write idiomatic Go code

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with verbose output
go test -v ./...
```

### Database Migrations

```bash
# Create new migration
migrate create -ext sql -dir migrations -seq migration_name

# Run migrations
migrate -path migrations -database "postgresql://user:pass@localhost:5432/music_awards?sslmode=disable" up

# Rollback migrations
migrate -path migrations -database "postgresql://user:pass@localhost:5432/music_awards?sslmode=disable" down
```

## 🧪 Testing

### Unit Tests

```bash
go test ./internal/services/...
go test ./internal/repositories/...
```

### Integration Tests

```bash
go test ./internal/handlers/... -tags=integration
```

## 🚢 Deployment

### Using Docker

1. **Build Docker image**
   ```bash
   docker build -t music-awards:latest .
   ```

2. **Run with Docker Compose**
   ```bash
   docker-compose up -d
   ```

### Manual Deployment

1. **Build the application**
   ```bash
   go build -o music-awards cmd/api/main.go
   ```

2. **Run the binary**
   ```bash
   ./music-awards
   ```

### Production Checklist

- [ ] Set strong JWT secret
- [ ] Enable HTTPS/TLS
- [ ] Configure CORS properly
- [ ] Set up database backups
- [ ] Enable rate limiting
- [ ] Configure logging
- [ ] Set up monitoring
- [ ] Use environment-specific configs
- [ ] Enable SSL for database connection
- [ ] Review and update security headers

## 🔒 Security Features

- **JWT Authentication**: Secure token-based authentication
- **Password Hashing**: bcrypt with cost factor 10
- **Input Validation**: Comprehensive request validation
- **SQL Injection Prevention**: GORM parameterized queries
- **CORS Protection**: Configurable CORS middleware
- **Email Enumeration Prevention**: Generic responses for password reset
- **Role-Based Access Control**: Admin and user roles
- **Cascade Deletes**: Database-level referential integrity

## 🤝 Contributing

Contributions are welcome! Please follow these steps:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

Please ensure your code:
- Follows Go best practices
- Includes appropriate tests
- Has clear commit messages
- Updates documentation as needed

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 👥 Authors

- **Your Name** - *Initial work* - [YourGitHub](https://github.com/nyashahama)

## 🙏 Acknowledgments

- Gin Web Framework
- GORM ORM
- PostgreSQL
- JWT-Go
- All contributors and supporters

## 📞 Support

For support, email nyashahama55@gmail.com or open an issue in the GitHub repository.

---

**Note**: This is a production-ready application, but always review and adjust security settings, environment variables, and configurations based on your specific deployment requirements.

Made by Nyasha Hama
