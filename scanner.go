package sqlto

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// GenericScanner is a somewhat hacky way of dealing with SQL database
// types that are not well represnted in Go.  For example, a DECIAML type
// has no matching Go datatype and if we do nothing, then it will be JSON
// Marshaled into a base64 string which isn't what we want for this library
// use case.
// This is all a big of "magic box" stuff here that isn't very idomatic Go
// so if someone knows a better solution to this, I'd love to hear it.
type GenericScanner struct {
	value  interface{}
	dbtype string
	want   string
}

func (scanner *GenericScanner) String() string {
	switch scanner.value.(type) {
	case []byte:
		return string(scanner.value.([]byte))
	default:
		return fmt.Sprintf("%v", scanner.value)

	}
}
func (scanner *GenericScanner) MarshalJSON() ([]byte, error) {
	return json.Marshal(scanner.value)
}

func (scanner *GenericScanner) Scan(src interface{}) error {
	switch src.(type) {
	case []byte:
		switch scanner.dbtype {
		case "DECIMAL":
			f, err := strconv.ParseFloat(string(src.([]byte)), 64)
			if err != nil {
				scanner.value = string(src.([]byte))
			} else {
				scanner.value = f
			}
		default:
			scanner.value = string(src.([]byte))
		}
	default:
		scanner.value = src

	}
	return nil
}
