package goauthmongo

import (
	goauthpkg "chotot/go_auth/pkg"
	goauthcrypto "chotot/go_auth/pkg/crypto"
	"fmt"

	"github.com/mongodb/mongo-go-driver/mongo/options"

	"github.com/mongodb/mongo-go-driver/bson"

	"github.com/mongodb/mongo-go-driver/mongo"
)

// UserCollection Holds a collection and other methods to manipulate it
type UserCollection struct {
	coll *mongo.Collection
	hash *goauthcrypto.HashStr
}

// NewUserCollection Creates new user service to holds the collection (table)
func NewUserCollection(client *Client, dbName string, collName string, hash *goauthcrypto.HashStr) *UserCollection {
	coll := client.GetCollection(dbName, collName)
	ctx, cancelFunc := CtxCreator(5)
	defer cancelFunc()
	coll.Indexes().CreateOne(ctx, userModelIndexCreator())
	return &UserCollection{coll, hash}
}

// CreateUser Insert new user to the user Collection
func (uc *UserCollection) CreateUser(u *goauthpkg.UserObj) error {
	ctx, cancelFunc := CtxCreator(5)
	defer cancelFunc()
	hash, hashErr := uc.hash.Generate(u.Password)
	if hashErr != nil {
		return hashErr
	}

	u.Password = hash
	res, err := uc.coll.InsertOne(ctx, u)
	fmt.Println("Inserted ID: ", res.InsertedID)
	return err
}

// GetByUsername get user from mongodb
func (uc *UserCollection) GetByUsername(username string) (goauthpkg.UserObj, error) {
	ctx, cancelFunc := CtxCreator(5)
	defer cancelFunc()

	var userObj goauthpkg.UserObj

	error := uc.coll.FindOne(ctx, bson.M{
		"username": username,
	}, &options.FindOneOptions{
		Projection: bson.M{"password": 0},
	}).Decode(&userObj)

	return userObj, error
}

// GetAllUser get all user object from mongodb
func (uc *UserCollection) GetAllUser() (*[]goauthpkg.UserObj, error) {
	ctx, cancelFunc := CtxCreator(5)
	defer cancelFunc()

	cursor, err := uc.coll.Find(ctx, bson.M{}, &options.FindOptions{
		Projection: bson.M{"password": 0},
	})
	if err != nil {
		return nil, err
	}

	userObjArr := []goauthpkg.UserObj{}

	for cursor.Next(ctx) {
		var userObj goauthpkg.UserObj
		err = cursor.Decode(&userObj)
		if err != nil {
			return nil, err
		}
		userObjArr = append(userObjArr, userObj)
	}

	return &userObjArr, cursor.Err()
}

// Login let user log into server
func (uc *UserCollection) Login(c goauthpkg.Credential) (goauthpkg.UserObj, error) {
	ctx, cancelFunc := CtxCreator(5)
	defer cancelFunc()
	model := userModel{}
	// Get back user from database and force them into userModel

	err := uc.coll.FindOne(ctx, bson.M{
		"username": c.Username,
	}).Decode(&model)
	if err != nil {
		return goauthpkg.UserObj{}, err
	}

	// Compares the password
	err = uc.hash.Compare(model.Password, c.Password)
	if err != nil {
		return goauthpkg.UserObj{}, err
	}

	return goauthpkg.UserObj{
		ID:       model.ID.Hex(),
		Username: model.Username,
		Password: "",
	}, err
}
