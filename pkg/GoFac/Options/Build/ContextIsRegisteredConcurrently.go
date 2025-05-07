package Build

import (
	"github.com/TaBSRest/GoFac/internal/BuildOption"
)

func RegisterSameContextConcurrently(option *BuildOption.BuildOption) {
	option.IsRegisterContextRunningConcurrently = true
}
