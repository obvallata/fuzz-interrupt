package injection

import (
	"time"
)

type SleepInjectionConfig struct {
	Duration time.Duration `json:"duration"`
}
