*Appendix 13*

### ITS INTERFACE

#### TABLE OF CONTENTS

- 1. INTRODUCTION
- 1.1. Scope
- 1.2. Acronyms and definitions
- 2. REFERENCED STANDARDS
- 3. ITS INTERFACE WORKING PRINCIPLES
- 3.1. Communication technology
- 3.2. Available services
- 3.3. Access through the ITS interface
- 3.4. Data available and need of driver consent
- 4. LIST OF DATA AVAILABLE THROUGH THE ITS INTERFACE AND PERSONAL/NOT PERSONAL CLASSIFICATION

## 1. INTRODUCTION

- 1.1. **Scope**
  - ITS\_01 This Appendix specifies the basics of the communication through the tachograph interface with Intelligent Transport Systems (ITS), requested in Articles 10 and 11 of Regulation (EU) No 165/2014.
  - ITS\_02 The ITS interface shall allow external devices to obtain data from the tachograph, to use tachograph services and also to provide data to the tachograph.

Other tachograph interfaces (e.g. CAN bus) may also be used for that purpose.

This Appendix does not specify:

- how data provided through the ITS interface are collected and managed within the tachograph,
- the form of presentation of collected data to applications hosted on the external device,
- the ITS security specification in addition to what provides Bluetooth®,
- the Bluetooth® protocols used by the ITS interface.

## 1.2. Acronyms and definitions

The following acronyms and definitions specific to this Appendix are used:

| GNSS     | Global Navigation Satellite System                 |
|----------|----------------------------------------------------|
| ITS      | Intelligent Transport System                       |
| OSI      | Open Systems Interconnection                       |
| VU       | Vehicle Unit                                       |
| ITS unit | an external device or application using the VU ITS |

## 2. REFERENCED STANDARDS

interface.

ITS\_03 This Appendix refers to and depends upon all or parts of the following regulations and standards. Within the clauses of this Appendix, the relevant standards, or relevant clauses of standards, are referred to. In the event of any contradiction the clauses of this Appendix shall take precedence.

Standards referenced to in this Appendix are:

— Bluetooth® – Core Version 5.0.

- ISO 16844-7: Road vehicles -Tachograph systems Part 7: Parameters
- ISO/IEC 7498-1:1994 Information technology Open Systems Interconnection - Basic Reference Model, the Basic Model

# 3. ITS INTERFACE WORKING PRINCIPLES

ITS\_04 The VU shall be responsible to keep updated and maintain tachograph data transmitted through the ITS interface, without any involvement of the ITS interface.

# 3.1. Communication technology

- ITS\_05 Communication through the ITS interface shall be performed via Bluetooth® interface and be compatible to Bluetooth® Low Energy according to Bluetooth version 5.0 or newer.
- ITS\_06 The communication between the VU and the ITS unit shall be established after a Bluetooth® pairing process has been completed.
- ITS\_07 A secure and encrypted communication shall be established between the VU and the ITS unit, in accordance with the Bluetooth® specification mechanisms. This Appendix does not specify encryption or other security mechanisms in addition to what Bluetooth® provides.
- ITS\_08 Bluetooth® is using a server/client model to control the transmission of data between devices, in which the VU shall be the server and the ITS unit shall be the client.

## 3.2. Available services

ITS\_09 The data to be transmitted through the ITS interface in accordance with point 4 shall be made available through the services specified in Appendix 7 and Appendix 8. In addition, the VU shall make available to the ITS unit the services that are necessary for manual data entry in accordance with requirement 61 of Annex IC, and optionally, for other data entries in real time.

### Figure 1

**partition of the communication through the ITS interface according to the OSI Model layers**

![](_page_2_Figure_6.jpeg)

- ITS\_10 When the download interface is used via the front connector, the VU shall not provide the download services specified in Appendix 7 via ITS Bluetooth® connection.
- ITS\_11 When the calibration interface is used via the front connector, the VU shall not provide the calibration services specified in Appendix 8 via ITS Bluetooth® connection.

## 3.3. Access through the ITS interface

- ITS\_12 The ITS interface shall provide a wireless access to all services specified in Appendix 7 and Appendix 8, in replacement of a cable connection to the front connector for calibration and download specified in Appendix 6.
- ITS\_13 The VU shall make the ITS interface available to the user according to the combination of valid tachograph cards inserted in the VU, as specified in Table 1.

