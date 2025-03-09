package main

import (
	"fmt"
	"github.com/go-delve/delve/service/api"
	"github.com/go-delve/delve/service/rpc2"
	"log"
)

func main() {
	client := getClient()
	defer client.Disconnect(true)

	bp, err := setBreakpoint(client)
	if err != nil {
		log.Fatalf("Error creating breakpoint: %v", err)
	}
	defer client.ClearBreakpoint(bp.ID)

	if err != continueToBreakpoint(client) {
		log.Fatalf("Error eval variable: %v", err)
	}

	for i := 0; i < 8; i++ {
		command, err := getStr(client, "command")
		if err != nil {
			log.Fatalf("Error get string: %v", err)
		}

		keypointName, err := getStr(client, "keypointName")
		if err != nil {
			log.Fatalf("Error get string: %v", err)
		}

		fmt.Printf("[%s] %s\n", keypointName, command)

		if err != continueToBreakpoint(client) {
			log.Fatalf("Error eval variable: %v", err)
		}
	}
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
