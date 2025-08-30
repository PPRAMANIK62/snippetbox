# Snippetbox

A secure web application for sharing code snippets, built with Go and MySQL. Users can create, view, and share code snippets with automatic expiration dates.

## Features

- **Snippet Management**: Create, view, and browse code snippets
- **User Authentication**: Secure user registration and login system
- **Session Management**: Session-based authentication with MySQL storage
- **Security Features**:
  - HTTPS/TLS encryption
  - CSRF protection
  - Secure password hashing with bcrypt
  - Input validation and sanitization
- **Responsive UI**: Clean, modern web interface
- **Automatic Expiration**: Snippets automatically expire after a set period

## Tech Stack

- **Backend**: Go 1.25
- **Database**: MySQL
- **Web Framework**: Custom HTTP handlers with httprouter
- **Session Store**: MySQL-based session management
- **Security**:
  - TLS/HTTPS
  - CSRF protection (nosurf)
  - Password hashing (bcrypt)
- **Templating**: Go's html/template
- **Middleware**: Custom middleware chain using Alice

## Dependencies

- `github.com/julienschmidt/httprouter` - HTTP request router
- `github.com/justinas/alice` - Middleware chaining
- `github.com/justinas/nosurf` - CSRF protection
- `github.com/alexedwards/scs/v2` - Session management
- `github.com/go-sql-driver/mysql` - MySQL driver
- `github.com/go-playground/form/v4` - Form data binding
- `golang.org/x/crypto` - Cryptography utilities

## Prerequisites

- Go 1.25 or higher
- MySQL 5.7 or higher
- TLS certificates (for HTTPS)

## Database Setup

Create a MySQL database and the following tables:

```sql
-- Snippets table
CREATE TABLE snippets (
    id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
    title VARCHAR(100) NOT NULL,
    content TEXT NOT NULL,
    created DATETIME NOT NULL,
    expires DATETIME NOT NULL
);

-- Users table
CREATE TABLE users (
    id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    hashed_password CHAR(60) NOT NULL,
    created DATETIME NOT NULL
);

-- Add unique constraint on email
ALTER TABLE users ADD CONSTRAINT users_uc_email UNIQUE (email);

-- Sessions table (for SCS session store)
CREATE TABLE sessions (
    token CHAR(43) PRIMARY KEY,
    data BLOB NOT NULL,
    expiry TIMESTAMP(6) NOT NULL
);

-- Add index on expiry for cleanup
CREATE INDEX sessions_expiry_idx ON sessions (expiry);
```

## Installation & Setup

1. **Clone the repository**
   ```bash
   git clone https://github.com/PPRAMANIK62/snippetbox.git
   cd snippetbox
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Set up environment variables**
   ```bash
   export MYSQL_DSN="username:password@/database_name?parseTime=true"
   ```

4. **Generate TLS certificates**
   ```bash
   mkdir tls
   cd tls
   # Generate self-signed certificate for development
   go run /usr/local/go/src/crypto/tls/generate_cert.go --rsa-bits=2048 --host=localhost
   # Rename generated files
   mv cert.pem cert.pem
   mv key.pem key.pem
   cd ..
   ```

5. **Run the application**
   ```bash
   go run ./cmd/web
   ```

The server will start on port 4000 by default. Visit `https://localhost:4000` to access the application.

## Usage

### Command Line Options

- `-addr`: HTTP network address (default: ":4000")

Example:
```bash
go run ./cmd/web -addr=":8080"
```

### Environment Variables

- `MYSQL_DSN`: MySQL Data Source Name for database connection

## Project Structure

```
snippetbox/
├── cmd/web/                 # Application entry point and web handlers
│   ├── main.go             # Main application setup
│   ├── handlers.go         # HTTP handlers
│   ├── middleware.go       # Custom middleware
│   ├── routes.go           # URL routing
│   ├── templates.go        # Template handling
│   └── helpers.go          # Helper functions
├── internal/
│   ├── models/             # Data models and database logic
│   │   ├── snippets.go     # Snippet model
│   │   └── users.go        # User model
│   ├── validator/          # Input validation
│   └── assert/             # Testing utilities
├── ui/
│   ├── html/               # HTML templates
│   │   ├── pages/          # Page templates
│   │   ├── partials/       # Partial templates
│   │   └── base.html       # Base template
│   ├── static/             # Static assets (CSS, JS, images)
│   └── efs.go              # Embedded file system
├── tls/                    # TLS certificates (gitignored)
├── go.mod                  # Go module definition
└── README.md               # This file
```

## API Endpoints

- `GET /` - Home page with latest snippets
- `GET /snippet/view/:id` - View a specific snippet
- `GET /snippet/create` - Create snippet form
- `POST /snippet/create` - Create new snippet
- `GET /user/signup` - User registration form
- `POST /user/signup` - Register new user
- `GET /user/login` - Login form
- `POST /user/login` - Authenticate user
- `POST /user/logout` - Logout user

## Security Features

- **HTTPS Only**: All traffic encrypted with TLS
- **CSRF Protection**: Cross-site request forgery protection on all forms
- **Secure Sessions**: HTTP-only, secure cookies with MySQL storage
- **Password Security**: Bcrypt hashing with cost factor 12
- **Input Validation**: Server-side validation and sanitization
- **SQL Injection Protection**: Prepared statements for all queries

## Testing

Run the test suite:

```bash
go test ./...
```

Run tests with coverage:

```bash
go test -cover ./...
```

## Development

To run in development mode with live reloading, you can use tools like `air`:

```bash
# Install air
go install github.com/cosmtrek/air@latest

# Run with air
air
```

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## Acknowledgments

- Built following best practices from Alex Edwards' "Let's Go" book
- Uses secure defaults and production-ready patterns
- Implements modern Go web development techniques
