package roc

import (
	"io"
	"net/http"
	"strings"

	proto "github.com/treethought/roc/proto/v1"
)

// TODO drop this and
const EndpointTypeHTTPClient string = "httpclient"

type HTTPClient struct {
	*BaseEndpoint
	client *http.Client
}

func NewHTTPClient(ed *proto.EndpointMeta) HTTPClient {
	if ed == nil {
		ed = &proto.EndpointMeta{
			Type:       EndpointTypeHTTPClient,
			Identifier: "httpclient",
			Grammar: &proto.Grammar{
				Base: "active:http",
				Active: &proto.ActiveElement{
					Identifier: "active:http", // not using res:/ because that is added via overlay rewrite
					Arguments: []*proto.ActiveArgument{
						{Name: "url"},
					},
				},
			},
		}
	}

	return HTTPClient{
		BaseEndpoint: NewBaseEndpoint(ed),
		client:       http.DefaultClient,
	}
}

func (e HTTPClient) Source(ctx *RequestContext) interface{} {
	urlValue, err := ctx.GetArgumentValue("url")
	if err != nil {
		return "failed to get url value"
	}
	s := new(proto.String)
	err = urlValue.To(s)
	if err != nil {
		log.Error("failed to convert url value to string", "err", err)
		return err
	}
	// TODO https:// in later half of url being replacedwith https:/
	url := strings.Replace(s.GetValue(), "https:/", "https://", 1)

	log.Debug("making HTTP request", "url", url)
	resp, err := e.client.Get(url)
	if err != nil {
		log.Error("failed to make http request", "url", url, "err", err)
		return err
	}
	data, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		log.Error("failed to read http response body", "err", err)
		return err
	}
	return string(data)
}
