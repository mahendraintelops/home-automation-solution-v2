package models

type Device struct {
	Id int64 `json:"id,omitempty"`

	Name string `json:"name,omitempty"`

	Volume int32 `json:"volume,omitempty"`
}
