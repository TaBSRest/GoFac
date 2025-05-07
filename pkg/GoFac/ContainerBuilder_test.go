package GoFac_test

import (
	ctx "context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	gf "github.com/TaBSRest/GoFac/pkg/GoFac"
	AsOptions "github.com/TaBSRest/GoFac/pkg/GoFac/Options/As"
	ss "github.com/TaBSRest/GoFac/tests/SampleStructs"
)

type MockRealUUIDProvider struct{}

func (MockRealUUIDProvider) New() uuid.UUID { return uuid.New() }

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

func TestNewContainerBuilder_DoesNotBuildTwice(t *testing.T) {
	assert := assert.New(t)

	containerBuilder := gf.NewContainerBuilder()
	_, err := containerBuilder.Build(MockRealUUIDProvider{})
	assert.True(containerBuilder.IsBuilt())
	assert.Nil(err)

	_, err = containerBuilder.Build(MockRealUUIDProvider{})
	assert.NotNil(err)
}

func TestBuild_AbleToResolveContainerAndTheDependencyInIt(t *testing.T) {
	assert := assert.New(t)

	containerBuilder := gf.NewContainerBuilder()
	gf.RegisterConstructor(containerBuilder, ss.NewIndependentStruct)
	gf.RegisterConstructor(
		containerBuilder,
		func(container *gf.Container) (ss.IIndependentStruct, error) {
			return gf.Resolve[*ss.IndependentStruct](ctx.Background(), container)
		},
		AsOptions.As[ss.IIndependentStruct],
	)
	container, err := containerBuilder.Build(MockRealUUIDProvider{})
	assert.True(containerBuilder.IsBuilt())
	if err != nil {
		assert.Fail(err.Error())
	}

	sample1, err := gf.Resolve[*ss.IndependentStruct](ctx.Background(), container)
	assert.NotNil(sample1)
	if err != nil {
		assert.Fail(err.Error())
	}

	sample2, err := gf.Resolve[ss.IIndependentStruct](ctx.Background(), container)
	assert.NotNil(sample2)
	if err != nil {
		assert.Fail(err.Error())
	}
}
