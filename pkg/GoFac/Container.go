package GoFac

import (
	"reflect"

	r "github.com/TaBSRest/GoFac/internal/Registration"
)

type Container struct {
	*ContainerBuilder
	SingletonCache map[*r.Registration]*reflect.Value
}
