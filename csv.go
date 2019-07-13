package gtfs

import (
	"encoding/csv"
	"fmt"
	"io"
)

func readCSVWithHeadings(r io.Reader, fields map[string]bool, strictMode bool) ([]map[string]string, error) {
	var headerFields []string
	var res []map[string]string

	csvFile := csv.NewReader(r)
	csvFile.FieldsPerRecord = -1 // Ignore mismatched numbers of fields
	csvFile.LazyQuotes = true    // Allow different quoting styles

	headers, err := csvFile.Read()
	if err != nil {
		return nil, err
	}

	for _, h := range headers {
		if _, ok := fields[h]; !ok {
			continue
		}

		headerFields = append(headerFields, h)
	}

	for {
		row, err := csvFile.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			return res, err
		}

		rowMap := map[string]string{}
		for i, v := range row {
			if i >= len(headerFields) {
				if strictMode {
					return res, fmt.Errorf("unexpected number of fields in row: %d", i+1)
				}

				continue
			}

			rowMap[headerFields[i]] = v
		}

		res = append(res, rowMap)
	}

	return res, nil
}
