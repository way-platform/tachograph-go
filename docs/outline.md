# Tachograph Regulation Documentation Outline

This document provides a structured outline of the tachograph regulation, summarizing key aspects of each section. It is intended for developers who need to understand the tachograph system, its data structures, and communication protocols to implement data parsing and processing solutions.

---

## Part 1: Main Body of the Regulation

### [Article 1: Subject Matter and Scope](regulation/02-appendix-subject-subject-matter-and-scope.html)

- **Summary**: This article establishes the scope of the regulation, which covers the uniform application of rules for tachographs.
- **Key Content**:
    - Recording of vehicle position.
    - Remote early detection of manipulation.
    - Interface with Intelligent Transport Systems (ITS).
    - Requirements for type-approval, construction, testing, installation, and security.
- **Developer Takeaway**: This section sets the overall context. The technical details for implementation are found in the Annexes and Appendices, particularly Annex IC.
- **Key Points (in order)**:
    - The regulation lays down provisions for the uniform application of rules regarding:
        - Recording of vehicle position at certain points.
        - Remote early detection of possible manipulation or misuse.
        - Interface with intelligent transport systems.
        - Administrative and technical requirements for type-approval procedures.
    - The construction, testing, installation, inspection, operation, and repair of smart tachographs must comply with Annex IC.
    - Tachographs other than smart tachographs must continue to comply with Annex I to Regulation (EU) No 165/2014 or Annex IB to Council Regulation (EEC) No 3821/85.
    - The remote early detection facility shall also transmit weight data for fraud detection.

### [Article 2: Definitions](regulation/03-appendix-definitions-definitions.html)

- **Summary**: Provides core definitions for terms used throughout the regulation.
- **Key Content**: Defines terms like `digital tachograph`, `smart tachograph`, `vehicle unit (VU)`, `tachograph card`, and other key components.
- **Developer Takeaway**: Essential reading. Understanding these definitions is fundamental to correctly interpreting the technical specifications. Refer to `docs/definitions.md` for a consolidated list.
- **Key Points (in order)**:
    - In addition to definitions in Regulation (EU) No 165/2014, this article defines:
        - `digital tachograph` or `first generation tachograph`: A digital tachograph other than a smart tachograph.
        - `external GNSS facility`: A facility containing the GNSS receiver when it is not integrated into the main vehicle unit.
        - `information folder` & `information package`: Documentation related to the type-approval process.
        - `remote early detection facility`: Equipment used for targeted roadside checks.
        - `smart tachograph` or `second generation tachograph`: A tachograph complying with Articles 8, 9, and 10 of Regulation (EU) No 165/2014 and Annex IC.
        - `tachograph component`: Includes the vehicle unit, motion sensor, record sheet, external GNSS facility, and external remote early detection facility.
        - `vehicle unit (VU)`: The tachograph excluding the motion sensor. It may be a single unit or distributed units and includes a processing unit, data memory, time measurement, card interfaces, printer, display, and other facilities.

### [Articles 3-6: Procedures and Timeline](regulation/04-appendix-location-based-location-based-services.html)

- **Summary**: These articles cover location-based services, type-approval procedures, modifications to approvals, and the entry into force of the regulation.
- **Key Content**:
    - **Article 3**: Mandates compatibility with Galileo and EGNOS for positioning.
    - **Article 4**: Outlines the administrative process for getting a tachograph or component type-approved.
    - **Article 5**: Describes how modifications to existing type-approvals are handled (revisions vs. extensions).
    - **Article 6**: Specifies the dates from which the different parts of the regulation apply.
- **Developer Takeaway**: These sections provide context on the administrative and legal framework. The dates in Article 6 are important for understanding which generation of equipment and which set of rules apply to a given vehicle.
- **Key Points (in order)**:
    - **Art. 3**: Smart tachographs must be compatible with Galileo and EGNOS positioning services.
    - **Art. 4**: A manufacturer must submit an application with an information folder to a designated type-approval authority to receive a type-approval certificate.
    - **Art. 5**: The manufacturer must inform the type-approval authority of any modification in software or hardware. Modifications can lead to a 'revision', an 'extension', or a 'new type-approval'.
    - **Art. 6**: The regulation applies from 2 March 2016, but Annex IC applies from 15 June 2019.

---

## Part 2: Annexes and Appendices (Technical Specifications)

### [ANNEX I C: Core Technical Requirements](regulation/08-annex-1c-requirements-for-construction-testing-installation-and-inspection.html)

