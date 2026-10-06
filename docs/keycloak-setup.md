# Keycloak Setup

OTA Sign authenticates and authorizes users exclusively through Keycloak. Moodle is not part of the OTA Sign login or authorization flow.

## Keycloak Client

1. Create a confidential OpenID Connect client named `ota-sign` in the OTA realm.
2. Enable Standard Flow and disable Direct Access Grants and Implicit Flow.
3. Set the redirect URI to `https://OTASIGN_HOST/auth/keycloak/callback`.
4. Save the client secret as OTA Sign's `KEYCLOAK_CLIENT_SECRET` deployment secret.
5. Keep the built-in `basic`, `profile`, and `email` client scopes assigned to the client.

## Identity Claims

In the `ota-sign` dedicated client scope, add User Attribute mappers that include these Keycloak user-profile attributes in the ID token:

```text
User attribute       Token claim
dod_id               dod_id
unit_uic             uic
rank                 rank
enterprise_email     army_email
```

Enable Add to ID token and Add to UserInfo for each mapper. The standard scopes provide `name`, `given_name`, `family_name`, and `email`.

## Authorization Roles

Create these client roles on the `ota-sign` client:

```text
viewown
viewunit
signascommander
configure
```

Add `viewown` as a composite of the realm's default role so every OTA realm user can view their own forms. Add `viewunit` as a composite of `signascommander`; this makes commander signing imply unit visibility without allowing unit viewers to sign.

Configure the client's role mapper to include client roles in the ID token as:

```json
{
  "resource_access": {
    "ota-sign": {
      "roles": ["viewown", "viewunit", "signascommander"]
    }
  }
}
```

OTA Sign denies login unless the validated ID token includes `viewown`. The UIC claim scopes unit viewing and commander signing: a user with `uic=WABC12` can access only WABC12 submissions.

## OTA Sign Configuration

Set these backend deployment values. Do not expose the client secret to the browser or commit it to Git.

```env
KEYCLOAK_ISSUER_URL=https://KEYCLOAK_HOST/realms/ota
KEYCLOAK_CLIENT_ID=ota-sign
KEYCLOAK_CLIENT_SECRET=KEYCLOAK_CLIENT_SECRET
KEYCLOAK_REDIRECT_URL=https://OTASIGN_HOST/auth/keycloak/callback
KEYCLOAK_DOD_ID_CLAIM=dod_id
KEYCLOAK_UIC_CLAIM=uic
KEYCLOAK_RANK_CLAIM=rank
KEYCLOAK_ARMY_EMAIL_CLAIM=army_email
SESSION_COOKIE_SECURE=true
ENFORCE_HTTPS=true
```

## Verify

1. Open `https://OTASIGN_HOST/auth/login`.
2. Sign in with a user that has `viewown`.
3. Confirm the dashboard shows the user's Keycloak name and UIC.
4. Confirm a user with `viewunit` sees only their claim's UIC.
5. Confirm a user with `signascommander` can sign only submissions for their claim's UIC.
6. Confirm a user without `viewown` is denied without an OTA Sign session.
