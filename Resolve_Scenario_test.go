package GoFac_test

import (
	ctx "context"
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/TaBSRest/GoFac"
	i "github.com/TaBSRest/GoFac/interfaces"
	"github.com/TaBSRest/GoFac/internal/BuildOption"
	"github.com/TaBSRest/GoFac/pkg/ContainerBuilder"
	BuildOptions "github.com/TaBSRest/GoFac/pkg/Options/Build"
	RegistrationOptions "github.com/TaBSRest/GoFac/pkg/Options/Registration"
	ss "github.com/TaBSRest/GoFac/tests/SampleStructs"
)

const num_GOROUTINES = 2500

type token string

var (
	scenarios  []scenario
	container  i.Container
	setupOnce  sync.Once
	setupError error
)

type scenario struct {
	name       string
	runResolve func(t *testing.T, container i.Container, scenarioID int)
}

func setupContainer(t *testing.T) {
	setupOnce.Do(func() {
		containerBuilder := ContainerBuilder.New()
		var registerError error

		registerError = containerBuilder.Register(
			ss.NewA,
			RegistrationOptions.As[ss.IIndependentStruct],
			RegistrationOptions.Named("Named_A"),
		)
		if registerError != nil {
			setupError = registerError
			return
		}

		registerError = containerBuilder.Register(
			ss.NewS,
			RegistrationOptions.As[ss.IIndependentStruct],
			RegistrationOptions.AsSingleton,
			RegistrationOptions.Named("Singleton_S"),
		)
		if registerError != nil {
			setupError = registerError
			return
		}

		registerError = containerBuilder.Register(
			ss.NewA,
			RegistrationOptions.As[ss.IIndependentStruct],
			RegistrationOptions.PerContext,
			RegistrationOptions.Named("PerContext_A"),
		)
		if registerError != nil {
			setupError = registerError
			return
		}

		registerError = containerBuilder.Register(
			ss.NewB,
			RegistrationOptions.As[ss.IIndependentStruct],
			RegistrationOptions.PerContext,
			RegistrationOptions.Named("PerContext_B"),
		)
		if registerError != nil {
			setupError = registerError
			return
		}

		registerError = containerBuilder.Register(
			func() ss.ISpecificPerContextStruct {
				return ss.NewSpecificPerContextStruct(ss.NewA())
			},
			RegistrationOptions.As[ss.ISpecificPerContextStruct],
			RegistrationOptions.PerContext,
			RegistrationOptions.Named("Specific_PerContext_A_Instance"),
		)
		if registerError != nil {
			setupError = registerError
			return
		}

		registerError = containerBuilder.Register(
			ss.NewA,
			RegistrationOptions.As[ss.IIndependentStruct],
			RegistrationOptions.Grouped[ss.IIndependentStruct]("Group"),
		)
		if registerError != nil {
			setupError = registerError
			return
		}

		registerError = containerBuilder.Register(
			ss.NewB,
			RegistrationOptions.As[ss.IIndependentStruct],
			RegistrationOptions.Grouped[ss.IIndependentStruct]("Group"),
		)
		if registerError != nil {
			setupError = registerError
			return
		}

		registerError = containerBuilder.Register(
			ss.NewStructRelyingOnIndependentStructs,
			RegistrationOptions.As[ss.IStructRelyingOnIndependentStructs],
			RegistrationOptions.PerContext,
			RegistrationOptions.Named("ArrayDependencyPerContextParent"),
		)
		if registerError != nil {
			setupError = registerError
			return
		}

		registerError = containerBuilder.Register(
			func(c i.Container) (ss.IStructRelyingOnIndependentStruct, error) {
				singletonDependency, err := GoFac.ResolveNamed[ss.IIndependentStruct](c, ctx.Background(), "Singleton_S")
				if err != nil {
					return nil, err
				}
				return ss.NewStructRelyingOnIndependentStruct(singletonDependency), nil
			},
			RegistrationOptions.As[ss.IStructRelyingOnIndependentStruct],
			RegistrationOptions.PerContext,
			RegistrationOptions.Named("PerContextWithSingletonDependency"),
		)
		if registerError != nil {
			setupError = registerError
			return
		}

		registerError = containerBuilder.Register(
			func(c i.Container, specificPerContextDependency ss.ISpecificPerContextStruct) (ss.ISingletonWithPerContextDependencyStruct, error) {
				singletonMain, err := GoFac.ResolveNamed[ss.IIndependentStruct](c, ctx.Background(), "Singleton_S")
				if err != nil {
					return nil, err
				}
				return ss.NewSingletonWithSpecificPerContextDependencyStruct(singletonMain, specificPerContextDependency), nil
			},
			RegistrationOptions.As[ss.ISingletonWithPerContextDependencyStruct],
			RegistrationOptions.AsSingleton,
			RegistrationOptions.Named("SingletonWithSpecificPerContextDependency"),
		)
		if registerError != nil {
			setupError = registerError
			return
		}

		var buildOptionFunction func(*BuildOption.BuildOption)
		if rand.Intn(2) == 0 {
			buildOptionFunction = BuildOptions.RegisterSameContextConcurrently
		}

		var buildContainer i.Container
		var buildError error

		if buildOptionFunction != nil {
			buildContainer, buildError = containerBuilder.Build(buildOptionFunction)
		} else {
			buildContainer, buildError = containerBuilder.Build()
		}

		container = buildContainer
		setupError = buildError
	})

	if setupError != nil {
		t.Fatalf("Failed to setup container: %v", setupError)
	}
}

