# OTA Sign Architecture

## System Roles

```text
Keycloak
- Login and identity
- DoD ID, UIC, rank, name, and email claims

Moodle
- OTA Sign roles/capabilities
- Optional launch point into OTA Sign

OTA Sign Backend
- Validates Keycloak OIDC tokens
- Asks Moodle connector for OTA Sign capabilities using the DoD ID
- Owns sessions and permissions
- Stores users, UIC roles, templates, submissions, signer state, and audit events
- Talks to DocuSeal API
- Receives DocuSeal webhooks
- Sends SMTP emails

OTA Sign Frontend
- Student "My Forms" dashboard
- Commander "Unit Forms" dashboard
- Start/view/sign/download actions through backend API only

DocuSeal
- Signing ceremony
- Signer authentication flow
- Audit trail
- Completed signed PDFs
```

## Launch Flow

```text
1. User opens OTA Sign directly or from the Moodle link.
2. OTA Sign redirects the user to Keycloak.
3. Keycloak returns a signed ID token containing the user identity attributes.
4. OTA Sign validates issuer, signature, audience, expiry, state, nonce, and PKCE.
5. OTA Sign sends the validated DoD ID to the Moodle connector service.
6. Moodle returns only the evaluated OTA Sign capabilities.
7. Backend creates a secure session.
8. Frontend loads /api/me and dashboard data from the backend.
```

## Security Rules

- The frontend never receives Moodle signing secrets, DocuSeal API keys, SMTP credentials, or database credentials.
- Every frontend request goes to the OTA Sign backend.
- Backend permission checks are UIC-scoped.
- A commander or commander representative can only see users and submissions for authorized UICs.
- A user can see their own submissions even when they also have commander permissions.
- DocuSeal webhooks must be verified before updating local state.

## Submission Status

OTA Sign should normalize DocuSeal status into these user-facing states:

```text
missing
pending
complete
canceled
failed
```

For commander dashboards, the backend should calculate the most current submission per:

```text
soldier_id + template_id
```

If no submission exists for that soldier/template, return `missing`.

## Commander Dashboard Shape

```text
Unit Forms
Search by soldier name, DoD ID, or form name

Soldier group
  Current submission for each relevant template
  Sign button when commander signer is waiting
  Download button when complete
  Missing state when no current submission exists
```

Commanders also have a normal "My Forms" view for their own submissions.
