package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/soundcloud/periskop/api"
	"github.com/soundcloud/periskop/config"
	"github.com/soundcloud/periskop/repository"
	"github.com/soundcloud/periskop/scraper"
	"github.com/soundcloud/periskop/servicediscovery"
)

const numOfProcessors = 8

func main() {

	var (
		port              = flag.String("port", os.Getenv("PORT"), "The server port")
		configurationFile = flag.String("config", os.Getenv("CONFIG_FILE"), "The configuration file")
	)

	flag.Parse()

	basePath, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Using baseDir %s", basePath)

	if _, err := os.Stat(*configurationFile); err != nil {
		log.Fatalf("Invalid configuration file %s", *configurationFile)
	}

	log.Printf("Using configFile %s", *configurationFile)
	cfg, err := config.LoadFile(*configurationFile)
	if err != nil {
		panic(err)
	}

	processor := scraper.NewProcessor(numOfProcessors)
	processor.Run()
	repo := repository.NewInMemory()
	for _, service := range cfg.Services {
		resolver := servicediscovery.NewResolver(service)
		s := scraper.NewScraper(resolver, &repo, service, processor)
		go s.Scrape()
	}

	http.HandleFunc("/-/health", healthHandler)

	webFolder := filepath.Join(basePath, "web/dist")
	log.Printf("Using webFolder %s", webFolder)

	fs := http.FileServer(http.Dir(webFolder))
	http.Handle("/", fs)

	http.Handle("/services/", api.NewHandler(repo))

	address := fmt.Sprintf(":%s", *port)
	log.Printf("Serving on address %s", address)
	log.Fatal(http.ListenAndServe(address, nil))
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("OK"))
	if err != nil {
		log.Fatalf("error running health handler")
	}
}
