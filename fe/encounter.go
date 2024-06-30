package fe

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

const dateFormat = "01/02/2006"

// var specialDate = time.Date(2100, time.December, 31, 0, 0, 0, 0, time.UTC)

// Encounter ...
type Encounter struct {
	UID                   sql.NullInt64
	FullName              sql.NullString
	Age                   sql.NullFloat64
	Gender                sql.NullString
	Race                  sql.NullString
	RaceImputed           sql.NullString
	RaceImputationProb    sql.NullFloat64
	ImageURL              sql.NullString
	InjuryDate            sql.NullTime
	Address               sql.NullString
	City                  sql.NullString
	State                 sql.NullString
	ZipCode               sql.NullString
	County                sql.NullString
	Latitude              sql.NullFloat64
	Longitude             sql.NullFloat64
	Agency                sql.NullString
	CauseOfDeath          sql.NullString
	Description           sql.NullString
	UseOfForce            sql.NullString
	DocumentationURL      sql.NullString
	VideoURL              sql.NullString
	XDispositionExclusion sql.NullString
	XMentalIllness        sql.NullString
}

// ParseRow ...
func ParseRow(values []string) (Encounter, error) {
	// Check for and skip development rows
	if values[uIDColumn] == "" {
		log.Printf("development row found: %q\n", strings.Join(values, ","))
		return Encounter{}, nil
	}

	uID, err := strconv.ParseInt(values[uIDColumn], 10, 64)
	if err != nil {
		return Encounter{}, fmt.Errorf("invalid unique id: %v", err)
	}

	ageStr := values[ageColumn]
	age, ageErr := strconv.ParseFloat(ageStr, 64)
	if ageStr != "" && ageErr != nil {
		log.Printf("uID:%d invalid age: %v\n", uID, ageErr)
	}

	raceStr := values[raceColumn]
	raceImputedStr := values[raceImputedColumn]
	probStr := values[raceImputationProbColumn]
	imputationProb, probErr := strconv.ParseFloat(probStr, 64)
	if probErr != nil {
		if strings.EqualFold(probStr, "Not imputed") || (!strings.EqualFold(raceStr, "Race unspecified") && strings.EqualFold(raceStr, raceImputedStr)) {
			imputationProb = 1.0
			probErr = nil
		} else if !strings.EqualFold(probStr, "NA") && !strings.EqualFold(probStr, "Race not determined") {
			log.Printf("uID:%d invalid imputation probability: %v\n", uID, probErr)
		}
	} else if strings.EqualFold(raceImputedStr, "Race unspecified") || strings.EqualFold(raceImputedStr, "NA") {
		log.Printf("uID:%d race:%s raceImputed:%s prob:%f\n", uID, raceStr, raceImputedStr, imputationProb)
	}

	injuryDate, injErr := time.ParseInLocation(dateFormat, values[injuryDateColumn], time.UTC)
	if injErr != nil {
		log.Printf("uID:%d invalid injury date: %v\n", uID, injErr)
	}

	latitude, latErr := strconv.ParseFloat(values[latitudeColumn], 64)
	if latErr != nil {
		log.Printf("uID:%d invalid latitude: %v\n", uID, latErr)
	}

	longitude, lonErr := strconv.ParseFloat(values[longitudeColumn], 64)
	if lonErr != nil {
		log.Printf("uID:%d invalid longitude: %v\n", uID, lonErr)
	}

	return Encounter{
		UID:                   sql.NullInt64{Int64: uID, Valid: true},
		FullName:              NewSQLString(values[fullNameColumn]),
		Age:                   sql.NullFloat64{Float64: age, Valid: ageErr == nil},
		Gender:                NewSQLString(values[genderColumn]),
		Race:                  NewSQLString(raceStr),
		RaceImputed:           NewSQLString(raceImputedStr),
		RaceImputationProb:    sql.NullFloat64{Float64: imputationProb, Valid: probErr == nil},
		ImageURL:              NewSQLString(values[imageURLColumn]),
		InjuryDate:            sql.NullTime{Time: injuryDate, Valid: injErr == nil},
		Address:               NewSQLString(values[addressColumn]),
		City:                  NewSQLString(values[cityColumn]),
		State:                 NewSQLString(values[stateColumn]),
		ZipCode:               NewSQLString(values[zipCodeColumn]),
		County:                NewSQLString(values[countyColumn]),
		Latitude:              sql.NullFloat64{Float64: latitude, Valid: latErr == nil},
		Longitude:             sql.NullFloat64{Float64: longitude, Valid: lonErr == nil},
		Agency:                NewSQLString(values[agencyColumn]),
		CauseOfDeath:          NewSQLString(values[causeOfDeathColumn]),
		Description:           NewSQLString(values[descriptionColumn]),
		UseOfForce:            NewSQLString(values[useOfForceColumn]),
		DocumentationURL:      NewSQLString(values[documentationURLColumn]),
		VideoURL:              NewSQLString(""),
		XDispositionExclusion: NewSQLString(values[xDispositionExclusionColumn]),
		XMentalIllness:        NewSQLString(values[xMentalIllnessColumn]),
	}, nil
}

