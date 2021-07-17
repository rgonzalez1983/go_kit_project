package endpoints

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"go_kit_project/internal/entity"
	"go_kit_project/internal/service"
)

type PersonEndpoints struct {
	CreatePersonEndpoint endpoint.Endpoint
	UpdatePersonEndpoint endpoint.Endpoint
	ListPersonsEndpoint  endpoint.Endpoint
	GetPersonEndpoint    endpoint.Endpoint
	DeletePersonEndpoint endpoint.Endpoint
}

func MakePersonEndpoints(ps service.PersonService) PersonEndpoints {
	var middlewares []endpoint.Middleware
	return PersonEndpoints{
		CreatePersonEndpoint: wrapEndpoint(makeCreatePersonEndpoint(ps), middlewares),
		UpdatePersonEndpoint: wrapEndpoint(makeUpdatePersonEndpoint(ps), middlewares),
		ListPersonsEndpoint:  wrapEndpoint(makeListPersonsEndpoint(ps), middlewares),
		GetPersonEndpoint:    wrapEndpoint(makeGetPersonEndpoint(ps), middlewares),
		DeletePersonEndpoint: wrapEndpoint(makeDeletePersonEndpoint(ps), middlewares),
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
	for _, m := range middlewares {
		e = m(e)
	}
	return e
}
