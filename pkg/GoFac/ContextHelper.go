package GoFac

import (
	ctx "context"
	"reflect"
	"runtime"
	"sync"

	r "github.com/TaBSRest/GoFac/internal/Registration"
)

type contextKey string

const gofac_UUID_WRAPPER_KEY = contextKey("GoFacUUIDWrapper")

type iGoFacUUIDWrapper interface {
	getContextID() string
}

type gofacUUIDWrapper struct {
	contextID string
}

func (uw *gofacUUIDWrapper) getContextID() string {
	return uw.contextID
}

type ContextRegistration struct {
	Instance *reflect.Value
	Once     *sync.Once
}

var (
	registry sync.Map // map[string]map[*r.Registration]ContextRegistration
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

	uuidString := c.BuildOption.UUIDProvider.New().String()
	metadata := make(map[*r.Registration]*ContextRegistration)
	for _, registration := range c.perContextRegistrations {
		metadata[registration] = &ContextRegistration{
			Instance: nil,
			Once:     &sync.Once{},
		}
	}
	registry.Store(uuidString, metadata)

	uuidWrapper := &gofacUUIDWrapper{contextID: uuidString}
	runtime.SetFinalizer(uuidWrapper, func(f *gofacUUIDWrapper) {
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
