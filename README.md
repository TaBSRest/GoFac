# GoFac

AutoFac (from C#) like open-source Dependency Injection Container. It contains the lifetime management features for dependencies.

## Usage

It follows a simple pattern of creating container, registering dependencies, building the container, and resolving the dependencies. Going back a step is not allowed.

1. Creating the Container

```golang
    var container GoFac.Container = GoFac.NewContainer()
```

2. Registering Dependencies

There are several methods for registering a dependency.

  * Normal Registration

```golang
    container.RegisterConstructor[T interface](
        GoFac Container;
        Constructor for the concrete implementation of the interface;
        Functions that accept *RegistrationOption, change the option, and returns error
    )
```

## Limitations

When registering context to GoFac, please do not cocurrently register context multiple times. It will result in multiple children context one per thread.

## Appendix 1. Registration Options

### Scope

There are four out-of-the-box options for controling the scope.

```golang
    import (
        o "github.com/pyj4104/GoFac/pkg/Options"
        ro "github.com/pyj4104/GoFac/internal/RegistrationOption"
    )

    option := ro.NewRegistrationOption()

    o.PerCall(option)
    o.PerRequest(option)
    o.PerScope(option)
    o.AsSingleton(option)
```

* PerCall: If registered as PerCall, the object will be created anew whenever it is needed.
* PerRequest: If registered as PerRequest, the object will be created anew per request for a dependency.
* PerScope: If registered as PerScope, the object will be created anew per context.
* AsSingleton: If registered as singleton, the object will be created at the build time.

### Registration Option Usage

```golang
    container.RegisterConstructor[T interface](
        container,
        somePackage.NewStruct(),
        PerCall, PerRequest, AsSingleton, PerScope
    )
```

The last one to be called will be used. (In this case, the dependency will be registered as PerScope)
