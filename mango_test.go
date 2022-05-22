package mango

import (
	"context"
	"reflect"
	"testing"

	"github.com/basjoofan/compose/model"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func newClient() *mongo.Client {
	uri := "mongodb://localhost:27017"
	// create a new client and connect to the server
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	return client
}

func disconnect(client *mongo.Client) {
	if err := client.Disconnect(context.TODO()); err != nil {
		panic(err)
	}
}

func TestFindById(t *testing.T) {
	client := newClient()
	defer disconnect(client)
	database := client.Database("basjoofan")
	space := NewCollection(database, "space")
	spaceModel := model.Space{Id: "6277cc2316b97479f315e797"}
	spaceResult, _ := space.FindById(context.TODO(), spaceModel)
	assert.Equal(t, "6277cc2316b97479f315e797", spaceResult.(model.Space).Id)
	assert.Equal(t, reflect.TypeOf(spaceModel), reflect.TypeOf(spaceResult))
	spaceModel = model.Space{Id: "6277cc2316b97479f315e7971"}
	spaceResult, _ = space.FindById(context.TODO(), spaceModel)
	assert.Nil(t, spaceResult)
	spaceModel = model.Space{Id: "6277cc2316b97479f315e798"}
	spaceResult, _ = space.FindById(context.TODO(), spaceModel)
	assert.Nil(t, spaceResult)
}

func TestFind(t *testing.T) {
	client := newClient()
	defer disconnect(client)
	database := client.Database("basjoofan")
	space := NewCollection(database, "space")
	spaceModel := model.Space{Name: "test"}
	spaceResults, _ := space.Find(context.TODO(), spaceModel)
	assert.Equal(t, reflect.TypeOf(spaceModel), reflect.TypeOf(spaceResults[0]))
}

func TestPage(t *testing.T) {
	client := newClient()
	defer disconnect(client)
	database := client.Database("basjoofan")
	space := NewCollection(database, "space")
	userId := "6273f351686dbd7a80c5ade9"
	filter := bson.M{"$or": []bson.M{{"founder": userId}, {"owners": userId}, {"members": userId}, {"guests": userId}}}
	sort := bson.D{{Key: "foundTime", Value: 1}}
	cursor, _ := space.Page(context.TODO(), filter, sort, "a", "-10")
	var spaces []model.Space
	err := cursor.All(context.TODO(), &spaces)
	assert.Nil(t, err)
	assert.NotNil(t, spaces)
}

func TestInsertOne(t *testing.T) {
	client := newClient()
	defer disconnect(client)
	database := client.Database("basjoofan")
	space := NewCollection(database, "space")
	spaceModel := model.Space{Name: "test2"}
	result, _ := space.InsertOne(context.TODO(), &spaceModel)
	assert.Equal(t, spaceModel.Id, result.InsertedID.(primitive.ObjectID).Hex())
	assert.NotNil(t, result.InsertedID.(primitive.ObjectID).Hex())
}

func TestUpdateById(t *testing.T) {
	client := newClient()
	defer disconnect(client)
	database := client.Database("basjoofan")
	space := NewCollection(database, "space")
	spaceModel := model.Space{Id: "6277cc2316b97479f315e798"}
	result, _ := space.UpdateById(context.TODO(), &spaceModel)
	assert.NotNil(t, result)
	result, _ = space.UpdateById(context.TODO(), spaceModel)
	assert.NotNil(t, result)
}

func TestReplaceById(t *testing.T) {
	client := newClient()
	defer disconnect(client)
	database := client.Database("basjoofan")
	space := NewCollection(database, "space")
	spaceModel := model.Space{Id: "62768dfab3abf1aadd39d56d", Name: "testReplace"}
	result, _ := space.ReplaceById(context.TODO(), &spaceModel)
	assert.NotNil(t, result)
	result, _ = space.ReplaceById(context.TODO(), spaceModel)
	assert.NotNil(t, result)
}

func TestDeleteById(t *testing.T) {
	client := newClient()
	defer disconnect(client)
	database := client.Database("basjoofan")
	space := NewCollection(database, "space")
	spaceModel := model.Space{Id: "6273fad54e19dbc0dd20f116", Name: "testReplace"}
	result, _ := space.DeleteById(context.TODO(), &spaceModel)
	assert.NotNil(t, result)
	result, _ = space.DeleteById(context.TODO(), spaceModel)
	assert.NotNil(t, result)
}
