DO $$
BEGIN
    IF EXISTS (
        SELECT 1
        FROM information_schema.columns
        WHERE table_schema = 'public'
          AND table_name = 'users'
          AND column_name = 'moodle_user_id'
    ) THEN
        ALTER TABLE users RENAME COLUMN moodle_user_id TO identity_user_id;
    END IF;
END $$;

ALTER TABLE user_unit_roles ALTER COLUMN source SET DEFAULT 'keycloak';
UPDATE user_unit_roles SET source = 'keycloak' WHERE source = 'moodle';
