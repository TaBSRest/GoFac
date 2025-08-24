# GoFac

## Introduction
GoFac is an AutoFac-like dependency injection container for Go. It provides lifetime management, named registrations, and grouping so you can describe how components are created and resolved in your application.

The typical workflow is:
1. create a container builder,
2. register constructors for your components,
3. build the container, and
4. resolve the dependencies you need.

## Registration
Use a container builder to register the constructors for the types you want to resolve later.

```go
import (
    cb "github.com/TaBSRest/GoFac/pkg/ContainerBuilder"
    ro "github.com/TaBSRest/GoFac/pkg/Options/Registration"
)

builder := cb.New()

// basic registration
builder.Register(NewFoo)

// registration with options
builder.Register(
    NewBar,
    ro.As[IBar],
    ro.Named("special"),
    ro.AsSingleton,
)
```

### Options for registration
Some useful registration options include:

* `ro.As[T]` – expose the registration as interface `T`.
* `ro.Named("name")` – give the registration a unique name.
* `ro.Grouped[T]("group")` – add the registration to a group.
* Scope options to control lifetime:
  * `ro.PerCall` – create a new instance on each resolve.
  * `ro.PerContext` – share one instance per context. Contexts must be
    registered with GoFac using `container.RegisterContext` before
    they are used for resolving so the container can track them.
  * `ro.AsSingleton` – create a single instance at build time.

Example using a scope option:

```go
builder.Register(NewCache, ro.As[ICache], ro.PerContext)
```

When using `ro.PerContext`, register each `context.Context` with the
container before resolving so GoFac can reuse instances within that
scope:

```go
ctx := container.RegisterContext(context.Background())
```

## Resolve
After registration, build the container and resolve the dependencies.

```go
container, err := builder.Build()
if err != nil {
    panic(err)
}

ctx := container.RegisterContext(context.Background())
cache, err := GoFac.Resolve[ICache](container, ctx)
if err != nil {
    panic(err)
}
```

## Installation
Add GoFac to your module with `go get`:

```bash
go get github.com/TaBSRest/GoFac
```
