package goauthmongo_test

import (
	goauthpkg "chotot/go_auth/pkg"
	goauthcrypto "chotot/go_auth/pkg/crypto"
	goauthmongo "chotot/go_auth/pkg/mongo"
	"log"
	"testing"

	"github.com/mongodb/mongo-go-driver/bson"
)

const (
	mongoURI     = "mongodb://localhost:27017"
	dbName       = "go_auth"
	userCollName = "user"
)

func Test_UserCollection(t *testing.T) {
	t.Run("InsertUserIntoMongoDB", createUser_test)
}

func createUser_test(t *testing.T) {
	// Arrange
	// Get new client
	client, err := goauthmongo.NewClient(mongoURI)
	if err != nil {
		log.Fatalf("Counldn't connect to mongo: %s", err)
	}
	// Disconnect when this function return
	defer func() {
		client.DropDB(dbName)
		client.Disconnect()
	}()

	// mockHash := &goauthmock.HashStr{}
	mockHash := &goauthcrypto.HashStr{}

	userCollection := goauthmongo.NewUserCollection(client, dbName, userCollName, mockHash)

	// Create some tesing username and password
	testUsername := "test_user"
	testPassword := "test_password"
	// Conver data into user object
	user := goauthpkg.UserObj{
		Username: testUsername,
		Password: testPassword,
	}

	// Act
	// Insert test user into database
	err = userCollection.CreateUser(&user)
	// Assert error
	if err != nil {
		t.Errorf("Unable to create new user: %s", err)
	}

	// Assert collections
	count := 0
	var userDocs []goauthpkg.UserObj
	ctx, cancelFunc := goauthmongo.CtxCreator(5)
	defer cancelFunc()
	cursor, err := client.GetCollection(dbName, userCollName).Find(ctx, bson.M{})
	if err != nil {
		t.Errorf("Error while counting user documents: %s", err)
	}

	for cursor.Next(ctx) {
		count++
		var user goauthpkg.UserObj
		err = cursor.Decode(&user)
		if err != nil {
			t.Errorf("Error while getting user: %s", err)
		}
		userDocs = append(userDocs, user)
	}

	if err = cursor.Err(); err != nil {
		t.Errorf("Error while getting user: %s", err)
	}

	if count != 1 {
		t.Errorf("Incorrect number of results. Expected 1, got %d", count)
	} else if userDocs[0].Username != testUsername {
		t.Errorf("Incorrect Username. Expected %s, Got: %s", testUsername, userDocs[0].Username)
	}
}
