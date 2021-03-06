CREATE TABLE IF NOT EXISTS fe.encounters (
    uid INTEGER PRIMARY KEY,
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

CREATE INDEX ON fe.encounters(age);
CREATE INDEX ON fe.encounters(gender);
CREATE INDEX ON fe.encounters(race_imputed);
CREATE INDEX ON fe.encounters(injury_date);
CREATE INDEX ON fe.encounters(state);
CREATE INDEX ON fe.encounters(zip_code);
CREATE INDEX ON fe.encounters(county);
CREATE INDEX ON fe.encounters(cause_of_death);
