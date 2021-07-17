package service

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"go_kit_project/internal/entity"
	"go_kit_project/internal/static"
	"net/http"
)

func DecodeRequest(_ context.Context, _ *http.Request) (interface{}, error) {
	var req entity.GenericRequest
	return req, nil
}

func DecodePersonRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req entity.PersonRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	return req, nil
}

func DecodeUpdatePersonRequest(_ context.Context, r *http.Request) (interface{}, error) {
	id := mux.Vars(r)[static.KeyId]
	var req entity.UpdatePersonRequest
	var values interface{}
	if err := json.NewDecoder(r.Body).Decode(&values); err != nil {
		return nil, err
	}
	req.ID = id
	req.Values = values
	return req, nil
}

func DecodeGetIDPersonRequest(_ context.Context, r *http.Request) (interface{}, error) {
	id := mux.Vars(r)[static.KeyId]
	var req entity.GetIDPersonRequest
	req.ID = id
	return req, nil
}

func RespondWithJSON(_ context.Context, w http.ResponseWriter, payload interface{}) error {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	_, err := w.Write(response)
	return err
}
