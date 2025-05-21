package Scope

type LifetimeScope int

const (
	PerCall LifetimeScope = iota
	PerContext
	Singleton
)

func (e LifetimeScope) String() string {
	switch e {
	case PerCall:
		return "PerCall"
	case PerContext:
		return "PerContext"
	case Singleton:
		return "Singleton"
	default:
		return "Not an option"
	}
}
