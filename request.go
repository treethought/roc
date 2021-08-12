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
	// pRep, err := repToProto(val)
	// if err != nil {
	// 	panic(err)
	// }

	// msg, ok := val.ProtoReflect().(protoreflect.ProtoMessage)
	// if !ok {
	// 	log.Error("rep value is not protoflect message")
	// }

	// m, err := anypb.New(msg)
	// if err != nil {
	// 	log.Error("failed to convert agument value to any", "err", err)
	// 	panic(err)
	// }

	log.Warn("setting argument value", "arg", name, "val_type", val.Any().TypeUrl)
	r.m.ArgumentValues[name] = val.Representation.Value

}

// // Identifier returns the identifier of the requested resource
// func (r Request) Identifier() Identifier {
// 	return r.identifier
// }

// // Verb returns the specified action to be taken when evaluating the request
// func (r Request) Verb() Verb {
// 	return r.verb
// }

// func (r *Request) Headers() http.Header {
// 	return r.headers
// }

// func (r Request) Arguments() url.Values {
// 	return r.argmuments
// }
