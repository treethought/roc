package main

import (
	"fmt"
	"net/http"

	"github.com/hashicorp/go-hclog"
	"github.com/treethought/roc"
)

var log = hclog.Default()

type HttpTransport struct {
	*roc.TransportImpl
}

func (t HttpTransport) handler(w http.ResponseWriter, req *http.Request) {
	identifier := roc.Identifier(fmt.Sprintf("http:/%s", req.URL.String()))
	log.Info("transport received request", "identifier", identifier, "url", req.URL.String())

	// TODO refactor to use roc.Source()?

	ctx := roc.NewRequestContext(identifier, roc.Source)
	ctx.Request.SetArgumentByValue("httpRequest", req)
	ctx.Request.SetArgumentByValue("httpResponse", w)

	resp, err := t.Dispatch(ctx)
	if err != nil {
		log.Error("failed to dispatch request", "err", err)
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
	}

	r := fmt.Sprintf("%+v", resp)
	log.Info("returning response to http client", "response", r)
	w.Write([]byte(r))

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
