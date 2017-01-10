package gaas

import (
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"path"
	"sync"
	"time"

	"github.com/manki/flickgo"
)

var goatDir = "goat-pics"

func Run() {
	rand.Seed(time.Now().UTC().UnixNano())

	if !isDir(goatDir) {
		mkdir(goatDir, 0700)
		getGoats()
	}
	http.HandleFunc("/", handler)
	// TODO: doesn't seem to be serving to remote hosts,
	// use http.Server instance instead?
	http.ListenAndServe("0.0.0.0:8080", nil)

}

func downloadFile(wg *sync.WaitGroup, url string) {
	defer wg.Done()

	filePath := goatDir + "/" + path.Base(url)

	// Create the file
	out, err := os.Create(filePath)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer resp.Body.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("Downloaded ", filePath)
}

// mkdir util function to create dir
// Creates dir if it doesn't exist, similar to mkdir -p
func mkdir(path string, perms os.FileMode) error {
	// create dir if it does not exist
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		err = os.MkdirAll(path, perms)
		if err != nil {
			return err
		}
	}
	return nil
}

// Downloads pictures of goats!
func getGoats() {
	c := flickgo.New(os.Getenv("FLICKR_KEY"), os.Getenv("FLICKR_SECRET"),
		http.DefaultClient)
	resp, err := c.Search(map[string]string{
		"tags":        "goat",
		"safe-search": "1",
	})
	if err != nil {
		panic(err.Error())
	}

	wg := &sync.WaitGroup{}
	wg.Add(len(resp.Photos))
	for _, p := range resp.Photos {
		// Note URL() doesn't seem to work with flickgo.SizeOriginal, get bad urls
		// We're Downloading these to disk but we probably want them in ram or
		// some more efficient storage
		go downloadFile(wg, p.URL(flickgo.SizeMedium500))
	}
	wg.Wait()
	fmt.Println("Finished downloading")
}

// isDir returns true if the path exists and is a directory.
func isDir(path string) bool {
	fi, err := os.Stat(path)
	return err == nil && fi.IsDir()
}

func getGoatPath() string {
	goats, err := ioutil.ReadDir(goatDir)
	if err != nil {
		panic(err.Error())
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
