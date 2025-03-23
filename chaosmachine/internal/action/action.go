package action

import (
	"diploma/chaosmachine/internal/interaction"
	"diploma/keypoint/schema"
)

type Action interface {
	BuildAutomaton()
	HandleNotification(request schema.NotifyRequest)
	HandleBreakpoint(injectionName string)
}

type action struct {
	clients interaction.Clients
}

func NewAction(clients interaction.Clients) Action {
	return &action{
		clients: clients,
	}
}

func (a *action) BuildAutomaton() {
	//TODO implement me
	panic("implement me")
}

func (a *action) HandleNotification(request schema.NotifyRequest) {
	//TODO implement me
	panic("implement me")
}

func (a *action) HandleBreakpoint(injectionName string) {
	//TODO implement me
	panic("implement me")
}