| Availability of the ITS<br>interface |               | Driver slot   |             |               |               |               |
|--------------------------------------|---------------|---------------|-------------|---------------|---------------|---------------|
|                                      |               | No card       | Driver card | Control card  | Workshop card | Company card  |
| Co – driverslot                      | No card       | Not available | Available   | Available     | Available     | Available     |
|                                      | Driver card   | Available     | Available   | Available     | Available     | Available     |
|                                      | Control card  | Available     | Available   | Available     | Not available | Not available |
|                                      | Workshop card | Available     | Available   | Not available | Available     | Not available |
|                                      | Company card  | Available     | Available   | Not available | Not available | Available     |

| Table 1                                                                                |  |
|----------------------------------------------------------------------------------------|--|
| Availability of ITS interface depending on the type of card inserted in the tachograph |  |

#### Table 2

#### Assignment of the ITS connection depending on the type of card inserted in the tachograph

| Assignment of the ITS<br>Bluetooth® connection |               | Driver slot   |                  |                  |                   |                  |
|------------------------------------------------|---------------|---------------|------------------|------------------|-------------------|------------------|
|                                                |               | No card       | Driver card      | Control card     | Workshop card     | Company card     |
| Co - driverslot                                | No card       | Not available | Driver card      | Control card     | Workshop card     | Company card     |
|                                                | Driver card   | Driver card   | Driver card (**) | Control card     | Workshop card     | Company card     |
|                                                | Control card  | Control card  | Control card     | Control card (*) | Not available     | Not available    |
|                                                | Workshop card | Workshop card | Workshop card    | Not available    | Workshop card (*) | Not available    |
|                                                | Company card  | Company card  | Company card     | Not available    | Not available     | Company card (*) |

(\*) The ITS Bluetooth® connection shall be assigned to the tachograph card in the driver slot of the VU.

(\*\*) The user shall select the card to which the ITS Bluetooth® connection shall be assigned (inserted in the driver or in the co-driver slot).

- ITS\_15 If a tachograph card is withdrawn, then the VU shall terminate the ITS Bluetooth® connection which is assigned to this card.
- ITS\_16 The VU shall support the ITS connection to at least one ITS unit and may support connections to multiple ITS units at the same time.
- ITS\_17 The access rights to the data and services available via the ITS interface shall comply with requirements 12 and 13 of Annex IC, in addition to the driver consent specified in section 3.4 of this Appendix.

ITS\_14 After a successful ITS Bluetooth® pairing, the VU shall assign the ITS Bluetooth® connection to the specific inserted tachograph card according to Table 2:

#### 3.4. Data available and need of driver consent

- ITS\_18 All tachograph data available via the services referred to in point 3.3 shall be classified as either personal or not personal for the driver, co-driver or both.
- ITS\_19 At least the list of data classified as mandatory in section 4 shall be made available through the ITS interface.
- ITS\_20 The data in section 4 that are classified as 'personal' shall only be accessible upon driver consent, accepting therefore that the personal data can leave the vehicle network, except in the case set out in requirement ITS\_25, for which the driver consent is not needed.
- ITS\_21 Data additional to those gathered in point 4 and considered as mandatory may be made available through the ITS interface. Additional data which are not included in point 4 shall be classified as 'personal' or not 'personal' by the VU manufacturer, being the driver consent requested for those data that have been classified as personal, except in the case set out in requirement ITS\_25, for which the driver consent is not needed.
- ITS\_22 Upon insertion of a driver card which is unknown to the vehicle unit, the cardholder shall be prompted by the tachograph to enter the consent for transmission of personal data output through the ITS interface, in accordance with requirement 61 of Annex IC.
- ITS\_23 The consent status (enabled/disabled) shall be recorded in the data memory of the vehicle unit.
- ITS\_24 In case of multiple drivers, only the personal data related to the drivers who gave their consent shall be accessible through the ITS interface. For instance, in a crew situation, if only the driver has given his/her consent, personal data related to the co-driver shall not be accessible.
- ITS\_25 When the VU is in control, company or calibration modes, the access rights through the ITS interface shall be managed according to requirements 12 and 13 of Annex IC, hence the driver consent not being needed.

#### 4. LIST OF DATA AVAILABLE THROUGH THE ITS INTERFACE AND PERSONAL/NOT PERSONAL CLASSIFICATION

