package models

type FormUrlQuery struct {
	RequestId string `json:"request_id"`
	RoomId    string `json:"room_id"`
	Number    string `json:"phone_number"`
	Name      string `json:"name"`
	Batch     string `json:"batch"`
}
