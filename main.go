package main

import (
	"github.com/ijijni/wechat-gpt/bootstrap"
	"github.com/ijijni/wechat-gpt/config"
	log "github.com/sirupsen/logrus"
)

func main() {
	err := config.LoadConfig()
	if err != nil {
		log.Warn("没有找到配置文件，尝试读取环境变量")
	}
	wechatEnv := config.GetWechat()
	telegramEnv := config.GetTelegram()
	if wechatEnv != nil && *wechatEnv == "true" {
		bootstrap.StartWebChat()
	} else if telegramEnv != nil {
		bootstrap.StartTelegramBot()
	}
}
