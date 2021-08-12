package roc

import (
	"github.com/treethought/roc/proto"
	protov2 "google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/structpb"
)

// Identifier is an opaque token that identifies a single resource
// a resource may have one or more identifiers
type Identifier struct {
	m *proto.String
}

func NewIdentifier(val string) Identifier {
	return Identifier{m: &proto.String{Value: val}}
}

func (i Identifier) ProtoReflect() protoreflect.Message {
	return i.m.ProtoReflect()
}
func (i Identifier) String() string {
	return i.m.GetValue()
}

// func (i Identifier) String() string {
// 	return i
// }

type RepresentationClass string

func (r RepresentationClass) String() string {
	return string(r)
}

// type RepresentationClass interface {
// 	String() string
// 	Identifier() Identifier
// }

type ComparibleRepresentation interface {
	Equals(interface{}) bool
	HashCode() int
}

type Representation struct {
	*proto.Representation
}

func NewRepresentation(val interface{}) Representation {
	log.Info("creating representation from", "val", val)

	var msg protoreflect.ProtoMessage

	switch v := val.(type) {
	case Representation:
		return v

	case protov2.Message:
		msg = v

	case string:
		msg = &proto.String{Value: v}

	case nil:
		msg = &proto.Empty{}

	default:
		sval, err := structpb.NewValue(val)
		if err != nil {
			log.Error("failed to convert representation to proto struct")
		}
		msg = sval

	}

	any, err := anypb.New(msg)
	if err != nil {
		log.Error("failed to construct any from value", "err", err)
		panic(err)

	}

	return Representation{Representation: &proto.Representation{Value: any}}

}

// type Representation interface {
// 	protov2.Message
// }
