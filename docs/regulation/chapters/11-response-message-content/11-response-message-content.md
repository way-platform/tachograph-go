DDP\_028 The IDE shall wait at least for a period of P3 min before beginning each transmission; the wait period shall be measured from the last calculated occurrence of a stop bit after the error was detected.

#### Figure 3

#### IDE error handling

![](_page_0_Figure_4.jpeg)

2.2.6 *Response Message content*

This paragraph specifies the content of the data fields of the various positive response messages.

Data elements are defined in Appendix 1 data dictionary.

Remark: For generation 2 downloads, each top-level data element is represented by a record array, even if it contains only one record. A record array starts with a header; this header contains the record type, the record size and the number of records. Record arrays are named by '…RecordArray' (with header) in the following tables.

# M3

- 2.2.6.1 P o s i t i v e R e s p o n s e T r a n s f e r D a t a D o w n l o a d I n t e r f a c e V e r s i o n
  - DDP\_028a The data field of the 'Positive Response Transfer Data Download Interface Version' message shall provide the following data in the following order under the SID 76 Hex, the TREP 00 Hex:

# B

Data structure generation 2, version 2 (TREP 00 Hex)

Data element Comment

DownloadInterfaceVersion

Comment

DownloadInterfaceVersion Generation and version of the VU: 02,02 Hex for Generation 2, version 2. Not supported by Generation 1 and Generation 2, version 1 VU, which shall respond negatively (Sub function not supported, see DDP\_018)

### 2.2.6.2 P o s i t i v e R e s p o n s e T r a n s f e r D a t a O v e r v i e w

DDP\_029 The data field of the 'Positive Response Transfer Data Overview' message shall provide the following data in the following order under the SID 76 Hex, the TREP 01, 21 or 31 Hex and appropriate sub message splitting and counting:

Data structure generation 1 (TREP 01 Hex)

| Data element                      | Comment                                                                                                                                    |
|-----------------------------------|--------------------------------------------------------------------------------------------------------------------------------------------|
| MemberStateCertificate            | VU Security certificates                                                                                                                   |
| VUCertificate                     |                                                                                                                                            |
| VehicleIdentificationNumber       | Vehicle identification                                                                                                                     |
| VehicleRegistrationIdentification |                                                                                                                                            |
| CurrentDateTime                   | VU current date and time                                                                                                                   |
| VuDownloadablePeriod              | Downloadable period                                                                                                                        |
| CardSlotsStatus                   | Type of cards inserted in the VU                                                                                                           |
| VuDownloadActivityData            | Previous VU download                                                                                                                       |
| VuCompanyLocksData                | All company locks stored. If the section is empty, only noOfLocks = 0 is sent.                                                             |
| VuControlActivityData             | All control records stored in the VU. If the section is empty, only noOfControls = 0 is sent                                               |
| Signature                         | RSA signature of all data (except certificates) starting from VehicleIdentificationNumber down to last byte of last VuControlActivityData. |

MemberStateCertificateRecordArray Member state certificate

VUCertificateRecordArray VU certificate

VehicleIdentificationNumberRecordArray Vehicle identification

| Data element      | Comment                  |
|-------------------|--------------------------|
| icateRecordArray  | Member state certificate |
| rdArray           | VU certificate           |
| NumberRecordArray | Vehicle identification   |

Data structure generation 2, version 2 (TREP 31 Hex)

| Data element                           | Comment                                                                                              |
|----------------------------------------|------------------------------------------------------------------------------------------------------|
| MemberStateCertificateRecordArray      | Member state certificate                                                                             |
| VUCertificateRecordArray               | VU certificate                                                                                       |
| VehicleIdentificationNumberRecordArray | Vehicle identification                                                                               |
| VehicleRegistrationNumberRecordArray   | Vehicle registration number                                                                          |
| CurrentDateTimeRecordArray             | VU current date and time                                                                             |
| VuDownloadablePeriodRecordArray        | Downloadable period                                                                                  |
| CardSlotsStatusRecordArray             | Type of cards inserted in the VU                                                                     |
| VuDownloadActivityDataRecordArray      | Previous VU download                                                                                 |
| VuCompanyLocksRecordArray              | All company locks stored. If the section is<br>an array header with noOfRecords                      |
| VuControlActivityRecordArray           | All control records stored in the<br>section is empty, an array header. If no<br>Records = 0 is sent |

