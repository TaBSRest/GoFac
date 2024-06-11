package GoFac

import "reflect"

type Container struct {
	*ContainerBuilder
	SingletonCache map[string]*reflect.Value
}
