package sqlto

import (
	"encoding/json"
	"fmt"
	"io"
)

func (s *SQLto) JSON(w io.Writer) error {
	var err error

	colNames, err := s.Rows.Columns()
	if err != nil {
		return fmt.Errorf("error fetching column names, %s\n", err)
	}
	length := len(colNames)

	pointers := make([]interface{}, length)
	container := make([]interface{}, length)
	for i := range pointers {
		pointers[i] = &container[i]
	}

	jenc := json.NewEncoder(w)
	w.Write([]byte("["))
	delim := []byte("")

	for s.Rows.Next() {
		err = s.Rows.Scan(pointers...)
		if err != nil {
			return fmt.Errorf("error scanning sql row, %s\n", err)
		}

		w.Write(delim)
		err = jenc.Encode(container)
		delim = []byte(",")
		if err != nil {
			return fmt.Errorf("error encoding data, %s", err)
		}

	}
	w.Write([]byte("]"))
	return nil
}
