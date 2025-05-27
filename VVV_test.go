package GoFac_test

// import (
// 	ctx "context"
// 	"fmt"
// 	"math/rand"
// 	"sync"
// 	"testing"
// 	"time"

// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/require"

// 	"github.com/TaBSRest/GoFac"
// 	i "github.com/TaBSRest/GoFac/interfaces"
// 	"github.com/TaBSRest/GoFac/internal/BuildOption"
// 	"github.com/TaBSRest/GoFac/pkg/ContainerBuilder"
// 	BuildOptions "github.com/TaBSRest/GoFac/pkg/Options/Build"
// 	RegistrationOptions "github.com/TaBSRest/GoFac/pkg/Options/Registration"
// 	ss "github.com/TaBSRest/GoFac/tests/SampleStructs"
// )

// const (
// 	num_GOROUTINES = 5000
// )

// // stressScenario defines a test case for stress testing.
// type stressScenario struct {
// 	name        string
// 	description string
// 	setupAndRun func(t *testing.T, scenarioID int) // scenarioID can be used for unique naming if needed
// }

// // Global list of scenarios
// var stressScenarios []stressScenario

// // Helper to get context ID from our wrapper (from pkg/Container/Container.go, simplified for test)
// type iGoFacUUIDWrapper interface {
// 	GetContextID() string
// }

// func getContextIDFromWrapper(t *testing.T, context ctx.Context) string {
// 	// Note: This key is defined in pkg/Container/Container.go and not exported.
// 	// For a real test utility, you might need to expose it or use a specific testing hook.
// 	// For this example, we'll assume a way to get it or simulate it.
// 	// Let's try to access it via its known string representation if your wrapper allows it.
// 	// This is a simplification; in practice, you'd want a robust way to get this.
// 	// We'll simulate by just returning a unique string for the context for now if direct access is hard.

// 	// Proper way would be to call an exported helper or make the key available for testing.
// 	// If your `gofac_UUID_WRAPPER_KEY` is `type contextKey string`, then:
// 	// wrapperVal := context.Value(contextKey("GoFacUUIDWrapper"))

// 	// For now, let's assume a simplified access or a passed-in ID for testing.
// 	// This part is tricky without modifying the original code to expose the key or a helper.
// 	// For the purpose of this test, we'll pass a generated ID for contexts.

// 	// If you have a way to retrieve the wrapper:
// 	// if wrapper, ok := context.Value(SOME_EXPORTED_KEY_FOR_TESTING).(iGoFacUUIDWrapper); ok {
// 	// 	return wrapper.GetContextID()
// 	// }
// 	// require.Fail(t, "Could not get context ID from wrapper")
// 	// return ""

// 	// Simplified: If context has a value set by us for testing ID
// 	id := context.Value("testContextID")
// 	if idStr, ok := id.(string); ok {
// 		return idStr
// 	}
// 	// Fallback if not set, though for PerContext tests, we must ensure it's set.
// 	return fmt.Sprintf("unknownCtx-%p", context)
// }

// func initStressScenarios(t *testing.T) {
// 	stressScenarios = []stressScenario{
// 		{
// 			name:        "ResolveSimple",
// 			description: "Registers and resolves a simple struct as an interface.",
// 			setupAndRun: func(t *testing.T, scenarioID int) {
// 				cb := ContainerBuilder.New()
// 				err := cb.Register(ss.NewA, RegistrationOptions.As[ss.IIndependentStruct])
// 				require.NoError(t, err)

// 				container, err := cb.Build(&i.RealUUIDProvider{})
// 				require.NoError(t, err)

// 				instance, err := GoFac.Resolve[ss.IIndependentStruct](ctx.Background(), container)
// 				require.NoError(t, err)
// 				require.NotNil(t, instance)
// 				assert.Equal(t, "IndependentStruct", instance.ReturnNameIndependentStruct())
// 			},
// 		},
// 		{
// 			name:        "ResolveNamed",
// 			description: "Registers with a name and resolves it.",
// 			setupAndRun: func(t *testing.T, scenarioID int) {
// 				cb := ContainerBuilder.New()
// 				serviceName := fmt.Sprintf("namedService_%d", scenarioID)
// 				err := cb.Register(ss.NewB, RegistrationOptions.Named(serviceName), RegistrationOptions.As[ss.IIndependentStruct])
// 				require.NoError(t, err)

// 				container, err := cb.Build(&i.RealUUIDProvider{})
// 				require.NoError(t, err)

