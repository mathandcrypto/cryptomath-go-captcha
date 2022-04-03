package configErrors

import "fmt"

type ValidationError struct {
	ConfigName  string
	ValidateErr error
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("failed to validate %s config: %s", e.ConfigName, e.ValidateErr.Error())
}
