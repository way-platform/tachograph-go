# Specification: Protobuf-driven Raw Card File Type Inference

## 1. Overview

### 1.1. Problem

The current implementation of `InferRawCardFileType` in `filetype.go` has several design flaws:

1.  **Brittle Field Mapping**: It uses a large `switch` statement (`getExpectedEF`) to map Go struct field names to Protobuf `ElementaryFileType` enums. This is fragile and error-prone; any renaming of fields in the `.proto` file will break the logic silently.
2.  **Hardcoded Conditional Logic**: The `isFieldConditional` function contains a hardcoded map of field names that are considered conditional. This logic is based on regulation knowledge that is not formally captured in the schema, making it difficult to verify and maintain.
3.  **Separation of Concerns**: The file structure—a core part of the tachograph specification—is defined in Go application logic instead of within the Protobuf schema itself. The schema should be the single source of truth.

### 1.2. Goal

This document specifies a new implementation that refactors the card type inference logic to be driven directly by the Protobuf schemas. By annotating our card type messages with their expected file structure, we can create a more robust, maintainable, and declarative system.

## 2. Proposed Design

The new design is centered around a custom Protobuf option, `(wayplatform.connect.tachograph.card.v1.file_structure)`. This option will be used to annotate each card type message (e.g., `DriverCardFile`, `WorkshopCardFile`) with a `FileDescriptor` message that describes its complete, ordered file hierarchy.

The `InferRawCardFileType` function in `filetype.go` will be rewritten to use Protobuf reflection:

1.  For each candidate card type, it will dynamically retrieve the `file_structure` annotation from the corresponding message descriptor.
2.  This annotation provides a tree of expected Elementary Files (EFs), including whether each file is mandatory or conditional.
3.  A new matching algorithm will perform a sequential, dual-cursor walk through the list of EFs from the annotation and the list of EFs from the input `RawCardFile`.
4.  This approach eliminates all hardcoded mappings and conditional logic from the Go code, delegating that responsibility to the Protobuf schema where it belongs.

## 3. Protobuf Schema Implementation

### 3.1. Defining the Custom Option

First, we must define the `file_structure` extension. This is done by extending `google.protobuf.MessageOptions`. We will place this definition in a new file, `proto/wayplatform/connect/tachograph/card/v1/options.proto`, for clarity.

**File: `proto/wayplatform/connect/tachograph/card/v1/options.proto`**

```protobuf
syntax = "proto3";

package wayplatform.connect.tachograph.card.v1;

import "google/protobuf/descriptor.proto";
import "wayplatform/connect/tachograph/card/v1/file_descriptor.proto";

extend google.protobuf.MessageOptions {
  // Describes the file system structure for a specific card type message.
  // The annotation provides a FileDescriptor tree that specifies the
  // expected sequence of Elementary Files for that card.
  FileDescriptor file_structure = 50000;
}
```

### 3.2. Annotating Card Type Messages

Next, we apply this option to each of the card type messages (e.g., `DriverCardFile`, `WorkshopCardFile`, etc.). The annotation declaratively defines the expected file structure, including conditional files as specified by the tachograph regulation.

Below is an example for the `DriverCardFile`. Similar annotations would be added for `WorkshopCardFile`, `ControlCardFile`, and `CompanyCardFile`.

**Example: `proto/wayplatform/connect/tachograph/card/v1/driver_card.proto`**

```protobuf
syntax = "proto3";

package wayplatform.connect.tachograph.card.v1;

// Import the new options and other necessary types
import "wayplatform/connect/tachograph/card/v1/options.proto";
import "wayplatform/connect/tachograph/card/v1/file_descriptor.proto";
import "wayplatform/connect/tachograph/card/v1/elementary_file_type.proto";
import "wayplatform/connect/tachograph/card/v1/dedicated_file_type.proto";
import "wayplatform/connect/tachograph/card/v1/file_type.proto";

// Other message imports...

message DriverCardFile {
  // This annotation defines the expected file structure for a Driver Card.
  option (wayplatform.connect.tachograph.card.v1.file_structure) = {
    type: DF,
    df: DF_TACHOGRAPH,
    files: [
      { type: EF, ef: EF_ICC },
      { type: EF, ef: EF_IC },
      { type: EF, ef: EF_APPLICATION_IDENTIFICATION },
      { type: EF, ef: EF_IDENTIFICATION },
      { type: EF, ef: EF_DRIVING_LICENCE_INFO },
      { type: EF, ef: EF_EVENTS_DATA },
      { type: EF, ef: EF_FAULTS_DATA },
      { type: EF, ef: EF_DRIVER_ACTIVITY_DATA },
      { type: EF, ef: EF_VEHICLES_USED },
      { type: EF, ef: EF_PLACES },
      { type: EF, ef: EF_CURRENT_USAGE },
      { type: EF, ef: EF_SPECIFIC_CONDITIONS },
      // Gen2 files are marked as conditional
      { type: EF, ef: EF_APPLICATION_IDENTIFICATION_V2, conditional: true },
      { type: EF, ef: EF_VEHICLE_UNITS_USED, conditional: true },
      { type: EF, ef: EF_GNSS_PLACES, conditional: true },
      { type: EF, ef: EF_PLACES_AUTHENTICATION, conditional: true },
      { type: EF, ef: EF_GNSS_PLACES_AUTHENTICATION, conditional: true },
      { type: EF, ef: EF_BORDER_CROSSINGS, conditional: true },
      { type: EF, ef: EF_LOAD_UNLOAD_OPERATIONS, conditional: true },
      { type: EF, ef: EF_LOAD_TYPE_ENTRIES, conditional: true }
    ]
  };

  // ... existing message fields for DriverCardFile ...
}
```

