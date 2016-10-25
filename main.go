package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"time"

	"github.com/Sirupsen/logrus"
)

var log = logrus.New()

func init() {
	log.Formatter = new(logrus.TextFormatter)
	log.Level = logrus.DebugLevel
}

func main() {
	var importFolder, exportFile string

	flag.StringVar(&importFolder, "input", "", "Path the text files of Appboy data to import")
	flag.StringVar(&exportFile, "output", "", "Path to the CSV file to export to")
	flag.Parse()

	halt := false
	if importFolder == "" {
		log.Error("Path to text files of Appboy data is required.")
		halt = true
	}

	if exportFile == "" {
		log.Error("Path to export CSV file is required.")
		halt = true
	}

	if halt {
		flag.Usage()
		os.Exit(2)
	}

	fmt.Println("===========================================")
	fmt.Println(fmt.Sprintf(" Use CPU(s) num      : %d", runtime.NumCPU()))
	fmt.Println(fmt.Sprintf(" Sleep time          : 0.2 second(s)"))
	fmt.Println("===========================================")

	runtime.GOMAXPROCS(runtime.NumCPU())

	// 1. Walk through the directory of files and make an array of files to iterate.
	importFiles := []string{}
	err := filepath.Walk(importFolder, func(path string, f os.FileInfo, err error) error {
		importFiles = append(importFiles, path)
		return nil
	})
	exitOnError(err)

	importRecords := []AppboyRecord{}
	// 2. Iterate through each file
	for _, path := range importFiles {
		func() {
			file, openErr := os.Open(path)
			logOnError(openErr)
			defer file.Close()

			// 2a. Parse current file contents
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				// Unmarshal current line into an `AppboyRecord`
				record := AppboyRecord{}
				jsonErr := json.Unmarshal(scanner.Bytes(), &record)
				logOnError(jsonErr)
				importRecords = append(importRecords, record)
			}
			log.WithFields(logrus.Fields{
				"count": len(importRecords),
			}).Info("Importing records")
		}()
	}

	// 3. Take the `importRecords` and generate a CSV out of it.
	log.Info(len(importRecords))

	exporter, err := os.Create(exportFile)
	exitOnError(err)
	defer exporter.Close()

	writer := csv.NewWriter(exporter)

	// 4. Write the headers
	writer.Write([]string{"idfv", "first_used", "last_used", "token", "sessions"})

	for _, record := range importRecords {
		var row []string

		log.Info(record)

		if len(record.Devices) > 0 {
			row = append(row, record.Devices[0].IDFV)
		} else {
			row = append(row, "")
		}

		if len(record.Apps) > 0 {
			t, err := time.Parse(time.RFC3339, record.Apps[0].FirstUsed)

			if err != nil {
				exitOnError(err)
			}

			row = append(row, strconv.FormatInt(t.Unix(), 10))
		} else {
			row = append(row, "")
		}

		if len(record.Apps) > 0 {
			t, err := time.Parse(time.RFC3339, record.Apps[0].LastUsed)

			if err != nil {
				exitOnError(err)
			}

			row = append(row, strconv.FormatInt(t.Unix(), 10))
		} else {
			row = append(row, "")
		}

		if len(record.Tokens) > 0 {
			row = append(row, record.Tokens[0].Token)
		} else {
			row = append(row, "")
		}

		if len(record.Apps) > 0 {
			row = append(row, strconv.Itoa(record.Apps[0].Sessions))
		} else {
			row = append(row, "")
		}
		log.Info(row)

		err := writer.Write(row)
		logOnError(err)
	}

	defer writer.Flush()
}
