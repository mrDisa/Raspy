# Raspy

> **Your college schedule, without the routine of checking it.**

Raspy is a Telegram Mini App for students that provides a personalized college schedule based on their group and subgroup.

Instead of manually opening the college schedule every day, a student can open Raspy directly from Telegram and immediately see what classes are happening today, tomorrow, or during the current week.

The main goal of the project is to go one step further than a simple schedule viewer: **Raspy is designed to detect schedule changes and notify students about them automatically.**

---

## Features

### Schedule

- Today's schedule
- Tomorrow's schedule
- Current week schedule
- Next week schedule
- Current and next lesson detection
- Lesson time, classroom, building and teacher information
- Subgroup-aware schedule filtering
- Automatic timezone handling

### Personalization

- Telegram-based authentication
- Automatic user creation
- College group selection with search
- Subgroup selection
- Persistent user settings
- Notification preferences

### Telegram Mini App

- Runs directly inside Telegram
- No separate account registration
- Telegram `initData` authentication
- Mobile-first interface
- Schedule and settings available from one application

### Backend

- REST API built with Go
- Layered architecture
- PostgreSQL persistence
- Repository pattern
- Service layer
- External college API client
- Telegram authentication middleware
- Request validation and error handling

### Planned

- Automatic schedule change detection
- Telegram notifications about schedule changes:
  - cancelled classes
  - added classes
  - changed classrooms
  - changed teachers
  - changed time
  - changed subjects
- Background workers for schedule synchronization
- Production deployment
- Automated tests and CI/CD

---

## Architecture

Raspy follows a client-server architecture where the Go backend acts as the central application layer.

```text
                         ┌─────────────────────┐
                         │    Telegram Bot     │
                         │      aiogram        │
                         └──────────┬──────────┘
                                    │
                                    │ launches
                                    ▼
                         ┌─────────────────────┐
                         │   Telegram Mini App │
                         │    React + TypeScript│
                         └──────────┬──────────┘
                                    │
                                    │ HTTP / JSON
                                    ▼
                    ┌──────────────────────────────┐
                    │          Go Backend          │
                    │                              │
                    │  Authentication              │
                    │  HTTP Handlers               │
                    │  Services                    │
                    │  Repositories                │
                    │  College API Client          │
                    └──────────────┬───────────────┘
                                   │
                    ┌──────────────┴──────────────┐
                    │                             │
                    ▼                             ▼
             ┌─────────────┐             ┌─────────────────┐
             │ PostgreSQL  │             │   College API   │
             │             │             │                 │
             │ Users       │             │ Schedule        │
             │ Groups      │             │ Groups          │
             └─────────────┘             └─────────────────┘
```

### Authentication flow

Raspy uses Telegram Mini App `initData` to identify users.

```text
Telegram
   │
   │ initData
   ▼
Mini App
   │
   │ X-Telegram-Init-Data
   ▼
Go Backend
   │
   ├── Validate Telegram signature
   │
   ├── Extract Telegram user ID
   │
   └── Find or create user
          │
          ▼
      PostgreSQL
```

The backend remains responsible for authentication and user identity. The frontend never sends a Telegram user ID as a trusted identity source.

---

## Project Structure

```text
Raspy/
│
├── backend/
│   ├── cmd/
│   │   └── server/
│   │       └── main.go
│   │
│   ├── internal/
│   │   ├── auth/
│   │   ├── collegeapi/
│   │   ├── handler/
│   │   ├── model/
│   │   ├── repository/
│   │   └── service/
│   │
│   ├── migrations/
│   ├── go.mod
│   └── go.sum
│
├── bot/
│   └── ...
│
├── frontend/
│   ├── src/
│   │   ├── api/
│   │   ├── pages/
│   │   ├── types/
│   │   └── ...
│   ├── package.json
│   └── vite.config.ts
│
├── docker-compose.yml
└── README.md
```

The backend is intentionally separated into layers:

```text
HTTP Handler
     │
     ▼
  Service
     │
     ├──────────────► College API
     │
     ▼
 Repository
     │
     ▼
 PostgreSQL
```

This keeps HTTP, business logic, persistence and external API integration independent from each other.

---

## Tech Stack

### Backend

- **Go**
- **Chi** — HTTP router
- **PostgreSQL** — persistent storage
- **pgx** — PostgreSQL driver
- **net/http** — HTTP client/server primitives

### Frontend

- **React**
- **TypeScript**
- **Vite**
- Telegram Mini Apps API

### Bot

- **Python**
- **aiogram 3**

### Infrastructure

- Docker
- Docker Compose
- ngrok for local Telegram Mini App development

---

## API

The backend exposes a REST API under:

```text
/api/v1
```

### Authentication

Telegram authentication is performed using:

```http
X-Telegram-Init-Data: <telegram-init-data>
```

The backend validates the data before processing authenticated requests.

### User

```http
GET /api/v1/me
```

Returns the currently authenticated user and their selected group/subgroup.

### Group selection

```http
GET /api/v1/groups?q=<query>
```

Searches available college groups.

Example:

```http
GET /api/v1/groups?q=КИС-24
```

