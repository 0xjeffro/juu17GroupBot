package actions

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func BanUser(bot *tgbotapi.BotAPI, chatId int64, userId int64, untilDate int64) {
	_, err := bot.Request(tgbotapi.BanChatMemberConfig{
		ChatMemberConfig: tgbotapi.ChatMemberConfig{
			ChatID: chatId,
			UserID: userId,
		},
		UntilDate: untilDate,
	})
	if err != nil {
		return
	}
}

// BanUsers 批量封禁用户
func BanUsers(bot *tgbotapi.BotAPI, chatId int64, userIds []int64, untilDate int64) {
	for _, userId := range userIds {
		BanUser(bot, chatId, userId, untilDate)
	}
}
