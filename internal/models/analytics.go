package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Analytics struct {
	UserId    primitive.ObjectID `bson:"user_id"`
	CompanyId string             `bson:"company_id"`
	Data      []byte             `bson:"data"`
}

type Company struct {
	SNo         int    `bson:"s_no" json:"s_no"`
	Company     string `bson:"company" json:"company"`
	Country     string `bson:"country" json:"country"`
	CountryCode string `bson:"country_code" json:"country_code"`
}
