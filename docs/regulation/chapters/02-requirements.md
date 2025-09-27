# *ANNEX I C*

## **Requirements for construction, testing, installation, and inspection**

INTRODUCTION

| INTRODUCTION |             |                                                                                                                                                                      |
|--------------|-------------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------|
|              | 1.          | DEFINITIONS                                                                                                                                                          |
|              | 2.          | GENERAL CHARACTERISTICS AND FUNCTIONS OF THE<br>RECORDING EQUIPMENT                                                                                                  |
|              | 2.1         | General characteristics                                                                                                                                              |
|              | 2.2         | Functions                                                                                                                                                            |
|              | 2.3         | Modes of operation                                                                                                                                                   |
|              | 2.4         | Security                                                                                                                                                             |
|              | 3.          | CONSTRUCTION AND FUNCTIONAL REQUIREMENTS FOR<br>RECORDING EQUIPMENT                                                                                                  |
|              | 3.1         | Monitoring cards insertion and withdrawal                                                                                                                            |
|              | 3.2         | Speed, position and distance measurement                                                                                                                             |
|              | 3.2.1       | Measurement of distance travelled                                                                                                                                    |
|              | 3.2.2       | Measurement of speed                                                                                                                                                 |
|              | 3.2.3       | Measurement of position                                                                                                                                              |
|              | 3.3         | Time measurement                                                                                                                                                     |
|              | 3.4         | Monitoring driver activities                                                                                                                                         |
|              | 3.5         | Monitoring driving status                                                                                                                                            |
|              | 3.6         | Driver's entries                                                                                                                                                     |
|              | 3.6.1       | Entry of places where daily work periods begin and/or end                                                                                                            |
|              | 3.6.2       | Manual entry of driver activities and driver consent for ITS<br>interface                                                                                            |
|              | 3.6.3       | Entry of specific conditions                                                                                                                                         |
| ▼M3<br>▼B    | 3.6.4       | Entry of load/unload operation                                                                                                                                       |
|              | 3.7         | Company locks management                                                                                                                                             |
|              | 3.8         | Monitoring control activities                                                                                                                                        |
|              | 3.9         | Detection of events and/or faults                                                                                                                                    |
|              | 3.9.1       | 'Insertion of a non-valid card' event                                                                                                                                |
|              | 3.9.2       | 'Card conflict' event                                                                                                                                                |
|              | 3.9.3       | 'Time overlap' event                                                                                                                                                 |
|              | 3.9.4       | 'Driving without an appropriate card' event                                                                                                                          |
|              | 3.9.5       | 'Card insertion while driving' event                                                                                                                                 |
|              | 3.9.6       | 'Last card session not correctly closed' event                                                                                                                       |
|              | 3.9.7       | 'Over speeding' event                                                                                                                                                |
|              | 3.9.8       | 'Power supply interruption' event                                                                                                                                    |
|              | 3.9.9       | 'Communication error with the remote communication facility'<br>event                                                                                                |
|              | 3.9.10      | 'Absence of position information from GNSS receiver' event                                                                                                           |
| ▼B           |             |                                                                                                                                                                      |
|              | 3.9.11      | 'Communication error with the external GNSS facility' event                                                                                                          |
|              | 3.9.12      | 'Motion data error' event                                                                                                                                            |
|              | 3.9.13      | 'Vehicle motion conflict' event                                                                                                                                      |
|              | 3.9.14      | 'Security breach attempt' event                                                                                                                                      |
|              | 3.9.15      | 'Time conflict' event                                                                                                                                                |
|              | 3.9.16      | 'Card' fault                                                                                                                                                         |
|              | 3.9.17      | 'Recording equipment' fault                                                                                                                                          |
| ▼M3          |             |                                                                                                                                                                      |
|              | 3.9.18      | 'GNSS anomaly' event                                                                                                                                                 |
| ▼B           | 3.10        | Built-in and self-tests                                                                                                                                              |
|              | 3.11        | Reading from data memory                                                                                                                                             |
|              | 3.12        | Recording and storing in the data memory                                                                                                                             |
|              | 3.12.1      | Equipment identification data                                                                                                                                        |
|              | 3.12.1.1    | Vehicle unit identification data                                                                                                                                     |
|              | 3.12.1.2    | Motion sensor identification data                                                                                                                                    |
|              | 3.12.1.3    | Global Navigation Satellite Systems identification data                                                                                                              |
|              | 3.12.2      | Keys and certificates                                                                                                                                                |
|              | 3.12.3      | Driver or workshop card insertion and withdrawal data                                                                                                                |
|              | 3.12.4      | Driver activity data                                                                                                                                                 |
| ▼M1          | 3.12.5      | Places and positions where daily work periods begin, end and/or<br>where 3 hours accumulated driving time is reached                                                 |
| ▼B           | 3.12.6      | Odometer data                                                                                                                                                        |
|              | 3.12.7      | Detailed speed data                                                                                                                                                  |
|              | 3.12.8      | Events data                                                                                                                                                          |
|              | 3.12.9      | Faults data                                                                                                                                                          |
|              | 3.12.10     | Calibration data                                                                                                                                                     |
|              | 3.12.11     | Time adjustment data                                                                                                                                                 |
|              | 3.12.12     | Control activity data                                                                                                                                                |
|              | 3.12.13     | Company locks data                                                                                                                                                   |
|              | 3.12.14     | Download activity data                                                                                                                                               |
|              | 3.12.15     | Specific conditions data                                                                                                                                             |
|              | 3.12.16     | Tachograph card data                                                                                                                                                 |
| ▼M3          |             |                                                                                                                                                                      |
|              | 3.12.17     | Border crossings                                                                                                                                                     |
|              | 3.12.18     | Load/unload operations                                                                                                                                               |
|              | 3.12.19     | Digital map                                                                                                                                                          |
| ▼B           | 3.13        | Reading from tachograph cards                                                                                                                                        |
|              | 3.14        | Recording and storing on tachograph cards                                                                                                                            |
|              | 3.14.1      | Recording and storing in first generation tachograph cards                                                                                                           |
|              |             |                                                                                                                                                                      |
| ▼B           | 3.14.2      | Recording and storing in second generation tachograph cards                                                                                                          |
|              | 3.15        | Displaying                                                                                                                                                           |
|              | 3.15.1      | Default display                                                                                                                                                      |
|              | 3.15.2      | Warning display                                                                                                                                                      |
|              | 3.15.3      | Menu access                                                                                                                                                          |
|              | 3.15.4      | Other displays                                                                                                                                                       |
|              | 3.16        | Printing                                                                                                                                                             |
|              | 3.17        | Warnings                                                                                                                                                             |
|              | 3.18        | Data downloading to external media                                                                                                                                   |
|              | 3.19        | Remote communication for targeted roadside checks                                                                                                                    |
| ▼M3          | 3.20        | Data exchanges with additional external devices                                                                                                                      |
| ▼B           | 3.21        | Calibration                                                                                                                                                          |
|              | 3.22        | Roadside calibration checking                                                                                                                                        |
|              | 3.23        | Time adjustment                                                                                                                                                      |
|              | 3.24        | Performance characteristics                                                                                                                                          |
|              | 3.25        | Materials                                                                                                                                                            |
|              | 3.26        | Markings                                                                                                                                                             |
| ▼M3          | 3.27        | Monitoring border crossings                                                                                                                                          |
|              | 3.28        | Software update                                                                                                                                                      |
| ▼B           | 4.          | CONSTRUCTION AND FUNCTIONAL REQUIREMENTS FOR<br>TACHOGRAPH CARDS                                                                                                     |
|              | 4.1         | Visible data                                                                                                                                                         |
|              | 4.2         | Security                                                                                                                                                             |
|              | 4.3         | Standards                                                                                                                                                            |
|              | 4.4         | Environmental and electrical specifications                                                                                                                          |
|              | 4.5         | Data storage                                                                                                                                                         |
|              | 4.5.1       | Elementary files for identification and card management                                                                                                              |
|              | 4.5.2       | IC card identification                                                                                                                                               |
|              | 4.5.2.1     | Chip identification                                                                                                                                                  |
|              | 4.5.2.2     | DIR (only present in second generation tachograph cards)                                                                                                             |
|              | 4.5.2.3     | ATR information (conditional, only present in second generation<br>tachograph cards)                                                                                 |
|              | 4.5.2.4     | Extended length information (conditional, only present in second<br>generation tachograph cards)                                                                     |
|              | 4.5.3       | Driver card                                                                                                                                                          |
|              | 4.5.3.1     | Tachograph application (accessible to first and second generation<br>vehicle units)                                                                                  |
|              | 4.5.3.1.1   | Application identification                                                                                                                                           |
|              | 4.5.3.1.2   | Key and certificates                                                                                                                                                 |
|              | 4.5.3.1.3   | Card identification                                                                                                                                                  |
| ▼B           | 4.5.3.1.5   | Card download                                                                                                                                                        |
|              | 4.5.3.1.6   | Driving licence information                                                                                                                                          |
|              | 4.5.3.1.7   | Events data                                                                                                                                                          |
|              | 4.5.3.1.8   | Faults data                                                                                                                                                          |
|              | 4.5.3.1.9   | Driver activity data                                                                                                                                                 |
|              | 4.5.3.1.10  | Vehicles used data                                                                                                                                                   |
|              | 4.5.3.1.11  | Places where daily work periods start and/or end                                                                                                                     |
|              | 4.5.3.1.12  | Card session data                                                                                                                                                    |
|              | 4.5.3.1.13  | Control activity data                                                                                                                                                |
|              | 4.5.3.1.14  | Specific conditions data                                                                                                                                             |
|              | 4.5.3.2     | Tachograph generation 2 application (not accessible to first<br>generation vehicle unit)                                                                             |
|              | 4.5.3.2.1   | Application identification                                                                                                                                           |
| ▼M3          | 4.5.3.2.1.1 | Additional application identification (not accessed by version 1<br>of second generation vehicle units)                                                              |
| ▼B           | 4.5.3.2.2   | Keys and certificates                                                                                                                                                |
|              | 4.5.3.2.3   | Card identification                                                                                                                                                  |
|              | 4.5.3.2.4   | Card holder identification                                                                                                                                           |
|              | 4.5.3.2.5   | Card download                                                                                                                                                        |
|              | 4.5.3.2.6   | Driving licence information                                                                                                                                          |
|              | 4.5.3.2.7   | Events data                                                                                                                                                          |
|              | 4.5.3.2.8   | Faults data                                                                                                                                                          |
|              | 4.5.3.2.9   | Driver activity data                                                                                                                                                 |
|              | 4.5.3.2.10  | Vehicles used data                                                                                                                                                   |
|              | 4.5.3.2.11  | Places and positions where daily work periods start and/or end                                                                                                       |
|              | 4.5.3.2.12  | Card session data                                                                                                                                                    |
|              | 4.5.3.2.13  | Control activity data                                                                                                                                                |
|              | 4.5.3.2.14  | Specific conditions data                                                                                                                                             |
|              | 4.5.3.2.15  | Vehicle units used data                                                                                                                                              |
| ▼M1          | 4.5.3.2.16  | Three hours accumulated driving places data                                                                                                                          |
| ▼M3          | 4.5.3.2.17  | Authentication status for positions related to places where daily<br>work periods start and/or end (not accessed by version 1<br>of second generation vehicle units) |
|              | 4.5.3.2.18  | Authentication status for positions where three hours accumulated<br>driving time are reached (not accessed by version 1 of second<br>generation vehicle units)      |
|              | 4.5.3.2.19  | Border crossings (not accessed by version 1 of second generation<br>vehicle units)                                                                                   |
| ▼M3          |             |                                                                                                                                                                      |
|              | 4.5.3.2.20  | Load/unload operations (not accessed by version 1 of second<br>generation vehicle units)                                                                             |
|              | 4.5.3.2.21  | Load type entries (not accessed by version 1 of second generation<br>vehicle units)                                                                                  |
|              | 4.5.3.2.22  | VU configurations (not accessed by version 1 of second generation<br>vehicle units)                                                                                  |
| ▼B           | 4.5.4       | Workshop card                                                                                                                                                        |
|              | 4.5.4.1     | Tachograph application (accessible to first and second generation<br>vehicle units)                                                                                  |
|              | 4.5.4.1.1   | Application identification                                                                                                                                           |
|              | 4.5.4.1.2   | Keys and certificates                                                                                                                                                |
|              | 4.5.4.1.3   | Card identification                                                                                                                                                  |
|              | 4.5.4.1.4   | Card holder identification                                                                                                                                           |
|              | 4.5.4.1.5   | Card download                                                                                                                                                        |
|              | 4.5.4.1.6   | Calibration and time adjustment data                                                                                                                                 |
|              | 4.5.4.1.7   | Events and faults data                                                                                                                                               |
|              | 4.5.4.1.8   | Driver activity data                                                                                                                                                 |
|              | 4.5.4.1.9   | Vehicles used data                                                                                                                                                   |
|              | 4.5.4.1.10  | Daily work periods start and/or end data                                                                                                                             |
|              | 4.5.4.1.11  | Card session data                                                                                                                                                    |
|              | 4.5.4.1.12  | Control activity data                                                                                                                                                |
|              | 4.5.4.1.13  | Specific conditions data                                                                                                                                             |
|              | 4.5.4.2     | Tachograph generation 2 application (not accessible to first<br>generation vehicle unit)                                                                             |
|              | 4.5.4.2.1   | Application identification                                                                                                                                           |
| ▼M3          | 4.5.4.2.1.1 | Additional application identification (not accessed by version 1<br>of second generation vehicle units)                                                              |
| ▼B           | 4.5.4.2.2   | Keys and certificates                                                                                                                                                |
|              | 4.5.4.2.3   | Card identification                                                                                                                                                  |
|              | 4.5.4.2.4   | Card holder identification                                                                                                                                           |
|              | 4.5.4.2.5   | Card download                                                                                                                                                        |
|              | 4.5.4.2.6   | Calibration and time adjustment data                                                                                                                                 |
|              | 4.5.4.2.7   | Events and faults data                                                                                                                                               |
|              | 4.5.4.2.8   | Driver activity data                                                                                                                                                 |
|              | 4.5.4.2.9   | Vehicles used data                                                                                                                                                   |
|              | 4.5.4.2.10  | Daily work periods start and/or end data                                                                                                                             |
|              |             |                                                                                                                                                                      |
| ▼B           |             |                                                                                                                                                                      |
|              | 4.5.4.2.12  | Control activity data                                                                                                                                                |
|              | 4.5.4.2.13  | Vehicle units used data                                                                                                                                              |
| ▼M1          | 4.5.4.2.14  | Three hours accumulated driving places data                                                                                                                          |
| ▼B           | 4.5.4.2.15  | Specific conditions data                                                                                                                                             |
| ▼M3          | 4.5.4.2.16  | Authentication status for positions related to places where daily<br>work periods start and/or end (not accessed by version 1<br>of second generation vehicle units) |
|              | 4.5.4.2.17  | Authentication status for positions where three hours accumulated<br>driving are reached (not accessed by version 1 of second<br>generation vehicle units)           |
|              | 4.5.4.2.18  | Border crossings (not accessed by version 1 of second generation<br>vehicle units)                                                                                   |
|              | 4.5.4.2.19  | Load/unload operations (not accessed by version 1 of second<br>generation vehicle units)                                                                             |
|              | 4.5.4.2.20  | Load type entries (not accessed by version 1 of second generation<br>vehicle units)                                                                                  |
|              | 4.5.4.2.21  | Calibration Additional Data (not accessed by version 1 of second<br>generation vehicle units)                                                                        |
|              | 4.5.4.2.22  | VU configurations (not accessed by version 1 of second generation<br>vehicle units)                                                                                  |
| ▼B           | 4.5.5       | Control card                                                                                                                                                         |
|              | 4.5.5.1     | Tachograph application (accessible to first and second generation<br>vehicle units)                                                                                  |
|              | 4.5.5.1.1   | Application identification                                                                                                                                           |
|              | 4.5.5.1.2   | Keys and certificates                                                                                                                                                |
|              | 4.5.5.1.3   | Card identification                                                                                                                                                  |
|              | 4.5.5.1.4   | Card holder identification                                                                                                                                           |
|              | 4.5.5.1.5   | Control activity data                                                                                                                                                |
|              | 4.5.5.2     | Tachograph G2 application (not accessible to first generation<br>vehicle unit)                                                                                       |
|              | 4.5.5.2.1   | Application identification                                                                                                                                           |
| ▼M3          | 4.5.5.2.1.1 | Additional application identification (not accessed by version 1<br>of second generation vehicle units)                                                              |
| ▼B           |             |                                                                                                                                                                      |
|              | 4.5.5.2.2   | Keys and certificates                                                                                                                                                |
|              | 4.5.5.2.3   | Card identification                                                                                                                                                  |
|              | 4.5.5.2.4   | Card holder identification                                                                                                                                           |
|              | 4.5.5.2.5   | Control activity data                                                                                                                                                |
| ▼M3          | 4.5.5.2.6   | VU configurations (not accessed by version 1 of second generation<br>vehicle units)                                                                                  |
| ▼B           |             |                                                                                                                                                                      |
| ▼B           | 4.5.6.1     | Tachograph application (accessible to first and second generation vehicle units)                                                                                     |
|              | 4.5.6.1.1   | Application identification                                                                                                                                           |
|              | 4.5.6.1.2   | Keys and certificates                                                                                                                                                |
|              | 4.5.6.1.3   | Card identification                                                                                                                                                  |
|              | 4.5.6.1.4   | Card holder identification                                                                                                                                           |
|              | 4.5.6.1.5   | Company activity data                                                                                                                                                |
|              | 4.5.6.2     | Tachograph G2 application (not accessible to first generation vehicle unit)                                                                                          |
|              | 4.5.6.2.1   | Application identification                                                                                                                                           |
| ▼M3          | 4.5.6.2.1.1 | Additional application identification (not accessed by version 1 of second generation vehicle units)                                                                 |
| ▼B           | 4.5.6.2.2   | Keys and certificates                                                                                                                                                |
|              | 4.5.6.2.3   | Card identification                                                                                                                                                  |
|              | 4.5.6.2.4   | Card holder identification                                                                                                                                           |
|              | 4.5.6.2.5   | Company activity data                                                                                                                                                |
| ▼M3          | 4.5.6.2.6   | VU configurations (not accessed by version 1 of second generation vehicle units)                                                                                     |
| ▼B           | 5.          | INSTALLATION OF RECORDING EQUIPMENT                                                                                                                                  |
|              | 5.1         | Installation                                                                                                                                                         |
|              | 5.2         | Installation plaque                                                                                                                                                  |
|              | 5.3         | Sealing                                                                                                                                                              |
|              | 6.          | CHECKS, INSPECTIONS AND REPAIRS                                                                                                                                      |
|              | 6.1         | Approval of fitters, workshops and vehicle manufacturers                                                                                                             |
| ▼M1          | 6.2         | Check of new or repaired components                                                                                                                                  |
| ▼B           | 6.3         | Installation inspection                                                                                                                                              |
|              | 6.4         | Periodic inspections                                                                                                                                                 |
|              | 6.5         | Measurement of errors                                                                                                                                                |
|              | 6.6         | Repairs                                                                                                                                                              |
|              | 7.          | CARD ISSUING                                                                                                                                                         |
|              | 8.          | TYPE-APPROVAL OF RECORDING EQUIPMENT AND TACHOGRAPH CARDS                                                                                                            |
|              | 8.1         | General points                                                                                                                                                       |
|              | 8.2         | Security certificate                                                                                                                                                 |
|              | 8.3         | Functional certificate                                                                                                                                               |
|              | 8.4         | Interoperability certificate                                                                                                                                         |
|              | 8.5         | Type-approval certificate                                                                                                                                            |
|              | 8.6         | Exceptional procedure: first interoperability certificates for 2nd generation recording equipment and tachograph cards                                               |

### **B**

4.5.4.2.11 Card session data

4.5.6 Company card

## INTRODUCTION

This Annex contains second generation recording equipment and tachograph cards requirements.

Since June 15th 2019, second generation recording equipment is being installed in vehicles registered in the Union for the first time, and second generation tachograph cards are being issued.

In order to smoothly implement the second generation tachograph system, second generation tachograph cards have been designed to be also used in first generation vehicle units built in accordance with Annex IB to Regulation (EEC) No 3821/85.

Reciprocally, first generation tachograph cards may be used in second generation vehicle units. Nevertheless, second generation vehicle units can only be calibrated using second generation workshop cards.

The requirements regarding the interoperability between the first and the second generation tachograph systems are specified in this Annex. In this respect, Appendix 15 contains additional details on the management of the co-existence of both generations.

In addition, due to the implementation of new functions such as the use of Galileo Open Signal Navigation Messages Authentication, detection of border crossings, entry of load and unload operations, and also to the need to increase the driver card capacity to 56 days of driver activities, this Regulation introduces the technical requirements for the second version of the second generation recording equipment and tachograph cards.

### **B**

List of Appendixes

- App 1: DATA DICTIONARY
- App 2: TACHOGRAPH CARDS SPECIFICATION
- App 3: PICTOGRAMS
- App 4: PRINTOUTS
- App 5: DISPLAY
- App 6: FRONT CONNECTOR FOR CALIBRATION AND DOWNLOAD
- App 7: DATA DOWNLOADING PROTOCOLS
- App 8: CALIBRATION PROTOCOL
- App 9: TYPE-APPROVAL AND LIST OF MINIMUM REQUIRED TESTS
- App 10: SECURITY REQUIREMENTS
- App 11: COMMON SECURITY MECHANISMS
- App 12: POSITIONING BASED ON GLOBAL NAVIGATION SATELLITE SYSTEM (GNSS)
- App 13: ITS INTERFACE
- App 14: REMOTE COMMUNICATION FUNCTION
- App 15: MIGRATION: MANAGING THE COEXISTENCE OF EQUIPMENT GENERATIONS
- App 16: ADAPTOR FOR M1 AND N1 CATEGORY VEHICLES

## 1. DEFINITIONS

In this Annex:

(a) 'activation' means:

the phase in which the tachograph becomes fully operational and implements all functions, including security functions, through the use of a workshop card;

(b) 'authentication' means:

a function intended to establish and verify a claimed identity;

(c) 'authenticity' means:

the property that information is coming from a party whose identity can be verified;

(d) 'built-in test (BIT)' means:

tests run at request, triggered by the operator or by external equipment;

(e) 'calendar day' means:

a day ranging from 00:00 hours to 24:00 hours. All calendar days relate to UTC time (Universal Time Coordinated);

**▼M3**

(f) 'calibration of a smart tachograph' means:

updating or confirming vehicle parameters to be held in the data memory. Vehicle parameters include vehicle identification (VIN, VRN and registering Member State) and vehicle characteristics (w, k, l, tyre size, speed-limiting device setting (if applicable), current UTC time, current odometer value, by-default load type); during the calibration of a recording equipment, the types and identifiers of all type-approval relevant seals in place shall also be stored in the data memory;

any update or confirmation of UTC time only, shall be considered as a time adjustment and not as a calibration, provided it does not contradict requirement 409 set out in point 6.4.

calibrating a recording equipment requires the use of a workshop card;

(g) 'card number' means:

a 16-alpha-numerical characters number that uniquely identifies a tachograph card within a Member State. The card number includes an identification, which consists in a driver identification, or in a card owner identification together with a card consecutive index, a card replacement index and a card renewal index;

a card is therefore uniquely identified by the code of the issuing Member State and the card number;

(h) 'card consecutive index' means:

the 14th alphanumerical character of a card number that is used to differentiate the different cards issued to a company, a workshop or a control authority entitled to be issued several tachograph cards. The company, the workshop or the control authority is uniquely identified by the 13 first characters of the card number;

### **B**

(i) 'card renewal index' means:

the 16th alpha-numerical character of a card number which is incremented each time a tachograph card corresponding to a given identification, i.e. driver identification or owner identification together with consecutive index, is renewed;

(j) 'card replacement index' means:

the 15th alpha-numerical character of a card number which is incremented each time a tachograph card corresponding to a given identification, i.e. driver identification or owner identification together with consecutive index, is replaced;

(k) 'characteristic coefficient of the vehicle' means:

