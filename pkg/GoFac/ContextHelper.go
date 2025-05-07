package GoFac

import (
	ctx "context"
	"reflect"
	"runtime"
	"sync"

	r "github.com/TaBSRest/GoFac/internal/Registration"
)

type ContextRegistration struct {
	Instance *reflect.Value
	Once     *sync.Once
}

type contextIDFinalizer struct {
	contextID string
}

type contextKey string

const (
	GOFAC_CONTEXT_ID_KEY = contextKey("GoFacContextID")
	GOFAC_FINALIZER_KEY  = contextKey("GoFacFinalizer")
)

var (
	registry sync.Map // map[string]map[*r.Registration]ContextRegistration
	mutex    sync.Mutex
)

func (c *Container) RegisterContext(context ctx.Context) ctx.Context {
	if c.BuildOption.IsRegisterContextRunningConcurrently {
		mutex.Lock()
		defer mutex.Unlock()
	}

	id, found := context.Value(GOFAC_CONTEXT_ID_KEY).(string)
	_, exists := registry.Load(id)
	if found && exists {
		return context
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
	context = ctx.WithValue(context, GOFAC_CONTEXT_ID_KEY, uuidString)

	runtime.SetFinalizer(context, func(f *contextIDFinalizer) {
		registry.Delete(f.contextID)
	})

	// finalizer := &contextIDFinalizer{contextID: uuidString}
	// runtime.SetFinalizer(finalizer, func(f *contextIDFinalizer) {
	// 	registry.Delete(f.contextID)
	// })

	// context = ctx.WithValue(context, GOFAC_CONTEXT_ID_KEY, uuidString)
	// context = ctx.WithValue(context, GOFAC_FINALIZER_KEY, finalizer)

	return context
}

func GetMetadataFromContext(context ctx.Context) (map[*r.Registration]*ContextRegistration, bool) {
	id, found := context.Value(GOFAC_CONTEXT_ID_KEY).(string)
	if !found {
		return nil, false
	}

	val, found := registry.Load(id)
	if !found {
		return nil, false
	}

	return val.(map[*r.Registration]*ContextRegistration), true
}
