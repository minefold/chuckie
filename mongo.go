package main

import (
	"fmt"
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

func readUrlForServer(id bson.ObjectId) (url string, err error) {
	session, db, err := openMongoSession(os.Getenv("MONGO_URL"))
	if err != nil {
		return
	}
	defer session.Close()

  fmt.Println(os.Getenv("MONGO_URL"), id)

	var server map[string]interface{}
	err = db.C("servers").
		FindId(id).
		One(&server)
	if err != nil {
		return
	}
	fmt.Println("server:", id, "snapshot_id:", server["snapshot_id"])

	var snapshot map[string]interface{}
	err = db.C("snapshots").
		FindId(server["snapshot_id"]).
		One(&snapshot)

	url = fmt.Sprintf("%v", snapshot["url"])
	fmt.Println("server:", id, "url:", url)

	return
}
