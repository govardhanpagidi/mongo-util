package mongo

import (
	"encoding/csv"
	"os"
)

func GenerateCSV(entries [][]string, fileName string) error {

	csvFile, err := os.Create(fileName + ".csv")

	if err != nil {
		return err
	}
	csvWriter := csv.NewWriter(csvFile)
	defer csvWriter.Flush()

	for _, entry := range entries {
		csvWriter.Write(entry)
	}
	return err
}
