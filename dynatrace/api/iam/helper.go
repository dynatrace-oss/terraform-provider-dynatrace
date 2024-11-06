package iam

import (
	"context"
	"encoding/json"
)

func GET(client IAMClient, ctx context.Context, url string, expectedResponseCode int, forceNewBearer bool, target any) error {
	data, err := client.GET(ctx, url, expectedResponseCode, forceNewBearer)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &target)
}
