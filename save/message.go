package save

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"juu17GroupBot/orm"
)

func SaveMsgToDB(update tgbotapi.Update) {
	db := orm.GetConn()
	row := orm.GroupMessage{
		MessageID:     update.Message.MessageID,
		UserID:        update.Message.From.ID,
		UserName:      update.Message.From.UserName,
		UserFirstName: update.Message.From.FirstName,
		UserLastName:  update.Message.From.LastName,
		Text:          update.Message.Text,
		Date:          update.Message.Date,
	}
	db.Create(&row)
}