SignatureRecordArray

| Data element                                 | Comment                                                                                                            |
|----------------------------------------------|--------------------------------------------------------------------------------------------------------------------|
| VehicleRegistrationIdentificationRecordArray | Vehicle registration number                                                                                        |
| CurrentDateTimeRecordArray                   | VU current date and time                                                                                           |
| VuDownloadablePeriodRecordArray              | Downloadable period                                                                                                |
| CardSlotsStatusRecordArray                   | Type of cards inserted in the VU                                                                                   |
| VuDownloadActivityDataRecordArray            | Previous VU download                                                                                               |
| VuCompanyLocksRecordArray                    | All company locks stored. If the section is empty,<br>an array header with noOfRecords = 0 is sent                 |
| VuControlActivityRecordArray                 | All control records stored in the VU. If the<br>section is empty, an array header with noOf<br>Records = 0 is sent |
| SignatureRecordArray                         | ECC signature of all preceding data except the<br>certificates.                                                    |

| Data element                           | Comment                                                                                                            |
|----------------------------------------|--------------------------------------------------------------------------------------------------------------------|
| MemberStateCertificateRecordArray      | Member state certificate                                                                                           |
| VUCertificateRecordArray               | VU certificate                                                                                                     |
| VehicleIdentificationNumberRecordArray | Vehicle identification                                                                                             |
| VehicleRegistrationNumberRecordArray   | Vehicle registration number                                                                                        |
| CurrentDateTimeRecordArray             | VU current date and time                                                                                           |
| VuDownloadablePeriodRecordArray        | Downloadable period                                                                                                |
| CardSlotsStatusRecordArray             | Type of cards inserted in the VU                                                                                   |
| VuDownloadActivityDataRecordArray      | Previous VU download                                                                                               |
| VuCompanyLocksRecordArray              | All company locks stored. If the section is empty,<br>an array header with noOfRecords = 0 is sent                 |
| VuControlActivityRecordArray           | All control records stored in the VU. If the<br>section is empty, an array header with noOf<br>Records = 0 is sent |

SignatureRecordArray ECC signature of all preceding data except the certificates.

## 2.2.6.3 P o s i t i v e R e s p o n s e T r a n s f e r D a t a A c t i v i t i e s

DDP\_030 The data field of the 'Positive Response Transfer Data Activities' message shall provide the following data in the following order under the SID 76 Hex, the TREP 02, 22 or 32 Hex and appropriate sub message splitting and counting:

Data structure generation 1 (TREP 02 Hex)

| Data element               | Comment                                                                                                                                                                                                                                                                                                                            |
|----------------------------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| TimeReal                   | Date of day downloaded                                                                                                                                                                                                                                                                                                             |
| OdometerValueMidnight      | Odometer at end of downloaded day                                                                                                                                                                                                                                                                                                  |
| VuCardIWData               | Cards insertion withdrawal cycles data.<br>— If this section contains no available data, only<br>noOfVuCardIWRecords = 0 is sent.<br>— When a VuCardIWRecord lies across 00:00 (card<br>insertion on previous day) or across 24:00 (card<br>withdrawal the following day) it shall appear in<br>full within the two days involved. |
| VuActivityDailyData        | Slots status at 00:00 and activity changes recorded for<br>the day downloaded.                                                                                                                                                                                                                                                     |
| VuPlaceDailyWorkPeriodData | Places related data recorded for the day downloaded. If<br>the section is empty, only noOfPlaceRecords = 0 is<br>sent.                                                                                                                                                                                                             |
| VuSpecificConditionData    | Specific conditions data recorded for the day down<br>loaded. If the section is empty, only noOfSpecificCon<br>ditionRecords=0 is sent                                                                                                                                                                                             |
| Signature                  | RSA signature of all data starting from TimeReal down<br>to last byte of last specific condition record.                                                                                                                                                                                                                           |

