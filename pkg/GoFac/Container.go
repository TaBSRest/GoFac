package GoFac

import (
	"sync"
)

type Container struct {
	ContainerBuilder *ContainerBuilder
	SingletonCache   sync.Map
}
