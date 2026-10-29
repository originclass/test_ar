-- 巡检问题表(go-admin 库)
CREATE TABLE IF NOT EXISTS app_patrol_issue (
	id          BIGSERIAL PRIMARY KEY,
	title       VARCHAR(255) NOT NULL,
	description VARCHAR(1000) DEFAULT '',
	priority    VARCHAR(10) NOT NULL,
	status      VARCHAR(20) NOT NULL,
	pos_x       DOUBLE PRECISION DEFAULT 0,
	pos_y       DOUBLE PRECISION DEFAULT 0,
	pos_z       DOUBLE PRECISION DEFAULT 0,
	create_by   INT DEFAULT 0,
	update_by   INT DEFAULT 0,
	created_at  TIMESTAMPTZ DEFAULT now(),
	updated_at  TIMESTAMPTZ DEFAULT now()
);
