package GoFac_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	gf "github.com/TaBSRest/GoFac/pkg/GoFac"
)

func TestContainer_Constructor_InitializedProperly(t *testing.T) {
	assert := assert.New(t)

	var gofac *gf.ContainerBuilder
	assert.NotPanics(
		func() {
			gofac = gf.NewContainerBuilder()
		},
		"Should not panic when creating new container",
	)
	assert.NotNil(gofac, "Initialized container should not be nil")
}
