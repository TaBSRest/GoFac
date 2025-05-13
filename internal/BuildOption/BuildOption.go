package BuildOption

import (
	i "github.com/TaBSRest/GoFac/interfaces"
	te "github.com/TaBSRest/GoFac/internal/TaBSError"
)

type BuildOption struct {
	IsRegisterContextRunningConcurrently bool
	UUIDProvider                         i.UUIDProvider
}

func New(uuidProvider i.UUIDProvider) (*BuildOption, error) {
	if uuidProvider == nil {
		return nil, te.New("uuidProvider is nil")
	}

	return &BuildOption{
		IsRegisterContextRunningConcurrently: false,
		UUIDProvider:                         uuidProvider,
	}, nil
}
