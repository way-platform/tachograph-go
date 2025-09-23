# Specification: Unmarshalling Completion

## 1. Introduction

This document outlines the necessary steps to complete the unmarshalling functionality for tachograph data files (`.DDD`). The current implementation in `unmarshal.go` is partial and supports only a subset of the data structures defined in the EU regulations for both driver card and vehicle unit (VU) files.

The goal is to achieve comprehensive parsing of all standard data blocks, using the existing code structure and patterns as a foundation.

## 2. Existing Architecture

The current unmarshalling logic is orchestrated from `unmarshal.go` and follows these patterns:

- **`UnmarshalFile(data []byte)`**: The main entry point. It infers the file type (Card or VU) and calls the appropriate handler.
- **`unmarshalCard(data []byte)`**: Handles driver card data. It reads data in a Tag-Length-Value (TLV) format. It iterates through the file, identifies Elementary Files (EFs) by their 3-byte tag (FID + appendix), and dispatches the data `value` to a specific unmarshalling function.
- **`unmarshalVU(data []byte)`**: Handles vehicle unit data. It reads data in a Tag-Value (TV) format. It identifies data blocks by their 2-byte tag and dispatches to a specific unmarshalling function based on the transfer type (TREP).
- **Protobuf Definitions**: The target data structures are defined as Protobuf messages in `proto/`.
- **Separation of Concerns**: Unmarshalling logic for specific data blocks is located in separate files (e.g., `unmarshal_card_icc.go`, `unmarshal_vu_overview.go`).

This is a solid foundation that should be extended.

## 3. Completing Card Data Unmarshalling

The `unmarshalCard` function currently handles only a few Elementary Files. The following EFs need to be implemented.

### 3.1. Missing Elementary Files (EFs)

Based on `tachoparser` and EU regulations, the following EFs are missing:

- **EF Card Chip Identification** (FID `0x0005`)
- **EF Application Identification** (FID `0x0501`)
- **EF Driver Activity Data** (FID `0x0504`)
- **EF Vehicles Used** (FID `0x0505`)
- **EF Places** (FID `0x0506`)
- **EF Current Usage** (FID `0x0507`)
- **EF Control Activity Data** (FID `0x0508`)
- **EF Specific Conditions** (FID `0x0522`)
- **EF Last Card Download** (FID `0x050E`)
- **EF Card Vehicle Units Used** (FID `0x0523`, Gen2)
- **EF GNSS Accumulated Driving** (FID `0x0524`, Gen2)
- And others related to Gen2 v2 (Border Crossings, Load/Unload Operations, etc.)

### 3.2. Implementation Plan (per EF)

For each missing Elementary File, the following steps should be taken:

1.  **Define Protobuf Message**: In the appropriate `.proto` file (e.g., `proto/wayplatform/connect/tachograph/card/v1/file.proto`), define the message(s) that represent the data structure of the EF. Use `benchmark/tachoparser/pkg/decoder/definitions.go` as a reference for the Go struct, which can be translated to a Protobuf message.
2.  **Create Unmarshal File**: Create a new file named `unmarshal_card_FILENAME.go` (e.g., `unmarshal_card_activity.go`).
3.  **Implement Unmarshal Function**: Inside the new file, implement the unmarshalling function (e.g., `UnmarshalCardActivityData(data []byte, target *cardv1.CardActivityData)`). This function will parse the raw byte slice into the target protobuf message.
4.  **Integrate into `unmarshalCard`**: Add a new `case` to the `switch` statement in `unmarshal.go:unmarshalCard`. This case will match the corresponding `ElementaryFileType` and call the newly created unmarshalling function.

    ```go
    // Example for CardDriverActivity in unmarshal.go
    case cardv1.ElementaryFileType_EF_DRIVER_ACTIVITY:
        activityData := &cardv1.DriverActivity{}
        if err := UnmarshalCardActivityData(value, activityData); err != nil {
            return nil, err
        }
        driverCard.SetDriverActivity(activityData)
    ```

## 4. Completing Vehicle Unit (VU) Data Unmarshalling

The `unmarshalVU` function currently only parses the overview and interface version blocks. The remaining data blocks need to be implemented.

### 4.1. Missing VU Data Blocks

Based on `docs/regulation/chapters/11-response-message-content/11-response-message-content.md`, the following data transfers are missing:

- **Activities** (TREPs `0x02`, `0x22`, `0x32`)
- **Events and Faults** (TREPs `0x03`, `0x23`, `0x33`)
- **Detailed Speed** (TREPs `0x04`, `0x24`)
- **Technical Data** (TREPs `0x05`, `0x25`, `0x35`)

### 4.2. Implementation Plan (per Block)

For each missing VU data block, the following steps should be taken:

1.  **Define Protobuf Message**: In `proto/wayplatform/connect/tachograph/vu/v1/file.proto`, define the messages for the missing data structures (e.g., `Activities`, `EventsAndFaults`, etc.). Again, `tachoparser` is an excellent reference.
2.  **Create Unmarshal File**: Create a new file, e.g., `unmarshal_vu_activities.go`.
3.  **Implement Unmarshal Function**: Implement the function (e.g., `UnmarshalVuActivities(r *bytes.Reader, target *vuv1.Activities, generation int)`). This function will read from the `bytes.Reader` and populate the target message. The function should be responsible for consuming the correct number of bytes for its block.
4.  **Integrate into `unmarshalVU`**: Add new `case` statements to the `switch` in `unmarshal.go:unmarshalVU` for the corresponding `TransferType` enums.

    ```go
    // Example for VU Activities in unmarshal.go
    case vuv1.TransferType_ACTIVITIES_GEN1, vuv1.TransferType_ACTIVITIES_GEN2_V1, vuv1.TransferType_ACTIVITIES_GEN2_V2:
        activities := &vuv1.Activities{}
        generation := 1 // Determine generation from transferType
        if transferType == vuv1.TransferType_ACTIVITIES_GEN2_V1 || transferType == vuv1.TransferType_ACTIVITIES_GEN2_V2 {
            generation = 2
        }
        _, err := UnmarshalVuActivities(r, activities, generation)
        if err != nil {
            return nil, err
        }
        transfer.SetActivities(activities)
    ```

## 5. General Considerations

- **Testing**: Each new unmarshalling function should be accompanied by unit tests with sample byte data to ensure correctness.
- **Incremental Development**: This work should be done incrementally. Implement one EF or VU block at a time to keep changes manageable and testable.
- **Reference Implementation**: The `benchmark/tachoparser` directory is an invaluable resource. Its `pkg/decoder/definitions.go` file contains Go structs that are very close to what the final protobuf messages should look like, and its `pkg/decoder/decoder.go` file shows the parsing logic.
