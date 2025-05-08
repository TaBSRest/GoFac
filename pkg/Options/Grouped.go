package Options

import (
	"reflect"

	gi "github.com/TaBSRest/GoFac/internal/RegistrationOption/GroupInfo"
	o "github.com/TaBSRest/GoFac/internal/RegistrationOption"
)

func Grouped[T any](groupName string) func(option *o.RegistrationOption) error {
	return func(option *o.RegistrationOption) error {
		option.RegistrationGroup = &gi.GroupInfo{
			Name: groupName,
			GroupType: reflect.TypeFor[T](),
		}
		return nil
	}
}
