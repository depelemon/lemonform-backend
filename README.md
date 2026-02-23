# 🍋 LemonForm Backend

The backend API for **LemonForm** — a form builder application (like Google Forms) with a lemon twist. Built with [Go](https://go.dev) and [Fiber](https://gofiber.io), connected to [Supabase](https://supabase.com) (PostgreSQL) as the database.

Authentication is handled entirely in the Go backend (JWT-based), using Supabase solely as a PostgreSQL database.

---

## 🚀 Features

- **Authentication** — Register & login with JWT-based auth (bcrypt password hashing)
- **Form CRUD** — Create, read, update, and delete forms (protected routes)
- Modular project architecture (handlers, middleware, config, platform)
- Environment configuration via `.env`
- Swagger API documentation
- Graceful shutdown
- CORS middleware for frontend integration

---

## 📁 Project Structure

```
lemonform-backend/
├── cmd/app/             # Application entrypoint, router, graceful shutdown
├── internal/
│   ├── common/          # Shared response structs
│   ├── config/          # Environment & stage config
│   └── test/            # Example test module (controller + routes)
├── platform/
│   └── database/        # PostgreSQL (Supabase) connection
├── sqlc/
│   └── postgresql/      # SQL schema & queries (for sqlc code generation)
├── docs/                # Swagger docs (auto-generated)
├── Makefile             # Build, run, lint, test commands
├── sqlc.yaml            # sqlc configuration
├── go.mod
└── go.sum
```

---

## 💠 Getting Started

### Prerequisites

- Go 1.23 or newer
- A [Supabase](https://supabase.com) project (for the PostgreSQL database)
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

Create a `.env` file in the project root:

```env
STAGE_STATUS=dev
SERVER_PORT=8080
DATABASE_TYPE=pgx
DATABASE_URL=postgresql://postgres:<password>@<host>:<port>/<database>
```

Set `DATABASE_URL` to your Supabase PostgreSQL connection string (found in Supabase → Settings → Database → Connection string → URI).

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
