package iam

import "encoding/json"

func GET(client IAMClient, url string, expectedResponseCode int, forceNewBearer bool, target any) error {
	data, err := client.GET(url, expectedResponseCode, forceNewBearer)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &target)
}
