CREATE TABLE metrics (
                         id SERIAL PRIMARY KEY,
                         name VARCHAR(255) NOT NULL,
                         type VARCHAR(20) NOT NULL,
                         delta BIGINT,
                         value DOUBLE PRECISION,
                         hash VARCHAR(64),
                         created_at TIMESTAMPTZ DEFAULT NOW(),
                         updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE UNIQUE INDEX ux_metrics_name ON metrics(name);
