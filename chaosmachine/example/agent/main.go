package main

import (
	"context"
	"fmt"
	"log"
	"math/rand/v2"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"diploma/keypoint"
)

func downloadAndSaveFile(ctx context.Context) (err error) {
	keypoint.State(ctx, "Start")
	defer func() {
		if err != nil {
			keypoint.State(ctx, "Error", keypoint.WithData(map[string]any{"error": err.Error()}))
			return
		}
		keypoint.State(ctx, "Finish")
	}()

	filePath, content, err := getFilePathWithContentWithTimeout(ctx)
	if err != nil {
		return err
	}

	keypoint.State(ctx, "FileAndContentGot", keypoint.WithData(map[string]any{"filepath": filePath}))

	if err = createFile(ctx, filePath); err != nil {
		return err
	}

	if err = writeInFile(ctx, filePath, content); err != nil {
		return err
	}

	keypoint.State(ctx, "ContentWritten")

	return nil
}

func createFile(ctx context.Context, filePath string) error {
	file, err := keypoint.WithInject(ctx, "fileCreate", os.Create)(filePath)
	if err != nil {
		return err
	}

	return file.Close()
}

func writeInFile(ctx context.Context, filePath string, content []byte) error {
	file, err := keypoint.WithInject(ctx, "fileOpen", os.OpenFile)(filePath, os.O_RDWR, 0740)
	if err != nil {
		return err
	}

	defer file.Close()

	_, err = keypoint.WithInject(ctx, "fileWrite", file.Write)(content)
	return err
}

func getFilePathWithContentWithTimeout(ctx context.Context) (string, []byte, error) {
	ctx, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
	defer cancel()

	var (
		file    string
		content []byte
		err     error
		ready   = make(chan struct{}, 1)
	)
	go func() {
		file, content, err = keypoint.WithInject(ctx, "getFileNameWithContent", getFileNameWithContent)(ctx)
		file, err = filepath.Abs(fmt.Sprintf("files/%s", file))

		ready <- struct{}{}
	}()

	select {
	case <-ready:
		return file, content, err
	case <-ctx.Done():
		return "", nil, fmt.Errorf("timeout")
	}
}

// getFileNameWithContent can be replaced by grpc/http call
func getFileNameWithContent(ctx context.Context) (string, []byte, error) {
	file, content := rand.N(100), rand.N(1_000_000_000)

	fileStr := fmt.Sprintf("%d.txt", file)
	contentBytes := []byte(fmt.Sprintf("%d", content))

	return fileStr, contentBytes, nil
}

// go build -a -x -gcflags="all=-N -l" && GOKEYPOINT_HTTP="127.0.0.1:1234" ./agent
// dlv attach --continue --headless --accept-multiclient --api-version 2 --listen 0.0.0.0:50080 <PID>
func main() {
	log.Printf("PID: %d\n", os.Getpid())

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		<-sigChan
		cancel()
	}()

	for {
		select {
		case <-ctx.Done():
			log.Println("Finished")
			return
		case <-time.Tick(3 * time.Second):
			err := downloadAndSaveFile(ctx)
			if err != nil {
				log.Println(err)
			}
		}
	}
}
