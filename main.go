package main

import (
	"bufio"
	"bytes"
	"compress/flate"
	"compress/gzip"
	"context"
	"encoding/csv"
	"fatal-encounters/fe"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/cheggaaa/pb"
)

const (
	dataURL           = "https://docs.google.com/spreadsheets/d/1dKmaV_JiWcG8XBoRgP8b4e9Eopkpgt7FL7nyspvzAsE/export?format=csv&id=1dKmaV_JiWcG8XBoRgP8b4e9Eopkpgt7FL7nyspvzAsE&gid=0"
	fileNameFormat    = "data/fatal-encounters-%s.csv"
	fileDateFormat    = "2006-01-02T15.04.05"
	headersPathFormat = "headers/v%d.csv"
)

// hmm aud_id: 265376

func main() {
	// Fetch data from Google Sheets
	res, err := http.Get(dataURL)
	if err != nil {
		log.Fatalf("couldn't get data from %s: %v\n", dataURL, err)
	}

	// Read data from response
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatalln("couldn't read data:", err)
	}
	defer res.Body.Close()

	// Parse CSV
	rows, err := csv.NewReader(bytes.NewReader(body)).ReadAll()
	if err != nil {
		log.Fatalln("couldn't parse data:", err)
	}

	{ // validate CSV
		if len(rows) == 0 {
			log.Fatalln("invalid data: must contain at least a header row")
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
			if err != nil {
				log.Fatalf("couldn't parse expected headers in %s: %v\n", headersPath, err)
			}

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
		if verErr != nil {
			log.Fatalln(verErr)
		}
		log.Printf("data version: v%d\n", v)
	}

	// Drop header row
	rows = rows[1:]

	encounters := make([]fe.Encounter, 0)
	for i, row := range rows {
		if e, err := fe.ParseRow(row); err != nil {
			log.Fatalf("parse error in row %d: %v\n", i+1, err)
		} else if e.UID.Valid {
			encounters = append(encounters, e)
		}
	}
	log.Printf("n encounters=%d\n", len(encounters))

	// Open database connection
	db, err := fe.OpenDB()
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		log.Fatalln("couldn't open transaction:", err)
	}
	defer tx.Rollback() // ok to call after Commit

	var count int
	pb := pb.StartNew(len(encounters))
	for _, e := range encounters {
		if e.InsertOrUpdate(tx) {
			count++
		}
		pb.Increment()
	}
	pb.Finish()
	log.Printf("n adds/updates=%d\n", count)

	if count == 0 {
		return
	}

	// Save data to a file (if any)
	var input string
	r := bufio.NewReader(os.Stdin)
	matchN, matchY := regexp.MustCompile(`^\s*[nN]\s*$`), regexp.MustCompile(`^\s*[yY]\s*$`)
	for {
		log.Print("commit (y/n)?")
		input, _ = r.ReadString('\n')
		if matchN.MatchString(input) {
			log.Println("rolling back changes...") // deferred
			return
		}
		if matchY.MatchString(input) {
			break
		}
	}

	filePath := fmt.Sprintf(fileNameFormat, time.Now().UTC().Format(fileDateFormat))
	var buf bytes.Buffer
	w, err := gzip.NewWriterLevel(&buf, flate.BestCompression)
	if err != nil {
		log.Fatalln("couldn't create gzip writer:", err)
	}
	if _, err := w.Write(body); err != nil {
		log.Fatalln(err)
	}
	// The writer must be closed to flush the bytes to the buffer
	if err := w.Close(); err != nil {
		log.Fatalln(err)
	}
	if err := os.WriteFile(filePath+".gz", buf.Bytes(), 0666); err != nil {
		log.Fatalln("couldn't save data to", filePath, ":", err)
	}
	log.Printf("%d bytes written to %s\n", len(buf.Bytes()), filePath)

	if err := tx.Commit(); err != nil {
		log.Fatalln(err)
	}
	log.Println("changes committed...")
}
