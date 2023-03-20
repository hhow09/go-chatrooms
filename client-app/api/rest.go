package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
)

func GetRoomList() ([]string, error) {
	u := url.URL{Scheme: "http", Host: fmt.Sprintf("localhost:%s", os.Getenv("WEB_HOST")), Path: "/rooms"}
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("Non-OK HTTP status: %v", res.StatusCode)
	}
	defer res.Body.Close()
	var rommlist []string
	err = json.NewDecoder(res.Body).Decode(&rommlist)
	if err != nil {
		return nil, err
	}
	return rommlist, nil
}
