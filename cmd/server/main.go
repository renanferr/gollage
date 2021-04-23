package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/renanferr/gollage/pkg/albums"

	"github.com/renanferr/gollage/pkg/collage"
	"github.com/renanferr/gollage/pkg/rest"
	"github.com/shkh/lastfm-go/lastfm"
	yaml "gopkg.in/yaml.v2"
)

type Credentials struct {
	LastFM LastFMCredentials `yaml:"lastfm"`
}

type LastFMCredentials struct {
	Key    string `yaml:"key"`
	Secret string `yaml:"secret"`
}

var (
	credentials *Credentials
)

func main() {
	err := loadCredentials()
	if err != nil {
		panic(err)
	}

	api := lastfm.New(credentials.LastFM.Key, credentials.LastFM.Secret)
	albums := albums.New(api)
	collage := collage.New()

	router := rest.Handler(collage, albums)

	port := os.Getenv("PORT")
	if !strings.HasPrefix(port, ":") {
		port = fmt.Sprintf(":%s", port)
	}

	go func() {
		log.Printf("starting server at %s", port)
		log.Fatal(http.ListenAndServe(port, router))
	}()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT)
	<-ch

}

func loadCredentials() error {
	credentialsPath := os.Getenv("CRED")
	if credentialsPath == "" {
		credentialsPath = "./.credentials.yaml"
	}
	log.Printf("reading credentials from %s\n", credentialsPath)
	b, err := ioutil.ReadFile(credentialsPath)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(b, &credentials)
}
