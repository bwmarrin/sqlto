package sqlto

import (
	"fmt"
	"io"
)

func (s *SQLto) HTML(w io.Writer) error {
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

	fmt.Fprintf(w, `<table>`)
	fmt.Fprintf(w, "<tr>")
	for _, v := range colNames {
		fmt.Fprintf(w, "<th>%s</th>", v)
	}
	fmt.Fprintf(w, "</tr>\r\n")

	for s.Rows.Next() {
		fmt.Fprintf(w, "<tr>")
		err := s.Rows.Scan(pointers...)
		if err != nil {
			fmt.Println("Failed to scan row", err)
			return err
		}

		for _, v := range container {
			if v == nil {
				fmt.Fprintf(w, "<td></td>")
			} else {
				fmt.Fprintf(w, "<td>%s</td>", v)
			}
		}
		fmt.Fprintf(w, "</tr>\r\n")
	}

	fmt.Fprintf(w, "</table>")

	return nil
}