Data structure generation 2, version 1 (TREP 22 Hex)

| Data element                      | Comment                                                                                                                                                                                                                                                                                                                                    |
|-----------------------------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| DateOfDayDownloadedRecordArray    | Date of day downloaded                                                                                                                                                                                                                                                                                                                     |
| OdometerValueMidnightRecordArray  | Odometer at end of downloaded day                                                                                                                                                                                                                                                                                                          |
| VuCardIWRecordArray               | Cards insertion withdrawal cycles data.<br>— If this section contains no available data, an array<br>header with noOfRecords = 0 is sent.<br>— When a VuCardIWRecord lies across 00:00 (card<br>insertion on previous day) or across 24:00 (card<br>withdrawal the following day) it shall appear in<br>full within the two days involved. |
| VuActivityDailyRecordArray        | Slots status at 00:00 and activity changes recorded for<br>the day downloaded.                                                                                                                                                                                                                                                             |
| VuPlaceDailyWorkPeriodRecordArray | Places related data recorded for the day downloaded. If<br>the section is empty, an array header with noOfRecords<br>= 0 is sent.                                                                                                                                                                                                          |

VuSpecificConditionRecordArray

SignatureRecordArray

Data element Comment

VuGNSSADRecordArray GNSS positions of the vehicle if the accumulated driving time of the vehicle reaches a multiple of three hours. If the section is empty, an array header with noOfRecords = 0 is sent.

VuSpecificConditionRecordArray Specific conditions data recorded for the day downloaded. If the section is empty, an array header with noOfRecords =0 is sent

SignatureRecordArray ECC signature of all preceding data.

Data structure generation 2, version 2 (TREP 32 Hex)

| Data element                      | Comment                                                                                                                                                                                                                                                                                                                                    |
|-----------------------------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| DateOfDayDownloadedRecordArray    | Date of day downloaded                                                                                                                                                                                                                                                                                                                     |
| OdometerValueMidnightRecordArray  | Odometer at end of downloaded day                                                                                                                                                                                                                                                                                                          |
| VuCardIWRecordArray               | Cards insertion withdrawal cycles data.<br>— If this section contains no available data, an array<br>header with noOfRecords = 0 is sent.<br>— When a VuCardIWRecord lies across 00:00 (card<br>insertion on previous day) or across 24:00 (card<br>withdrawal the following day) it shall appear in<br>full within the two days involved. |
| VuActivityDailyRecordArray        | Slots status at 00:00 and activity changes recorded for<br>the day downloaded.                                                                                                                                                                                                                                                             |
| VuPlaceDailyWorkPeriodRecordArray | Places related data recorded for the day downloaded. If<br>the section is empty, an array header with noOfRecords<br>= 0 is sent.                                                                                                                                                                                                          |
| VuGNSSADRecordArray               | GNSS positions of the vehicle if the accumulated<br>driving time of the vehicle reaches a multiple of three<br>hours. If the section is empty, an array header with<br>noOfRecords = 0 is sent.                                                                                                                                            |
| VuSpecificConditionRecordArray    | Specific conditions data recorded for the day down<br>loaded. If the section is empty, an array header with<br>noOfRecords =0 is sent                                                                                                                                                                                                      |
| VuBorderCrossingRecordArray       | Border crossings for the day downloaded. If the<br>section is empty, an array header with noOfRecords =<br>0 is sent.                                                                                                                                                                                                                      |
| VuLoadUnloadRecordArray           | Load/unload operations for the day downloaded. If the<br>section is empty, an array header with noOfRecords = 0<br>is sent.                                                                                                                                                                                                                |
| SignatureRecordArray              | ECC signature of all preceding data.                                                                                                                                                                                                                                                                                                       |

