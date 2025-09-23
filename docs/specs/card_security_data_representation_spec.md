# Specification: Representing Certificates and Signatures in Parsed Card Data

## 1. Introduction

This document specifies necessary extensions to the tachograph card's protobuf data model. The current implementation parses Elementary Files (EFs) into structured messages (e.g., `DriverCardFile`), but it omits critical security data: the certificates required for authentication and the digital signatures that ensure data integrity.

The goal of this specification is to extend the data model to include the raw data of all certificate EFs and the digital signatures for all signed EFs, ensuring that a fully parsed file is complete and verifiable.

## 2. Problem Analysis

According to EU regulations (Appendix 7 and 12), a downloaded card file is a concatenation of Tag-Length-Value (TLV) records. This stream contains:

1.  **Unsigned Certificate EFs**: Files like `EF_Card_Certificate`, `EF_CA_Certificate`, and `EF_Link_Certificate` are present as distinct TLV records.
2.  **Signed Data EFs**: Most application data files (e.g., `EF_Events_Data`, `EF_Vehicles_Used`) are represented by a data TLV record immediately followed by a signature TLV record.

The current high-level protobuf messages (`DriverCardFile`, `WorkshopCardFile`, etc.) only have fields for the application data itself. They lack fields to store the raw bytes of the certificate files and the signatures, leading to a loss of information during the unmarshalling process from a `RawCardFile`.

## 3. Proposed Solution

To create a complete and verifiable representation of a card file, we will introduce new messages and fields to hold the security data.

### 3.1. Certificate Representation

We will define a `Certificates` message to hold the raw data of all certificate files present on the card. This message will then be embedded in each top-level card file message.

**New Message Definition (`card/v1/certificates.proto`):**
```protobuf
syntax = "proto3";

package wayplatform.connect.tachograph.card.v1;

// Certificates holds the raw data of all certificate Elementary Files
// downloaded from a tachograph card.
message Certificates {
  // Raw data of the card's public key certificate.
  // Gen1: EF_Card_Certificate (FID 0xC100)
  // Gen2: EF_CardMA_Certificate (FID 0xC100)
  bytes card_certificate = 1;

  // Raw data of the Certification Authority's public key certificate.
  // EF_CA_Certificate (FID 0xC108)
  bytes ca_certificate = 2;

  // Raw data of the card's public key certificate for digital signatures (Gen2+).
  // EF_CardSignCertificate (FID 0xC101)
  bytes card_sign_certificate = 3;

  // Raw data of the link certificate for chaining root CAs (Gen2+).
  // EF_Link_Certificate (FID 0xC109)
  bytes link_certificate = 4;
}
```

This new `Certificates` message will be added to `DriverCardFile`, `WorkshopCardFile`, `ControlCardFile`, and `CompanyCardFile`.

**Example (`DriverCardFile.proto`):**
```protobuf
// ... imports
import "wayplatform/connect/tachograph/card/v1/certificates.proto";

message DriverCardFile {
  // ... existing fields
  LastCardDownload last_card_download = 15;
  Certificates certificates = 16; // New field
  // ... Gen2 specific fields
}
```

### 3.2. Signature Representation

For every message representing a signable EF, we will add a `bytes signature` field. This field will store the raw signature data that follows the EF data in the downloaded file.

**Example (`event_data.proto`):**
```protobuf
message EventData {
  repeated Record records = 1;
  // The digital signature for the entire EF_Events_Data file content.
  bytes signature = 2;
}
```

## 4. Detailed Schema Changes

The following changes are required to implement this specification.

1.  **Create `proto/wayplatform/connect/tachograph/card/v1/certificates.proto`** with the `Certificates` message defined in section 3.1.

2.  **Add the `Certificates` field** to the following messages:
    - `proto/wayplatform/connect/tachograph/card/v1/driver_card_file.proto`
    - `proto/wayplatform/connect/tachograph/card/v1/workshop_card_file.proto`
    - `proto/wayplatform/connect/tachograph/card/v1/control_card_file.proto`
    - `proto/wayplatform/connect/tachograph/card/v1/company_card_file.proto`

3.  **Add a `bytes signature = N;` field** to the protobuf messages corresponding to all signable EFs. This includes, but is not limited to:
    - `ApplicationIdentificationV2`
    - `BorderCrossings`
    - `Calibrations`
    - `CalibrationsAddData`
    - `CompanyActivityData`
    - `CompanyApplicationIdentification`
    - `ControlActivityData`
    - `ControllerActivityData`
    - `CurrentUsage`
    - `DriverActivity`
    - `DriverCardApplicationIdentification`
    - `DrivingLicenceInfo`
    - `EventData`
    - `FaultData`
    - `GnssPlaces`
    - `GnssPlacesAuthentication`
    - `CardIdentification` & `*HolderIdentification` (The `EF_Identification` file is signed as a whole)
    - `LoadTypeEntries`
    - `LoadUnloadOperations`
    - `Places`
    - `PlacesAuthentication`
    - `SpecificConditions`
    - `VehicleUnitsUsed`
    - `VehiclesUsed`
    - `VuConfiguration`
    - `WorkshopApplicationIdentification`

    *Note*: Since `EF_Identification` is signed as a single block but parsed into two separate messages (`CardIdentification` and a holder-specific identification), the `signature` field should be added to the top-level card file message (e.g., `DriverCardFile`) specifically for this EF.

    **Revised Proposal for `EF_Identification` signature:**
    Add `bytes identification_signature = N;` to `DriverCardFile`, `WorkshopCardFile`, etc.

## 5. Implementation Guidance

The unmarshalling logic in `unmarshal.go` must be updated to populate these new fields.

-   When iterating through the TLV records in `unmarshalCard`, `case` statements for the certificate `ElementaryFileType`s (`EF_CARD_CERTIFICATE`, `EF_CA_CERTIFICATE`, etc.) should be added. These cases will populate the appropriate field in the new `Certificates` message.
-   The logic should be adapted to handle signature records. When a data EF is unmarshalled, the system should check if the next TLV record is its corresponding signature (e.g., `EF_EVENTS_DATA` `0x0502` is followed by its signature `0x0502` with the signature bit set in the appendix). If so, the signature value should be stored in the `signature` field of the just-parsed data message (e.g., `EventData.signature`).
