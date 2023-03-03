package telegram

import (
	"github.com/ijijni/wechat-gpt/openai"
	log "github.com/sirupsen/logrus"
	"strings"
)

func Handle(msg string) *string {
	requestText := strings.TrimSpace(msg)
	reply, err := openai.Completions(requestText)
	if err != nil {
		log.Println(err)
	}
	return reply
}
