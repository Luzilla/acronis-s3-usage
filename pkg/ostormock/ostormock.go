package ostormock

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Luzilla/acronis-s3-usage/pkg/ostor"
	"github.com/gorilla/mux"
)

func bucketsHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ostor.OstorBucketListResponse{
		Buckets: []ostor.Bucket{
			{
				Name: "bucket.example.org",
			},
			{
				Name: "another",
			},
		},
	})
}

// PUT Handler for /?ostor-users&emailAddress=<value>
func createUserHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ostor.OstorCreateUserResponse{
		Email: "user@example.org",
		ID:    "hash-1",
		AccessKeys: []ostor.AccessKeyPair{
			{
				AccessKeyID:     "key-1",
				SecretAccessKey: "secret-1",
			},
		},
	})
}

func disableEnableUserHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func genKeyHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ostor.OstorCreateUserResponse{
		Email: "user@example.org",
		ID:    "806e7d49f2dd9763",
		AccessKeys: []ostor.AccessKeyPair{
			{
				AccessKeyID:     "key-1",
				SecretAccessKey: "secret-1",
			},
			{
				AccessKeyID:     "key-2",
				SecretAccessKey: "secret-2",
			},
		},
	})
}

// GET Handler for /?ostor-users
func getUserHandler(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("emailAddress")
	if email != "" {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(ostor.OstorUser{
			Email:        "user@example.org",
			ID:           "806e7d49f2dd9763",
			State:        "enabled",
			Owner:        "0000000000000000",
			AccountCount: "0",
			Flags:        []string{},
			AccessKeys: []ostor.AccessKeyPair{
				{
					AccessKeyID:     "key-1",
					SecretAccessKey: "secret-1",
				},
			},
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ostor.OstorUsersListResponse{
		Users: []ostor.OstorUser{
			{
				Email:        "user@example.org",
				ID:           "806e7d49f2dd9763",
				State:        "enabled",
				Owner:        "0000000000000000",
				AccountCount: "0",
				Flags:        []string{},
				AccessKeys:   []ostor.AccessKeyPair{},
			},
		},
	})
}

// StartMockServer starts a mock HTTP server on a random port and returns the server and its URL
func StartMockServer(t *testing.T) (*httptest.Server, string) {
	router := mux.NewRouter()
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			next.ServeHTTP(w, r)
		})
	})
	router.HandleFunc("/", bucketsHandler).Methods(http.MethodGet).Queries(
		"ostor-buckets", "",
	)
	router.HandleFunc("/", getUserHandler).Methods(http.MethodGet).Queries(
		"ostor-users", "",
	)
	router.HandleFunc("/", createUserHandler).Methods(http.MethodPut).Queries(
		"ostor-users", "",
		"emailAddress", "{email}",
	)
	router.HandleFunc("/", disableEnableUserHandler).Methods(http.MethodPost).Queries(
		"ostor-users", "",
		"emailAddress", "{email}",
		"disable", "",
	)
	router.HandleFunc("/", disableEnableUserHandler).Methods(http.MethodPost).Queries(
		"ostor-users", "",
		"emailAddress", "{email}",
		"enable", "",
	)
	router.HandleFunc("/", genKeyHandler).Methods(http.MethodPost).Queries(
		"ostor-users", "",
		"emailAddress", "{email}",
		"genKey", "",
	)

	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Logf("Route not found: %s (method: %s)", r.URL.String(), r.Method)
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte("NO, NO, NO!"))
	})

	// Create a new test server
	server := httptest.NewServer(router)

	// handle clean-up in the library after each test
	t.Cleanup(func() {
		server.CloseClientConnections()
		server.Close()
	})

	// Return the server and its URL
	return server, server.URL
}
