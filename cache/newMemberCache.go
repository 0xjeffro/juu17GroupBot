package cache

import (
	"sync"
	"time"
)

// NewMember 一个全局map，用来存储用户的状态
// 用户的状态是一个结构体，包括用户ID, 消息ID, 答题deadline
type NewMember struct {
	UserId    int64
	MessageId int
	Deadline  int64
}

var NewMemberCache sync.Map

// AddMember 添加用户到map
func AddMember(userId int64, messageId int, deadline int64) {
	NewMemberCache.Store(userId, NewMember{
		UserId:    userId,
		MessageId: messageId,
		Deadline:  deadline,
	})
}

// DeleteMember deleteMember 删除用户
func DeleteMember(userId int64) {
	NewMemberCache.Delete(userId)
}

// GetMember getMember 获取用户
func GetMember(userId int64) (NewMember, bool) {
	value, ok := NewMemberCache.Load(userId)
	if !ok {
		return NewMember{}, false
	}
	return value.(NewMember), true
}

func UpdateMember(userId int64, messageId int, deadline int64) {
	NewMemberCache.Store(userId, NewMember{
		UserId:    userId,
		MessageId: messageId,
		Deadline:  deadline,
	})
}

// GetAllMember getAllMember 获取所有用户
func GetAllMember() []NewMember {
	var members []NewMember
	NewMemberCache.Range(func(key, value interface{}) bool {
		members = append(members, value.(NewMember))
		return true
	})
	return members
}

// GetAllExpiredMemberID 获取所有过期用户的ID
func GetAllExpiredMemberID() []int64 {
	var userIds []int64
	NewMemberCache.Range(func(key, value interface{}) bool {
		if value.(NewMember).Deadline < time.Now().Unix() {
			userIds = append(userIds, key.(int64))
		}
		return true
	})
	return userIds
}

// PopAllExpiredMemberID 返回并删除所有过期用户的ID
func PopAllExpiredMemberID() []int64 {
	var userIds []int64
	NewMemberCache.Range(func(key, value interface{}) bool {
		if value.(NewMember).Deadline < time.Now().Unix() {
			userIds = append(userIds, key.(int64))
			NewMemberCache.Delete(key)
		}
		return true
	})
	return userIds
}
