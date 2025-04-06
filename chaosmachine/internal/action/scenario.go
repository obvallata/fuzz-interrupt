package action

import (
	"diploma/chaosmachine/internal/config"
	"diploma/keypoint/injection"
	"golang.org/x/exp/maps"
)

type scenario map[string]injection.Config

func newNextScenario(injections config.InjectionsConfig) chan scenario {
	// TODO: rewrite on iterators

	var (
		injectionNames = maps.Keys(injections)
		scenarios      = make(chan scenario)
		enumerate      func(i int, s scenario)
	)

	enumerate = func(i int, s scenario) {
		sCopy := make(scenario)
		maps.Copy(sCopy, s)

		if i == len(injectionNames) {
			scenarios <- sCopy
			return
		}

		injName := injectionNames[i]
		for _, inj := range append(
			[]injection.Config{injection.NewDefaultOffConfig()},
			injections[injName].Injections...,
		) {
			sCopy[injName] = inj
			enumerate(i+1, sCopy)
		}

		if i == 0 {
			close(scenarios)
		}
	}

	go enumerate(0, make(scenario))

	return scenarios
}
