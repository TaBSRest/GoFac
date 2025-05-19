package Group

import (
	r "github.com/TaBSRest/GoFac/internal/Registration"
	gi "github.com/TaBSRest/GoFac/internal/RegistrationOption/GroupInfo"
)

type Group struct {
	*gi.GroupInfo
	Registrations []*r.Registration
}
