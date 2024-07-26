package models

type FormUrlQuery struct {
	RoomId string `json:"room_id"`
	Number string `json:"phone_number"`
	Name   string `json:"name"`
}
