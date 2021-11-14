package configErrors

import "fmt"

type UnmarshalError struct {
	ConfigName string
	ViperErr   error
}

func (e *UnmarshalError) Error() string {
	return fmt.Sprintf("unable to decode into struct in %s config, error: %s", e.ConfigName, e.ViperErr.Error())
}
