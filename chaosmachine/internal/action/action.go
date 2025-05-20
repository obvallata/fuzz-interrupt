package action

import (
	"context"
	"fmt"
	"log"
	"maps"
	"strings"

	"diploma/chaosmachine/internal/breakpoint"
	"diploma/chaosmachine/internal/config"
	"diploma/chaosmachine/internal/interaction"
	"diploma/keypoint/injection"
	"diploma/keypoint/schema"
	"diploma/keypoint/utils/ptr"
	"github.com/looplab/fsm"
)

type Action interface {
	BuildAutomaton(ctx context.Context) error
	HandleNotification(request schema.NotifyRequest)
	HandleBreakpoint(injectionName string)
}

type action struct {
	clients       interaction.Clients
	notifications chan schema.NotifyRequest
	config        config.Config

	state state
}

type state struct {
	data               map[string]any
	scenario           scenario
	breakpointRollback []breakpoint.Rollback
}

func NewAction(clients interaction.Clients, config config.Config) Action {
	return &action{
		clients:       clients,
		notifications: make(chan schema.NotifyRequest),
		config:        config,
	}
}

func (a *action) HandleNotification(request schema.NotifyRequest) {
	a.notifications <- request
}

func (a *action) HandleBreakpoint(injectionName string) {
	rollback, err := breakpoint.HandleInjection(
		a.state.scenario[injectionName].BreakpointInjectionConfig,
		a.state.data,
	)
	if err != nil {
		log.Printf("breakpoint handling: %v", err)
	}

	if rollback != nil {
		a.state.breakpointRollback = append(a.state.breakpointRollback, rollback)
	}
}

func (a *action) BuildAutomaton(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	bpHandler := breakpointHandler{
		action: a,
		dlv:    a.clients.Dlv,
		config: a.config.Breakpoint,
	}
	if err := bpHandler.run(ctx); err != nil {
		return fmt.Errorf("run breakpoint handler: %w", err)
	}

	if err := a.initKeypoints(); err != nil {
		return fmt.Errorf("init keypoints: %w", err)
	}
	defer a.finalizeKeypoints()

	var events fsm.Events
	for s := range newNextScenario(a.config.Injections) {
		newEvents, err := a.runScenario(s)
		if err != nil {
			return err
		}

		events = append(events, newEvents...)
	}

	sumFsm := fsm.NewFSM("START", events, fsm.Callbacks{})
	fmt.Println(fsm.Visualize(sumFsm))
	return nil
}

func (a *action) runScenario(s scenario) (fsm.Events, error) {
	var events fsm.Events

	a.state = state{
		data:     make(map[string]any),
		scenario: s,
	}

	defer func() {
		for _, rollback := range a.state.breakpointRollback {
			if err := rollback(); err != nil {
				log.Printf("breakpoint rollback: %v")
			}
		}
	}()

	if err := a.setKeypoints(s); err != nil {
		return nil, fmt.Errorf("set keypoints: %w", err)
	}

	var (
		currentState = "START"
		currentEdge  []string
	)

	for notification := range a.notifications {
		switch notification.Type {
		case schema.NotifyStateType:
			maps.Copy(a.state.data, notification.Data)

			// flush current edge injections
			edge := strings.Join(currentEdge, "\n")
			events = append(events, fsm.EventDesc{
				Name: edge,
				Src:  []string{currentState},
				Dst:  notification.Name,
			})

			currentEdge = currentEdge[:0]
			currentState = notification.Name

		case schema.NotifyInjectionType:
			maps.Copy(a.state.data, notification.Data)

			if s[notification.Name].Type == injection.TypeOff {
				currentEdge = append(currentEdge, notification.Name)
			} else {
				currentEdge = append(currentEdge, fmt.Sprintf("%s + %s", notification.Name, s[notification.Name].Type))
			}
		}

		if a.config.States[currentState].Finish {
			break
		}
	}

	return events, nil
}

func (a *action) setKeypoints(s scenario) error {
	// TODO: batch
	for name, conf := range s {
		if err := a.clients.KeyPoint.EnableInjection(name, conf.Config); err != nil {
			return fmt.Errorf("enable injection: %w", err)
		}
	}
	return nil
}

func (a *action) initKeypoints() error {
	notifierURL := fmt.Sprintf("http://%s", a.config.Server.URL)

	if err := a.clients.KeyPoint.EnableMonitor(schema.EnableMonitorRequest{
		NotifierURL: ptr.T(notifierURL),
	}); err != nil {
		return err
	}

	return nil
}

func (a *action) finalizeKeypoints() error {
	return a.clients.KeyPoint.DisableMonitor()
}
