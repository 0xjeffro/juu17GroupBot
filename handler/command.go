package handler

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"juu17GroupBot/actions"
	"log"
	"os"
	"strconv"
	"strings"
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
	case "msg":
		// 获取消息内容
		msgText := update.Message.CommandArguments()
		// 获取消息发送者的ID
		msgFrom := update.Message.From.ID
		if msgFrom == 5563126596 {
			currentChatId := os.Getenv("CURRENT_CHAT_ID")
			currentChatId = strings.Split(currentChatId, "|")[0]
			currentChatIdInt64, _ := strconv.ParseInt(currentChatId, 10, 64)
			msg = tgbotapi.NewMessage(currentChatIdInt64, msgText)
			_, err := bot.Send(msg)
			if err != nil {
				log.Println(err)
			}
		}
	case "pin":
		msgText := update.Message.CommandArguments()
		msgFrom := update.Message.From.ID
		if msgFrom == 5563126596 {
			currentChatId := os.Getenv("CURRENT_CHAT_ID")
			currentChatId = strings.Split(currentChatId, "|")[0]
			currentChatIdInt64, _ := strconv.ParseInt(currentChatId, 10, 64)
			msg = tgbotapi.NewMessage(currentChatIdInt64, msgText)
			req, err := bot.Send(msg)
			if err != nil {
				log.Println(err)
			}
			actions.PinMessage(bot, currentChatIdInt64, req.MessageID, true)
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
