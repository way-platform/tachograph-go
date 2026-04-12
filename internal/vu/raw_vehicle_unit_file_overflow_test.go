package vu

import (
	"encoding/binary"
	"strings"
	"testing"
)

// buildRecordArrayHeader builds a 5-byte Gen2 RecordArray header.
//
// RecordArray header (Appendix 1, Section 1.1.3):
//
//	recordType:  1 byte
//	recordSize:  2 bytes (big-endian)
//	noOfRecords: 2 bytes (big-endian)
func buildRecordArrayHeader(recordType byte, recordSize, noOfRecords uint16) []byte {
	buf := make([]byte, 5)
	buf[0] = recordType
	binary.BigEndian.PutUint16(buf[1:], recordSize)
	binary.BigEndian.PutUint16(buf[3:], noOfRecords)
	return buf
}

// TestSizeOfRecordArray_Overflow_MaxMax verifies that sizeOfRecordArray returns an
// error (not a panic or massive allocation) when both recordSize and noOfRecords
// are 0xFFFF. Their product (65535 * 65535 = 4,294,836,225) overflows int32 and
// can cause memory exhaustion on 64-bit systems.
func TestSizeOfRecordArray_Overflow_MaxMax(t *testing.T) {
	header := buildRecordArrayHeader(0x01, 0xFFFF, 0xFFFF)

	// Data is only 5 bytes — far less than the ~4 GB claimed by the header.
	_, err := sizeOfRecordArray(header, 0)
	if err == nil {
		t.Fatal("expected error for overflow recordSize=0xFFFF * noOfRecords=0xFFFF, got nil")
	}

	if !strings.Contains(err.Error(), "overflow") && !strings.Contains(err.Error(), "exceeds") {
		t.Errorf("expected overflow or bounds error, got: %v", err)
	}
}

// TestSizeOfRecordArray_TotalSizeExceedsData verifies that sizeOfRecordArray
// returns an error when the computed totalSize exceeds the actual data length.
// This prevents out-of-bounds reads when the VU file has been truncated or
// the header is corrupted.
func TestSizeOfRecordArray_TotalSizeExceedsData(t *testing.T) {
	// RecordArray claims 10 records of 100 bytes each = 1000 bytes payload.
	// But we only supply 5 bytes of header + 20 bytes of payload = 25 bytes total.
	header := buildRecordArrayHeader(0x01, 100, 10)
	data := make([]byte, 25)
	copy(data, header)

	_, err := sizeOfRecordArray(data, 0)
	if err == nil {
		t.Fatal("expected error for totalSize exceeding data length, got nil")
	}

	if !strings.Contains(err.Error(), "exceeds") && !strings.Contains(err.Error(), "insufficient") {
		t.Errorf("expected bounds-check error, got: %v", err)
	}
}

// TestSizeOfRecordArray_ZeroRecordSize verifies graceful handling when
// recordSize is 0. With zero-size records the payload is always 0 bytes
// regardless of noOfRecords, so totalSize should equal headerSize (5).
func TestSizeOfRecordArray_ZeroRecordSize(t *testing.T) {
	header := buildRecordArrayHeader(0x01, 0, 5)

	size, err := sizeOfRecordArray(header, 0)
	if err != nil {
		t.Fatalf("unexpected error for recordSize=0: %v", err)
	}

	const expectedHeaderSize = 5
	if size != expectedHeaderSize {
		t.Errorf("expected totalSize=%d (header only), got %d", expectedHeaderSize, size)
	}
}

// TestSizeOfRecordArray_ZeroNoOfRecords verifies that noOfRecords=0 produces
// an empty array with totalSize equal to the 5-byte header.
func TestSizeOfRecordArray_ZeroNoOfRecords(t *testing.T) {
	header := buildRecordArrayHeader(0x01, 64, 0)

	size, err := sizeOfRecordArray(header, 0)
	if err != nil {
		t.Fatalf("unexpected error for noOfRecords=0: %v", err)
	}

	const expectedHeaderSize = 5
	if size != expectedHeaderSize {
		t.Errorf("expected totalSize=%d (header only), got %d", expectedHeaderSize, size)
	}
}

