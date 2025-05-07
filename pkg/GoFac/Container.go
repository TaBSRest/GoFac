package GoFac

import (
	"sync"

	"github.com/TaBSRest/GoFac/internal/BuildOption"
)

type Container struct {
	*ContainerBuilder
	BuildOption    *BuildOption.BuildOption
	SingletonCache sync.Map
}
