package interfaces

import "sync"

type Container interface {
	ContainerBuilder
	GetSingletonCache() *sync.Map
}
