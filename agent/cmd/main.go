package main

import (
	"context"
	"encoding/base64"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc"

	distributorpb "diploma/gen/distributor"
)

func downloadAndSaveFile(client distributorpb.DistributorServiceClient) error {
	agentUUID := uuid.New()
	response, err := client.DownloadFile(context.Background(), &distributorpb.DownloadFileRequest{AgentUuid: agentUUID.String()})
	if err != nil {
		return err
	}

	content, err := base64.StdEncoding.DecodeString(response.Content)
	if err != nil {
		return err
	}

	file, err := os.Create(response.Path)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(content)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	// as mvp interface to send pid to caller in fuzzing
	savePID()
	noCron := len(os.Args) > 1 && os.Args[1] == "no-cron"

	conn, err := grpc.Dial("localhost:778", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()

	client := distributorpb.NewDistributorServiceClient(conn)

	if noCron {
		go execSignals(client)
	} else {
		go execCron(client)
	}

	select {}
}

func savePID() {
	file, err := os.Create("/Users/ddr/fuzz-interrupt/agent/src/pid")
	if err != nil {
		log.Fatalf("Failed to create pid file: %v", err)
	}
	defer file.Close()

	_, err = file.Write([]byte(strconv.Itoa(os.Getpid())))
	if err != nil {
		log.Fatalf("Failed to write pid: %v", err)
	}
}

func execCron(client distributorpb.DistributorServiceClient) {
	for {
		err := downloadAndSaveFile(client)
		if err != nil {
			log.Printf("Failed to download and save file: %v", err)
		} else {
			log.Println("File downloaded and saved successfully.")
		}

		time.Sleep(5 * time.Second)
	}
}

func execSignals(client distributorpb.DistributorServiceClient) {
	sigs := make(chan os.Signal, 1)
	// ctrl T for mac
	signal.Notify(sigs, syscall.SIGINFO)
	for {
		log.Println("Listen signals")
		sig := <-sigs
		log.Printf("Received signal: %v", sig)

		err := downloadAndSaveFile(client)
		if err != nil {
			log.Printf("Failed to download and save file: %v", err)
		} else {
			log.Println("File downloaded and saved successfully.")
		}
	}
}