- **Summary**: This is the most critical annex, detailing the construction and functional requirements for the entire smart tachograph system.
- **Developer Takeaway**: This is the primary reference for any developer. Pay close attention to Section 3 (Construction and Functional Requirements) and Section 4 (Tachograph Cards) as they define the data you will be parsing.
- **Key Points (in order)**:
    - **Section 1: DEFINITIONS**: Contains a large set of technical definitions crucial for implementation (e.g., `w`, `k`, `l`, `vehicle characteristic constant`, `tyre size`, `UTC time`).
    - **Section 2: GENERAL CHARACTERISTICS AND FUNCTIONS**:
        - Defines the four modes of operation: `operational`, `control`, `calibration`, `company`.
        - Specifies data access rights for each mode, which is fundamental for security and privacy.
    - **Section 3: CONSTRUCTION AND FUNCTIONAL REQUIREMENTS FOR RECORDING EQUIPMENT**:
        - `3.2`: Mandates recording of speed, distance, and position via GNSS.
        - `3.4`: Defines the four main driver activities: `DRIVING`, `WORK`, `AVAILABILITY`, `BREAK/REST`.
        - `3.6`: Specifies manual entries for start/end of work period and special conditions like `OUT_OF_SCOPE` and `FERRY/TRAIN_CROSSING`.
        - `3.9`: Lists all events and faults that the VU must detect and record (e.g., `Card conflict`, `Time overlap`, `Power supply interruption`, `Motion data error`).
        - `3.12`: Details the data structures to be stored in the VU's permanent memory. This is a critical reference for parsing downloaded VU data. It covers equipment identification, security data, driver activities, locations, odometer readings, detailed speed, events, faults, and calibration data.
        - `3.14`: Specifies which data must be recorded onto inserted tachograph cards.
    - **Section 4: CONSTRUCTION AND FUNCTIONAL REQUIREMENTS FOR TACHOGRAPH CARDS**:
        - `4.1`: Defines the visible data that must be printed on the physical card.
        - `4.5`: Details the data storage structure on the card, organized into Elementary Files (EFs) within a Dedicated File (DF). This section is essential for parsing data read directly from a card. It specifies the file identifiers (FIDs) and data content for each file on each of the four card types.
    - **Section 5 & 6**: Cover installation and inspection procedures, including the requirements for sealing the equipment to prevent tampering.

### [Appendix 1: Data Dictionary](regulation/09-appendix-1-data-dictionary.html)

- **Summary**: The "Rosetta Stone" for tachograph data. It defines every single data type, its structure, encoding, and size.
- **Developer Takeaway**: **This is the most important document for a data parsing implementation.** Every field in a downloaded file or a real-time data stream is defined here. You will need to map your parsing logic directly to these definitions.
- **Key Points (in order)**:
    - Provides an alphabetical list of all data types used in the tachograph system.
    - For each data type, it specifies:
        - The data type definition (e.g., `OCTET STRING`, `INTEGER`, `BCDString`).
        - The size in bytes.
        - A detailed description of its content and encoding rules.
    - Examples of critical data types defined here include `TimeReal`, `VehicleIdentificationNumber`, `CardActivityDailyRecord`, `GNSSPlaceRecord`, and `FullCardNumber`.

### [Appendix 2: Tachograph Cards Specification](regulation/10-appendix-2-tachograph-cards-specification.html)

- **Summary**: Details the technical specifications of the four types of tachograph cards (Driver, Workshop, Control, Company).
- **Developer Takeaway**: Essential for anyone building a tool to read data directly from a tachograph card. You must implement the specified APDU commands and understand the file structure (`DF_Tachograph`, `EF_Application_Identification`, etc.) to navigate and extract data.
- **Key Points (in order)**:
    - **Section 2: ELECTRICAL AND PHYSICAL CHARACTERISTICS**: Specifies voltage, clock frequency, and other physical properties compliant with ISO/IEC 7816.
    - **Section 3: HARDWARE AND COMMUNICATION**:
        - Mandates support for both T=0 and T=1 transmission protocols.
        - Defines the structure of the Answer to Reset (ATR).
        - Lists the required APDU commands (e.g., `SELECT FILE`, `READ BINARY`, `UPDATE BINARY`, `VERIFY`, `INTERNAL AUTHENTICATE`, `GET CHALLENGE`).
        - Defines the access rules (security conditions) for each command and file (e.g., `ALW`, `NEV`, `PWD`, `SM-MAC`).
    - **Section 4: TACHOGRAPH CARDS STRUCTURE**:
        - Describes the hierarchical file system: Master File (MF) contains Dedicated Files (DFs), which in turn contain Elementary Files (EFs).
        - Provides detailed file structures, File IDs (FIDs), and access conditions for each of the four card types, covering both Gen1 and Gen2 applications on the same card.

