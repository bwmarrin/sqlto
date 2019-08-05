package sqlto

import (
	"fmt"
	"io"
	"strings"
)

func (s *SQLto) CSV(w io.Writer) error {

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

	delim := ""
	for _, v := range colNames {
		fmt.Fprintf(w, `%s"%v"`, delim, v)
		delim = ","
	}
	fmt.Fprintf(w, "\r\n")

	for s.Rows.Next() {

		err = s.Rows.Scan(pointers...)
		if err != nil {
			return fmt.Errorf("error scanning sql row, %s\n", err)
		}

		delim := ""

		for _, v := range container {

			if v == nil {
				fmt.Fprintf(w, `%s""`, delim)
			} else {
				s := fmt.Sprintf("%v", v)
				fmt.Fprintf(w, `%s"%s"`, delim, strings.Replace(s, `"`, `""`, -1))
			}
			delim = ","

		}

		fmt.Fprintf(w, "\r\n")
	}

	return nil
}
