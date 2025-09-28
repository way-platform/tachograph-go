package tachograph

import (
	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
)

// unmarshalNationNumeric unmarshals a nation code from a byte slice
func unmarshalNationNumeric(data []byte) (ddv1.NationNumeric, error) {
	if len(data) == 0 {
		return ddv1.NationNumeric_NATION_NUMERIC_UNSPECIFIED, nil
	}
	return ddv1.NationNumeric(data[0]), nil
}