#### 2.2.6.4 P o s i t i v e R e s p o n s e T r a n s f e r D a t a E v e n t s a n d F a u l t s

DDP\_031 The data field of the 'Positive Response Transfer Data Events and Faults' message shall provide the following data in the following order under the SID 76 Hex, the TREP 03, 23 or 33 Hex and appropriate sub message splitting and counting:

Data structure generation 1, (TREP 03 Hex)

| Data element              | Comment                                                                                                                                                         |
|---------------------------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------|
| VuFaultData               | All faults stored or on-going in the VU.<br>If the section is empty, only noOfVuFaults = 0 is sent.                                                             |
| VuEventData               | All events (except over speeding) stored or on-going in<br>the VU.<br>If the section is empty, only noOfVuEvents = 0 is sent.                                   |
| VuOverSpeedingControlData | Data related to last over speeding control (default value<br>if no data).                                                                                       |
| VuOverSpeedingEventData   | All over speeding events stored in the VU.<br>If the section is empty, only noOfVuOverSpeed<br>ingEvents = 0 is sent.                                           |
| VuTimeAdjustmentData      | All time adjustment events stored in the VU (outside the<br>frame of a full calibration).<br>If the section is empty, only noOfVuTimeAdjRecords =<br>0 is sent. |
| Signature                 | RSA signature of all data starting from noOfVuFaults<br>down to last byte of last time adjustment record                                                        |

#### Data structure generation 2, version 1 (TREP 23 Hex)

| Data element                         | Comment                                                                                                                                          |
|--------------------------------------|--------------------------------------------------------------------------------------------------------------------------------------------------|
| VuFaultRecordArray                   | All faults stored or on-going in the VU.<br>If the section is empty, an array header with noOf<br>Records = 0 is sent.                           |
| VuEventRecordArray                   | All events (except over speeding) stored or on-going in<br>the VU.<br>If the section is empty, an array header with noOf<br>Records = 0 is sent. |
| VuOverSpeedingControlDataRecordArray | Data related to last over speeding control (default value<br>if no data).                                                                        |
| VuOverSpeedingEventRecordArray       | All over speeding events stored in the VU.<br>If the section is empty, an array header with noOf<br>Records = 0 is sent.                         |

#### **▼M3**

| Data element                | Comment                                                                                                                                                                 |
|-----------------------------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| VuTimeAdjustmentRecordArray | All time adjustment events stored in the VU (outside the<br>frame of a full calibration).<br>If the section is empty, an array header with noOf<br>Records = 0 is sent. |
| SignatureRecordArray        | ECC signature of all preceding data.                                                                                                                                    |

Data structure generation 2, version 2 (TREP 33 Hex)

| Data element                         | Comment                                                                                                                                                                 |
|--------------------------------------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| VuFaultRecordArray                   | All faults stored or on-going in the VU.<br>If the section is empty, an array header with noOf<br>Records = 0 is sent.                                                  |
| VuEventRecordArray                   | All events (except over speeding) stored or on-going in<br>the VU.<br>If the section is empty, an array header with noOf<br>Records = 0 is sent.                        |
| VuOverSpeedingControlDataRecordArray | Data related to last over speeding control (default value<br>if no data).                                                                                               |
| VuOverSpeedingEventRecordArray       | All over speeding events stored in the VU.<br>If the section is empty, an array header with noOf<br>Records = 0 is sent.                                                |
| VuTimeAdjustmentRecordArray          | All time adjustment events stored in the VU (outside the<br>frame of a full calibration).<br>If the section is empty, an array header with noOf<br>Records = 0 is sent. |
| SignatureRecordArray                 | ECC signature of all preceding data.                                                                                                                                    |

