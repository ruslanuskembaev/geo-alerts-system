CREATE INDEX IF NOT EXISTS idx_incidents_active ON incidents (is_active);

CREATE INDEX IF NOT EXISTS idx_location_checks_checked_at ON location_checks (checked_at);
CREATE INDEX IF NOT EXISTS idx_location_checks_user_id ON location_checks (user_id);

CREATE INDEX IF NOT EXISTS idx_location_check_incidents_incident ON location_check_incidents (incident_id);
