package fe

const (
	getByIDStatement = `
	SELECT
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
	FROM fe.encounters
	WHERE uid = $1;`

	insertStatement = `
	INSERT INTO fe.encounters(
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
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24);`

	updateStatement = `
    UPDATE fe.encounters
    SET
        full_name = $2,
        age = $3,
        gender = $4,
        race = $5,
        race_imputed = $6,
        race_imputation_prob = $7,
        image_url = $8,
        injury_date = $9,
        address = $10,
        city = $11,
        state = $12,
        zip_code = $13,
        county = $14,
        latitude = $15,
        longitude = $16,
        agency = $17,
        cause_of_death = $18,
        description = $19,
        use_of_force = $20,
        documentation_url = $21,
        video_url = $22,
        x_disposition_exclusion = $23,
        x_mental_illness = $24
    WHERE uid = $1;`
)