### [Appendix 3: Pictograms](regulation/11-appendix-3-pictograms.html)

- **Summary**: Provides a visual reference for all pictograms used on the tachograph display and in printouts.
- **Developer Takeaway**: Useful if your application needs to display information in a user-friendly way that mimics the tachograph itself. You can map the binary codes from the data files to these images or textual descriptions.
- **Key Points (in order)**:
    - Lists basic pictograms for people (Company, Controller, Driver), activities (Available, Driving, Rest, Work), equipment (Card, Clock, Display, Printer), and miscellaneous concepts (Events, Faults, Location, Time).
    - Defines pictogram combinations used to represent more complex information, such as "Crew driving" or "Driver activities from card daily printout".

### [Appendix 4: Printouts](regulation/12-appendix-4-printouts.html)

- **Summary**: Defines the exact content and layout for all types of printouts that can be generated by the tachograph.
- **Developer Takeaway**: If you need to parse or generate human-readable reports from raw tachograph data, this appendix provides the official format.
- **Key Points (in order)**:
    - **Section 1: GENERALITIES**: Specifies general formatting rules, such as right-alignment for numbers, left-alignment for strings, and using spaces for unknown data.
    - **Section 2: DATA BLOCKS SPECIFICATION**: Defines the content and format of reusable "data blocks" that are combined to create a full printout. Each block has a unique identifier and structure.
    - **Section 3: PRINTOUT SPECIFICATIONS**: Lays out the structure of the six types of printouts by specifying the sequence of data blocks for each one:
        - Driver Activities from Card Daily Printout
        - Driver Activities from VU Daily Printout
        - Events and Faults from Card Printout
        - Events and Faults from VU Printout
        - Technical data Printout
        - Over speeding Printout

### [Appendix 5 & 6: Display and Connectors](regulation/13-appendix-5-display.html)

- **Summary**: These appendices specify the requirements for the tachograph's visual display and the physical front connector used for calibration and downloading.
- **Developer Takeaway**: Appendix 6 is important if you are developing hardware that interfaces directly with the tachograph's front port.
- **Key Points (in order)**:
    - **Appendix 5**: Defines the visual format for default displays, warnings, and other information screens, using the pictograms from Appendix 3.
    - **Appendix 6**:
        - Specifies the dimensions and pin allocation for the 6-pin front connector.
        - **Pinout**: Defines pins for power, ground, data downloading (RxD/TxD), and calibration (K-line).
        - **Downloading Interface**: Must comply with RS232 specifications, using an 8-E-1 format (8 data bits, even parity, 1 stop bit) at speeds from 9600 bps to 115200 bps.
        - **Calibration Interface**: Must comply with ISO 14230-1 (KWP2000) on the K-line.

### [Appendix 7: Data Downloading Protocols](regulation/15-appendix-7-data-downloading-protocols.html)

- **Summary**: Specifies the protocols for downloading data from both the Vehicle Unit and tachograph cards.
- **Developer Takeaway**: **Crucial for parsing downloaded files.** This appendix explains how the raw data is wrapped and structured for transfer and storage. You need to understand this to correctly extract and verify the data.
- **Key Points (in order)**:
    - **Section 2: VU DATA DOWNLOADING**:
        - Describes the download procedure, which requires a company, control, or workshop card to be inserted.
        - Defines the message-based protocol, which is based on Keyword Protocol 2000 (ISO 14230-2).
        - Specifies the message structure (Header, Data Field, Checksum) and message types (e.g., `Start Communication`, `Request Upload`, `Transfer Data`).
        - Defines the message flow for a complete download session.
    - **Section 3: TACHOGRAPH CARDS DOWNLOADING PROTOCOL**:
        - Describes the process for downloading a card directly using an external card reader.
        - Defines the sequence of APDU commands required to select and read the files from the card.
    - **Section 3.4: DATA STORAGE FORMAT**:
        - Specifies that downloaded data must be stored in files with a `.ddd` extension.
        - Defines the TLV (Tag-Length-Value) structure of the downloaded file, where each "card file" or "VU file" is stored with a specific tag, its length, and its raw content, followed by a digital signature.

