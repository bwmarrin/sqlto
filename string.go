package sqlto

import (
	"fmt"
	"io"
)

func (s *SQLto) String(w io.Writer) error {
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

	for s.Rows.Next() {
		err := s.Rows.Scan(values...)
		if err != nil {
			return err
		}

		for _, v := range values {
			if v == nil {
			} else {
				fmt.Fprintf(w, "%v", v)
			}
		}
	}
	return nil
}
