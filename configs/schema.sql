CREATE TYPE monitoring_types AS ENUM ('cpu', 'memory');

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
    target SMALLINT NOT NULL check(target BETWEEN 1 AND 100)
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