package entity

type PersonRequest struct {
	Name     string `json:"name,omitempty" bson:"name,omitempty"`
	LastName string `json:"lastname,omitempty" bson:"lastname,omitempty"`
	CI       string `json:"ci,omitempty" bson:"ci,omitempty"`
	Country  string `json:"country,omitempty" bson:"country,omitempty"`
	Age      int    `json:"age,omitempty" bson:"age,omitempty"`
	Gender   string `json:"gender,omitempty" bson:"gender,omitempty"`
	Address  string `json:"address,omitempty" bson:"address,omitempty"`
}

type UpdatePersonRequest struct {
	ID     string      `json:"id"`
	Values interface{} `json:"values"`
}

type GetIDPersonRequest struct {
	ID string `json:"id"`
}

type GenericRequest struct {
}
