package repository

import (
	"errors"
	"go_kit_project/internal/db"
	"go_kit_project/internal/entity"
	"go_kit_project/internal/static"
	"gopkg.in/mgo.v2/bson"
	"time"
)

var (
	ErrResponseObjectExists       = errors.New(static.MsgResponseObjectExists)
	ErrResponseServerErrorNoID    = errors.New(static.MsgResponseServerErrorNoID)
	ErrResponseServerErrorWrongID = errors.New(static.MsgResponseServerErrorWrongID)
	ErrResponseServerErrorNoData  = errors.New(static.MsgResponseServerErrorNoData)
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
}

func NewPersonRepository(connMgo *db.MongoConnection) PersonRepository {
	_ = connMgo.EnsureIndex(static.CollectionPerson, []string{static.FieldCi})
	return &personRepository{
		connMgo: connMgo,
	}
}

//Creating Person
func (p *personRepository) CreatePerson(r *interface{}) (error error) {
	object := p.ToEntityObject(*r)
	object.Created, object.Updated = time.Now(), time.Now()
	err := p.connMgo.InsertData(static.CollectionPerson, object)
	if err != nil {
		return ErrResponseObjectExists
	}
	return nil
}

//Updating Person
func (p *personRepository) UpdatePerson(id string, r *interface{}) (error error) {
	objectNew := p.ToEntityObject(*r)
	objectNew.Updated = time.Now()
	if bson.IsObjectIdHex(id) {
		filter := bson.M{static.Field__id: bson.ObjectIdHex(id)}
		found, _ := p.GetFindPersons(static.CollectionPerson, filter, bson.M{}, static.FieldUpdated, static.OrderDesc)
		if len(found) > 0 {
			document, _ := p.ToDocument(objectNew)
			update := bson.M{static.MongoDB__set: *document}
			err := p.connMgo.UpdateData(static.CollectionPerson, filter, update)
			if err != nil {
				return ErrResponseObjectExists
			}
			return nil
		}
		return ErrResponseServerErrorNoID
	}
	return ErrResponseServerErrorWrongID
}

//Listing Persons
func (p *personRepository) ListPersons() (templates []interface{}, error error) {
	found, _ := p.GetFindPersons(static.CollectionPerson, bson.M{}, bson.M{}, static.FieldLastname, static.OrderAsc)
	return found, nil
}

//Getting Person
func (p *personRepository) GetPerson(id string) (template interface{}, error error) {
	collection := static.CollectionPerson
	if bson.IsObjectIdHex(id) {
		filter := bson.M{static.Field__id: bson.ObjectIdHex(id)}
		found, _ := p.GetFindPersons(collection, filter, bson.M{}, static.FieldLastname, static.OrderAsc)
		if len(found) == 0 {
			return static.ValueEmpty, ErrResponseServerErrorNoID
		}
		return found[0], nil
	}
	return static.ValueEmpty, ErrResponseServerErrorWrongID
}

//Deleting Person
func (p *personRepository) DeletePerson(id string) (error error) {
	collection := static.CollectionPerson
	if bson.IsObjectIdHex(id) {
		filter := bson.M{static.Field__id: bson.ObjectIdHex(id)}
		found, _ := p.GetFindPersons(collection, filter, bson.M{}, static.FieldUpdated, static.OrderDesc)
		if len(found) > 0 {
			_ = p.connMgo.DeleteData(collection, id)
			return nil
		}
		return ErrResponseServerErrorNoData
	}
	return ErrResponseServerErrorWrongID
}

//FUNCIONES AUXILIARES

func (p *personRepository) GetFindPersons(collection string, query bson.M, selector bson.M, fieldSort string, orderSort string) (items []interface{}, err error) {
	items, err = p.connMgo.GetFindData(collection, query, selector, fieldSort, orderSort)
	return items, err
}

func (p *personRepository) ToDocument(v interface{}) (doc *bson.M, err error) {
	data, err := bson.Marshal(v)
	if err != nil {
		return
	}
	err = bson.Unmarshal(data, &doc)
	return
}

func (p *personRepository) ToEntityObject(i interface{}) entity.Person {
	person := entity.Person{}
	bsonBytes, _ := bson.Marshal(i)
	_ = bson.Unmarshal(bsonBytes, &person)
	return person
}
