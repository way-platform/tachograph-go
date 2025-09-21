## 3.3. **Card Downloading**

**▼M3**

- DDP\_035 The download of a tachograph card includes the following steps:
  - Download the common information of the card in the EFs ICC and IC. This information is optional and is not secured with a digital signature.
  - For first and second generation tachograph cards
    - Download EFs within Tachograph DF:
      - Download the EFs Card\_Certificate and CA\_Certificate. This information is not secured with a digital signature.

It is mandatory to download these files for each download session.

- Download the other application data EFs (within Tachograph DF) except EF Card\_Download. This information is secured with a digital signature, using Appendix 11 Common Security Mechanisms Part A.
- It is mandatory to download at least the EFs Application\_Identification and Identification for each download session.
- When downloading a driver card it is also mandatory to download the following EFs:

Events\_Data,

Faults\_Data,

Driver\_Activity\_Data,

Vehicles\_Used,

Places,

Control\_Activity\_Data,

Specific\_Conditions.

— For second generation tachograph cards only:

- Except when a download of a driver card inserted in a VU is performed during drivers' control by a non EU control authority, using a first generation control card, download EFs within Tachograph\_G2 DF:
  - Download the EFs CardSignCertificate, CA\_Certificate and Link\_Certificate. This information is not secured with a digital signature.
  - It is mandatory to download these files for each download session.

- Download the other application data EFs (within Tachograph\_G2 DF) except EF Card\_Download. This information is secured with a digital signature, using Appendix 11 Common Security Mechanisms Part B.
- It is mandatory to download at least the EFs Application\_Identification, Application\_Identification\_V2 (if present) and Identification for each download session.
- When downloading a driver card it is also mandatory to download the following EFs:

Events\_Data,

Faults\_Data,

Driver\_Activity\_Data,

Vehicles\_Used,

Places,

Control\_Activity\_Data,

Specific\_Conditions,

VehicleUnits\_Used,

GNSS\_Places,

Places\_Authentication, if present,

GNSS\_Places\_Authentication, if present,

Border\_Crossings, if present,

Load\_Unload\_Operations, if present,

Load\_Type\_Entries, if present.

- When downloading a driver card, update the Last-CardDownload date in EF Card\_Download, in the Tachograph and, if applicable, Tachograph\_G2 DFs.
- When downloading a workshop card, reset the calibration counter in EF Card\_Download in the Tachograph and, if applicable, Tachograph\_G2 DFs.
- When downloading a workshop card the EF Sensor\_Installation\_Data in the Tachograph and, if applicable, Tachograph\_G2 DFs shall not be downloaded.

#### 3.3.1 *Initialisation sequence*

DDP\_036 The IDE shall initiate the sequence as follows:

| Card       | Direction | IDE/IFD        | Meaning/Remarks |
|------------|-----------|----------------|-----------------|
|            | ←         | Hardware reset |                 |
| <b>ATR</b> | →         |                |                 |

It is optional to use PPS to switch to a higher baud rate as long as the ICC supports it.

### **M3**

#### 3.3.2 *Sequence for un-signed data files*

DDP\_037 **►M1** The sequence to download EFs ICC, IC, Card\_Certificate (or CardSignCertificate for DF Tachograph\_G2), CA\_Certificate and Link\_Certificate (for DF Tachograph\_G2 only) is as follows: ◄

| Card            | Direction | IDE/IFD           | Meaning/Remarks                                                                                                                                           |
|-----------------|-----------|-------------------|-----------------------------------------------------------------------------------------------------------------------------------------------------------|
|                 | ⇦         | Select File       | Select by File identifiers                                                                                                                                |
| OK              | ⇨         |                   |                                                                                                                                                           |
|                 | ⇦         | Read Binary       | If the file contains more data<br>than the buffer size of the<br>reader or the card the<br>command has to be repeated<br>until the complete file is read. |
| File Data<br>OK | ⇨         | Store data to ESM | according to 3.4 Data storage<br>format                                                                                                                   |

*Note 1:* Before selecting the Card\_Certificate (or CardSign-Certificate) EF, the Tachograph Application must be selected (selection by AID).

*Note 2:* Selecting and reading a file may also be performed in one step using a Read Binary command with a short EF identifier.

#### 3.3.3 *Sequence for Signed data files*

DDP\_038 The following sequence shall be used for each of the following files that has to be downloaded with their signature:

**▼M1**

| Card                                                          | Dir                                    | IDE / IFD            | Meaning / Remarks                                                                                                                                                                                                  |
|---------------------------------------------------------------|----------------------------------------|----------------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
|                                                               | <span style="font-size: 1em;">⬅</span> | Select File          |                                                                                                                                                                                                                    |
| OK                                                            | <span>➡</span>                         |                      |                                                                                                                                                                                                                    |
|                                                               | <span>⬅</span>                         | Perform Hash of File | — Calculates the hash value<br>over the data content of<br>the selected file using the<br>prescribed hash algo-<br>rithm in accordance with<br>Appendix 11, part A or<br>B. This command is not<br>an ISO-Command. |
| Calculate Hash of File<br>and store Hash value<br>temporarily |                                        |                      |                                                                                                                                                                                                                    |

