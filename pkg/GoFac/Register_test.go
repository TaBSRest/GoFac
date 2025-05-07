package GoFac

import (
	"testing"

	"github.com/stretchr/testify/assert"

	AsOptions "github.com/TaBSRest/GoFac/pkg/GoFac/Options/As"
	ss "github.com/TaBSRest/GoFac/tests/SampleStructs"
)

func TestRegister(t *testing.T) {
	assert := assert.New(t)

	containerBuilder := NewContainerBuilder()
	var err error
	assert.NotPanics(
		func() {
			err = RegisterConstructor(containerBuilder, ss.NewA, AsOptions.As[ss.IIndependentStruct])
		},
		"Should not have paniced when registering a constructor!",
	)

	assert.Nil(err, "No Error should have happened when registering")
}