// 				instance, err := GoFac.ResolveNamed[ss.IIndependentStruct](ctx.Background(), container, serviceName)
// 				require.NoError(t, err)
// 				require.NotNil(t, instance)
// 				assert.Equal(t, "IndependentStructB", instance.ReturnNameIndependentStruct())
// 			},
// 		},
// 		{
// 			name:        "ResolveSingleton",
// 			description: "Registers as singleton and verifies same instance.",
// 			setupAndRun: func(t *testing.T, scenarioID int) {
// 				cb := ContainerBuilder.New()
// 				err := cb.Register(ss.NewS, RegistrationOptions.As[ss.IIndependentStruct], RegistrationOptions.AsSingleton)
// 				require.NoError(t, err)

// 				container, err := cb.Build(&i.RealUUIDProvider{})
// 				require.NoError(t, err)

// 				instance1, err := GoFac.Resolve[ss.IIndependentStruct](ctx.Background(), container)
// 				require.NoError(t, err)
// 				require.NotNil(t, instance1)
// 				_ = instance1.ReturnNameIndependentStruct() // Call to potentially init internal sync.Once of NewS

// 				instance2, err := GoFac.Resolve[ss.IIndependentStruct](ctx.Background(), container)
// 				require.NoError(t, err)
// 				require.NotNil(t, instance2)

// 				assert.Same(t, instance1, instance2, "Singleton instances should be the same")
// 			},
// 		},
// 		{
// 			name:        "ResolvePerContext",
// 			description: "Registers as PerContext and verifies instance scoping.",
// 			setupAndRun: func(t *testing.T, scenarioID int) {
// 				cb := ContainerBuilder.New()
// 				// Optionally, make concurrent registration of context itself a part of the test
// 				// useConcurrentCtxReg := rand.Intn(2) == 0
// 				err := cb.Register(ss.NewA, RegistrationOptions.As[ss.IIndependentStruct], RegistrationOptions.PerContext)
// 				require.NoError(t, err)

// 				// In the "ResolvePerContext" scenario's setupAndRun function:
// 				var buildOptFunc func(*BuildOption.BuildOption)
// 				if rand.Intn(2) == 0 {
// 					buildOptFunc = BuildOptions.RegisterSameContextConcurrently
// 				}

// 				var container i.Container // Adjusted to use the actual type of container if it's from pkg/Container
// 				// Or more generally: var container i.Container
// 				// Assuming your cb.Build returns *Container.Container which implements i.Container

// 				if buildOptFunc != nil {
// 					// If you aliased "github.com/TaBSRest/GoFac/pkg/Container" as `pkgContainer` for example:
// 					// containerTyped, buildErr := cb.Build(&i.RealUUIDProvider{}, buildOptFunc)
// 					// Or just get the interface:
// 					var buildErr error
// 					container, buildErr = cb.Build(&i.RealUUIDProvider{}, buildOptFunc)
// 					err = buildErr

// 				} else {
// 					var buildErr error
// 					container, buildErr = cb.Build(&i.RealUUIDProvider{}) // Call without the (nil) config function
// 					err = buildErr
// 				}
// 				require.NoError(t, err, "Container build failed")
// 				require.NotNil(t, container, "Container should not be nil after build")

// 				// Create distinct contexts for this test
// 				// We use context.WithValue to attach a "testContextID" for reliable retrieval in getContextIDFromWrapper
// 				// In a real scenario, the container.RegisterContext embeds its own UUID.
// 				baseCtx1 := ctx.WithValue(ctx.Background(), "testContextID", fmt.Sprintf("ctx1_scenario%d", scenarioID))
// 				baseCtx2 := ctx.WithValue(ctx.Background(), "testContextID", fmt.Sprintf("ctx2_scenario%d", scenarioID))

// 				instrumentedCtx1 := container.RegisterContext(baseCtx1)
// 				instrumentedCtx2 := container.RegisterContext(baseCtx2)
// 				uninstrumentedCtx := ctx.WithValue(ctx.Background(), "testContextID", fmt.Sprintf("ctx3_scenario%d", scenarioID))

// 				// Instance 1a from context 1
// 				inst1a, err := GoFac.Resolve[ss.IIndependentStruct](instrumentedCtx1, container)
// 				require.NoError(t, err)
// 				require.NotNil(t, inst1a)

// 				// Instance 1b from context 1 (should be same as 1a)
// 				inst1b, err := GoFac.Resolve[ss.IIndependentStruct](instrumentedCtx1, container)
// 				require.NoError(t, err)
// 				require.NotNil(t, inst1b)
// 				assert.Same(t, inst1a, inst1b, "Instances from the same registered context should be the same")

