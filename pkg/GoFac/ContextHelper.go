package GoFac

import (
	ctx "context"
	"reflect"
	"runtime"
	"sync"

	"github.com/google/uuid"

	r "github.com/TaBSRest/GoFac/internal/Registration"
	// i "github.com/TaBSRest/TaBSCore/interfaces"
)

type ContextRegistration struct {
	Instance *reflect.Value
	Once     *sync.Once
}

type contextKey string

const (
	GOFAC_CONTEXT_ID_KEY = contextKey("GoFacContextID")
	GOFAC_FINALIZER_KEY  = contextKey("GoFacFinalizer")
)

var registry sync.Map // map[string]map[*r.Registration]ContextRegistration

type contextIDFinalizer struct {
	contextID string
}

func RegisterContextToGoFac(context ctx.Context) ctx.Context {
	if id, found := context.Value(GOFAC_CONTEXT_ID_KEY).(string); found {
		if _, found := registry.Load(id); found {
			return context
		}
	}

	// uuidString := i.uuidProvider.New().String()
	uuidString := uuid.NewString()

	registry.Store(uuidString, make(map[*r.Registration]ContextRegistration))

	finalizer := &contextIDFinalizer{contextID: uuidString}
	runtime.SetFinalizer(finalizer, func(f *contextIDFinalizer) {
		registry.Delete(f.contextID)
	})

	context = ctx.WithValue(context, GOFAC_CONTEXT_ID_KEY, uuidString)
	context = ctx.WithValue(context, GOFAC_FINALIZER_KEY, finalizer)

	return context
}

func GetMetadataFromContext(context ctx.Context) (map[*r.Registration]ContextRegistration, bool) {
	id, found := context.Value(GOFAC_CONTEXT_ID_KEY).(string)
	if !found {
		return nil, false
	}

	val, found := registry.Load(id)
	if !found {
		return nil, false
	}

	return val.(map[*r.Registration]ContextRegistration), true
}

// Optional: Explicit Cleanup
func CleanupContextMetadata(context ctx.Context) {
	id, found := context.Value(GOFAC_CONTEXT_ID_KEY).(string)
	if found {
		registry.Delete(id)
	}
}
