package roc

import (
	"github.com/treethought/roc/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

// Request represents request for a resouce
type Request struct {
	m *proto.Request
}

func NewRequest(i Identifier, verb proto.Verb, class RepresentationClass) *Request {
	classStr := ""
	if class != "" {
		classStr = class.String()
	}

	return &Request{
		m: &proto.Request{
			Identifier:          i.String(),
			Verb:                verb,
			RepresentationClass: classStr,
			Arguments:           make(map[string]*proto.StringSlice),
			ArgumentValues:      make(map[string]*anypb.Any),
		},
	}
}

func (r *Request) Identifier() Identifier {
	return NewIdentifier(r.m.Identifier)
}

// SetRepresentationClass sets the desired format of the representation response
func (r *Request) SetRepresentationClass(class string) {
	r.m.RepresentationClass = class
}

// SetArgument sets the value of an argument to an identifier
// The argument's representation can then be sources during evalutation
// This replaces any existing values, to append an identifier, use AddArgument
func (r *Request) SetArgument(name string, i Identifier) {
	r.m.Arguments[name] = &proto.StringSlice{
		Values: []string{i.String()},
	}
}

// AddArgument appends an identifier to any existing ones for the named arguement
func (r *Request) AddArgument(name string, i Identifier) {
	_, exists := r.m.Arguments[name]
	if !exists {
		r.m.Arguments[name] = &proto.StringSlice{}
	}
	r.m.Arguments[name].Values = append(r.m.Arguments[name].Values, i.String())
}

func (r *Request) SetArgumentByValue(name string, val Representation) {
	log.Info("setting argument value", "arg", name, "type", val.Type())
	r.m.ArgumentValues[name] = val.any()

}
