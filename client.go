package typesafe

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"fmt"
	"net/http"
	"strings"
	"time"
)

var url = `https://api.typesafe.ai/v1/systemone`
var baseTime = 100 * time.Millisecond
var ErrTooManyRequests = 429
var ErrOverloaded = 529

type NoQuestionsErr struct{}

func (err NoQuestionsErr) Error() string {
	return "There must be at least one question asked"
}

type Client struct {
	c   *http.Client
	key string
}

func NewClient(key string) *Client {
	c := Client{
		key: key,
		c:   http.DefaultClient,
	}

	return &c
}

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

	request, err := http.NewRequest("POST", url, strings.NewReader(sb.String()))
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
