package tachograph

// Shared constants for card data structures

// Event and Fault record sizes
const (
	// CardEventRecord and CardFaultRecord both use the same size
	cardEventFaultRecordSize = 24
)

// Place record size
const (
	// PlaceRecord size (includes 2-byte region and 1-byte reserved field)
	placeRecordSize = 12
)

// Specific conditions
const (
	// Total number of specific condition records
	specificConditionTotalRecords = 56
)
