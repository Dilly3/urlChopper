package msgpack

import (
	"encoding/json"

	"github.com/dilly3/urlshortner/internal"
	"github.com/pkg/errors"
)

type Redirect struct{}

func (r *Redirect) Decode(input []byte) (*internal.Redirect, error) {
	redirect := &internal.Redirect{}
	if err := json.Unmarshal(input, redirect); err != nil {
		return nil, errors.Wrap(err, "serializer.msgpack.Decode")
	}
	return redirect, nil
}

func (r *Redirect) Encode(input *internal.Redirect) ([]byte, error) {
	rawMsg, err := json.Marshal(input)
	if err != nil {
		return nil, errors.Wrap(err, "serializer.msgpack.Encode")
	}
	return rawMsg, nil
}
