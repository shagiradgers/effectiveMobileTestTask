package nationality_api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	Url = "https://api.nationalize.io/"
)

type NationalityNotFound struct{}

func (n *NationalityNotFound) Error() string {
	return "nationality not found"
}

func GetNationality(name string) (string, error) {
	resp, err := http.Get(fmt.Sprintf("%v?name=%v", Url, name))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	bodyByte, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var nationalityResp NationalityApiResponse
	err = json.Unmarshal(bodyByte, &nationalityResp)
	if err != nil {
		return "", err
	}

	if nationalityResp.Count == 0 {
		return "", &NationalityNotFound{}
	}
	return nationalityResp.Countries[0].CountryId, nil
}
