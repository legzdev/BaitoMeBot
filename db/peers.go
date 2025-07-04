package db

import (
	"errors"
	"strconv"

	"github.com/boltdb/bolt"
	"github.com/gotd/td/tg"
)

var Peers *peersDB

type peersDB struct {
	db *bolt.DB
}

func NewPeersDB() (*peersDB, error) {
	db, err := bolt.Open("cache.db", 0600, nil)
	if err != nil {
		return nil, err
	}

	tx, err := db.Begin(true)
	if err != nil {
		return nil, err
	}

	_, err = tx.CreateBucketIfNotExists([]byte("channels"))
	if err != nil {
		return nil, err
	}

	_, err = tx.CreateBucketIfNotExists([]byte("users"))
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return &peersDB{db: db}, nil
}

func (db *peersDB) Put(peer tg.InputPeerClass) error {
	tx, err := db.db.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var bucket *bolt.Bucket
	var id int64
	var accessHash int64

	switch v := peer.(type) {
	case *tg.InputPeerChannel:
		bucket = tx.Bucket([]byte("channels"))
		id, accessHash = v.ChannelID, v.AccessHash

	case *tg.InputPeerUser:
		bucket = tx.Bucket([]byte("users"))
		id, accessHash = v.UserID, v.AccessHash
	}

	if id == 0 || accessHash == 0 {
		return errors.New("peersDB: unexpected peer type")
	}

	if bucket == nil {
		return errors.New("peersDB: nil bucket")
	}

	key := strconv.FormatInt(id, 10)
	value := strconv.FormatInt(accessHash, 10)

	err = bucket.Put([]byte(key), []byte(value))
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (db *peersDB) GetChannel(id int64) (*tg.InputPeerChannel, error) {
	accessHash, err := db.get(id, "channels")
	if err != nil || accessHash == 0 {
		return nil, err
	}

	return &tg.InputPeerChannel{
		ChannelID:  id,
		AccessHash: accessHash,
	}, nil
}

func (db *peersDB) GetUser(id int64) (*tg.InputPeerUser, error) {
	accessHash, err := db.get(id, "users")
	if err != nil || accessHash == 0 {
		return nil, err
	}

	return &tg.InputPeerUser{
		UserID:     id,
		AccessHash: accessHash,
	}, nil
}

func (db *peersDB) get(id int64, bucketName string) (int64, error) {
	tx, err := db.db.Begin(false)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	bucket := tx.Bucket([]byte(bucketName))
	if bucket == nil {
		return 0, errors.New("peersDB: nil bucket")
	}

	key := strconv.FormatInt(id, 10)

	value := bucket.Get([]byte(key))
	if value == nil {
		return 0, nil
	}

	return strconv.ParseInt(string(value), 10, 64)
}
