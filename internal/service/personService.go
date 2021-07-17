package service

import (
	"context"
	"go_kit_project/internal/repository"
	"go_kit_project/internal/static"
	"net/http"
)

type PersonService struct {
	personRepository repository.PersonRepository
}

func NewPersonService(repo repository.PersonRepository) PersonService {
	return PersonService{
		personRepository: repo,
	}
}

func (ps PersonService) CreatePersonService(_ context.Context, data *interface{}) (string, error, int) {
	if err := ps.personRepository.CreatePerson(data); err != nil {
		return static.ERROR, err, http.StatusInternalServerError
	}
	return static.MsgResponseCreatingOne, nil, http.StatusCreated
}

func (ps PersonService) UpdatePersonService(_ context.Context, id string, data *interface{}) (string, error, int) {
	if err := ps.personRepository.UpdatePerson(id, data); err != nil {
		return static.ERROR, err, http.StatusInternalServerError
	}
	return static.MsgResponseUpdatingOne, nil, http.StatusCreated
}

func (ps PersonService) ListPersonsService(_ context.Context) (string, error, int, []interface{}) {
	list, err := ps.personRepository.ListPersons()
	if err != nil {
		return static.ERROR, err, http.StatusInternalServerError, nil
	}
	return static.MsgResponseListingAll, nil, http.StatusCreated, list
}

func (ps PersonService) GetPersonService(_ context.Context, id string) (string, error, int, interface{}) {
	object, err := ps.personRepository.GetPerson(id)
	if err != nil {
		return static.ERROR, err, http.StatusInternalServerError, nil
	}
	return static.MsgResponseGettingOne, nil, http.StatusCreated, object
}

func (ps PersonService) DeletePersonService(_ context.Context, id string) (string, error, int) {
	err := ps.personRepository.DeletePerson(id)
	if err != nil {
		return static.ERROR, err, http.StatusInternalServerError
	}
	return static.MsgResponseDeletingOne, nil, http.StatusCreated
}
