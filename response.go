package roc

import "net/http"

// Verb specifies the acton to be taken by endpoint when resolving request
type Expiry int

const (

	// Source retrieves representation of resource
	EXPIRY_ALWAYS Expiry = iota
	EXPIRY_CONSTANT
	EXPIRY_NEVER
	EXPIRY_DEPENDENT
	EXPIRY_MIN_CONSTANT_DEPENDENT
	EXPIRY_MAX_CONSTANT_DEPENDENT
	EXPIRY_FUNCTION
	EXPIRY_MIN_FUNCTION_DEPENDENT
	EXPIRY_POLLED_DEPENDENT
)

// String - Creating common behavior - give the type a String function
func (e Expiry) String() string {
	return [...]string{
		"ALWAYS",
		"CONSTANT",
		"NEVER",
		"DEPENDENT",
		"MIN_CONSTANT_DEPENDENT",
		"MAX_CONSTANT_DEPENDENT",
		"FUNCTION",
		"MIN_FUNCTION_DEPENDENT",
		"POLLED_DEPENDENT",
	}[e-1]
}

func (e Expiry) EnumIndex() int {
	return int(e)
}

type Response struct {
	representation Representation
	metadata       map[string]interface{}
	mimetype       string
	expiry         Expiry
	expiryTime     uint64
	noCache        bool
	header         http.Header
    Affirmation bool
}