the numerical characteristic giving the value of the output signal emitted by the part of the vehicle linking it with the recording equipment (gearbox output shaft or axle) while the vehicle travels a distance of one kilometre under standard test conditions as defined under requirement 414. The characteristic coefficient is expressed in impulses per kilometre (w = … imp/km);

(l) 'company card' means:

a tachograph card issued by the authorities of a Member State to a transport undertaking needing to operate vehicles fitted with a tachograph, which identifies the transport undertaking and allows for the displaying, downloading and printing of the data, stored in the tachograph, which have been locked by that transport undertaking;

(m) 'constant of the recording equipment' means:

the numerical characteristic giving the value of the input signal required to show and record a distance travelled of one kilometre; this constant shall be expressed in impulses per kilometre (k = … imp/km);

(n) 'continuous driving time' is computed within the recording equipment as (1):

the continuous driving time is computed as the current accumulated driving times of a particular driver, since the end of his last AVAILABILITY or BREAK/REST or UNKNOWN (2) period of 45 minutes or more (this period may have been split according to Regulation (EC) No 561/2006 of the European Parliament and of the Council (3)). The computations involved take into account, as needed, past activities stored on the driver card. When the driver has not inserted his card, the computations involved are based on the data memory recordings related to the current period where no card was inserted and related to the relevant slot;

### **M3**

<sup>(1)</sup> This way of computing the continuous driving time and the cumulative break time serves in the recording equipment for computing the continuous driving time warning. It does not prejudge the legal interpretation to be made of these times. Alternative ways of computing the continuous driving time and the cumulative break time may be used to replace these definitions if they have been made obsolete by updates in other relevant legislation.

<sup>(2)</sup> UNKNOWN periods correspond to periods where the driver card was not inserted in the recording equipment and for which no manual entry of driver activities was made.

<sup>(3)</sup> Regulation (EC) No 561/2006 of the European Parliament and of the Council of 15 March 2006 on the harmonisation of certain social legislation relating to road transport and amending Council Regulations (EEC) No 3821/85 and (EC) No 2135/98 and repealing Council Regulation (EEC) No 3820/85 (OJ L 102, 11.4.2006, p. 1).

(o) 'control card' means:

a tachograph card issued by the authorities of a Member State to a national competent control authority which identifies the control body and, optionally, the control officer, and which allows access to the data stored in the data memory or in the driver cards and, optionally, in the workshop cards for reading, printing and/or downloading;

It shall also give access to the roadside calibration checking function and to data on the remote early detection communication reader;

(p) 'cumulative break time' is computed within the recording equipment as (1):

> the cumulative break from driving time is computed as the current accumulated AVAILABILITY or BREAK/REST or UNKNOWN (2) times of 15 minutes or more of a particular driver, since the end of his last AVAILABILITY or BREAK/REST or UNKNOWN (2) period of 45 minutes or more (this period may have been split according to Regulation (EC) No 561/2006).

> The computations involved take into account, as needed, past activities stored on the driver card. Unknown periods of negative duration (start of unknown period > end of unknown period) due to time overlaps between two different sets of recording equipment, are not taken into account for the computation.

> When the driver has not inserted his card, the computations involved are based on the data memory recordings related to the current period where no card was inserted and related to the relevant slot;

(q) 'data memory' means:

an electronic data storage device built into the recording equipment;

(r) 'digital signature' means:

data appended to, or a cryptographic transformation of, a block of data that allows the recipient of the block of data to prove the authenticity and integrity of the block of data;

(s) 'downloading' means:

the copying, together with the digital signature, of a part, or of a complete set, of data files recorded in the data memory of the vehicle unit or in the memory of a tachograph card, provided that this process does not alter or delete any stored data;

<sup>(1)</sup> This way of computing the continuous driving time and the cumulative break time serves in the recording equipment for computing the continuous driving time warning. It does not prejudge the legal interpretation to be made of these times. Alternative ways of computing the continuous driving time and the cumulative break time may be used to replace these definitions if they have been made obsolete by updates in other relevant legislation.

<sup>(2)</sup> UNKNOWN periods correspond to periods where the driver card was not inserted in the recording equipment and for which no manual entry of driver activities was made.

manufacturers of smart tachograph vehicle units and manufacturers of equipment designed and intended to download data files shall take all reasonable steps to ensure that the downloading of such data can be performed with the minimum delay by transport undertakings or drivers;

The downloading of the detailed speed file may not be necessary to establish compliance with Regulation (EC) No 561/2006, but may be used for other purposes such as accident investigation;

(t) 'driver card' means:

a tachograph card, issued by the authorities of a Member State to a particular driver, which identifies the driver and allows for the storage of driver activity data;

(u) 'effective circumference of the wheels' means:

the average of the distances travelled by each of the wheels moving the vehicle (driving wheels) in the course of one complete rotation. The measurement of these distances shall be made under standard test conditions as defined under requirement 414 and is expressed in the form 'l = … mm'. Vehicle manufacturers may replace the measurement of these distances by a theoretical calculation which takes into account the distribution of the weight on the axles, vehicle unladen in normal running order (1). The methods for such theoretical calculation are subject to approval by a competent Member State authority and can take place only before tachograph activation;

(v) 'event' means:

an abnormal operation detected by the smart tachograph which may result from a fraud attempt;

(w) 'external GNSS facility' means

a facility which contains the GNSS receiver when the vehicle unit is not a single unit as well as other components needed to protect the communication of position data to the rest of the vehicle unit;

(x) 'fault' means:

abnormal operation detected by the smart tachograph which may come from an equipment malfunction or failure;

(y) 'GNSS receiver' means:

an electronic device that receives and digitally processes the signals from one or more Global Navigation Satellite System(s) (GNSS in English) in order to provide position, speed and time information;

<sup>(1)</sup> Commission Regulation (EU) No 1230/2012 of 12 December 2012 implementing Regulation (EC) No 661/2009 of the European Parliament and of the Council with regard to type-approval requirements for masses and dimensions of motor vehicles and their trailers and amending Directive 2007/46/EC of the European Parliament and of the Council (OJ L 353, 21.12.2012, p. 31) as last amended.

(z) 'installation' means:

the mounting of a tachograph in a vehicle;

(aa) 'interoperability' means:

the capacity of systems and the underlying business processes to exchange data and to share information;

(bb) 'interface' means:

a facility between systems which provides the media through which they can connect and interact;

(cc) 'position' means:

geographical coordinates of the vehicle at a given time;

(dd) 'motion sensor' means:

a part of the tachograph, providing a signal representative of vehicle speed and/or distance travelled;

## (ee) 'non-valid card' means:

a card detected as faulty, or which authentication failed, or whose start of validity date is not yet reached, or which expiry date has passed;

- a card is also considered as non-valid by the vehicle unit:
- if a card with the same card issuing Member State, the same identification, i.e. driver identification or owner identification together with consecutive index, and a higher renewal index has already been inserted in the vehicle unit, or
- if a card with the same card issuing Member State, the same identification, i.e. driver identification or owner identification together with consecutive index and renewal index but with a higher replacement index has already been inserted in the vehicle unit;

### **B**

**▼M3**

(ff) 'open standard' means:

a standard set out in a standard specification document available freely or at a nominal charge which it is permissible to copy, distribute or use for no fee or for a nominal fee;

(gg) 'out of scope' means:

when the use of the recording equipment is not required, according to the provisions of Regulation (EC) No 561/2006;

(hh) 'over speeding' means:

exceeding the authorised speed of the vehicle, defined as any period of more than 60 seconds during which the vehicle's measured speed exceeds the limit for setting the speed limitation device laid down in Council Directive 92/6/EEC (1), as last amended;

<sup>(1)</sup> Council Directive 92/6/EEC of 10 February 1992 on the installation and use of speed limitation devices for certain categories of motor vehicles in the Community (OJ L 57, 2.3.1992, p. 27).

|     | (ii) | 'periodic inspection' means:                                                                                                                                                                                                                                         |
|-----|------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
|     |      | a set of operations performed to check that the tachograph<br>works properly, that its settings correspond to the vehicle<br>parameters, and that no manipulation devices are attached to<br>the tachograph;                                                         |
|     | (jj) | 'printer' means:                                                                                                                                                                                                                                                     |
|     |      | component of the recording equipment which provides<br>printouts of stored data;                                                                                                                                                                                     |
|     | (kk) | 'remote early detection communication' means:                                                                                                                                                                                                                        |
|     |      | communication between the remote early detection communi-<br>cation facility and the remote early detection communication<br>reader during targeted roadside checks with the aim of<br>remotely detecting possible manipulation or misuse of<br>recording equipment; |
| ▼M3 | (ll) | 'remote communication facility', 'remote communication<br>module' or 'remote early detection facility' means:                                                                                                                                                        |
|     |      | the equipment of the vehicle unit which is used to perform<br>targeted roadside checks;                                                                                                                                                                              |
| ▼B  | (mm) | 'remote early detection communication reader' means:                                                                                                                                                                                                                 |
|     |      | the system used by control officers for targeted roadside<br>checks.                                                                                                                                                                                                 |
| ▼M3 | (nn) | 'card renewal' means:                                                                                                                                                                                                                                                |
|     |      | issue of a new tachograph card when an existing card reaches<br>its expiry date, or is malfunctioning and has been returned to<br>the issuing authority;                                                                                                             |
| ▼B  | (oo) | 'repair' means:                                                                                                                                                                                                                                                      |
|     |      | any repair of a motion sensor or of a vehicle unit or of a cable<br>that requires the disconnection of its power supply, or its<br>disconnection from other tachograph components, or the<br>opening of the motion sensor or vehicle unit;                           |
| ▼M3 | (pp) | 'card replacement' means:                                                                                                                                                                                                                                            |
|     |      | issue of a new tachograph card in replacement of an existing<br>card, which has been declared as lost, stolen or malfunc-<br>tioning and has not been returned to the issuing authority;                                                                             |
| ▼B  | (qq) | 'security certification' means:                                                                                                                                                                                                                                      |
|     |      | process to certify, by a common criteria certification body,<br>that the recording equipment (or component) or the<br>tachograph card under investigation fulfils the security<br>requirements defined in the relative protection profiles;                          |
|     | (rr) | 'self test' means:                                                                                                                                                                                                                                                   |
|     |      | tests run cyclically and automatically by the recording<br>equipment to detect faults;                                                                                                                                                                               |
|     | (ss) | 'time measurement' means:                                                                                                                                                                                                                                            |
|     |      | a permanent digital record of the coordinated universal date                                                                                                                                                                                                         |

and time (UTC);

## (tt) 'time adjustment' means:

an adjustment of current time; this adjustment can be automatic, using the time provided by the GNSS receiver as a reference, or performed in calibration mode;

## (uu) 'tyre size' means:

the designation of the dimensions of the tyres (external driving wheels) in accordance with Council Directive 92/23/EEC (1) as last amended;

(vv) 'vehicle identification' means:

numbers identifying the vehicle: vehicle registration number (VRN) with indication of the registering Member State and vehicle identification number (VIN) (2);

(ww) for computing sake in the recording equipment 'week' means:

the period between 00:00 hours UTC on Monday and 24:00 UTC on Sunday;

(xx) 'workshop card' means:

a tachograph card issued by the authorities of a Member State to designated staff of a tachograph manufacturer, a fitter, a vehicle manufacturer or a workshop, approved by that Member State, which identifies the cardholder and allows for the testing, calibration and activation of tachographs, and/or downloading from them;

(yy) 'adaptor' means:

a device, providing a signal permanently representative of vehicle speed and/or distance travelled, other than the one used for the independent movement detection, and which is:

**▼M3**

**▼B**

- installed and used only in M1 and N1 type vehicles, as defined in Article 4 of Regulation (EU) 2018/858 of the European Parliament and of the Council (3),
- installed where it is not mechanically possible to install any other type of existing motion sensor which is otherwise compliant with the provisions of this Annex and its Appendixes 1 to 15,

### **M3**

<sup>(1)</sup> Council Directive 92/23/EEC of 31 March 1992 relating to tyres for motor vehicles and their trailers and to their fitting (OJ L 129, 14.5.1992, p. 95).

<sup>(2)</sup> Council Directive 76/114/EEC of 18 December 1975 on the approximation of the laws of the Member States relating to statutory plates and inscriptions for motor vehicles and their trailers, and their location and method of attachment (OJ L 24, 30.1.1976, p. 1).

<sup>(3)</sup> Regulation (EU) 2018/858 of the European Parliament and of the Council of 30 May 2018 on the approval and market surveillance of motor vehicles and their trailers, and of systems, components and separate technical units intended for such vehicles, amending Regulations (EC) No 715/2007 and (EC) No 595/2009 and repealing Directive 2007/46/EC (OJ L 151, 14.6.2018, p.1)

|  | — installed between the vehicle unit and where the speed/ |  |  |  |  |
|--|-----------------------------------------------------------|--|--|--|--|
|  | distance impulses are generated by integrated sensors or  |  |  |  |  |
|  | alternative interfaces,                                   |  |  |  |  |

— seen from a vehicle unit, the adaptor behaviour is the same as if a motion sensor, compliant with the provisions of this Annex and its Appendixes 1 to 16, was connected to the vehicle unit;

use of such an adaptor in those vehicles described above shall allow for the installation and correct use of a vehicle unit compliant with all the requirements of this Annex,

for those vehicles, the smart tachograph includes cables, an adaptor, and a vehicle unit;

(zz) 'data integrity' means:

the accuracy and consistency of stored data, indicated by an absence of any alteration in data between two updates of a data record. Integrity implies that the data is an exact copy of the original version, e.g. that it has not been corrupted in the process of being written to, and read back from, a tachograph card or a dedicated equipment or during transmission via any communications channel;

**▼M3**

**▼B**

(aaa) reserved for future use;

(bbb) 'smart tachograph' system means:

the recording equipment, tachograph cards and the set of all directly or indirectly interacting equipment during their construction, installation, use, testing and control, such as cards, remote communication reader and any other equipment for data downloading, data analysis, calibration, generating, managing or introducing security elements, etc.;

**▼M3**

(ccc) 'introduction date' means:

the date set out in Regulation (EU) No 165/2014 as from which vehicles registered for the first time shall be fitted with a tachograph in accordance with this Regulation;

## **B**

(ddd) 'protection profile' means:

a document used as part of certification process according Common Criteria, providing implementation independent specification of information assurance security requirements;

(eee) 'GNSS accuracy':

in the context of recording the position from a Global Navigation Satellite System (GNSS) with tachographs, means the value of the horizontal dilution of precision (HDOP) calculated as the minimum of the HDOP values collected on the available GNSS systems;

## (fff) 'accumulated driving time' means:

a value representing the total accumulated number of minutes of driving of a particular vehicle.

The accumulated driving time value is a free running count of all minutes regarded as DRIVING by the monitoring of driving activities function of the recording equipment, and is only used for triggering the recording of the vehicle position, every time a multiple of three hours of accumulated driving is reached. The accumulation is started at the recording equipment activation. It is not affected by any other condition, like out of scope or ferry/train crossing.

The accumulated driving time value is not intended to be displayed, printed, or downloaded.

**▼B**

## 2. GENERAL CHARACTERISTICS AND FUNCTIONS OF THE RECORDING EQUIPMENT

### 2.1 **General characteristics**

The purpose of the recording equipment is to record, store, display, print, and output data related to driver activities.

Any vehicle fitted with the recording equipment complying with the provisions of this Annex, must include a speed display and an odometer. These functions may be included within the recording equipment.

- (1) The recording equipment includes cables, a motion sensor, and a vehicle unit.
- (2) The interface between motion sensors and vehicle units shall comply with the requirements specified in Appendix 11.
- (3) The vehicle unit shall be connected to global navigation satellite system(s), as specified in Appendix 12.
- (4) The vehicle unit shall communicate with remote early detection communication readers, as specified in Appendix 14.

## **M3**

(5) The vehicle unit shall include an ITS interface, which is specified in Appendix 13.

> The recording equipment may be connected to other facilities through additional interfaces and/or through the ITS interface.

### **B**

(6) Any inclusion in or connection to the recording equipment of any function, device, or devices, approved or otherwise, shall not interfere with, or be capable of interfering with, the proper and secure operation of the recording equipment and the provisions of this Regulation.

> Recording equipment users identify themselves to the equipment via tachograph cards.

(7) The recording equipment provides selective access rights to data and functions according to user's type and/or identity.

The recording equipment records and stores data in its data memory, in the remote communication facility and in tachograph cards.

This is done in accordance with the applicable Union legislation regarding data protection and in compliance with Article 7 of Regu-

- lation (EU) No 165/2014.
- **▼B**

**▼M3**

## 2.2 **Functions**

- (8) The recording equipment shall ensure the following functions:
  - monitoring cards insertions and withdrawals,
  - speed, distance and position measurement,
  - time measurement,
  - monitoring driver activities,
  - monitoring driving status,

**▼M3**

- drivers manual entries:
  - entry of places where daily work periods begin and/or end,
  - manual entry of driver activities and driver consent for ITS interface,
  - entry of specific conditions,
  - entry of load/unload operations,

**▼B**

- company locks management,
- monitoring control activities,
- detection of events and/or faults,
- built-in and self-tests,
- reading from data memory,
- recording and storing in data memory,
- reading from tachograph cards,
- recording and storing in tachograph cards,
- displaying,
- printing,
- warning,
- data downloading to external media,

— remote communication for targeted roadside checks,

- output data to additional facilities,
- calibration,
- roadside calibration check,
- time adjustment,

- **▼M3**
- monitoring border crossings,
  - software update.

**▼B**

### 2.3 **Modes of operation**

- (9) The recording equipment shall possess four modes of operation:
  - operational mode,
  - control mode,
  - calibration mode,
  - company mode.
- (10) The recording equipment shall switch to the following mode of operation according to the valid tachograph cards inserted into the card interface devices. In order to determine the mode of operation, the tachograph card generation is irrelevant, provided the inserted card is valid. A first generation workshop card shall always be considered as non-valid when it is inserted in a second generation VU.

| Mode of operation | Driver slot |             |              |                 |              |
|-------------------|-------------|-------------|--------------|-----------------|--------------|
| Co-driver slot    | No card     | Driver card | Control card | Workshop card   | Company card |
| No card           | Operational | Operational | Control      | Calibration     | Company      |
| Driver card       | Operational | Operational | Control      | Calibration     | Company      |
| Control card      | Control     | Control     | Control (*)  | Operational     | Operational  |
| Workshop card     | Calibration | Calibration | Operational  | Calibration (*) | Operational  |
| Company card      | Company     | Company     | Operational  | Operational     | Company (*)  |

(\*) In these situations the recording equipment shall use only the tachograph card inserted in the driver slot.

- (11) The recording equipment shall ignore non-valid cards inserted, except displaying, printing or downloading data held on an expired card which shall be possible.
- (12) All functions listed in 2.2. shall work in any mode of operation with the following exceptions:
  - the calibration function is accessible in the calibration mode only,
  - the roadside calibration checking function is accessible in the control mode only,
  - the company locks management function is accessible in the company mode only,

| ▼B        | — the monitoring of control activities function is operational in the control mode only,                                                                                                                                                                                                                                                                                                                                                                              |
|-----------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| ▼M3       | — The downloading function is not accessible in the operational mode, except:                                                                                                                                                                                                                                                                                                                                                                                         |
|           | (a) as provided for in requirement 193,                                                                                                                                                                                                                                                                                                                                                                                                                               |
|           | (b) downloading a driver card when no other card type is inserted into the VU.                                                                                                                                                                                                                                                                                                                                                                                        |
| ▼B        | (13) The recording equipment can output any data to display, printer or external interfaces with the following exceptions:                                                                                                                                                                                                                                                                                                                                            |
|           | — in the operational mode, any personal identification (surname and first name(s)) not corresponding to a tachograph card inserted shall be blanked and any card number not corresponding to a tachograph card inserted shall be partially blanked (every odd character — from left to right — shall be blanked),                                                                                                                                                     |
| ▼M3       | — in the company mode, driver related data (requirements 102, 105, 108, 133a and 133e) can be output only for periods where no lock exists or no other company holds a lock (as identified by the first 13 digits of the company card number),                                                                                                                                                                                                                        |
| ▼B        | — when no card is inserted in the recording equipment, driver related data can be output only for the current and 8 previous calendar days,                                                                                                                                                                                                                                                                                                                           |
| ▼M3       | — personal data recorded and produced by either the tachograph or the tachograph cards shall not be output through the ITS interface of the VU unless the consent of the driver to whom the data relates is verified,                                                                                                                                                                                                                                                 |
| ▼M1       | — the vehicle units have a normal operations validity period of 15 years, starting with the vehicle unit certificates effective date, but vehicle units can be used for additional 3 months, for data downloading only.                                                                                                                                                                                                                                               |
| ▼B<br>2.4 | Security                                                                                                                                                                                                                                                                                                                                                                                                                                                              |
| ▼M1       | The system security aims at protecting the data memory in such a way as to prevent unauthorised access to and manipulation of the data and detecting any such attempts, protecting the integrity and authenticity of data exchanged between the motion sensor and the vehicle unit, protecting the integrity and authenticity of data exchanged between the recording equipment and the tachograph cards, protecting the integrity and authenticity of data exchanged |

protecting the confidentiality, integrity and authenticity of data exchanged through the remote early detection communication for control purposes, and verifying the integrity and authenticity of

data downloaded.

(14) In order to achieve the system security, the following components shall meet the security requirements specified in their Protection Profiles, as required in Appendix 10: — vehicle unit, — tachograph card, — motion sensor, — external GNSS facility (this Profile is only needed and applicable for the external GNSS facility variant). 3. CONSTRUCTION AND FUNCTIONAL REQUIREMENTS FOR RECORDING EQUIPMENT 3.1 **Monitoring cards insertion and withdrawal** (15) The recording equipment shall monitor the card interface devices to detect card insertions and withdrawals. (16) Upon card insertion (or remote card authentication) the recording equipment shall detect whether the card is a valid tachograph card in accordance with definition (ee) in section 1, and in such a case identify the card type and the card generation. For checking if a card has already been inserted, the recording equipment shall use the tachograph card data stored in its data memory, as set out in requirement 133. (17) First generation tachograph cards shall be considered as non-valid by the recording equipment, after the possibility of using first generation tachograph cards has been suppressed by a workshop, in compliance with Appendix 15 (req. MIG003). (18) First generation workshop cards which are inserted in the second generation recording equipment shall be considered as non-valid.

- (19) The recording equipment shall be so designed that the tachograph cards are locked in position on their proper insertion into the card interface devices.
- **▼M3**
- (20) The withdrawal of tachograph cards may function only when the vehicle is stopped and after the relevant data have been stored on the cards. The withdrawal of the card shall require positive action by the user.
- **▼B**

## 3.2 **Speed, position and distance measurement**

- (21) The motion sensor (possibly embedded in the adaptor) is the main source for speed and distance measurement.
- (22) This function shall continuously measure and be able to provide the odometer value corresponding to the total distance travelled by the vehicle using the pulses provided by the motion sensor.

### **B**

**▼M3**

**▼B**

**▼M3**

- (23) This function shall continuously measure and be able to provide the speed of the vehicle using the pulses provided by the motion sensor.
- (24) The speed measurement function shall also provide the information whether the vehicle is moving or stopped. The vehicle shall be considered as moving as soon as the function detects more than 1 imp/sec for at least 5 seconds from the motion sensor, otherwise the vehicle shall be considered as stopped.
- (25) Devices displaying speed (speedometer) and total distance travelled (odometer) installed in any vehicle fitted with a recording equipment complying with the provisions of this Regulation, shall comply with the requirements relating to maximum tolerances (see 3.2.1 and 3.2.2) laid down in this Annex.

- (26) To detect manipulation of motion data, information from the motion sensor shall be corroborated by vehicle motion information derived from the GNSS receiver and by other source(s) independent from the motion sensor. At least another independent vehicle motion source shall be inside the VU without the need of an external interface.
- (27) This function shall measure the position of the vehicle in order to allow for the recording of:
  - positions where the driver and/or the co-driver begins his daily work period;
  - positions where the accumulated driving time reaches a multiple of three hours;
  - positions where the vehicle has crossed the border of a country;
  - positions where operations of load/unload have been carried out;
  - positions where the driver and/or the co-driver ends his daily work period.

## **B**

### 3.2.1 *Measurement of distance travelled*

- (28) The distance travelled may be measured either:
  - so as to cumulate both forward and reverse movements, or
  - so as to include only forward movement.
