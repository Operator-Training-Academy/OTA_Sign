<?php
// This file is part of Moodle - http://moodle.org/

$string['pluginname'] = 'OTA Sign Connector';
$string['open_otasign'] = 'OTA Sign';
$string['launch_url'] = 'OTA Sign login URL';
$string['launch_url_desc'] = 'OTA Sign login endpoint, normally https://otasign.example.mil/auth/login. Keycloak performs authentication; Moodle supplies authorization through the connector service.';
$string['signing_secret'] = 'Legacy launch signing secret';
$string['signing_secret_desc'] = 'Leave empty for Keycloak login. Set only when using the legacy Moodle-signed launch-token flow.';
$string['uic_profile_field'] = 'UIC profile field shortname';
$string['uic_profile_field_desc'] = 'Moodle custom profile field shortname containing the user UIC.';
$string['dodid_profile_field'] = 'DoD ID profile field shortname';
$string['dodid_profile_field_desc'] = 'Moodle custom profile field shortname containing the user DoD ID.';
$string['rank_profile_field'] = 'Rank profile field shortname';
$string['rank_profile_field_desc'] = 'Moodle custom profile field shortname containing the user rank abbreviation. OTA Sign derives pay grade from this rank.';
$string['army_email_profile_field'] = 'Army email profile field shortname';
$string['army_email_profile_field_desc'] = 'Moodle custom profile field shortname containing the user @army.mil email address. If empty or not an @army.mil address, OTA Sign falls back to the regular Moodle email only when it is an @army.mil address.';
$string['missingconfig'] = 'OTA Sign Connector is missing its launch URL.';
$string['accessdenied'] = 'OTA Sign access is not available for this user.';
$string['otasignconnector:viewown'] = 'Launch OTA Sign and view own forms';
$string['otasignconnector:viewunit'] = 'View OTA Sign unit forms';
$string['otasignconnector:signascommander'] = 'Sign OTA Sign forms as commander';
$string['otasignconnector:configure'] = 'Configure OTA Sign Connector';