### [Appendix 8: Calibration Protocol](regulation/16-appendix-8-calibration-protocol.html)

- **Summary**: Defines the protocol used by workshops to calibrate the tachograph.
- **Developer Takeaway**: Primarily for developers of calibration tools, but provides insight into which vehicle parameters are configurable and how they are set.
- **Key Points (in order)**:
    - The protocol uses the K-line interface on the front connector.
    - It is based on ISO 14229-1 (UDS).
    - Defines services for:
        - Starting and stopping communication (`StartCommunication`, `StopCommunication`).
        - Managing diagnostic sessions (`StartDiagnosticSession`).
        - Gaining security access (`SecurityAccess`).
        - Reading and writing data by identifier (`ReadDataByIdentifier`, `WriteDataByIdentifier`).
    - Lists the specific `dataRecord` identifiers for calibration parameters like `vehicleIdentificationNumber`, `w_CharacteristicCoefficient`, `k_ConstantOfRecordingEquipment`, and `tyreSize`.

### [Appendix 9 & 10: Type Approval and Security](regulation/17-appendix-9-type-approval-list-of-minimum-required-tests.html)

- **Summary**: These appendices outline the minimum tests required for type approval and the high-level security requirements.
- **Developer Takeaway**: Confirms that the tachograph is a security-certified device. The actual security implementation details are in Appendix 11.
- **Key Points (in order)**:
    - **Appendix 9**: Lists the minimum functional, environmental, EMC, and interoperability tests that a component must pass to receive type approval.
    - **Appendix 10**:
        - Mandates that the VU, tachograph card, motion sensor, and external GNSS facility must be security certified.
        - The certification must be based on the Common Criteria (CC) scheme.
        - The required assurance level is EAL4, augmented with ATE_DPT.2 and AVA_VAN.5.

### [Appendix 11: Common Security Mechanisms](regulation/19-appendix-11-common-security-mechanisms.html)

- **Summary**: Details the cryptographic algorithms and procedures that secure the entire tachograph system.
- **Developer Takeaway**: **Essential for parsing downloaded files.** You cannot verify the authenticity or integrity of the data without correctly implementing the signature verification process described here. You will also need access to the European public key.
- **Key Points (in order)**:
    - **Part A (First-Generation)**:
        - Uses RSA-1024 for public-key cryptography and Triple-DES for symmetric cryptography.
        - Defines the key hierarchy (European -> Member State -> Equipment).
        - Specifies the format of public key certificates.
        - Details the mutual authentication mechanism between VU and cards.
        - Describes the Secure Messaging (SM) protocol for protecting data transfer.
        - Defines the process for creating and verifying digital signatures on downloaded data using RSA and SHA-1.
    - **Part B (Second-Generation)**:
        - Upgrades the cryptographic algorithms to Elliptic Curve Cryptography (ECC) and AES.
        - Defines new cipher suites and standardized domain parameters for ECC.
        - Details the new key and certificate management, including certificate chains.
        - Specifies the updated mutual authentication and session key agreement protocol (Chip Authentication).
        - Describes the new Secure Messaging protocol based on AES.
        - Defines the updated digital signature process using ECDSA.

### [Appendix 12 & 13: GNSS and ITS Interface](regulation/20-appendix-12-positioning-based-on-global-navigation-satellite-system-gnss.html)

- **Summary**: These appendices cover the use of GNSS for positioning and the optional ITS interface for external devices.
- **Developer Takeaway**: Appendix 12 is key for parsing location data. Appendix 13 is relevant for developers creating mobile or other external applications that will communicate with the tachograph.
- **Key Points (in order)**:
    - **Appendix 12**:
        - Mandates compatibility with Galileo and EGNOS.
        - Requires the GNSS receiver to output standard NMEA sentences, specifically `RMC` (Recommended Minimum Specific) and `GSA` (GPS DOP and active satellites).
        - Defines two configurations: VU with an internal GNSS receiver, and VU with a secure external GNSS facility.
        - For the external facility, it specifies a secure communication protocol between the VU and the facility based on ISO/IEC 7816-4.
    - **Appendix 13**:
        - Specifies that the ITS interface must use Bluetooth® (version 5.0 or newer, Low Energy compatible).
        - The VU acts as the server, and the external ITS unit acts as the client.
        - Data is made available via the services defined in Appendices 7 and 8.
        - Access to personal data requires explicit driver consent, which is recorded by the VU.
        - Provides a table classifying various data points as personal or not personal.

