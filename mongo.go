package main

import (
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"net/url"
	"os"
	"strings"
)

func openMongoSession(mongoUrl string) (session *mgo.Session, db *mgo.Database, err error) {
	session, err = mgo.Dial(mongoUrl)
	if err != nil {
		return
	}

	url, err := url.Parse(mongoUrl)
	if err != nil {
		return
	}

	dbName := strings.TrimLeft(url.RequestURI(), "/")
	db = session.DB(dbName)
	return
}

func readS3KeyForWorld(id bson.ObjectId) (s3key string, err error) {
	session, db, err := openMongoSession(os.Getenv("MONGO_URI"))
	if err != nil {
		return
	}
	defer session.Close()

	var world map[string]string
	err = db.C("worlds").
		FindId(id).
		Select(bson.M{"world_data_file": 1}).
		One(&world)

	s3key = world["world_data_file"]
	return
}
