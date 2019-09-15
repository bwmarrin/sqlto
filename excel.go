package sqlto

import (
	"fmt"
	"io"

	"github.com/bwmarrin/xlsx"
)

func (s *SQLto) Excel(w io.Writer) error {
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

	xfile := xlsx.NewFile()
	xsheet, err := xfile.AddSheet("Sheet1")
	if err != nil {
		return fmt.Errorf("error adding sheet to xlsx file, %s\n", err)
	}

	colNames, err := s.Rows.Columns()
	if err != nil {
		return fmt.Errorf("error fetching column names, %s\n", err)
	}
	xrow := xsheet.AddRow()
	xrow.WriteSlice(&colNames, -1)

	for s.Rows.Next() {

		err = s.Rows.Scan(values...)
		if err != nil {
			return fmt.Errorf("error scanning sql row, %s\n", err)
		}

		xrow = xsheet.AddRow()

		for _, v := range values {

			xcell := xrow.AddCell()

			xcell.SetValue(v)

		}

	}

	err = xfile.Write(w)
	if err != nil {
		fmt.Printf("xlsx file.Write: error %s\n", err)
	}

	return nil
}
