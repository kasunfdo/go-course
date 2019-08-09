package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/syndtr/goleveldb/leveldb"
)

type LeveldbStore struct {
	path          string
	dataFileName  string
	indexFileName string
	indexName     string
	nextID        uint64
}

func NewLeveldbStore(path, dataFileName string) (*LeveldbStore, error) {
	store := &LeveldbStore{
		nextID:        1,
		path:          path,
		dataFileName:  dataFileName,
		indexFileName: "index.db",
		indexName:     "last_index",
	}

	db, err := openDB(store.path+"/"+store.indexFileName, true)
	if err != nil {
		return nil, NewError(ErrInternal, err)
	}
	defer closeDB(db)

	exists, err := db.Has([]byte(store.indexName), nil)
	if !exists {
		log.Printf("Initialized store with index: %v\n", store.nextID)
		return store, nil
	}
	if err != nil {
		return nil, NewError(ErrInternal, err)
	}

	value, err := db.Get([]byte(store.indexName), nil)
	if err != nil {
		return nil, NewError(ErrInternal, err)
	}

	store.nextID, err = strconv.ParseUint(string(value), 10, 64)
	if err != nil {
		return nil, NewError(ErrInternal, err)
	}

	log.Printf("Initialized store with index: %v", store.nextID)

	return store, nil
}

func (l *LeveldbStore) CreatePuppy(puppy Puppy) (uint64, error) {
	if puppy.Value < 0 {
		return 0, NewError(ErrInvalid, nil, "value of puppy is negative")
	}

	dataDB, err := openDB(l.path+"/"+l.dataFileName, false)
	if err != nil {
		return 0, NewError(ErrInternal, err)
	}
	defer closeDB(dataDB)

	indexDB, err := openDB(l.path+"/"+l.indexFileName, false)
	if err != nil {
		return 0, NewError(ErrInternal, err)
	}
	defer closeDB(indexDB)

	puppy.ID = l.nextID
	l.nextID++

	encodedPuppy, err := json.Marshal(puppy)
	if err != nil {
		return 0, NewError(ErrInternal, err)
	}

	err = dataDB.Put([]byte(strconv.FormatUint(puppy.ID, 10)), encodedPuppy, nil)
	if err != nil {
		return 0, NewError(ErrInternal, err)
	}

	err = indexDB.Put([]byte(l.indexName), []byte(strconv.FormatUint(l.nextID, 10)), nil)
	if err != nil {
		return 0, NewError(ErrInternal, err)
	}

	return puppy.ID, nil
}

func (l *LeveldbStore) ReadPuppy(id uint64) (Puppy, error) {
	db, err := openDB(l.path+"/"+l.dataFileName, true)
	if err != nil {
		return Puppy{}, NewError(ErrInternal, err)
	}
	defer closeDB(db)

	encoded, err := db.Get([]byte(strconv.FormatUint(id, 10)), nil)
	if err != nil {
		return Puppy{}, NewError(ErrNotFound, nil, fmt.Sprintf("puppy with id: %v is not found", id))
	}

	puppy := &Puppy{}
	err = json.Unmarshal(encoded, puppy)
	if err != nil {
		return Puppy{}, NewError(ErrInternal, err)
	}

	return *puppy, nil
}

func (l *LeveldbStore) UpdatePuppy(puppy Puppy) error {
	if puppy.Value < 0 {
		return NewError(ErrInvalid, nil, "value of puppy is negative")
	}

	db, err := openDB(l.path+"/"+l.dataFileName, false)
	if err != nil {
		return NewError(ErrInternal, err)
	}
	defer closeDB(db)

	exists, err := db.Has([]byte(strconv.FormatUint(puppy.ID, 10)), nil)
	if !exists {
		return NewError(ErrNotFound, nil, fmt.Sprintf("puppy with id: %v is not found", puppy.ID))
	}
	if err != nil {
		return NewError(ErrInternal, err)
	}

	encodedPuppy, err := json.Marshal(puppy)
	if err != nil {
		return NewError(ErrInternal, err)
	}

	err = db.Put([]byte(strconv.FormatUint(puppy.ID, 10)), encodedPuppy, nil)
	if err != nil {
		return NewError(ErrInternal, err)
	}

	return nil
}

func (l *LeveldbStore) DeletePuppy(id uint64) error {
	db, err := openDB(l.path+"/"+l.dataFileName, false)
	if err != nil {
		return NewError(ErrInternal, err)
	}
	defer closeDB(db)

	exists, err := db.Has([]byte(strconv.FormatUint(id, 10)), nil)
	if !exists {
		return NewError(ErrNotFound, nil, fmt.Sprintf("puppy with id: %v is not found", id))
	}
	if err != nil {
		return NewError(ErrInternal, err)
	}

	err = db.Delete([]byte(strconv.FormatUint(id, 10)), nil)
	if err != nil {
		return NewError(ErrInternal, err)
	}

	return nil
}

func openDB(dbPath string, isReadOnly bool) (*leveldb.DB, error) {
	db, err := leveldb.OpenFile(dbPath, nil)
	if err != nil {
		return nil, err
	}

	if isReadOnly {
		err = db.SetReadOnly()
		if err != nil {
			return nil, err
		}
	}
	return db, nil
}

func closeDB(db *leveldb.DB) {
	err := db.Close()
	if err != nil {
		panic(err)
	}
}
