package repository

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"social_todo_list/common"
	"social_todo_list/module/item/model"
)

type ListItemStorage interface {
	ListItem(
		ctx context.Context,
		filter *model.Filter,
		paging *common.Paging,
		moreKeys ...string,
	) ([]model.TodoItem, error)
}

type ItemLikeStorage interface {
	GetItemLikes(ctx context.Context, ids []int) (map[int]int, error)
}

type listItemRepo struct {
	store     ListItemStorage
	likeStore ItemLikeStorage
	requester common.Requester
	tracer    trace.Tracer
}

func NewListItemRepo(store ListItemStorage, likeStore ItemLikeStorage, requester common.Requester) *listItemRepo {
	return &listItemRepo{
		store:     store,
		likeStore: likeStore,
		requester: requester,
		tracer:    otel.Tracer("Item.Repo"),
	}
}

func (repo *listItemRepo) ListItem(
	ctx context.Context,
	filter *model.Filter,
	paging *common.Paging,
	moreKeys ...string,
) ([]model.TodoItem, error) {
	ctxStore := context.WithValue(ctx, common.CurrentUser, repo.requester)
	ctx, span := repo.tracer.Start(ctxStore, "item.repo.list")
	defer span.End()

	data, err := repo.store.ListItem(ctx, filter, paging, moreKeys...)

	if err != nil {
		return nil, common.ErrCannotListEntity(model.EntityName, err)
	}

	if len(data) == 0 {
		return data, nil
	}

	ids := make([]int, len(data))

	for i := range ids {
		ids[i] = data[i].Id
	}

	likeUserMap, err := repo.likeStore.GetItemLikes(ctx, ids)

	if err != nil {
		return data, nil
	}

	for i := range data {
		data[i].LikedCount = likeUserMap[data[i].Id]
	}

	return data, nil
}
