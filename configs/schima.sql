CREATE TABLE templates(
    id SERIAL PRIMARY KEY,
    image VARCHAR NOT NULL,
    cloud_config VARCHAR
);

CREATE TABLE groups(
    id SERIAL PRIMARY KEY,
    template_id INTEGER NOT NULL REFERENCES templates(id),
    zone VARCHAR NOT NULL,
    locations VARCHAR[] NOT NULL,
    server_types VARCHAR[] NOT NULL,
    min_size SMALLINT NOT NULL,
    desired_size SMALLINT NOT NULL,
    max_size SMALLINT NOT NULL
);

CREATE TABLE servers(
    id SERIAL PRIMARY KEY,
    private_ip INET NOT NULL
);