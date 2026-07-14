package testutil

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/mic615/chill-crate-api/internal/storage"
)

func SetupTestStorageClient(t *testing.T) *storage.Storage {
	t.Helper()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	t.Cleanup(server.Close)

	client := s3.NewFromConfig(aws.Config{
		Region: "us-west-1",
		Credentials: aws.NewCredentialsCache(
			credentials.NewStaticCredentialsProvider("test", "test", ""),
		),
	}, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(server.URL)
		o.UsePathStyle = true
		o.HTTPClient = server.Client()
	})

	return &storage.Storage{Client: client}
}
