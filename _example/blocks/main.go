package main

import (
	"fmt"

	"github.com/aggronmagi/promptx/v2/v2"
	"github.com/aggronmagi/promptx/v2/v2/blocks"
)

func main() {
	app := blocks.New(
		blocks.WithCommon(blocks.WithCommonOptionExec(func(ctx blocks.Context, command string) {
			fmt.Println("exec command", command)
			switch command {
			case "exit":
				ctx.Stop()
			case "help":
				fmt.Println("help")
			default:
				fmt.Println("unknown command", command)
				promptx.InputString(ctx, "input command:", "abc", func(input string) error {
					return nil
				})
			}
		})),
	)
	app.Run()
}
