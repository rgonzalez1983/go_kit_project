package service

import (
	"context"
	"github.com/go-kit/kit/log"
	"go_kit_project/internal/kafka"
	"go_kit_project/internal/repository"
	"go_kit_project/internal/static"
	"net/http"
)

type personService struct {
	personRepository repository.PersonRepository
	logger           log.Logger
}

type PersonService interface {
	CreatePersonService(context context.Context, data *interface{}) (string, error, int)
	UpdatePersonService(context context.Context, id string, data *interface{}) (string, error, int)
	ListPersonsService(context context.Context) (string, error, int, []interface{})
	GetPersonService(context context.Context, id string) (string, error, int, interface{})
	DeletePersonService(context context.Context, id string) (string, error, int)
}

func NewPersonService(repo repository.PersonRepository, logger log.Logger) PersonService {
	return &personService{
		personRepository: repo,
		logger:           logger,
	}
}

func (ps *personService) CreatePersonService(_ context.Context, data *interface{}) (string, error, int) {
	if err := ps.personRepository.CreatePerson(data); err != nil {
		return static.ERROR, err, http.StatusInternalServerError
	}
	kafka.SaveDataToKafka(data)
	return static.MsgResponseCreatingOne, nil, http.StatusCreated
}

func (ps *personService) UpdatePersonService(_ context.Context, id string, data *interface{}) (string, error, int) {
	if err := ps.personRepository.UpdatePerson(id, data); err != nil {
		return static.ERROR, err, http.StatusInternalServerError
	}
	return static.MsgResponseUpdatingOne, nil, http.StatusCreated
}

func (ps *personService) ListPersonsService(_ context.Context) (string, error, int, []interface{}) {
	list, err := ps.personRepository.ListPersons()
	if err != nil {
		return static.ERROR, err, http.StatusInternalServerError, nil
	}
	return static.MsgResponseListingAll, nil, http.StatusCreated, list
}

func (ps *personService) GetPersonService(_ context.Context, id string) (string, error, int, interface{}) {
	object, err := ps.personRepository.GetPerson(id)
	if err != nil {
		return static.ERROR, err, http.StatusInternalServerError, nil
	}
	return static.MsgResponseGettingOne, nil, http.StatusCreated, object
}

func (ps *personService) DeletePersonService(_ context.Context, id string) (string, error, int) {
	err := ps.personRepository.DeletePerson(id)
	if err != nil {
		return static.ERROR, err, http.StatusInternalServerError
	}
	return static.MsgResponseDeletingOne, nil, http.StatusCreated
}