- (29) The recording equipment shall measure distance from 0 to 9 999 999,9 km.

| ▼ <b>B</b>  |       | (30) | Distance measured shall be within the following tolerances<br>(distances of at least 1 000 m.):<br><br>— ± 1 % before installation,<br><br>— ± 2 % on installation and periodic inspection,<br><br>— ± 4 % in use.                                                                                                                                                                                                                                                                                                                                                                                                                                                      |
|-------------|-------|------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| ▼ <b>M3</b> |       |      | The tolerances shall not be used to intentionally alter the<br>distance measured.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                       |
| ▼ <b>B</b>  |       | (31) | Distance measured shall have a resolution better than or<br>equal to 0,1 km.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                            |
|             | 3.2.2 |      | <i>Measurement of speed</i>                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                             |
|             |       | (32) | The recording equipment shall measure speed from 0 to<br>220 km/h.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                      |
| ▼ <b>M3</b> |       | (33) | To ensure a maximum tolerance on speed displayed of<br>± 6 km/h in use, and taking into account:<br><br>— a ± 2 km/h tolerance for input variations (tyre<br>variations, ...),<br><br>— a ± 1 km/h tolerance in measurements made during<br>installation or periodic inspections,<br><br>the recording equipment shall, for speeds between 20<br>and 180 km/h, and for characteristic coefficients of the<br>vehicle between 2 400 and 25 000 imp/km, measure the<br>speed with a tolerance of ± 1 km/h (at constant speed).<br><br>Note: The resolution of data storage brings an additional<br>tolerance of ± 0,5 km/h to speed stored by the recording<br>equipment. |
| ▼ <b>B</b>  |       | (34) | The speed shall be measured correctly within the normal<br>tolerances within 2 seconds of the end of a speed change<br>when the speed has changed at a rate up to 2 m/s2.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                               |
|             |       | (35) | Speed measurement shall have a resolution better than or<br>equal to 1 km/h.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                            |
|             | 3.2.3 |      | <i>Measurement of position</i>                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                          |
|             |       | (36) | The recording equipment shall measure the absolute<br>position of the vehicle using the GNSS receiver.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                  |
| ▼ <b>M3</b> |       | (37) | The absolute position shall be measured in geographical<br>coordinates of latitude and longitude in degrees and<br>minutes with a resolution of 1/10 of a minute.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                       |
| ▼ <b>B</b>  | 3.3   |      | <i>Time measurement</i>                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                 |

(38) The time measurement function shall measure permanently and digitally provide UTC date and time.

- (39) UTC date and time shall be used for dating data inside the recording equipment (recordings, data exchange) and for all printouts specified in Appendix 4 'Printouts'.
- (40) In order to visualise the local time, it shall be possible to change the offset of the time displayed, in half hour steps. No other offsets than negative or positive multiples of half hours shall be allowed.
- (41) Time drift shall be ± 1 second per day or less, in temperature conditions in accordance with requirement 213, in the absence of any time adjustment.
- (41a) Time accuracy when time is adjusted by workshops in accordance with requirement 212 shall be 3 seconds or better.
- (41b) The vehicle unit shall include a drift counter, which computes the maximal time drift since the last time adjustment in accordance with point 3.23. The maximal time drift shall be defined by the vehicle unit manufacturer and shall not exceed 1 second per day, as set out in requirement 41.
- (41c) The drift counter shall be reset to 1 second after each time adjustment of the recording equipment in accordance with point 3.23. This includes:
  - automatic time adjustments,
  - time adjustments performed in calibration mode.

- (42) Time measured shall have a resolution better than or equal to 1 second.
- (43) Time measurement shall not be affected by an external power supply cut-off of less than 12 months in type approval conditions.

### 3.4 **Monitoring driver activities**

- (44) This function shall permanently and separately monitor the activities of one driver and one co-driver.
- (45) Driver activity shall be DRIVING, WORK, AVAILABIL-ITY or BREAK/REST.
- (46) It shall be possible for the driver and/or the co-driver to manually select WORK, AVAILABILITY or BREAK/REST.
- (47) When the vehicle is moving, DRIVING shall be selected automatically for the driver and AVAILABILITY shall be selected automatically for the co-driver.
- (48) When the vehicle stops, WORK shall be selected automatically for the driver.

### **B**

- (49) The first change of activity to BREAK/REST or AVAIL-ABILITY arising within 120 seconds of the automatic change to WORK due to the vehicle stop shall be assumed to have happened at the time of vehicle stop (therefore possibly cancelling the change to WORK).
- (50) This function shall output activity changes to the recording functions at a resolution of one minute.
- (51) Given a calendar minute, if DRIVING is registered as the activity of both the immediately preceding and the immediately succeeding minute, the whole minute shall be regarded as DRIVING.
- (52) Given a calendar minute that is not regarded as DRIVING according to requirement 051, the whole minute shall be regarded to be of the same type of activity as the longest continuous activity within the minute (or the latest of the equally long activities).
- (53) This function shall also permanently monitor the continuous driving time and the cumulative break time of the driver.

### 3.5 **Monitoring driving status**

- (54) This function shall permanently and automatically monitor the driving status.
- (55) The driving status CREW shall be selected when two valid driver cards are inserted in the equipment, the driving status SINGLE shall be selected in any other case.

### 3.6 **Driver's entries**

- 3.6.1 *Entry of places where daily work periods begin and/or end*
  - (56) This function shall allow for the entry of places where, according to the driver and/or the co-driver, his daily work periods begin and/or end.

### **M3**

- (57) Places are defined as the country and, in addition where applicable, the region.
- (58) Upon driver (or workshop) card withdrawal, the recording equipment shall display the current place of the vehicle on the basis of the GNSS information, and of stored digital map in accordance with point 3.12.19, and shall request the cardholder to confirm or to manually rectify the place.
- (59) The place entered in accordance with requirement 58 shall be considered as the place where the daily work period ends. It shall be recorded in the relevant driver (or workshop) card as a temporary record, and may therefore be later overwritten.

Under the following conditions temporary entry made at last card withdrawal is validated (i.e. shall not be overwritten anymore):

— entry of a place where the current daily work period begins during manual entry according to requirement (61);

### **M1**

| ▼M3         |      |                                                                                                                                                                                                                                                                                                         |
|-------------|------|---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
|             |      | — the next entry of a place where the current daily work<br>period begins if the card holder does not enter any<br>place where the work period begins or ended during<br>the manual entry according to requirement (61).                                                                                |
|             |      | Under the following conditions temporary entry made at<br>last card withdrawal is overwritten and the new value is<br>validated:                                                                                                                                                                        |
|             |      | — the next entry of a place where the current daily work<br>period ends if the card holder does not enter any place<br>where the work period begins or ended during the<br>manual input according to requirement (61).                                                                                  |
| ▼B          | (60) | It shall be possible to input places where daily work<br>periods begin and/or end through commands in the<br>menus. If more than one such input is done within one<br>calendar minute, only the last begin place input and the<br>last end place input done within that time shall be kept<br>recorded. |
| ▼M3         |      | The recording equipment shall display the current place of<br>the vehicle on the basis of the GNSS information, and of<br>stored digital map(s) in accordance with point 3.12.19 and<br>shall request the driver to confirm or to manually rectify<br>the place.                                        |
| ▼B<br>3.6.2 |      | <i>Manual entry of driver activities and driver consent for ITS<br/>interface</i>                                                                                                                                                                                                                       |
| ▼M3         | (61) | Upon driver (or workshop) card insertion, and only at this<br>time, the recording equipment shall allow manual entries<br>of activities. Manual entries of activities shall be<br>performed using local time and date values of the time<br>zone (UTC offset) currently set for the vehicle unit.       |
|             |      | At driver or workshop card insertion the cardholder shall<br>be reminded of:                                                                                                                                                                                                                            |
|             |      | — the date and time of his last card withdrawal;                                                                                                                                                                                                                                                        |
|             |      | — optionally: the local time offset currently set for the<br>vehicle unit.                                                                                                                                                                                                                              |

At the first insertion of a given driver card or workshop card currently unknown to the vehicle unit, the cardholder shall be invited to express his consent for tachograph related personal data output through the ITS interface. For checking if a card has already been inserted, the recording equipment shall use the tachograph card data stored in its data memory, as set out in requirement 133.

At any moment, the driver (resp. workshop) consent can be enabled or disabled through commands in the menu, provided the driver (resp. workshop) card is inserted.

It shall be possible to input activities with the following restrictions:

- Activity type shall be WORK, AVAILABILITY or BREAK/REST;
- Start and end times for each activity shall be within the period of the last card withdrawal – current insertion only;
- Activities shall not be allowed to overlap mutually in time.

It shall be possible to make manual entries, if required, at the first insertion of a previously unused driver (or workshop) card.

The procedure for manual entries of activities shall include as many consecutive steps as necessary to set a type, a start time and an end time for each activity. For any part of the time period between last card withdrawal and current card insertion, the cardholder shall have the option not to declare any activity.

During the manual entries associated with card insertion and if applicable, the card holder shall have the opportunity to input:

- a place where a previous daily work period ended, associated to the relevant time (thus overwriting and validating the entry made at the last card withdrawal),
- a place where the current daily work period begins, associated to the relevant time (thus validating a temporary entry made at last card withdrawal).

For the place where the current daily work period begins entered at the current card insertion, the recording equipment shall display the current place of the vehicle on the basis of the GNSS information, and of stored digital map(s) in accordance with point 3.12.19, and shall request the driver to confirm or to manually rectify the place.

If the card holder does not enter any place where the work period begins or ended, during the manual entries associated with card insertion, this shall be considered as a declaration that his work period has not changed since the last card withdrawal. The next entry of a place where a previous daily work period ends shall then overwrite the temporary entry made at the last card withdrawal.

If a place is entered, it shall be recorded in the relevant tachograph card.

Manual entries shall be interrupted if:

— the card is withdrawn or,

— the vehicle is moving and the card is in the driver slot.

Additional interruptions are allowed, e.g. a timeout after a certain period of user inactivity. If manual entries are interrupted, the recording equipment shall validate any complete place and activity entries (having either unambiguous place and time, or activity type, begin time and end time) already made.

If a second driver or workshop card is inserted while manual entries of activities are in progress for a previously inserted card, the manual entries for this previous card shall be allowed to be completed before manual entries start for the second card.

The cardholder shall have the option to insert manual entries according to the following minimum procedure:

- Enter activities manually, in chronological order, for the period last card withdrawal – current insertion.
- Begin time of the first activity shall be set to card withdrawal time. For each subsequent entry, the start time shall be preset to immediately follow the end time of the previous entry. Activity type and end time shall be selected for each activity.

The procedure shall end when the end time of a manually entered activity equals the card insertion time.

The recording equipment shall allow drivers and workshops to alternately upload manual entries that need to be entered during the procedure through the ITS interface specified in Appendix 13 and, optionally, through other interfaces.

The recording equipment shall allow the card holder to modify any activity manually entered, until validation by selection of a specific command. Thereafter, any such modification shall be forbidden.

**▼B** 3.6.3 *Entry of specific conditions*

**▼M3**

- (62) The recording equipment shall allow the driver to enter, in real time, the following two specific conditions:
  - 'OUT OF SCOPE' (begin, end),

— 'FERRY / TRAIN CROSSING' (begin, end).

A 'FERRY / TRAIN CROSSING' shall not occur if an 'OUT OF SCOPE' condition is opened. If an 'OUT OF SCOPE' condition is opened, the recording equipment shall not allow users to enter a 'FERRY / TRAIN CROSSING' begin flag.

An opened 'OUT OF SCOPE' condition must be automatically closed, by the recording equipment, if a driver card is inserted or withdrawn.

An opened 'OUT OF SCOPE' condition shall inhibit the following events and warnings:

- Driving without an appropriate card,
- Warnings associated with continuous driving time.

The driver shall enter the FERRY / TRAIN CROSSING begin flag immediately after selecting BREAK/REST on the ferry or train.

An opened FERRY / TRAIN CROSSING must be ended by the recording equipment when any of the following options occurs:

- the driver manually ends the FERRY/TRAIN CROSSING, which shall occur upon arrival to destination of the ferry/ train, before driving off the ferry/train,
- an 'OUT OF SCOPE' condition is opened,
- the driver ejects his card,
- driver activity is computed as DRIVING during a calendar minute in accordance with point 3.4.

If more than one specific conditions entry of the same type is done within one calendar minute, only the last one shall be kept recorded.

### 3.6.4 **Entry of load/unload operation**

(62a) The recording equipment shall allow the driver to enter and confirm, in real time, information indicating that the vehicle is being loaded, unloaded or that simultaneous load/unload operation is being performed.

> If more than one load/unload operation entry of the same type is done within one calendar minute, only the last one shall be kept recorded.

- (62b) Load, unload or simultaneous load/unload operations shall be recorded as separate events.
- (62c) The load/unload information shall be entered before the vehicle leaves the place where the load/unload operation is carried out.

## 3.7 **Company locks management**

- (63) This function shall allow the management of the locks placed by a company to restrict data access in company mode to itself.
- (64) Company locks consist in a start date/time (lock-in) and an end date/time (lock-out) associated with the identification of the company as denoted by the company card number (at lock-in).
- (65) Locks may be turned 'in' or 'out' in real time only.
- (66) Locking-out shall only be possible for the company whose lock is 'in' (as identified by the first 13 digits of the company card number), or,
- (67) Locking-out shall be automatic if another company locks in.
- (68) In the case where a company locks in and where the previous lock was for the same company, then it will be assumed that the previous lock has not been turned 'out' and is still 'in'.

### 3.8 **Monitoring control activities**

- (69) This function shall monitor DISPLAYING, PRINTING, VU and card DOWNLOADING, and ROADSIDE CALI-BRATION check activities carried while in control mode.
- (70) This function shall also monitor OVER SPEEDING CONTROL activities while in control mode. An over speeding control is deemed to have happened when, in control mode, the 'over speeding' printout has been sent to the printer or to the display, or when 'events and faults' data have been downloaded from the VU data memory.

### 3.9 **Detection of events and/or faults**

(71) This function shall detect the following events and/or faults:

### 3.9.1 *'Insertion of a non-valid card' event*

(72) This event shall be triggered at the insertion of any non-valid card, at the insertion of a driver card already replaced and/or when an inserted valid card expires.

## 3.9.2 *'Card conflict' event*

(73) This event shall be triggered when any of the valid cards combination noted X in the following table arises:

| Card conflict  |               | Driver slot |             |              |               |              |
|----------------|---------------|-------------|-------------|--------------|---------------|--------------|
|                |               | No card     | Driver card | Control card | Workshop card | Company card |
| Co-driver slot | No card       |             |             |              |               |              |
|                | Driver card   |             |             |              | X             |              |
|                | Control card  |             |             | X            | X             | X            |
|                | Workshop card |             | X           | X            | X             | X            |
|                | Company card  |             |             | X            | X             | X            |

### 3.9.3 *'Time overlap' event*

(74) This event shall be triggered when the date / time of last withdrawal of a driver card, as read from the card, is later than the current date / time of the recording equipment in which the card is inserted.

### 3.9.4 *'Driving without an appropriate card' event*

(75) This event shall be triggered for any valid tachograph cards combination noted X in the following table, when driver activity changes to DRIVING, or when there is a change of the mode of operation while driver activity is DRIVING:

| Driving without an appropriate<br>card |                        | Driver slot            |             |              |               |              |
|----------------------------------------|------------------------|------------------------|-------------|--------------|---------------|--------------|
|                                        |                        | No (or non-valid) card | Driver card | Control card | Workshop card | Company card |
| Co-driver slot                         | No (or non-valid) card | X                      | X           | X            | X             |              |
|                                        | Driver card            | X                      | X           | X            | X             |              |
|                                        | Control card           | X                      | X           | X            | X             |              |
|                                        | Workshop card          | X                      | X           | X            | X             |              |
|                                        | Company card           | X                      | X           | X            | X             |              |

## 3.9.5 *'Card insertion while driving' event*

(76) This event shall be triggered when a tachograph card is inserted in any slot, while driver activity is DRIVING.

## 3.9.6 *'Last card session not correctly closed' event*

(77) This event shall be triggered when at card insertion the recording equipment detects that, despite the provisions laid down in paragraph 3.1., the previous card session has not been correctly closed (the card has been withdrawn before all relevant data have been stored on the card). This event shall be triggered by driver and workshop cards only.

## 3.9.7 *'Over speeding' event*

(78) This event shall be triggered for each over speeding.

### 3.9.8 *'Power supply interruption' event*

- (79) This event shall be triggered, while not in calibration or control mode, in case of any interruption exceeding 200 milliseconds of the power supply of the motion sensor and/or of the vehicle unit. The interruption threshold shall be defined by the manufacturer. The drop in power supply due to the starting of the engine of the vehicle shall not trigger this event.
- 3.9.9 *'Communication error with the remote communication facility' event*
  - (80) This event shall be triggered, **while not in calibration mode**, when the remote communication facility does not acknowledge the successful reception of remote communication data sent from the vehicle unit for more than three attempts.

### 3.9.10 *'Absence of position information from GNSS receiver' event*

(81) This event shall be triggered, **while not in calibration mode**, in case of absence of position information originating from the GNSS receiver (whether internal or external) for more than three hours of accumulated driving time.

3.9.11 *'Communication error with the external GNSS facility' event*

(82) This event shall be triggered, **while not in calibration mode**, in case of interruption of the communication between the external GNSS facility and the vehicle unit for more than 20 continuous minutes, when the vehicle is moving.

3.9.12 *'Motion data error' event*

## **M3**

(83) This event shall be triggered, **while not in calibration mode**, in case of interruption of the normal data flow between the motion sensor and the vehicle unit and/or in case of data integrity or data authentication error during data exchange between the motion sensor and the vehicle unit. This event shall also be triggered, **while not in calibration mode**, in case the speed calculated from the motion sensor pulses increases from 0 to more than 40 km/h within 1 second, and then stays above 40km/h during at least 3 seconds.

**▼B**

**▼M3**

3.9.13 *'Vehicle motion conflict' event*

(84) This event shall be triggered, as specified in Appendix 12, **while not in calibration mode**, in case motion information calculated from the motion sensor is contradicted by motion information calculated from the internal GNSS receiver or from the external GNSS facility or by other independent source(s) in accordance with requirement 26. This event shall not be triggered during a ferry/train crossing.

### 3.9.14 *'Security breach attempt' event*

(85) This event shall be triggered for any other event affecting the security of the motion sensor and/or of the vehicle unit and/or the external GNSS facility as required in Appendix 10, while not in calibration mode.

## **M1**

3.9.15 *'Time conflict' event*

**▼M3**

(86) This event shall be triggered, **while not in calibration mode**, when the VU detects a discrepancy between the time of the vehicle unit's time measurement function and the time originating from the authenticated positions transmitted by the GNSS receiver or the external GNSS facility. A 'time discrepancy' is detected if the time difference exceeds ±3 seconds corresponding to the time accuracy set out in requirement 41a, the latter increased by the maximal time drift per day. This event shall be recorded together with the internal clock value of the recording equipment. The VU shall perform the check for triggering the 'time conflict' event right before the VU automatically re-adjusts the VU internal clock, in accordance with requirement 211.

### **B**

- 3.9.16 *'Card' fault*
  - (87) This fault shall be triggered when a tachograph card failure occurs during operation.
- 3.9.17 *'Recording equipment' fault*
  - (88) This fault shall be triggered for any of these failures, while not in calibration mode:
    - VU internal fault
    - Printer fault
    - Display fault
    - Downloading fault
    - Sensor fault
    - GNSS receiver or external GNSS facility fault
    - Remote Communication facility fault

- ITS interface fault.
- 3.9.18 *'GNSS anomaly' event*
  - (88a) This event shall be triggered, while not in calibration mode, when the GNSS receiver detects an attack, or when authentication of navigation messages has failed, as specified in Appendix 12. After a GNSS anomaly event has been triggered, the VU shall not generate other GNSS anomaly events for the next 10 minutes.

### 3.10 **Built-in and self-tests**

(89) **►M1** The recording equipment shall detect faults through self-tests and built-in-tests, according to the following table: ◄

| Sub-assembly to test                                  | Self-test            | Built-in-test          |
|-------------------------------------------------------|----------------------|------------------------|
| Software                                              |                      | Integrity              |
| Data memory                                           | Access               | Access, data integrity |
| Card interface devices                                | Access               | Access                 |
| Keyboard                                              |                      | Manual check           |
| Printer                                               | (up to manufacturer) | Printout               |
| Display                                               |                      | Visual check           |
| Downloading<br>(performed only during<br>downloading) | Proper operation     |                        |
| Sensor                                                | Proper operation     | Proper operation       |
| Remote<br>communication<br>facility                   | Proper operation     | Proper operation       |
| GNSS facility                                         | Proper operation     | Proper operation       |
| ITS interface                                         | Proper operation     |                        |

**▼M3**

**▼B**

## 3.11 **Reading from data memory**

(90) The recording equipment shall be able to read any data stored in its data memory.

## 3.12 **Recording and storing in the data memory**

## **M3**

For the purpose of this point,

- '365 days' is defined as 365 calendar days of average drivers activity in a vehicle. The average activity per day in a vehicle is defined as at least 6 drivers or co-drivers, 6 card insertion withdrawal cycles, and 256 activity changes. '365 days' therefore include at least 2 190 drivers or co-drivers, 2 190 card insertion withdrawal cycles, and 93 440 activity changes,
- the average number of place entries per day is defined as at least 6 entries where the daily work period begins and 6 entries where the daily work period ends, so that '365 days' include at least 4 380 place entries,
- the average number of positions per day when the accumulated driving time reaches a multiple of three hours is defined as at least 6 positions, so that '365 days' include at least 2 190 such positions,
- the average number of border crossings per day is defined as at least 20 crossings, so that '365 days' include at least 7 300 border crossings,

- the average number of load/unload operations per day is defined as at least 25 operations (irrespective of the type), so that '365 days' include at least 9 125 load/unload operations,
- times are recorded with a resolution of one minute, unless otherwise specified,
- odometer values are recorded with a resolution of one kilometre,
- speeds are recorded with a resolution of 1 km/h,
- positions (latitudes and longitudes) are recorded in degrees and minutes, with a resolution of 1/10 of minute, with the associated GNSS accuracy and acquisition time, and with a flag indicating whether the position has been authenticated.'.

- (91) Data stored into the data memory shall not be affected by an external power supply cut-off of less than twelve months in type approval conditions. In addition, data stored in the external remote communication facility, as defined in Appendix 14, shall not be affected by power-supply cut-off of less than 28 days.
- (92) The recording equipment shall be able to record and store implicitly or explicitly in its data memory the following:
- 3.12.1 *Equipment identification data*
- 3.12.1.1 V e h i c l e u n i t i d e n t i f i c a t i o n d a t a
  - (93) The recording equipment shall be able to store in its data memory the following vehicle unit identification data:
    - name of the manufacturer,
    - address of the manufacturer,
    - part number,
    - serial number,
    - VU generation,
    - ability to use first generation tachograph cards,
    - software version number,
    - software version installation date,
    - year of equipment manufacture,
    - approval number,

**▼M3**

- digital map version identifier (requirement 133l).
- (94) Vehicle unit identification data are recorded and stored once and for all by the vehicle unit manufacturer, except data which may be changed in case of software update in accordance with this Regulation, and the ability to use first generation tachograph cards.

| ▼B  | 3.12.1.2 | (95)  | Motion sensor identification data<br>The motion sensor shall be able to store in its memory the following identification data:<br>— name of the manufacturer,<br>— serial number,<br>— approval number,<br>— embedded security component identifier (e.g. internal chip/processor part number),<br>— operating system identifier (e.g. software version number).                                |
|-----|----------|-------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
|     |          | (96)  | Motion sensor identification data are recorded and stored once and for all in the motion sensor, by the motion sensor manufacturer.                                                                                                                                                                                                                                                             |
| ▼M3 |          | (97)  | The vehicle unit shall be able to record and store in its data memory the following data related to the 20 most recent successful pairings of motion sensors (if several pairings happen within one calendar day, only the first and the last one of the day shall be stored):                                                                                                                  |
| ▼B  |          |       | The following data shall be recorded for each of these pairings:<br>— motion sensor identification data:<br>— serial number<br>— approval number<br>— motion sensor pairing data:<br>— pairing date.                                                                                                                                                                                            |
|     | 3.12.1.3 | (98)  | Global Navigation Satellite Systems identification data<br>The external GNSS facility shall be able to store in its memory the following identification data:<br>— name of the manufacturer,<br>— serial number,<br>— approval number,<br>— embedded security component identifier (e.g. internal chip/processor part number),<br>— operating system identifier (e.g. software version number). |
|     |          | (99)  | The identification data are recorded and stored once and for all in the external GNSS facility, by the external GNSS facility manufacturer.                                                                                                                                                                                                                                                     |
| ▼M3 |          | (100) | The vehicle unit shall be able to record and store in its data memory the following data related to the 20 most                                                                                                                                                                                                                                                                                 |

