package main

import (
	"fmt"
	"net/http"

	"github.com/hashicorp/go-hclog"
	"github.com/treethought/roc"
	"github.com/treethought/roc/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

var log = hclog.Default()

type HttpTransport struct {
	*roc.TransportImpl
}

func getAsHttpRequest(rep roc.Representation) (*proto.HttpRequest, error) {
	// TODO: this needs to be handled with response vlass behnd the scenes
	any, err := anypb.New(rep)
	if err != nil {
		return nil, err
	}

	m := new(proto.HttpRequest)
	err = any.UnmarshalTo(m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (t HttpTransport) handler(w http.ResponseWriter, req *http.Request) {
	identifier := roc.NewIdentifier(fmt.Sprintf("http:/%s", req.URL.String()))
	log.Info("transport received request", "identifier", identifier, "url", req.URL.String())

	ctx := roc.NewRequestContext(identifier, proto.Verb_Source)

	// create dynamic httpRequest accessor that will provide acess to http request
	// rep := roc.NewHttpRequestDefinition(req)
	rep := roc.NewHttpRequestMessage(req)

	repr := roc.NewRepresentation(rep)

	log.Info("setting httpRequest request arg value")
	ctx.Request().SetArgumentByValue("httpRequest", repr)
	// ctx.Request.SetArgumentByValue("httpResponse", w)

	resp, err := t.Dispatch(ctx)
	if err != nil {
		log.Error("failed to dispatch request", "err", err)
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
	}

	m := new(proto.String)
	err = resp.MarshalTo(m)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	log.Info("returning response to http client", "response", m)
	w.Write([]byte(m.GetValue()))

}

func main() {

	transport := HttpTransport{
		TransportImpl: roc.NewTransport("http_transport"),
	}

	go func() {
		log.Info("starting transport rpc server")
		roc.ServeTransport(&transport)
	}()

	log.Info("starting transport http server")
	http.HandleFunc("/", transport.handler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}

}
