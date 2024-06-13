package scope

type LifetimeScope int

const (
	PerCall LifetimeScope = iota
	PerRequest
	PerScope
	Singleton
)

func (e LifetimeScope) String() string {
	switch e {
		case PerCall:
			return "PerCall"
		case PerRequest:
			return "PerRequest"
		case PerScope:
			return "PerScope"
		case Singleton:
			return "Singleton"
		default:
			return "Not an option"
	}
}
