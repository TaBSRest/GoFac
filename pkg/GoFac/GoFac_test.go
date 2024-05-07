package gofac

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Constructor_InitializedProperly(t *testing.T) {
	assert := assert.New(t)

	var gofac *Container
	assert.NotPanics(
		func() {
			gofac = NewContainer()
		},
		"Should not panic when creating new container",
	)
	assert.NotNil(gofac, "Initialized container should not be nil")
	assert.NotNil(gofac.cache, "The container's cache should not be nil")
}