| Data name                          | Data format | Source      | Data classification (personal/ not<br>personal) |              | Consent for the avail-<br>ability of the data | Availability |
|------------------------------------|-------------|-------------|-------------------------------------------------|--------------|-----------------------------------------------|--------------|
|                                    |             |             | driver                                          | co-driver    |                                               |              |
| VehicleIdentification-<br>Number   | Appendix 8  | VU          | not personal                                    | not personal | no need of<br>consent                         | mandatory    |
| CalibrationDate                    | ISO 16844-7 | VU          | not personal                                    | not personal | no need of<br>consent                         | mandatory    |
| TachographVehi-<br>cleSpeed        | ISO 16844-7 | VU          | personal                                        | N/A          | driver consent                                | mandatory    |
| Data name                          | Data format | Source      | Data classification (personal/ not personal)    |              | Consent for the availability of the data      | Availability |
|                                    |             |             | driver                                          | co-driver    |                                               |              |
| Driver1WorkingState                | ISO 16844-7 | VU          | personal                                        | N/A          | driver consent                                | mandatory    |
| Driver2WorkingState                | ISO 16844-7 | VU          | N/A                                             | personal     | co-driver consent                             | mandatory    |
| DriveRecognize                     | ISO 16844-7 | VU          | not personal                                    | not personal | no need of consent                            | mandatory    |
| Driver1TimeRelated-States          | ISO 16844-7 | VU          | personal                                        | N/A          | driver consent                                | mandatory    |
| Driver2TimeRelated-States          | ISO 16844-7 | VU          | N/A                                             | personal     | co-driver consent                             | mandatory    |
| DriverCardDriver1                  | ISO 16844-7 | VU          | personal                                        | N/A          | driver consent                                | mandatory    |
| DriverCardDriver2                  | ISO 16844-7 | VU          | N/A                                             | personal     | co-driver consent                             | mandatory    |
| OverSpeed                          | ISO 16844-7 | VU          | personal                                        | N/A          | driver consent                                | mandatory    |
| TimeDate                           | Appendix 8  | VU          | not personal                                    | not personal | no need of consent                            | mandatory    |
| HighResolutionTotalVehicleDistance | ISO 16844-7 | VU          | not personal                                    | not personal | no need of consent                            | mandatory    |
| HighResolutionTrip-Distance        | ISO 16844-7 | VU          | not personal                                    | not personal | no need of consent                            | mandatory    |
| ServiceComponentIdentification     | ISO 16844-7 | VU          | not personal                                    | not personal | no need of consent                            | mandatory    |
| ServiceDelayCalendarTimeBased      | ISO 16844-7 | VU          | not personal                                    | not personal | no need of consent                            | mandatory    |
| Driver1Identification              | ISO 16844-7 | Driver Card | personal                                        | N/A          | driver consent                                | mandatory    |
| Driver2Identification              | ISO 16844-7 | Driver Card | N/A                                             | personal     | co-driver consent                             | mandatory    |

| M3 |  |
|----|--|
|----|--|

| Data name                                                   | Data format | Source              | Data classification (personal/ not<br>personal) |              | Consent for the avail-<br>ability of the data | Availability |
|-------------------------------------------------------------|-------------|---------------------|-------------------------------------------------|--------------|-----------------------------------------------|--------------|
|                                                             |             |                     | driver                                          | co-driver    |                                               |              |
| NextCalibrationDate                                         | Appendix 8  | VU                  | not personal                                    | not personal | no need of<br>consent                         | mandatory    |
| Driver1Continuous-<br>DrivingTime                           | ISO 16844-7 | VU                  | personal                                        | N/A          | driver consent                                | mandatory    |
| Driver2Continuous-<br>DrivingTime                           | ISO 16844-7 | VU                  | N/A                                             | personal     | co-driver consent                             | mandatory    |
| Driver1Cumulative-<br>BreakTime                             | ISO 16844-7 | VU                  | personal                                        | N/A          | driver consent                                | mandatory    |
| Driver2Cumulative-<br>BreakTime                             | ISO 16844-7 | VU                  | N/A                                             | personal     | co-driver consent                             | mandatory    |
| Driver1CurrentDur-<br>ationOfSelectedAc-<br>tivity          | ISO 16844-7 | VU                  | personal                                        | N/A          | driver consent                                | mandatory    |
| Driver2CurrentDur-<br>ationOfSelectedAc-<br>tivity          | ISO 16844-7 | VU                  | N/A                                             | personal     | co-driver consent                             | mandatory    |
| SpeedAuthorised                                             | Appendix 8  | VU                  | not personal                                    | not personal | no need of<br>consent                         | mandatory    |
| TachographCardSlot1                                         | ISO 16844-7 | VU                  | not personal                                    | N/A          | no need of<br>consent                         | mandatory    |
| TachographCardSlot2                                         | ISO 16844-7 | VU                  | N/A                                             | not personal | no need of<br>consent                         | mandatory    |
| Driver1Name                                                 | ISO 16844-7 | Driv-<br>er<br>Card | personal                                        | N/A          | driver consent                                | mandatory    |
| Driver2Name                                                 | ISO 16844-7 | Driv-<br>er<br>Card | N/A                                             | personal     | co-driver consent                             | mandatory    |
| OutOfScopeCondition                                         | ISO 16844-7 | VU                  | not personal                                    | not personal | no need of<br>consent                         | mandatory    |
| ModeOfOperation                                             | ISO 16844-7 | VU                  | not personal                                    | not personal | no need of<br>consent                         | mandatory    |
| Driver1Cumulated-<br>DrivingTimePreviou-<br>sAndCurrentWeek | ISO 16844-7 | VU                  | personal                                        | N/A          | driver consent                                | mandatory    |

