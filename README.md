# 📱 AndroGate Server

Private outbound SMS transport server backend designed for **CONNECT ME** and the **AndroGate** Android client.

The server provides a high-performance REST API that enqueues outbound SMS verification messages, dispatches real-time event notifications to connected Android worker phones via Server-Sent Events (SSE), and tracks carrier delivery state.

---

## 📚 Table of Contents

- [About](#-about)
- [Architecture](#-architecture)
- [Features](#-features)
- [Prerequisites](#-prerequisites)
- [Quickstart with Docker](#-quickstart-with-docker)
- [Deployment on Render](#-deployment-on-render)
- [Configuration Reference](#-configuration-reference)
- [Authentication](#-authentication)
- [API Overview](#-api-overview)
- [License](#-license)

---

## 📖 About

AndroGate Server acts as the secure bridge between your core application backend (such as CONNECT ME) and dedicated Android devices running the AndroGate client app.

* **Outbound-only & Minimal**: Stripped of unnecessary legacy features (inbox harvesting, webhooks, Firebase/FCM cloud relays).
* **Direct Real-Time Push via SSE**: Communicates directly with the Android foreground service over persistent HTTP Server-Sent Events without external third-party push services.
* **Stateless & Scalable**: Compatible with managed cloud MySQL (Aiven, AWS RDS, PlanetScale) and container platforms like Render.

---

## 🏛 Architecture

```text
CONNECT ME Application / Backend
  │
  ├─ Enqueues SMS job: POST /api/3rdparty/v1/messages
  └─ Checks delivery state: GET /api/3rdparty/v1/messages/{id}
                │
                ▼
        AndroGate Server (Render)
    (MySQL Storage + Event Bus + SSE Broker)
                │
                │ Real-time Event: MessageEnqueued (over SSE)
                ▼
      AndroGate Android Handset
  ├─ Foreground service listening on /api/mobile/v1/events
  ├─ Fetches pending jobs: GET /api/mobile/v1/message?order=fifo
  ├─ Transmits SMS via Android Telephony APIs / active SIM
  └─ Reports state: PATCH /api/mobile/v1/message (Sent, Delivered, Failed)
```

---

## ⭐ Features

- **Authenticated Outbound SMS**: Supports text, data, and scheduled SMS dispatch.
- **Real-Time Push via SSE**: Immediate notification to connected devices using Server-Sent Events with automated keep-alive.
- **Granular Status Tracking**: Tracks `Pending`, `Processed`, `Sent`, `Delivered`, and `Failed` message lifecycles with recipient-level granularity.
- **Private Token Authentication**: Enforces `GATEWAY__PRIVATE_TOKEN` validation for device enrollment.
- **JWT & Basic Auth**: Secure API access with token refresh rotation and scoped authorization.
- **Container-Ready**: Root Dockerfile with automatic schema migrations and container health probes.
- **Swagger Documentation**: Interactive OpenAPI / Swagger UI served at `/api/docs/index.html`.

---

## 📦 Prerequisites

- **MySQL 8.0+** or **MariaDB 10.2+** (e.g., Aiven MySQL, local Docker MariaDB)
- **Docker** (for containerized deployment) or **Go 1.26+** (for local source builds)

---

## 🚀 Quickstart with Docker

### 1. Clone & Configure

Copy the environment template:

```bash
cp .env.example .env
```

Ensure your database connection and private tokens are configured in `.env`.

### 2. Run with Docker Compose

To start the server, background worker, and local MariaDB:

```bash
docker compose -f deployments/docker-compose/docker-compose.yml up --build
```

The server binds to `http://localhost:3000`.

---

## ☁️ Deployment on Render

This repository is optimized for one-click Docker web service deployment on [Render](https://render.com):

1. **Push this repository** to your personal GitHub account.
2. In the Render Dashboard, click **New +** → **Web Service**.
3. Select your repository.
4. Render will automatically detect the root [`Dockerfile`](Dockerfile).
5. Configure the service settings:
   - **Environment**: `Docker`
   - **Instance Type**: Starter / Standard (or Free tier)
6. Under **Environment Variables**, add the values from your `.env` file:
   - `GATEWAY__MODE` = `private`
   - `GATEWAY__PRIVATE_TOKEN` = `<your-secure-device-registration-token>`
   - `JWT__SECRET` = `<your-32-byte-hex-jwt-secret>`
   - `DATABASE__HOST` = `<your-mysql-host>`
   - `DATABASE__PORT` = `<your-mysql-port>`
   - `DATABASE__USER` = `<your-mysql-username>`
   - `DATABASE__PASSWORD` = `<your-mysql-password>`
   - `DATABASE__DATABASE` = `<your-database-name>`
   - `HTTP__LISTEN` = `0.0.0.0:3000`
   - `HTTP__OPENAPI__ENABLED` = `true`
7. Click **Deploy Web Service**.

> [!NOTE]
> DB migrations run automatically on container startup through the entrypoint script.
> Render automatically manages HTTPS certificates and routes external traffic to port `3000`.

---

## ⚙️ Configuration Reference

Every variable in `.env` can be supplied via environment variables:

| Variable | Default | Purpose |
| :--- | :--- | :--- |
| `GATEWAY__MODE` | `private` | Must be `private` for protected device registration. |
| `GATEWAY__PRIVATE_TOKEN` | *(Required)* | Secret token required by Android devices to register with the server. |
| `HTTP__LISTEN` | `0.0.0.0:3000` | Bind host and port for HTTP server. |
| `HTTP__API__HOST` | *(Empty)* | Public domain (leave empty to infer from incoming `Host` header). |
| `HTTP__OPENAPI__ENABLED` | `true` | Enables Swagger UI at `/api/docs/index.html`. |
| `DATABASE__HOST` | `localhost` | MySQL/MariaDB hostname. |
| `DATABASE__PORT` | `3306` | MySQL/MariaDB port. |
| `DATABASE__USER` | `sms` | Database username. |
| `DATABASE__PASSWORD` | `sms` | Database password. |
| `DATABASE__DATABASE` | `sms` | Database schema name. |
| `DATABASE__TIMEZONE` | `UTC` | Timezone for database timestamps and TTL calculations. |
| `JWT__SECRET` | *(Required for JWT)* | 32+ byte string used to sign JWT access and refresh tokens. |
| `JWT__ACCESS_TTL` | `15m` | Lifetime of access tokens. |
| `JWT__REFRESH_TTL` | `720h` | Lifetime of refresh tokens. |
| `CACHE__URL` | `memory://` | Cache provider (`memory://` or `redis://`). |
| `PUBSUB__URL` | `memory://` | PubSub provider for SSE events (`memory://` or `redis://`). |
| `SSE__KEEP_ALIVE_PERIOD_SECONDS` | `15` | Keep-alive ping interval for persistent SSE links. |

---

## 🔐 Authentication

### 1. Android Device Registration
The Android phone sends a registration request with your private token:

```http
POST /api/mobile/v1/device
Authorization: Bearer <GATEWAY__PRIVATE_TOKEN>
Content-Type: application/json

{
  "name": "Pixel 7 Pro",
  "simCards": [
    {
      "slotIndex": 0,
      "simNumber": 1,
      "carrierName": "Operator",
      "phoneNumber": "+1234567890"
    }
  ]
}
```

The server returns:
- Device `id`
- Device `token` (used for subsequent requests: `Authorization: Bearer <device_token>`)
- User credentials for direct API access

### 2. Third-Party API Authentication (Backend)
Authenticate using Basic Authentication or JWT Bearer tokens:

* **Issue JWT Token Pair**:
  ```http
  POST /api/3rdparty/v1/auth/token
  Authorization: Basic <base64(user:password)>
  ```
* **Refresh Token**:
  ```http
  POST /api/3rdparty/v1/auth/token/refresh
  Authorization: Bearer <refresh_token>
  ```

---

## 🔌 API Overview

### Third-Party API (for CONNECT ME Backend)

| Method | Endpoint | Description |
| :--- | :--- | :--- |
| `POST` | `/api/3rdparty/v1/messages` | Enqueue outbound SMS for dispatch |
| `GET` | `/api/3rdparty/v1/messages` | List sent and pending messages |
| `GET` | `/api/3rdparty/v1/messages/{id}` | Get status and recipient details of a message |
| `DELETE` | `/api/3rdparty/v1/messages/{id}` | Cancel a pending message |
| `GET` | `/api/3rdparty/v1/devices` | List registered Android devices |
| `DELETE` | `/api/3rdparty/v1/devices/{id}` | Unregister an Android device |
| `POST` | `/api/3rdparty/v1/auth/token` | Issue JWT access/refresh token pair |
| `GET` | `/api/3rdparty/v1/health` | Health status and database check |

### Mobile API (for AndroGate Android Client)

| Method | Endpoint | Description |
| :--- | :--- | :--- |
| `POST` | `/api/mobile/v1/device` | Register Android device (Bearer `GATEWAY__PRIVATE_TOKEN`) |
| `GET` | `/api/mobile/v1/device` | Get device connection info & external IP |
| `PATCH` | `/api/mobile/v1/device` | Update SIM cards or public encryption key |
| `GET` | `/api/mobile/v1/events` | **Server-Sent Events** real-time event stream |
| `GET` | `/api/mobile/v1/message` | Pull pending SMS jobs (`?order=fifo`) |
| `PATCH` | `/api/mobile/v1/message` | Report message state (`Sent`, `Delivered`, `Failed`) |
| `GET` | `/api/mobile/v1/settings` | Read device settings |

---

## 🛠 Local Development

```bash
make deps        # Download Go modules
make gen         # Regenerate Swagger OpenAPI docs
make test        # Run unit tests
make build       # Compile binary to tmp/sms-gateway
make air         # Live-reload development server
```

---

## ⚖️ License

Apache 2.0. See [LICENSE](LICENSE).
