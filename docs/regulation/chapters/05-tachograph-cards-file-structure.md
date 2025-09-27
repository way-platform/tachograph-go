### 4.1. **Master File MF**

TCS\_142 After its personalisation, the master file MF shall have the following permanent file structure and file access rules:

> *Note:* The short EF identifier SFID is given as decimal number, e.g. the value 30 corresponds to 11110 in binary.

| File                             | File ID | SFID | Read / Select | Update |
|----------------------------------|---------|------|---------------|--------|
| MF                               | '3F00h' |      |               |        |
| EF ICC                           | '0002h' |      | ALW           | NEV    |
| EF IC                            | '0005h' |      | ALW           | NEV    |
| EF DIR                           | '2F00h' | 30   | ALW           | NEV    |
| EF ATR/INFO (conditional)        | '2F01h' | 29   | ALW           | NEV    |
| EF Extended_Length (conditional) | '0006h' | 28   | ALW           | NEV    |
| — DF Tachograph                  | '0500h' |      | SC1           |        |
| DF Tachograph G2                 |         |      | SC1           |        |

The following abbreviation for the security condition is used in this table:

### **SC1** ALW OR SM-MAC-G2

### TCS\_143 All EF structures shall be transparent.

TCS\_144 The Master File MF shall have the following data structure:

| File / Data element        | No of<br>Records | Size (bytes) |     | Default<br>Values |
|----------------------------|------------------|--------------|-----|-------------------|
| MF                         |                  | Min          | Max |                   |
| EF ICC                     |                  | 63           | 184 |                   |
| └CardIccIdentification     |                  | 25           | 25  |                   |
| └clockStop                 |                  | 1            | 1   | {00}              |
| └cardExtendedSerialNumber  |                  | 8            | 8   | {00..00}          |
| └cardApprovalNumber        |                  | 8            | 8   | {20..20}          |
| └cardPersonaliserID        |                  | 1            | 1   | {00}              |
| └embedderIcAssemblerId     |                  | 5            | 5   | {00..00}          |
| └icIdentifier              |                  | 2            | 2   | {00 00}           |
| EF IC                      |                  | 8            | 8   |                   |
| └CardChipIdentification    |                  | 8            | 8   |                   |
| └icSerialNumber            |                  | 4            | 4   | {00..00}          |
| └icManufacturingReferences |                  | 4            | 4   | {00..00}          |
| EF DIR                     |                  | 20           | 20  |                   |
| └See TCS_145               |                  | 20           | 20  | {00..00}          |
| EF ATR/INFO                |                  | 7            | 128 |                   |
| └See TCS_146               |                  | 7            | 128 | {00..00}          |
| EF EXTENDED LENGTH         |                  | 3            | 3   |                   |
| └See TCS_147               |                  | 3            | 3   | {00..00}          |
| DF Tachograph              |                  |              |     |                   |
| DF Tachograph G2           |                  |              |     |                   |

#### TCS\_145 The elementary file EF DIR shall contain the following application related data objects: '61 08 4F 06 FF 54 41 43 48 4F 61 08 4F 06 FF 53 4D 52 44 54'

- TCS\_146 The elementary file EF ATR/INFO shall be present if the tachograph card indicates in its ATR that it supports extended length fields. In this case the EF ATR/INFO shall contain the extended length information data object (DO'7F66') as specified in ISO/IEC 7816-4:2013 clause 12.7.1.
- TCS\_147 The elementary file EF Extended\_Length shall be present if the tachograph card indicates in its ATR that it supports extended length fields. In this case the EF shall contain the following data object: '02 01 xx' where the value 'xx' indicates whether extended length fields are supported for the T = 1 and / or T = 0 protocol.

The value '01' indicates extended length field support for the T = 1 protocol.

The value '10' indicates extended length field support for the T = 0 protocol.

The value '11' indicates extended length field support for the T = 1 and the T = 0 protocol.

#### 4.2. **Driver card applications**

4.2.1 *Driver card application generation 1*

TCS\_148 After its personalisation, the driver card application generation 1 shall have the following permanent file structure and file access rules:

| File                          | File ID | Read | Select | Update |
|-------------------------------|---------|------|--------|--------|
| DF Tachograph                 | '0500h' | SC2  | SC1    | NEV    |
| EF Application_Identification | '0501h' | SC2  | SC1    | NEV    |
| EF Card_Certificate           | 'C100h' | SC2  | SC1    | NEV    |
| EF CA_Certificate             | 'C108h' | SC2  | SC1    | NEV    |
| EF Identification             | '0520h' | SC2  | SC1    | NEV    |
| EF Card_Download              | '050Eh' | SC2  | SC1    | SC1    |
| EF Driving_Licence_Info       | '0521h' | SC2  | SC1    | NEV    |
| EF Events_Data                | '0502h' | SC2  | SC1    | SC3    |
| EF Faults_Data                | '0503h' | SC2  | SC1    | SC3    |
| EF Driver_Activity_Data       | '0504h' | SC2  | SC1    | SC3    |
| EF Vehicles_Used              | '0505h' | SC2  | SC1    | SC3    |
| EF Places                     | '0506h' | SC2  | SC1    | SC3    |
| EF Current_Usage              | '0507h' | SC2  | SC1    | SC3    |
| EF Control_Activity_Data      | '0508h' | SC2  | SC1    | SC3    |
| EF Specific Conditions        | '0522h' | SC2  | SC1    | SC3    |

The following abbreviations for the security conditions are used in this table:

- **SC1** ALW OR SM-MAC-G2
- **SC2** ALW OR SM-MAC-G1 OR SM-MAC-G2
- **SC3** SM-MAC-G1 OR SM-MAC-G2
- TCS\_149 All EF structures shall be transparent.
- TCS\_150 The driver card application generation 1 shall have the following data structure:

