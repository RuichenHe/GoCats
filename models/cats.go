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

func CheckField(orgCatInfo, newCatInfo *Cat) {
	//Check if the field is provided or not, if not, filled with original value in the database
	if (*newCatInfo).Name == "" {
		(*newCatInfo).Name = (*orgCatInfo).Name
	}
	if (*newCatInfo).Brand == "" {
		(*newCatInfo).Brand = (*orgCatInfo).Brand
	}
	if (*newCatInfo).Age == 0 {
		(*newCatInfo).Age = (*orgCatInfo).Age
	}
	if (*newCatInfo).Gender == "" {
		(*newCatInfo).Gender = (*orgCatInfo).Gender
	}
	if (*newCatInfo).Color == "" {
		(*newCatInfo).Color = (*orgCatInfo).Color
	}
	if (*newCatInfo).Weight == "" {
		(*newCatInfo).Weight = (*orgCatInfo).Weight
	}
}
