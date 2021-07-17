package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"go_kit_project/internal/config"
	"go_kit_project/internal/handler"
	"go_kit_project/internal/static"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var a handler.App

func TestMain(m *testing.M) {
	if err, isConfigurable := config.ConfigEnv(); !isConfigurable {
		fmt.Printf(""+static.MsgResponseStartError+", %s", err)
	} else {
		a = App()
		code := m.Run()
		os.Exit(code)
	}
}

func App() handler.App {
	a := handler.App{}
	_ = a.Initialize(static.ValueEmpty, static.ValueEmpty)
	return a
}

func ResponseToJSON(responseBody string) (map[string]interface{}, error) {
	var JSON map[string]interface{}
	err := json.Unmarshal([]byte(responseBody), &JSON)
	return JSON, err
}

//func TestListPersons(t *testing.T) {
//	request, _ := http.NewRequest(static.HTTP_GET, static.URLListingAll, nil)
//	response := httptest.NewRecorder()
//	a.Router.ServeHTTP(response, request)
//	values := []interface{}{static.KeyType, static.TEST, static.KeyURL, static.URLListingAll, static.KeyMessage, static.MsgResponseListingAll}
//	a.LoggingOperation(values...)
//	responseBody := ResponseToJSON(response.Body.String())
//	length := len(responseBody[static.KeyResponseData].([]interface{}))
//	assert.Equal(t, response.Code, response.Code, static.MsgTestEXPECTED + " "+strconv.Itoa(response.Code))
//	assert.Equal(t, length, len(responseBody[static.KeyResponseData].([]interface{})), static.MsgTestEXPECTED + " "+strconv.Itoa(length))
//}
//
//func TestGetPerson(t *testing.T) {
//	url := static.URLGettingOne + "/60e63f2ebefb1fb4a19de900"
//	request, _ := http.NewRequest(static.HTTP_GET, url, nil)
//	response := httptest.NewRecorder()
//	a.Router.ServeHTTP(response, request)
//	message := func() string {
//		if response.Code == 500 {
//			return static.MsgResponseServerErrorNoData
//		} else {
//			return static.MsgResponseGettingOne
//		}
//	}()
//	values := []interface{}{static.KeyType, static.TEST, static.KeyURL, static.URLGettingOne, static.KeyMessage, message}
//	a.LoggingOperation(values...)
//	responseBody := ResponseToJSON(response.Body.String())
//	assert.Equal(t, response.Code, response.Code, static.MsgTestEXPECTED + " "+strconv.Itoa(response.Code))
//	assert.Equal(t, message, responseBody[static.KeyResponseMessage].(interface{}), static.MsgTestEXPECTED + " "+message)
//	assert.Equal(t, "83110715463", responseBody[static.KeyResponseData].(map[string]interface{})["ci"], "EXPECTED 83110715463")
//}

func TestCreatePerson(t *testing.T) {
	payload := []byte(`{"name" : "JUAN BRAULIO",
						"lastname" : "HERNANDEZ NAPOLES",
						"ci" : "96012326175679",
						"country" : "Cuba",
						"age" : 24,
						"gender" : "M",
						"address" : "Calle 2da, Buenos Aires, Camaguey"
						}`)
	request, _ := http.NewRequest(static.HTTP_POST, static.URLCreatingOne, bytes.NewBuffer(payload))
	response := httptest.NewRecorder()
	a.Router.ServeHTTP(response, request)
	responseBody, err := ResponseToJSON(response.Body.String())
	message, eval := func() (string, string) {
		if err != nil {
			return response.Body.String(), static.MsgResponseObjectExists
		} else {
			return responseBody[static.KeyMessage].(string), static.MsgResponseCreatingOne
		}
	}()
	assert.Equal(t, message, eval, static.MsgTestEXPECTED+" "+message)
}

//func TestUpdatePerson(t *testing.T) {
//	url := static.URLUpdatingOne + "/60e661b0befb1fb4a19df241"
//	payload := []byte(`{
//						"name" : "ANA M."
//						}`)
//	request, _ := http.NewRequest(static.HTTP_POST, url, bytes.NewBuffer(payload))
//	response := httptest.NewRecorder()
//	a.Router.ServeHTTP(response, request)
//	responseBody := ResponseToJSON(response.Body.String())
//	message := func() string {
//		if responseBody[static.KeyResponseStatusCode].(interface{}).(float64) == http.StatusConflict {
//			return static.MsgResponseObjectExists
//		} else if responseBody[static.KeyResponseStatusCode].(interface{}).(float64) == http.StatusInternalServerError {
//			return static.MsgResponseServerError
//		} else {
//			return static.MsgResponseUpdatingOne
//		}
//	}()
//	values := []interface{}{static.KeyType, static.TEST, static.KeyURL, static.URLUpdatingOne, static.KeyMessage, message}
//	a.LoggingOperation(values...)
//	assert.Equal(t, response.Code, response.Code, static.MsgTestEXPECTED + " "+strconv.FormatFloat(responseBody[static.KeyResponseStatusCode].(interface{}).(float64), 'E', -1, 64))
//	assert.Equal(t, message, responseBody[static.KeyResponseMessage].(interface{}), static.MsgTestEXPECTED + " "+message)
//	assert.Equal(t, "ANA M.", responseBody[static.KeyResponseData].(map[string]interface{})["name"], "EXPECTED ANA M.")
//}
//
//func TestDeletePerson(t *testing.T) {
//	url := static.URLDeletingOne + "/60de364abefb1fb4a19d8bb7"
//	request, _ := http.NewRequest(static.HTTP_DELETE, url, nil)
//	response := httptest.NewRecorder()
//	a.Router.ServeHTTP(response, request)
//	message := func() string {
//		if response.Code == 500 {
//			return static.MsgResponseServerErrorNoData
//		} else {
//			return static.MsgResponseUpdatingOne
//		}
//	}()
//	values := []interface{}{static.KeyType, static.TEST, static.KeyURL, static.URLDeletingOne, static.KeyMessage, message}
//	a.LoggingOperation(values...)
//	responseBody := ResponseToJSON(response.Body.String())
//	assert.Equal(t, response.Code, response.Code, static.MsgTestEXPECTED + " "+strconv.Itoa(response.Code))
//	assert.Equal(t, message, responseBody[static.KeyResponseMessage].(interface{}), static.MsgTestEXPECTED + " "+message)
//}
