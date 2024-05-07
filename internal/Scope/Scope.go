package scope

type LifetimeScope int

const (
	PerCall LifetimeScope = iota
	PerRequest
	PerScope
	Singleton
)
