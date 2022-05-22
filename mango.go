package mango

import (
	"context"
	"reflect"
	"strconv"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Collection struct {
	Driver *mongo.Collection
}

func NewCollection(database *mongo.Database, name string, opts ...*options.CollectionOptions) *Collection {
	driver := database.Collection(name, opts...)
	return &Collection{Driver: driver}
}

func (c *Collection) FindById(ctx context.Context, object interface{}, opts ...*options.FindOneOptions) (interface{}, error) {
	filter, err := reflectGetIdFilter(object)
	if err != nil {
		return nil, err
	}
	result := c.Driver.FindOne(ctx, filter, opts...)
	if result.Err() != nil {
		return nil, result.Err()
	}
	pointer := reflect.New(reflect.TypeOf(object))
	err = result.Decode(pointer.Interface())
	if err != nil {
		return nil, err
	}
	return pointer.Elem().Interface(), err
}

func (c *Collection) Find(ctx context.Context, object interface{}, opts ...*options.FindOptions) ([]interface{}, error) {
	filter, err := transform(object)
	if err != nil {
		return nil, err
	}
	cursor, err := c.Driver.Find(ctx, filter, opts...)
	if err != nil {
		return nil, err
	}
	var results []interface{}
	for cursor.Next(context.TODO()) {
		pointer := reflect.New(reflect.TypeOf(object))
		if err := cursor.Decode(pointer.Interface()); err != nil {
			return nil, err
		}
		results = append(results, pointer.Elem().Interface())
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	if err := cursor.Close(ctx); err != nil {
		return nil, err
	}
	return results, nil
}

func (c *Collection) Page(ctx context.Context, filter interface{}, sort interface{}, page string, limit string) (cur *mongo.Cursor, err error) {
	pageInt, err := strconv.ParseInt(page, 10, 64)
	if err != nil || pageInt < 1 {
		pageInt = 1
	}
	limitInt, err := strconv.ParseInt(limit, 10, 64)
	if err != nil || limitInt <= 0 {
		limitInt = 1
	}
	skip := (pageInt - 1) * limitInt
	opts := options.Find().SetSort(sort).SetSkip(skip).SetLimit(limitInt)
	return c.Driver.Find(ctx, filter, opts)
}

func (c *Collection) InsertOne(ctx context.Context, object interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	result, err := c.Driver.InsertOne(ctx, object, opts...)
	if err != nil {
		return nil, err
	}
	reflectSetId(object, result.InsertedID)
	return result, nil
}

func (c *Collection) UpdateById(ctx context.Context, object interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	filter, err := reflectGetIdFilter(object)
	if err != nil {
		return nil, err
	}
	update, err := transform(object)
	update = bson.D{{Key: "$set", Value: update}}
	if err != nil {
		return nil, err
	}
	return c.Driver.UpdateOne(ctx, filter, update, opts...)
}

func (c *Collection) ReplaceById(ctx context.Context, object interface{}, opts ...*options.ReplaceOptions) (*mongo.UpdateResult, error) {
	filter, err := reflectGetIdFilter(object)
	if err != nil {
		return nil, err
	}
	replacement := cloneWithoutId(object)
	return c.Driver.ReplaceOne(ctx, filter, replacement, opts...)
}

func (c *Collection) DeleteById(ctx context.Context, object interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	filter, err := reflectGetIdFilter(object)
	if err != nil {
		return nil, err
	}
	return c.Driver.DeleteOne(ctx, filter, opts...)
}
