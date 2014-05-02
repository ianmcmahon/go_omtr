package omtr

import (
	"encoding/json"
)

type Number int64

var _ = json.Unmarshaler(new(Number))

func (n *Number) UnmarshalJSON(data []byte) error {
	var num json.Number
	err := json.Unmarshal(data, &num)
	if err != nil { return err }

	*(*int64)(n), err = num.Int64()
	return err
}