(if several couplings happen within one calendar day, only the first and the last one of the day shall be stored).

The following data shall be recorded for each of these couplings:

- external GNSS facility identification data:
  - serial number,
  - approval number,
- external GNSS facility coupling data:
  - coupling date

### 3.12.2 *Keys and certificates*

- (101) The recording equipment shall be able to store a number of cryptographic keys and certificates, as specified in Appendix 11 part A and part B.
- 3.12.3 *Driver or workshop card insertion and withdrawal data*
  - (102) For each insertion and withdrawal cycle of a driver or workshop card in the equipment, the recording equipment shall record and store in its data memory:
    - the card holder's surname and first name(s) as stored in the card,
    - the card's number, issuing Member State and expiry date as stored in the card,
    - the card generation,
    - the insertion date and time,
    - the vehicle odometer value at card insertion,
    - the slot in which the card is inserted,
    - the withdrawal date and time,
    - the vehicle odometer value at card withdrawal,
    - the following information about the previous vehicle used by the driver, as stored in the card:
      - VRN and registering Member State,
      - VU generation (when available),
      - card withdrawal date and time,
    - a flag indicating whether, at card insertion, the card holder has manually entered activities or not.
  - (103) The data memory shall be able to hold these data for at least 365 days.
  - (104) When storage capacity is exhausted, new data shall replace oldest data.

### 3.12.4 *Driver activity data*

- (105) The recording equipment shall record and store in its data memory whenever there is a change of activity for the driver and/or the co-driver, and/or whenever there is a change of driving status, and/or whenever there is an insertion or withdrawal of a driver or workshop card:
  - the driving status (CREW, SINGLE),
  - the slot (DRIVER, CO-DRIVER),
  - the card status in the relevant slot (INSERTED, NOT INSERTED),
  - the activity (DRIVING, AVAILABILITY, WORK, BREAK/REST),
  - the date and time of the change.

INSERTED means that a valid driver or workshop card is inserted in the slot. NOT INSERTED means the opposite i.e. no valid driver or workshop card is inserted in the slot (e.g. a company card is inserted or no card is inserted)

Activity data manually entered by a driver are not recorded in the data memory.

- (106) The data memory shall be able to hold driver activity data for at least 365 days.
- (107) When storage capacity is exhausted, new data shall replace oldest data.

## **M1**

- 3.12.5 *Places and positions where daily work periods begin, end and/or where 3 hours accumulated driving time is reached*
  - (108) The recording equipment shall record and store in its data memory:
    - places and positions where the driver and/or co-driver begins his daily work period;
    - positions where the accumulated driving time reaches a multiple of three hours;
    - places and positions where the driver and/or the co-driver ends his daily work period.

## **B**

- (109) When the position of the vehicle is not available from the GNSS receiver at these times, the recording equipment shall use the latest available position, and the related date and time.
- (110) Together with each place or position, the recording equipment shall record and store in its data memory:

| ▼M3                                              | — the driver and/or co-driver card number and card issuing Member State,                                                                                                                                                                              |                                                                                                                                                                                                                                                                                                                                       |
|--------------------------------------------------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| ▼B                                               | — the card generation,                                                                                                                                                                                                                                |                                                                                                                                                                                                                                                                                                                                       |
|                                                  | — the date and time of the entry,                                                                                                                                                                                                                     |                                                                                                                                                                                                                                                                                                                                       |
| ▼M1                                              | — the type of entry (begin, end or 3 hours accumulated driving time),                                                                                                                                                                                 |                                                                                                                                                                                                                                                                                                                                       |
| ▼B                                               | — the related GNSS accuracy, date and time if applicable;                                                                                                                                                                                             |                                                                                                                                                                                                                                                                                                                                       |
|                                                  | — the vehicle odometer value,                                                                                                                                                                                                                         |                                                                                                                                                                                                                                                                                                                                       |
| ▼M3                                              | — a flag indicating whether the position has been authenticated.                                                                                                                                                                                      |                                                                                                                                                                                                                                                                                                                                       |
|                                                  | (110a) For places where the daily work period begins or ends entered during the manual entry procedure at card insertion in accordance with requirement 61, the current odometer value and position of the vehicle shall be stored.                   |                                                                                                                                                                                                                                                                                                                                       |
| ▼M1                                              | (111) The data memory shall be able to hold places and positions where daily work periods begin, end and/or where 3 hours accumulated driving time is reached for at least 365 days.                                                                  |                                                                                                                                                                                                                                                                                                                                       |
| ▼B                                               | (112) When storage capacity is exhausted, new data shall replace oldest data.                                                                                                                                                                         |                                                                                                                                                                                                                                                                                                                                       |
|                                                  | 3.12.6 Odometer data                                                                                                                                                                                                                                  |                                                                                                                                                                                                                                                                                                                                       |
|                                                  | (113) The recording equipment shall record in its data memory the vehicle odometer value and the corresponding date at midnight every calendar day.                                                                                                   |                                                                                                                                                                                                                                                                                                                                       |
|                                                  | (114) The data memory shall be able to store midnight odometer values for at least 365 calendar days.                                                                                                                                                 |                                                                                                                                                                                                                                                                                                                                       |
|                                                  | (115) When storage capacity is exhausted, new data shall replace oldest data.                                                                                                                                                                         |                                                                                                                                                                                                                                                                                                                                       |
|                                                  | 3.12.7 Detailed speed data                                                                                                                                                                                                                            |                                                                                                                                                                                                                                                                                                                                       |
| ▼M1                                              | (116) The recording equipment shall record and store in its data memory the instantaneous speed of the vehicle and the corresponding date and time at every second of at least the last 24 hours that the vehicle has been moving.                    |                                                                                                                                                                                                                                                                                                                                       |
| ▼B                                               | 3.12.8 Events data<br>For the purpose of this subparagraph, time shall be recorded with a resolution of 1 second.                                                                                                                                     |                                                                                                                                                                                                                                                                                                                                       |
| Event                                            | Storage rules                                                                                                                                                                                                                                         | Data to be recorded per event                                                                                                                                                                                                                                                                                                         |
| ▼B<br>Insertion of a non-valid<br>card           | — the 10 most recent events.                                                                                                                                                                                                                          | — date and time of event,<br>— card(s) type, number, issuing Member<br>State and generation of the card<br>creating the event.<br>— number of similar events that day                                                                                                                                                                 |
| Card conflict                                    | — the 10 most recent events.                                                                                                                                                                                                                          | — date and time of beginning of event,<br>— date and time of end of event,<br>— card(s) type, number, issuing Member<br>State and generation of the two cards<br>creating the conflict.                                                                                                                                               |
| Driving without an appro-<br>priate card         | — the longest event for each of the 10 last<br>days of occurrence,<br>— the 5 longest events over the last 365<br>days.                                                                                                                               | — date and time of beginning of event,<br>— date and time of end of event,<br>— card(s) type, number, issuing Member<br>State and generation of any card<br>inserted at beginning and/or end of the<br>event,<br>— number of similar events that day.                                                                                 |
| Card<br>insertion<br>while<br>driving            | — the last event for each of the 10 last<br>days of occurrence,                                                                                                                                                                                       | — date and time of the event,<br>— card(s) type, number, issuing Member<br>State and generation,<br>— number of similar events that day                                                                                                                                                                                               |
| ▼M3<br>Last card session not<br>correctly closed | — the 10 most recent events.                                                                                                                                                                                                                          | — date and time of card insertion,<br>— card(s) type, number, issuing Member<br>State and generation,<br>— last session data as read from the card<br>— date and time of card insertion.                                                                                                                                              |
| ▼B<br>Over speeding (1)                          | — the most serious event for each of the<br>10 last days of occurrence (i.e. the one<br>with the highest average speed),<br>— the 5 most serious events over the last<br>365 days.<br>— the first event having occurred after the<br>last calibration | — date and time of beginning of event,<br>— date and time of end of event,<br>— maximum speed measured during the<br>event,<br>— arithmetic average speed measured<br>during the event,<br>— card type, number, issuing Member State<br>and generation of the driver card (if<br>applicable),<br>— number of similar events that day. |

(117) The recording equipment shall record and store in its data memory the following data for each event detected according to the following storage rules:

| ▼B |  |
|----|--|
|    |  |

| ▼B  | Event                                                      | Storage rules                                                                                                                                                                                                                              | Data to be recorded per event                                                                                                                                                                                                                |
|-----|------------------------------------------------------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
|     | Power supply interruption (2)                              | — the longest event for each of the 10 last days of occurrence,<br>— the 5 longest events over the last 365 days.                                                                                                                          | — date and time of beginning of event,<br>— date and time of end of event,<br>— card(s) type, number, issuing Member State and generation of any card inserted at beginning and/or end of the event,<br>— number of similar events that day. |
|     | Communication error with the remote communication facility | — the longest event for each of the 10 last days of occurrence,<br>— the 5 longest events over the last 365 days.                                                                                                                          | — date and time of beginning of event,<br>— date and time of end of event,<br>— card(s) type, number, issuing Member State and generation of any card inserted at beginning and/or end of the event,<br>— number of similar events that day. |
|     | Absence of position information from GNSS receiver         | — the longest event for each of the 10 last days of occurrence,<br>— the 5 longest events over the last 365 days.                                                                                                                          | — date and time of beginning of event,<br>— date and time of end of event,<br>— card(s) type, number, issuing Member State and generation of any card inserted at beginning and/or end of the event,<br>— number of similar events that day. |
| ▼M1 | Communication error with the external GNSS facility        | — the longest event for each of the 10 last days of occurrence,<br>— the 5 longest events over the last 365 days.                                                                                                                          | — date and time of beginning of event,<br>— date and time of end of event,<br>— card(s) type, number, issuing Member State and generation of any card inserted at beginning and/or end of the event,<br>— number of similar events that day. |
| ▼B  | Motion data error                                          | — the longest event for each of the 10 last days of occurrence,<br>— the 5 longest events over the last 365 days.                                                                                                                          | — date and time of beginning of event,<br>— date and time of end of event,<br>— card(s) type, number, issuing Member State and generation of any card inserted at beginning and/or end of the event,<br>— number of similar events that day. |
|     | Vehicle motion conflict                                    | — the longest event for each of the 10 last days of occurrence,<br>— the 5 longest events over the last 365 days.                                                                                                                          | — date and time of beginning of event,<br>— date and time of end of event,<br>— card(s) type, number, issuing Member State and generation of any card inserted at beginning and/or end of the event,<br>— number of similar events that day. |
|     | Event                                                      | Storage rules                                                                                                                                                                                                                              | Data to be recorded per event                                                                                                                                                                                                                |
|     | Security breach attempt                                    | — the 10 most recent events per type of event.                                                                                                                                                                                             | — date and time of beginning of event,<br>— date and time of end of event (if relevant),<br>— card(s) type, number, issuing Member State and generation of any card inserted at beginning and/or end of the event,<br>— type of event.       |
| ▼M1 | Time conflict                                              | — the most serious event for each of the 10 last days of occurrence (i.e. the ones with the greatest difference between recording equipment date and time, and GNSS date and time).<br>— the 5 most serious events over the last 365 days. | — recording equipment date and time<br>— GNSS date and time,<br>— card(s) type, number, issuing Member State and generation of any card inserted at beginning and/or end of the event,<br>— number of similar events that day.               |
| ▼M3 | GNSS anomaly                                               | — the longest events for each of the 10 last days of occurrence,<br>— the 5 longest events over the last 365 days.                                                                                                                         | — date and time of beginning of event,<br>— date and time of end of event,<br>— card(s) type, number, issuing Member State and generation of any card inserted at beginning and/or end of the event,<br>— number of similar events that day. |

- (1) The recording equipment shall also record and store in its data memory:
  - the date and time of the last OVER SPEEDING CONTROL,
  - the date and time of the first over speeding following this OVER SPEEDING CONTROL,
  - the number of over speeding events since the last OVER SPEEDING CONTROL.
- (2) These data may be recorded at power supply reconnection only, times may be known with an accuracy to the minute.

3.12.9 *Faults data*

For the purpose of this subparagraph, time shall be recorded with a resolution of 1 second.

(118) The recording equipment shall attempt to record and store in its data memory the following data for each fault detected according to the following storage rules:

| ▼B |  |
|----|--|
|    |  |

| Fault                      | Storage rules                                                                                        | Data to be recorded per fault                                                                                                                                                                                               |
|----------------------------|------------------------------------------------------------------------------------------------------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| Card fault                 | — the 10 most recent driver card faults.                                                             | — date and time of beginning of fault,<br>— date and time of end of fault,<br>— card(s) type, number, issuing Member State and generation.                                                                                  |
| Recording equipment faults | — the 10 most recent faults for each type of fault,<br>— the first fault after the last calibration. | — date and time of beginning of fault,<br>— date and time of end of fault,<br>— type of fault,<br>— card(s) type, number and issuing Member State and generation of any card inserted at beginning and/or end of the fault. |

### 3.12.10 *Calibration data*

- (119) The recording equipment shall record and store in its data memory data relevant to:
  - known calibration parameters at the moment of activation,
  - its very first calibration following its activation,
  - its first calibration in the current vehicle (as identified by its VIN),
  - the 20 most recent calibrations (if several calibrations happen within one calendar day, only the first and the last one of the day shall be stored).
- (120) The following data shall be recorded for each of these calibrations:
  - purpose of calibration (activation, first installation, installation, periodic inspection),
  - workshop name and address,
  - workshop card number, card issuing Member State and card expiry date,
  - vehicle identification,
  - parameters updated or confirmed: w, k, l, tyre size, speed limiting device setting, odometer (old and new values), date and time (old and new values),
  - the types and identifiers of all the seals in place,
  - the serial numbers of the motion sensor, the external GNSS facility (if any), and the external remote communication facility (if any),
  - the by-default load type associated to the vehicle (load of either goods or passengers),

— the country in which the calibration has been performed, and the date time when the position used to determine this country was provided by the GNSS receiver.

**▼B**

(121) In addition, the recording equipment shall record and store in its data memory its ability to use first generation tachograph cards (still activated or not).

- (122) The motion sensor shall record and store in its memory the following motion sensor installation data:
  - first pairing with a VU (date, time, VU approval number, VU serial number),
  - last pairing with a VU (date, time, VU approval number, VU serial number).
- (123) The external GNSS facility shall record and store in its memory the following external GNSS facility installation data:
  - first coupling with a VU (date, time, VU approval number, VU serial number),
  - last coupling with a VU (date, time, VU approval number, VU serial number).
- 3.12.11 *Time adjustment data*
  - (124) The recording equipment shall record and store in its data memory data relevant to time adjustments performed in calibration mode outside the frame of a regular calibration (def. f)):
    - the most recent time adjustment,
    - the 5 largest time adjustments.
  - (125) The following data shall be recorded for each of these time adjustments:
    - date and time, old value,
    - date and time, new value,
    - workshop name and address,
    - workshop card number, card issuing Member State, card generation and card expiry date.

## 3.12.12 *Control activity data*

- (126) The recording equipment shall record and store in its data memory the following data relevant to the 20 most recent control activities:
  - date and time of the control,
  - control card number, card issuing Member State and card generation,
  - type of the control (displaying and/or printing and/or VU downloading and/or card downloading and/or roadside calibration checking).

- (127) In case of downloading, the dates of the oldest and of the most recent days downloaded shall also be recorded.
- 3.12.13 *Company locks data*
  - (128) The recording equipment shall record and store in its data memory the following data relevant to the 255 most recent company locks:
    - lock-in date and time,
    - lock-out date and time,
    - company card number, card issuing Member State and card generation,
    - company name and address.

Data previously locked by a lock removed from memory due to the limit above, shall be treated as not locked.

- 3.12.14 *Download activity data*
  - (129) The recording equipment shall record and store in its data memory the following data relevant to the last data memory downloading to external media while in company or in calibration mode:
    - date and time of downloading,
    - company or workshop card number, card issuing Member State and card generation,
    - company or workshop name.
- 3.12.15 *Specific conditions data*
  - (130) The recording equipment shall record in its data memory the following data relevant to specific conditions:
    - date and time of the entry,
    - type of specific condition.
  - (131) The data memory shall be able to hold specific conditions data for at least 365 days (with the assumption that on average, 1 condition is opened and closed per day). When storage capacity is exhausted, new data shall replace oldest data.

## 3.12.16 *Tachograph card data*

- (132) The recording equipment shall be able to store the following data related to the different tachograph cards in which had been used in the VU:
  - the tachograph card number and its serial number,

— the manufacturer of the tachograph card,

- the tachograph card type,
- the tachograph card version.
- (133) The recording equipment shall be able to store at least 88 such records.

3.12.17 *Border crossings*

- (133a) The recording equipment shall record and store in its data memory the following information about border crossings:
  - the country that the vehicle is leaving,
  - the country that the vehicle is entering,
  - the position where the vehicle has crossed the border.
- (133b) Together with countries and position, the recording equipment shall record and store in its data memory:
  - the driver and/or co-driver card number and card issuing Member State,
  - the card generation,
  - the related GNSS accuracy, date and time,
  - a flag indicating whether the position has been authenticated
  - the vehicle odometer value at the time of border crossing detection.
- (133c) The data memory shall be able to hold border crossings for at least 365 days.
- (133d) When storage capacity is exhausted, new data shall replace oldest data.
- 3.12.18 *Load/unload operations*
  - (133e) The recording equipment shall record and store in its data memory the following information about load and unload operations of the vehicle:
    - the type of operation (load, unload or simultaneous load/unload),
    - the position where the load/unload operation has occurred.
  - (133f) When the position of the vehicle is not available from the GNSS receiver at the time of the load/unload operation, the recording equipment shall use the latest available position, and the related date and time.
  - (133g) Together with the type of operation and position, the recording equipment shall record and store in its data memory:
    - the driver and/or co-driver card number and card issuing Member State,
    - the card generation,

- the date and time of the load/unload operation,
- the related GNSS accuracy, date and time if applicable,
- a flag indicating whether the position has been authenticated,
- the vehicle odometer value.
- (133h) The data memory shall be able to store load/unload operations for at least 365 calendar days.
- (133i) When storage capacity is exhausted, new data shall replace oldest data.

### 3.12.19 *Digital map*

- (133j) For the purpose of recording the position of the vehicle when the border of a country is crossed, the recording equipment shall store in its data memory a digital map.
- (133k) Allowed digital maps for supporting the border crossing monitoring function of the recording equipment shall be made available by the European Commission for download from a dedicated secured website, under various formats.
- (133l) For each of these maps, a version identifier and a hash value shall be available on the website.
- (133m) The maps shall feature:
  - a level of definition corresponding to NUTS level 0, according to the Nomenclature of Territorial Units for Statistics,
  - a scale of 1:1 million.
- (133n) Tachograph manufacturers shall select a map from the website and download it securely.
- (133o) Tachograph manufacturers shall only use a downloaded map from the website after having verified its integrity using the hash value of the map.
- (133p) The selected map shall be imported in the recording equipment by its manufacturer, under an appropriate format, but the semantic of the imported map shall remain unchanged.
- (133q) The manufacturer shall also store the version identifier of the map used in the recording equipment.
- (133r) It shall be possible to update or replace the stored digital map by a new one made available by the European Commission.
- (133s) Digital map updates shall be made using the software update mechanisms set up by the manufacturer, in application of requirements 226d and 226e, so that the recording equipment can verify the authenticity and integrity of a new imported map, before storing it and replacing the previous one.

| <b>▼M3</b> | (133t) | Tachograph manufacturers may add additional information to the basic map referred to in requirement (133m), for purposes other than recording border crossings, such as the borders of the EU regions, provided that the semantic of the basic map is not changed.                                                                                                                                                                                                                                                |
|------------|--------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| <b>▼B</b>  | 3.13   | Reading from tachograph cards                                                                                                                                                                                                                                                                                                                                                                                                                                                                                     |
|            | (134)  | The recording equipment shall be able to read from first and second generation tachograph cards, where applicable, the necessary data:<br>— to identify the card type, the card holder, the previously used vehicle, the date and time of the last card withdrawal and the activity selected at that time,<br>— to check that last card session was correctly closed,                                                                                                                                             |
| <b>▼M3</b> |        | — to compute the driver's continuous driving time, cumulative break time and accumulated driving times for the previous and the current week,                                                                                                                                                                                                                                                                                                                                                                     |
| <b>▼B</b>  |        | — to print requested printouts related to data recorded on a driver card,<br>— to download a driver card to external media.                                                                                                                                                                                                                                                                                                                                                                                       |
|            | (135)  | This requirement only applies to first generation tachograph cards if their use has not been suppressed by a workshop.<br>In case of a reading error, the recording equipment shall try again, three times maximum, the same read command, and then if still unsuccessful, declare the card faulty and non-valid.                                                                                                                                                                                                 |
| <b>▼M3</b> | (135a) | The structure in the 'TACHO_G2' application depends on the version. Version 2 cards contain additional Elementary Files to the ones of version 1 cards, in particular:<br>— in driver and workshop cards:<br>— EF Places_Authentication shall contain the authentication status of the vehicle positions stored in EF Places. A timestamp shall be stored with each authentication status, which shall be exactly the same as the date and time of the entry stored with the corresponding position in EF Places. |

— EF GNSS\_Places\_Authentication shall contain the authentication status of the vehicle positions stored in EF GNSS\_Places. A timestamp shall be stored with each authentication status, which shall be exactly the same as the date and time of the entry stored with the corresponding position in EF Places.

- EF Border\_Crossings, EF Load\_Unload\_Operations and EF Load\_Type\_Entries shall contain data related to border crossings, load/unload operations and load types.
- in workshop cards:
  - EF Calibration\_Add\_Data shall contain additional calibration data to the ones stored in EF Calibration. The old date and time value and the vehicle identification number shall be stored with each additional calibration data record, which shall be exactly the same as the old date and time value and the vehicle identification number stored with the corresponding calibration data in EF Calibration.
- in all tachograph cards:
  - EF VU\_Configuration shall contain the cardholder tachograph specific settings.

The vehicle unit shall ignore any authentication status found in EF Places\_Authentication or EF GNSS\_Places\_Authentication, when no vehicle position with the same timestamp is found in EF Places or EF GNSS\_Places.

The vehicle unit shall ignore the elementary file EF VU\_Configuration in all cards insofar as no specific rules have been provided with respect to the use of such elementary file. Those rules shall be set out through an amendment of Annex IC, which shall include the modification or deletion of this paragraph.

**▼B**

2 14 1

### 3.14 **Recording and storing on tachograph cards**

- 3.14.1 *Recording and storing in first generation tachograph cards*
  - (136) Provided first generation tachograph cards use has not been suppressed by a workshop, the recording equipment shall record and store data exactly in the same way as a first generation recording equipment would do.
  - (137) The recording equipment shall set the 'card session data' in the driver or workshop card right after the card insertion.
  - (138) The recording equipment shall update data stored on valid driver, workshop, company and/or control cards with all necessary data relevant to the period while the card is inserted and relevant to the card holder. Data stored on these cards are specified in Chapter 4.
  - (139) The recording equipment shall update driver activity and places data (as specified in 4.5.3.1.9 and 4.5.3.1.11), stored on valid driver and/or workshop cards, with activity and places data manually entered by the cardholder.
  - (140) All events and faults not defined for the first generation recording equipment shall not be stored on the first generation driver and workshop cards.

## **M3**

- (141) Tachograph cards data update shall be such that, when needed and taking into account card actual storage capacity, most recent data replace oldest data.
- (142) In the case of a writing error, the recording equipment shall try again, three times maximum, the same write command and then if still unsuccessful, declare the card faulty and non-valid.
- (143) Before releasing a driver or workshop card and after all relevant data have been stored on the card, the recording equipment shall reset the 'card session data'.

