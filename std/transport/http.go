package main

import (
	"fmt"
	"net/http"

	"github.com/hashicorp/go-hclog"
	"github.com/treethought/roc"
)

var log = hclog.New(&hclog.LoggerOptions{
	DisableTime: true,
})

type HttpTransport struct {
	*roc.TransportImpl
}

func (t HttpTransport) handler(w http.ResponseWriter, req *http.Request) {
	identifier := roc.Identifier(fmt.Sprintf("res:/%s", req.URL.Path))
	log.Info("mapped http request to identifier", "identifier", identifier)

	// TODO refactor to use roc.Source()?

	ctx := roc.NewRequestContext(identifier, roc.Source)
	// ctx.Scope = t.Scope

	resp, err := t.Dispatch(ctx)
	if err != nil {
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