## 4. Go Implementation Guide

The Go implementation will use the `google.golang.org/protobuf/proto` and `google.golang.org/protobuf/reflect/protoreflect` packages to access the annotations at runtime.

### 4.1. Accessing the Custom Annotation

Here is the key function for retrieving the `file_structure` annotation from a message type. This replaces the need for manual field mapping.

```go
import (
    "google.golang.org/protobuf/proto"
    "google.golang.org/protobuf/reflect/protoreflect"
    cardpb "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
)

// getFileStructure retrieves the file structure annotation from a message type.
// It uses proto.GetExtension to access the custom option we defined.
func getFileStructure(msgType protoreflect.MessageType) *cardpb.FileDescriptor {
    opts := msgType.Descriptor().Options()
    if !proto.HasExtension(opts, cardpb.E_FileStructure) {
        // This indicates a configuration error - the message is missing its annotation.
        return nil
    }
    ext := proto.GetExtension(opts, cardpb.E_FileStructure)
    fileDesc, ok := ext.(*cardpb.FileDescriptor)
    if !ok {
        // This indicates a type mismatch error.
        return nil
    }
    return fileDesc
}
```

### 4.2. New Inference and Matching Logic

The following code outlines the new, annotation-driven `InferRawCardFileType` function and its helpers.

```go
// InferRawCardFileType iterates through known card types, retrieves their
// file structure via protobuf annotations, and checks for a match.
func InferRawCardFileType(input *cardv1.RawCardFile) tachographv1.File_Type {
	if input == nil || len(input.GetRecords()) == 0 {
		return tachographv1.File_TYPE_UNSPECIFIED
	}

	cardTypes := []tachographv1.File_Type{
		tachographv1.File_DRIVER_CARD,
		tachographv1.File_WORKSHOP_CARD,
		tachographv1.File_CONTROL_CARD,
		tachographv1.File_COMPANY_CARD,
	}

	for _, cardType := range cardTypes {
		msgType := getMessageTypeForCardType(cardType)
		if msgType == nil {
			continue
		}

		fileStructure := getFileStructure(msgType)
		if fileStructure == nil {
			// Log an error: card type is missing its structure annotation
			continue
		}

		// Flatten the file structure to get an ordered list of expected EFs.
		expectedEFs := flattenFileStructure(fileStructure)

		if matchesStructure(input.GetRecords(), expectedEFs) {
			return cardType
		}
	}

	return tachographv1.File_TYPE_UNSPECIFIED
}

// flattenFileStructure recursively traverses the FileDescriptor tree
// and returns a flat, ordered list of file descriptors for all EFs.
func flattenFileStructure(desc *cardpb.FileDescriptor) []*cardpb.FileDescriptor {
	var flatList []*cardpb.FileDescriptor
	if desc.GetType() == cardpb.FileType_EF {
		flatList = append(flatList, desc)
	}
	for _, child := range desc.GetFiles() {
		flatList = append(flatList, flattenFileStructure(child)...)
	}
	return flatList
}

// matchesStructure uses a dual-cursor approach to match raw records against the expected EF structure.
func matchesStructure(records []*cardv1.RawCardFile_Record, expectedEFs []*cardpb.FileDescriptor) bool {
	recordCursor := 0
	efCursor := 0

	for recordCursor < len(records) && efCursor < len(expectedEFs) {
		record := records[recordCursor]
		expectedEF := expectedEFs[efCursor]

		if record.GetFile() == expectedEF.GetEf() {
			// Match found, advance both cursors.
			recordCursor++
			efCursor++
		} else if expectedEF.GetConditional() {
			// Expected EF is conditional and not present in the input file.
			// Advance the EF cursor to check the next expected file.
			efCursor++
		} else {
			// A required EF is missing or out of order. This card type does not match.
			return false
		}
	}

	// After consuming all records, ensure any remaining expected EFs are all conditional.
	// This handles cases where trailing files are conditional and not present.
	for efCursor < len(expectedEFs) {
		if !expectedEFs[efCursor].GetConditional() {
			return false // A required EF was not found in the input file.
		}
		efCursor++
	}

	return true
}
```

## 5. Advantages of the New Approach

*   **Single Source of Truth**: The file structure is formally defined in the `.proto` schema, directly alongside the data structures.
*   **Enhanced Maintainability**: When the tachograph regulation is updated to change a card's file structure, only the Protobuf annotation needs to be changed. The Go logic remains untouched.
*   **Improved Clarity and Readability**: The `file_structure` annotation provides a clear, declarative specification of each card type's layout, making the schema easier to understand.
*   **Reduced Code Complexity**: The new approach completely eliminates the brittle `switch` statements, hardcoded maps, and manual mappings in `filetype.go`, resulting in cleaner and more reliable code.