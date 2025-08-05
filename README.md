# Backend Go with MVC Architecture

A RESTful API backend built with Go, Gin router, and GORM ORM with PostgreSQL database.

## Features

- **MVC Architecture**: Clean separation of concerns with Models, Views (Controllers), and Routes
- **Gin Router**: Fast HTTP web framework
- **GORM**: Feature-rich ORM with auto-migration using standardized gorm.Model
- **PostgreSQL**: Robust relational database
- **CORS Support**: Cross-origin resource sharing enabled
- **Environment Variables**: Configuration through .env file
- **Standardized Models**: All models use gorm.Model for consistent ID, CreatedAt, UpdatedAt, and DeletedAt fields
- **Package-based Controllers**: Controllers organized in separate packages with consistent function names (GetAll, GetByID, Create, Update, Delete)

## Project Structure

```
backend-go/
├── config/
│   └── database.go           # Database configuration and connection
├── controllers/
│   ├── auth/
│   │   └── auth.go           # Authentication controllers (Register, Login, etc.)
│   ├── user/
│   │   └── user.go           # User controllers (GetAll, GetByID, Create, etc.)
│   └── trip/
│       └── trip.go           # Trip controllers (GetAll, GetByID, Create, etc.)
├── middleware/
│   └── auth.go               # Authentication and authorization middleware
├── models/
│   ├── models.go             # Auto-migration function
│   ├── user.go               # User model definition
│   ├── trip.go               # Trip model definition
│   └── image.go              # Image model definition
├── routes/
│   ├── router.go             # Main router setup
│   ├── auth/
│   │   └── auth_routes.go    # Authentication routes
│   ├── user/
│   │   └── user_routes.go    # User routes
│   └── trip/
│       └── trip_routes.go    # Trip routes
├── utils/
│   └── auth.go               # JWT and password utilities
├── .env                      # Environment variables
├── main.go                   # Application entry point
└── README.md
```

## Installation

1. **Clone the repository**

   ```bash
   git clone <repository-url>
   cd backend-go
   ```

2. **Install dependencies**

   ```bash
   go mod tidy
   ```

3. **Set up PostgreSQL database**

   - Install PostgreSQL
   - Create a database named `backend_go`
   - Update the `.env` file with your database credentials

4. **Configure environment variables**
   Copy and modify the `.env` file:
   ```env
   DB_HOST=localhost
   DB_PORT=5432
   DB_USER=postgres
   DB_PASSWORD=your_password
   DB_NAME=backend_go
   PORT=8080
   ```

## Running the Application

```bash
go run main.go
```

The server will start on `http://localhost:8080`

## API Endpoints

### Health Check

- `GET /health` - Check API status

### Authentication

- `POST /api/v1/auth/register` - Register new user
- `POST /api/v1/auth/login` - Login user
- `GET /api/v1/auth/profile` - Get current user profile (requires auth)
- `PUT /api/v1/auth/profile` - Update current user profile (requires auth)

### Users (Protected Routes)

- `GET /api/v1/users` - Get all users (admin only)
- `GET /api/v1/users/:id` - Get user by ID (authenticated users)
- `POST /api/v1/users` - Create new user (admin only)
- `PUT /api/v1/users/:id` - Update user (authenticated users)
- `DELETE /api/v1/users/:id` - Delete user (admin only)

### Trips

- `GET /api/v1/trips` - Get all trips (public)
- `GET /api/v1/trips/:id` - Get trip by ID (public)
- `POST /api/v1/trips` - Create new trip (requires auth)
- `PUT /api/v1/trips/:id` - Update trip (owner or admin only)
- `DELETE /api/v1/trips/:id` - Delete trip (owner or admin only)

## Authentication

This API uses JWT (JSON Web Token) for authentication. After successful login or registration, you'll receive a token that must be included in the Authorization header for protected routes.

### Authorization Header Format

```
Authorization: Bearer <your-jwt-token>
```

### User Roles

- **user**: Default role, can create and manage own trips
- **admin**: Can manage all users and trips

## Example API Usage

### Register a User

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"name": "John Doe", "email": "john@example.com", "password": "password123"}'
```

### Login

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "john@example.com", "password": "password123"}'
```

### Create a Trip (with authentication)

```bash
curl -X POST http://localhost:8080/api/v1/trips \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <your-jwt-token>" \
  -d '{"name": "Mountain Adventure", "description": "Hiking trip", "price": 299.99, "duration": 3, "start_latitude": 40.7128, "start_longitude": -74.0060, "end_latitude": 40.7589, "end_longitude": -73.9851}'
```

### Get User Profile

```bash
curl -X GET http://localhost:8080/api/v1/auth/profile \
  -H "Authorization: Bearer <your-jwt-token>"
```

-H "Content-Type: application/json" \
 -d '{"name": "Laptop", "description": "Gaming laptop", "price": 999.99, "stock": 10}'

````

## Dependencies

- **Gin**: Web framework
- **GORM**: ORM library
- **PostgreSQL Driver**: Database driver
- **GoDotEnv**: Environment variable loader

## Auto-Migration

The application automatically creates and updates database tables based on the model definitions when started. All models use `gorm.Model` which provides:

- `ID` (uint, primary key)
- `CreatedAt` (time.Time)
- `UpdatedAt` (time.Time)
- `DeletedAt` (gorm.DeletedAt, for soft deletes)

The migration process includes:

- Creating tables if they don't exist
- Adding new columns
- Creating indexes
- Setting up foreign key constraints

## Controller Package Structure

Controllers are organized into separate packages for better modularity and consistent naming. Each controller package provides standard CRUD operations with consistent function names:

### Standard Controller Functions:
- `GetAll(c *gin.Context)` - Retrieve all records
- `GetByID(c *gin.Context)` - Retrieve a single record by ID
- `Create(c *gin.Context)` - Create a new record
- `Update(c *gin.Context)` - Update an existing record
- `Delete(c *gin.Context)` - Delete a record

### Usage in Routes:
```go
import (
    "backend-go/controllers/user"
    "backend-go/controllers/trip"
)

// In route definitions:
userGroup.GET("/", user.GetAll)
tripGroup.GET("/", trip.GetAll)
```

This structure allows for consistent naming across different controllers while maintaining clear separation of concerns.

## Development

To add new models:

1. Create a new file in `models/` directory (e.g., `models/category.go`)
2. Define the model using `gorm.Model`:

   ```go
   package models

   import "gorm.io/gorm"

   type Category struct {
       gorm.Model
       Name string `json:"name" gorm:"not null"`
   }
````

3. Add the model to the `AutoMigrate()` function in `models/models.go`
4. Create controller functions in `controllers/`
5. Add routes in `routes/routes.go`
