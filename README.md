# 🍋 LemonForm Backend

[![Deploy to Koyeb](https://www.koyeb.com/static/images/deploy/button.svg)](https://app.koyeb.com/deploy?name=lemonform-backend&type=git&repository=depelemon%2Flemonform-backend&branch=main&regions=sin&env%5BSTAGE_STATUS%5D=dev&env%5BSERVER_PORT%5D=8080&env%5BDATABASE_URL%5D=postgresql%3A%2F%2Fpostgres%3AYOUR_PASSWORD%40YOUR_HOST%3A5432%2Fpostgres&env%5BJWT_SECRET%5D=your-secret-key&env%5BCORS_ORIGINS%5D=https%3A%2F%2Fyour-frontend.vercel.app&ports=8080%3Bhttp%3B%2F&hc_protocol%5B8080%5D=tcp&hc_grace_period%5B8080%5D=5&hc_interval%5B8080%5D=30&hc_restart_limit%5B8080%5D=3&hc_timeout%5B8080%5D=5&hc_path%5B8080%5D=%2F&hc_method%5B8080%5D=get)

The backend API for **LemonForm** — a form builder application (like Google Forms) with a lemon twist. Built with [Go](https://go.dev) and [Fiber](https://gofiber.io), connected to [Supabase](https://supabase.com) (PostgreSQL) as the database.

Authentication is handled entirely in the Go backend (JWT-based), using Supabase solely as a PostgreSQL database.

---

## 🚀 Features

- **Authentication** — Register & login with JWT-based auth (bcrypt password hashing)
- **Form CRUD** — Create, read, update, and delete forms (protected routes)
- **Question CRUD** — Add, update, and remove questions (short answer, radio, checkbox, dropdown)
- **Response submission** — Public endpoint for submitting form responses (no auth required)
- **Response listing** — View all responses for a form (owner only)
- **Public form access** — Fetch a form and its questions without auth (for respondents)
- **Search, filter & sort** — Search by title, filter by status/date range, sort by title or date
- **Pagination** — Configurable page size with total count metadata
- Modular project architecture (handlers, middleware, config, platform)
- Swagger API documentation (auto-generated from annotations)
- Configurable CORS origins via environment variable
- Graceful shutdown

---

## 📁 Project Structure

```
lemonform-backend/
├── cmd/app/             # Application entrypoint, router, graceful shutdown
├── internal/
│   ├── auth/            # Register, login, JWT middleware
│   ├── common/          # Shared response structs
│   ├── config/          # Environment & stage config
│   ├── form/            # Form & question CRUD (protected)
│   ├── models/          # GORM models (User, Form, Question, Response, Answer)
│   ├── response/        # Response submission & listing, public form access
│   └── test/            # Example test module
├── platform/
│   └── database/        # PostgreSQL (Supabase) connection
├── docs/                # Swagger docs (auto-generated)
├── sqlc/
│   └── postgresql/      # SQL schema (for reference)
├── Makefile             # Build, run, lint, test commands
├── go.mod
└── go.sum
```

---

## 💠 Getting Started

### Prerequisites

- Go 1.23 or newer
- A PostgreSQL Database (Local or [Supabase](https://supabase.com))
- [sqlc](https://sqlc.dev) (optional, for regenerating DB code from SQL)

### Installation

```bash
# Clone the repository
git clone https://github.com/depelemon/lemonform-backend.git
cd lemonform-backend

# Copy the environment file and fill in your values
cp .env.example .env

# Run the application
make run
# or directly:
go run cmd/app/*.go
```

### Environment Variables

Create a `.env` file in the project root based on the provided `.env.example` file:

```env
STAGE_STATUS=dev
SERVER_PORT=8080
DATABASE_URL=postgresql://postgres:<password>@<host>:<port>/<database>
JWT_SECRET=dogfood
CORS_ORIGINS=http://localhost:3000
```

| Variable | Description |
|---|---|
| `STAGE_STATUS` | `dev` or `prod` |
| `SERVER_PORT` | Port the server listens on (default `8080`) |
| `DATABASE_URL` | PostgreSQL connection string (e.g. from Supabase → Settings → Database → Connection string → URI) |
| `JWT_SECRET` | Secret key for signing JWT tokens |
| `CORS_ORIGINS` | Comma-separated allowed origins (default `http://localhost:3000`) |

---

## 🧪 Useful Commands

```bash
make run        # Run the application
make build      # Run tests & build binary
make test       # Lint, security check, and run tests
make lint       # Run golangci-lint
make gen.docs   # Regenerate Swagger docs
```

### Recommended Tools

- [Air](https://github.com/cosmtrek/air) — live reloading during development
- [sqlc](https://sqlc.dev) — type-safe SQL code generation
- [swag](https://github.com/swaggo/swag) — Swagger doc generation from annotations

---

## 📄 License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for more details.

---

## 🙏 Credits

Project structure and boilerplate based on [**go-fiber-template**](https://github.com/crlnravel/go-fiber-template) by [@crlnravel](https://github.com/crlnravel). Thank you for the excellent starter template!
