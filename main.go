package main

import (
	"bytes"
	"context"
	"encoding/csv"
	"errors"
	"fatal-encounters/fe"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

const dataURL = "https://docs.google.com/spreadsheets/d/1dKmaV_JiWcG8XBoRgP8b4e9Eopkpgt7FL7nyspvzAsE/export?format=csv&id=1dKmaV_JiWcG8XBoRgP8b4e9Eopkpgt7FL7nyspvzAsE&gid=0"
const fileNameFormat = "data/fatal-encounters-%s.csv"
const fileDateFormat = "2006-01-02T15.04.05"
const headersPathFormat = "headers/v%d.csv"
const perms = os.FileMode(0666)

// hmm aud_id: 265376

func main() {
	log.Println("go...")
	defer log.Println("...done.")

	// Fetch data from Google Sheets
	res, err := http.Get(dataURL)
	fe.PanicOnErrorWithReason(err, "couldn't get data from %s", dataURL)

	// Read data from response
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	fe.PanicOnErrorWithReason(err, "couldn't read data")

	// Parse CSV
	rows, err := csv.NewReader(bytes.NewReader(body)).ReadAll()
	fe.PanicOnErrorWithReason(err, "couldn't parse data")

	// Validate CSV
	{
		if len(rows) == 0 {
			fe.PanicOnError(errors.New("invalid data: must contain at least a header row"))
		}

		v := 1
		verErr := fmt.Errorf("couldn't infer version from headers: %s", strings.Join(rows[0], ","))
	VERSION:
		for {
			headersPath := fmt.Sprintf(headersPathFormat, v)
			headersFile, err := os.Open(headersPath)
			if err != nil {
				break VERSION
			}
			defer headersFile.Close()

			headers, err := csv.NewReader(headersFile).Read()
			fe.PanicOnErrorWithReason(err, "couldn't parse expected headers in %s", headersPath)

			if len(rows[0]) == len(headers) {
				for i, header := range rows[0] {
					if header != headers[i] {
						v++
						continue VERSION
					}
				}
				verErr = nil
				break VERSION
			}
			v++
		}
		fe.PanicOnError(verErr)
		log.Printf("data version: v%d\n", v)
	}

	// Drop header row
	rows = rows[1:]

	encounters := make([]fe.Encounter, 0, 0)
	for i, row := range rows {
		e, err := fe.ParseRow(row)
		fe.PanicOnErrorWithReason(err, "parse error in row %d", i+1)

		if e.UID.Valid {
			encounters = append(encounters, e)
		}
	}

	fe.PanicOnError(err)
	log.Printf("encounters=%d\n", len(encounters))

	// Open database connection
	db, err := fe.OpenDB()
	fe.PanicOnError(err)
	defer db.Close()

	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	defer tx.Rollback()

	count := 0
	for _, e := range encounters {
		if e.InsertOrUpdate(tx) {
			count++
		}
	}
	log.Printf("adds/updates=%d\n", count)

	// Save data to a file (if any)
	if count > 0 {
		filePath := fmt.Sprintf(fileNameFormat, time.Now().UTC().Format(fileDateFormat))
		file, err := os.OpenFile(filePath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, perms)
		fe.PanicOnErrorWithReason(err, "couldn't create file")
		defer file.Close()
		n, err := file.Write(body)
		fe.PanicOnErrorWithReason(err, "couldn't save data to %s", filePath)
		log.Printf("wrote %d chars to %s\n", n, filePath)
	}

	err = tx.Commit()
	fe.PanicOnError(err)
}
