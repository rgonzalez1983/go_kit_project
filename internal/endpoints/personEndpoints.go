package endpoints

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"go_kit_project/internal/entity"
	"go_kit_project/internal/middleware"
	"go_kit_project/internal/service"
	"go_kit_project/internal/static"
)

type PersonEndpoints struct {
	CreatePersonEndpoint endpoint.Endpoint
	UpdatePersonEndpoint endpoint.Endpoint
	ListPersonsEndpoint  endpoint.Endpoint
	GetPersonEndpoint    endpoint.Endpoint
	DeletePersonEndpoint endpoint.Endpoint
}

func MakePersonEndpoints(ps *service.PersonService, pm middleware.PersonMiddleware) PersonEndpoints {
	return PersonEndpoints{
		CreatePersonEndpoint: wrapEndpoint(makeCreatePersonEndpoint(*ps), []endpoint.Middleware{pm.AuthorizationCI(static.FuncCreatePerson)}),
		UpdatePersonEndpoint: wrapEndpoint(makeUpdatePersonEndpoint(*ps), []endpoint.Middleware{pm.AuthorizationCI(static.FuncUpdatePerson)}),
		ListPersonsEndpoint:  wrapEndpoint(makeListPersonsEndpoint(*ps), nil),
		GetPersonEndpoint:    wrapEndpoint(makeGetPersonEndpoint(*ps), []endpoint.Middleware{pm.AuthorizationID()}),
		DeletePersonEndpoint: wrapEndpoint(makeDeletePersonEndpoint(*ps), []endpoint.Middleware{pm.AuthorizationID()}),
	}
}

func makeCreatePersonEndpoint(ps service.PersonService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		ok, err, statusCode := ps.CreatePersonService(ctx, &request)
		return entity.OnePersonResponse{
			Message:    ok,
			StatusCode: statusCode,
		}, err
	}
}

func makeUpdatePersonEndpoint(ps service.PersonService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(entity.UpdatePersonRequest)
		ok, err, statusCode := ps.UpdatePersonService(ctx, req.ID, &req.Values)
		return entity.OnePersonResponse{
			Message:    ok,
			StatusCode: statusCode,
		}, err
	}
}

func makeListPersonsEndpoint(ps service.PersonService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		ok, err, statusCode, list := ps.ListPersonsService(ctx)
		return entity.ListPersonsResponse{
			Message:    ok,
			StatusCode: statusCode,
			Data:       list,
		}, err
	}
}

func makeGetPersonEndpoint(ps service.PersonService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(entity.GetIDPersonRequest)
		ok, err, statusCode, object := ps.GetPersonService(ctx, req.ID)
		return entity.ObjectPersonResponse{
			Message:    ok,
			StatusCode: statusCode,
			Data:       object,
		}, err
	}
}

func makeDeletePersonEndpoint(ps service.PersonService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(entity.GetIDPersonRequest)
		ok, err, statusCode := ps.DeletePersonService(ctx, req.ID)
		return entity.OnePersonResponse{
			Message:    ok,
			StatusCode: statusCode,
		}, err
	}
}

func wrapEndpoint(e endpoint.Endpoint, middlewares []endpoint.Middleware) endpoint.Endpoint {
	if middlewares != nil {
		for _, m := range middlewares {
			e = m(e)
		}
	}
	return e
}
