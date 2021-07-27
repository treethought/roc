package roc


type Request struct {
	identifier          Identifier
	verb                Verb
	representationClass RepresentationClass
}


func (r Request) Identifier() Identifier {
    return r.identifier
}

// RequestArment are tha arguments of the request mapped from the parsed grammar roups
type RequestArgument struct {
	// the name of the request argument
	Name string

	// the value of the named argument. can be passed by value or
	Value string
}

func (r Request) Arguments() {

}