The backend requests matching groups from the college API and persists the selected groups in PostgreSQL.

### Schedule

```http
GET /api/v1/schedule/today
GET /api/v1/schedule/tomorrow
GET /api/v1/schedule/week
GET /api/v1/schedule/week/next
```

Schedule requests use the authenticated user's group and subgroup.

---

## Data Model

The core database entities are users and groups.

```text
┌──────────────┐
│    groups    │
├──────────────┤
│ id           │
│ external_id  │
│ name         │
└──────┬───────┘
       │
       │ 1:N
       │
┌──────▼───────┐
│    users     │
├──────────────┤
│ id           │
│ telegram_id  │
│ group_id     │
│ subgroup     │
│ notifications_enabled
│ created_at   │
│ updated_at   │
└──────────────┘
```

A group is identified internally by its PostgreSQL ID while `external_id` stores the identifier used by the college API.

---

## External College API

Raspy integrates with the college's schedule API instead of maintaining a separate copy of the entire schedule.

The backend contains a dedicated client responsible for communicating with the external service.

```text
ScheduleService
       │
       ▼
 CollegeAPI Client
       │
       ▼
 College Schedule API
```

This keeps external API details outside the business logic.

The backend can then transform and filter the received schedule according to the authenticated user's group and subgroup.

---

## Local Development

### Requirements

Make sure you have:

- Go
- Node.js
- npm
- PostgreSQL
- Python 3
- Docker / Docker Compose
- Telegram account

### Clone

```bash
git clone https://github.com/mrDisa/Raspy.git
cd Raspy
```

### Backend

```bash
cd backend
go mod download
go run ./...
```

The development server runs on:

```text
http://localhost:3000
```

### Frontend

```bash
cd frontend
npm install
npm run dev
```

By default, Vite runs on:

```text
http://localhost:5173
```

The frontend proxies `/api` requests to the Go backend during local development.

### Bot

```bash
cd bot
pip install -r requirements.txt
python main.py
```

The bot is responsible for launching the Telegram Mini App.

---

## Environment Variables

Create the required `.env` files locally.

Example backend configuration:

```env
DATABASE_URL=postgres://user:password@localhost:5432/raspy
TELEGRAM_BOT_TOKEN=your_bot_token
DEV_MODE=true
```

Never commit real credentials, bot tokens or database passwords to the repository.

---

## Development with Telegram

Telegram Mini Apps require an HTTPS URL when opened from Telegram.

For local development, the project can be exposed through a tunnel such as ngrok:

```text
Telegram
    │
    ▼
HTTPS ngrok URL
    │
    ▼
Vite :5173
    │
    │ /api
    ▼
Go :3000
```

This makes it possible to test the real Telegram authentication flow locally instead of relying only on mocked data.

---

## Design Principles

Raspy is being developed with several engineering principles in mind:

### Separation of concerns

HTTP handlers should not contain business logic.

```text
Handler → Service → Repository
                 └→ External API
```

### Explicit dependencies

Services receive their dependencies through constructors rather than creating them internally.

### Thin handlers

Handlers are responsible primarily for:

- reading requests
- validating input
- calling services
- returning HTTP responses

### External API isolation

College API-specific logic lives inside the `collegeapi` package instead of leaking into services and handlers.

### Persistent user state

The backend stores user-specific settings so the Mini App can remain personalized between sessions.

---

## Roadmap

### Phase 1 — Core schedule

- [x] Telegram Mini App
- [x] Telegram authentication
- [x] User creation
- [x] PostgreSQL integration
- [x] College API integration
- [x] Group persistence
- [x] Subgroup filtering
- [x] Today's schedule
- [x] Tomorrow's schedule
- [x] Current week
- [x] Next week
- [x] Group search
- [x] Group selection
- [x] Subgroup selection
- [x] Settings page

### Phase 2 — Notifications

- [ ] Schedule snapshots
- [ ] Schedule comparison
- [ ] Change detection
- [ ] Background synchronization
- [ ] Telegram notifications
- [ ] Notification preferences
- [ ] Retry/error handling for background jobs

### Phase 3 — Production

- [ ] Production deployment
- [ ] Dockerized production environment
- [ ] HTTPS
- [ ] CI/CD
- [ ] Automated tests
- [ ] Monitoring and logging

---

## Why Raspy?

College schedules are often inconvenient to use:

- information is spread across different interfaces;
- students repeatedly check the same schedule;
- schedule changes can be easy to miss;
- the interface is not necessarily designed around a student's daily workflow.

Raspy focuses on one simple idea:

> **The student should not have to check the schedule. The schedule should come to the student.**

The Telegram Mini App is only the interface. The long-term goal is a backend service that continuously tracks the student's schedule and proactively informs them about relevant changes.

---

## Project Status

Raspy is currently under active development.

The core schedule functionality and personalized user flow are implemented. The next major engineering milestone is automatic schedule change detection and Telegram notifications.

---

## Author

**Даниил**

Backend-focused developer working with Go, Python and web technologies.

GitHub: [@mrDisa](https://github.com/mrDisa)

---

## License

This project is currently developed as a personal/educational project.