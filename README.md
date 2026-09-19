# Messenger

A full-stack web messenger built with **Go, React, TypeScript and PostgreSQL**.

The project implements user authentication, email verification, sessions, chat rooms, room invitations and real-time messaging over WebSocket.

## Features

* User registration and login
* Email verification
* Session-based authentication
* Secure password hashing with bcrypt
* Chat room creation
* Room membership management
* Room invitations
* Sending and receiving messages
* Real-time message delivery using WebSocket
* PostgreSQL database
* Database migrations
* Backend tests
* Docker and Docker Compose support
* React + TypeScript frontend
* Serving the built frontend from the Go backend

## Tech Stack

### Backend

* Go
* `net/http`
* PostgreSQL
* `pgx`
* `golang-migrate`
* bcrypt
* WebSocket

### Frontend

* React
* TypeScript
* Vite
* CSS

### Infrastructure

* Docker
* Docker Compose
* PostgreSQL

## Project Structure

```text
.
├── backend/
│   ├── auth/
│   │   └── ...
│   │
│   ├── check/
│   │   ├── validate/
│   │   └── verify/
│   │
│   ├── db/
│   │   └── migrations/
│   │
│   ├── internal/
│   │   ├── apperrors/
│   │   └── testutils/
│   │
│   ├── models/
│   │
│   ├── server/
│   │   ├── invite/
│   │   ├── login/
│   │   ├── message/
│   │   ├── page_main/
│   │   ├── register/
│   │   ├── room/
│   │   ├── root/
│   │   ├── static/
│   │   └── websocket/
│   │
│   ├── sharing/
│   ├── Dockerfile
│   └── main.go
│
├── frontend/
│   ├── src/
│   ├── dist/
│   └── package.json
│
└── compose.yaml
```

## Architecture

The application consists of two main parts:

```text
                ┌─────────────────┐
                │     Browser     │
                │ React + TS      │
                └────────┬────────┘
                         │
                 HTTP / WebSocket
                         │
                         ▼
                ┌─────────────────┐
                │   Go Backend    │
                │    net/http     │
                └───────┬─────────┘
                        │
              ┌─────────┴─────────┐
              │                   │
              ▼                   ▼
       ┌─────────────┐     ┌─────────────┐
       │ PostgreSQL  │     │  WebSocket  │
       │             │     │     Hub     │
       └─────────────┘     └─────────────┘
```

The Go backend handles HTTP requests, authentication, database operations and WebSocket connections.

The frontend is developed separately with React and TypeScript. Its production build is placed in `frontend/dist` and can be served by the backend as static files.

## Authentication

Authentication is implemented using server-side sessions.

After successful authentication, a session is created and associated with the user. The session identifier is stored server-side and sent to the browser using an HTTP cookie.

Passwords are hashed with bcrypt and are never stored as plaintext.

### Registration and Email Verification

Registration uses a separate pending-user flow.

```text
Register
   │
   ▼
Pending user
   │
   ├── verification token
   │
   ▼
Email verification
   │
   ▼
Active user
```

Verification tokens are generated using a cryptographically secure random source. Only the token hash is stored by the server.

Expired pending registrations can be removed automatically.

## Rooms

Users can create chat rooms and become their members.

Room membership is represented separately from the room itself, allowing users to participate in multiple rooms.

The project also contains an invitation mechanism for sharing access to rooms.

The relevant functionality is separated into the `room` and `invite` server packages.

## Messaging

Messages are stored in PostgreSQL and associated with both their author and the room in which they were sent.

The HTTP API is responsible for operations such as loading existing messages, while WebSocket is used for real-time communication.

## WebSocket

Real-time communication is implemented with WebSocket.

The WebSocket functionality is separated into its own server package and includes a hub responsible for managing active connections.

The general message flow is:

```text
Client
   │
   │ WebSocket
   ▼
WebSocket handler
   │
   ▼
Hub
   │
   ├── connected clients
   │
   └── message delivery
```

The WebSocket upgrade request is validated before the connection is accepted.

## Database

PostgreSQL is used as the primary database.

Database changes are managed through migrations located in:

```text
backend/db/migrations/
```

The database contains separate entities for users, rooms, room membership, messages, pending users and sessions.

Migrations are applied when the backend starts.

## Frontend

The frontend is built with React and TypeScript using Vite.

Development source code is located in:

```text
frontend/src/
```

A production build can be created with:

```bash
npm run build
```

The resulting files are placed in:

```text
frontend/dist/
```

The built frontend is then made available to the Go backend as static files.

## Running with Docker Compose

The project contains a `compose.yaml` for running the application using Docker Compose.

Build the services:

```bash
docker compose build
```

Build and start the application:

```bash
docker compose up --build
```

Run in the background:

```bash
docker compose up --build -d
```

Stop the application:

```bash
docker compose down
```

## Running Locally

### Backend

From the backend directory:

```bash
cd backend
go run .
```

### Frontend

Install dependencies:

```bash
cd frontend
npm install
```

Start the development server:

```bash
npm run dev
```

Build the frontend:

```bash
npm run build
```

## Testing

Backend tests can be run with:

```bash
cd backend
go test ./...
```

The project contains tests for database operations and server functionality.

## Configuration

The backend uses environment variables for configuration, including the database connection and other application services.

Secrets and credentials should be provided through environment variables and should not be committed to the repository.

## Development

The project is built primarily with standard Go HTTP functionality rather than relying on a large backend framework.

The application is structured into separate packages for authentication, validation, verification, database access, HTTP handlers, WebSocket communication and room functionality.

The frontend uses React and TypeScript, while the backend remains responsible for the main application logic and data persistence.
