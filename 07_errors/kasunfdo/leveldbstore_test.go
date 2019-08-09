package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewLeveldbStoreErrOpen(t *testing.T) {
	os.RemoveAll("./test")
	_, err := NewLeveldbStore("", "")
	fmt.Println(err.Error())
	assert.Equal(t, "internal error\n\tmkdir /index.db: permission denied", err.Error())
}

func TestNewLeveldbStoreInitialize(t *testing.T) {
	store, _ := NewLeveldbStore("test", "puppy.db")
	_, _ = store.CreatePuppy(Puppy{Breed: "Labrador", Colour: "Cream", Value: 2999.99})
	store, _ = NewLeveldbStore("test", "puppy.db")
	id, _ := store.CreatePuppy(Puppy{Breed: "Labrador", Colour: "Cream", Value: 2999.99})

	assert.True(t, id == 2)
	os.RemoveAll("./test")
}

func TestLeveldbStoreCreatePuppyErrOpen(t *testing.T) {
	store := LeveldbStore{path: ""}
	_, err := store.CreatePuppy(Puppy{})
	assert.Equal(t, "internal error\n\topen /LOCK: permission denied", err.Error())

	store.path = "test"
	_, err = store.CreatePuppy(Puppy{})
	assert.Equal(t, "internal error\n\tresource temporarily unavailable", err.Error())

	os.RemoveAll("./test")
}

func TestLeveldbStoreReadPuppyErrOpen(t *testing.T) {
	store := LeveldbStore{path: ""}
	_, err := store.ReadPuppy(1)
	assert.Equal(t, "internal error\n\topen /LOCK: permission denied", err.Error())

	os.RemoveAll("./test")
}

func TestLeveldbStoreUpdatePuppyErrOpen(t *testing.T) {
	store := LeveldbStore{path: ""}
	err := store.UpdatePuppy(Puppy{})
	assert.Equal(t, "internal error\n\topen /LOCK: permission denied", err.Error())

	os.RemoveAll("./test")
}

func TestLeveldbStoreDeletePuppyErrOpen(t *testing.T) {
	store := LeveldbStore{path: ""}
	err := store.DeletePuppy(1)
	assert.Equal(t, "internal error\n\topen /LOCK: permission denied", err.Error())

	os.RemoveAll("./test")
}
