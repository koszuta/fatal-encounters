CREATE TABLE IF NOT EXISTS audit.enc_aud (
    aud_id SERIAL PRIMARY KEY,
    user_id TEXT NOT NULL,
    action CHAR(1) NOT NULL CHECK (action IN ('I', 'U', 'D')),
    action_ts TIMESTAMP WITH TIME ZONE NOT NULL,
    uid INTEGER,
    full_name TEXT,
    age NUMERIC,
    gender TEXT,
    race TEXT,
    race_imputed TEXT,
    race_imputation_prob NUMERIC,
    image_url TEXT,
    injury_date DATE,
    address TEXT,
    city TEXT,
    state TEXT,
    zip_code TEXT,
    county TEXT,
    latitude NUMERIC,
    longitude NUMERIC,
    agency TEXT,
    cause_of_death TEXT,
    description TEXT,
    use_of_force TEXT,
    documentation_url TEXT,
    video_url TEXT,
    x_disposition_exclusion TEXT,
    x_mental_illness TEXT
);

CREATE INDEX ON audit.enc_aud(action);
CREATE INDEX ON audit.enc_aud(action_ts);
CREATE INDEX ON audit.enc_aud(uid);
