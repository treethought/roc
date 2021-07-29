package roc

import (
	"fmt"
	"log"
	"net/http"
)

type Resolver interface {
	Resolve(request *Request, ch chan (Evaluator))
}

type Evaluator interface {
	Evaluate(request *Request) Representation
}

type Dispatcher struct {
	resolvers   []Resolver
	evalutators []Evaluator
}

type Kernel struct {
	Spaces   []*Space
	receiver chan (*Request)
	server   http.Server
}

func NewKernel() *Kernel {
	return &Kernel{
		Spaces:   []*Space{},
		receiver: make(chan *Request),
	}
}

func (k *Kernel) Serve(port int) {
	http.HandleFunc("/issue", k.HandleIssue)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))

}

func (k *Kernel) HandleIssue(w http.ResponseWriter, req *http.Request) {
	req.Context()

	identifier := req.URL.Query().Get("identifier")
	if identifier == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

}

func (k *Kernel) Receiver() chan (*Request) {
	return k.receiver
}

func (k *Kernel) Register(s *Space) {
	s.channel = k.receiver
	k.Spaces = append(k.Spaces, s)
}

// func (k *Kernel) LoadFromFile(path string) error {

// 	data, err := ioutil.ReadFile(path)
// 	if err != nil {
// 		return err
// 	}

// 	spaces := []*Space{}

// 	err = yaml.Unmarshal(data, &spaces)
// 	if err != nil {
// 		return err
// 	}

// 	k.Spaces = append(k.Spaces, spaces...)
// 	return nil
// }

func (k Kernel) startReceiver() {
	for {
		incoming := <-k.receiver
		k.Dispatch(incoming)
	}
}

func (k Kernel) buildResolveRequest(request *Request) *Request {
	return NewRequest(request.Identifier(), Resolve, nil)

}

func (k Kernel) resolveEndpoint(request *Request) Endpoint {
	c := make(chan (Endpoint))
	for _, s := range k.Spaces {
		go s.Resolve(request, c)
	}

	return <-c
}

func (k Kernel) Dispatch(request *Request) Representation {
	log.Printf("dispatching request for identifer: %s", request.Identifier())

	endpoint := k.resolveEndpoint(request)
	log.Printf("resolved to endpoint: %s", endpoint)

	// TODO route verbs to methods
	rep := endpoint.Source(request)
	return rep

}