**▼M3**

3.14.2 *Recording and storing in second generation tachograph cards*

(144) Second generation tachograph cards shall contain 2 different card applications, the first of which shall be exactly the same as the TACHO application of first generation tachograph cards, and the second the 'TACHO\_G2' application, as specified in Chapter 4 and Appendix 2.

> The structure in the 'TACHO\_G2' application depends on the version. Version 2 cards contain additional Elementary Files to the ones of version 1 cards.

## **B**

**▼M3**

- (145) The recording equipment shall set the 'card session data' in the driver or workshop card right after the card insertion.
- (146) The recording equipment shall update data stored on the 2 card applications of valid driver, workshop, company and/or control cards with all necessary data relevant to the period while the card is inserted and relevant to the card holder. Data stored on these cards are specified in Chapter 4.
- (147) The recording equipment shall update driver activity places and positions data (as specified in 4.5.3.1.9, 4.5.3.1.11, 4.5.3.2.9 and 4.5.3.2.11), stored on valid driver and/or workshop cards, with activity and places data manually entered by the cardholder.

## **M3**

- (147a) On insertion of a driver or workshop card, the recording equipment shall store on the card the by-default load type of the vehicle.
- (147b) On insertion of a driver or workshop card, and after the manual entry procedure, the recording equipment shall check the last place where the daily work period begins or ends stored on the card. This place may be temporary, as specified in requirement 59. If this place is in a different country from the current one in which the vehicle is located, the recording equipment shall store on the card a border crossing record, with:

— the country that the driver left: not available,

— the country that the driver is entering: the current country in which the vehicle is located,

— the date and time when the driver has crossed the border: the card insertion time,

- the position of the driver when the border has been crossed: not available,
- the vehicle odometer value: not available.
- (148) Tachograph cards data update shall be such that, when needed and taking into account card actual storage capacity, most recent data replace oldest data.
- (149) In the case of a writing error, the recording equipment shall try again, three times maximum, the same write command and then if still unsuccessful, declare the card faulty and non-valid.
- (150) Before releasing a driver card and after all relevant data have been stored on the 2 card applications of the card, the recording equipment shall reset the 'card session data'.
- (150a) The vehicle unit shall ignore the elementary file EF VU\_Configuration in all cards insofar as no specific rules have been provided with respect to the use of such elementary file. Those rules shall be set out through an amendment of Annex IC, which shall include the modification or deletion of this paragraph.

## **B**

**▼M3**

## 3.15 **Displaying**

- (151) The display shall include at least 20 characters.
- (152) The minimum character size shall be 5 mm high and 3.5 mm wide.
- (153) The display shall support the characters specified in Appendix 1 Chapter 4 'Character sets'. The display may use simplified glyphs (e.g. accented characters may be displayed without accent, or lower case letters may be shown as upper case letters).
- (154) The display shall be provided with adequate non-dazzling lighting.
- (155) Indications shall be visible from outside the recording equipment.
- (156) The recording equipment shall be able to display:
  - default data,
  - data related to warnings,
  - data related to menu access,
  - other data requested by a user.

Additional information may be displayed by the recording equipment, provided that it is clearly distinguishable from information required above.

### **M3**

- (157) The display of the recording equipment shall use the pictograms or pictograms combinations listed in Appendix 3. Additional pictograms or pictograms combinations may also be provided by the display, if clearly distinguishable from the aforementioned pictograms or pictograms combinations.
- (158) The display shall always be ON when the vehicle is moving.
- (159) The recording equipment may include a manual or automatic feature to turn the display OFF when the vehicle is not moving.

Displaying format is specified in Appendix 5.

3.15.1 *Default display*

- (160) When no other information needs to be displayed, the recording equipment shall display, by default, the following:
  - the local time (as a result of UTC time + offset as set by the driver),
  - the mode of operation,
  - the current activity of the driver and the current activity of the co-driver,
  - information related to the driver:
  - if his current activity is DRIVING, his current continuous driving time and his current cumulative break time,
  - if his current activity is not DRIVING, the current duration of this activity (since it was selected) and his current cumulative break time.
- (161) Display of data related to each driver shall be clear, plain and unambiguous. In the case where the information related to the driver and the co-driver cannot be displayed at the same time, the recording equipment shall display by default the information related to the driver and shall allow the user to display the information related to the co-driver.
- (162) In the case where the display width does not allow displaying by default the mode of operation, the recording equipment shall briefly display the new mode of operation when it changes.
- (163) The recording equipment shall briefly display the card holder name at card insertion.
- (164) When an 'OUT OF SCOPE' or FERRY/TRAIN condition is opened, then the default display must show using the relevant pictogram that the particular condition is opened (it is acceptable that the driver's current activity may not be shown at the same time).

|     | 3.15.2 | Warning display |                                                                                                                                                                                                                                                                                    |
|-----|--------|-----------------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
|     |        | (165)           | The recording equipment shall display warning information<br>using primarily the pictograms of Appendix 3, completed<br>where needed by additional numerically coded information.<br>A literal description of the warning may also be added in<br>the driver's preferred language. |
|     | 3.15.3 | Menu access     |                                                                                                                                                                                                                                                                                    |
|     |        | (166)           | The recording equipment shall provide necessary commands<br>through an appropriate menu structure.                                                                                                                                                                                 |
|     | 3.15.4 | Other displays  |                                                                                                                                                                                                                                                                                    |
|     |        | (167)           | It shall be possible to display selectively on request:                                                                                                                                                                                                                            |
|     |        |                 | — the UTC date and time, and local time offset,                                                                                                                                                                                                                                    |
| ▼M3 |        |                 | — the content of any of the printouts listed in requirement<br>169 under the same formats as the printouts themselves,                                                                                                                                                             |
| ▼B  |        |                 | — the continuous driving time and cumulative break time<br>of the driver,                                                                                                                                                                                                          |
|     |        |                 | — the continuous driving time and cumulative break time<br>of the co-driver,                                                                                                                                                                                                       |
| ▼M3 |        |                 | — the accumulated driving time of the driver for the<br>previous and the current week,                                                                                                                                                                                             |
|     |        |                 | — the accumulated driving time of the co-driver for the<br>previous and the current week,                                                                                                                                                                                          |
| ▼B  |        |                 | optional:                                                                                                                                                                                                                                                                          |
|     |        |                 | — the current duration of co-driver activity (since it was<br>selected),                                                                                                                                                                                                           |
| ▼M3 |        |                 | — the accumulated driving time of the driver for the<br>current week,                                                                                                                                                                                                              |
|     |        |                 | — the accumulated driving time of the co-driver for the<br>current daily work period,                                                                                                                                                                                              |
|     |        |                 | — the accumulated driving time of the driver for the<br>current daily work period.                                                                                                                                                                                                 |
| ▼B  |        | (168)           | Printout content display shall be sequential, line by line. If<br>the display width is less than 24 characters the user shall<br>be provided with the complete information through an<br>appropriate mean (several lines, scrolling, …).                                           |

Printout lines devoted to hand-written information may be omitted for display.

(169) The recording equipment shall be able to print information from its data memory and/or from tachograph cards in accordance with the seven following printouts: — driver activities from card daily printout, — driver activities from Vehicle Unit daily printout, — events and faults from card printout, — events and faults from Vehicle Unit printout, — technical data printout, — over speeding printout.

> — tachograph card data history for a given VU (see chapter 3.12.16)

The detailed format and content of these printouts are specified in Appendix 4.

Additional data may be provided at the end of the printouts.

Additional printouts may also be provided by the recording equipment, if clearly distinguishable from the seven aforementioned printouts.

- (170) The 'driver activities from card daily printout' and 'Events and faults from card printout' shall be available only when a driver card or a workshop card is inserted in the recording equipment. The recording equipment shall update data stored on the relevant card before starting printing.
- (171) In order to produce the 'driver activities from card daily printout' or the 'events and faults from card printout', the recording equipment shall:
  - either automatically select the driver card or the workshop card if one only of these cards is inserted,
  - or provide a command to select the source card or select the card in the driver slot if two of these cards are inserted in the recording equipment.
- (172) The printer shall be able to print 24 characters per line.
- (173) The minimum character size shall be 2.1 mm high and 1.5 mm wide.
- (174) The printer shall support the characters specified in Appendix 1 Chapter 4 'Character sets'.
- (175) Printers shall be so designed as to produce these printouts with a degree of definition likely to avoid any ambiguity when they are read.

## **▼B**

3.16 **Printing**

- (176) Printouts shall retain their dimensions and recordings under normal conditions of humidity (10-90 %) and temperature.
- (177) The type approved paper used by the recording equipment shall bear the relevant type approval mark and an indication of the type(s) of recording equipment with which it may be used.
- (178) Printouts shall remain clearly legible and identifiable under normal conditions of storage, in terms of light intensity, humidity and temperature, for at least two years.
- (179) Printouts shall conform at least to the test specifications defined in Appendix 9.
- (180) It shall also be possible to add hand-written notes, such as the driver's signature, to these documents.
- (181) The recording equipment shall manage 'paper out' events while printing by, once paper has been re-loaded, restarting printing from printout beginning or by continuing printing and providing an unambiguous reference to previously printed part.

### 3.17 **Warnings**

- (182) The recording equipment shall warn the driver when detecting any event and/or fault.
- (183) Warning of a power supply interruption event may be delayed until the power supply is reconnected.
- (184) The recording equipment shall warn the driver 15 minutes before and at the time of exceeding the maximum allowed continuous driving time.
- (185) Warnings shall be visual. Audible warnings may also be provided in addition to visual warnings.
- (186) Visual warnings shall be clearly recognisable by the user, shall be situated in the driver's field of vision and shall be clearly legible both by day and by night.
- (187) Visual warnings may be built into the recording equipment and/or remote from the recording equipment.
- (188) In the latter case it shall bear a 'T' symbol.
- (189) Warnings shall have a duration of at least 30 seconds, unless acknowledged by the user by hitting one or more specific keys of the recording equipment. This first acknowledgement shall not erase warning cause display referred to in next paragraph.
- (190) Warning cause shall be displayed on the recording equipment and remain visible until acknowledged by the user using a specific key or command of the recording equipment.

| (191) | Additional warnings may be provided, as long as they do not confuse drivers in relation to previously defined ones. |
|-------|---------------------------------------------------------------------------------------------------------------------|
|-------|---------------------------------------------------------------------------------------------------------------------|

### 3.18 **Data downloading to external media**

- (192) The recording equipment shall be able to download on request data from its data memory or from a driver card to external storage media via the calibration/downloading connector. The recording equipment shall update data stored on the relevant card before starting downloading.
- (193) In addition and as an optional feature, the recording equipment may, in any mode of operation, download data through any other interface to a company authenticated through this channel. In such a case, company mode data access rights shall apply to this download.

**▼B**

**▼M3**

- (194) Downloading shall not alter or delete any stored data.
- (195) The calibration/downloading connector electrical interface is specified in Appendix 6.
- (196) Downloading protocols are specified in Appendix 7.

## **M3**

(196a) A transport undertaking which uses vehicles that are fitted with recording equipment complying with this Annex and fall within the scope of Regulation (EC) No 561/2006, shall ensure that all data are downloaded from the vehicle unit and driver cards.

> The maximum period within which the relevant data are downloaded shall not exceed:

- 90 days for data from vehicle unit;
- 28 days for data from the driver card.
- (196b) Transport undertakings shall keep the data downloaded from the vehicle unit and driver cards for at least twelve months following recording.

**▼B**

### 3.19 **Remote communication for targeted roadside checks**

- (197) When the ignition is on, the Vehicle Unit shall store every 60 seconds in the remote communication facility the most recent data necessary for the purpose of targeted roadside checks. Such data shall be encrypted and signed as specified in Appendix 11 and Appendix 14.
- (198) Data to be checked remotely shall be available to remote communication readers through wireless communication, as specified in Appendix 14.

- (199) Data necessary for the purpose of targeted roadside checks shall be related to:
  - the latest security breach attempt,
  - the longest power supply interruption,
  - sensor fault,
  - motion data error,
  - vehicle motion conflict,
  - driving without a valid card,
  - card insertion while driving,
  - time adjustment data,
  - calibration data including the dates of the two latest stored calibration records,
  - vehicle registration number,
  - speed recorded by the tachograph ,

- vehicle position,
- an indication if the driver may currently infringe the driving times.

### 3.20 **Data exchanges with additional external devices**

(200) The recording equipment shall also be equipped with an ITS interface in accordance with Appendix 13, allowing the data recorded or produced by either the tachograph or the tachograph cards to be used by an external facility.

> In operational mode, the driver consent shall be needed for the transmission of personal data through the ITS interface. Nevertheless, the driver consent shall not apply to tachograph or card data accessed in control, company or calibration mode. The data and functional access rights for these modes are specified in requirements 12 and 13.

> The following requirements shall apply to ITS data made available through that interface:

> — personal data shall only be available after the verifiable consent of the driver has been given, accepting that personal data can leave the vehicle network.

A set of selected existing data that can be available via the ITS interface, and the classification of the data as personal or not personal are specified in Appendix 13. Additional data may also be output in addition to the set of data provided in Appendix 13. The VU manufacturer shall classify those data as 'personal' or 'not personal', being the driver consent applicable to those data classified as 'personal',

- at any moment, the driver consent can be enabled or disabled through commands in the menu, provided the driver card is inserted,
- in any circumstances, the presence of the ITS interface shall not disturb or affect the correct functioning and the security of the vehicle unit.

Additional vehicle unit interfaces may co-exist, provided they fully comply with the requirements of Appendix 13 in terms of driver consent. The recording equipment shall have the capacity to communicate the driver consent status to other platforms in the vehicle network and to external devices.

For personal data injected in the vehicle network, which are further processed outside the vehicle network, it shall not be the responsibility of the tachograph manufacturer to have that personal data process compliant with the applicable Union legislation regarding data protection.

The ITS interface shall also allow for data entry during the manual entry procedure in accordance with requirement 61, for both the driver and the co-driver.

The ITS interface may also be used to enter additional information, in real time, such as:

- driver activity selection, in accordance with requirement 46,
- places in accordance with requirement 56,
- specific conditions, in accordance with requirement 62,
- load/unload operations, in accordance with requirement 62a.

This information may also be entered through other interfaces.

(201) The serial link interface as specified in Annex IB to Regulation (EEC) No 3821/85, as last amended, can continue to equip tachographs for back compatibility. The serial link is classified as a part of the vehicle network, in accordance with requirement 200.

### **B**

### 3.21 **Calibration**

- (202) The calibration function shall allow:
  - to automatically pair the motion sensor with the VU,
  - to automatically couple the external GNSS facility with the VU if applicable,
  - to digitally adapt the constant of the recording equipment (k) to the characteristic coefficient of the vehicle (w),

|     |       | — to adjust the current time within the validity period of<br>the inserted workshop card,                                                                                                                |
|-----|-------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
|     |       | — to adjust the current odometer value,                                                                                                                                                                  |
|     |       | — to update motion sensor identification data stored in<br>the data memory,                                                                                                                              |
|     |       | — to update, if applicable, external GNSS facility identi-<br>fication data stored in the data memory,                                                                                                   |
|     |       | — to update the types and identifiers of all the seals in<br>place,                                                                                                                                      |
| ▼M3 |       | — to update or confirm other parameters known to the<br>recording equipment: vehicle identification, w, l, tyre<br>size and speed limiting device setting if applicable, and<br>by-default load type,    |
|     |       | — to automatically store the country in which the cali-<br>bration has been performed, and the date time when<br>the position used to determine this country was<br>provided by the GNSS receiver.       |
| ▼B  | (203) | In addition, the calibration function shall allow to supress<br>the use of first generation tachograph cards in the<br>recording equipment, provided the conditions specified in<br>Appendix 15 are met. |
|     | (204) | Pairing the motion sensor to the VU shall consist, at least,<br>in:                                                                                                                                      |
|     |       | — updating motion sensor installation data held by the<br>motion sensor (as needed),                                                                                                                     |
|     |       | — copying from the motion sensor to the VU data<br>memory the necessary motion sensor identification<br>data.                                                                                            |
| ▼M3 | (205) | Coupling the external GNSS facility to the VU shall<br>consist, at least, in:                                                                                                                            |
|     |       | — updating external GNSS facility installation data held<br>by the external GNSS facility (as needed),                                                                                                   |
|     |       | — copying from the external GNSS facility to the VU<br>data memory the necessary external GNSS facility<br>identification data including the serial number of the<br>external GNSS facility.             |
| ▼B  | (206) | The calibration function shall be able to input necessary<br>data through the calibration/downloading connector in                                                                                       |

accordance with the calibration protocol defined in Appendix 8. The calibration function may also input

necessary data through other means.

## 3.22 **Roadside calibration checking**

- (207) The roadside calibration checking function shall allow reading the motion sensor serial number (possibly embedded in the adaptor) and the external GNSS facility serial number (when applicable), connected to the vehicle unit, at the time of the request.
- (208) This reading shall at least be possible on the vehicle unit display through commands in the menus.
- (209) The roadside calibration checking function shall also allow controlling the selection of the I/O mode of the calibration I/O signal line specified in Appendix 6, via the K-line interface. This shall be done through the ECUAdjustment-Session, as specified in Appendix 8, section 7 Control of Test Pulses — Input output control functional unit.
  - When the I/O mode of the calibration I/O signal line is active according to this requirement, the 'Driving without an appropriate card' warning (requirement 75) shall not be triggered by the vehicle unit.

3.23 **Time adjustment**

(210) The time adjustment function shall allow for automatically adjusting the current time. Two time sources are used in the recording equipment for time adjustment: 1) the internal VU clock, 2) the GNSS receiver.

### **M3**

- (211) The time setting of the VU internal clock shall be automatically re-adjusted at variable time intervals. The next automatic time re-adjustment shall be triggered between 72h and 168h after the previous one, and after the VU can access to GNSS time through a valid authenticated position message in accordance with Appendix 12. Nevertheless, the time adjustment shall never be bigger than the accumulated maximal time drift per day, as calculated by the VU manufacturer in accordance with requirement 41b. If the difference between internal VU clock time and GNSS receiver time is bigger than the accumulated maximum time drift per day, then the time adjustment shall bring the VU internal clock as close as possible to the GNSS receiver time. The time setting may only be done if the time provided by the GNSS receiver is obtained using authenticated position messages as set out in Appendix 12. The time reference for the automatic time setting of the VU internal clock shall be the time provided in the authenticated position message.
- (212) The time adjustment function shall also allow for triggered adjustment of the current time, in calibration mode.

Workshops may adjust time:

— either by writing a time value in the VU, using the WriteDataByIdentifier service in accordance with section 6.2 of Appendix 8,

### **B**

**▼M3**

— or by requesting an alignment of the VU clock to the time provided by the GNSS receiver. This may only be done if the time provided by the GNSS receiver is obtained using authenticated position messages. In this latter case, the RoutineControl service shall be used in accordance with section 8 of Appendix 8.

## 3.24 **Performance characteristics**

- (213) The Vehicle Unit shall be fully operational in the temperature range – 20 °C to 70 °C, the external GNSS facility in the temperature range – 20 °C to 70 °C, and the motion sensor in the temperature range – 40 °C to 135 °C. Data memory content shall be preserved at temperatures down to – 40 °C.
- (214) The tachograph shall be fully operational in the humidity range 10 % to 90 %.
- (215) The seals used in the smart tachograph shall withstand the same conditions than those applicable to the tachograph components to which they are affixed.
- (216) The recording equipment shall be protected against over-voltage, inversion of its power supply polarity, and short circuits.
- (217) Motion sensors shall either:
  - react to a magnetic field disturbing vehicle motion detection. In such circumstances, the vehicle unit will record and store a sensor fault (requirement 88) or,
  - have a sensing element that is protected from, or immune to, magnetic fields.
- (218) The recording equipment and the external GNSS facility shall conform to international regulation UN ECE R10 and shall be protected against electrostatic discharges and transients.

## 3.25 **Materials**

- (219) All the constituent parts of the recording equipment shall be made of materials of sufficient stability and mechanical strength and with stable electrical and magnetic characteristics.
- (220) For normal conditions of use, all the internal parts of the equipment shall be protected against damp and dust.
- (221) The Vehicle Unit and the external GNSS facility shall meet the protection grade IP 40 and the motion sensor shall meet the protection grade IP 64, as per standard IEC 60529:1989 including A1:1999 and A2:2013.
- (222) The recording equipment shall conform to applicable technical specifications related to ergonomic design.
- (223) The recording equipment shall be protected against accidental damage.

### **M3**

3.26 **Markings**

- (224) If the recording equipment displays the vehicle odometer value and speed, the following details shall appear on its display:
  - near the figure indicating the distance, the unit of measurement of distance, indicated by the abbreviation 'km',
  - near the figure showing the speed, the entry 'km/h'.

The recording equipment may also be switched to display the speed in miles per hour, in which case the unit of measurement of speed shall be shown by the abbreviation 'mph'. The recording equipment may also be switched to display the distance in miles, in which case the unit of measurement of distance shall be shown by the abbreviation 'mi'.

- (225) A descriptive plaque shall be affixed to each separate component of the recording equipment and shall show the following details:
  - name and address of the manufacturer,
  - manufacturer's part number and year of manufacture,
  - serial number,
  - type-approval mark.
- (226) When physical space is not sufficient to show all above mentioned details, the descriptive plaque shall show at least: the manufacturer's name or logo and the part number.

**▼M3**

## 3.27 **Monitoring border crossings**

- (226a) This function shall detect when the vehicle has crossed the border of a country, which country has been left and which country has been entered.
- (226b) The border crossing detection shall be based on the position measured by the recording equipment, and stored digital map in accordance with point 3.12.19.
- (226c) Border crossings related to the presence of the vehicle in a country for a shorter period than 120s shall not be recorded.

## 3.28 **Software update**

(226d) The vehicle unit shall incorporate a function for the implementation of software updates whenever such updates do not involve the availability of additional hardware resources beyond the resources set out in requirement 226f, and the type-approval authorities give their authorisation to the software updates based on the existing type-approved vehicle unit, in accordance with Article 12(5) of Regulation (EU) No 165/2014.

### **B**

- (226e) The software update function shall be designed for supporting the following functional features, whenever they are legally required:
  - modification of the functions referred to in point 2.2, except the software update function itself,
  - the addition of new functions directly related to the enforcement of Union legislation on road transport,
  - modification of the modes of operation in point 2.3,
  - modification of the file structure such as the addition of new data or the increase of the file size,
  - deployment of software patches to address software as well as security defects or reported attacks on the functions of the recording equipment.
- (226f) The vehicle unit shall provide free hardware resources of at least 35% for software and data needed for the implementation of requirement 226e and free hardware resources of at least 65% for the update of the digital map based on the hardware resources required for the NUTS 0 map version 2021.

## 4. CONSTRUCTION AND FUNCTIONAL REQUIREMENTS FOR TACHOGRAPH CARDS

## 4.1 **Visible data**

The front page shall contain:

- (227) the words 'Driver card' or 'Control card' or 'Workshop card' or 'Company card' printed in capital letters in the official language or languages of the Member State issuing the card, according to the type of the card.
- (228) the name of the Member State issuing the card (optional);
- (229) the distinguishing sign of the Member State issuing the card, printed in negative in a blue rectangle and encircled by 12 yellow stars. The distinguishing signs shall be as follows:

| B   | Belgium        | LV | Latvia          |
|-----|----------------|----|-----------------|
| BG  | Bulgaria       | L  | Luxembourg      |
| CZ  | Czech Republic | LT | Lithuania       |
| CY  | Cyprus         | M  | Malta           |
| DK  | Denmark        | NL | The Netherlands |
| D   | Germany        | A  | Austria         |
| EST | Estonia        | PL | Poland          |

| GR  | Greece  | P   | Portugal           |
|-----|---------|-----|--------------------|
|     |         | RO  | Romania            |
|     |         | SK  | Slovakia           |
|     |         | SLO | Slovenia           |
| E   | Spain   | FIN | Finland            |
| F   | France  | S   | Sweden             |
| HR  | Croatia |     |                    |
| H   | Hungary |     |                    |
| IRL | Ireland | UK  | The United Kingdom |
| I   | Italy   |     |                    |

