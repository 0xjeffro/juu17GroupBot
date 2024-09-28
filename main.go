package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"io"
	"juu17GroupBot/actions"
	"juu17GroupBot/cache"
	"juu17GroupBot/conditions"
	"juu17GroupBot/handler"
	"juu17GroupBot/orm"
	"juu17GroupBot/save"
	"juu17GroupBot/utils"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var bot *tgbotapi.BotAPI

const PROJ_NAME = "juu17GroupBot"

// 该机器人只服务于以下群组
// CurrentlyChatID 为当前服务的群组ID， 用于进群验证时的解封，移出等操作
var chatIDWhiteList = []int64{ // 白名单
	-1001924194112, // 正式群
	-1001832030593, // 测试群4
	-1001611670994, // 测试群2
	-1001661504220, // 测试群5
	5563126596,     // 管理员Jeffro
}

func main() {
	utils.CheckEnv()
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
	err := orm.Init()
	if err != nil {
		panic(err)
	}

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
	router.POST("/sendNewTwitter", sendNewTwitterHandler)

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

	// defer 踢掉所有过期未完成答题的用户
	defer func() {
		// 获得所有过期用户的ID
		userIds := cache.PopAllExpiredMemberID()
		// 把这些用户踢掉
		until := time.Now().Add(time.Minute * 1).Unix()
		CurrentChatID := os.Getenv("CURRENT_CHAT_ID")
		CurrentChatID = strings.Split(CurrentChatID, "|")[0]
		// 把CurrentChatID转换成int64
		CurrentChatIDInt64, _ := strconv.ParseInt(CurrentChatID, 10, 64)

		actions.BanUsers(bot, CurrentChatIDInt64, userIds, until)
	}()

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
			"first_name", update.Message.From.FirstName,
			"last_name", update.Message.From.LastName,
			"text", update.Message.Text,
			"date", update.Message.Date,
		)
		if update.Message.IsCommand() {
			handler.CommandHandler(bot, update)
		} else if update.Message.Text != "" {
			save.SaveMsgToDB(update)
			debug := os.Getenv("DEBUG") == "true"
			if debug && update.Message.Text == "ping" {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "pong")
				msg.ReplyToMessageID = update.Message.MessageID
				_, err := bot.Send(msg)
				if err != nil {
					zap.S().Error(err)
				}
			}
		}
	} else if update.ChatMember != nil {
		if conditions.NewMemberJoined(update) { // 如果检测到有新成员进群
			// 常量，答题限时（分钟）
			const quizTimeLimit = 20
			// 如果不在白名单内，直接返回
			if !utils.InArrayInt64(update.ChatMember.Chat.ID, chatIDWhiteList) {
				zap.S().Infow("Not in white list.", "chat_id", update.ChatMember.Chat.ID)
				return
			}

			// 回复入群成员消息，完成入群验证
			msg := tgbotapi.NewMessage(update.ChatMember.Chat.ID, "")

			msg.Text = fmt.Sprintf("👋🏻 嗨, [%s %s](https://t.me/%s)! 请在%s分钟内完成 👉🏻[入群测验](https://t.me/juu17_bot/Juu17Quiz?chatID=%s)。",
				update.ChatMember.NewChatMember.User.FirstName, update.ChatMember.NewChatMember.User.LastName,
				update.ChatMember.NewChatMember.User.UserName, fmt.Sprintf("%d", quizTimeLimit),
				fmt.Sprintf("%d", update.ChatMember.Chat.ID))
			msg.ParseMode = "Markdown"

			msg.DisableWebPagePreview = true

			req, err := bot.Send(msg)
			if err != nil {
				zap.S().Error(err)
			}

			// 把新成员禁言
			actions.RestrictUser(bot, update.ChatMember.Chat.ID, update.ChatMember.NewChatMember.User.ID)

			if _, ok := cache.GetMember(update.ChatMember.NewChatMember.User.ID); !ok {
				// 如果缓存中没有该用户的信息，就把该用户的信息存入缓存
				cache.AddMember(
					update.ChatMember.NewChatMember.User.ID,
					req.MessageID,
					time.Now().Add(time.Minute*quizTimeLimit).Unix())
			} else {
				// 如果缓存中有该用户的信息，说明该用户已经进入过群组，但是没有完成入群验证
				// 这时候需要把该用户的信息更新

				// 获得用户上一次进群欢迎消息的msgID
				member, _ := cache.GetMember(update.ChatMember.NewChatMember.User.ID)
				msgID := member.MessageId
				// 删除上一次的欢迎消息
				actions.DeleteMessage(bot, update.ChatMember.Chat.ID, msgID)
				// 更新缓存中的用户信息
				cache.UpdateMember(
					update.ChatMember.NewChatMember.User.ID,
					req.MessageID,
					time.Now().Add(time.Minute*quizTimeLimit).Unix())
			}
		}
	} else if update.InlineQuery != nil {
		handler.InlineQueryHandler(bot, update)
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

	CurrentChatID := os.Getenv("CURRENT_CHAT_ID")
	CurrentChatID = strings.Split(CurrentChatID, "|")[0]
	// 把CurrentChatID转换成int64
	CurrentChatIDInt64, err := strconv.ParseInt(CurrentChatID, 10, 64)
	if err != nil {
		log.Println(err)
		return
	}
	if !req.Pass {
		// 把用户移出群组
		until := time.Now().Add(time.Hour * 6).Unix()
		actions.BanUser(bot, CurrentChatIDInt64, req.UserID, until)
		// 删除用户的欢迎消息
		member, _ := cache.GetMember(req.UserID)
		msgID := member.MessageId
		actions.DeleteMessage(bot, CurrentChatIDInt64, msgID)
	} else {
		// 把用户解除禁言
		actions.UnrestrictUser(bot, CurrentChatIDInt64, req.UserID)

		// 删除用户的欢迎消息
		member, _ := cache.GetMember(req.UserID)
		msgID := member.MessageId
		actions.DeleteMessage(bot, CurrentChatIDInt64, msgID)

		// 把用户从缓存中删除
		cache.DeleteMember(req.UserID)
	}
	return
}