| Card                                                                                                       | Dir | IDE / IFD                                             | Meaning / Remarks                                                                                                                                              |
|------------------------------------------------------------------------------------------------------------|-----|-------------------------------------------------------|----------------------------------------------------------------------------------------------------------------------------------------------------------------|
| OK                                                                                                         | ⇒   |                                                       |                                                                                                                                                                |
|                                                                                                            | ⇐   | Read Binary                                           | If the file contains more data<br>than the buffer of the reader<br>or the card can hold, the<br>command has to be repeated<br>until the complete file is read. |
| File Data<br>OK                                                                                            | ⇒   | Store received data to ESM                            | according to 3.4 Data storage<br>format                                                                                                                        |
|                                                                                                            | ⇐   | PSO: Compute Digital<br>Signature                     |                                                                                                                                                                |
| Perform Security<br>Operation 'Compute<br>Digital Signature'<br>using the temporarily<br>stored Hash value |     |                                                       |                                                                                                                                                                |
| Signature<br>OK                                                                                            | ⇒   | Append data to the previous<br>stored data on the ESM | according to 3.4 Data storage<br>format                                                                                                                        |

**▼B**

*Note:* Selecting and reading a file may also be performed in one step using a Read Binary command with a short EF identifier. In this case the EF may be selected and read before the command Perform Hash of File is applied.

##### 3.3.4 *Sequence for resetting the calibration counter.*

DDP\_039 The sequence to reset the counter in the EF in a workshop card is the following:

| Card                           | Dir | IDE/IFD                                                       | Meaning/Remarks            |
|--------------------------------|-----|---------------------------------------------------------------|----------------------------|
|                                | ⇦   | Select File EF<br>Card_Download                               | Select by File identifiers |
| OK                             | ⇨   |                                                               |                            |
|                                | ⇦   | Update Binary<br>NoOfCalibrationsSince-<br>Download = '00 00' |                            |
| resets card download<br>number |     |                                                               |                            |
| OK                             | ⇨   |                                                               |                            |

*Note:* Selecting and updating a file may also be performed in one step using an Update Binary command with a short EF identifier.

# **M1**

### 3.4. **Data storage format**

- 3.4.1 *Introduction*
  - DDP\_040 The downloaded data has to be stored according to the following conditions:
    - The data shall be stored transparent. This means that the order of the bytes as well as the order of the bits inside the byte that are transferred from the card has to be preserved during storage.
    - All files of the card downloaded within a download session are stored in one file on the ESM.

#### 3.4.2 *File format*

DDP\_041 The file format is a concatenation of several TLV objects.

- DDP\_042 The tag for an EF shall be the FID plus the appendix '00'.
- DDP\_043 The tag of an EF's signature shall be the FID of the file plus the appendix '01'.
- DDP\_044 The length is a two byte value. The value defines the number of bytes in the value field. The value 'FF FF' in the length field is reserved for future use.
- DDP\_045 When a file is not downloaded nothing related to the file shall be stored (no tag and no zero length).

### **M1**

DDP\_046 A signature shall be stored as the next TLV object directly after the TLV object that contains the data of the file.

| Definition            | Meaning                                                                  | Length  |
|-----------------------|--------------------------------------------------------------------------|---------|
| FID (2 Bytes)    '00' | Tag for EF (FID) in the Tachograph or for common information of the card | 3 Bytes |
| FID (2 Bytes)    '01' | Tag for Signature of EF (FID) in the Tachograph DF                       | 3 Bytes |
| FID (2 Bytes)    '02' | Tag for EF (FID) in the Tachograph_G2 DF                                 | 3 Bytes |
| FID (2 Bytes)    '03' | Tag for Signature of EF (FID) in the Tachograph_G2 DF                    | 3 Bytes |
| xx xx                 | Length of Value field                                                    | 2 Bytes |

**▼M1**

| Tag      | Length | Value                                                 |
|----------|--------|-------------------------------------------------------|
| 00 02 00 | 00 11  | Data of EF ICC                                        |
| C1 00 00 | 00 C2  | Data of EF Card_Certificate                           |
|          |        | ...                                                   |
| 05 05 00 | 0A 2E  | Data of EF Vehicles_Used (in the Tachograph DF)       |
| 05 05 01 | 00 80  | Signature of EF Vehicles_Used (in the Tachograph DF)  |
| 05 05 02 | 0A 2E  | Data of EF Vehicles_Used in the Tachograph_G2 DF      |
| 05 05 03 | xx xx  | Signature of EF Vehicles_Used in the Tachograph_G2 DF |

Example of data in a download file on an ESM:

**▼B**

- 4. DOWNLOADING A TACHOGRAPH CARD VIA A VEHICLE UNIT.
  - DDP\_047 The VU must allow for downloading the content of a driver card inserted to a connected IDE.
  - DDP\_048 The IDE shall send a 'Transfer Data Request Card Download' message to the VU to initiate this mode (see 2.2.2.9).

**▼M1**

**▼B**

DDP\_049 First generation driver cards: Data shall be downloaded using the first generation data download protocol, and downloaded data shall have the same format as data downloaded from a first generation vehicle unit.

> Second generation driver cards: the VU shall then download the whole card, file by file, in accordance with the card downloading protocol defined in paragraph 3, and forward all data received from the card to the IDE within the appropriate TLV file format (see 3.4.2) and encapsulated within a 'Positive Response Transfer Data' message.

- DDP\_050 The IDE shall retrieve card data from the 'Positive Response Transfer Data' message (stripping all headers, SIDs, TREPs, sub message counters, and checksums) and store them within one single physical file as described in paragraph 2.3.
- DDP\_051 The VU shall then, as applicable, update the or the file of the driver card.