package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"tcp-chat/client/app"
	"tcp-chat/client/ui"
	"tcp-chat/common"
	"tcp-chat/config"
	"tcp-chat/transport"
)

func main() {
	// загружаем данные из конфига
	cfg, err := config.LoadConfig("config/config.yaml")
	if err != nil {
		fmt.Printf("failed to load config: %v", err)
	}

	// подключаемся к серверу
	conn, err := transport.DialTCP(cfg.Client.Address)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to connect: %v\n", err)
		os.Exit(1)
	}

	// принимаем имя от пользователя
	name, err := AskName()

	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to take name: %v\n", err)
		os.Exit(1)
	}

	//отправляем имя серверу
	common.SendMessage(conn, common.Message{
		FromName: name,
		To:       common.ServerId,
		Type:     common.MessageNameIssuing,
	})

	idMsg, err := common.ReceiveMessage(conn)

	if err != nil || idMsg.Type != common.MessageIdIssuing {
		fmt.Fprintf(os.Stderr, "failed to recieve id from server: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("your id is: %d\n", idMsg.To)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	errCh := make(chan error, 2)

	client := app.NewClient(idMsg.To, strings.TrimSpace(name), conn)
	go func() {
		if err := client.Start(ctx, errCh); err != nil {
			errCh <- err
		}
	}()

	cli := ui.NewCLI(client)
	go func() {
		if err := cli.Run(ctx); err != nil {
			errCh <- err
		}
	}()

	select {
	case err := <-errCh:
		fmt.Fprintf(os.Stderr, "fatal error: %v\n", err)
		cancel()
		os.Exit(1)
	case <-ctx.Done():
	}
}

// TODO: это не должно быть в контроллере
func AskName() (string, error) {
	ui.PrintBox("enter your name")
	fmt.Print("> ")
	var name string
	_, err := fmt.Scanln(&name)

	if err != nil {
		return "", fmt.Errorf("error reading name from user: %w", err)
	}

	return name, nil
}
