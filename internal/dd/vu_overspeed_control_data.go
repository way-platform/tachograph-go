package dd

import (
	"fmt"

	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
)

// UnmarshalVuOverspeedControlData parses a VuOverspeedControlData (Generation 1).
//
// The data type `VuOverSpeedingControlData` is specified in the Data Dictionary, Section 2.212.
//
// ASN.1 Specification (Gen1):
//
//	VuOverSpeedingControlData ::= SEQUENCE {
//	    lastOverspeedControlTime    TimeReal,           -- 4 bytes
//	    firstOverspeedSince         TimeReal,           -- 4 bytes
//	    numberOfOverspeedSince      OverspeedNumber     -- 1 byte
//	}
func (opts UnmarshalOptions) UnmarshalVuOverspeedControlData(data []byte) (*ddv1.VuOverspeedControlData, error) {
	const (
		idxLastOverspeedControlTime = 0
		lenLastOverspeedControlTime = 4
		idxFirstOverspeedSince      = 4
		lenFirstOverspeedSince      = 4
		idxNumberOfOverspeedSince   = 8
		lenVuOverspeedControlData   = 9
	)

	if len(data) != lenVuOverspeedControlData {
		return nil, fmt.Errorf("invalid length for VuOverspeedControlData: got %d, want %d", len(data), lenVuOverspeedControlData)
	}

	controlData := &ddv1.VuOverspeedControlData{}
	if opts.PreserveRawData {
		controlData.SetRawData(data)
	}

	// Parse lastOverspeedControlTime (4 bytes)
	lastControlTime, err := opts.UnmarshalTimeReal(data[idxLastOverspeedControlTime : idxLastOverspeedControlTime+lenLastOverspeedControlTime])
	if err != nil {
		return nil, fmt.Errorf("failed to parse last overspeed control time: %w", err)
	}
	controlData.SetLastOverspeedControlTime(lastControlTime)

	// Parse firstOverspeedSince (4 bytes)
	firstOverspeedSince, err := opts.UnmarshalTimeReal(data[idxFirstOverspeedSince : idxFirstOverspeedSince+lenFirstOverspeedSince])
	if err != nil {
		return nil, fmt.Errorf("failed to parse first overspeed since: %w", err)
	}
	controlData.SetFirstOverspeedSinceLastControl(firstOverspeedSince)

	// Parse numberOfOverspeedSince (1 byte)
	controlData.SetNumberOfOverspeedSinceLastControl(int32(data[idxNumberOfOverspeedSince]))

	return controlData, nil
}

// MarshalVuOverspeedControlData marshals a VuOverspeedControlData to binary format (Generation 1).
func (opts MarshalOptions) MarshalVuOverspeedControlData(controlData *ddv1.VuOverspeedControlData) ([]byte, error) {
	const lenVuOverspeedControlData = 9

	if controlData == nil {
		return nil, fmt.Errorf("controlData cannot be nil")
	}

	// Use raw data painting strategy if available
	var canvas [lenVuOverspeedControlData]byte
	if controlData.HasRawData() {
		if len(controlData.GetRawData()) != lenVuOverspeedControlData {
			return nil, fmt.Errorf(
				"invalid raw_data length for VuOverspeedControlData: got %d, want %d",
				len(controlData.GetRawData()), lenVuOverspeedControlData,
			)
		}
		copy(canvas[:], controlData.GetRawData())
	}

	// Paint semantic values over the canvas
	const (
		idxLastOverspeedControlTime = 0
		lenLastOverspeedControlTime = 4
		idxFirstOverspeedSince      = 4
		lenFirstOverspeedSince      = 4
		idxNumberOfOverspeedSince   = 8
	)

	// Marshal lastOverspeedControlTime (4 bytes)
	lastControlTime, err := opts.MarshalTimeReal(controlData.GetLastOverspeedControlTime())
	if err != nil {
		return nil, fmt.Errorf("failed to marshal last overspeed control time: %w", err)
	}
	if controlData.GetLastOverspeedControlTime() != nil {
		copy(canvas[idxLastOverspeedControlTime:idxLastOverspeedControlTime+lenLastOverspeedControlTime], lastControlTime)
	}

	// Marshal firstOverspeedSince (4 bytes)
	firstOverspeedSince, err := opts.MarshalTimeReal(controlData.GetFirstOverspeedSinceLastControl())
	if err != nil {
		return nil, fmt.Errorf("failed to marshal first overspeed since: %w", err)
	}
	if controlData.GetFirstOverspeedSinceLastControl() != nil {
		copy(canvas[idxFirstOverspeedSince:idxFirstOverspeedSince+lenFirstOverspeedSince], firstOverspeedSince)
	}

	// Marshal numberOfOverspeedSince (1 byte)
	canvas[idxNumberOfOverspeedSince] = byte(controlData.GetNumberOfOverspeedSinceLastControl())

	return canvas[:], nil
}
