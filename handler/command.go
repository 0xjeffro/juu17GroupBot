package handler

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"log"
)

func CommandHandler(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	cmd := update.Message.Command()
	zap.S().Info("Receive Command: \\" + cmd + ".")
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

	switch cmd {
	case "start":
		msg.Text = fmt.Sprintf("👋🏻 嗨, [%s %s](https://t.me/%s)! \n [枪宝](https://t.me/%s) 正在进化中...",
			update.Message.From.FirstName, update.Message.From.LastName, update.Message.From.UserName, bot.Self.UserName)
		msg.ParseMode = "Markdown"
		msg.DisableWebPagePreview = true
		msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonURL("➕把我加入群组➕",
					fmt.Sprintf("https://t.me/%s?startgroup=start", bot.Self.UserName)),
			),
		)

		_, err := bot.Send(msg)
		if err != nil {
			log.Println(err)
		}
		//case "ban":
		//	if update.Message.ReplyToMessage != nil {
		//		actions.BanUser(bot, update.Message.Chat.ID, update.Message.ReplyToMessage.From.ID, 0)
		//		msg.Text = fmt.Sprintf("用户 [%s %s](https://t.me/%s) 已被永久封禁.",
		//			update.Message.ReplyToMessage.From.FirstName,
		//			update.Message.ReplyToMessage.From.LastName,
		//			update.Message.ReplyToMessage.From.UserName)
		//		msg.ParseMode = "Markdown"
		//		msg.DisableWebPagePreview = true
		//		msg.ReplyToMessageID = update.Message.ReplyToMessage.MessageID
		//		_, err := bot.Send(msg)
		//		if err != nil {
		//			log.Println(err)
		//		}
		//	} else {
		//		msg.Text = "🚫请回复一条消息，然后再使用 /ban 命令。🚫"
		//		_, err := bot.Send(msg)
		//		if err != nil {
		//			log.Println(err)
		//		}
		//	}
		//case "unban":
		//	if update.Message.ReplyToMessage != nil {
		//		actions.UnbanUser(bot, update.Message.Chat.ID, update.Message.ReplyToMessage.From.ID)
		//		msg.Text = fmt.Sprintf("用户 [%s %s](https://t.me/%s) 已被解封.",
		//			update.Message.ReplyToMessage.From.FirstName,
		//			update.Message.ReplyToMessage.From.LastName,
		//			update.Message.ReplyToMessage.From.UserName)
		//		msg.ParseMode = "Markdown"
		//		msg.DisableWebPagePreview = true
		//		msg.ReplyToMessageID = update.Message.ReplyToMessage.MessageID
		//		_, err := bot.Send(msg)
		//		if err != nil {
		//			log.Println(err)
		//		}
		//	} else {
		//		msg.Text = "🚫请回复一条消息，然后再使用 /unban 命令。🚫"
		//		_, err := bot.Send(msg)
		//		if err != nil {
		//			log.Println(err)
		//		}
		//	}
		//case "kick":
		//	if update.Message.ReplyToMessage != nil {
		//		actions.KickUser(bot, update.Message.Chat.ID, update.Message.ReplyToMessage.From.ID)
		//		msg.Text = fmt.Sprintf("用户 [%s %s](https://t.me/%s) 已被踢出群组.",
		//			update.Message.ReplyToMessage.From.FirstName,
		//			update.Message.ReplyToMessage.From.LastName,
		//			update.Message.ReplyToMessage.From.UserName)
		//		msg.ParseMode = "Markdown"
		//		msg.DisableWebPagePreview = true
		//		msg.ReplyToMessageID = update.Message.ReplyToMessage.MessageID
		//		_, err := bot.Send(msg)
		//		if err != nil {
		//			log.Println(err)
		//		}
		//	} else {
		//		msg.Text = "🚫请回复一条消息，然后再使用 /kick 命令。🚫"
		//		_, err := bot.Send(msg)
		//		if err != nil {
		//			log.Println(err)
		//		}
		//	}
	}
}
