package roc

import (
	"reflect"

	proto "github.com/treethought/roc/proto/v1"
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

type RepresentationClass string

func (r RepresentationClass) String() string {
	return string(r)
}

type ComparibleRepresentation interface {
	Equals(interface{}) bool
	HashCode() int
}

type Representation struct {
	m *proto.Representation
}

func NewRepresentation(val interface{}) Representation {
	log.Trace("creating representation from", "type", reflect.TypeOf(val))

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
		m, err := v.ToMessage()
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
		log.Warn("creating unknown representation as struct")
		sval, err := structpb.NewValue(val)
		if err != nil {
			log.Error("failed to convert representation to proto struct", "err", err)
		}
		msg = sval

	}

	any, err := anypb.New(msg)
	if err != nil {
		log.Error("failed to construct any from value", "err", err)
		panic(err)

	}

	rep := Representation{m: &proto.Representation{Value: any}}

	log.Debug("created representation",
		"from", reflect.TypeOf(val),
		"type", rep.Type(),
	)

	return rep
}

func (r Representation) ProtoReflect() protoreflect.Message {
	return r.m.ProtoReflect()
}

// messag provides the wrapped Representation proto message
func (r Representation) message() *proto.Representation {
	return r.m
}

// String returns the deserialized representation of the underlying type
// if the underlying type has a custom JSON representation, that will be returned instead
func (r Representation) String() string {
	return r.any().String()
}

// any returns the underlying Any proto message
func (r *Representation) any() *anypb.Any {
	return r.m.GetValue()
}

// Is reports whether the representation is the same type as m
func (r *Representation) Is(m protov2.Message) bool {
	return r.any().MessageIs(m)
}

// ToMessage unmarshals the representation into a new message of the underlying type
func (r *Representation) ToMessage() (protov2.Message, error) {
	return r.any().UnmarshalNew()
}

// As unnmarshals the representation into m
func (r Representation) To(m protov2.Message) error {
	return r.any().UnmarshalTo(m)
}

// Type returns the name of underlying proto message
func (r Representation) Type() protoreflect.FullName {
	return r.any().MessageName()
}

// URL returns the URL that identifies the underlying type
func (r Representation) URL() string {
	return r.any().TypeUrl
}
