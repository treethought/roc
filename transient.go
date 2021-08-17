package roc

import (
	"fmt"

	"github.com/google/uuid"
	proto "github.com/treethought/roc/proto/v1"
)

const EndpointTypeTransient string = "transient"

// TransientEndpoint is dynamically generated in-memory endpoint
// these are typically used for internal temporary resources.
type TransientEndpoint struct {
	*BaseEndpoint
}

func NewTransientEndpoint(ed *proto.EndpointMeta) TransientEndpoint {
	ed.Type = EndpointTypeTransient

	if ed.GetIdentifier() == "" {
		uid := uuid.New()
		uri := fmt.Sprintf("transient://%s", uid.String())
		ed.Identifier = uri
	}

	if ed.GetLiteral() == nil {
		log.Warn("transient endpoint has nil represetation", "id", ed.Identifier)
	}

	if ed.GetGrammar() == nil {
		ed.Grammar = &proto.Grammar{Base: ed.Identifier}
	}

	repr := Representation{ed.Literal}
	log.Debug("creating transient endpoint",
		"type", repr.Type(),
		"uri", ed.Identifier,
	)
	log.Trace(repr.String())

	return TransientEndpoint{
		BaseEndpoint: NewBaseEndpoint(ed),
	}
}

func (e TransientEndpoint) Type() string {
	return EndpointTypeTransient
}

func (e TransientEndpoint) Source(ctx *RequestContext) interface{} {
	log.Debug("sourcing transient endpoint",
		"identifier", ctx.Request().Identifier(),
		"type", e.Type(),
	)

	rep := NewRepresentation(e.Meta().GetLiteral())
	log.Trace(rep.String())

	m, err := rep.ToMessage()
	if err != nil {
		log.Error("failed to construct transient message", "err", err)
		return err
	}

	return m
}
