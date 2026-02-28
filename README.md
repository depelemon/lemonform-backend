# 🍋 LemonForm Backend

[![Deploy to Koyeb](https://www.koyeb.com/static/images/deploy/button.svg)](https://app.koyeb.com/deploy?name=lemonform-backend&type=git&repository=depelemon%2Flemonform-backend&branch=main&regions=sin&env%5BSTAGE_STATUS%5D=dev&env%5BSERVER_PORT%5D=8080&env%5BDATABASE_URL%5D=postgresql%3A%2F%2Fpostgres.fakhpdpvdpyjswgnsxxc%3Alemonleminristekdatabase%40aws-1-ap-southeast-1.pooler.supabase.com%3A5432%2Fpostgres&env%5BJWT_SECRET%5D=pepapepethisisaverysecretkeylemonleminpepapepe2506601956depelemonlemonformrizztek&ports=8080%3Bhttp%3B%2F&hc_protocol%5B8080%5D=tcp&hc_grace_period%5B8080%5D=5&hc_interval%5B8080%5D=30&hc_restart_limit%5B8080%5D=3&hc_timeout%5B8080%5D=5&hc_path%5B8080%5D=%2F&hc_method%5B8080%5D=get)

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
```

Set `DATABASE_URL` to your PostgreSQL connection string (For example, in Supabase go to Settings → Database → Connection string → URI).

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
