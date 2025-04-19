package GoFac

import (
	ctx "context"
	"reflect"
	"sync"

	r "github.com/TaBSRest/GoFac/internal/Registration"
)

type contextKey string

const (
	GOFAC_CACHE_KEY = contextKey("GoFacCache")
	GOFAC_ONCES_KEY = contextKey("GoFacOnces")
)

func GetContextWithGoFacCache(context ctx.Context) ctx.Context {
	if !doesContextHaveGoFacCache(context) {
		context = ctx.WithValue(context, GOFAC_CACHE_KEY, make(map[*r.Registration]*reflect.Value))
	}

	return context
}

func GetContextWithGoFacOnces(context ctx.Context) ctx.Context {
	if !doesContextHaveGoFacOnces(context) {
		context = ctx.WithValue(context, GOFAC_ONCES_KEY, make(map[*r.Registration]*sync.Once))
	}

	return context
}

func doesContextHaveGoFacCache(context ctx.Context) bool {
	return context.Value(GOFAC_CACHE_KEY) != nil
}

func doesContextHaveGoFacOnces(context ctx.Context) bool {
	return context.Value(GOFAC_ONCES_KEY) != nil
}
