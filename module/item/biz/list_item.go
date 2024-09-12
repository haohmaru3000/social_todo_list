package biz

import (
	"context"
	"to_do_list/common"
	"to_do_list/module/item/model"
)

type ListItemsStorage interface {
	ListItems(
		ctx context.Context,
		filter *model.Filter,
		paging *common.Paging,
		moreKeys ...string,
	) ([]model.TodoItem, error)
}

type listItemsBiz struct {
	store ListItemsStorage
}

func NewListItemsBiz(store ListItemsStorage) *listItemsBiz {
	return &listItemsBiz{store: store}
}

func (biz *listItemsBiz) ListItems(
	ctx context.Context,
	filter *model.Filter,
	paging *common.Paging,
) ([]model.TodoItem, error) {
	data, err := biz.store.ListItems(ctx, filter, paging)

	if err != nil {
		return nil, common.ErrCannotListEntity(model.EntityName, err)
	}

	return data, nil
}
