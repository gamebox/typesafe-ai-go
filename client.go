package typesafe

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"fmt"
	"net/http"
	"strings"
	"time"
)

var defaultUrl = `https://api.typesafe.ai/v1/systemone`
var baseTime = 100 * time.Millisecond
var ErrTooManyRequests = 429
var ErrOverloaded = 529

// This error is thrown when no questions have been supplied in a [Request]
type NoQuestionsErr struct{}

func (err NoQuestionsErr) Error() string {
	return "There must be at least one question asked"
}

type Client struct {
	c   *http.Client
	url string
	key string
}

// Creates a new Typesafe API [Client] using a provided valid API key
func NewClient(key string) *Client {
	c := Client{
		key: key,
		url: defaultUrl,
		c:   http.DefaultClient,
	}

	return &c
}

// Updates the [http.Client] used, the default comes from [http.DefaultClient]
func (c *Client) WithHttpClient(httpClient *http.Client) {
	c.c = httpClient
}

// Updates the url used to access the api, defaults to Typesafe AI's own API
func (c *Client) WithUrl(url string) {
	c.url = url
}

// Send a [Request] to the endpoint. Can return a [NoQuestionsError], or several other
// errors if the request fails for any reason.
//
// This method should handle redirects correctly, and should use exponential backoff
// when the endpoint returns a 429 or 529 status code.
//
// If successful, use the methods of the [Response] to read the answers to your questions
func (c Client) Ask(req Request) (response Response, err error) {
	if req.Model == "" {
		req.Model = ModelLatest
	}
	if len(req.Questions) == 0 {
		return response, NoQuestionsErr{}
	}

	sb := strings.Builder{}
	enc := jsontext.NewEncoder(&sb)

	err = json.MarshalEncode(enc, req)
	if err != nil {
		return response, fmt.Errorf("Could not encode request body: %e\n", err)
	}

	request, err := http.NewRequest("POST", defaultUrl, strings.NewReader(sb.String()))
	if err != nil {
		return response, fmt.Errorf("Could not create request: %e\n", err)
	}
	request.Header.Add("Authorization", "Bearer "+c.key)
	request.Header.Add("Content-Type", "application/json")

	resp, err := c.askWithBackoff(request)
	if err != nil {
		return response, fmt.Errorf("Failed to make request: %s\n", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		switch resp.StatusCode {
		case 422:
			var body any
			err = json.UnmarshalRead(resp.Body, &body)
			return response, fmt.Errorf("Bad Request. This should never happen:\n%#v", body)
		default:
			return response, fmt.Errorf("Had an error code of %d", resp.StatusCode)
		}
	}

	err = json.UnmarshalRead(resp.Body, &response)
	if err != nil {
		return response, fmt.Errorf("Failed to read the response body: %e\n", err)
	}

	return response, nil
}

func (c Client) askWithBackoff(req *http.Request) (resp *http.Response, err error) {
	attempt := 1
	for err == nil && attempt < 10 {
		resp, err = c.c.Do(req)
		if err != nil && resp.StatusCode != ErrOverloaded && resp.StatusCode != ErrTooManyRequests {
			// This is a conventional error
			return resp, err
		}
		if err != nil {
			// We need to backoff with this error code
			backoff := time.Duration(int64(2^attempt)*baseTime.Milliseconds()) * time.Millisecond
			time.Sleep(backoff)
			attempt += 1
			continue
		}
		// The request succeeded
		return resp, nil
	}

	// Too many attempts
	return resp, fmt.Errorf("Failed after sending 10 unhandled request to busy server")
}
