// Command client_http_go is a typed HTTP client for the paired user API in
// ../server_http_go, built with the standard library only.
//
// It demonstrates:
//   - A reusable HttpClient wrapping *http.Client with a base URL and bearer key.
//   - Context-aware requests (http.NewRequestWithContext) with per-call timeouts
//     inherited from the shared client.
//   - A fan-out worker pool: 50 users are pushed onto a job channel and POSTed
//     concurrently by a fixed set of workers, with results collected on resCh.
//   - Sentinel-error handling via errors.Is (ErrNotFound) and %w wrapping so
//     callers can distinguish a 404 from a transport error.
//
// Run the server first, then this client:
//
//	go run ../server_http_go   # terminal 1
//	go run .                   # terminal 2
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"
)

// type Config struct {
// 	BaseUrl string
// 	ApiKey  string
// 	TimeOut time.Duration
// }

// type HttpClient struct {
// 	client *http.Client
// 	config Config
// }

// User mirrors the server's User resource shape for request/response JSON.
type User struct {
	Id      string    `json:"id"`
	Name    string    `json:"name"`
	Email   string    `json:"email"`
	Created time.Time `json:"created"`
}

const (
	baseUrl = "http://localhost:8080" // where the paired server listens
	apiKey  = "abcxyz"                // bearer token expected by the server's auth middleware
)

// HttpClient is a thin wrapper over *http.Client that carries the base URL and
// API key so callers don't repeat them on every request.
type HttpClient struct {
	client  *http.Client
	baseUrl string
	apiKey  string
}

// NewHttpClient returns an HttpClient with a 20s request timeout and the
// package-level baseUrl/apiKey defaults.
func NewHttpClient() *HttpClient {
	return &HttpClient{
		client:  &http.Client{Timeout: 20 * time.Second},
		baseUrl: baseUrl,
		apiKey:  apiKey,
	}

}

// Result pairs a user id with the error (if any) from its POST, so the worker
// pool can report per-user outcomes on a single channel.
type Result struct {
	Id  string
	Err error
}

// main generates 50 users, POSTs them concurrently through a worker pool, then
// demonstrates a single get and a list against the server.
func main() {

	users := make([]User, 0, 50)
	for i := 1; i <= 50; i++ {
		users = append(users, User{
			Id:    fmt.Sprintf("user-%d", i),
			Name:  fmt.Sprintf("Test User %d", i),
			Email: fmt.Sprintf("user%d@example.com", i),
		})
	}

	wg := sync.WaitGroup{}
	ctx := context.Background()
	c := NewHttpClient()

	// Buffered channels sized to the workload so producers never block.
	jobCh := make(chan User, len(users))
	resCh := make(chan Result, len(users))
	workers := 5

	// Enqueue every user, then close so workers exit their range loop when drained.
	for _, u := range users {
		jobCh <- u
	}
	close(jobCh)

	// Fan out: each worker pulls jobs until jobCh is closed and empty.
	for w := 1; w <= workers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for job := range jobCh {
				err := c.postUser(ctx, job)
				resCh <- Result{
					Id:  job.Id,
					Err: err,
				}

			}
		}()

	}

	// Close resCh once all workers finish so the collector loop below terminates.
	go func() {
		wg.Wait()
		close(resCh)
	}()

	// Fan in: drain results as they arrive.
	for r := range resCh {
		if r.Err != nil {
			log.Printf("post user error: %v\n", r.Err)
			continue
		}
		fmt.Printf("user created successfully:%s\n", r.Id)
	}

	// u := User{
	// 	Id:    "abc123",
	// 	Name:  "ashwin",
	// 	Email: "ashwin@example.com",
	// }

	// err := c.postUser(ctx, u)
	// if err != nil {
	// 	log.Fatalf("post user %s: %v", u.Id, err)
	// }

	userId := "user-1"
	userData, err := c.getUser(ctx, userId)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			fmt.Printf("user %s not found", userId)
			return
		}
		log.Fatalf("get user: %v", err)
	}

	fmt.Printf("user data: %s - value: %v\n", userId, userData)

	data, err := c.listUser(ctx, 10)
	if err != nil {
		log.Fatalf("list users: %v", err)
	}

	for _, u := range data {
		fmt.Printf("user data - id:%s, name:%s, email:%s, created:%v\n", u.Id, u.Name, u.Email, u.Created)
	}

}