func sendNewTwitterHandler(c *gin.Context) {

	type PostData struct {
		Token      string `json:"token"`
		Text       string `json:"text"`
		TwitterURL string `json:"twitter_url"`
		Pin        bool   `json:"pin"` // 是否置顶
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
		c.JSON(http.StatusOK, gin.H{
			"message": err,
		})
		return
	}

	if req.Token != os.Getenv("BOT_TOKEN") {
		c.JSON(http.StatusOK, gin.H{
			"message": "token error",
		})
		return
	} else {
		CurrentChatIDStr := os.Getenv("CURRENT_CHAT_ID")
		// 用｜分割字符串，然后遍历
		CurrentChatIDs := strings.Split(CurrentChatIDStr, "|")
		for _, CurrentChatID := range CurrentChatIDs {
			// 把CurrentChatID转换成int64
			CurrentChatIDInt64, err := strconv.ParseInt(CurrentChatID, 10, 64)
			if err != nil {
				log.Println(err)
				return
			}
			if req.Pin {
				msg := tgbotapi.NewMessage(CurrentChatIDInt64, req.Text+"\n"+req.TwitterURL)
				msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonURL("查看原文",
							req.TwitterURL),
					),
				)
				msg.DisableWebPagePreview = false
				resp, err := bot.Send(msg)
				if err != nil {
					log.Println(err)
					return
				}
				actions.PinMessage(bot, CurrentChatIDInt64, resp.MessageID, true)
			} else {
				msg := tgbotapi.NewMessage(CurrentChatIDInt64, fmt.Sprintf("%s\n%s", req.Text, req.TwitterURL))
				_, err := bot.Send(msg)
				if err != nil {
					log.Println(err)
					return
				}
			}
		}
	}
}
