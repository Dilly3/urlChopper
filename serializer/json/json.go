package json

import (
	"encoding/json"

	"github.com/dilly3/urlshortner/internal"
	"github.com/pkg/errors"
)

type Redirect struct{}

func (r *Redirect) Decode(input []byte) (*internal.Redirect, error) {
	redirect := &internal.Redirect{}

	if err := json.Unmarshal(input, redirect); err != nil {
		return nil, errors.Wrap(err, "serializer.json.Decode")
	}
	return redirect, nil

}

func (r *Redirect) Encode(input *internal.Redirect) ([]byte, error) {
	jsonStr, err := json.Marshal(input)
	if err != nil {
		return nil, errors.Wrap(err, "serializer.json.Encode")
	}
	return jsonStr, nil
}
