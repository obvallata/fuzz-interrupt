package main

import (
	"fmt"
	"github.com/go-delve/delve/service/api"
	"github.com/go-delve/delve/service/rpc2"
	"github.com/kr/pretty"
	"log"
	"os"
)

// execute agent with go build main_test.go && ./main
// dlv attach --log --continue --headless --accept-multiclient --api-version 2 --listen 0.0.0.0:50080 <agent process pid>

// Что происходит?
// Мы ставим точку останову на строке с созданием файла, стопаем там агента, меняем права. У агента возникают проблемы.
// Проблемы продолжаются до того момента, как мы не вернем прошлые права.

// Апишка делва неинтуитивная: продолжение работы программы после остановки происходит не при client.Continue(), а каждый
// раз при чтении из канала, который вернул client.Continue().
// С другой стороны спасибо, что апишка вообще есть.
func main() {
	serverAddr := "localhost:50080"
	client := rpc2.NewClient(serverAddr)

	defer client.Disconnect(true)

	bp := &api.Breakpoint{
		File:       "/Users/ddr/fuzz-interrupt/agent/cmd/main_test.go",
		Line:       31,
		Tracepoint: true,
	}

	_, err := client.Halt()
	if err != nil {
		log.Fatalf("Error halting: %v", err)
	}

	createdBp, err := client.CreateBreakpoint(bp)
	if err != nil {
		log.Fatalf("Error creating breakpoint: %v", err)
	}

	defer client.ClearBreakpoint(createdBp.ID)

	stateChan := client.Continue()
	pretty.Println(<-stateChan)

	fmt.Printf("Created breakpoint at %s:%d\n", createdBp.FunctionName, createdBp.Line)

	err = os.Chmod("/Users/ddr/fuzz-interrupt/exp/text.txt", 0444)
	if err != nil {
		log.Fatalf("Error chmod file: %v", err)
	}

	pretty.Println(<-stateChan)

	index := 0
	for state := range stateChan {
		pretty.Println(state)
		index++
		if index > 3 {
			break
		}
	}

	defer os.Chmod("/Users/ddr/fuzz-interrupt/exp/text.txt", 0744)
}