| (230) | information specific to the card issued, numbered as follows: |  |  |  |  |
|-------|---------------------------------------------------------------|--|--|--|--|
|-------|---------------------------------------------------------------|--|--|--|--|

|      | Driver card                                                                                | Control Card                                       | Company or Workshop card                        |
|------|--------------------------------------------------------------------------------------------|----------------------------------------------------|-------------------------------------------------|
| 1.   | surname of the driver                                                                      | control body name                                  | company or workshop name                        |
| 2.   | first name(s) of the driver                                                                | surname of the controller<br>(if applicable)       | surname of card holder<br>(if applicable)       |
| 3.   | birth date of the driver                                                                   | first name(s) of the controller<br>(if applicable) | first name(s) of card holder<br>(if applicable) |
| 4.a  | card start of validity date                                                                |                                                    |                                                 |
| 4.b  | card expiry date                                                                           |                                                    |                                                 |
| 4.c  | the name of the issuing authority (may be printed on reverse page)                         |                                                    |                                                 |
| 4.d  | a different number from the one under heading 5, for administrative purposes<br>(optional) |                                                    |                                                 |
| 5. a | Driving licence number<br>(at the date of issue of the<br>driver card)                     | —                                                  | —                                               |
| 5. b | Card number                                                                                |                                                    |                                                 |
| 6.   | Photograph of the driver                                                                   | photograph of the controller<br>(optional)         | photograph of the fitter<br>(optional)-         |
| 7.   | Signature of the holder (optional)                                                         |                                                    |                                                 |
| 8.   | Normal place of residence,<br>or postal address of the<br>holder (optional).               | Postal address of control<br>body                  | postal address of company<br>or workshop        |

(231) dates shall be written using a 'dd/mm/yyyy' or 'dd.mm.yyyy' format (day, month, year).

The reverse page shall contain:

(232) an explanation of the numbered items which appear on the front page of the card;

- (233) with the specific written agreement of the holder, information which is not related to the administration of the card may also be added, such addition will not alter in any way the use of the model as a tachograph card.
- (234) Tachograph cards shall be printed with the following background predominant colours:
  - driver card: white,
  - control card: blue,
  - workshop card: red,
  - company card: yellow.
- (235) Tachograph cards shall bear at least the following features for protection of the card body against counterfeiting and tampering:
  - a security design background with fine guilloche patterns and rainbow printing,
  - in the area of the photograph, the security design background and the photograph shall overlap,
  - at least one two-coloured microprint line.

(1)

FRONT

DRIVER CARD
MEMBER STATE

1.

MS

2.

3.

4a.

4b.

4c.

6.

A
(4d.)

5a.

5b.

(7.)

G2
(8.)

B

B

(2)

CONTROL CARD
MEMBER STATE

1.

MS

(2.)

(3.)

43.

(40.)

(6.)

40.

(41.)

50.

(7.)

B.

G2

WORKSHOP CARD
MEMBER STATE

1.

MB

(2.)

(3.)

4a.

40.

40.

(40.)

50.

(7.)

B.

G2

COMPANY CARD
MEMBER STATE

1.

MS

(2.)

(3.)

43.

G2

40.

(4d.)

50.

(7.)

B.

40.

REVERSE

1. Surname 2. First name(s) 3. Birth date

4a. Date of start of validity of card

4b. Administrative expiry date of card

4c. Issuing authority

(4d.) No for national administrative purposes

5a. Driving license number 5b. Card number

6. Photograph

(7.) Signature (8.) Address

Please return to:

NAME OF AUTHORITY AND ADDRESS

1. Control Body (2.) Surname (3.) First name(s)

4a. Date of start of validity of card

4b. Administrative expiry date of card

4c. Issuing authority

(4d.) No for national administrative purposes

Sb. Card number

(6.) Photograph

(7.) Signature 8. Address

Please return to:

NAME OF AUTHORITY AND ADDRESS

1. Workshop Name (2.) Surname (3.) First name(s)

4a. Date of start of validity of card

40. Administrative expiry date of card
41. Issuing authority

(4d.) No for national administrative purposes

5b. Card number

(7.) Signature 8. Address

Please return to:

NAME OF AUTHORITY AND ADDRESS

이

1. Company Name (2.) surname (3.) First name(s)

4a. Date of start of validity of card

40. Administrative expiry date of card
41. Issuing authority

(4d.) No for natbnal administrative purposes

50. Card number

(7.) Signature 8. Address

Please return to:

NAME OF AUTHORITY AND ADDRESS

A

COMMUNITY MODEL TACHOGRAPH CARDS

## **(1) M1**

## **(2) M3**

(236) After consulting the Commission, Member States may add colours or markings, such as national symbols and security features, without prejudice to the other provisions of this Annex.

(237) Temporary cards referred to in Article 26.4 of Regulation (EU) No. 165/2014 shall comply with the provisions of this Annex.

### 4.2 **Security**

The system security aims at protecting integrity and authenticity of data exchanged between the cards and the recording equipment, protecting the integrity and authenticity of data downloaded from the cards, allowing certain write operations onto the cards to recording equipment only, decrypting certain data, ruling out any possibility of falsification of data stored in the cards, preventing tampering and detecting any attempt of that kind.

- (238) In order to achieve the system security, the tachograph cards shall meet the security requirements defined in Appendixes 10 and 11.
- (239) Tachograph cards shall be readable by other equipment such as personal computers.

### 4.3 **Standards**

- (240) Tachograph cards shall comply with the following standards:
  - ISO/IEC 7810 Identification cards Physical characteristics,
  - ISO/IEC 7816 Identification cards Integrated circuit cards:
    - Part 1: Physical characteristics,
    - Part 2: Dimensions and position of the contacts (ISO/IEC 7816-2:2007),
    - Part 3: Electrical interface and transmission protocols (ISO/IEC 7816-3:2006),
    - Part 4: Organisation, security and commands for interchange (ISO/IEC 7816-4:2013 + Cor 1:2014),
    - Part 6: Interindustry data elements for interchange (ISO/IEC 7816-6:2004 + Cor 1:2006),
    - Part 8: Commands for security operations (ISO/IEC 7816-8:2004).
  - Tachograph cards shall be tested in accordance to ISO/IEC 10373-3:2010 Identification cards — Test methods — Part 3: Integrated circuit cards with contacts and related interface devices.

### 4.4 **Environmental and electrical specifications**

- (241) Tachograph cards shall be capable of operating correctly in all the climatic conditions normally encountered in Community territory and at least in the temperature range – 25 °C to + 70 °C with occasional peaks of up to + 85 °C, 'occasional' meaning not more than 4 hours each time and not over 100 times during the life time of the card.
- (242) Tachograph cards shall be capable of operating correctly in the humidity range 10 % to 90 %.
- (243) Tachograph cards shall be capable of operating correctly for a five-year period if used within the environmental and electrical specifications.
- (244) During operation, tachograph cards shall conform to ECE R10, related to electromagnetic compatibility, and shall be protected against electrostatic discharges.

## 4.5 **Data storage**

For the purpose of this paragraph,

- times are recorded with a resolution of one minute, unless otherwise specified,
- odometer values are recorded with a resolution of one kilometre,
- speeds are recorded with a resolution of 1 km/h,
- positions (latitudes and longitudes) are recorded in degrees and minutes with a resolution of 1/10 of minute.

The tachograph cards functions, commands and logical structures, fulfilling data storage requirements are specified in Appendix 2.

If not otherwise specified, data storage on tachograph cards shall be organized in such a way, that new data replaces stored oldest data in case the foreseen memory size for the particular records is exhausted.

(245) This paragraph specifies minimum storage capacity for the various application data files. Tachograph cards shall be able to indicate to the recording equipment the actual storage capacity of these data files.

**▼M3**

(246) Any additional data may be stored on tachograph cards, provided that the storage of those data complies with the applicable legislation regarding data protection.

**▼B**

**▼M3**

**▼B**

- (247) Each Master File (MF) of any tachograph card shall contain up to five Elementary Files (EF) for card management, application and chip identifications, and two Dedicated Files (DF):
  - DF Tachograph, which contains the application accessible to first generation vehicle units, which is also present in first generation tachograph cards,
  - DF Tachograph\_G2, which contains the application only accessible to second generation vehicle units, which is only present in second generation tachograph cards.
  - Note: version 2 of second generation cards contains additional Elementary Files in DF Tachograph\_G2.
  - The full details of the tachograph cards structure are specified in Appendix 2.
- 4.5.1 *Elementary files for identification and card management*
- 4.5.2 *IC card identification*
  - (248) Tachograph cards shall be able to store the following smart card identification data:
    - clock stop,

— card serial number (including manufacturing references),

- card type approval number,
- card personaliser identification (ID),
- embedder ID,
- IC identifier.
- 4.5.2.1 C h i p i d e n t i f i c a t i o n
  - (249) Tachograph cards shall be able to store the following Integrated Circuit (IC) identification data:
    - IC serial number,
    - IC manufacturing references.
- 4.5.2.2 D I R ( o n l y p r e s e n t i n s e c o n d g e n e r a t i o n t a c h o g r a p h c a r d s )
  - (250) Tachograph cards shall be able to store the application identification data objects specified in Appendix 2.
- 4.5.2.3 A T R i n f o r m a t i o n ( c o n d i t i o n a l , o n l y p r e s e n t i n s e c o n d g e n e r a t i o n t a c h o g r a p h c a r d s )
  - (251) Tachograph cards shall be able to store the following extended length information data object:
    - in the case the tachograph card supports extended length fields, the extended length information data object specified in Appendix 2.
- 4.5.2.4 E x t e n d e d l e n g t h i n f o r m a t i o n ( c o n d i t i o n a l , o n l y p r e s e n t i n s e c o n d g e n e r a t i o n t a c h o g r a p h c a r d s )
  - (252) Tachograph cards shall be able to store the following extended length information data objects:
    - in the case the tachograph card supports extended length fields, the extended length information data objects specified in Appendix 2.
- 4.5.3 *Driver card*
- 4.5.3.1 T a c h o g r a p h a p p l i c a t i o n ( a c c e s s i b l e t o f i r s t a n d s e c o n d g e n e r a t i o n v e h i c l e u n i t s )
- 4.5.3.1.1 Application identification
  - (253) The driver card shall be able to store the following application identification data:
    - tachograph application identification,
    - type of tachograph card identification.

### 4.5.3.1.2 Key and certificates

- (254) The driver card shall be able to store a number of cryptographic keys and certificates, as specified in Appendix 11 part A.
- 4.5.3.1.3 Card identification
  - (255) The driver card shall be able to store the following card identification data:
    - card number,
    - issuing Member State, issuing authority name, issue date,
    - card beginning of validity date, card expiry date.

### 4.5.3.1.4 Card holder identification

- (256) The driver card shall be able to store the following card holder identification data:
  - surname of the holder,
  - first name(s) of the holder,
  - date of birth,
  - preferred language.

### 4.5.3.1.5 Card download

- (257) The driver card shall be able to store the following data related to card download:
  - date and time of last card download (for other purposes than control).
- (258) The driver card shall be able to hold one such record.
- 4.5.3.1.6 Driving licence information
  - (259) The driver card shall be able to store the following driving licence data:
    - issuing Member State, issuing authority name,
    - driving licence number (at the date of the issue of the card).

## 4.5.3.1.7 Events data

For the purpose of this subparagraph, time shall be stored with a resolution of 1 second.

- (260) The driver card shall be able to store data related to the following events detected by the recording equipment while the card was inserted:
  - Time overlap (where this card is the cause of the event),
  - Card insertion while driving (where this card is the subject of the event),

— Last card session not correctly closed (where this card is the subject of the event),

- Power supply interruption,
- Motion data error,
- Security breach attempts.
- (261) The driver card shall be able to store the following data for these events:
  - Event code,
  - Date and time of beginning of the event (or of card insertion if the event was on-going at that time),
  - Date and time of end of the event (or of card withdrawal if the event was on-going at that time),
  - VRN and registering Member State of vehicle in which the event happened.
  - Note: For the 'Time overlap' event:
  - Date and time of beginning of the event shall correspond to the date and time of the card withdrawal from the previous vehicle,
  - Date and time of end of the event shall correspond to the date and time of card insertion in current vehicle,
  - Vehicle data shall correspond to the current vehicle raising the event.

Note: For the 'Last card session not correctly closed' event:

- date and time of beginning of event shall correspond to the card insertion date and time of the session not correctly closed,
- date and time of end of event shall correspond to the card insertion date and time of the session during which the event was detected (current session),
- Vehicle data shall correspond to the vehicle in which the session was not correctly closed.
- (262) The driver card shall be able to store data for the six most recent events of each type (i.e. 36 events).

| ▼B  | 4.5.3.1.8 | Faults data                                                                                 |                                                                                                                                                                         |
|-----|-----------|---------------------------------------------------------------------------------------------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
|     |           | For the purpose of this subparagraph, time shall be recorded with a resolution of 1 second. |                                                                                                                                                                         |
|     |           | (263)                                                                                       | The driver card shall be able to store data related to the following faults detected by the recording equipment while the card was inserted:                            |
| ▼M1 |           |                                                                                             | — Card fault (where this card is the subject of the fault),                                                                                                             |
| ▼B  |           |                                                                                             | — Recording equipment fault.                                                                                                                                            |
|     |           | (264)                                                                                       | The driver card shall be able to store the following data for these faults:                                                                                             |
|     |           |                                                                                             | — Fault code,                                                                                                                                                           |
|     |           |                                                                                             | — Date and time of beginning of the fault (or of card insertion if the fault was on-going at that time),                                                                |
|     |           |                                                                                             | — Date and time of end of the fault (or of card with-drawal if the fault was on-going at that time),                                                                    |
|     |           |                                                                                             | — VRN and registering Member State of vehicle in which the fault happened.                                                                                              |
|     |           | (265)                                                                                       | The driver card shall be able to store data for the twelve most recent faults of each type (i.e. 24 faults).                                                            |
|     | 4.5.3.1.9 |                                                                                             | Driver activity data                                                                                                                                                    |
|     |           | (266)                                                                                       | The driver card shall be able to store, for each calendar day where the card has been used or for which the driver has entered activities manually, the following data: |
|     |           |                                                                                             | — the date,                                                                                                                                                             |
|     |           |                                                                                             | — a daily presence counter (increased by one for each of these calendar days),                                                                                          |
|     |           |                                                                                             | — the total distance travelled by the driver during this day,                                                                                                           |
|     |           |                                                                                             | — a driver status at 00:00,                                                                                                                                             |
|     |           |                                                                                             | — whenever the driver has changed of activity, and/or has changed of driving status, and/or has inserted or withdrawn his card:                                         |
|     |           |                                                                                             | — the driving status (CREW, SINGLE),                                                                                                                                    |
|     |           |                                                                                             | — the slot (DRIVER, CO-DRIVER),                                                                                                                                         |
|     |           |                                                                                             | — the card status (INSERTED, NOT INSERTED),                                                                                                                             |
|     |           |                                                                                             | — the activity (DRIVING, AVAILABILITY, WORK, BREAK/REST),                                                                                                               |

## the time of the change.

- (267) The driver card memory shall be able to hold driver activity data for at least 28 days (the average activity of a driver is defined as 93 activity changes per day).
- (268) The data listed under requirements 261, 264 and 266 shall be stored in a way allowing the retrieval of activities in the order of their occurrence, even in case of a time overlap situation.
- 4.5.3.1.10 Vehicles used data
  - (269) The driver card shall be able to store, for each calendar day where the card has been used, and for each period of use of a given vehicle that day (a period of use includes all consecutive insertion / withdrawal cycle of the card in the vehicle, as seen from the card point of view), the following data:
    - date and time of first use of the vehicle (i.e. first card insertion for this period of use of the vehicle, or 00h00 if the period of use is on-going at that time),
    - vehicle odometer value at that time,
    - date and time of last use of the vehicle, (i.e. last card withdrawal for this period of use of the vehicle, or 23h59 if the period of use is on-going at that time),
    - vehicle odometer value at that time,
    - VRN and registering Member State of the vehicle.
  - (270) The driver card shall be able to store at least 84 such records.
- 4.5.3.1.11 Places where daily work periods start and/or end
  - (271) The driver card shall be able to store the following data related to places where daily work periods begin and/or end, entered by the driver:
    - the date and time of the entry (or the date/time related to the entry if the entry is made during the manual entry procedure),
    - the type of entry (begin or end, condition of entry),
    - the country and region entered,
    - the vehicle odometer value.
  - (272) The driver card memory shall be able to hold at least 42 pairs of such records.

### 4.5.3.1.12 Card session data

(273) The driver card shall be able to store data related to the vehicle which opened its current session:

— date and time the session was opened (i.e. card insertion) with a resolution of one second,

- VRN and registering Member State.
- 4.5.3.1.13 Control activity data
  - (274) The driver card shall be able to store the following data related to control activities:
    - date and time of the control,
    - control card number and card issuing Member State,
    - type of the control (displaying and/or printing and/or VU downloading and/or card downloading (see note)),
    - Period downloaded, in case of downloading,
    - VRN and registering Member State of the vehicle in which the control happened.

Note: card downloading will only be recorded if performed through a recording equipment.

- (275) The driver card shall be able to hold one such record.
- 4.5.3.1.14 Specific conditions data
  - (276) The driver card shall be able to store the following data related to specific conditions entered while the card was inserted (whatever the slot):
    - Date and time of the entry,
    - Type of specific condition.
  - (277) The driver card shall be able to store at least 56 such records.

## **M3**

**▼B**

4.5.3.2 T a c h o g r a p h g e n e r a t i o n 2 a p p l i c a t i o n ( n o t a c c e s s i b l e t o f i r s t g e n e r a t i o n v e h i c l e u n i t s , a c c e s s i b l e t o v e r s i o n 1 a n d v e r s i o n 2 o f s e c o n d g e n e r a t i o n v e h i c l e u n i t s )

### 4.5.3.2.1 Application identification

- (278) The driver card shall be able to store the following application identification data:
  - tachograph application identification,
  - type of tachograph card identification.

## **M3**

- 4.5.3.2.1.1 Additional application identification (not accessed by version 1 of second generation vehicle units)
  - (278a) The driver card shall be able to store additional application identification data only applicable for version 2.

### 4.5.3.2.2 Keys and certificates

- (279) The driver card shall be able to store a number of cryptographic keys and certificates, as specified in Appendix 11 part B.
- 4.5.3.2.3 Card identification
  - (280) The driver card shall be able to store the following card identification data:
    - card number,
    - issuing Member State, issuing authority name, issue date,
    - card beginning of validity date, card expiry date.
- 4.5.3.2.4 Card holder identification
  - (281) The driver card shall be able to store the following card holder identification data:
    - surname of the holder,
    - first name(s) of the holder,
    - date of birth,
    - preferred language.

### 4.5.3.2.5 Card download

- (282) The driver card shall be able to store the following data related to card download:
  - date and time of last card download (for other purposes than control).
- (283) The driver card shall be able to hold one such record.

### 4.5.3.2.6 Driving licence information

- (284) The driver card shall be able to store the following driving licence data:
  - issuing Member State, issuing authority name,

— driving licence number (at the date of the issue of the card).

4.5.3.2.7 Events data

For the purpose of this subparagraph, time shall be stored with a resolution of 1 second.

- (285) The driver card shall be able to store data related to the following events detected by the recording equipment while the card was inserted:
  - Time overlap (where this card is the cause of the event),

- Card insertion while driving (where this card is the subject of the event),
- Last card session not correctly closed (where this card is the subject of the event),
- Power supply interruption,
- Communication error with the remote communication facility,
- Absence of position information from GNSS receiver event,
- Communication error with the external GNSS facility
- Motion data error,
- Vehicle motion conflict,
- Security breach attempts,
- Time conflict.
- (286) The driver card shall be able to store the following data for these events:
  - Event code,
  - Date and time of beginning of the event (or of card insertion if the event was on-going at that time),
  - Date and time of end of the event (or of card withdrawal if the event was on-going at that time),
  - VRN and registering Member State of vehicle in which the event happened.

Note: For the 'Time overlap' event:

- Date and time of beginning of the event shall correspond to the date and time of the card withdrawal from the previous vehicle,
- Date and time of end of the event shall correspond to the date and time of card insertion in current vehicle,
- Vehicle data shall correspond to the current vehicle raising the event.

Note: For the 'Last card session not correctly closed' event:

- date and time of beginning of event shall correspond to the card insertion date and time of the session not correctly closed,
- date and time of end of event shall correspond to the card insertion date and time of the session during which the event was detected (current session),

| ▼B  |           |       | — Vehicle data shall correspond to the vehicle in which<br>the session was not correctly closed.                                                                              |
|-----|-----------|-------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| ▼M3 |           | (287) | The driver card shall be able to store data for the 12 most<br>recent events of each type (i.e. 132 events).                                                                  |
| ▼B  | 4.5.3.2.8 |       | Faults data<br>For the purpose of this subparagraph, time shall be recorded with a<br>resolution of 1 second.                                                                 |
|     |           | (288) | The driver card shall be able to store data related to the<br>following faults detected by the recording equipment while<br>the card was inserted:                            |
| ▼M1 |           |       | — Card fault (where this card is the subject of the fault),                                                                                                                   |
| ▼B  |           |       | — Recording equipment fault.                                                                                                                                                  |
|     |           | (289) | The driver card shall be able to store the following data for<br>these faults:                                                                                                |
|     |           |       | — Fault code,                                                                                                                                                                 |
|     |           |       | — Date and time of beginning of the fault (or of card<br>insertion if the fault was on-going at that time),                                                                   |
|     |           |       | — Date and time of end of the fault (or of card with-<br>drawal if the fault was on-going at that time),                                                                      |
|     |           |       | — VRN and registering Member State of vehicle in which<br>the fault happened.                                                                                                 |
| ▼M3 |           | (290) | The driver card shall be able to store data for the 24 most<br>recent faults of each type (i.e. 48 faults).                                                                   |
| ▼B  | 4.5.3.2.9 |       | Driver activity data                                                                                                                                                          |
|     |           | (291) | The driver card shall be able to store, for each calendar<br>day where the card has been used or for which the driver<br>has entered activities manually, the following data: |
|     |           |       | — the date,                                                                                                                                                                   |
|     |           |       | — a daily presence counter (increased by one for each of<br>these calendar days),                                                                                             |
|     |           |       | — the total distance travelled by the driver during this<br>day,                                                                                                              |
|     |           |       | — a driver status at 00:00,                                                                                                                                                   |
|     |           |       | — whenever the driver has changed of activity, and/or has<br>changed of driving status, and/or has inserted or<br>withdrawn his card:                                         |
|     |           |       | — the driving status (CREW, SINGLE)                                                                                                                                           |

— the card status (INSERTED, NOT INSERTED),

— the activity (DRIVING, AVAILABILITY, WORK, BREAK/REST).

— the time of the change,

(292) The driver card memory shall be able to hold driver activity data for 56 days (the average activity of a driver is defined for this requirement as 117 activity changes per day).

(293) The data listed under requirements 286, 289 and 291 shall be stored in a way allowing the retrieval of activities in the order of their occurrence, even in case of a time overlap situation.

4.5.3.2.10 Vehicles used data

(294) The driver card shall be able to store, for each calendar day where the card has been used, and for each period of use of a given vehicle that day (a period of use includes all consecutive insertion / withdrawal cycle of the card in the vehicle, as seen from the card point of view), the following data:

> — date and time of first use of the vehicle (i.e. first card insertion for this period of use of the vehicle, or 00h00 if the period of use is on-going at that time),

— vehicle odometer value at that first use time,

- date and time of last use of the vehicle, (i.e. last card withdrawal for this period of use of the vehicle, or 23h59 if the period of use is on-going at that time),
- vehicle odometer value at that last use time,
- VRN and registering Member State of the vehicle,
- VIN of the vehicle.

## **M3**

(295) The driver card shall be able to store 200 such records.

### **B**

- 4.5.3.2.11 Places and positions where daily work periods start and/or end
  - (296) The driver card shall be able to store the following data related to places where daily work periods begin and/or end, entered by the driver:
    - the date and time of the entry (or the date/time related to the entry if the entry is made during the manual entry procedure),

### **B**

**▼B**

