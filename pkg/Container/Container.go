package Container

import (
	"sync"
	
	i "github.com/TaBSRest/GoFac/interfaces"
)

type Container struct {
	i.ContainerBuilder
	SingletonCache *sync.Map
}

func (c *Container) GetSingletonCache() *sync.Map {
	return c.SingletonCache
}
