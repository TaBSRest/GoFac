package Container

import (
	ctx "context"
	"reflect"
	"runtime"
	"sync"

	"github.com/google/uuid"

	i "github.com/TaBSRest/GoFac/interfaces"
	"github.com/TaBSRest/GoFac/internal/BuildOption"
	r "github.com/TaBSRest/GoFac/internal/Registration"
)

type contextKey string

const gofac_UUID_WRAPPER_KEY = contextKey("GoFacUUIDWrapper")

type Container struct {
	i.ContainerBuilder
	BuildOption    *BuildOption.BuildOption
	SingletonCache *sync.Map
}

func (c *Container) GetSingletonCache() *sync.Map {
	return c.SingletonCache
}

type iGoFacUUIDWrapper interface {
	getContextID() string
}

var _ iGoFacUUIDWrapper = (*goFacUUIDWrapper)(nil)

type goFacUUIDWrapper struct {
	contextID string
}

func (uw *goFacUUIDWrapper) getContextID() string {
	return uw.contextID
}

type ContextRegistration struct {
	Instance *reflect.Value
	Once     *sync.Once
}

var (
	registry sync.Map // Data Type: map[string]map[*r.Registration]*ContextRegistration
	mutex    sync.Mutex
)

func (c *Container) RegisterContext(context ctx.Context) ctx.Context {
	if c.BuildOption.IsRegisterContextRunningConcurrently {
		mutex.Lock()
		defer mutex.Unlock()
	}

	if uuidWrapper, ok := context.Value(gofac_UUID_WRAPPER_KEY).(iGoFacUUIDWrapper); ok {
		if _, exists := registry.Load(uuidWrapper.getContextID()); exists {
			return context
		}
	}

	uuidString := uuid.New().String()
	metadata := make(map[*r.Registration]*ContextRegistration)
	perContextRegistrations := c.GetPerContextRegistrations()
	for _, registration := range perContextRegistrations {
		metadata[registration] = &ContextRegistration{
			Instance: nil,
			Once:     &sync.Once{},
		}
	}
	registry.Store(uuidString, metadata)

	uuidWrapper := &goFacUUIDWrapper{contextID: uuidString}
	runtime.SetFinalizer(uuidWrapper, func(f *goFacUUIDWrapper) {
		registry.Delete(f.contextID)
	})

	context = ctx.WithValue(context, gofac_UUID_WRAPPER_KEY, uuidWrapper)

	return context
}

func GetMetadataFromContext(context ctx.Context) (map[*r.Registration]*ContextRegistration, bool) {
	uuidWrapper, ok := context.Value(gofac_UUID_WRAPPER_KEY).(iGoFacUUIDWrapper)
	if !ok {
		return nil, false
	}

	val, ok := registry.Load(uuidWrapper.getContextID())
	if !ok {
		return nil, false
	}

	return val.(map[*r.Registration]*ContextRegistration), true
}