// GetByID ...
func (e *Encounter) GetByID(tx *sql.Tx) (Encounter, bool) {

	var existing Encounter
	row := tx.QueryRow(getByIDStatement, e.UID)
	err := row.Scan(&existing.UID, &existing.FullName, &existing.Age, &existing.Gender, &existing.Race, &existing.RaceImputed, &existing.RaceImputationProb, &existing.ImageURL, &existing.InjuryDate, &existing.Address, &existing.City, &existing.State, &existing.ZipCode, &existing.County, &existing.Latitude, &existing.Longitude, &existing.Agency, &existing.CauseOfDeath, &existing.Description, &existing.UseOfForce, &existing.DocumentationURL, &existing.VideoURL, &existing.XDispositionExclusion, &existing.XMentalIllness)
	if err != nil && err != sql.ErrNoRows {
		log.Printf("uID:%d couldn't get row: %v\n", e.UID.Int64, err)
		return existing, false
	}
	return existing, err != sql.ErrNoRows
}

// InsertOrUpdate ...
func (e *Encounter) InsertOrUpdate(tx *sql.Tx) bool {
	// log.Printf("uID:%d insert or update\n", e.UID.Int64)

	existing, found := e.GetByID(tx)
	if !found {
		return doInsertOrUpdate(e, insertStatement, "insert", tx)
	}
	if *e != existing {
		log.Printf("uID:%d diff: %s\n", existing.UID.Int64, StructDiff(existing, *e))
		return doInsertOrUpdate(e, updateStatement, "update", tx)
	}
	// log.Printf("uID:%d no change to existing data\n", e.UID.Int64)
	return false
}

func doInsertOrUpdate(e *Encounter, statement, action string, tx *sql.Tx) bool {
	var actionPastTense string
	if strings.HasSuffix(action, "e") {
		actionPastTense = action + "d"
	} else {
		actionPastTense = action + "ed"
	}
	// log.Printf("uID:%d data %s required", e.UID.Int64, action)
	result, err := tx.Exec(statement, e.UID, e.FullName, e.Age, e.Gender, e.Race, e.RaceImputed, e.RaceImputationProb, e.ImageURL, e.InjuryDate, e.Address, e.City, e.State, e.ZipCode, e.County, e.Latitude, e.Longitude, e.Agency, e.CauseOfDeath, e.Description, e.UseOfForce, e.DocumentationURL, e.VideoURL, e.XDispositionExclusion, e.XMentalIllness)
	if err != nil {
		log.Printf("uID:%d couldn't %s row: %v\n", e.UID.Int64, action, err)
		return false
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("uID:%d couldn't get number of rows %s: %v\n", e.UID.Int64, actionPastTense, err)
	}
	if rowsAffected != 1 {
		log.Printf("uID:%d unexpected number of rows %s: wanted 1, got %d\n", e.UID.Int64, actionPastTense, rowsAffected)
	}
	// log.Printf("uID:%d data %s\n", e.UID.Int64, actionPastTense)
	return true
}
