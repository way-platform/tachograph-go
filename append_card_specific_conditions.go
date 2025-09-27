package tachograph

import (
	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
)

// AppendCardSpecificConditions appends specific conditions data to a byte slice.
func AppendCardSpecificConditions(data []byte, conditions *cardv1.SpecificConditions) ([]byte, error) {
	if conditions == nil {
		return data, nil
	}

	records := conditions.GetRecords()
	for _, record := range records {
		if record == nil {
			continue
		}

		// Entry time (4 bytes)
		data = appendTimeReal(data, record.GetEntryTime())

		// Specific condition type (1 byte) - convert enum to protocol value
		conditionTypeProtocol := GetSpecificConditionTypeProtocolValue(record.GetSpecificConditionType(), 0)
		data = append(data, byte(conditionTypeProtocol))
	}

	return data, nil
}
