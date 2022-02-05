package models

import (
	"UserMicroservice/database"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

var ctx = context.TODO()
var collection = database.Connect()

type User struct {
	ID        primitive.ObjectID `bson:"_id" json:"-"`
	FirstName string             `json:"first_name" validate:"required,min=2,max=100"`
	LastName  string             `json:"last_name" validate:"min=2,max=100"`
	Password  string             `json:"-"` // need to do static validation for this.
	Email     string             `json:"email" validate:"required,email"`
	Token     *string            `json:"token"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
	UserId    string             `json:"user_id"`
}

// hash the password

func (user *User) SignUp() error {
	_, err := collection.InsertOne(ctx, user)
	if err != nil {
		return err
	}
	return nil
}

// CheckExistingUser check existing user
func CheckExistingUser(email string) (bool, error) {
	filter := bson.M{"email": email}
	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}
	if count > 0 {
		return true, nil
	}
	return false, nil
}

// login user

//
