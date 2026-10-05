<?php
// This file is part of Moodle - http://moodle.org/

defined('MOODLE_INTERNAL') || die();

$functions = [
    'local_otasignconnector_get_access' => [
        'classname' => 'local_otasignconnector\\external\\get_access',
        'methodname' => 'execute',
        'description' => 'Returns OTA Sign capabilities for an active Moodle user identified by DoD ID.',
        'type' => 'read',
        'ajax' => false,
    ],
];

$services = [
    'OTA Sign authorization service' => [
        'functions' => ['local_otasignconnector_get_access'],
        'restrictedusers' => 1,
        'enabled' => 1,
        'shortname' => 'otasign_authorization',
    ],
];
