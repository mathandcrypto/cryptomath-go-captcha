package configErrors

import "fmt"

type ReadConfigError struct {
	ConfigName string
	ViperErr   error
}

func (e *ReadConfigError) Error() string {
	return fmt.Sprintf("failed to read %s config: %s", e.ConfigName, e.ViperErr.Error())
}
