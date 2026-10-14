# OTA Sign Architecture

## System Roles

```text
Keycloak
- Target login, identity, and authorization provider
- DoD ID, UIC, rank, name, and email claims

Moodle (temporary compatibility mode)
- Authenticated signed launch token
- DoD ID, UIC, rank, name, email, and capabilities

OTA Sign Backend
- Validates Keycloak OIDC tokens or Moodle signed launch tokens
- Maps Keycloak client roles or Moodle capabilities to OTA Sign permissions
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
1. In Keycloak mode, OTA Sign redirects the user to Keycloak, validates the returned ID token, and maps `ota-sign` client roles to capabilities.
2. In temporary Moodle mode, the Moodle connector sends a signed launch token to OTA Sign, which validates the signature and expiry.
3. OTA Sign uses the verified UIC as the authorization scope.
4. Backend creates a secure session.
5. Frontend loads /api/me and dashboard data from the backend.
```

## Security Rules

- The frontend never receives Keycloak client secrets, Moodle launch signing secrets, DocuSeal API keys, SMTP credentials, or database credentials.
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
