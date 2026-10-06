# OTA Sign

OTA Sign is a forms control portal authenticated and authorized by Keycloak, and backed by DocuSeal for signing.

The intended architecture is:

```text
Keycloak
  -> OTA Sign backend (OIDC login)

OTA Sign backend
  -> OTA Sign frontend
  -> DocuSeal API/webhooks
```

Keycloak remains the identity and authorization source. DocuSeal remains the signing engine. OTA Sign is the control layer for UIC-scoped access, form visibility, submissions, commander signatures, status tracking, downloads, and commander access lifecycle.

## Repository Layout

```text
backend/                    Go API service
frontend/                   React + Vite dashboard
docs/                       Architecture and API notes
```

## MVP Goal

The first milestone is:

```text
User opens OTA Sign
-> OTA Sign validates Keycloak OIDC login
-> Keycloak client roles determine OTA Sign capabilities
-> Backend creates session
-> User lands on the right dashboard
```

After that, the DocuSeal integration can be added behind the backend API without exposing DocuSeal secrets to the browser.

## Backend Quick Start

Go is required to run the backend.

```bash
cd backend
cp .env.example .env
go run ./cmd/server
```

## Frontend Quick Start

```bash
cd frontend
npm install
npm run dev
```

By default the frontend expects the backend at `http://localhost:8080`.

## Production Images

Production is intended to run from prebuilt GHCR images with no Portainer-side
build step:

```text
ghcr.io/YOUR_GITHUB_ORG/otasign-backend:latest
ghcr.io/YOUR_GITHUB_ORG/otasign-frontend:latest
```

Use `docker-compose.prod.example.yml` as the Portainer stack starting point.
See `docs/operations.md` for deployment, backup, restore, monitoring, and
secret rotation procedures.
Use `.env.prod.example` as the production stack variable template.
The frontend image is runtime-configurable with:

```text
OTASIGN_API_BASE_URL
```

The backend image is a compiled Go binary and expects production env vars such
as `DATABASE_URL`, `FRONTEND_URL`, `KEYCLOAK_ISSUER_URL`,
`KEYCLOAK_CLIENT_SECRET`, `DOCUSEAL_API_KEY`, and webhook secrets. See
[`docs/keycloak-setup.md`](docs/keycloak-setup.md) for the complete Keycloak
and OTA Sign setup procedure.

## Notification Webhook Demo

After setting `NOTIFICATION_WEBHOOK_URL` in `backend/.env`, send sample n8n
payloads with:

```bash
node scripts/send-demo-notification.js
```

The script sends one `commander_signature_requested` payload and one
`submission_completed` payload. If `NOTIFICATION_WEBHOOK_SECRET` is set, it
signs requests with `X-OTA-Signature`.
