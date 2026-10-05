# OTA Sign Connector

Moodle local plugin component:

```text
local_otasignconnector
```

Display name:

```text
OTA Sign Connector
```

## Keycloak Authorization Purpose

The plugin provides an optional Moodle launch link and the Moodle authorization service used by OTA Sign after Keycloak login.

```text
OTA Sign validates the Keycloak user identity
-> backend sends DoD ID to this plugin's REST service
-> plugin finds active Moodle user by user.idnumber
-> plugin returns OTA Sign capabilities only
```

## Required Settings

- OTA Sign login URL, for example `https://sign.example.mil/auth/login`
- Leave Legacy launch signing secret empty for Keycloak authentication.
- Ensure each active Moodle user has a unique `user.idnumber` equal to the Keycloak `dod_id` claim.

## Moodle Navbar Link

Add a Moodle custom menu/navbar link that points to the plugin launch page:

```text
OTA Sign|/local/otasignconnector/launch.php" target="_blank|Digitally sign your documents for OTA.
```

In Keycloak mode, the plugin launch page requires the Moodle session and then redirects to OTA Sign `/auth/login`. Keycloak SSO normally returns the user without another password prompt.

## Authorization Service

The plugin registers the restricted `OTA Sign authorization service` with one function: `local_otasignconnector_get_access`. Create a dedicated Moodle service account and token for that service, then set the token only in OTA Sign's `MOODLE_OTA_SIGN_SERVICE_TOKEN` deployment secret. See [`docs/keycloak-setup.md`](../../../docs/keycloak-setup.md) for the complete setup.

## Install Path

Copy this directory into Moodle as:

```text
local/otasignconnector
```
