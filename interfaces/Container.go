package interfaces

import (
	ctx "context"
	"sync"
)

type Container interface {
	ContainerBuilder
	GetSingletonCache() *sync.Map
	RegisterContext(context ctx.Context) ctx.Context
}
