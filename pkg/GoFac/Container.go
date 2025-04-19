package GoFac

import "reflect"

type Container struct {
	*ContainerBuilder
	SingletonCache map[reflect.Type]*reflect.Value
}
