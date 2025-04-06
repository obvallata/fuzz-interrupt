package action

import (
	"fmt"
	"strings"

	"diploma/chaosmachine/internal/config"
	"diploma/chaosmachine/internal/interaction"
	"diploma/keypoint/injection"
	"diploma/keypoint/schema"
	"diploma/keypoint/utils/ptr"
	"github.com/looplab/fsm"
)

type Action interface {
	BuildAutomaton() error
	HandleNotification(request schema.NotifyRequest)
	HandleBreakpoint(injectionName string)
}

type action struct {
	clients       interaction.Clients
	notifications chan schema.NotifyRequest
	config        config.Config
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
	//TODO implement me
	panic("implement me")
}

func (a *action) BuildAutomaton() error {
	//bpHandler := breakpointHandler{
	//	action: a,
	//	dlv:    a.clients.Dlv,
	//	config: a.config.Breakpoint,
	//}
	//if err := bpHandler.run(); err != nil {
	//	return fmt.Errorf("run breakpoint handler: %w", err)
	//}

	if err := a.initKeypoints(); err != nil {
		return fmt.Errorf("init keypoints: %w", err)
	}

	var events fsm.Events
	for s := range newNextScenario(a.config.Injections) {
		if err := a.setKeypoints(s); err != nil {
			return fmt.Errorf("set keypoints: %w", err)
		}

		var (
			currentState = "START"
			currentEdge  []string
		)

		for notification := range a.notifications {
			switch notification.Type {
			case schema.NotifyStateType:
				edge := strings.Join(currentEdge, "\n")
				events = append(events, fsm.EventDesc{
					Name: edge,
					Src:  []string{currentState},
					Dst:  notification.Name,
				})

				currentEdge = currentEdge[:0]
				currentState = notification.Name

			case schema.NotifyInjectionType:
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
	}

	sumFsm := fsm.NewFSM("START", events, fsm.Callbacks{})
	fmt.Println(fsm.Visualize(sumFsm))
	return nil
}

func (a *action) setKeypoints(s scenario) error {
	// TODO: batch
	for name, conf := range s {
		if err := a.clients.KeyPoint.EnableInjection(name, conf); err != nil {
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
