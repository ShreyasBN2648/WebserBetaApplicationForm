package models

type UserInfo struct {
	Name     string `json:"name" bson:"name"`
	Email    string `json:"email" bson:"email"`
	Age      string `json:"age-range" bson:"age-range"`
	Location string `json:"location" bson:"location"`
	Gpu      string `json:"gpu" bson:"gpu"`
	Cpu      string `json:"cpu" bson:"cpu"`
}
