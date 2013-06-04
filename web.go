package main

import (
	"archive/zip"
	"fmt"
	"io"
	"labix.org/v2/mgo/bson"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sync"
	"time"
)

const urlPrefix = "/worlds/"

func worldsHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path[len(urlPrefix):]

	// match /508b5c0ab75f04080000007b.zip
	route := regexp.MustCompile(`^([\w]{24})\.zip$`)

	if match := route.FindStringSubmatch(path); match != nil {
		filename := r.URL.Query().Get("name")
		if filename == "" {
			filename = match[1]
		}
		streamWorld(w, r, match[1], filename)
		fmt.Println("[GET]", r.URL.Path, "– 200")
	} else {
		fmt.Println("[GET]", r.URL.Path, "– 404")
		http.NotFound(w, r)
	}
}

func streamWorld(w http.ResponseWriter, r *http.Request, worldId string, filename string) {
	w.Header().Add("Content-Type", "application/zip")
	w.Header().Add("Content-Disposition",
		fmt.Sprintf("attachment; filename=%s.zip", filename))

	hexId := bson.ObjectIdHex(worldId)
	url, err := readUrlForServer(hexId)
	if err != nil {
		fmt.Println("no url for world", hexId.Hex(), err)
		http.NotFound(w, r)
	}

	tempPath, err := createTempDir(worldId)
	if err != nil {
		fmt.Println("failed to create tempdir", err)
		http.Error(w, "failed to create tempdir", 500)
	}

	zip := zip.NewWriter(w)
	defer zip.Close()

	// Zip files as they are downloaded
	watcher := NewWatcher(tempPath, 5*time.Second)
	go watcher.Watch()
	var wg sync.WaitGroup
	wg.Add(1)
	go zipFilesAsTheyAppear(tempPath, filename, watcher.C, zip, &wg)

	// start the download
	err = restoreDir(url, tempPath)
	if err != nil {
		fmt.Println("failed to download archive", err)
		http.Error(w, "failed to download archive", 500)
	}

	watcher.Cancel()
	wg.Wait()
}

func zipFilesAsTheyAppear(root, prefix string, files chan string, zip *zip.Writer, wg *sync.WaitGroup) {
	defer wg.Done()
	for path := range files {
		fileName := path[len(root)+1:]
		zipF, err := zip.Create(prefix + "/" + fileName)
		if err != nil {
			fmt.Println("failed to write zip header", fileName)
			return
		}

		f, err := os.Open(path)
		if err != nil {
			fmt.Println("failed to open path", path, err)
			return
		}
		_, err = io.Copy(zipF, f)
		if err != nil {
			fmt.Println("failed to zip path", path, err)
			return
		}
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

	fmt.Println("listening on " + port)

	http.ListenAndServe(":"+port, nil) // listen for connections at port 9999 on the local machine
}
