package tachograph

import (
	datadictionaryv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/datadictionary/v1"
)

// unmarshalNationNumeric unmarshals a nation code from a byte slice
func unmarshalNationNumeric(data []byte) (datadictionaryv1.NationNumeric, error) {
	if len(data) == 0 {
		return datadictionaryv1.NationNumeric_NATION_NUMERIC_UNSPECIFIED, nil
	}
	return datadictionaryv1.NationNumeric(data[0]), nil
}
