package circleci

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
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

// ListArtifacts returns all artifacts for given build
func (client *Client) ListArtifacts(query *ListArtifactInput) ([]*Artifact, error) {
	var (
		artifacts = []*Artifact{}
		path      = filepath.Join(
			"project",
			"github",
			query.Username,
			query.Project,
			query.Build,
			"artifacts")
	)

	response, err := client.do("GET", path, url.Values{}, nil)
	if err != nil {
		return artifacts, err
	}

	if err = client.unmarshal(response, &artifacts); err != nil {
		return artifacts, err
	}

	return artifacts, nil
}

// DownloadArtifacts downloads the artifacts
func (client *Client) DownloadArtifact(artifact *Artifact, dir string) error {
	response, err := client.do("GET", artifact.URL, url.Values{}, nil)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	path := filepath.Join(dir, artifact.Path)
	dir, _ = filepath.Split(path)

	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, response.Body)
	return err
}

// ListRecentBuilds returns all recent builds
func (client *Client) ListRecentBuilds() ([]*Build, error) {
	builds := []*Build{}

	response, err := client.do("GET", "recent-builds", url.Values{}, nil)
	if err != nil {
		return builds, err
	}

	if err = client.unmarshal(response, &builds); err != nil {
		return builds, err
	}

	return builds, nil
}

// SearchBuilds returns all recent builds
func (client *Client) SearchBuilds(query *SearchBuildInput) ([]*Build, error) {
	var (
		builds  = []*Build{}
		encoder = form.NewEncoder()
		path    = filepath.Join(
			"project",
			"github",
			query.Username,
			query.Project)
	)

	values, err := encoder.Encode(query)
	if err != nil {
		return builds, err
	}

	response, err := client.do("GET", path, values, nil)
	if err != nil {
		return builds, err
	}

	if err = client.unmarshal(response, &builds); err != nil {
		return builds, err
	}

	return builds, nil
}

func (client *Client) do(method, path string, values url.Values, data interface{}) (*http.Response, error) {
	values.Set("circle-token", client.Token)

	uri, err := url.Parse(path)
	if err == nil && uri.Scheme == "https" {
		uri.RawQuery = values.Encode()
	} else {
		uri = addr.ResolveReference(&url.URL{Path: path, RawQuery: values.Encode()})
	}

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