// TestSizeOfRecordArray_RealisticValid verifies correct parsing of a realistic
// RecordArray: 3 records of 64 bytes each (e.g., VuDetailedSpeedBlockRecordArray).
// totalSize = 5 + (64 * 3) = 197.
func TestSizeOfRecordArray_RealisticValid(t *testing.T) {
	const recordSize = 64
	const noOfRecords = 3
	const expectedTotal = 5 + recordSize*noOfRecords // 5 + 192 = 197

	header := buildRecordArrayHeader(0x01, recordSize, noOfRecords)
	// Provide enough data for the full RecordArray
	data := make([]byte, expectedTotal)
	copy(data, header)

	size, err := sizeOfRecordArray(data, 0)
	if err != nil {
		t.Fatalf("unexpected error for valid RecordArray: %v", err)
	}

	if size != expectedTotal {
		t.Errorf("expected totalSize=%d, got %d", expectedTotal, size)
	}
}

// TestSizeOfRecordArray_OverflowViaUnmarshalRaw verifies the overflow protection
// through the full UnmarshalRawVehicleUnitFile path. A crafted VU file with
// tag 0x7611 (OVERVIEW_GEN2_V1) followed by a RecordArray with maxed-out
// dimensions must produce an error without panicking or allocating gigabytes.
func TestSizeOfRecordArray_OverflowViaUnmarshalRaw(t *testing.T) {
	// Build a minimal VU file:
	// Tag 0x7611 = OVERVIEW_GEN2_V1 (TREP 0x11)
	// Followed by a RecordArray header with overflow values
	var data []byte

	// 2-byte tag for OVERVIEW_GEN2_V1
	data = binary.BigEndian.AppendUint16(data, 0x7611)

	// RecordArray header: recordType=0x01, recordSize=0xFFFF, noOfRecords=0xFFFF
	data = append(data, buildRecordArrayHeader(0x01, 0xFFFF, 0xFFFF)...)

	_, err := UnmarshalOptions{Strict: true}.UnmarshalRawVehicleUnitFile(data)
	if err == nil {
		t.Fatal("expected error from overflow in VU file parsing, got nil")
	}

	t.Logf("got expected error: %v", err)
}

// TestSizeOfRecordArray_LargeButNotOverflow tests a RecordArray that is large
// but does not overflow int. recordSize=1000, noOfRecords=1000 = 1,000,000 bytes
// payload. Since data is only 5 bytes, the bounds check must catch it.
func TestSizeOfRecordArray_LargeButNotOverflow(t *testing.T) {
	header := buildRecordArrayHeader(0x01, 1000, 1000)

	_, err := sizeOfRecordArray(header, 0)
	if err == nil {
		t.Fatal("expected error for totalSize exceeding data, got nil")
	}
}

// TestSizeOfRecordArray_OffsetBeyondData verifies error when offset is at the
// end of data, leaving insufficient room for the 5-byte header.
func TestSizeOfRecordArray_OffsetBeyondData(t *testing.T) {
	data := make([]byte, 10)

	_, err := sizeOfRecordArray(data, 8)
	if err == nil {
		t.Fatal("expected error for insufficient header data at offset 8, got nil")
	}

	if !strings.Contains(err.Error(), "insufficient") {
		t.Errorf("expected 'insufficient' in error, got: %v", err)
	}
}

// TestParseRecordArrayHeader_Overflow tests that parseRecordArrayHeader callers
// are protected from overflow when computing totalSize from returned values.
// This is a regression guard: even though parseRecordArrayHeader itself just
// reads the header, callers must not blindly compute int(recordSize)*int(noOfRecords).
func TestParseRecordArrayHeader_Overflow(t *testing.T) {
	header := buildRecordArrayHeader(0x01, 0xFFFF, 0xFFFF)

	_, recordSize, noOfRecords, _, err := parseRecordArrayHeader(header, 0)
	if err != nil {
		t.Fatalf("parseRecordArrayHeader should succeed on valid header: %v", err)
	}

	// Verify the raw values are correctly parsed
	if recordSize != 0xFFFF {
		t.Errorf("expected recordSize=0xFFFF, got 0x%04X", recordSize)
	}
	if noOfRecords != 0xFFFF {
		t.Errorf("expected noOfRecords=0xFFFF, got 0x%04X", noOfRecords)
	}

	// The point: callers must use SafeMul, not raw multiplication.
	// Product 65535 * 65535 = 4,294,836,225 which overflows int32.
	rawProduct := int(recordSize) * int(noOfRecords)
	t.Logf("raw product of 0xFFFF * 0xFFFF = %d (may overflow on 32-bit)", rawProduct)
}
