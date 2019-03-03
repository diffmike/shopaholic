package cmd

import (
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
)

// All commands should implement this interfaces
type Commander interface {
	Execute(args []string) error
}

// responseError returns error with status and response body
func responseError(resp *http.Response) error {
	body, e := ioutil.ReadAll(resp.Body)
	if e != nil {
		body = []byte("")
	}
	return errors.Errorf("error response %q, %s", resp.Status, body)
}
