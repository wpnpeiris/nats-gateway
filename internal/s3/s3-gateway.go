package s3

import (
	"log"

	"github.com/gorilla/mux"

	s3api "github.com/seaweedfs/seaweedfs/weed/s3api"
	s3cont "github.com/seaweedfs/seaweedfs/weed/s3api/s3_constants"

	"github.com/wpnpeiris/nats-gateway/internal/client"
)

// Provide S3 API implementation
type S3Gateway struct {
	*client.Client
}

func NewS3Gateway(natsServers string) (s3Gateway *S3Gateway) {
	comp := client.NewClient("s3-gateway")

	err := comp.SetupConnectionToNATS(natsServers)
	if err != nil {
		log.Fatal(err)
	}

	return &S3Gateway{
		comp,
	}
}

func (s3 S3Gateway) RegisterS3Routes(router *mux.Router) {

	s3Router := router.PathPrefix("/").Subrouter()

	iam := &s3api.IdentityAccessManagement{}

	s3Router.Methods("OPTIONS").HandlerFunc(s3.SetOptionHeaders)

	s3Router.Methods("HEAD").Path("/{bucket}/{key}").HandlerFunc(iam.Auth(s3.HeadObject, s3cont.ACTION_LIST))

	s3Router.Methods("GET").Path("/").HandlerFunc(s3.ListBuckets)
	s3Router.Methods("GET").Path("/{bucket}").HandlerFunc(s3.ListObjects)
	s3Router.Methods("GET").Path("/{bucket}/{key}").HandlerFunc(s3.Download)

	s3Router.Methods("PUT").Path("/{bucket}/{key}").HandlerFunc(s3.Upload)
}
