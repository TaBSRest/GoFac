package interfaces

import (
	"github.com/google/uuid"
)

type UUIDProvider interface {
	New() uuid.UUID
}

type RealUUIDProvider struct{}

func (RealUUIDProvider) New() uuid.UUID { return uuid.New() }
