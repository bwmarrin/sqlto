package sqlto

import (
	"fmt"
	"io"
)

func (s *SQLto) HTML(w io.Writer) error {
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

	fmt.Fprintf(w, `<table>`)
	fmt.Fprintf(w, "<tr>")
	for _, v := range cols {
		fmt.Fprintf(w, "<th>%s</th>", v.Name())
	}
	fmt.Fprintf(w, "</tr>\r\n")

	for s.Rows.Next() {
		fmt.Fprintf(w, "<tr>")
		err := s.Rows.Scan(values...)
		if err != nil {
			fmt.Println("Failed to scan row", err)
			return err
		}

		for _, v := range values {
			if v == nil {
				fmt.Fprintf(w, "<td></td>")
			} else {
				fmt.Fprintf(w, "<td>%v</td>", v)
			}
		}
		fmt.Fprintf(w, "</tr>\r\n")
	}

	fmt.Fprintf(w, "</table>")

	return nil
}
