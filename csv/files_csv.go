package csv

import (
	"capitalbank/config"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

func LoadCSVfiles(records []CSVRecord, allfiles []CSVfiles) ([]CSVRecord, []CSVfiles, error) {
	dir := config.Config["csvimportdir"].(string)
	//dir := "\\\\pay\\c$\\IMPORT" // specify your folder path here
	files, err := os.ReadDir(dir)
	if err != nil {
		return records, allfiles, nil
	}

	for _, f := range files {
		if !f.IsDir() && filepath.Ext(f.Name()) == ".csv" {
			fmt.Println("Processing file:", f.Name())
			file, err := os.Open(filepath.Join(dir, f.Name()))
			if err != nil {
				fmt.Println("Can't open file: ", f.Name(), err)
			}
			allfiles = append(allfiles, CSVfiles{FileName: filepath.Join(dir, f.Name())})

			// Create a new reader that converts Windows-1252 to UTF-8
			win1252Reader := transform.NewReader(file, charmap.Windows1251.NewDecoder())
			// Create a new CSV reader
			reader := csv.NewReader(win1252Reader)
			reader.Comma = ';'
			reader.FieldsPerRecord = -1 // don't check number of fields

			// skip first line
			_, err = reader.Read()
			if err != nil {
				panic(err)
			}

			for {
				line, err := reader.Read()
				if err == io.EOF {
					break
				} else if err != nil {
					//panic(err)
					fmt.Println("Error reading line", err)
					continue
				}
				record := NewCSVRecord(line)
				records = append(records, *record)
			}
			// At this point, 'records' holds the data of the CSV file
			// Close the file before attempting to delete it
			if err := file.Close(); err != nil {
				panic(err)
			}

		}
	}
	return records, allfiles, nil
}

func DeleteCSVfiles(allfiles []CSVfiles) {
	// Delete the file
	for _, f := range allfiles {
		if err := os.Remove(f.FileName); err != nil {
			fmt.Println("File cannot deleted: ", err.Error())
		}
		fmt.Println("File processed and deleted: ", f.FileName)
	}
}