### [Appendix 14: Remote Communication Function](regulation/26-appendix-14-remote-communication-function.html)

- **Summary**: Details the DSRC (Dedicated Short-Range Communication) interface used by control authorities for remote early detection.
- **Developer Takeaway**: Primarily for developers of control authority equipment, but provides context on how vehicles are pre-selected for inspection.
- **Key Points (in order)**:
    - The function is for pre-selecting vehicles and does not replace a formal inspection.
    - Communication uses 5.8 GHz DSRC.
    - The communication is initiated only by a control authority's reader (REDCR).
    - Data is secured to ensure integrity and is limited to what is necessary for targeting checks.
    - Defines the `RtmData` structure, which contains the 19 data items to be transmitted (e.g., `Time of last entry`, `Overspeed`, `Driving without a valid card`).
    - Specifies the DSRC protocol layers and transaction sequence (`GET` command).

### [Appendix 15 & Addendum: Migration and Time Calculation](regulation/28-appendix-15-migration-managing-the-co-existence-of-equipment-generations-and-versions.html)

- **Summary**: These sections cover the transition between tachograph generations and the rules for calculating driving times.
- **Developer Takeaway**: Appendix 15 is vital for any tool that needs to support both tachograph generations. The Addendum is essential for any application that performs compliance analysis of driver activity.
- **Key Points (in order)**:
    - **Appendix 15 (Migration)**:
        - Gen1 cards can be used in Gen2 VUs (but not vice-versa for all functions).
        - Gen2 cards can be used in Gen1 VUs, where they will function as Gen1 cards.
        - Gen2 VUs cannot be paired with Gen1 motion sensors.
        - Data download equipment must be able to handle both generations. It specifies how a Gen2 card/VU should respond to a Gen1 download request.
    - **Addendum (Time Calculation)**:
        - Defines concepts like `RTM-shift`, `accumulated driving time`, `daily rest period`, etc., for the purpose of calculation.
        - Provides detailed rules for how the VU must compute daily, weekly, and fortnightly driving times.
        - Explains how to handle special conditions like `OUT OF SCOPE` and `FERRY/TRAIN CROSSING` when calculating rest periods.

### [Appendix 16 & 17: Adaptor and Transitional Provisions](regulation/29-appendix-16-adaptor-for-m1-and-n1-category-vehicles.html)

- **Summary**: These cover special cases: the use of an "adaptor" for M1/N1 vehicles and the transitional rules for the introduction of OSNMA.
- **Developer Takeaway**: These are edge cases. Appendix 16 is only relevant for a small subset of vehicles. The transitional provisions are important for understanding the behavior of early Gen2 VUs.
- **Key Points (in order)**:
    - **Appendix 16 (Adaptor)**:
        - An adaptor is for M1/N1 vehicles where a standard motion sensor cannot be mechanically installed.
        - It contains a type-approved motion sensor within a sealed, yellow housing.
        - It converts incoming speed pulses from the vehicle into a signal that the embedded motion sensor can read.
        - The interface between the adaptor and the VU is the same as a standard motion sensor.
    - **Appendix 17 & related files (31-35) (Transitional Provisions)**:
        - Defines a `Transitional vehicle unit` as a Gen2 VU manufactured before the Galileo OSNMA service is fully operational.
        - These units use the OSNMA public test phase signals.
        - Crucially, until a software update is applied, these units **assume that standard GNSS positions are authenticated**.
        - Specifies how requirements related to position authentication (e.g., `Time conflict` event, data recording flags) are to be interpreted for these transitional units.

### [ANNEX II: Approval Mark and Certificate](regulation/36-annex-2-approval-mark-and-certificate.html)

- **Summary**: Describes the approval mark and the format of the type-approval certificate.
- **Developer Takeaway**: Provides context on the regulatory markings found on equipment. The approval number is part of the equipment's identification data.
- **Key Points (in order)**:
    - The approval mark consists of a rectangle with an 'e' followed by the country code (e.g., 'e1' for Germany).
    - An approval number is placed near the rectangle.
    - The mark must be shown on the descriptive plaque of the equipment and on each tachograph card.
    - Provides the template for the official EC type-approval certificate for analogue, digital, and smart tachographs.