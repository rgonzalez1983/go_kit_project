package repository

import (
	"errors"
	"github.com/go-kit/kit/log"
	"go_kit_project/internal/db"
	"go_kit_project/internal/entity"
	"go_kit_project/internal/middleware"
	"go_kit_project/internal/static"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type PersonRepository interface {
	CreatePerson(r *interface{}) (error error)
	UpdatePerson(id string, r *interface{}) (error error)
	ListPersons() (templates []interface{}, error error)
	GetPerson(id string) (template interface{}, error error)
	DeletePerson(id string) (error error)
}

type personRepository struct {
	connMgo *db.MongoConnection
	logger  log.Logger
}

func NewPersonRepository(connMgo *db.MongoConnection, logger log.Logger) PersonRepository {
	_ = connMgo.EnsureIndex(static.CollectionPerson, []string{static.FieldCi})
	return &personRepository{
		connMgo: connMgo,
		logger:  logger,
	}
}

//Creating Person
func (p personRepository) CreatePerson(r *interface{}) (error error) {
	object := p.ToEntityObject(*r)
	object.Created, object.Updated = time.Now(), time.Now()
	err := p.connMgo.InsertData(static.CollectionPerson, object)
	var values []interface{}
	result_err := errors.New(static.ValueEmpty)
	if err != nil {
		values = []interface{}{static.KeyType, static.ERROR, static.KeyURL, static.URLCreatingOne, static.KeyMessage, static.MsgResponseObjectExists}
		result_err = errors.New(static.MsgResponseObjectExists)
	} else {
		values = []interface{}{static.KeyType, static.SUCCESS, static.KeyURL, static.URLCreatingOne, static.KeyMessage, static.MsgResponseCreatingOne}
		result_err = nil
	}
	middleware.LoggingOperation(p.logger, values...)
	return result_err
}

//Updating Person
func (p personRepository) UpdatePerson(id string, r *interface{}) (error error) {
	objectNew := p.ToEntityObject(*r)
	objectNew.Updated = time.Now()
	filter := bson.M{static.Field__id: bson.ObjectIdHex(id)}
	found, _ := p.GetFindPersons(static.CollectionPerson, filter, bson.M{}, static.FieldUpdated, static.OrderDesc)
	var values []interface{}
	resultErr := errors.New(static.ValueEmpty)
	if len(found) > 0 {
		document, _ := p.ToDocument(objectNew)
		update := bson.M{static.MongoDB__set: *document}
		err := p.connMgo.UpdateData(static.CollectionPerson, filter, update)
		if err != nil {
			values = []interface{}{static.KeyType, static.ERROR, static.KeyURL, static.URLUpdatingOne, static.KeyMessage, static.MsgResponseObjectExists}
			resultErr = errors.New(static.MsgResponseObjectExists)
		} else {
			values = []interface{}{static.KeyType, static.SUCCESS, static.KeyURL, static.URLUpdatingOne, static.KeyMessage, static.MsgResponseUpdatingOne}
			resultErr = nil
		}
	} else {
		values = []interface{}{static.KeyType, static.ERROR, static.KeyURL, static.URLUpdatingOne, static.KeyMessage, static.MsgResponseServerErrorNoID}
		resultErr = errors.New(static.MsgResponseServerErrorNoID)
	}
	middleware.LoggingOperation(p.logger, values...)
	return resultErr
}

//Listing Persons
func (p personRepository) ListPersons() (templates []interface{}, error error) {
	found, _ := p.GetFindPersons(static.CollectionPerson, bson.M{}, bson.M{}, static.FieldLastname, static.OrderAsc)
	values := []interface{}{static.KeyType, static.SUCCESS, static.KeyURL, static.URLListingAll, static.KeyMessage, static.MsgResponseListingAll}
	middleware.LoggingOperation(p.logger, values...)
	return found, nil
}

//Getting Person
func (p personRepository) GetPerson(id string) (template interface{}, error error) {
	collection := static.CollectionPerson
	var values []interface{}
	var result interface{}
	resultErr := errors.New(static.ValueEmpty)
	if bson.IsObjectIdHex(id) {
		filter := bson.M{static.Field__id: bson.ObjectIdHex(id)}
		found, _ := p.GetFindPersons(collection, filter, bson.M{}, static.FieldLastname, static.OrderAsc)
		if len(found) == 0 {
			values = []interface{}{static.KeyType, static.ERROR, static.KeyURL, static.URLGettingOne, static.KeyMessage, static.MsgResponseServerErrorNoID}
			resultErr = errors.New(static.MsgResponseServerErrorNoID)
			result = static.ValueEmpty
		} else {
			values = []interface{}{static.KeyType, static.SUCCESS, static.KeyURL, static.URLGettingOne, static.KeyMessage, static.MsgResponseGettingOne}
			resultErr = nil
			result = found[0]
		}
	} else {
		values = []interface{}{static.KeyType, static.ERROR, static.KeyURL, static.URLGettingOne, static.KeyMessage, static.MsgResponseServerErrorWrongID}
		resultErr = errors.New(static.MsgResponseServerErrorWrongID)
		result = static.ValueEmpty
	}
	middleware.LoggingOperation(p.logger, values...)
	return result, resultErr

}

//Deleting Person
func (p personRepository) DeletePerson(id string) (error error) {
	collection := static.CollectionPerson
	var values []interface{}
	resultErr := errors.New(static.ValueEmpty)
	if bson.IsObjectIdHex(id) {
		filter := bson.M{static.Field__id: bson.ObjectIdHex(id)}
		found, _ := p.GetFindPersons(collection, filter, bson.M{}, static.FieldUpdated, static.OrderDesc)
		if len(found) > 0 {
			_ = p.connMgo.DeleteData(collection, id)
			values = []interface{}{static.KeyType, static.SUCCESS, static.KeyURL, static.URLDeletingOne, static.KeyMessage, static.MsgResponseDeletingOne}
			resultErr = nil
		} else {
			values = []interface{}{static.KeyType, static.ERROR, static.KeyURL, static.URLDeletingOne, static.KeyMessage, static.MsgResponseServerErrorNoData}
			resultErr = errors.New(static.MsgResponseServerErrorNoData)
		}
	} else {
		values = []interface{}{static.KeyType, static.ERROR, static.KeyURL, static.URLDeletingOne, static.KeyMessage, static.MsgResponseServerErrorWrongID}
		resultErr = errors.New(static.MsgResponseServerErrorWrongID)
	}
	middleware.LoggingOperation(p.logger, values...)
	return resultErr
}

//FUNCIONES AUXILIARES

func (p personRepository) GetFindPersons(collection string, query bson.M, selector bson.M, fieldSort string, orderSort string) (items []interface{}, err error) {
	items, err = p.connMgo.GetFindData(collection, query, selector, fieldSort, orderSort)
	return items, err
}

func (p personRepository) ToDocument(v interface{}) (doc *bson.M, err error) {
	data, err := bson.Marshal(v)
	if err != nil {
		return
	}
	err = bson.Unmarshal(data, &doc)
	return
}

func (p personRepository) ToEntityObject(i interface{}) entity.Person {
	person := entity.Person{}
	bsonBytes, _ := bson.Marshal(i)
	_ = bson.Unmarshal(bsonBytes, &person)
	return person
}
