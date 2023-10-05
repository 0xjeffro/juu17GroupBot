package actions

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func RestrictUser(bot *tgbotapi.BotAPI, chatId int64, userId int64) {
	_, err := bot.Request(tgbotapi.RestrictChatMemberConfig{
		ChatMemberConfig: tgbotapi.ChatMemberConfig{
			ChatID: chatId,
			UserID: userId,
		},
		Permissions: &tgbotapi.ChatPermissions{
			CanSendMessages:       false,
			CanSendMediaMessages:  false,
			CanSendPolls:          false,
			CanSendOtherMessages:  false,
			CanAddWebPagePreviews: false,
			CanChangeInfo:         false,
			CanInviteUsers:        false,
			CanPinMessages:        false,
		},
	})
	if err != nil {
		return
	}
}

func UnrestrictUser(bot *tgbotapi.BotAPI, chatId int64, userId int64) {
	_, err := bot.Request(tgbotapi.RestrictChatMemberConfig{
		ChatMemberConfig: tgbotapi.ChatMemberConfig{
			ChatID: chatId,
			UserID: userId,
		},
		Permissions: &tgbotapi.ChatPermissions{
			CanSendMessages:       true,
			CanSendMediaMessages:  true,
			CanSendPolls:          true,
			CanSendOtherMessages:  true,
			CanAddWebPagePreviews: true,
			CanChangeInfo:         false,
			CanInviteUsers:        true,
			CanPinMessages:        false,
		},
	})
	if err != nil {
		return
	}
}
