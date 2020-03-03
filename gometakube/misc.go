package gometakube

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func unexpectedResponseError(r *http.Response) error {
	reply := make(map[string]interface{})
	json.NewDecoder(r.Body).Decode(&reply)
	return fmt.Errorf("unexpected response %d: %v", r.StatusCode, reply)
}
