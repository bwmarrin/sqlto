package sqlto

import (
	"errors"
	"fmt"
	"io"

	"github.com/bwmarrin/xlsx"
)

func (s *SQLto) Excel(w io.Writer) error {
	if s == nil {
		return errors.New("SQLto is nil")
	}

	if s.Rows == nil {
		return errors.New("SQL Rows is nil")
	}

	var err error

	colTypes, err := s.Rows.ColumnTypes()
	if err != nil {
		return fmt.Errorf("error fetching column types, %s", err)
	}

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

	xfile := xlsx.NewFile()
	xsheet, err := xfile.AddSheet("Sheet1")
	if err != nil {
		return fmt.Errorf("error adding sheet to xlsx file, %s\n", err)
	}

	xrow := xsheet.AddRow()
	xrow.WriteSlice(&colNames, -1)

	for s.Rows.Next() {

		err = s.Rows.Scan(pointers...)
		if err != nil {
			return fmt.Errorf("error scanning sql row, %s\n", err)
		}

		xrow = xsheet.AddRow()

		for k, v := range container {

			xcell := xrow.AddCell()

			ct := colTypes[k].DatabaseTypeName()
			switch ct {
			case `SQLT_NUM`:
				fallthrough
			case `INT`:
				fallthrough
			case `DECIMAL`:
				fallthrough
			case `TINYINT`:
				xcell.SetNumeric(fmt.Sprintf("%s", v))
				break
			default:
				xcell.SetValue(v)
			}

		}

	}

	err = xfile.Write(w)
	if err != nil {
		fmt.Printf("xlsx file.Write: error %s\n", err)
	}

	return nil
}
