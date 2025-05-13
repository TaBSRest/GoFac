package GoFac

import (
	ctx "context"
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	RegistrationOption "github.com/TaBSRest/GoFac/internal/RegistrationOption"
	Scope "github.com/TaBSRest/GoFac/internal/Scope"
	mi "github.com/TaBSRest/GoFac/mocks/interfaces"
	AsOptions "github.com/TaBSRest/GoFac/pkg/GoFac/Options/As"
	ss "github.com/TaBSRest/GoFac/tests/SampleStructs"
)

var _ iGoFacUUIDWrapper = (*gofacUUIDWrapper)(nil)

func TestRegisterContext_ReturnsNewAndDifferentContext(t *testing.T) {
	assert := assert.New(t)

	mockUUID := uuid.New()
	mockUUIDProvider := mi.NewUUIDProvider(t)
	mockUUIDProvider.EXPECT().New().Return(mockUUID)

	containerBuilder := NewContainerBuilder()
	RegisterConstructor(containerBuilder, ss.NewA, AsOptions.As[ss.IIndependentStruct])
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

	containerBuilder := NewContainerBuilder()
	RegisterConstructor(
		containerBuilder,
		ss.NewA,
		func(option *RegistrationOption.RegistrationOption) error {
			option.Scope = Scope.PerContext
			return AsOptions.As[ss.IIndependentStruct](option)
		},
	)
	container, _ := containerBuilder.Build(mockUUIDProvider)

	context := ctx.Background()
	firstNewContext := container.RegisterContext(context)

	secondNewContext := container.RegisterContext(firstNewContext)
	assert.Equal(firstNewContext, secondNewContext, "Should return the same context")
}

func TestRegisterContext_SetsContextValueCorrectly(t *testing.T) {
	assert := assert.New(t)

	mockUUID := uuid.New()
	mockUUIDProvider := mi.NewUUIDProvider(t)
	mockUUIDProvider.EXPECT().New().Return(mockUUID)

	containerBuilder := NewContainerBuilder()
	RegisterConstructor(
		containerBuilder,
		ss.NewA,
		func(option *RegistrationOption.RegistrationOption) error {
			option.Scope = Scope.PerContext
			return AsOptions.As[ss.IIndependentStruct](option)
		},
	)
	container, _ := containerBuilder.Build(mockUUIDProvider)

	context := ctx.Background()
	newContext := container.RegisterContext(context)

	wrapperInNewContext := newContext.Value(gofac_UUID_WRAPPER_KEY)
	assert.NotNil(wrapperInNewContext, "The registered context should contain a non-nil wrapper under gofac_UUID_WRAPPER_KEY")

	typeInfo := reflect.TypeOf(wrapperInNewContext)
	assert.Equal("*GoFac.gofacUUIDWrapper", typeInfo.String(),
		"The wrapper in the registered context should be of type *GoFac.gofacUUIDWrapper",
	)
	contextID := wrapperInNewContext.(iGoFacUUIDWrapper).getContextID()
	assert.Equal(mockUUID.String(), contextID,
		"The wrapper in the registered context should contain the mock UUID string",
	)
}

func TestGetMetadataFromContext_ReturnsCorrectMetadataAndTrue(t *testing.T) {
	assert := assert.New(t)

	mockUUID := uuid.New()
	mockUUIDProvider := mi.NewUUIDProvider(t)
	mockUUIDProvider.EXPECT().New().Return(mockUUID)

	containerBuilder := NewContainerBuilder()
	RegisterConstructor(
		containerBuilder,
		ss.NewA,
		func(option *RegistrationOption.RegistrationOption) error {
			option.Scope = Scope.PerContext
			return AsOptions.As[ss.IIndependentStruct](option)
		},
	)

	container, _ := containerBuilder.Build(mockUUIDProvider)

	context := ctx.Background()
	newContext := container.RegisterContext(context)

	metadata, ok := GetMetadataFromContext(newContext)
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

	metadata, found := GetMetadataFromContext(context)
	assert.Nil(metadata, "Should return nil")
	assert.False(found, "Should return false")
}

func TestGetMetadataFromContext_ReturnsFalse_IfIDIsNotInRegistry(t *testing.T) {
	assert := assert.New(t)

	context := ctx.Background()
	contextWithInvalidID := ctx.WithValue(context, gofac_UUID_WRAPPER_KEY, "invalid-id")

	metadata, found := GetMetadataFromContext(contextWithInvalidID)
	assert.Nil(metadata, "Should return nil")
	assert.False(found, "Should return false")
}