func initializeScenarios(t *testing.T) {
	assert := assert.New(t)

	scenarios = []scenario{
		{
			name: "Resolve Simple",
			runResolve: func(t *testing.T, c i.Container, scenarioID int) {
				instance, err := GoFac.Resolve[ss.IIndependentStruct](c, ctx.Background())
				assert.Nil(err)
				assert.NotNil(instance)
				assert.Equal("IndependentStructB", instance.ReturnNameIndependentStruct())
			},
		},
		{
			name: "Resolve Named",
			runResolve: func(t *testing.T, c i.Container, scenarioID int) {
				instance, err := GoFac.ResolveNamed[ss.IIndependentStruct](c, ctx.Background(), "Named_A")
				assert.Nil(err)
				assert.NotNil(instance)
				assert.Equal("IndependentStruct", instance.ReturnNameIndependentStruct())
			},
		},
		{
			name: "Resolve Singleton",
			runResolve: func(t *testing.T, c i.Container, scenarioID int) {
				instance1, err := GoFac.ResolveNamed[ss.IIndependentStruct](c, ctx.Background(), "Singleton_S")
				assert.Nil(err)
				assert.NotNil(instance1)

				instance2, err := GoFac.ResolveNamed[ss.IIndependentStruct](c, ctx.Background(), "Singleton_S")
				assert.Nil(err)
				assert.NotNil(instance2)

				assert.Same(instance1, instance2, "Singleton instances should be the same")
			},
		},
		{
			name: "Resolve PerContext",
			runResolve: func(t *testing.T, c i.Container, scenarioID int) {
				contextKey1 := token(fmt.Sprintf("ContextKey1_Scenario%d", scenarioID))
				contextKey2 := token(fmt.Sprintf("ContextKey2_Scenario%d", scenarioID))
				contextKey3 := token(fmt.Sprintf("ContextKey3_Scenario%d", scenarioID))

				unregisteredContext1 := ctx.WithValue(ctx.Background(), contextKey1, "ContextValue1")
				unregisteredContext2 := ctx.WithValue(ctx.Background(), contextKey2, "ContextValue2")
				unregisteredContext3 := ctx.WithValue(ctx.Background(), contextKey3, "ContextValue3")

				registeredContext1 := c.RegisterContext(unregisteredContext1)
				registeredContext2 := c.RegisterContext(unregisteredContext2)

				instance1a, err := GoFac.ResolveNamed[ss.IIndependentStruct](c, registeredContext1, "PerContext_A")
				assert.Nil(err)
				assert.NotNil(instance1a)

				instance1b, err := GoFac.ResolveNamed[ss.IIndependentStruct](c, registeredContext1, "PerContext_A")
				assert.Nil(err)
				assert.NotNil(instance1b)

				assert.Same(instance1a, instance1b, "Instances from the same context should be the same (PerContext_A)")

				instance2, err := GoFac.ResolveNamed[ss.IIndependentStruct](c, registeredContext2, "PerContext_A")
				assert.Nil(err)
				assert.NotNil(instance2)

				assert.NotSame(instance1a, instance2, "Instances from different contexts should not be the same (PerContext_A)")
				_, err = GoFac.ResolveNamed[ss.IIndependentStruct](c, unregisteredContext3, "PerContext_A")
				assert.NotNil(err)
				assert.Contains(err.Error(), "GoFac.resolvePerContext: The context is not registered to GoFac")
			},
		},
		{
			name: "Resolve Array with Registered Context",
			runResolve: func(t *testing.T, c i.Container, scenarioID int) {
				contextKey := token(fmt.Sprintf("ContextKey_Scenario%d", scenarioID))
				registeredContext := c.RegisterContext(ctx.WithValue(ctx.Background(), contextKey, "ContextValue"))

				instances, err := GoFac.ResolveMultiple[ss.IIndependentStruct](c, registeredContext)
				assert.Nil(err)
				assert.NotNil(instances)

				expectedNumberOfInstances := 6
				assert.Len(instances, expectedNumberOfInstances, fmt.Sprintf("Should resolve all %d IIndependentStruct instances", expectedNumberOfInstances))
			},
		},
		{
			name: "Resolve Array with Unregistered Context Expecting Error",
			runResolve: func(t *testing.T, c i.Container, scenarioID int) {
				_, err := GoFac.ResolveMultiple[ss.IIndependentStruct](c, ctx.Background())
				assert.NotNil(err)
				assert.Contains(err.Error(), "GoFac.resolvePerContext: The context is not registered to GoFac")
			},
		},
		{
			name: "Resolve Group",
			runResolve: func(t *testing.T, c i.Container, scenarioID int) {
				instances, err := GoFac.ResolveGroup[ss.IIndependentStruct](c, ctx.Background(), "Group")
				assert.Nil(err)
				assert.Len(instances, 2)

				var names []string
				for _, instance := range instances {
					assert.NotNil(instance)
					names = append(names, instance.ReturnNameIndependentStruct())
				}
				assert.Contains(names, "IndependentStruct")
				assert.Contains(names, "IndependentStructB")
			},
		},
		{
			name: "Resolve PerContext Parent with Array of Context-Aware Items",
			runResolve: func(t *testing.T, c i.Container, scenarioID int) {
				contextValue1 := fmt.Sprintf("ContextValue1_Scenario%d", scenarioID)
				contextValue2 := fmt.Sprintf("ContextValue2_Scenario%d", scenarioID)

				registeredContext1 := c.RegisterContext(ctx.WithValue(ctx.Background(), token("ContextKey"), contextValue1))
				registeredContext2 := c.RegisterContext(ctx.WithValue(ctx.Background(), token("ContextKey"), contextValue2))

				parent1a, err1a := GoFac.ResolveNamed[ss.IStructRelyingOnIndependentStructs](c, registeredContext1, "ArrayDependencyPerContextParent")
				assert.Nil(err1a)

				parent1b, err1b := GoFac.ResolveNamed[ss.IStructRelyingOnIndependentStructs](c, registeredContext1, "ArrayDependencyPerContextParent")
				assert.Nil(err1b)

				assert.Same(parent1a, parent1b, "Parent (ArrayDependencyPerContextParent) from same context should be same")

				parent2, err2 := GoFac.ResolveNamed[ss.IStructRelyingOnIndependentStructs](c, registeredContext2, "ArrayDependencyPerContextParent")
				assert.Nil(err2)

				assert.NotSame(parent1a, parent2, "Parent (ArrayDependencyPerContextParent) from different contexts should be different")
			},
		},
		{
			name: "Resolve PerContext with Singleton Dependency",
			runResolve: func(t *testing.T, c i.Container, scenarioID int) {
				contextValue1 := fmt.Sprintf("ContextValue1_Scenario%d", scenarioID)
				contextValue2 := fmt.Sprintf("ContextValue2_Scenario%d", scenarioID)

				registeredContext1 := c.RegisterContext(ctx.WithValue(ctx.Background(), token("ContextKey"), contextValue1))
				registeredContext2 := c.RegisterContext(ctx.WithValue(ctx.Background(), token("ContextKey"), contextValue2))

				instance1a, err1a := GoFac.ResolveNamed[ss.IStructRelyingOnIndependentStruct](c, registeredContext1, "PerContextWithSingletonDependency")
				assert.Nil(err1a)

				instance1b, err1b := GoFac.ResolveNamed[ss.IStructRelyingOnIndependentStruct](c, registeredContext1, "PerContextWithSingletonDependency")
				assert.Nil(err1b)

				assert.Same(instance1a, instance1b, "PerContext instance with Singleton dependency from same context should be same")

				instance2, err2 := GoFac.ResolveNamed[ss.IStructRelyingOnIndependentStruct](c, registeredContext2, "PerContextWithSingletonDependency")
				assert.Nil(err2)

				assert.NotSame(instance1a, instance2, "PerContext service with Singleton dependency from different contexts should be different")
			},
		},
		{
			name: "Resolve Singleton with Specific PerContext Dependency",
			runResolve: func(t *testing.T, c i.Container, scenarioID int) {
				contextValue1 := fmt.Sprintf("s%d_cValA_gr%d", scenarioID, rand.Int())
				contextValue2 := fmt.Sprintf("s%d_cValB_gr%d", scenarioID, rand.Int())
				registeredContext1 := c.RegisterContext(ctx.WithValue(ctx.Background(), token("cKey1"), contextValue1))
				registeredContext2 := c.RegisterContext(ctx.WithValue(ctx.Background(), token("cKey2"), contextValue2))

				var firstInstance ss.ISingletonWithPerContextDependencyStruct
				var initialPerContextDependencyID string

				instance1, err := GoFac.ResolveNamed[ss.ISingletonWithPerContextDependencyStruct](c, registeredContext1, "SingletonWithSpecificPerContextDependency")
				assert.Nil(err)
				assert.NotNil(instance1)
				firstInstance = instance1
				initialPerContextDependencyID = instance1.GetInitialPerContextDepSpecificID()
				assert.NotEmpty(initialPerContextDependencyID, "Initial PerContext Depndency ID should be captured")

				instance1Again, err := GoFac.ResolveNamed[ss.ISingletonWithPerContextDependencyStruct](c, registeredContext1, "SingletonWithSpecificPerContextDependency")
				assert.Nil(err)
				assert.Same(firstInstance, instance1Again, "Singleton instance should be same on subsequent calls with same context")
				assert.Equal(initialPerContextDependencyID, instance1Again.GetInitialPerContextDepSpecificID(), "Captured PC dep ID should be consistent")

				instance2, err := GoFac.ResolveNamed[ss.ISingletonWithPerContextDependencyStruct](c, registeredContext2, "SingletonWithSpecificPerContextDependency")
				assert.Nil(err)
				assert.Same(firstInstance, instance2, "Singleton instance should be same even with different context")
				assert.Equal(initialPerContextDependencyID, instance2.GetInitialPerContextDepSpecificID(), "Captured PC dep ID should remain from the very first resolution")
			},
		},
	}
}

func TestResolve_DoesNotOccurRaceCondition(t *testing.T) {
	rand.New(rand.NewSource(time.Now().UnixNano()))

	initializeScenarios(t)

	if len(scenarios) == 0 {
		t.Fatal("No scenarios defined.")
	}

	setupContainer(t)

	if container == nil {
		t.Fatalf("Container was not built. Setup error: %v", setupError)
	}

	var waitGroup sync.WaitGroup
	waitGroup.Add(num_GOROUTINES)

	for i := 0; i < num_GOROUTINES; i++ {
		go func(goroutineID int) {
			defer waitGroup.Done()

			scenarioIndex := rand.Intn(len(scenarios))
			selectedScenario := scenarios[scenarioIndex]

			selectedScenario.runResolve(t, container, goroutineID)
		}(i)
	}

	waitGroup.Wait()

	t.Logf("Completed %d concurrent resolutions across %d scenarios.", num_GOROUTINES, len(scenarios))
}
