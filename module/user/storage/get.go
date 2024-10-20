package storage

import (
	"context"
	"gorm.io/gorm"

	"to_do_list/common"
	"to_do_list/module/user/model"
)

func (s *sqlStore) FindUser(ctx context.Context, conditions map[string]interface{}, moreInfo ...string) (*model.User, error) {
	// _, span := trace.StartSpan(ctx, "user.storage.find")
	// defer span.End()
	_, span := s.tracer.Start(ctx, "user.storage.find")
	defer span.End()

	db := s.db.Table(model.User{}.TableName())

	for i := range moreInfo {
		db = db.Preload(moreInfo[i])
	}

	var user model.User

	if err := db.Where(conditions).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, common.RecordNotFound
		}

		return nil, common.ErrDB(err)
	}

	return &user, nil
}
