package models

type User struct {
	ID string `json:"id,omitempty" bson:"_id,omitempty"`

	Name string `json:"name,omitempty" bson:"name,omitempty"`
}
