package Container_test

import (
	ctx "context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/TaBSRest/GoFac/internal/RegistrationOption"
	"github.com/TaBSRest/GoFac/internal/Scope"
	mi "github.com/TaBSRest/GoFac/mocks/interfaces"
	"github.com/TaBSRest/GoFac/pkg/Container"
	"github.com/TaBSRest/GoFac/pkg/ContainerBuilder"
	RegistrationOptions "github.com/TaBSRest/GoFac/pkg/Options/Registration"
	ss "github.com/TaBSRest/GoFac/tests/SampleStructs"
)

func TestRegisterContext_ReturnsNewAndDifferentContext(t *testing.T) {
	assert := assert.New(t)

	mockUUID := uuid.New()
	mockUUIDProvider := mi.NewUUIDProvider(t)
	mockUUIDProvider.EXPECT().New().Return(mockUUID)

	containerBuilder := ContainerBuilder.New()
	containerBuilder.RegisterConstructor(ss.NewA, RegistrationOptions.As[ss.IIndependentStruct])
	container, _ := containerBuilder.Build(mockUUIDProvider)

	context := ctx.Background()
	newContext := container.RegisterContext(context)
	assert.NotEqual(context, newContext, "Should return a new and different context")
}

func TestRegisterContext_ReturnsSameContext_IfContextHasAlreadyBeenRegistered(t *testing.T) {
	assert := assert.New(t)

	mockUUID := uuid.New()
	mockUUIDProvider := mi.NewUUIDProvider(t)
	mockUUIDProvider.EXPECT().New().Return(mockUUID)

	containerBuilder := ContainerBuilder.New()
	containerBuilder.RegisterConstructor(
		ss.NewA,
		func(option *RegistrationOption.RegistrationOption) error {
			option.Scope = Scope.PerContext
			return RegistrationOptions.As[ss.IIndependentStruct](option)
		},
	)
	container, _ := containerBuilder.Build(mockUUIDProvider)

	context := ctx.Background()
	firstNewContext := container.RegisterContext(context)

	secondNewContext := container.RegisterContext(firstNewContext)
	assert.Equal(firstNewContext, secondNewContext, "Should return the same context")
}

func TestGetMetadataFromContext_ReturnsCorrectMetadataAndTrue(t *testing.T) {
	assert := assert.New(t)

	mockUUID := uuid.New()
	mockUUIDProvider := mi.NewUUIDProvider(t)
	mockUUIDProvider.EXPECT().New().Return(mockUUID)

	containerBuilder := ContainerBuilder.New()
	containerBuilder.RegisterConstructor(
		ss.NewA,
		func(option *RegistrationOption.RegistrationOption) error {
			option.Scope = Scope.PerContext
			return RegistrationOptions.As[ss.IIndependentStruct](option)
		},
	)

	container, _ := containerBuilder.Build(mockUUIDProvider)

	context := ctx.Background()
	newContext := container.RegisterContext(context)

	metadata, ok := Container.GetMetadataFromContext(newContext)
	assert.NotNil(metadata, "Should return non-nil metadata")
	assert.True(ok, "Should return true")

	assert.Len(metadata, 1, "The metadata should contain 1 registration")
	for registration, contextRegistration := range metadata {
		assert.NotNil(registration, "Registration key should not be nil")
		assert.Nil(contextRegistration.Instance, "Instance should be nil")
		assert.NotNil(contextRegistration.Once, "Once should not be nil")
	}
}

func TestGetMetadataFromContext_ReturnsNilAndFalse_IfContextHasNotBeenRegistered(t *testing.T) {
	assert := assert.New(t)

	context := ctx.Background()

	metadata, found := Container.GetMetadataFromContext(context)
	assert.Nil(metadata, "Should return nil")
	assert.False(found, "Should return false")
}