- the type of entry (begin or end, condition of entry),
- the country and region entered,
- the vehicle odometer value,
- the vehicle position,
- the GNSS accuracy, date and time when the position was determined.

(297) The driver card memory shall be able to hold 112 such records.

## **B**

4.5.3.2.12 Card session data

- (298) The driver card shall be able to store data related to the vehicle which opened its current session:
  - date and time the session was opened (i.e. card insertion) with a resolution of one second,
  - VRN and registering Member State.
- 4.5.3.2.13 Control activity data
  - (299) The driver card shall be able to store the following data related to control activities:
    - date and time of the control,
    - control card number and card issuing Member State,
    - type of the control (displaying and/or printing and/or VU downloading and/or card downloading (see note)),
    - Period downloaded, in case of downloading,
    - VRN and registering Member State of the vehicle in which the control happened.

Note: security requirements imply that card downloading will only be recorded if performed through a recording equipment.

- (300) The driver card shall be able to hold one such record.
- 4.5.3.2.14 Specific conditions data
  - (301) The driver card shall be able to store the following data related to specific conditions entered while the card was inserted (whatever the slot):
    - Date and time of the entry,
    - Type of specific condition.

| ▼M3 | (302)      | The driver card shall be able to store 112 such records.                                                                                                                                                                                                                                                                                                                                                      |
|-----|------------|---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| ▼B  | 4.5.3.2.15 | Vehicle units used data                                                                                                                                                                                                                                                                                                                                                                                       |
|     | (303)      | The driver card shall be able to store the following data related to the different vehicle units in which the card was used:<br>— the date and time of the beginning of the period of use of the vehicle unit (i.e. first card insertion in the vehicle unit for the period),<br>— the manufacturer of the vehicle unit,<br>— the vehicle unit type,<br>— the vehicle unit software version number.           |
| ▼M3 | (304)      | The driver card shall be able to store 200 such records.                                                                                                                                                                                                                                                                                                                                                      |
| ▼M1 | 4.5.3.2.16 | Three hours accumulated driving places data                                                                                                                                                                                                                                                                                                                                                                   |
|     | (305)      | The driver card shall be able to store the following data related to the position of the vehicle where the accumulated driving time reaches a multiple of three hours:<br>— the date and time when the accumulated driving time reaches a multiple of three hours,<br>— the position of the vehicle,<br>— the GNSS accuracy, date and time when the position was determined,<br>— the vehicle odometer value. |
| ▼M3 | (306)      | The driver card shall be able to store 336 such records.                                                                                                                                                                                                                                                                                                                                                      |
|     | 4.5.3.2.17 | Authentication status for positions related to places where daily work periods start and/or end (not accessed by version 1 of second generation vehicle units)                                                                                                                                                                                                                                                |
|     | (306a)     | The driver card shall be able to store additional data related to places where daily work periods begin and/or end, entered by the driver in accordance with point 4.5.3.2.11:<br>— the date and time of the entry which shall be exactly                                                                                                                                                                     |

- the same date and time as the one stored in EF Places under the DF Tachograph\_G2,
- a flag indicating whether the position has been authenticated.
- (306b) The driver card memory shall be able to hold 112 such records.

- 4.5.3.2.18 Authentication status for positions where three hours accumulated driving time are reached (not accessed by version 1 of second generation vehicle units)
  - (306c) The driver card shall be able to store additional data related to the position of the vehicle where the accumulated driving time reaches a multiple of three hours in accordance with point 4.5.3.2.16:
    - the date and time when the accumulated driving time reaches a multiple of three hours, which shall be exactly the same date and time as the one stored in EF GNSS\_Places under the DF Tachograph\_G2,
    - a flag indicating whether the position has been authenticated.
  - (306d) The driver card shall be able to store 336 such records.
- 4.5.3.2.19 Border crossings (not accessed by version 1 of second generation vehicle units)
  - (306e) The driver card shall be able to store the following data related to border crossings either upon card insertion in accordance with requirement 147b or with the card already inserted:
    - the country that the vehicle is leaving,
    - the country that the vehicle is entering,
    - the date and time when the vehicle has crossed the border,
    - the position of the vehicle when the border was crossed,
    - the GNSS accuracy,
    - a flag indicating whether the position has been authenticated,
    - the vehicle odometer value.
  - (306f) The driver card memory shall be able to store 1120 such records.
- 4.5.3.2.20 Load/unload operations (not accessed by version 1 of second generation vehicle units)
  - (306g) The driver card shall be able to store the following data related to load/unload operations:
    - operation type (load, unload or simultaneous load/ unload),
    - the date and time of the load/unload operation,
    - the position of the vehicle,
    - the GNSS accuracy, date and time when the position was determined,

- a flag indicating whether the position has been authenticated,
- the vehicle odometer value.
- (306h) The driver card shall be able to store 1624 load/unload operations.
- 4.5.3.2.21 Load type entries (not accessed by version 1 of second generation vehicle units)
  - (306i) The driver card shall be able to store the following data related to load type automatically entered by the VU at each card insertion:
    - the load type entered (goods or passengers),
    - the date and time of the entry.
  - (306j) The driver card shall be able to store 336 such records.
- 4.5.3.2.22 VU configurations (not accessed by version 1 of second generation vehicle units)
  - (306k) The driver card shall be able to store the cardholder tachograph specific settings.
  - (306l) The driver card storage capacity for cardholder tachograph specific settings shall be 3072 bytes.

| 4.5.4 | Workshop card |
|-------|---------------|
|-------|---------------|

- 4.5.4.1 T a c h o g r a p h a p p l i c a t i o n ( a c c e s s i b l e t o f i r s t a n d s e c o n d g e n e r a t i o n v e h i c l e u n i t s )
- 4.5.4.1.1 Application identification
  - (307) The workshop card shall be able to store the following application identification data:
    - tachograph application identification,
    - type of tachograph card identification.
- 4.5.4.1.2 Keys and certificates
  - (308) The workshop card shall be able to store a number of cryptographic keys and certificates, as specified in Appendix 11 part A.
  - (309) The workshop card shall be able to store a Personal Identification Number (PIN code).
- 4.5.4.1.3 Card identification
  - (310) The workshop card shall be able to store the following card identification data:
    - card number,
    - issuing Member State, issuing authority name, issue date,
    - card beginning of validity date, card expiry date.

### 4.5.4.1.4 Card holder identification

- (311) The workshop card shall be able to store the following card holder identification data:
  - workshop name,
  - workshop address,
  - surname of the holder,
  - first name(s) of the holder,
  - preferred language.

### 4.5.4.1.5 Card download

- (312) The workshop card shall be able to store a card download data record in the same manner as a driver card.
- 4.5.4.1.6 Calibration and time adjustment data
  - (313) The workshop card shall be able to hold records of calibrations and/or time adjustments performed while the card is inserted in a recording equipment.
  - (314) Each calibration record shall be able to hold the following data:
    - Purpose of calibration (activation, first installation, installation, periodic inspection,),
    - Vehicle identification,
    - Parameters updated or confirmed (w, k, l, tyre size, speed limiting device setting, odometer (new and old values), date and time (new and old values)),
    - Recording equipment identification (VU part number, VU serial number, motion sensor serial number).
  - (315) The workshop card shall be able to store at least 88 such records.
  - (316) The workshop card shall hold a counter indicating the total number of calibrations performed with the card.
  - (317) The workshop card shall hold a counter indicating the number of calibrations performed since its last download.

### 4.5.4.1.7 Events and faults data

- (318) The workshop card shall be able to store events and faults data records in the same manner as a driver card.
- (319) The workshop card shall be able to store data for the three most recent events of each type (i.e. 18 events) and the six most recent faults of each type (i.e. 12 faults).

| ▼B  |                                                                                                                                                                          |                                                                                                                                      |
|-----|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------|--------------------------------------------------------------------------------------------------------------------------------------|
|     | 4.5.4.1.8                                                                                                                                                                | Driver activity data                                                                                                                 |
|     |                                                                                                                                                                          | (320) The workshop card shall be able to store driver activity data in the same manner as a driver card.                             |
|     |                                                                                                                                                                          | (321) The workshop card shall be able to hold driver activity data for at least 1 day of average driver activity.                    |
|     | 4.5.4.1.9                                                                                                                                                                | Vehicles used data                                                                                                                   |
|     |                                                                                                                                                                          | (322) The workshop card shall be able to store vehicles used data records in the same manner as a driver card.                       |
|     |                                                                                                                                                                          | (323) The workshop card shall be able to store at least 4 such records.                                                              |
|     | 4.5.4.1.10                                                                                                                                                               | Daily work periods start and/or end data                                                                                             |
|     |                                                                                                                                                                          | (324) The workshop card shall be able to store daily works period start and/or end data records in the same manner as a driver card. |
|     |                                                                                                                                                                          | (325) The workshop card shall be able to hold at least 3 pairs of such records.                                                      |
|     | 4.5.4.1.11                                                                                                                                                               | Card session data                                                                                                                    |
|     |                                                                                                                                                                          | (326) The workshop card shall be able to store a card session data record in the same manner as a driver card.                       |
|     | 4.5.4.1.12                                                                                                                                                               | Control activity data                                                                                                                |
|     |                                                                                                                                                                          | (327) The workshop card shall be able to store a control activity data record in the same manner as a driver card.                   |
|     | 4.5.4.1.13                                                                                                                                                               | Specific conditions data                                                                                                             |
|     |                                                                                                                                                                          | (328) The workshop card shall be able to store data relevant to specific conditions in the same manner as the driver card.           |
|     |                                                                                                                                                                          | (329) The workshop card shall be able to store at least 2 such records.                                                              |
| ▼M3 |                                                                                                                                                                          |                                                                                                                                      |
|     | 4.5.4.2 Tachograph Generation 2 application (not accessible to first generation vehicle units, accessible to version 1 and version 2 of second generation vehicle units) |                                                                                                                                      |
| ▼B  | 4.5.4.2.1                                                                                                                                                                | Application identification                                                                                                           |
|     |                                                                                                                                                                          | (330) The workshop card shall be able to store the following application identification data:                                        |
|     |                                                                                                                                                                          | — tachograph application identification,                                                                                             |

- 4.5.4.2.1.1 Additional application identification (not accessed by version 1 of second generation vehicle units)
  - (330a) The workshop card shall be able to store additional application identification data only applicable for version 2.

### 4.5.4.2.2 Keys and certificates

- (331) The workshop card shall be able to store a number of cryptographic keys and certificates, as specified in Appendix 11 part B.
- (332) The workshop card shall be able to store a Personal Identification Number (PIN code).

### 4.5.4.2.3 Card identification

- (333) The workshop card shall be able to store the following card identification data:
  - card number,
  - issuing Member State, issuing authority name, issue date,
  - card beginning of validity date, card expiry date.
- 4.5.4.2.4 Card holder identification
  - (334) The workshop card shall be able to store the following card holder identification data:
    - workshop name,
    - workshop address,
    - surname of the holder,
    - first name(s) of the holder,
    - preferred language.

## 4.5.4.2.5 Card download

