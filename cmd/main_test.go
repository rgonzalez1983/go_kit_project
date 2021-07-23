package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"go_kit_project/internal/app"
	"go_kit_project/internal/config"
	"go_kit_project/internal/static"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var a app.App

func TestMain(m *testing.M) {
	if err, isConfigurable := config.ConfigEnv(); !isConfigurable {
		fmt.Printf(""+static.MsgResponseStartError+", %s", err)
	} else {
		a = App()
		code := m.Run()
		os.Exit(code)
	}
}

func App() app.App {
	a := app.App{}
	_ = a.Initialize(static.ValueEmpty, static.ValueEmpty)
	return a
}

func ResponseToJSON(responseBody string) (map[string]interface{}, error) {
	var JSON map[string]interface{}
	err := json.Unmarshal([]byte(responseBody), &JSON)
	return JSON, err
}

func TestListPersons(t *testing.T) {
	t.Run("PROBANDO OBTENER LISTADO DE OBJETOS", func(t *testing.T) {
		url := static.URLListingAll
		message, eval := ResponsePersonTest(http.MethodGet, url, bytes.NewBuffer(nil), static.MsgResponseListingAll, static.MsgResponseListingAll)
		assert.Equal(t, message, eval, static.MsgTestEXPECTED+" "+message)
	})
}

func TestGetPerson(t *testing.T) {
	url := static.ValueEmpty
	t.Run("PROBANDO OBTENER OBJETO", func(t *testing.T) {
		url = static.URLGettingOne + "/60e63f2ebefb1fb4a19de900"
		message, eval := ResponsePersonTest(http.MethodGet, url, bytes.NewBuffer(nil), static.MsgResponseServerErrorNoData, static.MsgResponseGettingOne)
		if message != eval {
			eval = static.MsgResponseServerErrorNoID
		}
		assert.Equal(t, message, eval, static.MsgTestEXPECTED+" "+message)
	})
	t.Run("PROBANDO MIDDLEWARE ID NO VALIDO", func(t *testing.T) {
		url = static.URLGettingOne + "/60*63f2ebefb1fb4a19de900"
		message, eval := ResponsePersonTest(http.MethodGet, url, bytes.NewBuffer(nil), static.MsgUnauthorizatedAlphanumericID, static.MsgResponseGettingOne)
		assert.Equal(t, message, eval, static.MsgTestEXPECTED+" "+message)
	})

}

func TestCreatePerson(t *testing.T) {
	t.Run("PROBANDO MIDDLEWARE LONGITUD EXCESIVA EN CREAR", func(t *testing.T) {
		payload := []byte(`{"name" : "JUAN BRAULIO",
						"lastname" : "HERNANDEZ NAPOLES",
						"ci" : "96012326175679",
						"country" : "Cuba",
						"age" : 24,
						"gender" : "M",
						"address" : "Calle 2da, Buenos Aires, Camaguey"
						}`)
		message, eval := ResponsePersonTest(http.MethodPost, static.URLCreatingOne, bytes.NewBuffer(payload), static.MsgUnauthorizatedLengthCI, static.MsgResponseCreatingOne)
		assert.Equal(t, message, eval, static.MsgTestEXPECTED+" "+message)
	})
	t.Run("PROBANDO MIDDLEWARE CI SOLO DIGITOS EN CREAR", func(t *testing.T) {
		payload := []byte(`{"name" : "JUAN BRAULIO",
						"lastname" : "HERNANDEZ NAPOLES",
						"ci" : "96P123261",
						"country" : "Cuba",
						"age" : 24,
						"gender" : "M",
						"address" : "Calle 2da, Buenos Aires, Camaguey"
						}`)
		message, eval := ResponsePersonTest(http.MethodPost, static.URLCreatingOne, bytes.NewBuffer(payload), static.MsgUnauthorizatedDigitCI, static.MsgResponseCreatingOne)
		assert.Equal(t, message, eval, static.MsgTestEXPECTED+" "+message)
	})
	t.Run("PROBANDO INSERCION U OBJETO EXISTENTE", func(t *testing.T) {
		payload := []byte(`{"name" : "AMALIA",
						"lastname" : "CORRALES BETANCOURT",
						"ci" : "96012326175",
						"country" : "Cuba",
						"age" : 24,
						"gender" : "M",
						"address" : "Calle 2da, Buenos Aires, Camaguey"
						}`)
		message, eval := ResponsePersonTest(http.MethodPost, static.URLCreatingOne, bytes.NewBuffer(payload), static.MsgResponseCreatingOne, static.MsgResponseCreatingOne)
		if message != eval {
			eval = static.MsgResponseObjectExists
		}
		assert.Equal(t, message, eval, static.MsgTestEXPECTED+" "+message)
	})
}

