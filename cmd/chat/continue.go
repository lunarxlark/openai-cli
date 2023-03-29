package chat

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/lunarxlark/oai-go/cmd/model"
	"github.com/urfave/cli/v2"
)

func CmdContinue(ctx *cli.Context) error {
	model, err := model.List()
	if err != nil {
		return err
	}

	log, err := os.Create(filepath.Join(os.Getenv("XDG_CONFIG_HOME"), "oai", "chat", fmt.Sprintf("%d.json", time.Now().Unix())))
	if err != nil {
		return err
	}
	defer log.Close()

	conversation := new(Conversation)

	sc := bufio.NewScanner(os.Stdin)
	for {
		fmt.Printf("%s > ", user)
		if !sc.Scan() {
			conversation.Summary, err = CreateSummay(model, conversation.Messages)
			b, err := json.MarshalIndent(conversation, "", "  ")
			if err != nil {
				return err
			}
			if _, err := log.Write(b); err != nil {
				return err
			}
		}
		statement := []Message{{
			Role:    user,
			Content: sc.Text(),
		}}
		conversation.Messages = append(conversation.Messages, statement...)

		res, err := CreateReq(model, statement).Request()
		if err != nil {
			return err
		}

		for _, choice := range res.Choices {
			fmt.Println(choice.Message.String())
			conversation.Messages = append(conversation.Messages, choice.Message)
		}
	}
}
