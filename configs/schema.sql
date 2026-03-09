CREATE TYPE monitoring_types AS ENUM ('cpu', 'memory');
CREATE TYPE scaling_algorithms AS ENUM ('simple', 'target');

CREATE TABLE templates(
    id SERIAL PRIMARY KEY,
    image_id BIGINT NOT NULL,
    networks BIGINT[] NOT NULL,
    SSH_keys BIGINT[],
    public_ipv4 BOOL NOT NULL,
    public_ipv6 BOOL NOT NULL,
    firewalls BIGINT[],
    cloud_config VARCHAR
);
-- `target` semantics depend on `scaling_algorithm`:
-- `simple`: target[0] is the scale-up threshold; target[1] is the scale-down threshold.
-- `target`: target[0] is the single target threshold.
CREATE TABLE groups(
    id SERIAL PRIMARY KEY,
    name VARCHAR NOT NULL,
    template_id INTEGER NOT NULL REFERENCES templates(id),
    zone VARCHAR NOT NULL,
    locations INTEGER[] NOT NULL,
    server_type VARCHAR NOT NULL,
    min_size SMALLINT NOT NULL,
    desired_size SMALLINT NOT NULL,
    max_size SMALLINT NOT NULL,
    monitoring_type monitoring_types NOT NULL,
    scaling_algorithm scaling_algorithms NOT NULL,
    target_threshold SMALLINT,
    scale_up_threshold SMALLINT,
    scale_down_threshold SMALLINT,
    CHECK (
        (scaling_algorithm = 'target' 
        AND target_threshold BETWEEN 1 AND 100
        AND target_threshold IS NOT NULL
        )
        OR
        (scaling_algorithm = 'simple' 
        AND scale_up_threshold BETWEEN 1 AND 100 
        AND scale_down_threshold BETWEEN 1 AND 100
        AND scale_up_threshold IS NOT NULL
        AND scale_down_threshold IS NOT NULL
        AND scale_down_threshold < scale_up_threshold
        )
    )
);

CREATE TABLE servers(
    id SERIAL PRIMARY KEY,
    name VARCHAR NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    group_id INTEGER NOT NULL REFERENCES groups(id),
    type VARCHAR NOT NULL,
    location INTEGER NOT NULL,
    private_ip INET NOT NULL
);
CREATE INDEX ON servers(group_id);
