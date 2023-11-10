package handler

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"juu17GroupBot/orm"
	"juu17GroupBot/utils"
	"strconv"
)

func InlineQueryHandler(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	inlineQuery := update.InlineQuery
	if inlineQuery == nil {
		return
	} else {
		zap.S().Infof("InlineQueryHandler: %s", inlineQuery.Query)
		db := orm.GetConn()
		type Result struct {
			MessageID     int
			UserID        int64
			UserName      string
			UserFirstName string
			UserLastName  string
			Text          string
			Date          int
		}
		var results []Result

		if inlineQuery.Query != "" {
			db.Select("message_id, user_id, user_name, user_first_name, user_last_name, text, date").
				Where("text LIKE ?", "%"+inlineQuery.Query+"%").Order("date desc").Limit(100).Find(&orm.GroupMessage{}).Scan(&results)
		}

		articleResults := make([]interface{}, 0)
		for _, result := range results {
			line := tgbotapi.NewInlineQueryResultArticle(
				string(rune(result.MessageID)),
				utils.Int2Date(result.Date)+"  @"+result.UserName,
				"🔍搜索结果：\n"+"https://t.me/djdjsnsnssj/"+strconv.Itoa(result.MessageID),
			)
			line.Description = result.Text
			line.ThumbURL = "https://cdn-icons-png.flaticon.com/128/8377/8377294.png" //"https://cdn-icons-png.flaticon.com/128/8455/8455397.png"
			articleResults = append(articleResults, line)
		}

		_, err := bot.Request(tgbotapi.InlineConfig{
			InlineQueryID: inlineQuery.ID,
			Results:       articleResults,
		})
		if err != nil {
			zap.S().Error(err)
		}
	}
}
