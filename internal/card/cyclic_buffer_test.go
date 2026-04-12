package card

import (
	"encoding/binary"
	"strings"
	"testing"
)

// buildRecord constructs a single cyclic buffer record with the 4-byte header
// (prevRecordLength, currentRecordLength) followed by payload bytes.
func buildRecord(prevLen, currentLen uint16, payload []byte) []byte {
	rec := make([]byte, currentLen)
	binary.BigEndian.PutUint16(rec[0:2], prevLen)
	binary.BigEndian.PutUint16(rec[2:4], currentLen)
	copy(rec[4:], payload)
	return rec
}

func TestCyclicIterator_SingleRecord(t *testing.T) {
	// One record with prevRecordLength=0 (signals no previous record).
	payload := []byte{0xAA, 0xBB}
	rec := buildRecord(0, 6, payload) // 4-byte header + 2-byte payload = 6

	it := NewCyclicRecordIterator(rec, 0, 0)

	if !it.Next() {
		t.Fatal("expected first Next() to return true")
	}
	gotBytes, gotPos, gotLen := it.Record()
	if gotPos != 0 {
		t.Errorf("position: got %d, want 0", gotPos)
	}
	if gotLen != 6 {
		t.Errorf("length: got %d, want 6", gotLen)
	}
	if len(gotBytes) != 6 {
		t.Errorf("record bytes length: got %d, want 6", len(gotBytes))
	}
	// Payload should be at bytes [4:6].
	if gotBytes[4] != 0xAA || gotBytes[5] != 0xBB {
		t.Errorf("payload: got [%02x %02x], want [aa bb]", gotBytes[4], gotBytes[5])
	}

	if it.Next() {
		t.Fatal("expected second Next() to return false (prevRecordLength=0 terminates)")
	}
	if err := it.Err(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCyclicIterator_TwoRecordsNoWrap(t *testing.T) {
	// Two records laid out linearly: [rec0][rec1].
	// rec1 is the newest (startPos), rec0 is the oldest (oldestPos).
	// rec0: prevLen=0 (first record), currentLen=6, payload=0x11 0x22
	// rec1: prevLen=6 (points back 6 bytes to rec0), currentLen=8, payload=0x33 0x44 0x55 0x66
	rec0 := buildRecord(0, 6, []byte{0x11, 0x22})
	rec1 := buildRecord(6, 8, []byte{0x33, 0x44, 0x55, 0x66})

	buf := append(rec0, rec1...)
	startPos := 6  // rec1 starts at offset 6
	oldestPos := 0 // rec0 starts at offset 0

	it := NewCyclicRecordIterator(buf, startPos, oldestPos)

	// First iteration: newest record (rec1).
	if !it.Next() {
		t.Fatal("expected first Next() to return true (rec1)")
	}
	gotBytes, gotPos, gotLen := it.Record()
	if gotPos != 6 {
		t.Errorf("rec1 position: got %d, want 6", gotPos)
	}
	if gotLen != 8 {
		t.Errorf("rec1 length: got %d, want 8", gotLen)
	}
	if gotBytes[4] != 0x33 {
		t.Errorf("rec1 first payload byte: got %02x, want 33", gotBytes[4])
	}

	// Second iteration: oldest record (rec0).
	if !it.Next() {
		t.Fatal("expected second Next() to return true (rec0)")
	}
	gotBytes, gotPos, gotLen = it.Record()
	if gotPos != 0 {
		t.Errorf("rec0 position: got %d, want 0", gotPos)
	}
	if gotLen != 6 {
		t.Errorf("rec0 length: got %d, want 6", gotLen)
	}
	if gotBytes[4] != 0x11 {
		t.Errorf("rec0 first payload byte: got %02x, want 11", gotBytes[4])
	}

	// No more records (rec0 is oldestPos, iterator stops after it).
	if it.Next() {
		t.Fatal("expected third Next() to return false")
	}
	if err := it.Err(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCyclicIterator_WrapAround(t *testing.T) {
	// Simulate a cyclic buffer where the newest record is near the start and the
	// oldest record is near the end. prevRecordLength on the newest record causes
	// a wrap-around to the older record at the end of the buffer.
	//
	// Buffer layout (20 bytes total):
	//   [0..7]   = rec1 (newest), 8 bytes, prevLen=8 → wraps to offset 12
	//   [8..11]  = padding (zeros)
	//   [12..19] = rec0 (oldest), 8 bytes, prevLen=0
	buf := make([]byte, 20)

	// rec0 at offset 12: prevLen=0, currentLen=8, payload=0xDD 0xEE 0xFF 0xAA
	binary.BigEndian.PutUint16(buf[12:14], 0)   // prevRecordLength
	binary.BigEndian.PutUint16(buf[14:16], 8)   // currentRecordLength
	buf[16], buf[17], buf[18], buf[19] = 0xDD, 0xEE, 0xFF, 0xAA

	// rec1 at offset 0: prevLen=8, currentLen=8, payload=0x01 0x02 0x03 0x04
	// currentPos(0) - prevRecordLength(8) = -8 → -8 + 20 = 12 (wraps to rec0)
	binary.BigEndian.PutUint16(buf[0:2], 8)     // prevRecordLength
	binary.BigEndian.PutUint16(buf[2:4], 8)     // currentRecordLength
	buf[4], buf[5], buf[6], buf[7] = 0x01, 0x02, 0x03, 0x04

	startPos := 0   // newest record
	oldestPos := 12  // oldest record

	it := NewCyclicRecordIterator(buf, startPos, oldestPos)

	// First: newest (rec1 at offset 0).
	if !it.Next() {
		t.Fatal("expected first Next() to return true")
	}
	_, gotPos, _ := it.Record()
	if gotPos != 0 {
		t.Errorf("rec1 position: got %d, want 0", gotPos)
	}

	// Second: oldest (rec0 at offset 12, reached via wrap-around).
	if !it.Next() {
		t.Fatal("expected second Next() to return true (wrap-around to rec0)")
	}
	_, gotPos, _ = it.Record()
	if gotPos != 12 {
		t.Errorf("rec0 position: got %d, want 12", gotPos)
	}

	// No more.
	if it.Next() {
		t.Fatal("expected third Next() to return false")
	}
	if err := it.Err(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCyclicIterator_EmptyBuffer(t *testing.T) {
	it := NewCyclicRecordIterator(nil, 0, 0)
	if it.Next() {
		t.Fatal("expected Next() to return false on nil buffer")
	}
	if err := it.Err(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	it2 := NewCyclicRecordIterator([]byte{}, 0, 0)
	if it2.Next() {
		t.Fatal("expected Next() to return false on empty buffer")
	}
	if err := it2.Err(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCyclicIterator_ZeroLengthRecord(t *testing.T) {
	// A record with currentRecordLength=0 signals end of chain.
	buf := make([]byte, 8)
	binary.BigEndian.PutUint16(buf[0:2], 0) // prevRecordLength
	binary.BigEndian.PutUint16(buf[2:4], 0) // currentRecordLength = 0 → stop
	// Remaining bytes are irrelevant.

	it := NewCyclicRecordIterator(buf, 0, 0)
	if it.Next() {
		t.Fatal("expected Next() to return false when currentRecordLength is 0")
	}
	if err := it.Err(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCyclicIterator_InvalidRecordLength(t *testing.T) {
	// A record with currentRecordLength < 4 is invalid (cannot hold the header).
	buf := make([]byte, 8)
	binary.BigEndian.PutUint16(buf[0:2], 0) // prevRecordLength
	binary.BigEndian.PutUint16(buf[2:4], 3) // currentRecordLength = 3 → invalid

	it := NewCyclicRecordIterator(buf, 0, 0)
	if it.Next() {
		t.Fatal("expected Next() to return false for invalid record length")
	}
	if err := it.Err(); err == nil {
		t.Fatal("expected an error for currentRecordLength < 4, got nil")
	} else {
		wantSubstr := "invalid record length"
		if got := err.Error(); !strings.Contains(got, wantSubstr) {
			t.Errorf("error message: got %q, want substring %q", got, wantSubstr)
		}
	}
}

func TestCyclicIterator_MaxRecordsBackstop(t *testing.T) {
	// Build a small cyclic buffer where two records point at each other,
	// creating an infinite loop. The 5000-record backstop must stop it.
	//
	// However, cycle detection (seen map) will fire first for a 2-record loop.
	// To truly test the backstop, we need records that keep visiting new positions.
	//
	// Strategy: build a large buffer where each record points to the next position
	// forward (backwards walk wraps around the full buffer). With a buffer of
	// 5001 * 5 = 25005 bytes, each record is 5 bytes (4 header + 1 payload),
	// prevLen=5 always. This creates 5001 unique positions, so cycle detection
	// won't fire before the 5000 backstop.
	const recordSize = 5
	const numRecords = 5001
	bufSize := numRecords * recordSize
	buf := make([]byte, bufSize)

	for i := 0; i < numRecords; i++ {
		off := i * recordSize
		binary.BigEndian.PutUint16(buf[off:off+2], recordSize) // prevRecordLength
		binary.BigEndian.PutUint16(buf[off+2:off+4], recordSize) // currentRecordLength
		buf[off+4] = byte(i & 0xFF) // payload
	}
	// Fix the first record so it wraps to the last (it already does via modular
	// arithmetic: pos 0 - 5 = -5 → -5 + bufSize wraps correctly).

	// Start at position 0, set oldestPos to a position we'll never naturally reach
	// (use bufSize-recordSize but since the chain is circular, the backstop should
	// trigger before we get through all 5001 records).
	// To prevent the oldestPos termination, set it to a position that isn't on
	// the backward chain from 0. Since we go 0 → bufSize-5 → bufSize-10 → ...,
	// set oldestPos to 1 (not aligned to recordSize boundaries).
	it := NewCyclicRecordIterator(buf, 0, 1)

	count := 0
	for it.Next() {
		count++
	}

	if count != 5000 {
		t.Errorf("record count: got %d, want 5000 (backstop)", count)
	}
	if err := it.Err(); err == nil {
		t.Fatal("expected an error from the max records backstop, got nil")
	} else {
		wantSubstr := "exceeded maximum record count"
		if got := err.Error(); !strings.Contains(got, wantSubstr) {
			t.Errorf("error message: got %q, want substring %q", got, wantSubstr)
		}
	}
}

func TestCyclicIterator_CycleDetection(t *testing.T) {
	// Two records that form a cycle: each points to the other.
	// Buffer is exactly 16 bytes (2 * 8-byte records):
	//   rec0 at 0: prevLen=8 → 0-8=-8 → +16=8 (rec1)
	//   rec1 at 8: prevLen=8 → 8-8=0 (rec0)
	// Chain: start at rec1(8) → rec0(0) → rec1(8, already seen!) → stop.
	// oldestPos set off-chain so it doesn't trigger first.
	buf2 := make([]byte, 16)
	binary.BigEndian.PutUint16(buf2[0:2], 8) // rec0 prevLen
	binary.BigEndian.PutUint16(buf2[2:4], 8) // rec0 currentLen
	buf2[4], buf2[5], buf2[6], buf2[7] = 0xA0, 0xA1, 0xA2, 0xA3

	binary.BigEndian.PutUint16(buf2[8:10], 8) // rec1 prevLen
	binary.BigEndian.PutUint16(buf2[10:12], 8) // rec1 currentLen
	buf2[12], buf2[13], buf2[14], buf2[15] = 0xB0, 0xB1, 0xB2, 0xB3

	// Start at rec1 (offset 8). oldestPos at some position not on chain.
	it := NewCyclicRecordIterator(buf2, 8, 99)

	// First: rec1.
	if !it.Next() {
		t.Fatal("expected first Next() to return true (rec1)")
	}
	_, pos1, _ := it.Record()
	if pos1 != 8 {
		t.Errorf("first record position: got %d, want 8", pos1)
	}

	// Second: rec0.
	if !it.Next() {
		t.Fatal("expected second Next() to return true (rec0)")
	}
	_, pos2, _ := it.Record()
	if pos2 != 0 {
		t.Errorf("second record position: got %d, want 0", pos2)
	}

	// Third: would be rec1 again (offset 8), but it's in the seen map.
	if it.Next() {
		t.Fatal("expected third Next() to return false (cycle detected)")
	}

	// Cycle detection is silent (no error set).
	if err := it.Err(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCyclicIterator_RecordBytesWrapAround(t *testing.T) {
	// A record whose body physically wraps around the buffer boundary.
	// The header read (currentPos..currentPos+4) is direct (no modular arithmetic),
	// so the header must fit contiguously. But the byte extraction loop uses
	// buffer[(currentPos+i) % len(buffer)], so the payload can wrap.
	//
	// Buffer: 16 bytes. Record at offset 12, currentLen=8:
	//   header at [12..15] (contiguous), payload wraps to [0..3].
	buf2 := make([]byte, 16)
	binary.BigEndian.PutUint16(buf2[12:14], 0) // prevRecordLength = 0
	binary.BigEndian.PutUint16(buf2[14:16], 8) // currentRecordLength = 8

	// The payload is 4 bytes. Due to the modular read, bytes [0..3] of the
	// buffer become payload. Set them to a known pattern.
	buf2[0] = 0xCA
	buf2[1] = 0xFE
	buf2[2] = 0xBA
	buf2[3] = 0xBE

	it := NewCyclicRecordIterator(buf2, 12, 12)

	if !it.Next() {
		t.Fatal("expected Next() to return true")
	}
	gotBytes, _, gotLen := it.Record()
	if gotLen != 8 {
		t.Errorf("length: got %d, want 8", gotLen)
	}
	if len(gotBytes) != 8 {
		t.Fatalf("record bytes length: got %d, want 8", len(gotBytes))
	}

	// First 4 bytes are the header (from offset 12..15).
	wantPrev := uint16(0)
	gotPrev := binary.BigEndian.Uint16(gotBytes[0:2])
	if gotPrev != wantPrev {
		t.Errorf("prevRecordLength in extracted bytes: got %d, want %d", gotPrev, wantPrev)
	}
	wantCurr := uint16(8)
	gotCurr := binary.BigEndian.Uint16(gotBytes[2:4])
	if gotCurr != wantCurr {
		t.Errorf("currentRecordLength in extracted bytes: got %d, want %d", gotCurr, wantCurr)
	}

	// Payload bytes should be [0xCA, 0xFE, 0xBA, 0xBE] (wrapped from start of buffer).
	wantPayload := []byte{0xCA, 0xFE, 0xBA, 0xBE}
	for i, w := range wantPayload {
		if gotBytes[4+i] != w {
			t.Errorf("payload[%d]: got %02x, want %02x", i, gotBytes[4+i], w)
		}
	}

	if it.Next() {
		t.Fatal("expected second Next() to return false (prevRecordLength=0)")
	}
	if err := it.Err(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCyclicIterator_StopsAtOldestPos(t *testing.T) {
	// Three records in linear order: [rec0][rec1][rec2].
	// Start at rec2 (newest), oldestPos at rec1 (middle).
	// The iterator should return rec2 and rec1, then stop — not follow rec1's
	// prevRecordLength to rec0.
	rec0 := buildRecord(0, 6, []byte{0x00, 0x01})   // offset 0
	rec1 := buildRecord(6, 8, []byte{0x10, 0x11, 0x12, 0x13}) // offset 6
	rec2 := buildRecord(8, 6, []byte{0x20, 0x21})    // offset 14

	buf := make([]byte, 0, len(rec0)+len(rec1)+len(rec2))
	buf = append(buf, rec0...)
	buf = append(buf, rec1...)
	buf = append(buf, rec2...)

	startPos := 14 // rec2
	oldestPos := 6 // rec1

	it := NewCyclicRecordIterator(buf, startPos, oldestPos)

	// First: rec2.
	if !it.Next() {
		t.Fatal("expected first Next() to return true (rec2)")
	}
	_, gotPos, _ := it.Record()
	if gotPos != 14 {
		t.Errorf("first record position: got %d, want 14", gotPos)
	}

	// Second: rec1 (which is at oldestPos).
	if !it.Next() {
		t.Fatal("expected second Next() to return true (rec1 at oldestPos)")
	}
	_, gotPos, _ = it.Record()
	if gotPos != 6 {
		t.Errorf("second record position: got %d, want 6", gotPos)
	}

	// Third: should be false — rec1 was at oldestPos, so iterator doesn't continue.
	if it.Next() {
		_, gotPos, _ = it.Record()
		t.Fatalf("expected third Next() to return false (should stop after oldestPos), but got record at %d", gotPos)
	}
	if err := it.Err(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCyclicIterator_HeaderOutOfBounds(t *testing.T) {
	// If startPos is close to the end of the buffer such that the 4-byte header
	// doesn't fit, Next() should return false without error.
	buf := make([]byte, 6)
	// Place valid data at offset 0 so the buffer isn't all zeros.
	binary.BigEndian.PutUint16(buf[0:2], 0)
	binary.BigEndian.PutUint16(buf[2:4], 6)

	// startPos = 4 → needs bytes [4..7] for header, but buffer is only 6 bytes.
	it := NewCyclicRecordIterator(buf, 4, 0)
	if it.Next() {
		t.Fatal("expected Next() to return false when header doesn't fit")
	}
	if err := it.Err(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

