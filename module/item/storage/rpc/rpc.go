package rpc

import (
	"context"

	"social_todo_list/demogrpc/demo"
)

type rpcClient struct {
	client demo.ItemLikeServiceClient
}

func NewClient(client demo.ItemLikeServiceClient) *rpcClient {
	return &rpcClient{client: client}
}

func (c *rpcClient) GetItemLikes(ctx context.Context, ids []int) (map[int]int, error) {
	// log.Println("gRPC client calling")
	reqIds := make([]int32, len(ids))

	for i := range ids {
		reqIds[i] = int32(ids[i])
	}

	resp, err := c.client.GetItemLikes(ctx, &demo.GetItemLikesReq{Ids: reqIds})
	if err != nil {
		return nil, err
	}

	rs := make(map[int]int)
	for k, v := range resp.Result {
		rs[int(k)] = int(v)
	}

	return rs, nil
}
