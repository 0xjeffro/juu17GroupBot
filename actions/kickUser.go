package actions

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func KickUser(bot *tgbotapi.BotAPI, chatId int64, userId int64) {
	// 把用户移除群组
	_, err := bot.Request(tgbotapi.KickChatMemberConfig{
		ChatMemberConfig: tgbotapi.ChatMemberConfig{
			ChatID: chatId,
			UserID: userId,
		},
	})
	if err != nil {
		return
	}
}
