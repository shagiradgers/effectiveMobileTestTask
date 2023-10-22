package gender_api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	Url = "https://api.genderize.io/"
)

type GenderNotFound struct{}

func (g *GenderNotFound) Error() string {
	return "Gender not found"
}

func GetGender(name string) (string, error) {
	resp, err := http.Get(fmt.Sprintf("%v?name=%v", Url, name))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	bodyByte, err := io.ReadAll(resp.Body)
	if err != nil {
		return "nil", err
	}

	var genderResp GenderApiResponse
	err = json.Unmarshal(bodyByte, &genderResp)
	if err != nil {
		return "nil", err
	}

	if genderResp.Probability == 0.0 {
		return "", &GenderNotFound{}
	}

	return genderResp.Gender, nil
}
