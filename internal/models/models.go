package models

type Note struct {
	Data string `json:"data"`
	Once bool   `json:"once"`
	Id   string `json:"id"`
}