## 2.2.6.5 P o s i t i v e R e s p o n s e T r a n s f e r D a t a D e t a i l e d S p e e d

DDP\_032 The data field of the 'Positive Response Transfer Data Detailed Speed' message shall provide the following data in the following order under the SID 76 Hex, the TREP 04 or 24 Hex and appropriate sub message splitting and counting:

Data structure generation 1 (TREP 04 Hex)

| Data element        | Comment                                                                                                                                                |
|---------------------|--------------------------------------------------------------------------------------------------------------------------------------------------------|
| VuDetailedSpeedData | All detailed speed stored in the VU (one speed block per minute during which the vehicle has been moving) 60 speed values per minute (one per second). |
| Signature           | RSA signature of all data starting from noOf SpeedBlocks down to last byte of last speed block.                                                        |

## M3

Data structure generation 2 (TREP 24 Hex)

Data elem

VuDetailedSpeedBlockRecordArray

Data element Comment

VuDetailedSpeedBlockRecordArray All detailed speed stored in the VU (one speed block per minute during which the vehicle has been moving) 60 speed values per minute (one per second).

SignatureRecordArray

SignatureRecordArray ECC signature of all preceding data.

#### 2.2.6.6 P o s i t i v e R e s p o n s e T r a n s f e r D a t a T e c h n i c a l D a t a

DDP\_033 The data field of the 'Positive Response Transfer Data Technical Data' message shall provide the following data in the following order under the SID 76 Hex, the TREP 05, 25 or 35 Hex and appropriate sub message splitting and counting:

Data structure generation 1 (TREP 05 Hex)

| Data element      | Comment                                                                                                    |
|-------------------|------------------------------------------------------------------------------------------------------------|
| VuIdentification  |                                                                                                            |
| SensorPaired      |                                                                                                            |
| VuCalibrationData | All calibration records stored in the VU.                                                                  |
| Signature         | RSA signature of all data starting from vuManufacturerName down to last byte of last VuCalibration-Record. |

Data structure generation 2, version 1 (TREP 25 Hex)

| Data element                           | Comment                                               |
|----------------------------------------|-------------------------------------------------------|
| VuIdentificationRecordArray            |                                                       |
| VuSensorPairedRecordArray              | All MS pairings stored in the VU                      |
| VuSensorExternalGNSSCoupledRecordArray | All external GNSS facility couplings stored in the VU |
| VuCalibrationRecordArray               | All calibration records stored in the VU.             |
| VuCardRecordArray                      | All card insertion data stored in the VU.             |
| VuITSConsentRecordArray                |                                                       |
| VuPowerSupplyInterruptionRecordArray   |                                                       |
| SignatureRecordArray                   | ECC signature of all preceding data.                  |

Data structure generation 2, version 2 (TREP 35 Hex)

![](_page_8_Figure_3.jpeg)

# B

## 2.3. ESM File storage

DDP\_034 When a download session has included a VU data transfer, the IDE shall store within one single physical file all data received from the VU during the download session within Positive Response Transfer Data messages. Data stored excludes message headers, sub-message counters, empty sub-messages and checksums but include the SID and TREP (of the first sub-message only if several submessages).

### 3. TACHOGRAPH CARDS DOWNLOADING PROTOCOL

### 3.1. Scope

This paragraph describes the direct card data downloading of a tachograph card to an IDE. The IDE is not part of the secure environment; therefore no authentication between the card and the IDE is performed.

## 3.2. **Definitions**

**Download session**: Each time a download of the ICC data is performed. The session covers the complete procedure from the reset of the ICC by an IFD until the deactivation of the ICC (withdraw of the card or next reset).

**Signed Data File**: A file from the ICC. The file is transferred to the IFD in plain text. On the ICC the file is hashed and signed and the signature is transferred to the IFD.