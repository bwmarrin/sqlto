package sqlto

import (
	"fmt"
	"io"
)

func (s *SQLto) String(w io.Writer) error {
	var err error

	// XXX: don't really need col names here
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

	for s.Rows.Next() {
		err := s.Rows.Scan(pointers...)
		if err != nil {
			return err
		}

		for _, v := range container {
			if v == nil {
			} else {
				fmt.Fprintf(w, "%s", v)
			}
		}
	}
	return nil
}
