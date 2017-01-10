package gaas

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"time"
)

var goatDir = "goat-pics"

func Run() {
	rand.Seed(time.Now().UTC().UnixNano())

	if !isDir(goatDir) {
		getGoats()
	}
	http.HandleFunc("/", handler)
	// TODO: doesn't seem to be serving to remote hosts,
	// use http.Server instance instead?
	http.ListenAndServe("0.0.0.0:8080", nil)

}

// Downloads pictures of goats!
// TODO: Try the flickr API, or at a push:
// curl -sA "Chrome" -L "https://www.google.com/search?hl=en&tbm=isch&q=goats"
// and scrape

func getGoats() {
}

// isDir returns true if the path exists and is a directory.
func isDir(path string) bool {
	fi, err := os.Stat(path)
	return err == nil && fi.IsDir()
}

func getGoatPath() string {
	goats, err := ioutil.ReadDir(goatDir)
	if err != nil {
		fmt.Errorf(err.Error())
	}
	numGoats := len(goats)

	if numGoats == 0 {
		panic("critical lack of goats error")
	}
	i := rand.Intn(numGoats)
	return goatDir + "/" + goats[i].Name()
}

func handler(w http.ResponseWriter, r *http.Request) {
	var path = getGoatPath()
	fmt.Println("New query, serving", path)
	http.ServeFile(w, r, path)
	fmt.Fprintf(w, "Served a goat, path was %s", r.URL.Path[1:])
}
