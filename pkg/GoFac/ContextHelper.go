package GoFac

import (
	ctx "context"
	"reflect"
	"sync"

	r "github.com/TaBSRest/GoFac/internal/Registration"
)

type ContextCache struct {
	Cache *reflect.Value
	Once  *sync.Once
}

type contextKey string

const GOFAC_KEY = contextKey("GoFac")

func GetContextWithGoFacContextCache(context ctx.Context) ctx.Context {
	if !isContextRegisteredToGoFac(context) {
		context = ctx.WithValue(context, GOFAC_KEY, make(map[*r.Registration]ContextCache))
	}

	return context
}

func isContextRegisteredToGoFac(context ctx.Context) bool {
	return context.Value(GOFAC_KEY) != nil
}
