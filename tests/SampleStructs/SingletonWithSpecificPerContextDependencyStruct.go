package SampleStructs

type SingletonWithSpecificPerContextDependencyStruct struct {
	IIndependentStruct
	initialPerContextDependencyName       string
	initialPerContextDependencySpecificID string
}

func NewSingletonWithSpecificPerContextDependencyStruct(singleton IIndependentStruct, perContextDependency ISpecificPerContextStruct) ISingletonWithPerContextDependencyStruct {
	dependencyName := ""
	dependencySpecificID := ""
	if perContextDependency != nil {
		dependencyName = perContextDependency.ReturnNameIndependentStruct()
		dependencySpecificID = perContextDependency.GetSpecificID()
	}

	return &SingletonWithSpecificPerContextDependencyStruct{
		IIndependentStruct:                    singleton,
		initialPerContextDependencyName:       dependencyName,
		initialPerContextDependencySpecificID: dependencySpecificID,
	}
}

func (s *SingletonWithSpecificPerContextDependencyStruct) GetInitialPerContextDepName() string {
	return s.initialPerContextDependencyName
}

func (s *SingletonWithSpecificPerContextDependencyStruct) GetInitialPerContextDepSpecificID() string {
	return s.initialPerContextDependencySpecificID
}
