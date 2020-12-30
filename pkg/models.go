package resource

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

func NewConcourse(
	stdin io.Reader,
	stderr io.Writer,
	stdout io.Writer,
	args []string,
) *Concourse {
	return &Concourse{
		stdin:  stdin,
		stderr: stderr,
		stdout: stdout,
		args:   args,
	}
}

type Concourse struct {
	stdin  io.Reader
	stderr io.Writer
	stdout io.Writer
	args   []string
}

type Response struct {
	Version Version `json:"version"`
}

type Version struct {
	Version string `json:"version"`
}

type Source struct {
	DefectDojoUrl string `json:"defectdojo_url"`
	Username      string `json:"username"`
	Password      string `json:"password,omitempty"`
	ApiKey        string `json:"api_key"`
	Debug         bool   `json:"debug"`
}

func (c *Concourse) WriteResponseToConcourse(response Response) error {
	return json.NewEncoder(c.stdout).Encode(response)
}

func (s *Source) Validate() error {
	if s.DefectDojoUrl == "" {
		return fmt.Errorf("Required `defectdojo_url` not supplied.")
	}
	if !strings.HasPrefix(s.DefectDojoUrl, "http://") && !strings.HasPrefix(s.DefectDojoUrl, "https://") {
		return fmt.Errorf("Please provide http(s):// prefix in `defectdojo_url`.")
	}
	if s.Username == "" {
		return fmt.Errorf("Required `username` not supplied.")
	}
	if s.ApiKey == "" {
		return fmt.Errorf("Required `api_key` not supplied.")
	}
	return nil
}
