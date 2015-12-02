package db

import (
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	mongoConnectionTimeout = 5 * time.Second
)

type MongoDb struct {
	sess *mgo.Session
}

func (db *MongoDb) Connect(dsn string) error {
	var err error

	db.sess, err = mgo.DialWithTimeout(dsn, mongoConnectionTimeout)

	return err
}
func (db *MongoDb) Disconnect() {
	db.sess.Close()
}

func (db *MongoDb) Insert(coll string, v ...interface{}) error {
	sess := db.sess.Copy()
	defer sess.Close()

	return sess.DB("").C(coll).Insert(v...)
}

func (db *MongoDb) FindById(coll string, id string, v interface{}) bool {
	sess := db.sess.Copy()
	defer sess.Close()

	return mgo.ErrNotFound != sess.DB("").C(coll).FindId(id).One(v)
}

func (db *MongoDb) FindAll(coll string, v interface{}) error {
	sess := db.sess.Copy()
	defer sess.Close()

	return sess.DB("").C(coll).Find(bson.M{}).All(v)
}

func (db *MongoDb) Update(coll string, id, v interface{}) error {
	sess := db.sess.Copy()
	defer sess.Close()

	return sess.DB("").C(coll).Update(bson.M{"_id": id}, bson.M{"$set": v})
}

var db *MongoDb = &MongoDb{}

func Db() *MongoDb { return db }
