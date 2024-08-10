package space

import (
	"fmt"

	"github.com/google/uuid"
	"go.etcd.io/bbolt"
)

var (
	dbname = "space"
)

type BoltFunc func(tx *bbolt.Tx) error

type M map[string]string

type Collection struct {
	*bbolt.Bucket
}

type Space struct {
	db *bbolt.DB
}

func New() (*Space, error) {
	db, err := bbolt.Open(
		fmt.Sprintf("%s.db", dbname),
		0666,
		nil,
	)

	if err != nil {
		return nil, err
	}

	return &Space{
		db,
	}, nil
}

// func (s *Space) CreateCollection(name string) (*Collection, error) {

// 	tx, err := s.db.Begin(true)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer tx.Rollback()

// 	bucket, err := tx.CreateBucketIfNotExists([]byte(name))
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &Collection{Bucket: bucket}, nil
// }

func (s *Space) Insert(bucketName string, v M) (uuid.UUID, error) {
	id := uuid.New()

	tx, err := s.db.Begin(true)
	if err != nil {
		return id, err
	}

	bucket, err := tx.CreateBucketIfNotExists([]byte(bucketName))
	if err != nil {
		return id, err
	}

	for k, _v := range v {
		if err := bucket.Put([]byte(k), []byte(_v)); err != nil {
			tx.Rollback()
			return id, err
		}
	}

	err = bucket.Put([]byte("id"), []byte(id.String()))
	if err != nil {
		tx.Rollback()
		return id, err
	}

	tx.Commit()
	return id, nil
}

func (s *Space) Get(bucketName, key string, query any) (result M, err error) {
	result = make(M)
	tx, err := s.db.Begin(true)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	bucket := tx.Bucket([]byte(bucketName))
	if bucket == nil {
		return nil, fmt.Errorf("bucket (%s) not found", bucketName)
	}

	err = bucket.ForEach(func(k, v []byte) error {
		result[string(k)] = string(v)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

// var userdata = make(map[string]string)
// if err := db.View(func(tx *bbolt.Tx) error {
// 	bucket := tx.Bucket([]byte(userBucket))
// 	if bucket == nil {
// 		return fmt.Errorf("bucket (%s) not found", userBucket)
// 	}

// 	bucket.ForEach(func(k, v []byte) error {
// 		userdata[string(k)] = string(v)
// 		return nil
// 	})

// 	return nil
// }); err != nil {
// 	log.Fatal(err)
// }

// fmt.Println(userdata)
