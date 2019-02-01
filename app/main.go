package main

import (
	goauthcrypto "chotot/go_auth/pkg/crypto"
	goauthmongo "chotot/go_auth/pkg/mongo"
	goauthredis "chotot/go_auth/pkg/redis"
	goauthserver "chotot/go_auth/pkg/server"
	"log"
	"os"
)

func main() {
	var redisPass = ""
	if len(os.Args) == 2 {
		redisPass = os.Args[1]
	}

	client, err := goauthmongo.NewClient("mongodb://127.0.0.1:27017")
	if err != nil {
		log.Fatalln("Unable to connect to MongoDB server")
	}
	defer client.Disconnect()

	redisClient, err := goauthredis.StartRedis("localhost:6379", redisPass)
	if err != nil {
		log.Fatalln("Unable to connect to Redis server")
	}
	defer goauthredis.CloseConn(redisClient)

	hash := goauthcrypto.HashStr{}
	userCollection := goauthmongo.NewUserCollection(client, "go_auth", "users", &hash)

	// Load all users to Redis
	userObjArr, err := userCollection.GetAllUser()
	if err != nil {
		log.Fatalf("Problem retreive users from mongoDB")
	}
	err = redisClient.SaveUsers(userObjArr)
	if err != nil {
		log.Fatalf("Problem adding users to Redis")
	}

	server := goauthserver.NewServer(userCollection, redisClient)

	server.Start()
}