// 				// Instance 2 from context 2 (should be different from 1a)
// 				inst2, err := GoFac.Resolve[ss.IIndependentStruct](instrumentedCtx2, container)
// 				require.NoError(t, err)
// 				require.NotNil(t, inst2)
// 				assert.NotSame(t, inst1a, inst2, "Instances from different registered contexts should be different")

// 				// Attempt with uninstrumented context
// 				_, err = GoFac.Resolve[ss.IIndependentStruct](uninstrumentedCtx, container)
// 				assert.Error(t, err, "Resolving with an uninstrumented context should fail")
// 			},
// 		},
// 		{
// 			name:        "ResolveArray",
// 			description: "Registers multiple services for an interface and resolves them as a slice.",
// 			setupAndRun: func(t *testing.T, scenarioID int) {
// 				cb := ContainerBuilder.New()
// 				require.NoError(t, cb.Register(ss.NewA, RegistrationOptions.As[ss.IIndependentStruct]))
// 				require.NoError(t, cb.Register(ss.NewB, RegistrationOptions.As[ss.IIndependentStruct]))

// 				container, err := cb.Build(&i.RealUUIDProvider{})
// 				require.NoError(t, err)

// 				instances, err := GoFac.ResolveMultiple[ss.IIndependentStruct](ctx.Background(), container)
// 				require.NoError(t, err)
// 				require.Len(t, instances, 2)

// 				var names []string
// 				for _, inst := range instances {
// 					require.NotNil(t, inst)
// 					names = append(names, inst.ReturnNameIndependentStruct())
// 				}
// 				assert.Contains(t, names, "IndependentStruct")
// 				assert.Contains(t, names, "IndependentStructB")
// 			},
// 		},
// 		{
// 			name:        "ResolveGroup",
// 			description: "Registers services to a group and resolves the group.",
// 			setupAndRun: func(t *testing.T, scenarioID int) {
// 				cb := ContainerBuilder.New()
// 				groupName := fmt.Sprintf("group_%d", scenarioID)
// 				require.NoError(t, cb.Register(ss.NewA, RegistrationOptions.As[ss.IIndependentStruct], RegistrationOptions.Grouped[ss.IIndependentStruct](groupName)))
// 				require.NoError(t, cb.Register(ss.NewB, RegistrationOptions.As[ss.IIndependentStruct], RegistrationOptions.Grouped[ss.IIndependentStruct](groupName)))
// 				// Add one not in the group
// 				require.NoError(t, cb.Register(ss.NewS, RegistrationOptions.As[ss.IIndependentStruct]))

// 				container, err := cb.Build(&i.RealUUIDProvider{})
// 				require.NoError(t, err)

// 				instances, err := GoFac.ResolveGroup[ss.IIndependentStruct](ctx.Background(), container, groupName)
// 				require.NoError(t, err)
// 				require.Len(t, instances, 2)

// 				var names []string
// 				for _, inst := range instances {
// 					require.NotNil(t, inst)
// 					names = append(names, inst.ReturnNameIndependentStruct())
// 				}
// 				assert.Contains(t, names, "IndependentStruct")
// 				assert.Contains(t, names, "IndependentStructB")
// 				assert.NotContains(t, names, "SingletonStruct")
// 			},
// 		},
// 		// TODO: Add more complex scenarios as requested by the team lead.
// 		// Examples:
// 		// - ResolveWithArrayDependencyThatHasContext:
// 		//   - ServiceX (PerContext) implements IServiceX
// 		//   - ServiceY (PerContext) implements IServiceX
// 		//   - ServiceZ(dep []IServiceX)
// 		// - ResolveWithArrayDependencyThatHasNamed
// 		// - ResolveWithArrayDependencyThatHasContextThatHasSimpleDependency
// 		// - ResolveContextWithSingleton (A PerContext service depending on a Singleton service)
// 		// - ResolveSingletonWithContext (A Singleton service depending on a PerContext service - this is tricky and often an anti-pattern, how would GoFac handle the context for the singleton's dependency?)
// 	}
// }

// // Sample struct for more complex dependency tests
// type IServiceX interface{ ID() string }
// type ServiceX struct{ id string }

// func (s *ServiceX) ID() string         { return s.id }
// func NewServiceXSimple() IServiceX     { return &ServiceX{id: "ServiceX_Simple_" + NextStr()} } // Assuming NextStr() provides unique strings
// func NewServiceXPerContext() IServiceX { return &ServiceX{id: "ServiceX_PerContext_" + NextStr()} }
// func NewServiceXSingleton() IServiceX  { return &ServiceX{id: "ServiceX_Singleton_" + NextStr()} }

