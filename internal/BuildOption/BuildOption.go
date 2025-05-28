package BuildOption

type BuildOption struct {
	IsRegisterContextRunningConcurrently bool
}

func New() *BuildOption {
	return &BuildOption{IsRegisterContextRunningConcurrently: false}
}
