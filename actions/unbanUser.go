package actions

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func UnbanUser(bot *tgbotapi.BotAPI, chatId int64, userId int64) {
	_, err := bot.Request(tgbotapi.UnbanChatMemberConfig{
		ChatMemberConfig: tgbotapi.ChatMemberConfig{
			ChatID: chatId,
			UserID: userId,
		},
	})
	if err != nil {
		return
	}
}
