package circleci

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path/filepath"

	"github.com/go-playground/form"
)

var (
	addr = &url.URL{Host: "circleci.com", Scheme: "https", Path: "/api/v1.1/"}
)

// Client for CircleCI API
type Client struct {
	Token string
}

// ListRecentBuilds returns all recent builds
func (client *Client) ListRecentBuilds() ([]*Build, error) {
	build := []*Build{}

	response, err := client.do("GET", "recent-builds", url.Values{}, nil)
	if err != nil {
		return build, err
	}

	if err = client.unmarshal(response, &build); err != nil {
		return build, err
	}

	return build, nil
}

// SearchBuilds returns all recent builds
func (client *Client) SearchBuilds(query *Query) ([]*Build, error) {
	var (
		build   = []*Build{}
		path    = filepath.Join("project", "github", query.Username, query.Project)
		encoder = form.NewEncoder()
	)

	values, err := encoder.Encode(query)
	if err != nil {
		return build, err
	}

	response, err := client.do("GET", path, values, nil)
	if err != nil {
		return build, err
	}

	if err = client.unmarshal(response, &build); err != nil {
		return build, err
	}

	return build, nil
}

func (client *Client) do(method, path string, values url.Values, data interface{}) (*http.Response, error) {
	values.Set("circle-token", client.Token)

	uri := addr.ResolveReference(&url.URL{Path: path, RawQuery: values.Encode()})
	body := &bytes.Buffer{}

	if data != nil {
		if err := json.NewEncoder(body).Encode(data); err != nil {
			return nil, err
		}
	}

	request, err := http.NewRequest(method, uri.String(), body)
	if err != nil {
		return nil, err
	}

	request.Header.Add("Accept", "application/json")
	request.Header.Add("Content-Type", "application/json")

	return http.DefaultClient.Do(request)
}

func (client *Client) unmarshal(response *http.Response, obj interface{}) error {
	defer response.Body.Close()

	if response.StatusCode >= 300 {
		return client.error(response)
	}

	if obj != nil {
		return json.NewDecoder(response.Body).Decode(obj)
	}

	return nil
}

func (client *Client) error(response *http.Response) error {
	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return &Error{
			HTTPStatusCode: response.StatusCode,
			Message:        fmt.Sprintf("unable to parse response: %s", err),
		}
	}

	if len(body) > 0 {
		message := struct {
			Message string `json:"message"`
		}{}

		if err = json.Unmarshal(body, &message); err != nil {
			return &Error{
				HTTPStatusCode: response.StatusCode,
				Message:        fmt.Sprintf("unable to parse API response: %s", err),
			}
		}

		return &Error{
			HTTPStatusCode: response.StatusCode,
			Message:        message.Message,
		}
	}

	return &Error{HTTPStatusCode: response.StatusCode}
}
