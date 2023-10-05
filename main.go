package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"io"
	"juu17GroupBot/actions"
	"juu17GroupBot/conditions"
	"juu17GroupBot/handler"
	"juu17GroupBot/utils"
	"log"
	"net/http"
	"os"
	"strings"
)

var bot *tgbotapi.BotAPI

// 该机器人只服务于以下群组
// CurrentlyChatID 为当前服务的群组ID， 用于进群验证时的解封，移出等操作
// CurrentlyChatID 默认为 chatIDWhiteList[0]，所以在测试时，只需要把测试群的ID放在 chatIDWhiteList[0] 即可
var chatIDWhiteList = []int64{ // 白名单
	-1001924194112, // 正式群
	-1001832030593, // 测试群4
	-1001611670994, // 测试群2

}
var CurrentlyChatID int64 = chatIDWhiteList[0]

func main() {
	token := func() string {
		if os.Getenv("BOT_TOKEN") == "" {
			panic("BOT_TOKEN is not set")
		} else {
			return os.Getenv("BOT_TOKEN")
		}
	}()
	webhook := func() string {
		if os.Getenv("WEBHOOK") == "" {
			panic("WEBHOOK is not set")
		}
		return strings.TrimSuffix(os.Getenv("WEBHOOK"), "/")
	}()
	port := func() string {
		if os.Getenv("PORT") == "" {
			return "8080"
		}
		return os.Getenv("PORT")
	}()
	debug := os.Getenv("DEBUG") == "true"
	utils.InitLogger()

	webhookSuffix := utils.MD5(token)
	bot = utils.InitBot(token, webhook+"/"+webhookSuffix, debug)
	startGin(webhookSuffix, port, debug)
}

func startGin(webhookSuffix string, port string, debug bool) {

	router := gin.New()
	router.Use(utils.Cors())
	if debug {
		router.Use(gin.Logger())
	}
	router.POST("/"+webhookSuffix, webhookHandler)
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "success",
		})
	})
	router.POST("/testResult", testResultHandler)

	err := router.Run(":" + port)
	if err != nil {
		log.Println(err)
	}
}

func webhookHandler(c *gin.Context) {
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println(err)
		}
	}(c.Request.Body)

	bytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Println(err)
		return
	}

	var update tgbotapi.Update
	err = json.Unmarshal(bytes, &update)
	if err != nil {
		log.Println(err)
		return
	}

	if update.Message != nil {
		if !utils.InArrayInt64(update.Message.Chat.ID, chatIDWhiteList) {
			zap.S().Infow("Not in white list.", "chat_id", update.Message.Chat.ID)
			return
		}
		zap.S().Infow("Received a message.",
			"chat_id", update.Message.Chat.ID,
			"message_id", update.Message.MessageID,
			"from", update.Message.From,
			"text", update.Message.Text)
		if update.Message.IsCommand() {
			handler.CommandHandler(bot, update)
		} else if update.Message.Text != "" {
			// handler.RepeatTextHandler(bot, update)
		}
	} else if update.ChatMember != nil {
		if conditions.NewMemberJoined(update) { // 如果检测到有新成员进群
			// 如果不在白名单内，直接返回
			if !utils.InArrayInt64(update.ChatMember.Chat.ID, chatIDWhiteList) {
				zap.S().Infow("Not in white list.", "chat_id", update.ChatMember.Chat.ID)
				return
			}
			// 把新成员禁言
			actions.RestrictUser(bot, update.ChatMember.Chat.ID, update.ChatMember.NewChatMember.User.ID)
			// 回复入群成员消息，完成入群验证
			msg := tgbotapi.NewMessage(update.ChatMember.Chat.ID, "")

			msg.Text = fmt.Sprintf("👋🏻 嗨, [%s %s](https://t.me/%s)! 请在30分钟内完成 👉🏻[入群测验](https://t.me/juu17_bot/Juu17Quiz?chatID=%s)。",
				update.ChatMember.NewChatMember.User.FirstName, update.ChatMember.NewChatMember.User.LastName,
				update.ChatMember.NewChatMember.User.UserName, fmt.Sprintf("%d", update.ChatMember.Chat.ID))
			msg.ParseMode = "Markdown"
			msg.DisableWebPagePreview = true

			_, err = bot.Send(msg)
			if err != nil {
				zap.S().Error(err)
			}
		}
	}
}

func testResultHandler(c *gin.Context) {

	type PostData struct {
		UserID int64 `json:"user_id"`
		Pass   bool  `json:"pass"`
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println(err)
		}
	}(c.Request.Body)

	bytes, err := io.ReadAll(c.Request.Body)

	if err != nil {
		log.Println(err)
		return
	}

	var req PostData
	err = json.Unmarshal(bytes, &req)
	fmt.Println(req)
	if err != nil {
		log.Println(err)
		return
	}

	// 如果不在白名单内，直接返回
	if !utils.InArrayInt64(CurrentlyChatID, chatIDWhiteList) {
		return
	}
	if !req.Pass {
		// 把用户移出群组
		actions.BanUser(bot, CurrentlyChatID, req.UserID, 0)
		actions.UnbanUser(bot, CurrentlyChatID, req.UserID)

	} else {
		// 把用户解除禁言
		actions.UnrestrictUser(bot, CurrentlyChatID, req.UserID)
	}
	return
}
