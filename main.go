package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/soundcloud/periskop-go"

	"github.com/soundcloud/periskop/api"
	"github.com/soundcloud/periskop/config"
	"github.com/soundcloud/periskop/metrics"
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

	serverURL, envServer := os.LookupEnv("SERVER_URL")
	if !envServer || serverURL == "" {
		serverURL = "localhost"
	}

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

	webFolder := filepath.Join(basePath, "web/dist")
	log.Printf("Using webFolder %s", webFolder)
	fs := http.FileServer(http.Dir(webFolder))

	// API routing
	r := mux.NewRouter()
	r.Handle("/services/",
		api.NewServicesListHandler(&repo)).Methods(http.MethodGet)
	r.Handle("/services/{service_name}/errors/",
		api.NewErrorsListHandler(&repo)).Methods(http.MethodGet)
	r.Handle("/services/{service_name}/errors/{error_key:.*}/",
		api.NewErrorResolveHandler(&repo)).Methods(http.MethodDelete, http.MethodOptions)
	r.PathPrefix("/").Handler(http.StripPrefix("/", fs))
	r.Use(api.CORSLocalhostMiddleware(r))
	http.Handle("/", r)

	// Telemetry endpoints
	errorExporter := periskop.NewErrorExporter(&metrics.ErrorCollector)
	periskopHandler := periskop.NewHandler(errorExporter)

	http.Handle("/metrics", promhttp.Handler())
	http.Handle("/errors", periskopHandler)
	http.HandleFunc("/-/health", healthHandler)

	address := fmt.Sprintf("0.0.0.0:%s", *port)
	log.Printf("Serving on address %s", address)
	log.Fatal(http.ListenAndServe(address, nil))
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("OK"))
	if err != nil {
		log.Fatalf("error running health handler")
	}
}
