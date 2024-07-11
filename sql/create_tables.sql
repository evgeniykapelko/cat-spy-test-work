CREATE TABLE IF NOT EXISTS cats (
                                    id SERIAL PRIMARY KEY,
                                    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    name VARCHAR(255) NOT NULL,
    years_of_experience INT NOT NULL,
    breed VARCHAR(255) NOT NULL,
    salary NUMERIC(10, 2) NOT NULL
    );

CREATE TABLE IF NOT EXISTS missions (
                                        id SERIAL PRIMARY KEY,
                                        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    cat_id INT NOT NULL REFERENCES cats(id),
    complete BOOLEAN NOT NULL DEFAULT FALSE
    );

CREATE TABLE IF NOT EXISTS targets (
                                       id SERIAL PRIMARY KEY,
                                       created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    mission_id INT NOT NULL REFERENCES missions(id),
    name VARCHAR(255) NOT NULL,
    country VARCHAR(255) NOT NULL,
    notes TEXT,
    complete BOOLEAN NOT NULL DEFAULT FALSE
    );
