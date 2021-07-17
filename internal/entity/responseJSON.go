package entity

type ObjectPersonResponse struct {
	Message    string      `json:"message"`
	StatusCode int         `json:"statusCode"`
	Data       interface{} `json:"data"`
}

type OnePersonResponse struct {
	Message    string `json:"message"`
	StatusCode int    `json:"statusCode"`
}

type ListPersonsResponse struct {
	Message    string        `json:"message"`
	StatusCode int           `json:"statusCode"`
	Data       []interface{} `json:"data"`
}

type GenericResponse struct {
}
