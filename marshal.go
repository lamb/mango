package mango

import (
	"errors"
	"reflect"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/mgocompat"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	STRUCT_TAG_BSON      = "bson"
	OBJECT_ID_FIELD_NAME = "_id"
)

var ErrMustHaveIdField = errors.New("must have id field")

func transform(object interface{}) (bson.D, error) {
	data, err := bson.MarshalWithRegistry(mgocompat.RegistryRespectNilValues, object)
	if err != nil {
		return nil, err
	}
	var d bson.D
	err = bson.UnmarshalWithRegistry(mgocompat.RegistryRespectNilValues, data, &d)
	if err != nil {
		return nil, err
	}
	filter := bson.D{}
	for _, e := range d {
		valueOf := reflect.ValueOf(e.Value)
		if valueOf.IsValid() && !valueOf.IsZero() && e.Key != OBJECT_ID_FIELD_NAME {
			filter = append(filter, e)
		}
	}
	return filter, nil
}

func reflectGetId(object interface{}) string {
	var id string
	typeOfObject := reflect.TypeOf(object)
	valueOfObject := reflect.ValueOf(object)
	if typeOfObject.Kind() == reflect.Ptr {
		typeOfObject, valueOfObject = typeOfObject.Elem(), valueOfObject.Elem()
	}
	for i := 0; i < typeOfObject.NumField(); i++ {
		fieldType := typeOfObject.Field(i)
		fieldValue := valueOfObject.Field(i)
		tag := fieldType.Tag.Get(STRUCT_TAG_BSON)
		if strings.Contains(tag, OBJECT_ID_FIELD_NAME) {
			id = fieldValue.String()
		}
	}
	return id
}

func reflectSetId(object interface{}, id interface{}) {
	typeOfObject := reflect.TypeOf(object)
	valueOfObject := reflect.ValueOf(object)
	if typeOfObject.Kind() == reflect.Ptr {
		typeOfObject, valueOfObject = typeOfObject.Elem(), valueOfObject.Elem()
		for i := 0; i < typeOfObject.NumField(); i++ {
			fieldType := typeOfObject.Field(i)
			fieldValue := valueOfObject.Field(i)
			tag := fieldType.Tag.Get(STRUCT_TAG_BSON)
			if fieldValue.CanSet() && strings.Contains(tag, OBJECT_ID_FIELD_NAME) {
				if objectId, ok := id.(primitive.ObjectID); ok {
					fieldValue.SetString(objectId.Hex())
				} else {
					fieldValue.Set(reflect.ValueOf(id))
				}
			}
		}
	}
}

func cloneWithoutId(object interface{}) interface{} {
	srcType, srcValue := reflect.TypeOf(object), reflect.ValueOf(object)
	if srcType.Kind() == reflect.Ptr {
		srcType, srcValue = srcType.Elem(), srcValue.Elem()
	}
	dstValuePointer := reflect.New(srcType)
	dstValue := dstValuePointer.Elem()
	dstValue.Set(srcValue)
	for i := 0; i < srcType.NumField(); i++ {
		fieldType := srcType.Field(i)
		fieldValue := dstValue.Field(i)
		tag := fieldType.Tag.Get(STRUCT_TAG_BSON)
		if fieldValue.CanSet() && strings.Contains(tag, OBJECT_ID_FIELD_NAME) {
			fieldValue.Set(reflect.New(fieldType.Type).Elem())
			break
		}
	}
	return dstValue.Interface()
}

func reflectGetIdFilter(object interface{}) (bson.D, error) {
	id := reflectGetId(object)
	if len(id) == 0 {
		return nil, ErrMustHaveIdField
	}
	var filter bson.D
	if objectId, err := primitive.ObjectIDFromHex(id); err != nil {
		filter = bson.D{{Key: OBJECT_ID_FIELD_NAME, Value: id}}
	} else {
		filter = bson.D{{Key: OBJECT_ID_FIELD_NAME, Value: objectId}}
	}
	return filter, nil
}
