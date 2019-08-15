package router

import (
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/martinomburajr/gcp-vision-api/app"
	"google.golang.org/api/option"
	"net/http"
	"os"
)

// GCPVisionAPIServer is the base server that handles all incoming requests.
// It is equipped with a router and a reference to Vision API
type GCPVisionAPIServer struct {
	Router *mux.Router
	StorageClient *storage.Client
}

// Init ensures all relevant clients are called and set up before requests can begin
func (g *GCPVisionAPIServer) Init() error {
	clientOption := option.WithCredentialsFile(app.CredentialsLocalPath)
	ctx := context.Background()

	storageClient, err := storage.NewClient(ctx, clientOption)
	if err != nil {
		return fmt.Errorf("failed to initialize cloud storage | %v",
			err)
	}
	g.StorageClient = storageClient

	return nil
}

// Routes returns a Handler that acts as a multiplexer. We use the gorilla.Mux router.
func (g *GCPVisionAPIServer) Routes() *mux.Router {
	router := mux.NewRouter()

	router.Methods(http.MethodGet).Path("/").
		Handler(AppHandler(HealthCheckHandler()))
	router.Methods(http.MethodPost).Path("/vision/ocr").
		Handler(AppHandler(OCRHandler(g)))
	router.Methods(http.MethodPost).Path("/vision/ocr/bucket").
		Handler(AppHandler(OCRBucketDirHandler(g)))
	// Logger
	router.Handle("/", handlers.CombinedLoggingHandler(os.Stderr, router))
	return router
}

// Initializable interface houses the Init method used for initializing types.
type Initializable interface {
	// Init is used to initialize a given server with clients, environment variables and any other logic
	Init() error
}
