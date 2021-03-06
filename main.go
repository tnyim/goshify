package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/boltdb/bolt"
	"github.com/gorilla/mux"
	"github.com/russross/blackfriday"
	"github.com/satori/go.uuid"
)

var db *bolt.DB

func GetContent(id string) ([]byte, error) {
	var text []byte
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Texts"))
		c := b.Get([]byte(id))
		if c == nil {
			return fmt.Errorf("not found")
		}
		// byte slices returned by bolt are only valid inside their transaction
		text = make([]byte, len(c))
		copy(text, c)
		return nil
	})
	return text, err
}

func PutContent(content []byte) (string, error) {
	id := uuid.NewV4().String()

	err := db.Batch(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Texts"))
		return b.Put([]byte(id), content)
	})
	return id, err
}

func SendHTML(w http.ResponseWriter, r *http.Request, b64 string) {
	d, err := base64.StdEncoding.DecodeString(b64)
	w.Write([]byte("<html><head><link rel=stylesheet href=\"/style.css\" /></head><body><div id=\"markup\" style=\"margin-left: 0px;\">"))
	if err != nil {
		// it's probably markdown already, just show as-is
		w.Write(blackfriday.MarkdownCommon([]byte(b64)))
	} else {
		w.Write(blackfriday.MarkdownCommon(d))
	}
	w.Write([]byte("</div></body>"))
}

func URLtoMarkdown(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	SendHTML(w, r, vars["base64"])
}

func StoreContent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, err := PutContent([]byte(vars["base64"]))
	if err == nil {
		w.Write([]byte("<html><body>Success, your markdown's ID is "))
		w.Write([]byte(id))
		w.Write([]byte("</body>"))
	} else {
		w.WriteHeader(500)
		fmt.Println(err)
	}
}

func StoreContentPost(w http.ResponseWriter, r *http.Request) {
	content, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	// only allow text uploads
	if !strings.HasPrefix(http.DetectContentType(content), "text/") {
		w.WriteHeader(415)
		return
	}
	id, err := PutContent(content)
	if err == nil {
		w.Write([]byte("<html><body>Success, your markdown's ID is " + id + "</body>"))
	} else {
		w.WriteHeader(500)
		fmt.Println(err)
	}
}

func LoadContent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	text, err := GetContent(vars["id"])

	if err == nil {
		SendHTML(w, r, string(text))
	} else {
		w.WriteHeader(404)
	}
}

func LoadRawContent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	text, err := GetContent(vars["id"])

	if err == nil {
		w.Write(text)
	} else {
		w.WriteHeader(404)
	}
}

func LoadHome(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadFile("home.txt")

	if err == nil {
		SendHTML(w, r, string(b))
	} else {
		fmt.Println(err)
		w.WriteHeader(404)
	}
}

func main() {
	var err error
	db, err = bolt.Open("goshify.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("Texts"))
		return err
	})
	if err != nil {
		panic(err)
	}

	sigs := make(chan os.Signal)
	signal.Notify(sigs, syscall.SIGTERM)
	signal.Notify(sigs, syscall.SIGINT)

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/d/{base64:(?s).*}", URLtoMarkdown).Methods("GET")
	router.HandleFunc("/s/{base64:(?s).*}", StoreContent).Methods("GET")
	router.HandleFunc("/l/{id}", LoadContent).Methods("GET")
	router.HandleFunc("/r/{id}", LoadRawContent).Methods("GET")
	router.HandleFunc("/", LoadHome).Methods("GET")
	router.HandleFunc("/style.css", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, r.URL.Path[1:])
	})
	router.StrictSlash(false)
	router.HandleFunc("/s/", StoreContentPost).Methods("POST")
	router.HandleFunc("/s", StoreContentPost).Methods("POST")
	go http.ListenAndServe(":8089", router)
	fmt.Println("Listening...")
	<-sigs
}
