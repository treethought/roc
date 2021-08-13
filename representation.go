package roc

import (
	"reflect"

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
	log.Debug("creating representation from", "type", reflect.TypeOf(val))

	var msg protoreflect.ProtoMessage

	switch v := val.(type) {

	// unmarshal into the underlying remote type
	case *anypb.Any:
		m, err := v.UnmarshalNew()
		if err != nil {
			log.Error("failed ot unmarshal new concreate any", "err", err)
			panic(err)
		}
		msg = m

		// just marshal the underlaying any into it's remote type
	case *proto.Representation:
		m, err := v.Value.UnmarshalNew()
		if err != nil {
			log.Error("failed ot unmarshal new concreate any", "err", err)
			panic(err)
		}
		// return NewRepresentation(m)
		msg = m

		// unwrap to any and unmarshal to remote type, don't want nested Representations
	case Representation:
		log.Debug("representation is already representation, unmarshalling")
		m, err := v.Any().UnmarshalNew()
		if err != nil {
			log.Error("failed ot unmarshal new conreate any", "err", err)
			panic(err)
		}
		msg = m

		return v

		// already a message so we can create Any directly from it
	case protov2.Message:
		msg = v

		// convert to proto msg
	case string:
		msg = &proto.String{Value: v}

		// TODO:
	case nil:
		msg = &proto.Empty{}

		// allow to return errors from endpoints
	case error:
		msg = &proto.ErrorMessage{Message: v.Error()}

		// not sure what we have, create a new Struct proto msg
	default:
		log.Warn("creating representation with struct")
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

	log.Info("created representation",
		"from_type", reflect.TypeOf(val),
		"any_url", any.TypeUrl,
	)

	return Representation{Representation: &proto.Representation{Value: any}}
}

func (r *Representation) Any() *anypb.Any {
	return r.Representation.GetValue()
}

func (r *Representation) Is(m protov2.Message) bool {
	return r.Any().MessageIs(m)
}

func (r *Representation) ToMessage() (protov2.Message, error) {
	return r.Any().UnmarshalNew()
}

func (r Representation) MarshalTo(m protov2.Message) error {
	return r.Any().UnmarshalTo(m)
}
func (r Representation) Name() protoreflect.FullName {
	return r.Any().MessageName()
}
func (r Representation) Type() string {
	return r.Any().TypeUrl
}

// type Representation interface {
// 	protov2.Message
// }
