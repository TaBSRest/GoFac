package Scope

type LifetimeScope string

const (
	PerCall    = LifetimeScope("PerCall")
	PerContext = LifetimeScope("PerContext")
	Singleton  = LifetimeScope("Singleton")
)

func (ls LifetimeScope) String() string {
	return string(ls)
}
