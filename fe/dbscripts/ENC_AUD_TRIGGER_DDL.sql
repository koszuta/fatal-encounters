CREATE OR REPLACE FUNCTION audit.encounter_modified() RETURNS TRIGGER AS
$body$
DECLARE v_data RECORD;
BEGIN
    IF (TG_OP = 'INSERT' OR TG_OP = 'UPDATE') THEN
        v_data := NEW;
    ELSIF (TG_OP = 'DELETE') THEN
        v_data := OLD;
    ELSE
        RAISE WARNING '[AUDIT.ENCOUNTER_MODIFIED] - USER: % Unhandled action occurred: %, at %', SESSION_USER::TEXT, TG_OP, NOW();
        RETURN NULL;
    END IF;

    INSERT INTO audit.enc_aud (
        user_id,
        action,
        action_ts,
        uid,
        full_name,
        age,
        gender,
        race,
        race_imputed,
        race_imputation_prob,
        image_url,
        injury_date,
        address,
        city,
        state,
        zip_code,
        county,
        latitude,
        longitude,
        agency,
        cause_of_death,
        description,
        use_of_force,
        documentation_url,
        video_url,
        x_disposition_exclusion,
        x_mental_illness
    ) VALUES (
        SESSION_USER::TEXT,
        SUBSTRING(TG_OP, 1, 1),
        NOW(),
        v_data.uid,
        v_data.full_name,
        v_data.age,
        v_data.gender,
        v_data.race,
        v_data.race_imputed,
        v_data.race_imputation_prob,
        v_data.image_url,
        v_data.injury_date,
        v_data.address,
        v_data.city,
        v_data.state,
        v_data.zip_code,
        v_data.county,
        v_data.latitude,
        v_data.longitude,
        v_data.agency,
        v_data.cause_of_death,
        v_data.description,
        v_data.use_of_force,
        v_data.documentation_url,
        v_data.video_url,
        v_data.x_disposition_exclusion,
        v_data.x_mental_illness
    );

    RETURN v_data;

EXCEPTION
    WHEN data_exception THEN
        RAISE WARNING '[AUDIT.ENCOUNTER_MODIFIED] - UDF ERROR [DATA EXCEPTION] - USER: % SQLSTATE: %, SQLERRM: %', SESSION_USER::TEXT, SQLSTATE, SQLERRM;
        RETURN NULL;
    WHEN unique_violation THEN
        RAISE WARNING '[AUDIT.ENCOUNTER_MODIFIED] - UDF ERROR [UNIQUE] - USER: % SQLSTATE: %, SQLERRM: %', SESSION_USER::TEXT, SQLSTATE, SQLERRM;
        RETURN NULL;
    WHEN others THEN
        RAISE WARNING '[AUDIT.ENCOUNTER_MODIFIED] - UDF ERROR [OTHER] - USER: % SQLSTATE: %, SQLERRM: %', SESSION_USER::TEXT, SQLSTATE, SQLERRM;
        RETURN NULL;
END;
$body$
LANGUAGE plpgsql;

CREATE TRIGGER encounter_modified_trig
    AFTER INSERT OR UPDATE OR DELETE ON fe.encounters
    FOR EACH ROW EXECUTE PROCEDURE audit.encounter_modified();
