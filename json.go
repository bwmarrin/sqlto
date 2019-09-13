package sqlto

import (
	"encoding/json"
	"fmt"
	"io"
)

// JSONArray converts sql rows into an array of JSON arrays
func (s *SQLto) JSONArray(w io.Writer) error {
	var err error

	cols, err := s.Rows.ColumnTypes()
	if err != nil {
		return fmt.Errorf("error fetching column types, %s\n", err)
	}
	colsLen := len(cols)

	values := make([]interface{}, colsLen)
	for i := range values {
		values[i] = &GenericScanner{dbtype: cols[i].DatabaseTypeName()}
	}

	j := json.NewEncoder(w)
	w.Write([]byte("["))
	delim := []byte("")

	for s.Rows.Next() {
		err = s.Rows.Scan(values...)
		if err != nil {
			return fmt.Errorf("error scanning sql row, %s\n", err)
		}

		w.Write(delim)
		err = j.Encode(values)
		delim = []byte(",")
		if err != nil {
			return fmt.Errorf("error encoding data, %s", err)
		}

	}
	w.Write([]byte("]"))
	return nil
}

// JSONObject converts sql rows into an array of JSON objects
func (s *SQLto) JSONObject(w io.Writer) error {
	var err error

	cols, err := s.Rows.ColumnTypes()
	if err != nil {
		return fmt.Errorf("error fetching column information, %s\n", err)
	}
	colsLen := len(cols)

	// scan is needed here because Rows.Scan requires a []interface{} argument
	scan := make([]interface{}, colsLen)
	values := make(map[string]interface{})
	for i, v := range cols {
		values[v.Name()] = &GenericScanner{dbtype: v.DatabaseTypeName()}
		scan[i] = values[v.Name()]
	}

	j := json.NewEncoder(w)
	w.Write([]byte("["))
	delim := []byte("")

	for s.Rows.Next() {
		err = s.Rows.Scan(scan...)
		if err != nil {
			return fmt.Errorf("error scanning sql row, %s\n", err)
		}

		w.Write(delim)
		err = j.Encode(values)
		delim = []byte(",")
		if err != nil {
			return fmt.Errorf("error encoding data, %s", err)
		}

	}
	w.Write([]byte("]"))
	return nil
}
