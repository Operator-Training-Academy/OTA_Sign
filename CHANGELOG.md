# Changelog

## [1.1.0] - 2026-10-14

### Added

- Restored temporary Moodle signed-launch compatibility, selected explicitly with `AUTH_PROVIDER=moodle`.
- Preserved Keycloak OIDC mode, selected explicitly with `AUTH_PROVIDER=keycloak`.

## [1.0.1] - 2026-10-14

### Fixed

- Restricted Unit Forms submissions and completion calculations to the viewing user's UIC, hiding mismatched and unassigned historical submissions.

### Added

- Added regression coverage for cross-UIC submission filtering.
