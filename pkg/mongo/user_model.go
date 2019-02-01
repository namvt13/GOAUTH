package goauthmongo

import (
	goauthpkg "chotot/go_auth/pkg"

	"github.com/mongodb/mongo-go-driver/bson/primitive"

	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/options"
)

type userModel struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Username string
	Password string
}

func userModelIndexCreator() mongo.IndexModel {
	indexOptions := options.Index()
	indexOptions.SetBackground(true)
	indexOptions.SetUnique(true)
	indexOptions.SetSparse(true)

	return mongo.IndexModel{
		Keys:    []string{"username"},
		Options: indexOptions,
	}
}

func toUserModel(u *goauthpkg.UserObj) *userModel {
	return &userModel{
		Username: u.Username,
		Password: u.Password,
	}
}

func (u *userModel) toUserObj() *goauthpkg.UserObj {
	return &goauthpkg.UserObj{
		ID:       u.ID.Hex(),
		Username: u.Username,
		Password: u.Password,
	}
}
