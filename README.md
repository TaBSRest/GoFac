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
