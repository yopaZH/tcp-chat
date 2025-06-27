package ui

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"tcp-chat/client/app"
)

type CLI struct {
	client app.ClientAPI
}

func NewCLI(client app.ClientAPI) *CLI {
	return &CLI{client: client}
}

func (cli *CLI) Run(ctx context.Context) error {
	go cli.printIncoming(ctx)

	reader := bufio.NewReader(os.Stdin)

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			fmt.Print("To [id]> ")
			toLine, err := reader.ReadString('\n')
			if err != nil {
				PrintBox("error parsing id. try again")
				continue
			}
			toLine = strings.TrimSpace(toLine)
			var to uint64
			n, err := fmt.Sscanf(toLine, "%d", &to)
			if n != 1 || err != nil {
				PrintBox("error parsing id. try again")
				continue
			}

			fmt.Print("Message> ")
			msg, _ := reader.ReadString('\n')
			msg = strings.TrimSpace(msg)

			cli.client.SendUserMessage(to, msg)
		}		
	}
}

func (cli *CLI) printIncoming(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-cli.client.IncomingMessages():
			if !ok {
				return
			}
			fmt.Printf("\n[%s]: %s\n", msg.FromName, msg.Content)
		}
	}
}

func PrintBox(text string) {
	textLen := len(text)
	fmt.Printf("+%s+\n", strings.Repeat("-", textLen+2))
	fmt.Println("|", text, "|")
	fmt.Printf("+%s+\n", strings.Repeat("-", textLen+2))
}
