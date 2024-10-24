package storage

import (
	"context"

	"social_todo_list/common"
	"social_todo_list/module/userlikeitem/model"
)

func (s *sqlStore) Create(ctx context.Context, data *model.Like) error {
	if err := s.db.Create(data).Error; err != nil {
		return common.ErrDB(err)
	}

	return nil
}
