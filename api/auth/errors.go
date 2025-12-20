package auth

import "fmt"

func ErrFieldMissing(fieldName string) error {
	return fmt.Errorf("missing required value: %s", fieldName)
}
