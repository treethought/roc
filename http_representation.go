package roc

import (
	"encoding/json"
	"net/http"
)

type HttpRequestRepresentation struct {
	request *http.Request
}

func NewHttpRequestRepresentation(req *http.Request) *HttpRequestRepresentation {
	return &HttpRequestRepresentation{
		request: req,
	}
}

func (r *HttpRequestRepresentation) MarshalJSON() ([]byte, error) {
	type Alias http.Request
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(r.request),
	})

}
func (r *HttpRequestRepresentation) UnmarshalJSON(data []byte) error {
	type Alias http.Request
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(r.request),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	// u.LastSeen = time.Unix(aux.LastSeen, 0)
	return nil
}
