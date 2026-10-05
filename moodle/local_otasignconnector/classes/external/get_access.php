<?php
// This file is part of Moodle - http://moodle.org/

namespace local_otasignconnector\external;

defined('MOODLE_INTERNAL') || die();

require_once($CFG->dirroot . '/local/otasignconnector/locallib.php');
require_once($CFG->libdir . '/externallib.php');

use context_system;
use external_api;
use external_function_parameters;
use external_multiple_structure;
use external_single_structure;
use external_value;

class get_access extends external_api {
    public static function execute_parameters(): external_function_parameters {
        return new external_function_parameters([
            'dodid' => new external_value(PARAM_ALPHANUMEXT, 'The Keycloak-validated DoD ID.'),
        ]);
    }

    public static function execute(string $dodid): array {
        global $DB;

        $params = self::validate_parameters(self::execute_parameters(), ['dodid' => $dodid]);
        $dodid = trim($params['dodid']);
        if ($dodid === '') {
            throw new \invalid_parameter_exception('DoD ID is required.');
        }

        $users = $DB->get_records('user', [
            'idnumber' => $dodid,
            'deleted' => 0,
            'suspended' => 0,
        ], '', 'id, idnumber');
        if (count($users) !== 1) {
            throw new \moodle_exception('accessdenied', 'local_otasignconnector');
        }

        $user = reset($users);
        self::validate_context(context_system::instance());

        return [
            'capabilities' => local_otasignconnector_user_capabilities($user),
        ];
    }

    public static function execute_returns(): external_single_structure {
        return new external_single_structure([
            'capabilities' => new external_multiple_structure(
                new external_value(PARAM_ALPHANUMEXT, 'OTA Sign capability name.'),
                'OTA Sign capabilities granted by Moodle.',
            ),
        ]);
    }
}
