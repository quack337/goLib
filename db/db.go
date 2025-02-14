package db

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

func Jsonify(rows *sql.Rows) (string, error) {
	columns, err := rows.Columns()
	if err != nil {
		return "", err
	}

	values := make([]any, len(columns))
	scanArgs := make([]any, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	data := []string{"["}
	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			panic(err.Error())
		}

		data = append(data, "{")
		for i, value := range values {
			if value == nil {
				data = append(data, fmt.Sprintf("\"%s\": null", columns[i]))
			} else {
				switch val := value.(type) {
				case int64:
					data = append(data, fmt.Sprintf("\"%s\": %d", columns[i], val))
				case float64:
					data = append(data, fmt.Sprintf("\"%s\": %s", columns[i], strconv.FormatFloat(val, 'f', -1, 64)))
				default:
					data = append(data, fmt.Sprintf("\"%s\": \"%v\"", columns[i], val))
				}
			}
			data = append(data, ", ")
		}
		if data[len(data)-1] == ", " {
			data[len(data)-1] = "}, "
		} else {
			data = append(data, "}, ")
		}
	}
	if data[len(data)-1] == "}, " {
		data[len(data)-1] = "}]"
	} else {
		data = append(data, "]")
	}
	return strings.Join(data, ""), nil
}
