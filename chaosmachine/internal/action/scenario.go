package action

import (
	"diploma/chaosmachine/internal/config"
	"diploma/keypoint/injection"
	"golang.org/x/exp/maps"
)

type scenario map[string]config.InjectionListElement

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

		// force append off injection
		for _, inj := range append(
			[]config.InjectionListElement{{
				Config: injection.NewDefaultOffConfig(),
			}},
			injections[injName].InjectionsList...,
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
