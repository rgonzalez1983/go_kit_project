package middleware

import (
	"context"
	"errors"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"go_kit_project/internal/entity"
	"go_kit_project/internal/static"
	"regexp"
)

var (
	ErrUnauthorizatedLengthCI       = errors.New(static.MsgUnauthorizatedLengthCI)
	ErrUnauthorizatedDigitCI        = errors.New(static.MsgUnauthorizatedDigitCI)
	ErrUnauthorizatedAlphanumericID = errors.New(static.MsgUnauthorizatedAlphanumericID)
	IsDigit                         = regexp.MustCompile(`^[0-9]+$`).MatchString
	IsAlphanumeric                  = regexp.MustCompile(`^[0-9A-Za-z]+$`).MatchString
)

type personMiddleware struct {
	logger log.Logger
}

type (
	PersonMiddleware interface {
		AuthorizationCI(method string) endpoint.Middleware
		AuthorizationID() endpoint.Middleware
	}
)

func NewPersonMiddleware(logger log.Logger) PersonMiddleware {
	return &personMiddleware{
		logger: logger,
	}
}

func (pm *personMiddleware) AuthorizationCI(method string) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, req interface{}) (response interface{}, err error) {
			traces := []interface{}{}
			err = errors.New(static.ValueEmpty)
			CI := func() string {
				if method == static.FuncCreatePerson {
					return req.(entity.PersonRequest).CI
				} else {
					objectNew := req.(entity.UpdatePersonRequest).Values.(map[string]interface{})
					if objectNew[static.FieldCi] != nil {
						return objectNew[static.FieldCi].(string)
					}
					return static.ValueEmpty
				}
			}()
			if CI != static.ValueEmpty {
				if len(CI) < 6 || len(CI) > 11 {
					traces = []interface{}{static.KeyType, static.ERROR, static.KeyMessage, ErrUnauthorizatedLengthCI.Error()}
					err = ErrUnauthorizatedLengthCI
				} else if !IsDigit(CI) {
					traces = []interface{}{static.KeyType, static.ERROR, static.KeyMessage, ErrUnauthorizatedDigitCI.Error()}
					err = ErrUnauthorizatedDigitCI
				}
				if len(traces) > 0 {
					LoggingOperation(pm.logger, traces...)
					return nil, err
				}
			}
			return next(ctx, req)
		}
	}
}

func (pm *personMiddleware) AuthorizationID() endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, req interface{}) (response interface{}, err error) {
			traces := []interface{}{}
			err = errors.New(static.ValueEmpty)
			ID := req.(entity.GetIDPersonRequest).ID
			if !IsAlphanumeric(ID) {
				traces = []interface{}{static.KeyType, static.ERROR, static.KeyMessage, ErrUnauthorizatedAlphanumericID.Error()}
				err = ErrUnauthorizatedAlphanumericID
			}
			if len(traces) > 0 {
				LoggingOperation(pm.logger, traces...)
				return nil, err
			}
			return next(ctx, req)
		}
	}
}
