package model

import (
	"time"

	"social_todo_list/common"
)

const (
	EntityName = "UserLikeItem"
)

type Like struct {
	UserId    int                `json:"user_id" gorm:"column:user_id;"`
	ItemId    int                `json:"item_id" gorm:"column:item_id;"`
	CreatedAt *time.Time         `json:"created_at" gorm:"column:created_at;"`
	User      *common.SimpleUser `json:"-" gorm:"foreignKey:UserId;"`
}

func (l *Like) GetItemId() int { return l.ItemId }
func (l *Like) GetUserId() int { return l.UserId }

func (Like) TableName() string { return "user_like_items" }

func ErrCannotLikeItem(err error) *common.AppError {
	return common.NewCustomError(
		err,
		"Cannot like this item",
		"ErrCannotLikeItem",
	)
}

func ErrCannotUnlikeItem(err error) *common.AppError {
	return common.NewCustomError(
		err,
		"Cannot dislike this item",
		"ErrCannotUnlikeItem",
	)
}

func ErrDidNotLikeItem(err error) *common.AppError {
	return common.NewCustomError(
		err,
		"You have not liked this item",
		"ErrCannotDidNotLikeItem",
	)
}
