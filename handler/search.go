package handler

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"juu17GroupBot/orm"
	"juu17GroupBot/utils"
	"strconv"
	"strings"
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
			// 以空格分割，获取最后一个字符串
			lastWord := strings.Split(inlineQuery.Query, " ")[len(strings.Split(inlineQuery.Query, " "))-1]
			// 如果最后一个字符串是数字
			offset := 0
			queryText := inlineQuery.Query
			if _, err := strconv.Atoi(lastWord); err == nil {
				offset, _ = strconv.Atoi(lastWord)
				// queryText是去掉最后一个字符串的字符串
				queryText = strings.TrimSuffix(inlineQuery.Query, " "+lastWord)
			}
			db.Select("message_id, user_id, user_name, user_first_name, user_last_name, text, date").
				Where("text LIKE ?", "%"+queryText+"%").Order("date desc").Offset(offset * 50).Limit(50).Find(&orm.GroupMessage{}).Scan(&results)
		}

		articleResults := make([]interface{}, 0)
		for _, result := range results {
			line := tgbotapi.NewInlineQueryResultArticle(
				string(rune(result.MessageID)),
				utils.Int2Date(result.Date)+"  @"+result.UserName,
				"🔍搜索结果：\n"+"https://t.me/juu17_fan/"+strconv.Itoa(result.MessageID),
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