- (335) The workshop card shall be able to store a card download data record in the same manner as a driver card.
- 4.5.4.2.6 Calibration and time adjustment data
  - (336) The workshop card shall be able to hold records of calibrations and/or time adjustments performed while the card is inserted in a recording equipment.
  - (337) Each calibration record shall be able to hold the following data:
    - purpose of calibration (activation, first installation, installation, periodic inspection,),
    - vehicle identification,
    - parameters updated or confirmed (w, k, l, tyre size, speed limiting device setting, odometer (new and old values), date and time (new and old values),
    - recording equipment identification (VU part number, VU serial number, motion sensor serial number, remote communication facility serial number and external GNSS facility serial number, if applicable),

| <b>▼B</b>  |            |        | — seal type and identifier of all seals in place,                                                                                                                                                                                                                                                                                                                                                                                                              |
|------------|------------|--------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
|            |            |        | — ability of the VU to use first generation tachograph<br>cards (enabled or not).                                                                                                                                                                                                                                                                                                                                                                              |
| <b>▼M3</b> |            | (338)  | The workshop card shall be able to store 255 such records.                                                                                                                                                                                                                                                                                                                                                                                                     |
| <b>▼B</b>  |            | (339)  | The workshop card shall hold a counter indicating the total<br>number of calibrations performed with the card.                                                                                                                                                                                                                                                                                                                                                 |
|            |            | (340)  | The workshop card shall hold a counter indicating the<br>number of calibrations performed since its last download.                                                                                                                                                                                                                                                                                                                                             |
|            | 4.5.4.2.7  |        | Events and faults data                                                                                                                                                                                                                                                                                                                                                                                                                                         |
|            |            | (341)  | The workshop card shall be able to store events and faults<br>data records in the same manner as a driver card.                                                                                                                                                                                                                                                                                                                                                |
|            |            | (342)  | The workshop card shall be able to store data for the three<br>most recent events of each type (i.e. 33 events) and the six<br>most recent faults of each type (i.e. 12 faults).                                                                                                                                                                                                                                                                               |
|            | 4.5.4.2.8  |        | Driver activity data                                                                                                                                                                                                                                                                                                                                                                                                                                           |
|            |            | (343)  | The workshop card shall be able to store driver activity<br>data in the same manner as a driver card.                                                                                                                                                                                                                                                                                                                                                          |
| <b>▼M3</b> |            | (344)  | The workshop card shall be able to hold driver activity<br>data for 1 day containing 240 activity changes.                                                                                                                                                                                                                                                                                                                                                     |
| <b>▼B</b>  | 4.5.4.2.9  |        | Vehicles used data                                                                                                                                                                                                                                                                                                                                                                                                                                             |
|            |            | (345)  | The workshop card shall be able to store vehicles used<br>data records in the same manner as a driver card.                                                                                                                                                                                                                                                                                                                                                    |
| <b>▼M3</b> |            | (346)  | The workshop card shall be able to store 8 such records.                                                                                                                                                                                                                                                                                                                                                                                                       |
|            | 4.5.4.2.10 |        | Places and positions where daily work periods start and/or end data                                                                                                                                                                                                                                                                                                                                                                                            |
|            |            | (347)  | The workshop card shall be able to store places and<br>positions where daily work periods begin and/or end<br>data records in the same manner as a driver card.                                                                                                                                                                                                                                                                                                |
|            |            | (348)  | The workshop card shall be able to store 4 pairs of such<br>records.                                                                                                                                                                                                                                                                                                                                                                                           |
| <b>▼B</b>  | 4.5.4.2.11 |        | Card session data                                                                                                                                                                                                                                                                                                                                                                                                                                              |
| <b>B</b>   | 4.5.4.2.12 | (350)  | Control activity data<br>The workshop card shall be able to store a control activity data record in the same manner as a driver card.                                                                                                                                                                                                                                                                                                                          |
|            | 4.5.4.2.13 | (351)  | Vehicle units used data<br>The workshop card shall be able to store the following data related to the different vehicle units in which the card was used:<br>— the date and time of the beginning of the period of use of the vehicle unit (i.e. first card insertion in the vehicle unit for the period),<br>— the manufacturer of the vehicle unit,<br>— the vehicle unit type,<br>— the vehicle unit software version number.                               |
| <b>▼M3</b> |            | (352)  | The workshop card shall be able to store 8 such records.                                                                                                                                                                                                                                                                                                                                                                                                       |
| <b>▼M1</b> | 4.5.4.2.14 | (353)  | Three hours accumulated driving places data<br>The workshop card shall be able to store the following data related to the position of the vehicle where the accumulated driving time reaches a multiple of three hours:<br>— the date and time when the accumulated driving time reaches a multiple of three hours,<br>— the position of the vehicle,<br>— the GNSS accuracy, date and time when the position was determined,<br>— the vehicle odometer value. |
| <b>▼M3</b> |            | (354)  | The workshop card shall be able to store 24 such records.                                                                                                                                                                                                                                                                                                                                                                                                      |
| <b>▼B</b>  | 4.5.4.2.15 | (355)  | Specific conditions data<br>The workshop card shall be able to store data relevant to specific conditions in the same manner as the driver card.                                                                                                                                                                                                                                                                                                               |
| <b>▼M3</b> |            | (356)  | The workshop card shall be able to store 4 such records.                                                                                                                                                                                                                                                                                                                                                                                                       |
|            | 4.5.4.2.16 | (356a) | Authentication status for positions related to places where daily work periods start and/or end (not accessed by version 1 of second generation vehicle units)<br>The workshop card shall be able to store additional data related to places where daily work periods start and/or end                                                                                                                                                                         |

(349) The workshop card shall be able to store a card session data record in the same manner as a driver card.

- (356b) The workshop card memory shall be able to store 4 pairs of such records.
- 4.5.4.2.17 Authentication status for positions where three hours accumulated driving are reached (not accessed by version 1 of second generation vehicle units)
  - (356c) The workshop card shall be able to store additional data related to the position of the vehicle where the accumulated driving time reaches a multiple of three hours in the same manner as a driver card.
  - (356d) The workshop card shall be able to store 24 such records.
- 4.5.4.2.18 Border crossings (not accessed by version 1 of second generation vehicle units)
  - (356e) The workshop card shall be able to store the border crossings in the same manner as a driver card.
  - (356f) The workshop card memory shall be able to store 4 such records.
- 4.5.4.2.19 Load/unload operations (not accessed by version 1 of second generation vehicle units)
  - (356g) The workshop card shall be able to store the load/unload operations in the same manner as a driver card.
  - (356h) The workshop card shall be able to store 8 load, unload or simultaneous load/unload operations.
- 4.5.4.2.20 Load type entries (not accessed by version 1 of second generation vehicle units)
  - (356i) The workshop card shall be able to store the load type entries in the same manner as a driver card.
  - (356j) The workshop card shall be able to store 4 such records.
- 4.5.4.2.21 Calibration Additional Data (not accessed by version 1 of second generation vehicle units)
  - (356k) The workshop card shall be able to store additional calibration data only applicable for version 2:
    - the old date and time and the vehicle identification number, which shall be exactly the same values as the one stored in EF Calibration under the DF Tachograph\_G2,
    - the by-default load type entered during this calibration,
    - the country in which the calibration has been performed, and the date time when the position used to determine this country was provided by the GNSS receiver.

(356l) The workshop card shall be able to store 255 such records.

4.5.4.2.22 VU configurations (not accessed by version 1 of second generation vehicle units)

- (356m) The workshop card shall be able to store the cardholder tachograph specific settings.
- (356n) The workshop card storage capacity for cardholder tachograph specific settings shall be 3072 bytes.

## **B**

| 4.5.5<br>Control card |
|-----------------------|
|-----------------------|

- 4.5.5.1 T a c h o g r a p h a p p l i c a t i o n ( a c c e s s i b l e t o f i r s t a n d s e c o n d g e n e r a t i o n v e h i c l e u n i t s )
- 4.5.5.1.1 Application identification
  - (357) The control card shall be able to store the following application identification data:
    - tachograph application identification,
    - type of tachograph card identification.
- 4.5.5.1.2 Keys and certificates
  - (358) The control card shall be able to store a number of cryptographic keys and certificates, as specified in Appendix 11 part A.

## 4.5.5.1.3 Card identification

- (359) The control card shall be able to store the following card identification data:
  - card number,
  - issuing Member State, issuing authority name, issue date,
  - card beginning of validity date, card expiry date (if any).

## 4.5.5.1.4 Card holder identification

- (360) The control card shall be able to store the following card holder identification data:
  - control body name,
  - control body address,
  - surname of the holder,
  - first name(s) of the holder,
  - preferred language.

## 4.5.5.1.5 Control activity data

- (361) The control card shall be able to store the following control activity data:
  - date and time of the control,

**▼M3**

— type of the control (displaying and/or printing and/or VU downloading and/or card downloading),

- VRN and Member State registering authority of the controlled vehicle,
- card number and card issuing Member State of the driver card controlled.
- (362) The control card shall be able to hold at least 230 such records.
- 4.5.5.2 T a c h o g r a p h G 2 a p p l i c a t i o n ( n o t a c c e s s i b l e t o f i r s t g e n e r a t i o n v e h i c l e u n i t )
- 4.5.5.2.1 Application identification
  - (363) The control card shall be able to store the following application identification data:
    - tachograph application identification,
    - type of tachograph card identification.

- 4.5.5.2.1.1 Additional application identification (not accessed by version 1 of second generation vehicle units)
  - (363a) The control card shall be able to store additional application identification data only applicable for version 2.

## **B**

- 4.5.5.2.2 Keys and certificates
  - (364) The control card shall be able to store a number of cryptographic keys and certificates, as specified in Appendix 11 part B.

### 4.5.5.2.3 Card identification

- (365) The control card shall be able to store the following card identification data:
  - card number,
  - issuing Member State, issuing authority name, issue date,
  - card beginning of validity date, card expiry date (if any).
- 4.5.5.2.4 Card holder identification
  - (366) The control card shall be able to store the following card holder identification data:
    - control body name,
    - control body address,
    - surname of the holder,
    - first name(s) of the holder,
    - preferred language.

## 4.5.5.2.5 Control activity data

- (367) The control card shall be able to store the following control activity data:
  - date and time of the control,

— type of the control (displaying and/or printing and/or VU downloading and/or card downloading and/or roadside calibration checking)

- period downloaded (if any),
- VRN and Member State registering authority of the controlled vehicle,
- card number and card issuing Member State of the driver card controlled.
- (368) The control card shall be able to hold at least 230 such records.

### **M3**

- 4.5.5.2.6 VU configurations (not accessed by version 1 of second generation vehicle units)
  - (368a) The control card shall be able to store the cardholder tachograph specific settings.
  - (368b) The control card storage capacity for cardholder tachograph specific settings shall be 3072 bytes.

## **B**

- 4.5.6 *Company card*
- 4.5.6.1 T a c h o g r a p h a p p l i c a t i o n ( a c c e s s i b l e t o f i r s t a n d s e c o n d g e n e r a t i o n v e h i c l e u n i t s )
- 4.5.6.1.1 Application identification
  - (369) The company card shall be able to store the following application identification data:
    - tachograph application identification,
    - type of tachograph card identification.

### 4.5.6.1.2 Keys and Certificates

(370) The company card shall be able to store a number of cryptographic keys and certificates, as specified in Appendix 11 part A.

### 4.5.6.1.3 Card identification

- (371) The company card shall be able to store the following card identification data:
  - card number,
  - issuing Member State, issuing authority name, issue date,
  - card beginning of validity date, card expiry date (if any).

## 4.5.6.1.4 Card holder identification

- (372) The company card shall be able to store the following card holder identification data:
  - company name,
  - company address.

### 4.5.6.1.5 Company activity data

- (373) The company card shall be able to store the following company activity data:
  - date and time of the activity,
  - type of the activity (VU locking in and/or out, and/or VU downloading and/or card downloading)
  - period downloaded (if any),
  - VRN and Member State registering authority of vehicle,
  - card number and card issuing Member State (in case of card downloading).
- (374) The company card shall be able to hold at least 230 such records.
- 4.5.6.2 T a c h o g r a p h G 2 a p p l i c a t i o n ( n o t a c c e s s i b l e t o f i r s t g e n e r a t i o n v e h i c l e u n i t )
- 4.5.6.2.1 Application identification
  - (375) The company card shall be able to store the following application identification data:
    - tachograph application identification,
    - type of tachograph card identification.

### **M3**

- 4.5.6.2.1.1 Additional application identification (not accessed by version 1 of second generation vehicle units)
  - (375a) The company card shall be able to store additional application identification data only applicable for version 2.

## **B**

- 4.5.6.2.2 Keys and certificates
  - (376) The company card shall be able to store a number of cryptographic keys and certificates, as specified in Appendix 11 part B.

## 4.5.6.2.3 Card identification

- (377) The company card shall be able to store the following card identification data:
  - card number,
  - issuing Member State, issuing authority name, issue date,
  - card beginning of validity date, card expiry date (if any).
- 4.5.6.2.4 Card holder identification
  - (378) The company card shall be able to store the following card holder identification data:
    - company name,
    - company address.

### 4.5.6.2.5 Company activity data

- (379) The company card shall be able to store the following company activity data:
  - date and time of the activity,
  - type of the activity (VU locking in and/or out, and/or VU downloading and/or card downloading)
  - period downloaded (if any),
  - VRN and Member State registering authority of vehicle,
  - card number and card issuing Member State (in case of card downloading).
- (380) The company card shall be able to hold at least 230 such records.

### **M3**

- 4.5.6.2.6 VU configurations (not accessed by version 1 of second generation vehicle units)
  - (380a) The company card shall be able to store the cardholder tachograph specific settings.
  - (380b) The company card storage capacity for cardholder tachograph specific settings shall be 3072 bytes.

## **B**

### 5. INSTALLATION OF RECORDING EQUIPMENT

## 5.1 **Installation**

- (381) New recording equipment shall be delivered non-activated to fitters or vehicle manufacturers, with all calibration parameters, as listed in Chapter 3.21, set to appropriate and valid default values. Where no particular value is appropriate, literal parameters shall be set to strings of '?' and numeric parameters shall be set to '0'. Delivery of security relevant parts of the recording equipment can be restricted if required during security certification.
- (382) Before its activation, the recording equipment shall give access to the calibration function even if not in calibration mode.
- (383) Before its activation, the recording equipment shall neither record nor store data referred to by the requirements 102 to 133 inclusive. Nevertheless, before its activation, the recording equipment may record and store the security breach attempt events in accordance with requirement 117, and the recording equipment faults in accordance with requirement 118.
- (384) During installation, vehicle manufacturers shall pre-set all known parameters.
- (385) Vehicle manufacturers or fitters shall activate the installed recording equipment at the latest before the vehicle is used in scope of Regulation (EC) No 561/2006.

### **B**

**▼B**

- (386) The activation of the recording equipment shall be triggered automatically by the first insertion of a valid workshop card in either of its card interface devices.
- (387) Specific pairing operations required between the motion sensor and the vehicle unit, if any, shall take place automatically before or during activation.
- (388) In a similar way, specific coupling operations between the external GNSS facility and the vehicle unit, if any, shall take place automatically before or during activation.
- (389) After its activation, the recording equipment shall fully enforce functions and data access rights.
- (390) After its activation, the recording equipment shall communicate to the remote communication facility the secured data necessary for the purpose of targeted roadside checks.
- (391) The recording and storing functions of the recording equipment shall be fully operational after its activation.
- **▼M3**

- (392) Installation shall be followed by a calibration. The first calibration may not necessarily include entry of the vehicle registration identification (VRN and Member State), when it is not known by the approved workshop having to undertake this calibration. In these circumstances, it shall be possible, for the vehicle owner, and at this time only, to enter the VRN and the Member State using his company card prior to using the vehicle in scope of Regulation (EC) No 561/2006 (e.g by using commands through an appropriate menu structure of the vehicle unit's man-machine interface). Any update or confirmation of this entry shall only be possible using a workshop card.
- (393) The installation of an external GNSS facility requires the coupling with the vehicle unit and the subsequent verification of the GNSS position information.
- (394) The recording equipment must be positioned in the vehicle in such a way as to allow the driver to access the necessary functions from his seat.

## 5.2 **Installation plaque**

(395) **►M3** After the recording equipment has been checked on installation, an installation plaque, engraved or printed in a permanent way, which is clearly visible and easily accessible shall be affixed onto the recording equipment. In cases where this is not possible, the plaque shall be affixed to the vehicle's 'B' pillar so that it is clearly visible. For vehicles that do not have a 'B' pillar, the installation plaque should be affixed in the area of the door of the vehicle and be clearly visible in all cases. ◄

> After every inspection by an approved fitter or workshop, a new plaque shall be affixed in place of the previous one.

| ▼M1 | (396)                                                                                                                                                                                                                                                                                                                                                                                                                          | The plaque shall bear at least the following details:<br><br>— name, address or trade name of the approved fitter or workshop,<br><br>— characteristic coefficient of the vehicle, in the form<br>'w = ... imp/km',<br><br>— constant of the recording equipment, in the form<br>'k = ... imp/km',<br><br>— effective circumference of the wheel tyres in the form<br>'l = ... mm',<br><br>— tyre size,<br><br>— the date on which the characteristic coefficient of the vehicle and the effective circumference of the wheel tyres were measured,<br><br>— the vehicle identification number,<br><br>— the presence (or not) of an external GNSS facility,<br><br>— the serial number of the external GNSS facility, if applicable, |
|-----|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| ▼M3 |                                                                                                                                                                                                                                                                                                                                                                                                                                | — the serial number of the remote communication facility, if any,                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                    |
| ▼M1 |                                                                                                                                                                                                                                                                                                                                                                                                                                | — the serial number of all the seals in place,                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                       |
|     |                                                                                                                                                                                                                                                                                                                                                                                                                                | — the part of the vehicle where the adaptor, if any, is installed,                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                   |
|     |                                                                                                                                                                                                                                                                                                                                                                                                                                | — the part of the vehicle where the motion sensor is installed, if not connected to the gear-box or an adaptor is not being used,                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                    |
|     |                                                                                                                                                                                                                                                                                                                                                                                                                                | — a description of the colour of the cable between the adaptor and that part of the vehicle providing its incoming impulses,                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                         |
|     |                                                                                                                                                                                                                                                                                                                                                                                                                                | — the serial number of the embedded motion sensor of the adaptor.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                    |
| ▼M3 |                                                                                                                                                                                                                                                                                                                                                                                                                                | — the by-default load type associated to the vehicle.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                |
| ▼B  | (397)                                                                                                                                                                                                                                                                                                                                                                                                                          | For M1 and N1 vehicles only, and which are fitted with an adaptor in conformity with Commission Regulation (EC) No 68/2009 (1) as last amended and where it is not possible to include all the information necessary, as described in Requirement 396, a second, additional, plaque may be used. In such cases, this additional                                                                                                                                                                                                                                                                                                                                                                                                      |
| ▼B  | This second, additional plaque, if used, shall be affixed next to or beside the first primary plaque described in Requirement 396, and shall have the same protection level. Furthermore the secondary plaque shall also bear the name, address or trade name of the approved fitter or workshop that carried out the installation, and the date of installation.                                                              |                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                      |
|     | 5.3 Sealing                                                                                                                                                                                                                                                                                                                                                                                                                    |                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                      |
|     | (398) The following parts shall be sealed:                                                                                                                                                                                                                                                                                                                                                                                     |                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                      |
|     | — Any connection which, if disconnected, would cause undetectable alterations to be made or undetectable data loss (this may e.g. apply for the motion sensor fitting on the gearbox, the adaptor for M1/N1 vehicles, the external GNSS connection or the vehicle unit);                                                                                                                                                       |                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                      |
|     | — The installation plaque, unless it is attached in such a way that it cannot be removed without the markings thereon being destroyed.                                                                                                                                                                                                                                                                                         |                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                      |
| ▼M1 | (398a) The seals mentioned above shall be certified according to the standard EN 16882:2016.                                                                                                                                                                                                                                                                                                                                   |                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                      |
| ▼B  | (399) The seals mentioned above may be removed:                                                                                                                                                                                                                                                                                                                                                                                |                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                      |
|     | — In case of emergency,                                                                                                                                                                                                                                                                                                                                                                                                        |                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                      |
|     | — To install, to adjust or to repair a speed limitation device or any other device contributing to road safety, provided that the recording equipment continues to function reliably and correctly and is resealed by an approved fitter or workshop (in accordance with Chapter 6) immediately after fitting the speed limitation device or any other device contributing to road safety or within seven days in other cases. |                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                      |
|     | (400) On each occasion that these seals are broken a written statement giving the reasons for such action shall be prepared and made available to the competent authority.                                                                                                                                                                                                                                                     |                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                      |
|     | (401) Seals shall hold an identification number, allocated by its manufacturer. This number shall be unique and distinct from any other seal number allocated by any other seals manufacturer.                                                                                                                                                                                                                                 |                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                      |
| ▼M1 | This unique identification number is defined as: MMNNNNNNNN by non-removable marking, with MM as unique manufacturer identification (database registration to be managed by EC) and NNNNNNNN seal alpha-numeric number, unique in the manufacturer domain.                                                                                                                                                                     |                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                      |
| ▼B  | (402) The seals shall have a free space where approved fitters, workshops or vehicle manufacturers can add a special mark according the Article 22(3) of Regulation (EU)                                                                                                                                                                                                                                                       |                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                      |

described in Requirement 396.

<sup>(1)</sup> Commission Regulation (EC) No 68/2009 of 23 January 2009 adapting for the ninth time to technical progress Council Regulation (EEC) No 3821/85 on recording equipment in road transport (OJ L 21, 24.1.2009, p. 3).

No 165/2014.

This mark shall not cover the seal identification number.

- (403) Seals manufacturers shall be registered in a dedicated database when they get a seal model certified according to EN 16882:2016 and shall make their identification seals numbers public through a procedure to be established by the European Commission.
- (404) Approved workshops and vehicle manufacturers shall, in the frame of Regulation (EU) No 165/2014, only use seals certified according to EN 16882:2016 from those of the seals manufacturers listed in the database mentioned above.
- (405) Seal manufacturers and their distributors shall maintain full traceability records of the seals sold to be used in the frame of Regulation (EU) No 165/2014 and shall be prepared to produce them to competent national authorities whenever need be.
- (406) Seals unique identification numbers shall be visible on the installation plaque.

### 6. CHECKS, INSPECTIONS AND REPAIRS

Requirements on the circumstances in which seals may be removed, as referred to in Article 22(5) of Regulation (EU) No 165/2014, are defined in Chapter 5.3 of this annex.

### 6.1 **Approval of fitters, workshops and vehicle manufacturers**

The Member States approve, regularly control and certify the bodies to carry out:

- installations,
- checks,
- inspections,
- repairs.

Workshop cards shall be issued only to fitters and/or workshops approved for the activation and/or the calibration of recording equipment in conformity with this annex and, unless duly justified:

- who are not eligible for a company card;
- and whose other professional activities do not present a potential compromise of the overall security of the system as required in Appendix 10.

**▼M1**

## 6.2 **Check of new or repaired components**

(407) Every individual device, whether new or repaired, shall be checked in respect of its proper operation and the accuracy of its reading and recordings, within the limits laid down in Chapter 3.2.1, 3.2.2, 3.2.3 and 3.3.

### **B**

**▼M1**

| ▼B  | 6.3   | Installation inspection                                                                                                                                                                                                                                                                                                                                                                                                      |  |
|-----|-------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|--|
| ▼M1 | (408) | When being fitted to a vehicle, the whole installation<br>(including the recording equipment) shall comply with<br>the provisions relating to maximum tolerances laid down<br>in Chapter 3.2.1, 3.2.2, 3.2.3 and 3.3. The whole instal-<br>lation shall be sealed in accordance with Chapter 5.3 and it<br>shall include a calibration.                                                                                      |  |
| ▼B  | 6.4   | Periodic inspections                                                                                                                                                                                                                                                                                                                                                                                                         |  |
| ▼M3 | (409) | Periodic inspections of the equipment fitted to the vehicles<br>shall take place after any repair of the equipment, or after<br>any alteration of the characteristic coefficient of the vehicle<br>or of the effective circumference of the tyres, or after<br>equipment UTC time is wrong by more than 5 minutes,<br>or when the VRN has changed, and at least once within<br>two years (24 months) of the last inspection. |  |
| ▼B  | (410) | These inspections shall include the following checks:                                                                                                                                                                                                                                                                                                                                                                        |  |
|     |       | — that the recording equipment is working properly,<br>including the data storage in tachograph cards function<br>and the communication with remote communication<br>readers,                                                                                                                                                                                                                                                |  |
|     |       | — that compliance with the provisions of chapter 3.2.1<br>and 3.2.2 on the maximum tolerances on installation<br>is ensured,                                                                                                                                                                                                                                                                                                 |  |
|     |       | — that compliance with the provisions of chapter 3.2.3<br>and 3.3 is ensured,                                                                                                                                                                                                                                                                                                                                                |  |
|     |       | — that the recording equipment carries the type approval<br>mark,                                                                                                                                                                                                                                                                                                                                                            |  |
|     |       | — that the installation plaque, as defined by Requirement<br>396, and the descriptive plaque, as defined by<br>Requirement 225, are affixed,                                                                                                                                                                                                                                                                                 |  |
|     |       | — the tyre size and the actual circumference of the tyres,                                                                                                                                                                                                                                                                                                                                                                   |  |
|     |       | — that there are no manipulation devices attached to the<br>equipment,                                                                                                                                                                                                                                                                                                                                                       |  |
|     |       | — that seals are correctly placed, in good state, that their<br>identification numbers are valid (referenced seal manu-<br>facturer in the EC database) and that their identification<br>numbers correspond to the installation plaque markings<br>(see requirement 401)                                                                                                                                                     |  |

— that the version identifier of the stored digital map is the most recent one.

(410a) In case of detection of a manipulation by the competent national authorities, the vehicle may be sent to an authorised workshop for a recalibration of the recording equipment.

- (411) If one of the events listed in Chapter 3.9 (Detection of Events and/or Faults) is found to have occurred since the last inspection and is considered by tachograph manufacturers and/or national authorities as potentially putting the security of the equipment at risk, the workshop shall:
  - a. make a comparison between the motion sensor identification data of the motion sensor plugged into the gearbox with that of the paired motion sensor registered in the vehicle unit;
  - b. check if the information recorded on the installation plaque matches with the information contained within the vehicle unit record;
  - c. check if the motion sensor serial number and approval number, if printed on the body of the motion sensor, matches the information stored in the recording equipment data memory;
  - d. compare identification data marked on the descriptive plaque of the external GNSS facility, if any, to the ones stored in the vehicle unit data memory;
- (412) Workshops shall keep traces in their inspection reports of any findings concerning broken seals or manipulations devices. These reports shall be kept by workshops for at least 2 years and made available to the Competent Authority whenever requested to do so.
- (413) These inspections shall include a calibration and a preventive replacement of the seals whose fitting is under the responsibility of workshops.

### 6.5 **Measurement of errors**

- (414) The measurement of errors on installation and during use shall be carried out under the following conditions, which are to be regarded as constituting standard test conditions:
  - vehicle unladen, in normal running order,
  - tyre pressures in accordance with the manufacturer's instructions,
  - tyre wear, within the limits allowed by national law,
  - vehicle movement:
  - the vehicle shall advance under its own engine power in a straight line on level ground and at a speed of 50 ± 5 km/h. The measuring distance shall be at least 1 000 m.
  - provided that it is of comparable accuracy, alternative methods, such as a suitable test bench, may also be used for the test.

### 6.6 **Repairs**

- (415) Workshops shall be able to download data from the recording equipment to give the data back to the appropriate transport company.
- (416) Approved workshops shall issue to transport companies a certificate of data un-downloadability where the malfunction of the recording equipment prevents previously recorded data to be downloaded, even after repair by this workshop. The workshops will keep a copy of each issued certificate for at least two years.

### 7. CARD ISSUING

The card issuing processes set-up by the Member States shall conform to the following:

- (417) The card number of the first issue of a tachograph card to an applicant shall have a consecutive index (if applicable) and a replacement index and a renewal index set to '0'.
- (418) The card numbers of all non-personal tachograph cards issued to a single control body or a single workshop or a single transport company shall have the same first 13 digits, and shall all have a different consecutive index.
- (419) A tachograph card issued in replacement of an existing tachograph card shall have the same card number than the replaced one except the replacement index which shall be raised by '1' (in the order 0, …, 9, A, …, Z).
- (420) A tachograph card issued in replacement of an existing tachograph card shall have the same card expiry date as the replaced one.
- (421) A tachograph card issued in renewal of an existing tachograph card shall have the same card number as the renewed one except the replacement index which shall be reset to '0' and the renewal index which shall be raised by '1' (in the order 0, …, 9, A, …, Z).
- (422) The exchange of an existing tachograph card, in order to modify administrative data, shall follow the rules of the renewal if within the same Member State, or the rules of a first issue if performed by another Member State.
- (423) The 'card holder surname' for non-personal workshop or control cards shall be filled with workshop or control body name or with the fitter or control officer's name would Member States so decide.
- (424) Member States shall exchange data electronically in order to ensure the uniqueness of driver cards that they issue in accordance with Article 31 of Regulation (EU) No 165/2014.

### 8. TYPE-APPROVAL OF RECORDING EQUIPMENT AND TACHOGRAPH CARDS

### 8.1 **General points**

**▼M1**

For the purpose of this chapter, the words 'recording equipment' mean 'recording equipment or its components'. No type approval is required for the cable(s) linking the motion sensor to the VU, the external GNSS facility to the VU or the external remote communication facility to the VU. The paper, for use by the recording equipment, shall be considered as a component of the recording equipment.

Any manufacturer may ask for type approval of recording equipment component(s) with any other recording equipment component(s), provided each component complies with the requirements of this annex. Alternately, manufacturers may also ask for type approval of recording equipment.

As described in definition (10) in Article 2 of this Regulation, vehicle units have variants in components assembly. Whatever the vehicle unit components assembly, the external antenna and (if applicable) the antenna splitter connected to the GNSS receiver or to the remote communication facility are not part of the vehicle unit type approval.

Nevertheless, manufacturers having obtained type approval for recording equipment shall maintain a publicly available list of compatible antennas and splitters with each type approved vehicle unit, external GNSS facility and external remote communication facility.

- (425) Recording equipment shall be submitted for approval complete with any integrated additional devices.
- (426) Type approval of recording equipment and of tachograph cards shall include security related tests, functional tests and interoperability tests. Positive results to each of these tests are stated by an appropriate certificate.
- (427) Member States type approval authorities will not grant a type approval certificate as long as they do not hold:
  - a security certificate (if requested by this Annex),
  - a functional certificate,
  - and an interoperability certificate (if requested by this Annex)

for the recording equipment or the tachograph card, subject of the request for type approval.

- **▼B**
- (428) Any modification in software or hardware of the equipment or in the nature of materials used for its manufacture shall, before being used, be notified to the authority which granted type-approval for the equipment. This authority shall confirm to the manufacturer the extension of the type approval, or may require an update or a confirmation of the relevant functional, security and/or interoperability certificates.

### **B**

**▼B**

- (429) Procedures to update in-situ recording equipment software shall be approved by the authority which granted type approval for the recording equipment. Software update must not alter nor delete any driver activity data stored in the recording equipment. Software may be updated only under the responsibility of the equipment manufacturer.
- (430) Type approval of software modifications aimed to update a previously type approved recording equipment may not be refused if such modifications only apply to functions not specified in this Annex. Software update of a recording equipment may exclude the introduction of new character sets, if not technically feasible.

### 8.2 **Security certificate**

- (431) The security certificate is delivered in accordance with the provisions of Appendix 10 of this Annex. Recording equipment components to be certified are vehicle unit, motion sensor, external GNSS facility and tachograph cards.
- (432) In the exceptional circumstance that the security certification authorities refuse to certify new equipment on the ground of obsolescence of the security mechanisms, type approval shall continue to be granted only in these specific and exceptional circumstances, and when no alternative solution, compliant with the Regulation, exists.
- (433) In this circumstance the Member State concerned shall, without delay, inform the European Commission, which shall, within twelve calendar months of the grant of the type approval, launch a procedure to ensure that the level of security is restored to its original levels.

### 8.3 **Functional certificate**

- (434) Each candidate for type approval shall provide the Member State's type approval authority with all the material and documentation that the authority deems necessary.
- (435) Manufacturers shall provide the relevant samples of type approval candidate products and associated documentation required by laboratories appointed to perform functional tests, and within one month of the request being made. Any costs resulting from this request shall be borne by the requesting entity. Laboratories shall treat all commercially sensitive information in confidence.
- (436) A functional certificate shall be delivered to the manufacturer only after all functional tests specified in Appendix 9, at least, have been successfully passed.
- (437) The type approval authority delivers the functional certificate. This certificate shall indicate, in addition to the name of its beneficiary and the identification of the model, a detailed list of the tests performed and the results obtained.

- (438) The functional certificate of any recording equipment component shall also indicate the type approval numbers of the other type approved compatible recording equipment components tested for its certification.
- (439) The functional certificate of any recording equipment component shall also indicate the ISO or CEN standard against which the functional interface has been certified.

### 8.4 **Interoperability certificate**

- (440) Interoperability tests are carried out by a single laboratory under the authority and responsibility of the European Commission.
- (441) The laboratory shall register interoperability test requests introduced by manufacturers in the chronological order of their arrival.
- (442) Requests will be officially registered only when the laboratory is in possession of:
  - the entire set of material and documents necessary for such interoperability tests,
  - the corresponding security certificate,
  - the corresponding functional certificate,

The date of the registration of the request shall be notified to the manufacturer.

## **M3**

- (443) No interoperability tests shall be carried out by the laboratory, for recording equipment or tachograph cards that have not passed the vulnerability analysis of their security evaluation and a functional evaluation, except in the exceptional circumstances described in requirement 432.
- **▼B**
- (444) Any manufacturer requesting interoperability tests shall commit to leave to the laboratory in charge of these tests the entire set of material and documents which he provided

to carry out the tests.

- (445) The interoperability tests shall be carried out, in accordance with the provisions of Appendix 9 of this Annex, with respectively all the types of recording equipment or tachograph cards:
  - for which type approval is still valid or,
  - for which type approval is pending and that have a valid interoperability certificate.
- (446) The interoperability tests shall cover all generations of recording equipment or tachograph cards still in use.
- (447) The interoperability certificate shall be issued by the laboratory to the manufacturer only after all required interoperability tests have been successfully passed and after the manufacturer has shown that both a valid functional certificate and a valid security certificate for the product has been granted, except in the exceptional circumstances described in requirement 432.

### **B**

- (448) If the interoperability tests are not successful with one or more of the recording equipment or tachograph card(s), the interoperability certificate shall not be delivered, until the requesting manufacturer has realised the necessary modifications and has succeeded the interoperability tests. The laboratory shall identify the cause of the problem with the help of the manufacturers concerned by this interoperability fault and shall attempt to help the requesting manufacturer in finding a technical solution. In the case where the manufacturer has modified its product, it is the manufacturer's responsibility to ascertain from the relevant authorities that the security certificate and the functional certificates are still valid.
- (449) The interoperability certificate is valid for six months. It is revoked at the end of this period if the manufacturer has not received a corresponding type approval certificate. It is forwarded by the manufacturer to the type approval authority of the Member State who has delivered the functional certificate.
- (450) Any element that could be at the origin of an interoperability fault shall not be used for profit or to lead to a dominant position.

## 8.5 **Type-approval certificate**

- (451) The type approval authority of the Member State may deliver the type approval certificate as soon as it holds the three required certificates.
- (452) The type approval certificate of any recording equipment component shall also indicate the type approval numbers of the other type approved interoperable recording equipment.
- (453) The type approval certificate shall be copied by the type approval authority to the laboratory in charge of the interoperability tests at the time of deliverance to the manufacturer.
- (454) The laboratory competent for interoperability tests shall run a public web site on which will be updated the list of recording equipment or tachograph cards models:
  - for which a request for interoperability tests have been registered,
  - having received an interoperability certificate (even provisional),
  - having received a type approval certificate.

## 8.6 **Exceptional procedure: first interoperability certificates for 2nd generation recording equipment and tachograph cards**

(455) Until four months after a first couple of 2nd generation recording equipment and 2nd generation tachograph cards (driver, workshop, control and company cards) have been certified to be interoperable, any interoperability certificate delivered (including the first ones), regarding requests registered during this period, shall be considered provisional.

- (456) If at the end of this period, all products concerned are mutually interoperable, all corresponding interoperability certificates shall become definitive.
- (457) If during this period, interoperability faults are found, the laboratory in charge of interoperability tests shall identify the causes of the problems with the help of all manufacturers involved and shall invite them to realize the necessary modifications.
- (458) If at the end of this period, interoperability problems still remain, the laboratory in charge of interoperability tests, with the collaboration of the manufacturers concerned and with the type approval authorities who delivered the corresponding functional certificates shall find out the causes of the interoperability faults and establish which modifications should be made by each of the manufacturers concerned. The search for technical solutions shall last for a maximum of two months, after which, if no common solution is found, the Commission, after having consulted the laboratory in charge of interoperability tests, shall decide which equipment(s) and cards get a definitive interoperability certificate and state the reasons why.
- (459) Any request for interoperability tests, registered by the laboratory between the end of the four month period after the first provisional interoperability certificate has been delivered and the date of the decision by the Commission referred to in requirement 455, shall be postponed until the initial interoperability problems have been solved. Those requests are then processed in the chronological order of their registration.