// type IServiceY interface{ Val() string }
// type ServiceY struct{ val string }

// func (s *ServiceY) Val() string { return s.val }
// func NewServiceYWithDepX(depX IServiceX) IServiceY {
// 	return &ServiceY{val: "ServiceY_With_" + depX.ID()}
// }

// func addComplexScenarios(t *testing.T) {
// 	// Example: PerContext service depending on a Singleton service
// 	stressScenarios = append(stressScenarios, stressScenario{
// 		name:        "ResolvePerContextWithSingletonDep",
// 		description: "PerContext service depending on a Singleton service.",
// 		setupAndRun: func(t *testing.T, scenarioID int) {
// 			cb := ContainerBuilder.New()
// 			// Singleton dependency
// 			require.NoError(t, cb.Register(NewServiceXSingleton, RegistrationOptions.As[IServiceX], RegistrationOptions.AsSingleton))
// 			// PerContext service that depends on IServiceX
// 			require.NoError(t, cb.Register(NewServiceYWithDepX, RegistrationOptions.As[IServiceY], RegistrationOptions.PerContext))

// 			container, err := cb.Build(&i.RealUUIDProvider{})
// 			require.NoError(t, err)

// 			baseCtx1 := ctx.WithValue(ctx.Background(), "testContextID", fmt.Sprintf("complex_ctx1_scenario%d", scenarioID))
// 			instrumentedCtx1 := container.RegisterContext(baseCtx1)

// 			// Resolve singleton first to establish its instance
// 			singletonX, err := GoFac.Resolve[IServiceX](instrumentedCtx1, container) // Context doesn't matter for singleton resolution itself
// 			require.NoError(t, err)
// 			require.NotNil(t, singletonX)

// 			// Resolve PerContext service
// 			serviceY1, err := GoFac.Resolve[IServiceY](instrumentedCtx1, container)
// 			require.NoError(t, err)
// 			require.NotNil(t, serviceY1)
// 			assert.Contains(t, serviceY1.Val(), singletonX.ID(), "ServiceY should contain the ID of the singleton ServiceX")

// 			// Resolve PerContext service again from same context
// 			serviceY2, err := GoFac.Resolve[IServiceY](instrumentedCtx1, container)
// 			require.NoError(t, err)
// 			require.NotNil(t, serviceY2)
// 			assert.Same(t, serviceY1, serviceY2, "PerContext instances from same context should be same")
// 			assert.Contains(t, serviceY2.Val(), singletonX.ID())

// 			// Resolve singleton again, should be the same original singleton
// 			singletonX2, err := GoFac.Resolve[IServiceX](ctx.Background(), container)
// 			require.NoError(t, err)
// 			assert.Same(t, singletonX, singletonX2, "Singleton instance should remain the same")
// 		},
// 	})

// 	// Example: Array dependency where elements are PerContext
// 	// For this, ss.StructRelyingOnIndependentStructs takes []ss.IIndependentStruct
// 	// We need to register multiple ss.IIndependentStruct as PerContext
// 	stressScenarios = append(stressScenarios, stressScenario{
// 		name:        "ResolveArrayDepPerContext",
// 		description: "Service with []IIndependentStruct dependency, where IIndependentStructs are PerContext.",
// 		setupAndRun: func(t *testing.T, scenarioID int) {
// 			cb := ContainerBuilder.New()
// 			require.NoError(t, cb.Register(ss.NewA, RegistrationOptions.As[ss.IIndependentStruct], RegistrationOptions.PerContext)) // Impl A
// 			require.NoError(t, cb.Register(ss.NewB, RegistrationOptions.As[ss.IIndependentStruct], RegistrationOptions.PerContext)) // Impl B
// 			require.NoError(t, cb.Register(ss.NewStructRelyingOnIndependentStructs, RegistrationOptions.As[ss.IStructRelyingOnIndependentStructs]))

// 			container, err := cb.Build(&i.RealUUIDProvider{})
// 			require.NoError(t, err)

// 			baseCtx := ctx.WithValue(ctx.Background(), "testContextID", fmt.Sprintf("arrayDepCtx_scenario%d", scenarioID))
// 			instrumentedCtx := container.RegisterContext(baseCtx)

// 			mainService, err := GoFac.Resolve[ss.IStructRelyingOnIndependentStructs](instrumentedCtx, container)
// 			require.NoError(t, err)
// 			require.NotNil(t, mainService)

