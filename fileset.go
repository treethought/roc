package roc

import (
	"io/ioutil"
	"strings"

	proto "github.com/treethought/roc/proto/v1"
)

const EndpointTypeFileset string = "fileset"

type Fileset struct {
	*BaseEndpoint
	Mutable bool
}

func NewFilesetRegex(ed *proto.EndpointMeta) Fileset {
	return Fileset{
		BaseEndpoint: NewBaseEndpoint(ed),
		Mutable:      false,
	}
}

func (e Fileset) Source(ctx *RequestContext) interface{} {
	path := strings.Replace(ctx.Request().Identifier().String(), "res:/", "", 1)
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Error("failed to read fileset path", "path", path, "err", err)
		return err
	}
	return string(data)
}
