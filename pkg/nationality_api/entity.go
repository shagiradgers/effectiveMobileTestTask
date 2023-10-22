package nationality_api

type NationalityApiResponse struct {
	Count     int                  `json:"count"`
	Name      string               `json:"name"`
	Countries []CountryApiResponse `json:"country"`
}

type CountryApiResponse struct {
	CountryId   string  `json:"country_id"`
	Probability float64 `json:"probability"`
}
