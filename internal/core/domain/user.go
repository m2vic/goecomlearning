package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	UserId         primitive.ObjectID `json:"userid" bson:"_id"`
	Email          string             `json:"email" bson:"email"`
	FirstName      string             `json:"firstname" bson:"firstname"`
	LastName       string             `json:"lastname" bson:"lastname"`
	Username       string             `json:"username" bson:"username"`
	Password       string             `json:"password" bson:"password"`
	Token          string             `json:"token" bson:"token"`
	AddressDetails string             `json:"address" bson:"address"`
	Role           string             `json:"role" bson:"role"`
	Cart           []UsersProduct     `json:"cart" bson:"cart"`
}
