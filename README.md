# Raspy

> **Your college schedule, without the routine of checking it.**

**Raspy** is a Telegram Mini App for students that provides a personalized college schedule based on their group and subgroup.

Instead of manually opening the college schedule every day, students can open Raspy directly from Telegram and instantly see their schedule for today, tomorrow, or the current week.

### 🚀 Try Raspy

**[Open Raspy in Telegram](https://t.me/@raspy_smk_bot)**


---

## Features

### 📅 Schedule

- Today's schedule
- Tomorrow's schedule
- Current week schedule
- Next week schedule
- Current and next lesson detection
- Lesson time, classroom, building and teacher information
- Subgroup-aware schedule filtering
- Automatic timezone handling

### 👤 Personalization

- Telegram-based authentication
- Automatic user creation
- College group search
- Group and subgroup selection
- Persistent user settings
- Notification preferences

### 📱 Telegram Mini App

- Runs directly inside Telegram
- No separate account registration
- Telegram `initData` authentication
- Mobile-first interface
- Schedule and settings available from one application

### ⚙️ Backend

- REST API built with Go
- Layered architecture
- PostgreSQL persistence
- Repository pattern
- Service layer
- External college API client
- Telegram authentication middleware
- Request validation and error handling

### 🤖 Telegram Bot

The bot acts as the entry point to the Mini App.

- `/start` — open the main menu
- `/schedule` — open the schedule
- `/help` — show available commands
- Telegram Web App integration

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

Raspy uses Telegram Mini App `initData` to authenticate users.

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

The backend is responsible for authentication and user identity. The frontend never sends a Telegram user ID as a trusted identity source.

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
│   ├── main.py
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

The backend follows a layered architecture:

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

This keeps HTTP handling, business logic, persistence, and external API integration separated.

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
- **Telegram Mini Apps API**

### Bot

- **Python**
- **aiogram 3**

### Infrastructure

- **Docker**
- **Docker Compose**
- **nginx**
- **HTTPS**
- **ngrok** for local Telegram Mini App development

---

## API

The backend exposes a REST API under:

```text
/api/v1
```

### Authentication

Authenticated requests use:

```http
X-Telegram-Init-Data: <telegram-init-data>
```

The backend validates Telegram's signed `initData` before processing the request.

### User

```http
GET /api/v1/me
```

Returns the currently authenticated user and their selected group/subgroup.

### Group search

```http
GET /api/v1/groups?q=<query>
```

Searches available college groups.

Example:

```http
GET /api/v1/groups?q=КИС-24
```

### Group selection

```http
PUT /api/v1/me/group
```

Updates the authenticated user's group and subgroup.

### Schedule

```http
GET /api/v1/schedule/today
GET /api/v1/schedule/tomorrow
GET /api/v1/schedule/week
GET /api/v1/schedule/week/next
```

Schedule requests use the authenticated user's selected group and subgroup.

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

A group is identified internally by its PostgreSQL ID, while `external_id` stores the identifier used by the college API.

---

## External College API

Raspy integrates with the college schedule API instead of maintaining a separate copy of the entire schedule.

The backend contains a dedicated client responsible for communication with the external service.

```text
ScheduleService
       │
       ▼
 College API Client
       │
       ▼
 College Schedule API
```

The service layer filters the received schedule according to the authenticated user's group and subgroup.

---

## Local Development

### Requirements

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
go run ./cmd/server
```

The backend runs on:

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

### Bot

```bash
cd bot
pip install -r requirements.txt
python main.py
```

The bot launches the Telegram Mini App and provides the main entry point for users.

---

## Environment Variables

The backend uses environment variables for configuration.

Example:

```env
DATABASE_URL=postgres://user:password@localhost:5432/raspy
BOT_TOKEN=your_bot_token
```

The bot uses:

```env
BOT_TOKEN=your_bot_token
WEB_APP_URL=https://your-mini-app-url
PROXY_URL=optional-proxy-url
```

Never commit real credentials, bot tokens, or database passwords to the repository.

---

## Production

Raspy is deployed as a containerized application using Docker Compose.

```text
                    Internet
                       │
                       ▼
                  ┌─────────┐
                  │  nginx  │
                  │ HTTPS   │
                  └────┬────┘
                       │
              ┌────────┴────────┐
              │                 │
              ▼                 ▼
        ┌───────────┐     ┌───────────┐
        │ Frontend  │     │  Backend  │
        │ React     │     │    Go     │
        └───────────┘     └─────┬─────┘
                                │
                                ▼
                         ┌────────────┐
                         │ PostgreSQL │
                         └────────────┘
```

The production environment includes:

- Docker Compose
- Go backend
- React frontend
- PostgreSQL
- nginx reverse proxy
- HTTPS
- Telegram Bot
- external college schedule API

---

## Development with Telegram

Telegram Mini Apps require an HTTPS URL when opened from Telegram.

For local development, the application can be exposed through a tunnel such as ngrok:

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

This allows testing the real Telegram authentication flow during development.

---

## Design Principles

### Separation of concerns

HTTP handlers should not contain business logic.

```text
Handler → Service → Repository
                 └→ External API
```

### Explicit dependencies

Services receive dependencies through constructors rather than creating them internally.

### Thin handlers

Handlers are responsible primarily for:

- reading requests
- validating input
- calling services
- returning HTTP responses

### External API isolation

College API-specific logic lives inside the `collegeapi` package instead of leaking into services and handlers.

### Persistent user state

The backend stores user-specific settings so the Mini App remains personalized between sessions.

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

### Phase 2 — Schedule intelligence

- [ ] Schedule snapshots
- [ ] Schedule comparison
- [ ] Automatic change detection
- [ ] Background synchronization
- [ ] Retry and error handling for background jobs

### Phase 3 — Notifications

- [ ] Telegram notifications
- [ ] Cancelled class notifications
- [ ] Added class notifications
- [ ] Classroom change notifications
- [ ] Teacher change notifications
- [ ] Time change notifications
- [ ] Notification preferences

### Phase 4 — Engineering

- [ ] Automated tests
- [ ] CI/CD
- [ ] Monitoring
- [ ] Production observability

---

## Why Raspy?

College schedules are often inconvenient to use:

- information is spread across different interfaces;
- students repeatedly check the same schedule;
- schedule changes can be easy to miss;
- existing interfaces are not always designed around a student's daily workflow.

Raspy focuses on a simple idea:

> **The student should not have to check the schedule. The schedule should come to the student.**

The Telegram Mini App is the interface. The long-term goal is a backend service that continuously tracks a student's schedule and proactively informs them about relevant changes.

---

## Project Status

Raspy is currently under active development.

The core personalized schedule experience is implemented and deployed:

- Telegram authentication
- group and subgroup selection
- personalized schedules
- today's, tomorrow's, and weekly schedules
- Telegram Mini App
- production Docker deployment
- HTTPS
- PostgreSQL persistence

The next major milestone is automatic schedule change detection and Telegram notifications.

---

## Author

**Даниил**

Backend-focused developer working with Go, Python and web technologies.

GitHub: [@mrDisa](https://github.com/mrDisa)

---

## License

This project is currently developed as a personal/educational project.
