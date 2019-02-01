package goauthmongo

import (
	"context"
	"time"

	"github.com/mongodb/mongo-go-driver/mongo"
)

// Client Used for client type
type Client struct {
	client *mongo.Client
}

// CtxCreator Create context return context and cancelFunc
func CtxCreator(timeout int) (context.Context, context.CancelFunc) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	return ctx, cancelFunc
}

// NewClient Create new client and return the client
func NewClient(uri string) (*Client, error) {
	ctx, cancelFunc := CtxCreator(5)
	defer cancelFunc()
	client, err := mongo.Connect(ctx, uri)
	return &Client{client}, err
}

// GetCollection Get a collection from a database
func (c *Client) GetCollection(dbName string, collName string) *mongo.Collection {
	coll := c.client.Database(dbName).Collection(collName)
	return coll
}

// Disconnect Close the connection to the server
func (c *Client) Disconnect() {
	if c.client != nil {
		ctx, cancelFunc := CtxCreator(5)
		defer cancelFunc()
		c.client.Disconnect(ctx)
	}
}

// DropDB Drop the named database
func (c *Client) DropDB(dbName string) error {
	ctx, cancelFunc := CtxCreator(5)
	defer cancelFunc()
	return c.client.Database(dbName).Drop(ctx)
}
