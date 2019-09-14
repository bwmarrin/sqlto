package sqlto

import (
	"fmt"
	"io"
	"strings"
)

func (s *SQLto) CSV(w io.Writer) error {
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

	delim := ""
	for _, v := range cols {
		fmt.Fprintf(w, `%s"%v"`, delim, v.Name())
		delim = ","
	}
	fmt.Fprintf(w, "\r\n")

	for s.Rows.Next() {

		err = s.Rows.Scan(values...)
		if err != nil {
			return fmt.Errorf("error scanning sql row, %s\n", err)
		}

		delim := ""

		for _, v := range values {

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