// ErrNotFound is a sentinel returned by getUser when the server responds 404,
// letting callers branch with errors.Is instead of inspecting status codes.
var ErrNotFound = errors.New("user not found")

// getUser issues GET /users/{id}. It returns ErrNotFound (wrapped) on a 404,
// a wrapped error for any other non-200 status or transport/decoding failure,
// and the decoded User on success.
func (c *HttpClient) getUser(ctx context.Context, id string) (User, error) {
	url := fmt.Sprintf("%s/%s/%s", c.baseUrl, "users", id)
	token := fmt.Sprintf("Bearer %s", c.apiKey)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return User{}, fmt.Errorf("new request error %w", err)
	}

	req.Header.Set("Authorization", token)

	resp, err := c.client.Do(req)
	if err != nil {
		return User{}, fmt.Errorf("get user error in url %s:%w", url, err)
	}

	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {

		return User{}, fmt.Errorf("get user %s: %w", id, ErrNotFound)
	}

	if resp.StatusCode != http.StatusOK {
		return User{}, fmt.Errorf("server error in %s with statuscode %d", url, resp.StatusCode)
	}
	var u User
	if err := json.NewDecoder(resp.Body).Decode(&u); err != nil {
		return User{}, fmt.Errorf("invalud data: %w", err)
	}

	return u, nil

}

// listUser issues GET /users?limit=N and returns the decoded slice of users,
// or a wrapped error on a non-200 status or transport/decoding failure.
func (c *HttpClient) listUser(ctx context.Context, limit int) ([]User, error) {
	q := url.Values{}
	q.Set("limit", strconv.Itoa(limit))

	url := fmt.Sprintf("%s/users?%s", c.baseUrl, q.Encode())
	token := fmt.Sprintf("Bearer %s", c.apiKey)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("new request error %w", err)
	}

	req.Header.Set("Authorization", token)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("list user error url %s:%w", url, err)
	}

	defer resp.Body.Close()
	// if resp.StatusCode == http.StatusNotFound {

	// 	return nil, fmt.Errorf("no users found: %w", ErrNotFound)
	// }

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server error in %s with statuscode %d", url, resp.StatusCode)
	}
	var u []User
	if err := json.NewDecoder(resp.Body).Decode(&u); err != nil {
		return nil, fmt.Errorf("invalud data: %w", err)
	}

	return u, nil

}

// postUser issues POST /users with the JSON-encoded user. It returns nil on a
// 201 Created and a wrapped error otherwise (marshal, transport, or unexpected
// status).
func (c *HttpClient) postUser(ctx context.Context, u User) error {

	url := fmt.Sprintf("%s/%s", c.baseUrl, "users")
	token := fmt.Sprintf("Bearer %s", c.apiKey)

	body, err := json.Marshal(u)
	if err != nil {
		return fmt.Errorf("marshal user %s: %w", u.Id, err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("new request error %w", err)
	}

	req.Header.Set("Authorization", token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("post user error for id %s :%w", u.Id, err)
	}

	defer resp.Body.Close()
	// if resp.StatusCode == http.StatusNotFound {

	// 	return fmt.Errorf("get user %s: %w", id, ErrNotFound)
	// }

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("post user %s with statuscode %d", u.Id, resp.StatusCode)
	}

	return nil

}

// func LoadConfig() Config {
// 	return Config{
// 		BaseUrl: getEnv("BASE_URL", "http://localhost:8080/"),
// 		ApiKey:  getEnv("API_KEY", "abcxyz"),
// 		TimeOut: 20 * time.Second,
// 	}
// }

// func getEnv(conf, fallback string) string {
// 	if v := os.Getenv(conf); v != "" {
// 		return v
// 	}
// 	return fallback

// }

// func NewHttpClient(conf Config) *HttpClient {
// 	return &HttpClient{
// 		client: &http.Client{Timeout: conf.TimeOut},
// 		config: conf,
// 	}
// }
