package conditions

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func extractStatusChange(update tgbotapi.Update) (wasMember bool, isMember bool) {
	fmt.Printf("oldStatus: %s, newStatus: %s", update.ChatMember.OldChatMember.Status, update.ChatMember.NewChatMember.Status)
	oldStatus := update.ChatMember.OldChatMember.Status
	newStatus := update.ChatMember.NewChatMember.Status
	oldIsMember := update.ChatMember.OldChatMember.IsMember
	newIsMember := update.ChatMember.NewChatMember.IsMember

	//if oldStatus == newStatus {
	//	return
	//}

	fmt.Printf("oldStatus: %s, newStatus: %s", oldStatus, newStatus)
	fmt.Printf("oldIsMember: %t, newIsMember: %t", oldIsMember, newIsMember)
	wasMember = oldStatus == "member" || oldStatus == "administrator" || oldStatus == "creator" || (oldStatus == "restricted" && oldIsMember)
	isMember = newStatus == "member" || newStatus == "administrator" || newStatus == "creator" || (newStatus == "restricted" && newIsMember)

	return wasMember, isMember
}

func NewMemberJoined(update tgbotapi.Update) bool {
	wasMember, isMember := extractStatusChange(update)
	fmt.Println(wasMember, isMember)
	return !wasMember && isMember
}

func MemberLeft(update tgbotapi.Update) bool {
	wasMember, isMember := extractStatusChange(update)
	return wasMember && !isMember
}