func TestUpdatePerson(t *testing.T) {
	url := static.ValueEmpty
	t.Run("PROBANDO ERROR OBJETO INEXISTENTE", func(t *testing.T) {
		url = static.URLUpdatingOne + "/60e661b0befb1fb4a19df241"
		payload := []byte(`{
						"name" : "ANA M."
						}`)
		message, eval := ResponsePersonTest(http.MethodPost, url, bytes.NewBuffer(payload), static.MsgResponseServerErrorNoID, static.MsgResponseUpdatingOne)
		if message != eval {
			eval = static.MsgResponseServerErrorNoID
		}
		assert.Equal(t, message, eval, static.MsgTestEXPECTED+" "+message)
	})
	t.Run("PROBANDO ACTUALIZAR OBJETO", func(t *testing.T) {
		url = static.URLUpdatingOne + "/60f74ab279f9e73baffdb9de"
		payload := []byte(`{
						"ci" : "96092017065"
						}`)
		message, eval := ResponsePersonTest(http.MethodPost, url, bytes.NewBuffer(payload), static.MsgResponseUpdatingOne, static.MsgResponseUpdatingOne)
		assert.Equal(t, message, eval, static.MsgTestEXPECTED+" "+message)
	})
	t.Run("PROBANDO MIDDLEWARE LONGITUD EXCESIVA EN ACTUALIZAR", func(t *testing.T) {
		url = static.URLUpdatingOne + "/60f74ab279f9e73baffdb9de"
		payload := []byte(`{
						"ci" : "9609201706587"
						}`)
		message, eval := ResponsePersonTest(http.MethodPost, url, bytes.NewBuffer(payload), static.MsgUnauthorizatedLengthCI, static.MsgResponseUpdatingOne)
		assert.Equal(t, message, eval, static.MsgTestEXPECTED+" "+message)
	})
	t.Run("PROBANDO MIDDLEWARE CI SOLO DIGITOS EN ACTUALIZAR", func(t *testing.T) {
		url = static.URLUpdatingOne + "/60f74ab279f9e73baffdb9de"
		payload := []byte(`{
						"ci" : "9609201T065"
						}`)
		message, eval := ResponsePersonTest(http.MethodPost, url, bytes.NewBuffer(payload), static.MsgUnauthorizatedDigitCI, static.MsgResponseUpdatingOne)
		assert.Equal(t, message, eval, static.MsgTestEXPECTED+" "+message)
	})
}

func TestDeletePerson(t *testing.T) {
	url := static.ValueEmpty
	t.Run("PROBANDO ELIMINAR OBJETO", func(t *testing.T) {
		url = static.URLDeletingOne + "/60e63f2ebefb1fb4a19de900"
		message, eval := ResponsePersonTest(http.MethodPost, url, bytes.NewBuffer(nil), static.MsgResponseServerErrorNoData, static.MsgResponseGettingOne)
		if message != eval {
			eval = static.MsgResponseServerErrorNoID
		}
		assert.Equal(t, message, eval, static.MsgTestEXPECTED+" "+message)
	})
	t.Run("PROBANDO MIDDLEWARE ID NO VALIDO", func(t *testing.T) {
		url = static.URLDeletingOne + "/60*63f2ebefb1fb4a19de900"
		message, eval := ResponsePersonTest(http.MethodPost, url, bytes.NewBuffer(nil), static.MsgUnauthorizatedAlphanumericID, static.MsgResponseGettingOne)
		assert.Equal(t, message, eval, static.MsgTestEXPECTED+" "+message)
	})
}

//FUNCIONES EXTRAS
func ResponsePersonTest(method string, url string, body *bytes.Buffer, eval string, defaultMessage string) (string, string) {
	request, _ := http.NewRequest(method, url, body)
	response := httptest.NewRecorder()
	a.Router.ServeHTTP(response, request)
	responseBody, err := ResponseToJSON(response.Body.String())
	return func() (string, string) {
		if err != nil {
			return response.Body.String(), eval
		}
		return responseBody[static.KeyMessage].(string), defaultMessage
	}()
}