// 			subStructNames := mainService.ReturnSubStructNames()
// 			require.Len(t, subStructNames, 2) // Expecting NewA and NewB per context

// 			// To verify that the dependencies were indeed per-context and from *this* context,
// 			// this test would be more robust if we could inspect the actual instances of dependencies.
// 			// For now, we trust ResolveMultiple used by []IIndependentStruct works with PerContext correctly.
// 			// A deeper check would involve resolving []IIndependentStruct from instrumentedCtx and comparing.
// 			resolvedDeps, err := GoFac.ResolveMultiple[ss.IIndependentStruct](instrumentedCtx, container)
// 			require.NoError(t, err)
// 			require.Len(t, resolvedDeps, 2)

// 			// Resolve mainService again from the same context
// 			mainService2, err := GoFac.Resolve[ss.IStructRelyingOnIndependentStructs](instrumentedCtx, container)
// 			require.NoError(t, err)
// 			require.NotNil(t, mainService2) // Ensure mainService2 is not nil

// 			// Since StructRelyingOnIndependentStructs is likely PerCall here,
// 			// mainService and mainService2 are different instances.
// 			// However, their resolved PerContext dependencies should be the same set of instances for 'instrumentedCtx'.
// 			// We can check if the names of the resolved dependencies are consistent.
// 			subStructNames2 := mainService2.ReturnSubStructNames()
// 			require.Len(t, subStructNames2, 2, "mainService2 should also have 2 sub-structs")
// 			assert.ElementsMatch(t, subStructNames, subStructNames2, "Sub-struct names should be consistent for instances resolved with the same context, even if parent is PerCall")

// 			// This check is tricky: StructRelyingOnIndependentStructs creates a new slice of dependencies
// 			// even if the underlying dependencies are the same. So mainService != mainService2 (if per-call for mainService)
// 			// The internal dependencies (elements of the slice) should be the same PerContext instances.
// 			// For ss.StructRelyingOnIndependentStructs, since it's not PerContext itself, a new instance is made.
// 			// To test this properly, ss.StructRelyingOnIndependentStructs itself should be PerContext or Singleton.
// 			// Let's assume ss.StructRelyingOnIndependentStructs is PerCall for this specific scenario.
// 			// The key is that its dependencies ([]ss.IIndependentStruct) are resolved using the provided context.
// 		},
// 	})
// }

// // TestMainStress is the entry point for the stress test.
// func TestMainStress_ResolveConcurrency(t *testing.T) {
// 	// It's good practice to seed the global random source once, or use rand.New with a source for local randomness.
// 	// Using math/rand/v2.New(math/rand/v2.NewPCG(...)) is the modern way for non-global rand.
// 	// For simplicity with `math/rand` used in scenarios:
// 	rand.Seed(time.Now().UnixNano())

// 	initStressScenarios(t) // Initialize basic scenarios
// 	addComplexScenarios(t) // Add more complex ones

// 	if len(stressScenarios) == 0 {
// 		t.Fatal("No stress scenarios defined.")
// 	}

// 	var wg sync.WaitGroup
// 	wg.Add(num_GOROUTINES)

// 	for i := 0; i < num_GOROUTINES; i++ {
// 		go func(goroutineID int) {
// 			defer wg.Done()

// 			// Randomly select a scenario
// 			scenarioIndex := rand.Intn(len(stressScenarios))
// 			selectedScenario := stressScenarios[scenarioIndex]

// 			// Create a subtest for better logging and isolation of failures per iteration/scenario
// 			// This can make output very verbose with 5000 iterations.
// 			// For a less verbose output initially, you might run the scenario directly.
// 			// t.Run(fmt.Sprintf("Goroutine%d_Scenario%s", goroutineID, selectedScenario.name), func(subTest *testing.T) {
// 			//    selectedScenario.setupAndRun(subTest, goroutineID)
// 			// })
// 			// Running directly for now to reduce verbosity during stress. Errors will still point to line numbers.
// 			selectedScenario.setupAndRun(t, goroutineID)

// 		}(i)
// 	}

// 	wg.Wait()
// 	t.Logf("Completed %d concurrent resolutions with random scenarios.", num_GOROUTINES)
// }

// // --- Helper for unique strings, if needed by sample constructors ---
// var (
// 	strCounter int32
// 	strMutex   sync.Mutex
// )

// func NextStr() string {
// 	strMutex.Lock()
// 	defer strMutex.Unlock()
// 	strCounter++
// 	return fmt.Sprintf("%d", strCounter)
// }
