package tachograph

import (
	"bytes"

	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
)

// UnmarshalFaultsData parses the binary data for an EF_Faults_Data record.
func UnmarshalFaultsData(data []byte, fd *cardv1.FaultData) error {
	const recordSize = 24
	r := bytes.NewReader(data)
	for r.Len() >= recordSize {
		recordData := make([]byte, recordSize)
		r.Read(recordData)

		isEmpty := true
		for _, b := range recordData {
			if b != 0 {
				isEmpty = false
				break
			}
		}
		if isEmpty {
			continue
		}

		rec := &cardv1.FaultData_Record{}
		if err := UnmarshalFaultRecord(recordData, rec); err != nil {
			return err
		}
		records := fd.GetRecords()
		records = append(records, rec)
		fd.SetRecords(records)
	}
	return nil
}

// UnmarshalFaultRecord parses a single 24-byte fault record.
func UnmarshalFaultRecord(data []byte, rec *cardv1.FaultData_Record) error {
	r := bytes.NewReader(data)
	faultType, _ := r.ReadByte()

	// Convert raw fault type to enum using protocol annotations
	SetEventFaultType(int32(faultType), rec.SetFaultType, rec.SetUnrecognizedFaultType)

	rec.SetFaultBeginTime(readTimeReal(r))
	rec.SetFaultEndTime(readTimeReal(r))
	// TODO: Read BCD nation code
	r.ReadByte()
	rec.SetVehicleRegistrationNumber(readString(r, 14))
	return nil
}
