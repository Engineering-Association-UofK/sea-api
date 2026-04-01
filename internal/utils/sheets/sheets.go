package sheets

import (
	"io"
	"reflect"
	"sea-api/internal/errs"
	"strconv"
	"strings"

	"github.com/xuri/excelize/v2"
)

func ParseExcelToStructs[T any](fileReader io.Reader) ([]T, error) {
	f, err := excelize.OpenReader(fileReader)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	rows, err := f.GetRows("Sheet1")
	if err != nil || len(rows) < 2 {
		return nil, errs.New(errs.BadRequest, "sheet is empty or missing headers. Did you check sheet name?", nil)
	}

	// Get header names to use in parsing
	headerMap := make(map[string]int)
	for i, name := range rows[0] {
		headerMap[strings.TrimSpace(name)] = i
	}

	var result []T
	tType := reflect.TypeOf((*T)(nil)).Elem()

	for i := 1; i < len(rows); i++ {
		var item T
		vItem := reflect.ValueOf(&item).Elem()

		// Match struct tags to header map
		for j := 0; j < tType.NumField(); j++ {
			field := tType.Field(j)
			tag := field.Tag.Get("excel")

			if colIdx, ok := headerMap[tag]; ok && colIdx < len(rows[i]) {
				val := rows[i][colIdx]
				fValue := vItem.Field(j)

				// Basic type conversion
				switch fValue.Kind() {
				case reflect.String:
					fValue.SetString(val)
				case reflect.Int:
					intValue, _ := strconv.Atoi(val)
					fValue.SetInt(int64(intValue))
				}
			}
		}
		result = append(result, item)
	}

	return result, nil
}
