# Raspy

> Your college schedule, without the routine of checking it.

Raspy is a personal college schedule service built around a Telegram Mini App.

Instead of manually checking the college schedule every day, users can open Raspy and immediately see their schedule for the current or next day. The service is designed to keep the schedule personalized to the user's group and subgroup and eventually notify users when the college schedule changes.

## Features

### Current

- Telegram Mini App authentication
- Automatic user creation on first authorization
- User's college group and subgroup stored in PostgreSQL
- Schedule retrieval from the college API
- Schedule filtering by subgroup
- Today's schedule
- Tomorrow's schedule
- Centralized schedule service layer
- PostgreSQL repositories for users and groups
- HTTP API built with Chi

### Planned

- Group and subgroup selection
- Change group and subgroup from the Mini App
- Current lesson and next lesson
- Break detection
- Weekly schedule
- Automatic schedule change detection
- Telegram notifications about:
  - cancelled classes
  - new classes
  - changed classrooms
  - changed teachers
  - changed time
  - changed subjects
- Notification settings
- Production deployment

## Architecture

Raspy follows a client-server architecture where the Go backend acts as the central application layer.

```text
                    ┌─────────────────────┐
                    │    Telegram Bot     │
                    │      aiogram        │
                    └──────────┬──────────┘
                               │
                               │
┌─────────────────┐            │
│ Telegram Mini   │            │
│      App        │            │
│ React + TS      │            │
└────────┬────────┘            │
         │                     │
         │ HTTP API             │
         ▼                     ▼
                ┌──────────────────────┐
                │      Go Backend      │
                │                      │
                │  Authentication      │
                │  Handlers            │
                │  Services            │
                │  Repositories        │
                │  College API client  │
                └───────┬───────┬──────┘
                        │       │
                        ▼       ▼
                 ┌──────────┐  ┌─────────────────┐
                 │PostgreSQL│  │   College API   │
                 └──────────┘  └─────────────────┘