| ▼M3 |
|-----|
|     |

| Data name                                                 | Data format | Source             | Data classification (personal/ not<br>personal) |              | Consent for the avail                    | Availability |
|-----------------------------------------------------------|-------------|--------------------|-------------------------------------------------|--------------|------------------------------------------|--------------|
|                                                           |             |                    | driver                                          | co-driver    | ability of the data                      |              |
| Driver2Cumulated<br>DrivingTimePreviou<br>sAndCurrentWeek | ISO 16844-7 | VU                 | N/A                                             | personal     | co-driver consent                        | mandatory    |
| EngineSpeed                                               | ISO 16844-7 | VU                 | personal                                        | N/A          | driver consent                           | optional     |
| RegisteringMem<br>berState                                | Appendix 8  | VU                 | not personal                                    | not personal | no need of<br>consent                    | mandatory    |
| VehicleRegistration<br>Number                             | Appendix 8  | VU                 | not personal                                    | not personal | no need of<br>consent                    | mandatory    |
| Driver1EndOfLast<br>DailyRestPeriod                       | ISO 16844-7 | VU                 | personal                                        | N/A          | driver consent                           | optional     |
| Driver2EndOfLast<br>DailyRestPeriod                       | ISO 16844-7 | VU                 | N/A                                             | personal     | co-driver consent                        | optional     |
| Driver1EndOfLast<br>WeeklyRestPeriod                      | ISO 16844-7 | VU                 | personal                                        | N/A          | driver consent                           | optional     |
| Driver2EndOfLast<br>WeeklyRestPeriod                      | ISO 16844-7 | VU                 | N/A                                             | personal     | co-driver consent                        | optional     |
| Driver1EndOfSecond<br>LastWeeklyRestPeriod                | ISO 16844-7 | VU                 | personal                                        | N/A          | driver consent                           | optional     |
| Driver2EndOfSecond<br>LastWeeklyRestPeriod                | ISO 16844-7 | VU                 | N/A                                             | Personal     | co-driver consent                        | optional     |
| Driver1TimeLastLoa<br>dUnloadOperation                    | ISO 16844-7 | VU                 | personal                                        | N/A          | driver consent                           | optional     |
| Driver2TimeLastLoa<br>dUnloadOperation                    | ISO 16844-7 | VU                 | N/A                                             | personal     | co-driver consent                        | optional     |
| Driver1CurrentDaily<br>DrivingTime                        | ISO 16844-7 | VU                 | personal                                        | N/A          | driver consent                           | optional     |
| Driver2CurrentDaily<br>DrivingTime                        | ISO 16844-7 | VU                 | N/A                                             | personal     | co-driver consent                        | optional     |
| Driver1CurrentWeekly<br>DrivingTime                       | ISO 16844-7 | VU                 | personal                                        | N/A          | driver consent                           | optional     |
| Driver2CurrentWeekly<br>DrivingTime                       | ISO 16844-7 | VU                 | N/A                                             | personal     | co-driver consent                        | optional     |
| Data name                                                 | Data format | Source             | Data classification (personal/ not<br>personal) |              | Consent for the avail                    | Availability |
|                                                           |             |                    | driver                                          | co-driver    | ability of the data                      |              |
| Driver1TimeLeftUntil<br>NewDailyRestPeriod                | ISO 16844-7 | VU                 | personal                                        | N/A          | driver consent                           | optional     |
| Driver2TimeLeftUntil<br>NewDailyRestPeriod                | ISO 16844-7 | VU                 | N/A                                             | personal     | co-driver consent                        | optional     |
| Driver1Card<br>ExpiryDate                                 | ISO 16844-7 | Driv<br>er<br>Card | personal                                        | N/A          | driver consent                           | optional     |
| Driver2Card<br>ExpiryDate                                 | ISO 16844-7 | Driv<br>er<br>Card | N/A                                             | personal     | co-driver consent                        | optional     |
| Driver1CardNextMan<br>datoryDownloadDate                  | ISO 16844-7 | VU                 | personal                                        | N/A          | driver consent                           | optional     |
| Driver2CardNextMan<br>datoryDownloadDate                  | ISO 16844-7 | VU                 | N/A                                             | personal     | co-driver consent                        | optional     |
| TachographNextMan<br>datoryDownloadDate                   | ISO 16844-7 | VU                 | not personal                                    | not personal | no need of<br>consent                    | optional     |
| Driver1TimeLeftUntil<br>NewWeeklyRestPeriod               | ISO 16844-7 | VU                 | personal                                        | N/A          | driver consent                           | optional     |
| Driver2TimeLeftUntil<br>NewWeeklyRestPeriod               | ISO 16844-7 | VU                 | N/A                                             | personal     | co-driver consent                        | optional     |
| Driver1Numbe<br>rOfTimes9hDailyDriv<br>ingTimesExceeded   | ISO 16844-7 | VU                 | personal                                        | N/A          | driver consent                           | optional     |
| Driver2Numbe<br>rOfTimes9hDailyDriv<br>ingTimesExceeded   | ISO 16844-7 | VU                 | N/A                                             | personal     | co-driver consent                        | optional     |
| Driver1CumulativeUn<br>interruptedRestTime                | ISO 16844-7 | VU                 | personal                                        | N/A          | driver consent                           | optional     |
| Driver2CumulativeUn<br>interruptedRestTime                | ISO 16844-7 | VU                 | N/A                                             | personal     | co-driver consent                        | optional     |
| Driver1Minimum<br>DailyRest                               | ISO 16844-7 | VU                 | personal                                        | N/A          | driver consent                           | optional     |
| Driver2Minimum<br>DailyRest                               | ISO 16844-7 | VU                 | N/A                                             | personal     | co-driver consent                        | optional     |
| Data name                                                 | Data format | Source             | Data classification (personal/ not<br>personal) |              | Consent for the availability of the data | Availability |
|                                                           |             |                    | driver                                          | co-driver    |                                          |              |
| Driver1Minimum-<br>WeeklyRest                             | ISO 16844-7 | VU                 | personal                                        | N/A          | driver consent                           | optional     |
| Driver2Minimum-<br>WeeklyRest                             | ISO 16844-7 | VU                 | N/A                                             | personal     | co-driver consent                        | optional     |
| Driver1Maximum-<br>DailyPeriod                            | ISO 16844-7 | VU                 | personal                                        | N/A          | driver consent                           | optional     |
| Driver2Maximum-<br>DailyPeriod                            | ISO 16844-7 | VU                 | N/A                                             | personal     | co-driver consent                        | optional     |
| Driver1Maximum-<br>DailyDrivingTime                       | ISO 16844-7 | VU                 | personal                                        | N/A          | driver consent                           | optional     |
| Driver2Maximum-<br>DailyDrivingTime                       | ISO 16844-7 | VU                 | N/A                                             | personal     | co-driver consent                        | optional     |
| Driver1NumberOfUse-<br>dReducedDailyRest-<br>Periods      | ISO 16844-7 | VU                 | personal                                        | N/A          | driver consent                           | optional     |
| Driver2NumberOfUse-<br>dReducedDailyRest-<br>Periods      | ISO 16844-7 | VU                 | N/A                                             | personal     | co-driver consent                        | optional     |
| Driver1RemainingCur-<br>rentDrivingTime                   | ISO 16844-7 | VU                 | personal                                        | N/A          | driver consent                           | optional     |
| Driver2RemainingCur-<br>rentDrivingTime                   | ISO 16844-7 | VU                 | N/A                                             | personal     | co-driver consent                        | optional     |
| VehiclePosition                                           | Appendix 8  | VU                 | personal                                        | personal     | driver and<br>co-driver consent          | mandatory    |
| ByDefaultLoadType                                         | Appendix 8  | VU                 | personal                                        | personal     | driver and<br>co-driver consent          | mandatory    |