<?php
// This file is part of Moodle - http://moodle.org/

require_once(__DIR__ . '/../../config.php');
require_once($CFG->libdir . '/filelib.php');
require_once(__DIR__ . '/locallib.php');

require_login();

$launchurl = get_config('local_otasignconnector', 'launch_url');
$secret = get_config('local_otasignconnector', 'signing_secret');

if (empty($launchurl)) {
	throw new moodle_exception('missingconfig', 'local_otasignconnector');
}

// Keycloak mode uses Moodle only as the authorization authority. The user is
// still required to have a Moodle session before being sent to OTA Sign.
if (empty($secret)) {
    redirect(new moodle_url($launchurl));
}

$payload = local_otasignconnector_build_launch_payload($USER);
$token = local_otasignconnector_sign_payload($payload, $secret);

$redirect = new moodle_url($launchurl, ['token' => $token]);
redirect($redirect);
