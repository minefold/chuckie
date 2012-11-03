package main

import (
	"fmt"
	"labix.org/v2/mgo/bson"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
)

const urlPrefix = "/worlds/"

func worldsHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path[len(urlPrefix):]

	// match /508b5c0ab75f04080000007b.zip
	route := regexp.MustCompile(`^([\w]{24})\.zip$`)

	if match := route.FindStringSubmatch(path); match != nil {
		streamWorld(w, r, match[1])
		fmt.Println("[GET]", r.URL.Path, "– 200")
	} else {
		fmt.Println("[GET]", r.URL.Path, "– 404")
		http.NotFound(w, r)
	}
}

func streamWorld(w http.ResponseWriter, r *http.Request, worldId string) {
	hexId := bson.ObjectIdHex(worldId)
	s3key, err := readS3KeyForWorld(hexId)
	if err != nil {
		fmt.Println("no s3 key for world", hexId.Hex())
		http.NotFound(w, r)
	}

	tempPath, err := createTempDir(worldId)
	if err != nil {
		fmt.Println("failed to create tempdir", err)
		http.Error(w, "failed to create tempdir", 500)
	}

	// TODO currently waits until all files have been extracted
	// it could start zipping files straight away
	err = restoreDir(s3key, tempPath)
	if err != nil {
		fmt.Println("failed to download archive", err)
		http.Error(w, "failed to download archive", 500)
	}

	err = zipPath(w, tempPath)
	if err != nil {
		fmt.Println("failed to zip path", tempPath)
		http.Error(w, "failed zip path", 500)
	}
}

func createTempDir(name string) (string, error) {
	path := filepath.Join(os.TempDir(), name)
	err := exec.Command("rm", "-rf", path).Run()
	if err != nil {
		return "", err
	}
	return path, exec.Command("mkdir", "-p", path).Run()
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[GET]", r.URL.Path, "– 404")
	http.NotFound(w, r)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}
	http.HandleFunc(urlPrefix, worldsHandler) // redirect all urls to the handler function
	http.HandleFunc("/", notFoundHandler)

	fmt.Println("listening")

	http.ListenAndServe(":"+port, nil) // listen for connections at port 9999 on the local machine
}
