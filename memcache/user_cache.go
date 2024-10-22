/**** Implementation of cache.go to cache FindUser() ****/
package memcache

import (
	"context"
	"fmt"
	"sync"

	"to_do_list/module/user/model"
)

type RealStore interface {
	FindUser(ctx context.Context, conditions map[string]interface{}, moreInfo ...string) (*model.User, error)
}

type userCaching struct {
	store     Cache
	realStore RealStore
	once      *sync.Once
}

func NewUserCaching(store Cache, realStore RealStore) *userCaching {
	return &userCaching{
		store:     store,
		realStore: realStore,
		once:      new(sync.Once),
	}
}

func (uc *userCaching) FindUser(ctx context.Context, conditions map[string]interface{}, moreInfo ...string) (*model.User, error) {
	userId := conditions["id"].(int)
	key := fmt.Sprintf("user-%d", userId)

	userInCache := uc.store.Read(key)

	if userInCache != nil {
		return userInCache.(*model.User), nil
	}

	// log.Println("Query user in real store with wrk")

	// uc.once: to avoid Data-Racing
	uc.once.Do(func() {
		user, err := uc.realStore.FindUser(ctx, conditions, moreInfo...)

		if err != nil {
			panic(err)
		}
		// Update cache
		uc.store.Write(key, user)
	})

	return uc.store.Read(key).(*model.User), nil
}
