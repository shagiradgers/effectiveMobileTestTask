package age_api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	Url = "https://api.agify.io/"
)

type AgeNotFound struct{}

func (a *AgeNotFound) Error() string {
	return "age not found"
}

func GetAge(name string) (int, error) {
	resp, err := http.Get(fmt.Sprintf("%v?name=%v", Url, name))
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	bodyByte, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	var ageResp AgeApiResponse
	err = json.Unmarshal(bodyByte, &ageResp)
	if err != nil {
		return 0, err
	}

	if ageResp.Count == 0 {
		return 0, &AgeNotFound{}
	}
	return ageResp.Age, nil
}
