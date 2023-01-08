package models

type Cat struct {
	ID     string `json:"id,omitempty" bson:"_id,omitempty"` //json is for postman test, bson is for database communication
	Name   string `json:"name"`
	Brand  string `json:"brand"`
	Age    int    `json:"age"`
	Gender string `json:"gender"`
	Color  string `json:"color"`
	Weight string `json:"weight"`
}
