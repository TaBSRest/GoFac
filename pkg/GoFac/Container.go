package GoFac

import (
	"sync"
)

type Container struct {
	*ContainerBuilder
	SingletonCache sync.Map
}
