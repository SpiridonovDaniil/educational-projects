package models

type User struct {
	Name string    `json:"name"`
	Age  int       `json:"age"`
	Friends []string `json:"friends"`
}

type MakeFriends struct {
	SourceId string `json:"sourceId"`
	TargetId string `json:"targetId"`
}

type UpdateUser struct {
	NewName string    `json:"name"`
	NewAge  int       `json:"age"`
}

type Answer struct {
	Message string `json:"message,omitempty"`
	Error string `json:"error,omitempty"`
	Id string `json:"id,omitempty"`
}