| File / Data element                  | No of<br>Records | Size (bytes)<br>Min | Size (bytes)<br>Max | Default<br>Values |
|--------------------------------------|------------------|---------------------|---------------------|-------------------|
| DF Tachograph                        |                  | 11378               | 24926               |                   |
| EF Application Identification        |                  |                     |                     |                   |
| └DriverCardApplicationIdentification |                  | 10                  | 10                  |                   |
| └typeOfTachographCardId              |                  | 1                   | 1                   | {00}              |
| └cardStructureVersion                |                  | 2                   | 2                   | {00 00}           |
| └noOfEventsPerType                   |                  | 1                   | 1                   | {00}              |
| └noOfFaultsPerType                   |                  | 1                   | 1                   | {00}              |
| └activityStructureLength             |                  | 2                   | 2                   | {00 00}           |
| └noofCardVehicleRecords              |                  | 2                   | 2                   | {00 00}           |
| └noOfCardPlaceRecords                |                  | 1                   | 1                   | {00}              |
| EF Card Certificate                  |                  |                     |                     |                   |
| └CardCertificate                     |                  | 194                 | 194                 | {00..00}          |
| EF CA Certificate                    |                  |                     |                     |                   |
| └MemberStateCertificate              |                  | 194                 | 194                 | {00..00}          |
| EF Identification                    |                  | 143                 | 143                 |                   |
| └CardIdentification                  |                  | 65                  | 65                  |                   |
| └cardIssuingMemberState              |                  | 1                   | 1                   | {00}              |
| └cardNumber                          |                  | 16                  | 16                  | {20..20}          |
| └cardIssuingAuthorityName            |                  | 36                  | 36                  | {00, 20..2        |
| └cardIssueDate                       |                  | 4                   | 4                   | {00..00}          |
| └cardValidityBegin                   |                  | 4                   | 4                   | {00..00}          |
| └cardExpiryDate                      |                  | 4                   | 4                   | {00..00}          |
| └DriverCardHolderIdentification      |                  | 78                  | 78                  |                   |
| └cardHolderName                      |                  | 72                  | 72                  |                   |
| └holderSurname                       |                  | 36                  | 36                  | {00, 20..2        |
| └holderFirstNames                    |                  | 36                  | 36                  | {00, 20..2        |
| └cardHolderBirthDate                 |                  | 4                   | 4                   | {00..00}          |
| └cardHolderPreferredLanguage         |                  | 2                   | 2                   | {20 20}           |
| EF Card Download                     |                  | 4                   | 4                   |                   |
| ▶└LastCardDownload                   |                  | 4                   | 4                   | {00..00}          |
| EF Driving_Licence_Info              |                  | 53                  | 53                  |                   |
| └CardDrivingLicenceInformation       |                  | 53                  | 53                  |                   |
| └drivingLicenceIssuingAuthority      |                  | 36                  | 36                  | {00, 20..20       |
| └drivingLicenceIssuingNation         |                  | 1                   | 1                   | {00}              |
| └drivingLicenceNumber                |                  | 16                  | 16                  | {20..20}          |
| EF Events Data                       |                  | 864                 | 1728                |                   |
| └CardEventData                       |                  | 864                 | 1728                |                   |
| └cardEventRecords                    | 6                | 144                 | 288                 |                   |
| └ CardEventRecord                    | n1               | 24                  | 24                  |                   |
| └eventType                           |                  | 1                   | 1                   | {00}              |
| └eventBeginTime                      |                  | 4                   | 4                   | {00..00}          |
| └eventEndTime                        |                  | 4                   | 4                   | {00..00}          |
| └eventVehicleRegistration            |                  |                     |                     |                   |
| └vehicleRegistrationNation           |                  | 1                   | 1                   | {00}              |
| └vehicleRegistrationNumber           |                  | 14                  | 14                  | {00, 20..2        |
| EF Faults Data                       |                  | 576                 | 1152                |                   |
| └CardFaultData                       |                  | 576                 | 1152                |                   |
| └cardFaultRecords                    | 2                | 288                 | 576                 |                   |
| └ CardFaultRecord                    | n2               | 24                  | 24                  |                   |
| └faultType                           |                  | 1                   | 1                   | {00}              |
| └faultBeginTime                      |                  | 4                   | 4                   | {00..00}          |

**►**(1) (2) **M3**

| ▼ | B |
|---|---|
|   | _ |

|                                | -vehicleRegistrationNation      | 1    | 1     | {00}         |              |
|--------------------------------|---------------------------------|------|-------|--------------|--------------|
|                                | -vehicleRegistrationNumber      | 14   | 14    | {00, 20..20} |              |
| EF Driver Activity Data        |                                 | 5548 | 13780 |              |              |
| └CardDriverActivity            |                                 | 5548 | 13780 |              |              |
|                                | -activityPointerOldestDayRecord | 2    | 2     | {00 00}      |              |
|                                | -activityPointerNewestRecord    | 2    | 2     | {00 00}      |              |
|                                | -activityDailyRecords           | n6   | 5544  | 13776        | {00..00}     |
| EF Vehicles Used               |                                 | 2606 | 6202  |              |              |
| └CardVehiclesUsed              |                                 | 2606 | 6202  |              |              |
|                                | -vehiclePointerNewestRecord     | 2    | 2     | {00 00}      |              |
|                                | -cardVehicleRecords             |      | 2604  | 6200         |              |
|                                | └CardVehicleRecord              | n3   | 31    | 31           |              |
|                                | -vehicleOdometerBegin           |      | 3     | 3            | {00..00}     |
|                                | -vehicleOdometerEnd             |      | 3     | 3            | {00..00}     |
|                                | -vehicleFirstUse                |      | 4     | 4            | {00..00}     |
|                                | -vehicleLastUse                 |      | 4     | 4            | {00..00}     |
|                                | -vehicleRegistration            |      |       |              |              |
|                                | -vehicleRegistrationNation      |      | 1     | 1            | {00}         |
|                                | -vehicleRegistrationNumber      |      | 14    | 14           | {00, 20..20} |
|                                | -vuDataBlockCounter             |      | 2     | 2            | {00 00}      |
| EF Places                      |                                 |      | 841   | 1121         |              |
| └CardPlaceDailyWorkPeriod      |                                 |      | 841   | 1121         |              |
|                                | -placePointerNewestRecord       |      | 1     | 1            | {00}         |
|                                | -placeRecords                   |      | 840   | 1120         |              |
|                                | └PlaceRecord                    | n4   | 10    | 10           |              |
|                                | -entryTime                      |      | 4     | 4            | {00..00}     |
|                                | -entryTypeDailyWorkPeriod       |      | 1     | 1            | {00}         |
|                                | -dailyWorkPeriodCountry         |      | 1     | 1            | {00}         |
|                                | -dailyWorkPeriodRegion          |      | 1     | 1            | {00}         |
|                                | -vehicleOdometerValue           |      | 3     | 3            | {00..00}     |
| EF Current Usage               |                                 |      | 19    | 19           |              |
| └CardCurrentUse                |                                 |      | 19    | 19           |              |
|                                | -sessionOpenTime                |      | 4     | 4            | {00..00}     |
|                                | -sessionOpenVehicle             |      |       |              |              |
|                                | -vehicleRegistrationNation      |      | 1     | 1            | {00}         |
|                                | -vehicleRegistrationNumber      |      | 14    | 14           | {00, 20..20} |
| EF Control Activity Data       |                                 |      | 46    | 46           |              |
| └CardControlActivityDataRecord |                                 |      | 46    | 46           |              |
|                                | -controlType                    |      | 1     | 1            | {00}         |
|                                | -controlTime                    |      | 4     | 4            | {00..00}     |
|                                | -controlCardNumber              |      |       |              |              |
|                                | cardType                        |      | 1     | 1            | {00}         |
|                                | -cardIssuingMemberState         |      | 1     | 1            | {00}         |
|                                | -cardNumber                     |      | 16    | 16           | {20..20}     |
|                                | -controlVehicleRegistration     |      |       |              |              |
|                                | -vehicleRegistrationNation      |      | 1     | 1            | {00}         |
|                                | -vehicleRegistrationNumber      |      | 14    | 14           | {00, 20..20} |
|                                | -controlDownloadPeriodBegin     |      | 4     | 4            | {00..00}     |
|                                | -controlDownloadPeriodEnd       |      | 4     | 4            | {00..00}     |
| EF Specific_Conditions         |                                 |      | 280   | 280          |              |
| └SpecificConditionRecord       |                                 | 56   | 5     | 5            |              |
|                                | -entryTime                      |      | 4     | 4            | {00..00}     |
|                                | -SpecificConditionType          |      | 1     | 1            | {00}         |

TCS\_151 The following values, used to provide sizes in the table above, are the minimum and maximum record number values the driver card data structure must use for a generation 1 application:

|       |                         | Min                                            | Max                                              |
|-------|-------------------------|------------------------------------------------|--------------------------------------------------|
| $n_1$ | NoOfEventsPerType       | 6                                              | 12                                               |
| $n_2$ | NoOfFaultsPerType       | 12                                             | 24                                               |
| $n_3$ | NoOfCardVehicleRecords  | 84                                             | 200                                              |
| $n_4$ | NoOfCardPlaceRecords    | 84                                             | 112                                              |
| $n_6$ | CardActivityLengthRange | 5 544 bytes<br>(28 days * 93 activity changes) | 13 776 Bytes<br>(28 days * 240 activity changes) |

#### 4.2.2 *Driver card application generation 2*

### **▼M3**

TCS\_152 After its personalisation, the driver card application generation 2 shall have the following permanent file structure and file access rules:

#### *Notes:*

- The short EF identifier SFID is given as decimal number, e.g. the value 30 corresponds to 11110 in binary.
- EF Application\_Identification\_V2, EF Places\_Authentication, EF GNSS\_Places\_Authentication, EF Border\_Crossings, EF Load\_Unload\_Operations, EF VU\_Configuration and EF Load\_Type\_Entries are only present in version 2 of the generation 2 driver card.
- cardStructureVersion in EF Application\_Identification is equal to {01 01} for version 2 of the generation 2 driver card, while it was equal to {01 00} for version 1 of the generation 2 driver card.

| File                            | File ID | SFID | Access rules<br>Read / Select | Update   |
|---------------------------------|---------|------|-------------------------------|----------|
| └─DF Tachograph G2              |         |      | SC1                           |          |
| └─EF Application Identification | '0501h  | 1    | SC1                           | NEV      |
| └─EF CardMA Certificate         | 'C100h  | 2    | SC1                           | NEV      |
| └─EF CardSignCertificate        | 'C101h  | 3    | SC1                           | NEV      |
| └─EF CA Certificate             | 'C108h  | 4    | SC1                           | NEV      |
| └─EF Link Certificate           | 'C109h  | 5    | SC1                           | NEV      |
| └─EF Identification             | '0520h  | 6    | SC1                           | NEV      |
| └─EF Card Download              | '050Eh  | 7    | SC1                           | SC1      |
| └─EF Driving Licence Info       | '0521h  | 10   | SC1                           | NEV      |
| └─EF Events Data                | '0502h  | 12   | SC1                           | SM-MAC-G |
| └─EF Faults Data                | '0503h  | 13   | SC1                           | SM-MAC-G |
| └─EF Driver Activity Data       | '0504h  | 14   | SC1                           | SM-MAC-G |
| └─EF Vehicles Used              | '0505h  | 15   | SC1                           | SM-MAC-G |
| └─EF Places                     | '0506h  | 16   | SC1                           | SM-MAC-G |
| └─EF Current Usage              | '0507h  | 17   | SC1                           | SM-MAC-G |
| └─EF Control Activity Data      | '0508h  | 18   | SC1                           | SM-MAC-G |
| └─EF Specific Conditions        | '0522h  | 19   | SC1                           | SM-MAC-G |
| └─EF VehicleUnits Used          | '0523h  | 20   | SC1                           | SM-MAC-G |
| └─EF GNSS Places                | '0524h  | 21   | SC1                           | SM-MAC-G |
| └─EF Application Identification | '0525h  | 22   | SC1                           | NEV      |
| └─EF Places Authentication      | '0526h  | 23   | SC1                           | SM-MAC-G |
| └─EF GNSS Places Authentication | '0527h  | 24   | SC1                           | SM-MAC-G |
| └─EF Border Crossings           | '0528h  | 25   | SC1                           | SM-MAC-G |
| └─EF Load Unload Operations     | '0529h  | 26   | SC1                           | SM-MAC-G |
| └─EF Load Type Entries          | '0530h  | 27   | SC1                           | SM-MAC-G |
| EF Vu Configuration             | '0540h  | 30   | SC5/SC1                       | SM-MAC-G |

The following abbreviations for the security condition are used in this table:

**SC1** ALW OR SM-MAC-G2

**SC5** For the Read Binary command with even INS byte: SM-C-MAC-G2 AND SM-R-ENC-MAC-G2

> For the Read Binary command with odd INS byte (if supported): NEV

**▼B**

TCS\_153 All EF structures shall be transparent.

**▼M3**

TCS\_154 The driver card application generation 2 shall have the following data structure:

| File Element / Data                 | No of<br>Records | Size<br>(bytes) |      | Default Values |   |   |   |          |          |
|-------------------------------------|------------------|-----------------|------|----------------|---|---|---|----------|----------|
|                                     |                  | Min             | Max  |                |   |   |   |          |          |
| DF Tachograph_G2                    | 9830             | 988             |      |                |   |   |   |          |          |
| EF Application_Identification       | 0                |                 | 48   |                |   |   |   |          |          |
| DriverCardApplicationIdentification |                  | 17              | 17   |                |   |   |   |          |          |
| typeOfTachographCardId              |                  | 1               | 1    | {00}           |   |   |   |          |          |
| cardStructureVersion                |                  | 2               | 2    | {01 01}        |   |   |   |          |          |
| noOfEventsPerType                   |                  | 1               | 1    | {00}           |   |   |   |          |          |
| noOfFaultsPerType                   |                  | 1               | 1    | {00}           |   |   |   |          |          |
| activityStructureLength             |                  | 2               | 2    | {00 00}        |   |   |   |          |          |
| noofCardVehicleRecords              |                  | 2               | 2    | {00 00}        |   |   |   |          |          |
| noOfCardPlaceRecords                |                  | 2               | 2    | {00 00}        |   |   |   |          |          |
| noOfGNSSADRecords                   |                  | 2               | 2    | {00 00}        |   |   |   |          |          |
| noofSpecificConditionRecords        |                  | 2               | 2    | {00 00}        |   |   |   |          |          |
| noOfCardVehicleUnitRecords          |                  | 2               | 2    | {00 00}        |   |   |   |          |          |
| EF CardMA Certificate               |                  | 204             | 341  |                |   |   |   |          |          |
| CardMA Certificate                  |                  | 204             | 341  | {00..00}       |   |   |   |          |          |
| EF CardSignCertificate              |                  | 204             | 341  |                |   |   |   |          |          |
| CardSignCertificate                 |                  | 204             | 341  | {00..00}       |   |   |   |          |          |
| EF CA Certificate                   |                  | 204             | 341  |                |   |   |   |          |          |
| MemberStateCertificate              |                  | 204             | 341  | {00..00}       |   |   |   |          |          |
| EF Link Certificate                 |                  | 204             | 341  |                |   |   |   |          |          |
| LinkCertificate                     |                  | 204             | 341  | {00..00}       |   |   |   |          |          |
| EF Identification                   |                  | 143             | 143  |                |   |   |   |          |          |
| CardIdentification                  |                  | 65              | 65   |                |   |   |   |          |          |
| cardIssuingMemberState              |                  | 1               | 1    | {00}           |   |   |   |          |          |
| cardNumber                          |                  | 16              | 16   | {20..20}       |   |   |   |          |          |
| cardIssuingAuthorityName            |                  | 36              | 36   | {00, 20..20}   |   |   |   |          |          |
| cardIssueDate                       |                  | 4               | 4    | {00..00}       |   |   |   |          |          |
| cardValidityBegin                   |                  | 4               | 4    | {00..00}       |   |   |   |          |          |
| cardExpiryDate                      |                  | 4               | 4    | {00..00}       |   |   |   |          |          |
| DriverCardHolderIdentification      |                  | 78              | 78   |                |   |   |   |          |          |
| cardHolderName                      |                  | 72              | 72   |                |   |   |   |          |          |
| holderSurname                       |                  | 36              | 36   | {00, 20..20}   |   |   |   |          |          |
| holderFirstNames                    |                  | 36              | 36   | {00, 20..20}   |   |   |   |          |          |
| cardHolderBirthDate                 |                  | 4               | 4    | {00..00}       |   |   |   |          |          |
| cardHolderPreferredLanguage         |                  | 2               | 2    | {20 20}        |   |   |   |          |          |
| EF Card_Download                    |                  | 4               | 4    |                |   |   |   |          |          |
| LastCardDownload                    |                  | 4               | 4    | {00..00}       |   |   |   |          |          |
| EF Driving_Licence_Info             |                  | 53              | 53   |                |   |   |   |          |          |
| CardDrivingLicenceInformation       |                  | 53              | 53   |                |   |   |   |          |          |
| drivingLicenceIssuingAuthority      |                  | 36              | 36   | {00, 20..20}   |   |   |   |          |          |
| drivingLicenceIssuingNation         |                  | 1               | 1    | {00}           |   |   |   |          |          |
| drivingLicenceNumber                |                  | 16              | 16   | {20..20}       |   |   |   |          |          |
| EF Events_Data                      |                  | 3168            | 3168 |                |   |   |   |          |          |
| CardEventData                       |                  | 3168            | 3168 |                |   |   |   |          |          |
| cardEventRecords                    | 11               | 288             | 288  |                |   |   |   |          |          |
| CardEventRecord                     | n1               | 24              | 24   |                |   |   |   |          |          |
| eventType                           |                  | 1               | 1    | {00}           |   |   |   |          |          |
| eventBeginTime                      |                  | 4               | 4    | {00..00}       |   |   |   |          |          |
| eventEndTime                        |                  | 4               | 4    | {00..00}       |   |   |   |          |          |
| eventVehicleRegistration            |                  |                 |      |                |   |   |   |          |          |
| vehicleRegistration<br>Nation       |                  | 1               | 1    | {00}           |   |   |   |          |          |
| vehicleRegistration<br>Number       |                  | 14              | 14   | {00, 20..20}   |   |   |   |          |          |
| EF Faults_Data                      |                  | 1152            | 1152 |                |   |   |   |          |          |
| CardFaultData                       |                  | 1152            | 1152 |                |   |   |   |          |          |
| cardFaultRecords                    | 2                | 576             | 576  |                |   |   |   |          |          |
| CardFaultRecord                     | n2               | 24              | 24   |                |   |   |   |          |          |
| faultType                           |                  | 1               | 1    | {00}           |   |   |   |          |          |
| faultBeginTime                      | faultBeginTime   |                 |      | 4              | 4 | 4 | 4 | {00..00} | {00..00} |
| faultEndTime                        |                  | 4               | 4    | {00..00}       |   |   |   |          |          |
| faultVehicleRegistration            |                  |                 |      |                |   |   |   |          |          |
| vehicleRegistration                 |                  |                 |      |                |   |   |   |          |          |
| Nation                              |                  | 1               | 1    | {00}           |   |   |   |          |          |

| File Element / Data           | vehicleRegistration Number     | No of Records | Size (bytes) Min | Max   | Default Values |      |      |      |      |  |  |
|-------------------------------|--------------------------------|---------------|------------------|-------|----------------|------|------|------|------|--|--|
| EF Driver Activity Data       |                                |               | 13780            | 13780 |                |      |      |      |      |  |  |
| CardDriverActivity            |                                |               | 13780            | 13780 |                |      |      |      |      |  |  |
|                               | activityPointerOldestDayRecord |               | 2                | 2     | {00 00}        |      |      |      |      |  |  |
|                               | activityPointerNewestRecord    |               | 2                | 2     | {00 00}        |      |      |      |      |  |  |
|                               | activityDailyRecords           | n6            | 13776            | 13776 | (00..00)       |      |      |      |      |  |  |
| EF Vehicles Used              |                                |               | 9602             | 9602  |                |      |      |      |      |  |  |
| CardVehiclesUsed              |                                |               | 9602             | 9602  |                |      |      |      |      |  |  |
|                               | vehiclePointerNewestRecord     |               | 2                | 2     | {00 00}        |      |      |      |      |  |  |
|                               | cardVehicleRecords             |               | 9600             | 9600  |                |      |      |      |      |  |  |
|                               | cardVehicleRecord              | n3            | 48               | 48    |                |      |      |      |      |  |  |
|                               | vehicleOdometerBegin           |               | 3                | 3     | {00..00}       |      |      |      |      |  |  |
|                               | vehicleOdometerEnd             |               | 3                | 3     | (00..00)       |      |      |      |      |  |  |
|                               | vehicleFirstUse                |               | 4                | 4     | {00..00}       |      |      |      |      |  |  |
|                               | vehicleLastUse                 |               | 4                | 4     | {00..00}       |      |      |      |      |  |  |
|                               | vehicleRegistration            |               |                  |       |                |      |      |      |      |  |  |
|                               | vehicleRegistrationNation      |               | 1                | 1     | {00}           |      |      |      |      |  |  |
|                               | vehicleRegistrationNumber      |               | 14               | 14    | {00, 20..20}   |      |      |      |      |  |  |
|                               | vuDataBlockCounter             |               | 2                | 2     | {00 00}        |      |      |      |      |  |  |
|                               | vehicleIdentificationNumber    |               | 17               | 17    | (20..20)       |      |      |      |      |  |  |
| EF Places                     |                                |               | 2354             | 2354  |                |      |      |      |      |  |  |
| CardPlaceDailyWorkPeriod      |                                |               | 2354             | 2354  |                |      |      |      |      |  |  |
|                               | placePointerNewestRecord       |               | 2                | 2     | {00 00}        |      |      |      |      |  |  |
|                               | placeRecords                   |               | 2352             | 2352  |                |      |      |      |      |  |  |
|                               | PlaceRecord                    | n4            | 21               | 21    |                |      |      |      |      |  |  |
|                               | entryTime                      |               | 4                | 4     | {00..00}       |      |      |      |      |  |  |
|                               | entryTypeDailyWorkPeriod       |               | 1                | 1     | {00}           |      |      |      |      |  |  |
|                               | dailyWorkPeriodCountry         |               | 1                | 1     | {00}           |      |      |      |      |  |  |
|                               | dailyWorkPeriodRegion          |               | 1                | 1     | {00}           |      |      |      |      |  |  |
|                               | vehicleOdometerValue           |               | 3                | 3     | {00..00}       |      |      |      |      |  |  |
|                               | entryGNSSPlaceRecord           |               | 11               | 11    |                |      |      |      |      |  |  |
|                               | timeStamp                      |               | 4                | 4     | {00..00}       |      |      |      |      |  |  |
|                               | gnssAccuracy                   |               | 1                | 1     | {00}           |      |      |      |      |  |  |
|                               | geoCoordinates                 |               | 6                | 6     | (00..00}       |      |      |      |      |  |  |
| EF Current Usage              |                                |               | 19               | 19    |                |      |      |      |      |  |  |
| CardCurrentUse                |                                |               | 19               | 19    |                |      |      |      |      |  |  |
|                               | sessionOpenTime                |               | 4                | 4     | {00..00}       |      |      |      |      |  |  |
|                               | sessionOpenVehicle             |               |                  |       |                |      |      |      |      |  |  |
|                               | vehicleRegistrationNation      |               | 1                | 1     | {00}           |      |      |      |      |  |  |
|                               | vehicleRegistrationNumber      |               | 14               | 14    | {00, 20..20}   |      |      |      |      |  |  |
| EF Control_Activity_Data      |                                |               | 46               | 46    |                |      |      |      |      |  |  |
| CardControlActivityDataRecord |                                |               | 46               | 46    |                |      |      |      |      |  |  |
|                               | controlType                    |               | 1                | 1     | {00}           |      |      |      |      |  |  |
|                               | controlTime                    |               | 4                | 4     | {00..00}       |      |      |      |      |  |  |
|                               | controlCardNumber              |               |                  |       |                |      |      |      |      |  |  |
|                               | cardType                       |               | 1                | 1     | {00}           |      |      |      |      |  |  |
|                               | cardIssuingMemberState         |               | 1                | 1     | {00}           |      |      |      |      |  |  |
|                               | cardNumber                     |               | 16               | 16    | (20..20}       |      |      |      |      |  |  |
|                               | controlVehicleRegistration     |               |                  |       |                |      |      |      |      |  |  |
|                               | vehicleRegistrationNation      |               | 1                | 1     | {00}           |      |      |      |      |  |  |
|                               | vehicleRegistrationNumber      |               | 14               | 14    | {00, 20..20}   |      |      |      |      |  |  |
|                               | controlDownloadPeriodBegin     |               | 4                | 4     | [00..00}       |      |      |      |      |  |  |
|                               | controlDownloadPeriodEnd       |               | 4                | 4     | {00..00}       |      |      |      |      |  |  |
| EF Specific Conditions        |                                |               | 562              | 562   |                |      |      |      |      |  |  |
| SpecificConditions            |                                |               | 562              | 562   |                |      |      |      |      |  |  |
|                               | conditionPointerNewestRecord   |               | 2                | 2     | {00 00}        |      |      |      |      |  |  |
|                               | specificConditionRecords       |               | 560              | 560   |                |      |      |      |      |  |  |
|                               | SpecificConditionRecord        | n9            | 5                | 5     |                |      |      |      |      |  |  |
|                               | entryTime                      |               | 4                | 4     | {00..00}       |      |      |      |      |  |  |
|                               | specificConditionType          |               | 1                | 1     | {00}           |      |      |      |      |  |  |
| EF VehicleUnits Used          |                                |               | 2002             | 2002  |                |      |      |      |      |  |  |
| CardVehicleUnitslised         | CardVehicleUnitsUsed           |               |                  |       |                | 2002 | 2002 | 2002 | 2002 |  |  |
|                               | vehicleUnitPointerNewestRecord |               | 2                | 2     | {00 00}        |      |      |      |      |  |  |
|                               | cardVehicleUnitRecords         |               | 2000             | 2000  |                |      |      |      |      |  |  |
|                               | CardVehicleUnitRecord          | n7            | 10               | 10    |                |      |      |      |      |  |  |
|                               | timeStamp                      |               | 4                | 4     | {00..00}       |      |      |      |      |  |  |
|                               | manufacturerCode               |               | 1                | 1     | {00}           |      |      |      |      |  |  |
|                               | deviceID                       |               | 1                | 1     | {00}           |      |      |      |      |  |  |
|                               | vuSoftwareVersion              |               | 4                | 4     | {00..00}       |      |      |      |      |  |  |
| EF GNSS_Places                |                                |               | 6050             | 6050  |                |      |      |      |      |  |  |
|                               | GNSSAccumulatedDriving         |               | 6050             | 6050  |                |      |      |      |      |  |  |
|                               | gnssADPointerNewestRecord      |               | 2                | 2     | {00 00}        |      |      |      |      |  |  |

| File Element / Data                                                       | No of<br>Records         | Size (bytes)<br>Min Max | Default Values |          |      |      |      |  |
|---------------------------------------------------------------------------|--------------------------|-------------------------|----------------|----------|------|------|------|--|
| gnssAccumulatedDrivingRecords<br>GNSSAccumulatedDrivingRecord             | n8                       | 6048 6048               |                |          |      |      |      |  |
| timeStamp                                                                 |                          | 18 18                   |                |          |      |      |      |  |
| gnssPlaceRecord                                                           |                          | 4 4                     | {00..00}       |          |      |      |      |  |
| timeStamp                                                                 |                          | 14 14                   |                |          |      |      |      |  |
| gnssAccuracy                                                              |                          | 4 4                     | {00..00}       |          |      |      |      |  |
| geoCoordinates                                                            |                          | 1 1                     | {00}           |          |      |      |      |  |
| vehicleOdometerValue                                                      |                          | 6 6                     | {00..00}       |          |      |      |      |  |
| EF Application_Identification_V2<br>DriverCardApplicationIdentificationV2 |                          | 10 10                   | {00..00}       |          |      |      |      |  |
| lengthOfFollowingData                                                     |                          | 10 10                   |                |          |      |      |      |  |
| noOfBorderCrossingRecords                                                 |                          | 2 2                     | {00 00}        |          |      |      |      |  |
| noOfLoadUnloadRecords                                                     |                          | 2 2                     | {00 00}        |          |      |      |      |  |
| noOfLoadTypeEntryRecords                                                  |                          | 2 2                     | {00 00}        |          |      |      |      |  |
| VuConfigurationLengthRange                                                |                          | 2 2                     | {00 00}        |          |      |      |      |  |
| EF Places_Authentication<br>CardPlaceAuthDailyWorkPeriod                  |                          | 562 562                 |                |          |      |      |      |  |
| placeAuthPointerNewestRecord                                              |                          | 562 562                 |                |          |      |      |      |  |
| placeAuthStatusRecords                                                    |                          | 2 2                     | {00 00}        |          |      |      |      |  |
| PlaceAuthStatusRecord                                                     | n4                       | 560 560                 |                |          |      |      |      |  |
| entryTime                                                                 |                          | 5 5                     |                |          |      |      |      |  |
| authenticationStatus                                                      |                          | 4 4                     | {00..00}       |          |      |      |      |  |
| EF GNSS Places Authentication<br>GNSSAuthAccumulatedDriving               |                          | 1 1                     | {00}           |          |      |      |      |  |
| gnssAuthADPointerNewestRecord                                             |                          | 1682 1682               |                |          |      |      |      |  |
| gnssAuthStatusADRecords                                                   |                          | 1682 1682               |                |          |      |      |      |  |
| GNSSAuthStatusADRecord                                                    | n8                       | 2 2                     | {00 00}        |          |      |      |      |  |
| timeStamp                                                                 |                          | 1680 1680               |                |          |      |      |      |  |
| authenticationStatus                                                      |                          | 5 5                     |                |          |      |      |      |  |
| EF Border_Crossings<br>CardBorderCrossings                                |                          | 4 4                     | {00..00}       |          |      |      |      |  |
| borderCrossingPointerNewestRecord                                         |                          | 1 1                     | {00}           |          |      |      |      |  |
| cardBorderCrossingRecords                                                 |                          | 19042 19042             |                |          |      |      |      |  |
| CardBorderCrossingRecord                                                  | n10                      | 19042 19042             |                |          |      |      |      |  |
| countryLeft                                                               |                          | 2 2                     | {00 00}        |          |      |      |      |  |
| countryEntered                                                            |                          | 19040 19040             |                |          |      |      |      |  |
| gnssPlaceAuthRecord                                                       |                          | 17 17                   |                |          |      |      |      |  |
| timeStamp                                                                 |                          | 1 1                     | {00}           |          |      |      |      |  |
| gnssAccuracy                                                              |                          | 1 1                     | {00}           |          |      |      |      |  |
| geoCoordinates                                                            |                          | 12 12                   |                |          |      |      |      |  |
| authenticationStatus                                                      |                          | 4 4                     | {00..00}       |          |      |      |      |  |
| vehicleOdometerValue                                                      |                          | 1 1                     | {00}           |          |      |      |      |  |
| EF Load Unload Operations<br>CardLoadUnloadOperations                     |                          | 6 6                     | {00..00}       |          |      |      |      |  |
| loadUnloadPointerNewestRecord                                             |                          | 1 1                     | {00}           |          |      |      |      |  |
| cardloadUnloadRecords                                                     |                          | 3 3                     | {00..00}       |          |      |      |      |  |
| CardLoadUnloadRecord                                                      | n11                      | 32482 32482             |                |          |      |      |      |  |
| timestamp                                                                 |                          | 32482 32482             |                |          |      |      |      |  |
| operationType                                                             |                          | 2 2                     | {00 00}        |          |      |      |      |  |
| gnssPlaceAuthRecord                                                       |                          | 32480 32480             |                |          |      |      |      |  |
| timeStamp                                                                 |                          | 20 20                   |                |          |      |      |      |  |
| gnssAccuracy                                                              |                          | 4 4                     | {00}           |          |      |      |      |  |
| geoCoordinates                                                            |                          | 1 1                     | {00..00}       |          |      |      |      |  |
| authenticationStatus                                                      |                          | 12 12                   |                |          |      |      |      |  |
| vehicleOdometerValue                                                      |                          | 4 4                     | {00..00}       |          |      |      |      |  |
| EF Load_Type_Entries<br>CardLoadTypeEntries                               |                          | 1 1                     | {00}           |          |      |      |      |  |
| loadtypeEntryPointerNewestRecord                                          |                          | 6 6                     | {00..00}       |          |      |      |      |  |
| cardLoadTypeEntryRecords                                                  | cardLoadTypeEntryRecords |                         |                | 1 1      | 1680 | {00} | 1680 |  |
| CardLoadTypeEntryRecord                                                   | n12                      | 5                       | 5              |          |      |      |      |  |
| timestamp                                                                 |                          | 4                       | 4              | {00..00} |      |      |      |  |
| loadTypeEntered                                                           |                          | 1                       | 1              | {00}     |      |      |      |  |
| EF VU Configuration                                                       |                          | 3072                    | 3072           |          |      |      |      |  |
| VuConfigurations                                                          | n13                      | 3072                    | 3072           |          |      |      |      |  |

TCS\_155 The following values, used to provide sizes in the table above, are the minimum and maximum record number values the driver card data structure must use for a generation 2 application:

|     |                              | Min                                                | Max                                                |
|-----|------------------------------|----------------------------------------------------|----------------------------------------------------|
| n1  | NoOfEventsPerType            | 12                                                 | 12                                                 |
| n2  | NoOfFaultsPerType            | 24                                                 | 24                                                 |
| n3  | NoOfCardVehicleRecords       | 200                                                | 200                                                |
| n4  | NoOfCardPlaceRecords         | 112                                                | 112                                                |
| n6  | CardActivityLengthRange      | 13776 Bytes<br>(56 days * 117 activity<br>changes) | 13776 Bytes<br>(56 days * 117 activity<br>changes) |
| n7  | NoOfCardVehicleUnitRecords   | 200                                                | 200                                                |
| n8  | NoOfGNSSADRecords            | 336                                                | 336                                                |
| n9  | NoOfSpecificConditionRecords | 112                                                | 112                                                |
| n10 | NoOfBorderCrossingRecords    | 1120                                               | 1120                                               |
| n11 | NoOfLoadUnloadRecords        | 1624                                               | 1624                                               |
| n12 | NoOfLoadTypeEntryRecords     | 336                                                | 336                                                |
| n13 | VuConfigurationLengthRange   | 3072 Bytes                                         | 3072 Bytes                                         |

### **▼B**

### 4.3. **Workshop card applications**

TCS\_156 After its personalisation, the workshop card application generation 1 shall have the following permanent file structure and file access rules:

|                                |         |      | Access rules |        |  |
|--------------------------------|---------|------|--------------|--------|--|
| File                           | File ID | Read | Select       | Update |  |
| -DF Tachograph                 | '0500h' |      | SC1          |        |  |
| -EF Application_Identification | '0501h' | SC2  | SC1          | NEV    |  |
| -EF Card_Certificate           | 'C100h' | SC2  | SC1          | NEV    |  |
| -EF CA_Certificate             | 'C108h' | SC2  | SC1          | NEV    |  |
| -EF Identification             | '0520h' | SC2  | SC1          | NEV    |  |
| -EF Card_Download              | '0509h' | SC2  | SC1          | SC1    |  |
| -EF Calibration                | '050Ah' | SC2  | SC1          | SC3    |  |
| -EF Sensor_Installation_Data   | '050Bh' | SC4  | SC1          | NEV    |  |
| -EF Events_Data                | '0502h' | SC2  | SC1          | SC3    |  |
| -EF Faults_Data                | '0503h' | SC2  | SC1          | SC3    |  |
| -EF Driver_Activity_Data       | '0504h' | SC2  | SC1          | SC3    |  |
| -EF Vehicles_Used              | '0505h' | SC2  | SC1          | SC3    |  |
| -EF Places                     | '0506h' | SC2  | SC1          | SC3    |  |
| -EF Current_Usage              | '0507h' | SC2  | SC1          | SC3    |  |
| -EF Control_Activity_Data      | '0508h' | SC2  | SC1          | SC3    |  |
| -EF Specific Conditions        | '0522h' | SC2  | SC1          | SC3    |  |

<sup>4.3.1</sup> *Workshop card application generation 1*

The following abbreviations for the security conditions are used in this table:

**SC1** ALW OR SM-MAC-G2

**SC2** ALW OR SM-MAC-G1 OR SM-MAC-G2

**SC3** SM-MAC-G1 OR SM-MAC-G2

byte (if supported): NEV

**▼M1 SC4** For the READ BINARY command with even INS byte: (SM-C-MAC-G1 AND SM-R-ENC-MAC-G1) OR (SM-C-MAC-G2 AND SM-R-ENC-MAC-G2) For the READ BINARY command with odd INS

**▼B**

TCS\_157 All EF structures shall be transparent.

TCS\_158 The workshop card application generation 1 shall have the following data structure:

▼<u>B</u>

| File / Data element                      | No of<br>Records | Size (Bytes)<br>Min | Max   | Default<br>Values |
|------------------------------------------|------------------|---------------------|-------|-------------------|
| LDF Tachograph                           |                  | 11055               | 29028 |                   |
| EF Application_Identification            |                  | 11                  | 11    |                   |
| └─ WorkshopCardApplicationIdentification |                  | 11                  | 11    |                   |
| typeOfTachographCardId                   |                  | 1                   | 1     | {00}              |
| -cardStructureVersion                    |                  | 2                   | 2     | {00 00}           |
| -noOfEventsPerType                       |                  | 1                   | 1     | {00}              |
| noOfFaultsPerType                        |                  | 1                   | 1     | {00}              |
| activityStructureLength                  |                  | 2                   | 2     | {00 00}           |
| -noOfCardVehicleRecords                  |                  | 2                   | 2     | {00 00}           |
| -noOfCardPlaceRecords                    |                  | 1                   | 1     | {00}              |
| noofCalibrationRecords                   |                  | 1                   | 1     | {00}              |
| EF Card_Certificate                      |                  | 194                 | 194   |                   |
| └─ CardCertificate                       |                  | 194                 | 194   | {00..00}          |
| EF CA_Certificate                        |                  | 194                 | 194   |                   |
| MemberStateCertificate                   |                  | 194                 | 194   | {00..00}          |
| EF Identification                        |                  | 211                 | 211   |                   |
| -CardIdentification                      |                  | 65                  | 65    |                   |
| └── cardIssuingMemberState               |                  | 1                   | 1     | {00}              |
| - cardNumber                             |                  | 16                  | 16    | {20..20}          |
| _ cardIssuingAuthorityName               |                  | 36                  | 36    | {00, 20..2        |
| cardIssue Date                           |                  | 4                   | 4     | {00..00}          |
| - cardValidityBegin                      |                  | 4                   | 4     | {00..00}          |
| CardExpiryDate                           |                  | 4                   | 4     | {00..00}          |
| WorkshopCardHolderIdentification         |                  | 146                 | 146   |                   |
| workshopName                             |                  | 36                  | 36    | {00, 20..2        |
| — workshopAddress                        |                  | 36                  | 36    | {00, 20..2}       |
| - cardHolderName                         |                  |                     |       |                   |
| holderSurname                            |                  | 36                  | 36    | {00, 20..2}       |
| L holderFirstNames                       |                  | 36                  | 36    | {00, 20..2}       |
| CardHolderPreferredLanguage              |                  | 2                   | 2     | {20 20}           |
| EF Card_Download                         |                  | 2                   | 2     |                   |
| LNoOfCalibrationsSinceDownload           |                  | 2                   | 2     | {00 00}           |
| EF Calibration                           |                  | 9243                | 26778 |                   |
| WorkshopCardCalibrationData              |                  | 9243                | 26778 |                   |
| -calibrationTotalNumber                  |                  | 2                   | 2     | {00 00}           |
| - calibrationPointerNewestRecord         |                  | 1                   | 1     | {00}              |
| L_calibrationRecords                     |                  | 9240                | 26775 |                   |
| WorkshopCardCalibrationRecord            | n5               | 105                 | 105   |                   |
| -calibrationPurpose                      |                  | 1                   | 1     | {00}              |
| -vehicleIdentificationNumber             |                  | 17                  | 17    | {20..20}          |
| -vehicleRegistration                     |                  |                     |       |                   |
| -vehicleRegistrationNation               |                  | 1                   | 1     | {00}              |
| vehicleRegistrationNumber                |                  | 14                  | 14    | {00, 20..2}       |
| wVehicleCharacteristicConstant           |                  | 2                   | 2     | {00 00}           |
| - kConstantOfRecordingEquipment          |                  | 2                   | 2     | {00 00}           |
| - lTyreCircumference                     |                  | 2                   | 2     | {00 00}           |
| tyreSize                                 |                  | 15                  | 15    | {20..20}          |
| -authorisedSpeed                         |                  | 1                   | 1     | {00}              |
| oldOdometerValue                         |                  | 3                   | 3     | {00..00}          |
| newOdometerValue                         |                  | 3                   | 3     | {00..00}          |
| -oldTimeValue                            |                  | 4                   | 4     | {00..00}          |
| newTimeValue                             |                  | 4                   | 4     | {00..00}          |
| -nextCalibration Date                    |                  | 4                   | 4     | {00..00}          |
| vuPartNumber                             |                  | 16                  | 16    | {20..20}          |
| vuSerialNumber                           |                  | 8                   | 8     | {00..00}          |
| EF Sensor_Installation_Data              |                  | 16                  | 16    |                   |
| └ SensorInstallationSecData              |                  | 16                  | 16    | {00..00}          |
| EF Events_Data                           |                  | 432                 | 432   |                   |
| └ CardEventData                          |                  | 432                 | 432   |                   |
| └─cardEventRecords                       | 6                | 72                  | 72    |                   |
| └ CardEventRecord                        | $n_1$            | 24                  | 24    |                   |
| - eventType                              |                  | 1                   | 1     | {00}              |
| eventBeginTime                           |                  | 4                   | 4     | {00..00}          |
| - eventEndTime                           |                  | 4                   | 4     | {00..00}          |
| eventVehicleRegistration                 |                  |                     |       |                   |
| -vehicleRegistrationNation               |                  | 1                   | 1     | {00}              |
| -vehicleRegistrationNumber               |                  | 14                  | 14    | {00, 20..20}      |
| EF Faults_Data                           |                  | 288                 | 288   |                   |
| └ CardFaultData                          |                  | 288                 | 288   |                   |
| └─ cardFaultRecords                      | 2                | 144                 | 144   |                   |
| CardFaultRecord                          | $n_2$            | 24                  | 24    |                   |
| - faultType                              |                  | 1                   | 1     | {00}              |
| faultBeginTime                           |                  | 4                   | 4     | {00..00}          |
| — faultEndTime                           |                  | 4                   | 4     | {00..00}          |
| faultVehicleRegistration                 |                  |                     |       |                   |
| -vehicleRegistrationNation               |                  | 1                   | 1     | {00}              |
| vehicleRegistrationNumber                |                  | 14                  | 14    | {00, 20..20}      |
| _EF Driver_Activity_Data                 |                  | 202                 | 496   |                   |
| └─ CardDriverActivity                    |                  | 202                 | 496   |                   |
| activityPointerOldestDayRecord           |                  | 2                   | 2     | {00 00}           |
| activityPointerNewestRecord              |                  | 2                   | 2     | {00 00}           |
| LactivityDailyRecords                    | $n_6$            | 198                 | 492   | {00..00}          |
| _EF Vehicles Used                        |                  | 126                 | 250   |                   |
| L_CardVehiclesUsed                       |                  | 126                 | 250   |                   |
| -vehiclePointerNewestRecord              |                  | 2                   | 2     | {00 00}           |
| cardVehicleRecords                       |                  | 124                 | 248   |                   |
| CardVehicleRecord                        | $n_3$            | 31                  | 31    |                   |
| - vehicleOdometerBegin                   |                  | 3                   | 3     | {00..00}          |
| vehicleOdometerEnd                       |                  | 3                   | 3     | {00..00}          |
| vehicleFirstUse                          |                  | 4                   | 4     | {00..00}          |
| -vehicleLastUse                          |                  | 4                   | 4     | {00..00}          |
| - vehicleRegistration                    |                  |                     |       |                   |
| -vehicleRegistrationNation               |                  | 1                   | 1     | {00}              |
| -vehicleRegistrationNumber               |                  | 14                  | 14    | {00, 20..20}      |
| -vuDataBlockCounter                      |                  | 2                   | 2     | {00 00}           |
| _EF Places                               |                  | 61                  | 81    |                   |
| └ CardPlaceDailyWorkPeriod               |                  | 61                  | 81    |                   |
| -placePointerNewestRecord                |                  | 1                   | 1     | {00}              |
| placeRecords                             |                  | 60                  | 80    |                   |
| └ PlaceRecord                            | $n_4$            | 10                  | 10    |                   |
| - entryTime                              |                  | 4                   | 4     | {00..00}          |
| entryTypeDailyWorkPeriod                 |                  | 1                   | 1     | {00}              |
| -dailyWorkPeriodCountry                  |                  | 1                   | 1     | {00}              |
| dailyWorkPeriodRegion                    |                  | 1                   | 1     | {00}              |
| vehicleOdometerValue                     |                  | 3                   | 3     | {00..00}          |
| EF Current Usage                         |                  | 19                  | 19    |                   |
| CardCurrentUse                           |                  | 19                  | 19    |                   |
| -sessionOpenTime                         |                  | 4                   | 4     | {00..00}          |
| sessionOpenVehicle                       |                  |                     |       |                   |
| -vehicleRegistrationNation               |                  | 1                   | 1     | {00}              |
| -vehicleRegistrationNumber               |                  | 14                  | 14    | {00, 20..20}      |
| EF Control_Activity_Data                 |                  | 46                  | 46    |                   |
| CardControlActivityDataRecord            |                  | 46                  | 46    |                   |
| - controlType                            |                  | 1                   | 1     | {00}              |
| - controlTime                            |                  | 4                   | 4     | {00..00}          |
| - controlCardNumber                      |                  |                     |       |                   |
| └─ cardType                              |                  | 1                   | 1     | {00}              |
| - cardIssuingMemberState                 |                  | 1                   | 1     | {00}              |
| └─ cardNumber                            |                  | 16                  | 16    | {20..20}          |
| - controlVehicleRegistration             |                  |                     |       |                   |
| - vehicleRegistrationNation              |                  | 1                   | 1     | {00}              |
| - vehicleRegistrationNumber              |                  | 14                  | 14    | {00, 20..20}      |
| - controlDownloadPeriodBegin             |                  | 4                   | 4     | {00..00}          |
| - controlDownloadPeriodEnd               |                  | 4                   | 4     | {00..00}          |
| EF Specific_Conditions                   |                  | 10                  | 10    |                   |
| └─ SpecificConditionRecord               |                  | 2                   | 5     | 5                 |
| - entryTime                              |                  | 4                   | 4     | {00..00}          |
| - SpecificConditionType                  |                  | 1                   | 1     | {00}              |

### **B**

#### TCS\_159 The following values, used to provide sizes in the table above, are the minimum and maximum record number values the workshop card data structure must use for a generation 1 application:

|    |                         | Min                                        | Max                                         |
|----|-------------------------|--------------------------------------------|---------------------------------------------|
| n1 | NoOfEventsPerType       | 3                                          | 3                                           |
| n2 | NoOfFaultsPerType       | 6                                          | 6                                           |
| n3 | NoOfCardVehicleRecords  | 4                                          | 8                                           |
| n4 | NoOfCardPlaceRecords    | 6                                          | 8                                           |
| n5 | NoOfCalibrationRecords  | 88                                         | 255                                         |
| n6 | CardActivityLengthRange | 198 bytes (1 day *<br>93 activity changes) | 492 bytes (1 day *<br>240 activity changes) |

4.3.2 *Workshop card application generation 2*

**▼M3**

TCS\_160 After its personalisation, the workshop card application generation 2 shall have the following permanent file structure and file access rules.

*Notes:*

- The short EF identifier SFID is given as decimal number, e.g. the value 30 corresponds to 11110 in binary.
- EF Application\_Identification\_V2, EF Places\_Authentication, EF GNSS\_Places\_Authentication, EF Border\_Crossings, EF Load\_Unload\_Operations, EF Load\_Type\_Entries, EF VU\_Configuration and EF Calibration\_Add\_Data are only present in version 2 of the generation 2 workshop card.
- cardStructureVersion in EF Application\_Identification is equal to {01 01} for version 2 of the generation 2 workshop card, while it was equal to {01 00} for version 1 of the generation 2 workshop card.

| ▼M3 |
|-----|
|     |

| File                              | File ID | SFID | Access rules |           |           |
|-----------------------------------|---------|------|--------------|-----------|-----------|
|                                   |         |      | Read         | Select    | Update    |
| └─DF Tachograph_G2                |         |      | SC1          | SC1       |           |
| -EF Application_Identification    | '0501h' | 1    | SC1          | SC1       | NEV       |
| -EF CardMA_Certificate            | 'C100h' | 2    | SC1          | SC1       | NEV       |
| -EF CardSignCertificate           | 'C101h' | 3    | SC1          | SC1       | NEV       |
| -EF CA_Certificate                | 'C108h' | 4    | SC1          | SC1       | NEV       |
| -EF Link_Certificate              | 'C109h' | 5    | SC1          | SC1       | NEV       |
| -EF Identification                | '0520h' | 6    | SC1          | SC1       | NEV       |
| -EF Card_Download                 | '0509h' | 7    | SC1          | SC1       | SC1       |
| -EF Calibration                   | '050Ah' | 10   | SC1          | SC1       | SM-MAC-G2 |
| -EF Sensor_Installation_Data      | '050Bh' | 11   | SC5          | SM-MAC-G2 | NEV       |
| -EF Events_Data                   | '0502h' | 12   | SC1          | SC1       | SM-MAC-G2 |
| -EF Faults_Data                   | '0503h' | 13   | SC1          | SC1       | SM-MAC-G2 |
| -EF Driver_Activity_Data          | '0504h' | 14   | SC1          | SC1       | SM-MAC-G2 |
| -EF Vehicles_Used                 | '0505h' | 15   | SC1          | SC1       | SM-MAC-G2 |
| -EF Places                        | '0506h' | 16   | SC1          | SC1       | SM-MAC-G2 |
| -EF Current_Usage                 | '0507h' | 17   | SC1          | SC1       | SM-MAC-G2 |
| -EF Control_Activity_Data         | '0508h' | 18   | SC1          | SC1       | SM-MAC-G2 |
| -EF Specific_Conditions           | '0522h' | 19   | SC1          | SC1       | SM-MAC-G2 |
| -EF VehicleUnits_Used             | '0523h' | 20   | SC1          | SC1       | SM-MAC-G2 |
| -EF GNSS_Places                   | '0524h' | 21   | SC1          | SC1       | SM-MAC-G2 |
| -EF Application_Identification_V2 | '0525h' | 22   | SC1          | SC1       | NEV       |
| -EF Places_Authentication         | '0526h' | 23   | SC1          | SC1       | SM-MAC-G2 |
| -EF GNSS_Places_Authentication    | '0527h' | 24   | SC1          | SC1       | SM-MAC-G2 |
| -EF Border_Crossings              | '0528h' | 25   | SC1          | SC1       | SM-MAC-G2 |
| -EF Load_Unload_Operations        | '0529h' | 26   | SC1          | SC1       | SM-MAC-G2 |
| -EF Load_Type_Entries             | '0530h' | 27   | SC1          | SC1       | SM-MAC-G2 |
| -EF Calibration_Add_Data          | '0531h' | 28   | SC1          | SC1       | SM-MAC-G2 |
| -EF VU_Configuration              | '0540h' | 30   | SC5          | SC1       | SM-MAC-G2 |

The following abbreviations for the security conditions are used in this table:

#### **SC1** ALW OR SM-MAC-G2

**SC5** For the Read Binary command with even INS byte: SM-C-MAC-G2 AND SM-R-ENC-MAC-G2

> For the Read Binary command with odd INS byte (if supported): NEV

### **B**

TCS\_161 All EFs structures shall be transparent.

TCS\_162 The workshop card application generation 2 shall have the following data structure:

| File / Data Element                   | No of Records                 | Size (bytes)<br>Min | Max   | Default Value |   |   |   |         |         |
|---------------------------------------|-------------------------------|---------------------|-------|---------------|---|---|---|---------|---------|
| DF Tachograph G2                      |                               | 59582               | 60214 |               |   |   |   |         |         |
| EF Application_Identification         |                               | 19                  | 19    |               |   |   |   |         |         |
| WorkshopCardApplicationIdentification |                               | 19                  | 19    |               |   |   |   |         |         |
| typeOfTachographCardId                |                               | 1                   | 1     | {00}          |   |   |   |         |         |
| cardStructureVersion                  |                               | 2                   | 2     | {01 01}       |   |   |   |         |         |
| noOfEventsPerType                     |                               | 1                   | 1     | {00}          |   |   |   |         |         |
| noOfFaultsPerType                     |                               | 1                   | 1     | {00}          |   |   |   |         |         |
| activityStructureLength               |                               | 2                   | 2     | {00 00}       |   |   |   |         |         |
| noOfCardVehicleRecords                |                               | 2                   | 2     | {00 00}       |   |   |   |         |         |
| noOfCardPlaceRecords                  |                               | 2                   | 2     | {00 00}       |   |   |   |         |         |
| noOfCalibrationRecords                |                               | 2                   | 2     | {00 00}       |   |   |   |         |         |
| noOfGNSSADRecords                     |                               | 2                   | 2     | {00 00}       |   |   |   |         |         |
| noOfSpecificConditionRecords          |                               | 2                   | 2     | {00 00}       |   |   |   |         |         |
| noOfCardVehicleUnitRecords            |                               | 2                   | 2     | {00 00}       |   |   |   |         |         |
| EF CardMA Certificate                 |                               | 204                 | 341   |               |   |   |   |         |         |
| CardMA Certificate                    |                               | 204                 | 341   | {00..00}      |   |   |   |         |         |
| EF CardSignCertificate                |                               | 204                 | 341   |               |   |   |   |         |         |
| CardSignCertificate                   |                               | 204                 | 341   | {00..00}      |   |   |   |         |         |
| EF CA Certificate                     |                               | 204                 | 341   |               |   |   |   |         |         |
| MemberStateCertificate                |                               | 204                 | 341   | (00..00}      |   |   |   |         |         |
| EF Link Certificate                   |                               | 204                 | 341   |               |   |   |   |         |         |
| LinkCertificate                       |                               | 204                 | 341   | {00..00}      |   |   |   |         |         |
| EF Identification                     |                               | 211                 | 211   |               |   |   |   |         |         |
| CardIdentification                    |                               | 65                  | 65    |               |   |   |   |         |         |
| cardIssuingMemberState                |                               | 1                   | 1     | {00}          |   |   |   |         |         |
| cardNumber                            |                               | 16                  | 16    | {20..20}      |   |   |   |         |         |
| cardIssuingAuthorityName              |                               | 36                  | 36    | {00, 20..20}  |   |   |   |         |         |
| cardIssueDate                         |                               | 4                   | 4     | {00..00}      |   |   |   |         |         |
| cardValidityBegin                     |                               | 4                   | 4     | {00..00}      |   |   |   |         |         |
| cardExpiryDate                        |                               | 4                   | 4     | {00..00}      |   |   |   |         |         |
| WorkshopCardHolderIdentification      |                               | 146                 | 146   |               |   |   |   |         |         |
| workshopName                          |                               | 36                  | 36    |               |   |   |   |         |         |
| workshopAddress                       |                               | 36                  | 36    |               |   |   |   |         |         |
| cardHolderName                        |                               | 72                  | 72    |               |   |   |   |         |         |
| holderSurname                         |                               | 36                  | 36    | {00, 20..20}  |   |   |   |         |         |
| holderFirstNames                      |                               | 36                  | 36    | {00, 20..20}  |   |   |   |         |         |
| cardHolderPreferredLanguage           |                               | 2                   | 2     | {20 20}       |   |   |   |         |         |
| EF Card Download                      |                               | 2                   | 2     |               |   |   |   |         |         |
| NoOfCalibrationsSinceDownload         |                               | 2                   | 2     | {00 00}       |   |   |   |         |         |
| EF Calibration                        |                               | 45394               | 45394 |               |   |   |   |         |         |
| WorkshopCardCalibrationData           |                               | 45394               | 45394 |               |   |   |   |         |         |
| calibrationTotalNumber                |                               | 2                   | 2     | {00 00}       |   |   |   |         |         |
| calibrationPointerNewestRecord        |                               | 2                   | 2     | {00 00}       |   |   |   |         |         |
| calibrationRecords                    |                               | 45390               | 45390 |               |   |   |   |         |         |
| WorkshopCardCalibrationRecord         | n5                            | 178                 | 178   |               |   |   |   |         |         |
| calibrationPurpose                    |                               | 1                   | 1     | {00}          |   |   |   |         |         |
| vehicleIdentificationNumber           |                               | 17                  | 17    | {20..20}      |   |   |   |         |         |
| vehicleRegistration                   |                               |                     |       |               |   |   |   |         |         |
| vehicleRegistrationNation             |                               | 1                   | 1     | {00}          |   |   |   |         |         |
| vehicleRegistrationNumber             |                               | 14                  | 14    | {00, 20..20}  |   |   |   |         |         |
| wVehicleCharacteristicConstant        | kConstantOfRecordingEquipment |                     |       | 2             | 2 | 2 | 2 | {00 00} | {00 00} |
| lTyreCircumference                    |                               | 2                   | 2     | {00 00}       |   |   |   |         |         |
| tyreSize                              |                               | 15                  | 15    | {20..20}      |   |   |   |         |         |
| authorisedSpeed                       |                               | 1                   | 1     | {00}          |   |   |   |         |         |
| oldOdometerValue                      |                               | 3                   | 3     | {00..00}      |   |   |   |         |         |
| newOdometerValue                      |                               | 3                   | 3     | {00..00}      |   |   |   |         |         |
| oldTimeValue                          |                               | 4                   | 4     | {00..00}      |   |   |   |         |         |

| File / Data Element            | No of Records             | Size (bytes) Min | Max | Default Values |        |  |  |
|--------------------------------|---------------------------|------------------|-----|----------------|--------|--|--|
| newTimeValue                   |                           | 4                | 4   | {00..00}       |        |  |  |
| nextCalibrationDate            |                           | 4                | 4   | {00..00}       |        |  |  |
| vuPartNumber                   |                           | 16               | 16  | {20..20}       |        |  |  |
| vuSerialNumber                 |                           | 8                | 8   | {00..00}       |        |  |  |
| sensorSerialNumber             |                           | 8                | 8   | {00..00}       |        |  |  |
| sensorGNSSSerialNumber         |                           | 8                | 8   | {00..00}       |        |  |  |
| rcmSerialNumber                |                           | 8                | 8   | {00..00}       |        |  |  |
| vuAbility                      |                           | 1                | 1   | {00}           |        |  |  |
| sealDataCard                   |                           | 56               | 56  |                |        |  |  |
| noOfSealRecords                |                           | 1                | 1   | {00}           |        |  |  |
| SealRecords                    |                           | 55               | 55  |                |        |  |  |
| SealRecord                     | 5                         | 11               | 11  |                |        |  |  |
| equipmentType                  |                           | 1                | 1   | {00}           |        |  |  |
| extendedSealIdentifier         |                           | 10               | 10  | {00..00}       |        |  |  |
| EF Sensor Installation Data    |                           | 18               | 102 |                |        |  |  |
| SensorInstallationSecData      |                           | 18               | 102 | {00..00}       |        |  |  |
| EF Events Data                 |                           | 792              | 792 |                |        |  |  |
| CardEventData                  |                           | 792              | 792 |                |        |  |  |
| cardEventRecords               | 11                        | 72               | 72  |                |        |  |  |
| CardEventRecord                | n1                        | 24               | 24  |                |        |  |  |
| eventType                      |                           | 1                | 1   | {00}           |        |  |  |
| eventBeginTime                 |                           | 4                | 4   | {00..00}       |        |  |  |
| eventEndTime                   |                           | 4                | 4   | {00..00}       |        |  |  |
| eventVehicleRegistration       |                           |                  |     |                |        |  |  |
| vehicleRegistrationNation      |                           | 1                | 1   | {00}           |        |  |  |
| vehicleRegistrationNumber      |                           | 14               | 14  | {00, 20..20}   |        |  |  |
| EF Faults Data                 |                           | 288              | 288 |                |        |  |  |
| CardFaultData                  |                           | 288              | 288 |                |        |  |  |
| cardFaultRecords               | 2                         | 144              | 144 |                |        |  |  |
| CardFaultRecord                | n2                        | 24               | 24  |                |        |  |  |
| faultType                      |                           | 1                | 1   | {00}           |        |  |  |
| faultBeginTime                 |                           | 4                | 4   | {00..00}       |        |  |  |
| faultEndTime                   |                           | 4                | 4   | {00..00}       |        |  |  |
| faultVehicleRegistration       |                           |                  |     |                |        |  |  |
| vehicleRegistrationNation      |                           | 1                | 1   | {00}           |        |  |  |
| vehicleRegistrationNumber      |                           | 14               | 14  | {00, 20..20}   |        |  |  |
| EF Driver Activity Data        |                           | 496              | 496 |                |        |  |  |
| CardDriverActivity             |                           | 496              | 496 |                |        |  |  |
| activityPointerOldestDayRecord |                           | 2                | 2   | {00 00}        |        |  |  |
| activityPointerNewestRecord    |                           | 2                | 2   | {00 00}        |        |  |  |
| activityDailyRecords           | n6                        | 492              | 492 | {00..00}       |        |  |  |
| EF Vehicles Used               |                           | 386              | 386 |                |        |  |  |
| CardVehiclesUsed               |                           | 386              | 386 |                |        |  |  |
| vehiclePointerNewestRecord     |                           | 2                | 2   | {00 00}        |        |  |  |
| cardVehicleRecords             |                           | 384              | 384 |                |        |  |  |
| cardVehicleRecord              | n3                        | 48               | 48  |                |        |  |  |
| vehicleOdometerBegin           |                           | 3                | 3   | {00..00}       |        |  |  |
| vehicleOdometerEnd             |                           | 3                | 3   | {00..00}       |        |  |  |
| vehicleFirstUse                |                           | 4                | 4   | {00..00}       |        |  |  |
| vehicleLastUse                 |                           | 4                | 4   | {00..00}       |        |  |  |
| vehicleRegistration            | vehicleRegistrationNation |                  | 1   |                | 1 {00} |  |  |
| vehicleRegistrationNumber      | 14                        | 14 {00, 20..20}  |     |                |        |  |  |
| vuDataBlockCounter             | 2                         | 2 {00 00}        |     |                |        |  |  |

| File / Data Element                     | vehicleIdentificationNumber | No of Records | Size (bytes) Min | Max | Default Value |   |   |   |   |          |          |
|-----------------------------------------|-----------------------------|---------------|------------------|-----|---------------|---|---|---|---|----------|----------|
| EF Places                               |                             |               | 17               | 17  | {20..20}      |   |   |   |   |          |          |
| CardPlaceDailyWorkPeriod                |                             |               | 170              | 170 |               |   |   |   |   |          |          |
| placePointerNewestRecord                |                             |               | 2                | 2   | {00 00}       |   |   |   |   |          |          |
| placeRecords                            |                             |               | 168              | 168 |               |   |   |   |   |          |          |
| PlaceRecord                             |                             | n4            | 21               | 21  |               |   |   |   |   |          |          |
| entryTime                               |                             |               | 4                | 4   | {00..00}      |   |   |   |   |          |          |
| entryTypeDailyWorkPeriod                |                             |               | 1                | 1   | {00}          |   |   |   |   |          |          |
| dailyWorkPeriodCountry                  |                             |               | 1                | 1   | {00}          |   |   |   |   |          |          |
| dailyWorkPeriodRegion                   |                             |               | 1                | 1   | {00}          |   |   |   |   |          |          |
| vehicleOdometerValue                    |                             |               | 3                | 3   | {00..00}      |   |   |   |   |          |          |
| entryGNSSPlaceRecord                    |                             |               | 11               | 11  |               |   |   |   |   |          |          |
| timeStamp                               |                             |               | 4                | 4   | {00..00}      |   |   |   |   |          |          |
| gnssAccuracy                            |                             |               | 1                | 1   | {00}          |   |   |   |   |          |          |
| geoCoordinates                          |                             |               | 6                | 6   | {00..00}      |   |   |   |   |          |          |
| EF Current_Usage                        |                             |               | 19               | 19  |               |   |   |   |   |          |          |
| CardCurrentUse                          |                             |               | 19               | 19  |               |   |   |   |   |          |          |
| sessionOpenTime                         |                             |               | 4                | 4   | {00..00}      |   |   |   |   |          |          |
| sessionOpenVehicle                      |                             |               |                  |     |               |   |   |   |   |          |          |
| vehicleRegistrationNation               |                             |               | 1                | 1   | {00}          |   |   |   |   |          |          |
| vehicleRegistrationNumber               |                             |               | 14               | 14  | {00, 20..20}  |   |   |   |   |          |          |
| EF Control_Activity_Data                |                             |               | 46               | 46  |               |   |   |   |   |          |          |
| CardControlActivityDataRecord           |                             |               | 46               | 46  |               |   |   |   |   |          |          |
| controlType                             |                             |               | 1                | 1   | {00}          |   |   |   |   |          |          |
| controlTime                             |                             |               | 4                | 4   | {00..00}      |   |   |   |   |          |          |
| controlCardNumber                       |                             |               |                  |     |               |   |   |   |   |          |          |
| cardType                                |                             |               | 1                | 1   | {00}          |   |   |   |   |          |          |
| cardIssuingMemberState                  |                             |               | 1                | 1   | {00}          |   |   |   |   |          |          |
| cardNumber                              |                             |               | 16               | 16  | {20..20}      |   |   |   |   |          |          |
| controlVehicleRegistration              |                             |               |                  |     |               |   |   |   |   |          |          |
| vehicleRegistrationNation               |                             |               | 1                | 1   | {00}          |   |   |   |   |          |          |
| vehicleRegistrationNumber               |                             |               | 14               | 14  | {00, 20..20}  |   |   |   |   |          |          |
| controlDownloadPeriodBegin              |                             |               | 4                | 4   | {00..00}      |   |   |   |   |          |          |
| controlDownloadPeriodEnd                |                             |               | 4                | 4   | {00..00}      |   |   |   |   |          |          |
| EF VehicleUnits Used                    |                             |               | 82               | 82  |               |   |   |   |   |          |          |
| CardVehicleUnitsUsed                    |                             |               | 82               | 82  |               |   |   |   |   |          |          |
| vehicleUnitPointerNewestRecord          |                             |               | 2                | 2   | {00 00}       |   |   |   |   |          |          |
| cardVehicleUnitRecords                  |                             |               | 80               | 80  |               |   |   |   |   |          |          |
| CardVehicleUnitRecord                   |                             | n7            | 10               | 10  |               |   |   |   |   |          |          |
| timeStamp                               |                             |               | 4                | 4   | {00..00}      |   |   |   |   |          |          |
| manufacturerCode                        |                             |               | 1                | 1   | {00}          |   |   |   |   |          |          |
| deviceID                                |                             |               | 1                | 1   | {00}          |   |   |   |   |          |          |
| vuSoftwareVersion                       |                             |               | 4                | 4   | {00..00}      |   |   |   |   |          |          |
| EF GNSS Places                          |                             |               | 434              | 434 |               |   |   |   |   |          |          |
| GNSSAccumulatedDriving                  |                             |               | 434              | 434 |               |   |   |   |   |          |          |
| gnssADPointerNewestRecord               |                             |               | 2                | 2   | {00 00}       |   |   |   |   |          |          |
| gnssAccumulatedDrivingRecords           |                             |               | 432              | 432 |               |   |   |   |   |          |          |
| GNSSAccumulatedDrivingRecord            |                             | n8            | 18               | 18  |               |   |   |   |   |          |          |
| timeStamp                               |                             |               | 4                | 4   | {00..00}      |   |   |   |   |          |          |
| gnssPlaceRecord                         |                             |               | 14               | 14  |               |   |   |   |   |          |          |
| timeStamp                               |                             |               | 4                | 4   | {00..00}      |   |   |   |   |          |          |
| gnssAccuracy                            |                             |               | 1                | 1   | {00}          |   |   |   |   |          |          |
| geoCoordinates                          | vehicleOdometerValue        |               |                  |     |               | 6 | 3 | 6 | 3 | {00..00} | {00..00} |
| EF Specific_Conditions                  |                             |               | 22               | 22  |               |   |   |   |   |          |          |
| SpecificConditions                      |                             |               | 22               | 22  |               |   |   |   |   |          |          |
| conditionPointerNewestRecord            |                             |               | 2                | 2   | {00 00}       |   |   |   |   |          |          |
| specificConditionRecords                |                             |               | 20               | 20  |               |   |   |   |   |          |          |
| SpecificConditionRecord                 | n9                          |               | 5                | 5   |               |   |   |   |   |          |          |
| entryTime                               |                             |               | 4                | 4   | {00..00}      |   |   |   |   |          |          |
| specificConditionType                   |                             |               | 1                | 1   | {00}          |   |   |   |   |          |          |
| EF Application_Identification_V2        |                             |               | 10               | 10  |               |   |   |   |   |          |          |
| WorkshopCardApplicationIdentificationV2 |                             |               | 10               | 10  |               |   |   |   |   |          |          |

| File / Data Element                  | No of<br>Records                     | Size (bytes) |             | Default Value |    |    |    |  |
|--------------------------------------|--------------------------------------|--------------|-------------|---------------|----|----|----|--|
| LengthOfFollowingData                |                                      | 2            | 2           | {00 00}       |    |    |    |  |
| noOfBorderCrossingRecords            |                                      | 2            | 2           | {00 00}       |    |    |    |  |
| noOfLoadUnloadRecords                |                                      | 2            | 2           | {00 00}       |    |    |    |  |
| noofLoadTypeEntryRecords             |                                      | 2            | 2           | {00 00}       |    |    |    |  |
| VuConfigurationLengthRange           |                                      | 2            | 2           | {00 00}       |    |    |    |  |
| EF Places Authentication             |                                      | 42           | 42          |               |    |    |    |  |
| CardPlaceAuthDailyWorkPeriod         |                                      | 42           | 42          |               |    |    |    |  |
| placeAuthPointerNewestRecord         |                                      | 2            | 2           | {00 00}       |    |    |    |  |
| placeAuthStatusRecords               |                                      | 40           | 40          |               |    |    |    |  |
| PlaceAuthStatusRecord                | n4                                   | 5            | 5           |               |    |    |    |  |
| entryTime                            |                                      | 4            | 4           | {00..00}      |    |    |    |  |
| authenticationStatus                 |                                      | 1            | 1           | {00}          |    |    |    |  |
| EF GNSS_Places Authentication        |                                      | 122          | 122         |               |    |    |    |  |
| GNSSAuthAccumulatedDriving           |                                      | 122          | 122         |               |    |    |    |  |
| gnssAuthADPointerNewestRecord        |                                      | 2            | 2           | {00 00}       |    |    |    |  |
| gnssAuthStatusADRecords              |                                      | 120          | 120         |               |    |    |    |  |
| GNSSAuthStatusADRecord               | n8                                   | 5            | 5           |               |    |    |    |  |
| timeStamp                            |                                      | 4            | 4           | {00..00}      |    |    |    |  |
| authenticationStatus                 |                                      | 1            | 1           | {00}          |    |    |    |  |
| EF Border Crossings                  |                                      | 70           | 70          |               |    |    |    |  |
| CardBorderCrossings                  |                                      | 70           | 70          |               |    |    |    |  |
| borderCrossingPointerNewestRecord    |                                      | 2            | 2           | {00 00}       |    |    |    |  |
| cardBorderCrossingRecords            |                                      | 68           | 68          |               |    |    |    |  |
| CardBorderCrossingRecord             | n10                                  | 17           | 17          |               |    |    |    |  |
| countryLeft                          |                                      | 1            | 1           | {00}          |    |    |    |  |
| countryEntered                       |                                      | 1            | 1           | {00}          |    |    |    |  |
| gnssPlaceAuthRecord                  |                                      | 12           | 12          |               |    |    |    |  |
| timeStamp                            |                                      | 4            | 4           | {00..00}      |    |    |    |  |
| gnssAccuracy                         |                                      | 1            | 1           | {00}          |    |    |    |  |
| geoCoordinates                       |                                      | 6            | 6           | {00..00}      |    |    |    |  |
| authenticationStatus                 |                                      | 1            | 1           | {00}          |    |    |    |  |
| vehicleOdometerValue                 |                                      | 3            | 3           | {00..00}      |    |    |    |  |
| EF Load_Unload Operations            |                                      | 162          | 162         |               |    |    |    |  |
| CardLoadUnloadOperations             |                                      | 162          | 162         |               |    |    |    |  |
| loadUnloadPointerNewestRecord        |                                      | 2            | 2           | {00 00}       |    |    |    |  |
| cardloadUnloadRecords                |                                      | 160          | 160         |               |    |    |    |  |
| CardLoadUnloadRecord                 | n11                                  | 20           | 20          |               |    |    |    |  |
| timestamp                            |                                      | 4            | 4           | {00}          |    |    |    |  |
| operationType                        |                                      | 1            | 1           | {00..00}      |    |    |    |  |
| gnssPlaceAuthRecord                  |                                      | 12           | 12          |               |    |    |    |  |
| timeStamp                            |                                      | 4            | 4           | {00..00}      |    |    |    |  |
| gnssAccuracy                         |                                      | 1            | 1           | {00}          |    |    |    |  |
| geoCoordinates                       |                                      | 6            | 6           | {00..00}      |    |    |    |  |
| authenticationStatus                 |                                      | 1            | 1           | {00}          |    |    |    |  |
| vehicleOdometerValue                 |                                      | 3            | 3           | {00..00}      |    |    |    |  |
| EF Load Type_Entries                 |                                      | 22           | 22          |               |    |    |    |  |
| CardLoadTypeEntries                  |                                      | 22           | 22          |               |    |    |    |  |
| loadtypeEntryPointerNewestRecord     |                                      | 2            | 2           | {00 00}       |    |    |    |  |
| cardLoadTypeEntryRecords             |                                      | 20           | 20          |               |    |    |    |  |
| CardLoadTypeEntryRecord              | n12                                  | 5            | 5           |               |    |    |    |  |
| timestamp                            |                                      | 4            | 4           | {00..00}      |    |    |    |  |
| loadTypeEntered                      |                                      | 1            | 1           | {00}          |    |    |    |  |
| EF Calibration Add Data              |                                      | 6887         | 6887        |               |    |    |    |  |
| WorkshopCardCalibrationAddData       |                                      | 6887         | 6887        |               |    |    |    |  |
| calibrationPointerNewestRecord       |                                      | 2            | 2           | {00 00}       |    |    |    |  |
| workshopCardCalibrationA             |                                      |              |             |               |    |    |    |  |
| ddDataRecords                        |                                      | 6885         | 6885        |               |    |    |    |  |
| WorkshopCardCalibrationAddDataRecord | WorkshopCardCalibrationAddDataRecord | n5           | n5          | 27            | 27 | 27 | 27 |  |
| oldTimeValue                         |                                      | 4            | 4 {00..00}  |               |    |    |    |  |
| vehicleIdentificationNumber          |                                      | 17           | 17 {20..20} |               |    |    |    |  |
| byDefaultLoadType                    |                                      | 1            | 1 {00}      |               |    |    |    |  |
| calibrationCountry                   |                                      | 1            | 1 {00}      |               |    |    |    |  |
| calibrationCountryTimestamp          |                                      | 4            | 4 {00..00}  |               |    |    |    |  |
| EF VU_Configuration                  |                                      | 3072         | 3072        |               |    |    |    |  |
| VuConfigurations                     | n13                                  | 3072         | 3072        |               |    |    |    |  |

TCS\_163 The following values, used to provide sizes in the table above, are the minimum and maximum record number values the workshop card data structure must use for a generation 2 application:

|     |                              | Min                                      | Max                                      |
|-----|------------------------------|------------------------------------------|------------------------------------------|
| n1  | NoOfEventsPerType            | 3                                        | 3                                        |
| n2  | NoOfFaultsPerType            | 6                                        | 6                                        |
| n3  | NoOfCardVehicleRecords       | 8                                        | 8                                        |
| n4  | NoOfCardPlaceRecords         | 8                                        | 8                                        |
| n5  | NoOfCalibrationRecords       | 255                                      | 255                                      |
| n6  | CardActivityLengthRange      | 492 bytes (1 day * 240 activity changes) | 492 bytes (1 day * 240 activity changes) |
| n7  | NoOfCardVehicleUnitRecords   | 8                                        | 8                                        |
| n8  | NoOfGNSSADRecords            | 24                                       | 24                                       |
| n9  | NoOfSpecificConditionRecords | 4                                        | 4                                        |
| n10 | NoOfBorderCrossingRecords    | 4                                        | 4                                        |
| n11 | NoOfLoadUnloadRecords        | 8                                        | 8                                        |
| n12 | NoOfLoadTypeEntryRecords     | 4                                        | 4                                        |
| n13 | VuConfigurationLengthRange   | 3072 Bytes                               | 3072 Bytes                               |

**▼B**

#### 4.4. **Control card applications**

- 4.4.1 *Control Card application generation 1*
  - TCS\_164 After its personalisation, the control card application generation 1 shall have the following permanent file structure and file access rules:

| File                          | File ID | Read | Select | Update |
|-------------------------------|---------|------|--------|--------|
| DF Tachograph                 | '0500h' |      |        |        |
| EF Application_Identification | '0501h' | SC2  | SC1    | NEV    |
| EF Card_Certificate           | 'C100h' | SC2  | SC1    | NEV    |
| EF CA_Certificate             | 'C108h' | SC2  | SC1    | NEV    |
| EF Identification             | '0520h' | SC6  | SC1    | NEV    |
| EF Controller_Activity_Data   | '050Ch' | SC2  | SC1    | SC3    |

The following abbreviations for the security conditions are used in this table:

- **SC1** ALW OR SM-MAC-G2
- **SC2** ALW OR SM-MAC-G1 OR SM-MAC-G2
- **SC3** SM-MAC-G1 OR SM-MAC-G2
- **SC6** EXT-AUT-G1 OR SM-MAC-G1 OR SM-MAC-G2
- TCS\_165 All EF structures shall be transparent.
- TCS\_166 The control card application generation 1 shall have the following data structure:

| File / Data element                   | No of<br>Records | Size (Bytes)<br>Min | Size (Bytes)<br>Max |
|---------------------------------------|------------------|---------------------|---------------------|
| DF Tachograph                         |                  | 11186               | 24526               |
| └EF Application Identification        |                  | 5                   | 5                   |
| └ControlCardApplicationIdentification |                  | 5                   | 5                   |
| └typeOfTachographCardId               |                  | 1                   | 1 {00}              |
| └cardStructureVersion                 |                  | 2                   | 2 {00 00}           |
| └noOfControlActivityRecords           |                  | 2                   | 2 {00 00}           |
| └EF Card Certificate                  |                  | 194                 | 194                 |
| └CardCertificate                      |                  | 194                 | 194 {00..00}        |
| └EF CA Certificate                    |                  | 194                 | 194                 |
| └MemberStateCertificate               |                  | 194                 | 194 {00..00}        |
| └EF Identification                    |                  | 211                 | 211                 |
| └CardIdentification                   |                  | 65                  | 65                  |
| └cardIssuingMemberState               |                  | 1                   | 1 {00}              |
| └cardNumber                           |                  | 16                  | 16 {20..20}         |
| └cardIssuingAuthorityName             |                  | 36                  | 36 {00, 20..2}      |
| └cardIssueDate                        |                  | 4                   | 4 {00..00}          |
| └cardValidityBegin                    |                  | 4                   | 4 {00..00}          |
| └cardExpiryDate                       |                  | 4                   | 4 {00..00}          |
| └ControlCardHolderIdentification      |                  | 146                 | 146                 |
| └controlBodyName                      |                  | 36                  | 36 {00, 20..2}      |
| └controlBodyAddress                   |                  | 36                  | 36 {00, 20..2}      |
| └cardHolderName                       |                  |                     |                     |
| └holderSurname                        |                  | 36                  | 36 {00, 20..2}      |
| └holderFirstNames                     |                  | 36                  | 36 {00, 20..2}      |
| └cardHolderPreferredLanguage          |                  | 2                   | 2 {20 20}           |
| └EF Controller Activity Data          |                  | 10582               | 23922               |
| └ControlCardControlActivityData       |                  | 10582               | 23922               |
| └controlPointerNewestRecord           |                  | 2                   | 2 {00 00}           |
| └controlActivityRecords               |                  | 10580               | 23920               |
| └controlActivityRecord                | n7               | 46                  | 46                  |
| └controlType                          |                  | 1                   | 1 {00}              |
| └controlTime                          |                  | 4                   | 4 {00..00}          |
| └controlledCardNumber                 |                  |                     |                     |
| └cardType                             |                  | 1                   | 1 {00}              |
| └cardIssuingMemberState               |                  | 1                   | 1 {00}              |
| └cardNumber                           |                  | 16                  | 16 {20..20}         |
| └controlledVehicleRegistration        |                  |                     |                     |
| └vehicleRegistrationNation            |                  | 1                   | 1 {00}              |
| └vehicleRegistrationNumber            |                  | 14                  | 14 {00, 20..2}      |
| └controlDownloadPeriodBegin           |                  | 4                   | 4 {00..00}          |
| └controlDownloadPeriodEnd             |                  | 4                   | 4 {00..00}          |

#### TCS\_167 The following values, used to provide sizes in the table above, are the minimum and maximum record number values the control card data structure must use for a generation 1 application:

|    |                            | Min | Max |
|----|----------------------------|-----|-----|
| n7 | NoOfControlActivityRecords | 230 | 520 |

#### 4.4.2 *Control card application generation 2*

TCS\_168 After its personalisation, the control card application generation 2 shall have the following permanent file structure and file access rules.

### **B**

*Notes:*

- the short EF identifier SFID is given as decimal number, e.g. the value 30 corresponds to 11110 in binary,
- EF Application\_Identification\_V2, and EF VU\_Configuration are only present in version 2 of the generation 2 control card,
- cardStructureVersion in EF Application\_Identification is equal to {01 01} for version 2 of the generation 2 control card, while it was equal to {01 00} for version 1 of the generation 2 control card.

| File                              | File ID | SFID | Access rules  |           |
|-----------------------------------|---------|------|---------------|-----------|
| DF Tachograph_G2                  |         |      | Read / Select | Update    |
| -EF Application_Identification    | '0501h' | 1    | SC1           | NEV       |
| -EF CardMA_Certificate            | 'C100h' | 2    | SC1           | NEV       |
| -EF CA_Certificate                | 'C108h' | 4    | SC1           | NEV       |
| -EF Link_Certificate              | 'C109h' | 5    | SC1           | NEV       |
| -EF Identification                | '0520h' | 6    | SC1           | NEV       |
| -EF Controller_Activity_Data      | '050Ch' | 14   | SC1           | SM-MAC-G2 |
| -EF Application_Identification_V2 | '0525h' | 22   | SC1           | NEV       |
| -EF VU_Configuration              | '0540h' | 30   | SC5/SC1       | SM-MAC-G2 |

The following abbreviations for the security condition are used in this table:

- **SC1** ALW OR SM-MAC-G2
- **SC5** For the Read Binary command with even INS byte: SM-C-MAC-G2 AND SM-R-ENC-MAC-G2

For the Read Binary command with odd INS byte (if supported): NEV

**▼B**

- TCS\_169 All EF structures shall be transparent.
- TCS\_170 The control card application generation2 shall have the following data structure:

| File / Date |
|-------------|
| Element     |

| File / Data<br>Element               | No of Records | Min   | Max   | Default Values |
|--------------------------------------|---------------|-------|-------|----------------|
| DF Tachograph G2                     |               | 14486 | 28237 |                |
| EF Application_Identification        |               | 5     | 5     |                |
| ControlCardApplicationIdentification |               | 5     | 5     |                |
| typeOfTachographCardId               |               | 1     | 1     | {00}           |
| cardStructureVersion                 |               | 2     | 2     | {01 01} V2     |
| noOfControlActivityRecords           |               | 2     | 2     | {00 00}        |
| EF CardMA_Certificate                |               | 204   | 341   |                |
| CardMA_Certificate                   |               | 204   | 341   | {00..00}       |
| EF CA_Certificate                    |               | 204   | 341   |                |
| MemberStateCertificate               |               | 204   | 341   | {00..00}       |
| EF Link_Certificate                  |               | 204   | 341   |                |
| LinkCertificate                      |               | 204   | 341   | {00..00}       |
| EF Identification                    |               | 211   | 211   |                |
| CardIdentification                   |               | 65    | 65    |                |
| cardIssuingMemberState               |               | 1     | 1     | {00}           |
| cardNumber                           |               | 16    | 16    | {20..20}       |
| cardIssuingAuthorityName             |               | 36    | 36    | {00, 20..20}   |
| cardIssueDate                        |               | 4     | 4     | {00..00}       |
| cardValidityBegin                    |               | 4     | 4     | {00..00}       |
| cardExpiryDate                       |               | 4     | 4     | {00..00}       |
| ControlCardHolderIdentification      |               | 146   | 146   |                |
| controlBodyName                      |               | 36    | 36    | {00, 20..20}   |
| controlBodyAddress                   |               | 36    | 36    | {00, 20..20}   |
| cardHolderName                       |               |       |       |                |
| holderSurname                        |               | 36    | 36    | {00, 20..20}   |
| holderFirstNames                     |               | 36    | 36    | {00, 20..20}   |
| cardHolderPreferredLanguage          |               | 2     | 2     | {20 20}        |
| EF Controller_Activity_Data          |               | 10582 | 23922 |                |
| ControlCardControlActivityData       |               | 10582 | 23922 |                |
| controlPointerNewestRecord           |               | 2     | 2     | {00 00}        |
| controlActivityRecords               |               | 10580 | 23920 |                |
| controlActivityRecord                | n7            | 46    | 46    |                |
| controlType                          |               | 1     | 1     | {00}           |
| controlTime                          |               | 4     | 4     | {00..00}       |
| controlledCardNumber                 |               |       |       |                |
| cardType                             |               | 1     | 1     | {00}           |

**▼B**

TCS\_171 The following values, used to provide sizes in the table above, are the minimum and maximum record number values the control card data structure must use for a generation 2 application:

|     |                            | Min        | Max        |
|-----|----------------------------|------------|------------|
| n7  | NoOfControlActivityRecords | 230        | 520        |
| n13 | VuConfigurationLengthRange | 3072 Bytes | 3072 Bytes |

### **B**

#### 4.5. **Company card applications**

4.5.1 *Company card application generation 1*

TCS\_172 After its personalisation, the company card application generation 1 shall have the following permanent file structure and file access rules:

| File                          | File ID | Read       | Select | Update |
|-------------------------------|---------|------------|--------|--------|
| DF Tachograph                 | '0500h' |            | SC1    |        |
| EF Application_Identification | '0501h' | SC2        | SC1    | NEV    |
| EF Card_Certificate           | 'C100h' | SC2        | SC1    | NEV    |
| EF CA_Certificate             | 'C108h' | SC2        | SC1    | NEV    |
| EF Identification             | '0520h' | <b>SC6</b> | SC1    | NEV    |
| EF Company Activity Data      | '050Dh' | SC2        | SC1    | SC3    |

The following abbreviations for the security conditions are used in this table:

- **SC1** ALW OR SM-MAC-G2
- **SC2** ALW OR SM-MAC-G1 OR SM-MAC-G2
- **SC3** SM-MAC-G1 OR SM-MAC-G2
- **SC6** EXT-AUT-G1 OR SM-MAC-G1 OR SM-MAC-G2

TCS\_173 All EF structures shall be transparent.

#### TCS\_174 The company card application generation 1 shall have the following data structure:

| File / Data element                   | No of<br>Records | Size (bytes) |       | Default<br>Values |
|---------------------------------------|------------------|--------------|-------|-------------------|
| DF Tachograph                         |                  | 11114        | 24454 |                   |
| EF Application_Identification         |                  | 5            | 5     |                   |
| └CompanyCardApplicationIdentification |                  | 5            | 5     |                   |
| —typeOfTachographCardId               |                  | 1            | 1     | {00}              |
| — cardStructureVersion                |                  | 2            | 2     | {00 00}           |
| —noOfCompanyActivityRecords           |                  | 2            | 2     | {00 00}           |
| EF Card Certificate                   |                  | 194          | 194   |                   |
| └─CardCertificate                     |                  | 194          | 194   | {00..00}          |
| EF CA Certificate                     |                  | 194          | 194   |                   |
| └MemberStateCertificate               |                  | 194          | 194   | {00..00}          |
| EF Identification                     |                  | 139          | 139   |                   |
| └CardIdentification                   |                  | 65           | 65    |                   |
| - cardIssuingMemberState              |                  | 1            | 1     | {00}              |
| -cardNumber                           |                  | 16           | 16    | {20..20}          |
| — cardIssuingAuthorityName            |                  | 36           | 36    | {00, 20..20}      |
| -cardIssueDate                        |                  | 4            | 4     | {00..00}          |
| - cardValidityBegin                   |                  | 4            | 4     | {00..00}          |
| -cardExpiryDate                       |                  | 4            | 4     | {00..00}          |
| CompanyCardHolderIdentification       |                  | 74           | 74    |                   |
| └companyName                          |                  | 36           | 36    | {00, 20..20}      |
| └ companyAddress                      |                  | 36           | 36    | {00, 20..20}      |
| └cardHolderPreferredLanguage          |                  | 2            | 2     | {20 20}           |
| EF Company_Activity_Data              |                  | 10582        | 23922 |                   |
| L_CompanyActivityData                 |                  | 10582        | 23922 |                   |
| -companyPointerNewestRecord           |                  | 2            | 2     | {00 00}           |
| └companyActivityRecords               |                  | 10580        | 23920 |                   |
| └ companyActivityRecord               | nB               | 46           | 46    |                   |
| - companyActivityType                 |                  | 1            | 1     | {00}              |
| - companyActivityTime                 |                  | 4            | 4     | {00..00}          |
| - cardNumberInformation               |                  |              |       |                   |
| └cardType                             |                  | 1            | 1     | {00}              |
| └ cardIssuingMemberState              |                  | 1            | 1     | {00}              |
| └cardNumber                           |                  | 16           | 16    | {20..20}          |
| - vehicleRegistrationInformation      |                  |              |       |                   |
| └vehicleRegistrationNation            |                  | 1            | 1     | {00}              |
| └vehicleRegistrationNumber            |                  | 14           | 14    | {00, 20..20}      |
| - downloadPeriodBegin                 |                  | 4            | 4     | {00..00}          |
| downloadPeriodEnd                     |                  | 4            | 4     | {00..00}          |

TCS\_175 The following values, used to provide sizes in the table above, are the minimum and maximum record number values the company card data structure must use for a generation 1 application:

|    |                            | Min | Max |
|----|----------------------------|-----|-----|
| n8 | NoOfCompanyActivityRecords | 230 | 520 |

4.5.2 *Company card application generation 2*

TCS\_176 After its personalisation, the company card application generation 2 shall have the following permanent file structure and file access rules.

*Notes:*

- the short EF identifier SFID is given as decimal number, e.g. the value 30 corresponds to 11110 in binary,
- EF Application\_Identification\_V2, and EF VU\_Configuration are only present in version 2 of the generation 2 company card,
- cardStructureVersion in EF Application\_Identification is equal to {01 01} for version 2 of the generation 2 company card, while it was equal to {01 00} for version 1 of the generation 2 company card.

|  | File                              | File ID | SFID | Read / Select | Update    |
|--|-----------------------------------|---------|------|---------------|-----------|
|  | └─DF Tachograph_G2                |         |      | SC1           |           |
|  | └EF Application_Identification    | '0501h' | 1    | SC1           | NEV       |
|  | └EF CardMA_Certificate            | 'C100h' | 2    | SC1           | NEV       |
|  | └EF CA_Certificate                | 'C108h' | 4    | SC1           | NEV       |
|  | └EF Link_Certificate              | 'C109h' | 5    | SC1           | NEV       |
|  | └EF Identification                | '0520h' | 6    | SC1           | NEV       |
|  | └EF Company_Activity_Data         | '050Dh' | 14   | SC1           | SM-MAC-G2 |
|  | └EF Application_Identification_V2 | '0525h' | 22   | SC1           | NEV       |
|  | └EF VU_Configuration              | '0540h' | 30   | SC5/SC1       | SM-MAC-G2 |

The following abbreviations for the security condition are used in this table:

- **SC1** ALW OR SM-MAC-G2
- **SC5** For the Read Binary command with even INS byte: SM-C-MAC-G2 AND SM-R-ENC-MAC-G2

For the Read Binary command with odd INS byte (if supported):NEV

- **▼B**
- TCS\_177 All EF structures shall be transparent.
- TCS\_178 The company card application generation 2 shall have the following data structure:

| File / Data Element                    | No of Records | Min   | Max   | Default Values |
|----------------------------------------|---------------|-------|-------|----------------|
| DF Tachograph_G2                       |               | 14414 | 28165 |                |
| EF Application_Identification          |               | 5     | 5     |                |
| CompanyCardApplicationIdentification   |               | 5     | 5     |                |
| typeOfTachographCardId                 |               | 1     | 1     | {00}           |
| cardStructureVersion                   |               | 2     | 2     | {01 01} V2     |
| noOfCompanyActivityRecords             |               | 2     | 2     | {00.00}        |
| EF CardMA Certificate                  |               | 204   | 341   |                |
| CardMA_Certificate                     |               | 204   | 341   | {00..00}       |
| EF CA Certificate                      |               | 204   | 341   |                |
| MemberStateCertificate                 |               | 204   | 341   | {00..00}       |
| EF Link Certificate                    |               | 204   | 341   |                |
| LinkCertificate                        |               | 204   | 341   | {00..00}       |
| EF Identification                      |               | 139   | 139   |                |
| CardIdentification                     |               | 65    | 65    |                |
| cardIssuingMemberState                 |               | 1     | 1     | {00}           |
| cardNumber                             |               | 16    | 16    | {20..20}       |
| cardIssuingAuthorityName               |               | 36    | 36    | {00, 20..20}   |
| cardIssueDate                          |               | 4     | 4     | {00..00}       |
| cardValidityBegin                      |               | 4     | 4     | {00..00}       |
| cardExpiryDate                         |               | 4     | 4     | {00..00}       |
| CompanyCardHolderIdentification        |               | 74    | 74    |                |
| companyName                            |               | 36    | 36    | {00, 20..20}   |
| companyAddress                         |               | 36    | 36    | {00, 20..20}   |
| cardHolderPreferredLanguage            |               | 2     | 2     | {20 20}        |
| EF Company Activity _Data              |               | 10582 | 23922 |                |
| CompanyActivityData                    |               | 10582 | 23922 |                |
| companyPointerNewestRecord             |               | 2     | 2     | {00 00}        |
| companyActivityRecords                 |               | 10580 | 23920 |                |
| companyActivityRecord                  | n8            | 46    | 46    |                |
| companyActivityType                    |               | 1     | 1     | {00}           |
| companyActivityTime                    |               | 4     | 4     | {00..00}       |
| cardNumberInformation                  |               |       |       |                |
| cardType                               |               | 1     | 1     | {00}           |
| cardIssuingMemberState                 |               | 1     | 1     | {00}           |
| cardNumber                             |               | 16    | 16    | {20..20}       |
| vehicleRegistrationInformation         |               |       |       |                |
| vehicleRegistrationNation              |               | 1     | 1     | {00}           |
| vehicleRegistrationNumber              |               | 14    | 14    | {00, 20..20}   |
| download PeriodBegin                   |               | 4     | 4     | {00..00}       |
| downloadPeriodEnd                      |               | 4     | 4     | {00..00}       |
| EF Application Identification V2       |               | 4     | 4     |                |
| CompanyCardApplicationIdentificationV2 |               | 4     | 4     |                |
| lengthOfFollowingData                  |               | 2     | 2     | {00.00}        |
| VuConfigurationLengthRange             |               | 2     | 2     | {00.00}        |
| EF VuConfiguration                     |               | 3072  | 3072  |                |
| VuConfigurations                       | n13           | 3072  | 3072  |                |

**▼B**

TCS\_179 The following values, used to provide sizes in the table above, are the minimum and maximum record number values the company card data structure must use for a generation 2 application:

|     |                            | Min        | Max        |
|-----|----------------------------|------------|------------|
| n8  | NoOfCompanyActivityRecords | 230        | 520        |
| n13 | VuConfigurationLengthRange | 3072 Bytes | 3072 Bytes |