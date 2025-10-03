package main

import (
	"fmt"
	"log"

	"diploma/keypoint/client"
	"diploma/keypoint/injection"
	"github.com/go-delve/delve/service/api"
	"github.com/go-delve/delve/service/rpc2"
	fuzz "github.com/google/gofuzz"
	"github.com/looplab/fsm"
)

type srcDst struct {
	Src string
	Dst string
}

func main() {
	var (
		keypointClient = client.NewKeyPointClient(client.Config{URL: "http://127.0.0.1:1234"})
		fz             = fuzz.New()
	)

	dlvClient := getClient()
	defer dlvClient.Disconnect(true)

	bp, err := setBreakpoint(dlvClient)
	if err != nil {
		log.Fatalf("Error creating breakpoint: %v", err)
	}
	defer dlvClient.ClearBreakpoint(bp.ID)

	if _, err = dlvClient.ToggleBreakpoint(bp.ID); err != nil {
		log.Fatalf("Error toggle breakpoint: %v", err)
	}
	dlvClient.Continue()

	events := make(map[srcDst]string)

	for iter := 0; iter < 1000; iter++ {
		fmt.Println(iter)

		// Mock sum
		var result int
		fz.Fuzz(&result)
		if err = keypointClient.EnableInjection("sum", injection.Config{
			Type: injection.TypeMock,
			Mock: &injection.MockInjectionConfig{Outs: []any{result}},
		}); err != nil {
			log.Fatalf("Error enable keypoint: %v", err)
		}

		if err = toggleBreakpoint(dlvClient, bp); err != nil {
			log.Fatalf("Error toggle breakpoint: %v", err)
		}

		var keypoint, command string
		for keypoint != "start" || command != "notify_start" {
			keypoint, command, err = doStep(dlvClient)
			if err != nil {
				log.Fatalf("Error do step: %v", err)
			}
		}

		// Run flow
		for keypoint != "finish" || command != "notify_success" {
			nextKeypoint, nextCommand, err := doStep(dlvClient)
			if err != nil {
				log.Fatalf("Error do step: %v", err)
			}

			if nextKeypoint != keypoint {
				events[srcDst{
					Src: keypoint,
					Dst: nextKeypoint,
				}] = fmt.Sprintf("%d", result)
			}

			keypoint = nextKeypoint
			command = nextCommand
		}

		if _, err = dlvClient.ToggleBreakpoint(bp.ID); err != nil {
			log.Fatalf("Error toggle breakpoint: %v", err)
		}
		dlvClient.Continue()
	}

	buildFsm(events)
}

func toggleBreakpoint(client *rpc2.RPCClient, bp *api.Breakpoint) error {
	// Stop program
	_, err := client.Halt()
	if err != nil {
		log.Fatalf("Error halting: %v", err)
	}

	_, err = client.ToggleBreakpoint(bp.ID)
	return err
}

func buildFsm(events map[srcDst]string) {
	var fsmEvents fsm.Events
	for k, v := range events {
		fsmEvents = append(fsmEvents, fsm.EventDesc{
			Name: v,
			Src:  []string{k.Src},
			Dst:  k.Dst,
		})
	}

	sumFsm := fsm.NewFSM("start", fsmEvents, fsm.Callbacks{})
	fmt.Println(fsm.Visualize(sumFsm))
}

func getClient() *rpc2.RPCClient {
	serverAddr := "localhost:50080"
	client := rpc2.NewClient(serverAddr)
	return client
}

func setBreakpoint(client *rpc2.RPCClient) (*api.Breakpoint, error) {
	// Stop program
	_, err := client.Halt()
	if err != nil {
		log.Fatalf("Error halting: %v", err)
	}

	// Set breakpoint
	bp := &api.Breakpoint{
		File:       "/Users/ddr/fuzz-interrupt/keypoint/injection/breakpoint.go",
		Line:       3,
		Tracepoint: true,
	}

	return client.CreateBreakpoint(bp)
}

// continueToBreakpoint is used because client.Continue has weird api
//
// Обычное апи подразумевет, что нужно вызвать Continue, которое вернет канал.
// Каждый раз при чтении значения из канала delve будет делать фактический continue.
// Прочитанное значение будет содержать состояние программы на момент ПРЕДЫДУЩЕЙ остановки.
// Нам не нравится, хотим текущее состояние (хотя бы для того, чтобы видеть актуальне значения перемнных).
func continueToBreakpoint(client *rpc2.RPCClient) error {
	out := new(rpc2.CommandOut)
	return client.CallAPI("Command", &api.DebuggerCommand{Name: api.Continue, ReturnInfoLoadConfig: &normalLoadConfig}, &out)
}

func doStep(client *rpc2.RPCClient) (string, string, error) {
	if err := continueToBreakpoint(client); err != nil {
		log.Fatalf("Error eval variable: %v", err)
	}

	return getKeypointAndCommand(client)
}

func getKeypointAndCommand(client *rpc2.RPCClient) (string, string, error) {
	keypointName, err := getStr(client, "keypointName")
	if err != nil {
		return "", "", err
	}

	command, err := getStr(client, "command")
	if err != nil {
		return "", "", err
	}

	return keypointName, command, nil
}

func getStr(client *rpc2.RPCClient, varName string) (string, error) {
	varCommand, err := client.EvalVariable(api.EvalScope{GoroutineID: -1}, varName, normalLoadConfig)
	if err != nil {
		return "", err
	}

	return varCommand.Value, nil
}

var normalLoadConfig = api.LoadConfig{
	FollowPointers:     true,
	MaxVariableRecurse: 1,
	MaxStringLen:       64,
	MaxArrayValues:     64,
	MaxStructFields:    -1,
}
