# Keycloak and Moodle Authorization Setup

OTA Sign authenticates users with Keycloak and obtains OTA Sign capabilities from Moodle. Keycloak must provide identity attributes; Moodle must only decide the capabilities.

## 1. Configure Keycloak

1. Create or select the realm used by Moodle.
2. Create a confidential OpenID Connect client named `ota-sign`.
3. Enable Standard Flow and disable Direct Access Grants and Implicit Flow.
4. Set Valid Redirect URIs to `https://OTASIGN_HOST/auth/keycloak/callback`.
5. Copy the client secret into the OTA Sign deployment secret `KEYCLOAK_CLIENT_SECRET`.
6. Add client-scope mappers so the ID token contains these string claims:

```text
name
given_name
family_name
email
dod_id
uic
rank
army_email
```

7. Ensure the `dod_id` value exactly matches the Moodle user's `idnumber`. It must be unique for active Moodle users.
8. Configure Moodle to use this same Keycloak realm as its authentication source.

OTA Sign rejects login when `name`, `email`, `dod_id`, or `uic` are missing. `rank` and `army_email` are optional.

## 2. Install and Configure the Moodle Plugin

1. Install or upgrade `moodle/local_otasignconnector` in Moodle.
2. Complete the plugin upgrade so Moodle registers the `OTA Sign authorization service`.
3. Go to Site administration > Plugins > Local plugins > OTA Sign Connector.
4. Set OTA Sign login URL to `https://OTASIGN_HOST/auth/login`.
5. Leave Legacy launch signing secret empty.
6. Add a Moodle custom-menu link to `/local/otasignconnector/launch.php`.
7. Assign Moodle roles/capabilities in the system context:

```text
local/otasignconnector:viewown
local/otasignconnector:viewunit
local/otasignconnector:signascommander
local/otasignconnector:configure
```

The connector finds the person by `user.idnumber` and returns only these capabilities. It does not return profile data to OTA Sign.

## 3. Create the Moodle Service Token

1. Go to Site administration > Server > Web services > Overview and enable web services and the REST protocol.
2. Go to External services and open `OTA Sign authorization service`.
3. Authorize a dedicated service account. Do not use a human administrator account.
4. Create a token for that account and this service.
5. Store the token only as the OTA Sign secret `MOODLE_OTA_SIGN_SERVICE_TOKEN`.
6. Restrict Moodle network access to the OTA Sign backend where infrastructure permits.

## 4. Configure OTA Sign

Set these backend environment variables. Keep the two secret values in the deployment secret store, not Git or browser configuration.

```env
AUTH_PROVIDER=keycloak
KEYCLOAK_ISSUER_URL=https://KEYCLOAK_HOST/realms/REALM
KEYCLOAK_CLIENT_ID=ota-sign
KEYCLOAK_CLIENT_SECRET=KEYCLOAK_CLIENT_SECRET
KEYCLOAK_REDIRECT_URL=https://OTASIGN_HOST/auth/keycloak/callback
KEYCLOAK_DOD_ID_CLAIM=dod_id
KEYCLOAK_UIC_CLAIM=uic
KEYCLOAK_RANK_CLAIM=rank
KEYCLOAK_ARMY_EMAIL_CLAIM=army_email
MOODLE_WEBSERVICE_URL=https://MOODLE_HOST/webservice/rest/server.php
MOODLE_OTA_SIGN_SERVICE_TOKEN=MOODLE_SERVICE_TOKEN
```

Set `SESSION_COOKIE_SECURE=true` and `ENFORCE_HTTPS=true` in production. Deploy the backend and frontend together because the frontend must proxy `/auth/` to the backend.

## 5. Verify

1. Visit `https://OTASIGN_HOST/auth/login` directly and complete Keycloak login.
2. Confirm Keycloak returns to `/auth/keycloak/callback` and OTA Sign opens its dashboard.
3. Confirm the Keycloak DoD ID matches exactly one active Moodle `idnumber`.
4. Confirm a user with `viewown` can view personal forms.
5. Confirm `viewunit` controls unit dashboard access and `signascommander` controls commander signing.
6. Remove a capability in Moodle, sign out of OTA Sign, and sign in again. The feature must be denied.
7. Confirm the Moodle custom-menu link takes an already logged-in Moodle user through Keycloak SSO without a second password prompt.
