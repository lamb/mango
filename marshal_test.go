package mango

import (
	"testing"

	"github.com/basjoofan/compose/model"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestTransform(t *testing.T) {
	d, err := transform(model.Space{Name: "test"})
	assert.Nil(t, err)
	assert.Len(t, d, 1)
}

func TestReflectGetId(t *testing.T) {
	space := model.Space{Id: "6277cc2316b97479f315e797"}
	id := reflectGetId(space)
	assert.NotNil(t, id)
	assert.Equal(t, "6277cc2316b97479f315e797", id)
	id = reflectGetId(&space)
	assert.NotNil(t, id)
	assert.Equal(t, "6277cc2316b97479f315e797", id)
}

func TestReflectSetId(t *testing.T) {
	space := model.Space{}
	reflectSetId(space, "6277cc2316b97479f315e797")
	assert.Zero(t, space.Id)
	reflectSetId(&space, "6277cc2316b97479f315e797")
	assert.NotNil(t, space.Id)
	assert.Equal(t, "6277cc2316b97479f315e797", space.Id)
}

func TestCloneWithoutId(t *testing.T) {
	src := model.Space{Id: "6277cc2316b97479f315e797", Name: "test"}
	dst := cloneWithoutId(src)
	assert.NotEqual(t, src, dst)
	assert.Zero(t, dst.(model.Space).Id)
	dst = cloneWithoutId(&src)
	assert.NotEqual(t, src, dst)
	assert.Zero(t, dst.(model.Space).Id)
}

func TestReflectGetIdFilter(t *testing.T) {
	objectId, err := primitive.ObjectIDFromHex("6277cc2316b97479f315e797")
	assert.Nil(t, err)
	space := model.Space{Id: "6277cc2316b97479f315e797"}
	filter, _ := reflectGetIdFilter(space)
	assert.NotNil(t, filter)
	assert.Equal(t, objectId, filter.Map()[OBJECT_ID_FIELD_NAME])
	filter, _ = reflectGetIdFilter(&space)
	assert.NotNil(t, filter)
	assert.Equal(t, objectId, filter.Map()[OBJECT_ID_FIELD_NAME])
	// test for not an object id
	notAnObjectId := "notanobjectid"
	space = model.Space{Id: notAnObjectId}
	filter, _ = reflectGetIdFilter(space)
	assert.NotNil(t, filter)
	assert.Equal(t, notAnObjectId, filter.Map()[OBJECT_ID_FIELD_NAME])
	filter, _ = reflectGetIdFilter(&space)
	assert.NotNil(t, filter)
	assert.Equal(t, notAnObjectId, filter.Map()[OBJECT_ID_FIELD_NAME])
	notHaveId := NotHaveId{Name: "test"}
	filter, err = reflectGetIdFilter(notHaveId)
	assert.NotNil(t, err)
	assert.Nil(t, filter)
}

type NotHaveId struct {
	Name string
}
