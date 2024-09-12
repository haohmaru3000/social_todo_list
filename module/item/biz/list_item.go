package biz

import (
	"context"
	"to_do_list/common"
	"to_do_list/module/item/model"
)

type ListItemStorage interface {
	ListItems(
		ctx context.Context,
		filter *model.Filter,
		paging *common.Paging,
		moreKeys ...string,
	) ([]model.TodoItem, error)
}

type listItemBiz struct {
	store ListItemStorage
}

func NewListItemBiz(store ListItemStorage) *listItemBiz {
	return &listItemBiz{store: store}
}

func (biz *listItemBiz) ListItem(
	ctx context.Context,
	filter *model.Filter,
	paging *common.Paging,
) ([]model.TodoItem, error) {
	data, err := biz.store.ListItems(ctx, filter, paging)

	if err != nil {
		return nil, err
	}

	return data, nil
}
