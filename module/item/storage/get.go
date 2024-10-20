package storage

import (
	"context"
	"gorm.io/gorm"

	"to_do_list/common"
	"to_do_list/module/item/model"
)

func (s *sqlStore) GetItem(ctx context.Context, cond map[string]interface{}) (*model.TodoItem, error) {
	// _, span := trace.StartSpan(ctx, "item.storage.find")
	// defer span.End()
	_, span := s.tracer.Start(ctx, "item.storage.find")
	defer span.End()

	var data model.TodoItem

	if err := s.db.Where(cond).First(&data).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, common.RecordNotFound
		}

		return nil, common.ErrDB(err)
	}

	return &data, nil
}
