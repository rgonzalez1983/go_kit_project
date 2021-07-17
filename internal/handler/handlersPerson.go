package handler

import (
	httptransport "github.com/go-kit/kit/transport/http"
	"go_kit_project/internal/service"
)

// CRUD

// CreatePerson godoc
// @Summary Create one person
// @Description Create one person
// @Tags CRUD
// @Accept  plain
// @Produce  json
// @Param parameters body entity.PersonRequest true "Create Person"
// @Success 200 {object} entity.InterfaceAPI
// @Router /create_person [post]
func (a *App) CreatePerson(options []httptransport.ServerOption) *httptransport.Server {
	return httptransport.NewServer(
		a.PersonEndpoints.CreatePersonEndpoint,
		service.DecodePersonRequest,
		service.RespondWithJSON,
		options...,
	)
}

// UpdatePerson godoc
// @Summary Update one person
// @Description Update of one person
// @Tags CRUD
// @Accept  plain
// @Produce  json
// @Param id path string true "ID Person"
// @Param parameters body entity.UpdatePersonRequest true "Update Person"
// @Success 200 {object} entity.InterfaceAPI
// @Router /update_person/{id} [post]
func (a *App) UpdatePerson(options []httptransport.ServerOption) *httptransport.Server {
	return httptransport.NewServer(
		a.PersonEndpoints.UpdatePersonEndpoint,
		service.DecodeUpdatePersonRequest,
		service.RespondWithJSON,
		options...,
	)
}

// ListPerson godoc
// @Summary Get details of all persons
// @Description Get details of all persons
// @Tags CRUD
// @Accept  plain
// @Produce  json
// @Success 200 {object} entity.ListPersonsResponse
// @Router /list_persons [get]
func (a *App) ListPersons(options []httptransport.ServerOption) *httptransport.Server {
	return httptransport.NewServer(
		a.PersonEndpoints.ListPersonsEndpoint,
		service.DecodeRequest,
		service.RespondWithJSON,
		options...,
	)
}

// GetPerson godoc
// @Summary Get details of one person
// @Description Get details of one person
// @Tags CRUD
// @Accept  plain
// @Produce  json
// @Param id path string true "ID Person"
// @Success 200 {object} entity.InterfaceAPI
// @Router /get_person/{id} [get]
func (a *App) GetPerson(options []httptransport.ServerOption) *httptransport.Server {
	return httptransport.NewServer(
		a.PersonEndpoints.GetPersonEndpoint,
		service.DecodeGetIDPersonRequest,
		service.RespondWithJSON,
		options...,
	)
}

// DeletePerson godoc
// @Summary Delete one person
// @Description Delete of one person
// @Tags CRUD
// @Accept  plain
// @Produce  json
// @Param id path string true "ID Person"
// @Success 200 {object} entity.InterfaceAPI
// @Router /delete_person/{id} [delete]
func (a *App) DeletePerson(options []httptransport.ServerOption) *httptransport.Server {
	return httptransport.NewServer(
		a.PersonEndpoints.DeletePersonEndpoint,
		service.DecodeGetIDPersonRequest,
		service.RespondWithJSON,
		options...,
	)
}
