#### *Appendix 11*

#### **COMMON SECURITY MECHANISMS**

#### TABLE OF CONTENTS

#### PREAMBLE

| PART A | FIRST-GENERATION TACHOGRAPH SYSTEM |
|--------|------------------------------------|
|--------|------------------------------------|

| 1.     | INTRODUCTION                                                                           |
|--------|----------------------------------------------------------------------------------------|
|        |                                                                                        |
| 6.2.   | Signature verification                                                                 |
| PART B | SECOND-GENERATION TACHOGRAPH SYSTEM                                                    |
| 7.     | INTRODUCTION                                                                           |
| 7.1.   | References                                                                             |
| 7.2.   | Notations and Abbreviations                                                            |
| 7.3.   | Definitions                                                                            |
| 8.     | CRYPTOGRAPHIC SYSTEMS AND ALGORITHMS                                                   |
| 8.1.   | Cryptographic Systems                                                                  |
| 8.2.   | Cryptographic Algorithms                                                               |
| 8.2.1  | Symmetric Algorithms                                                                   |
| 8.2.2  | Asymmetric Algorithms and Standardized Domain Parameters                               |
| 8.2.3  | Hashing algorithms                                                                     |
| 8.2.4  | Cipher Suites                                                                          |
| 9.     | KEYS AND CERTIFICATES                                                                  |
| 9.1.   | Asymmetric Key Pairs and Public Key Certificates                                       |
| 9.1.1  | General                                                                                |
| 9.1.2  | European Level                                                                         |
| 9.1.3  | Member State Level                                                                     |
| 9.1.4  | Equipment Level: Vehicle Units                                                         |
| 9.1.5  | Equipment Level: Tachograph Cards                                                      |
| 9.1.6  | Equipment Level: External GNSS Facilities                                              |
| 9.1.7  | Overview: Certificate Replacement                                                      |
| 9.2.   | Symmetric Keys                                                                         |
| 9.2.1  | Keys for Securing VU — Motion Sensor Communication                                     |
| 9.2.2  | Keys for Securing DSRC Communication                                                   |
| 9.3.   | Certificates                                                                           |
| 9.3.1  | General                                                                                |
| 9.3.2  | Certificate Content                                                                    |
| 9.3.3  | Requesting Certificates                                                                |
| 10.    | VU- CARD MUTUAL AUTHENTICATION AND SECURE MESSAGING                                    |
| 10.1.  | General                                                                                |
| 10.2.  | Mutual Certificate Chain Verification                                                  |
| 10.2.1 | Card Certificate Chain Verification by VU                                              |
|        |                                                                                        |
|        |                                                                                        |
| 10.2.2 | VU Certificate Chain Verification by Card                                              |
| 10.3.  | VU Authentication                                                                      |
| 10.4.  | Chip Authentication and Session Key Agreement                                          |
| 10.5.  | Secure Messaging                                                                       |
| 10.5.1 | General                                                                                |
| 10.5.2 | Secure Message Structure                                                               |
| 10.5.3 | Secure Messaging Session Abortion                                                      |
| 11.    | VU — EXTERNAL GNSS FACILITY COUPLING, MUTUAL<br>AUTHENTICATION AND SECURE MESSAGING    |
| 11.1.  | General                                                                                |
| 11.2.  | VU and External GNSS Facility Coupling                                                 |
| 11.3.  | Mutual Certificate Chain Verification                                                  |
| 11.3.1 | General                                                                                |
| 11.3.2 | During VU — EGF Coupling                                                               |
| 11.3.3 | During Normal Operation                                                                |
| 11.4.  | VU<br>Authentication,<br>Chip<br>Authentication<br>and<br>Session<br>Key<br>Agreement  |
| 11.5.  | Secure Messaging                                                                       |
| 12.    | VU — MOTION SENSOR PAIRING AND COMMUNICATION                                           |
| 12.1.  | General                                                                                |
| 12.2.  | VU — Motion Sensor Pairing Using Different Key Generations                             |
| 12.3.  | VU — Motion Sensor Pairing and Communication using AES                                 |
| 12.4.  | VU<br>—<br>Motion<br>Sensor<br>Pairing<br>For<br>Different<br>Equipment<br>Generations |
| 13.    | SECURITY FOR REMOTE COMMUNICATION OVER DSRC                                            |
| 13.1.  | General                                                                                |
| 13.2.  | Tachograph Payload Encryption and MAC Generation                                       |
| 13.3.  | Verification and Decryption of Tachograph Payload                                      |
| 14.    | SIGNING<br>DATA<br>DOWNLOADS<br>AND<br>VERIFYING<br>SIGNATURES                         |
| 14.1.  | General                                                                                |
| 14.2.  | Signature generation                                                                   |
|        |                                                                                        |

- 1.1. References
- 1.2. Notations and abbreviated terms
- 2. CRYPTOGRAPHIC SYSTEMS AND ALGORITHMS
- 2.1. Cryptographic systems
- 2.2. Cryptographic algorithms
- 2.2.1 RSA algorithm
- 2.2.2 Hash algorithm
- 2.2.3 Data Encryption Algorithm
- 3. KEYS AND CERTIFICATES
- 3.1. Keys generation and distribution
- 3.1.1 RSA Keys generation and distribution
- 3.1.2 RSA Test keys
- 3.1.3 Motion sensor keys
- 3.1.4 T-DES session keys generation and distribution
- 3.2. Keys
- 3.3. Certificates
- 3.3.1 Certificates content
- 3.3.2 Certificates issued
- 3.3.3 Certificate verification and unwrapping
- 4. MUTUAL AUTHENTICATION MECHANISM
- 5. VU-CARDS DATA TRANSFER CONFIDENTIALITY, INTEGRITY AND AUTHENTICATION MECHANISMS
- 5.1. Secure Messaging
- 5.2. Treatment of Secure Messaging errors
- 5.3. Algorithm to compute Cryptographic Checksums
- 5.4. Algorithm to compute cryptograms for confidentiality DOs
- 6. DATA DOWNLOAD DIGITAL SIGNATURE MECHANISMS
- 6.1. Signature generation

#### PREAMBLE

This Appendix specifies the security mechanisms ensuring

- mutual authentication between different components of the tachograph system.
- confidentiality, integrity, authenticity and/or non-repudiation of data transferred between different components of the tachograph system or downloaded to external storage media.

This Appendix consists of two parts. Part A defines the security mechanisms for the first-generation tachograph system (digital tachograph). Part B defines the security mechanisms for the second-generation tachograph system (smart tachograph).

The mechanisms specified in Part A of this Appendix shall apply if at least one of the components of the tachograph system involved in a mutual authentication and/or data transfer process is of the first generation.

The mechanisms specified in Part B of this Appendix shall apply if both components of the tachograph system involved in the mutual authentication and/or data transfer process are of the second generation.

Appendix 15 provides more information regarding the use of first generation components in combination with second-generation components.

#### PART A

#### **FIRST-GENERATION TACHOGRAPH SYSTEM**

#### 1. INTRODUCTION

#### 1.1. **References**

The following references are used in this Appendix:

- SHA-1 National Institute of Standards and Technology (NIST). *FIPS Publication 180-1: Secure Hash Standard*. April 1995.
- PKCS1 RSA Laboratories. PKCS # 1: *RSA Encryption Standard*. Version 2.0. October 1998.
- TDES National Institute of Standards and Technology (NIST). *FIPS Publication 46-3: Data Encryption Standard*. Draft 1999.
- TDES-OP ANSI X9.52, Triple Data Encryption Algorithm Modes of Operation. 1998.
- ISO/IEC 7816-4 Information Technology Identification cards Integrated circuit(s) cards with contacts — Part 4: Interindustry commands for interexchange. First edition: 1995 + Amendment 1: 1997.
- ISO/IEC 7816-6 Information Technology Identification cards Integrated circuit(s) cards with contacts — Part 6: Interindustry data elements. First edition: 1996 + Cor 1: 1998.
- ISO/IEC 7816-8 Information Technology Identification cards Integrated circuit(s) cards with contacts — Part 8: Security related interindustry commands. First edition 1999.
- ISO/IEC 9796-2 Information Technology Security techniques Digital signature schemes giving message recovery — Part 2: Mechanisms using a hash function. First edition: 1997.

|                                                                          | ISO/IEC 9798-3 Information Technology — Security techniques Entity authentication mechanisms — Part 3: Entity authentication using a public key algorithm. Second edition 1998.                                                                                                                                                                      |                                                                        |
|--------------------------------------------------------------------------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|------------------------------------------------------------------------|
|                                                                          | ISO 16844-3                                                                                                                                                                                                                                                                                                                                          | Road vehicles — Tachograph systems — Part 3: Motion sensor interface.  |
| 1.2.                                                                     | Notations and abbreviated terms                                                                                                                                                                                                                                                                                                                      |                                                                        |
| The following notations and abbreviated terms are used in this Appendix: |                                                                                                                                                                                                                                                                                                                                                      |                                                                        |
|                                                                          | (Ka, Kb, Kc)                                                                                                                                                                                                                                                                                                                                         | a key bundle for use by the Triple Data Encryption Algorithm,          |
|                                                                          | CA                                                                                                                                                                                                                                                                                                                                                   | Certification Authority,                                               |
|                                                                          | CAR                                                                                                                                                                                                                                                                                                                                                  | Certification Authority Reference,                                     |
|                                                                          | CC                                                                                                                                                                                                                                                                                                                                                   | Cryptographic Checksum,                                                |
|                                                                          | CG                                                                                                                                                                                                                                                                                                                                                   | Cryptogram,                                                            |
|                                                                          | CH                                                                                                                                                                                                                                                                                                                                                   | Command Header,                                                        |
|                                                                          | CHA                                                                                                                                                                                                                                                                                                                                                  | Certificate Holder Authorisation,                                      |
|                                                                          | CHR                                                                                                                                                                                                                                                                                                                                                  | Certificate Holder Reference,                                          |
|                                                                          | D()                                                                                                                                                                                                                                                                                                                                                  | Decryption with DES,                                                   |
|                                                                          | DE                                                                                                                                                                                                                                                                                                                                                   | Data Element,                                                          |
|                                                                          | DO                                                                                                                                                                                                                                                                                                                                                   | Data Object,                                                           |
|                                                                          | d                                                                                                                                                                                                                                                                                                                                                    | RSA private key, private exponent,                                     |
|                                                                          | e                                                                                                                                                                                                                                                                                                                                                    | RSA public key, public exponent,                                       |
|                                                                          | E()                                                                                                                                                                                                                                                                                                                                                  | Encryption with DES,                                                   |
|                                                                          | EQT                                                                                                                                                                                                                                                                                                                                                  | Equipment,                                                             |
|                                                                          | Hash()                                                                                                                                                                                                                                                                                                                                               | hash value, an output of Hash,                                         |
|                                                                          | Hash                                                                                                                                                                                                                                                                                                                                                 | hash function,                                                         |
|                                                                          | KID                                                                                                                                                                                                                                                                                                                                                  | Key Identifier,                                                        |
|                                                                          | Km                                                                                                                                                                                                                                                                                                                                                   | TDES key. Master Key defined in ISO 16844-3.                           |
|                                                                          | KmVU                                                                                                                                                                                                                                                                                                                                                 | TDES key inserted in vehicle units.                                    |
|                                                                          | KmWC                                                                                                                                                                                                                                                                                                                                                 | TDES key inserted in workshop cards.                                   |
|                                                                          | m                                                                                                                                                                                                                                                                                                                                                    | message representative, an integer between 0 and n-1,                  |
|                                                                          | n                                                                                                                                                                                                                                                                                                                                                    | RSA keys, modulus,                                                     |
|                                                                          | PB                                                                                                                                                                                                                                                                                                                                                   | Padding Bytes,                                                         |
|                                                                          | PI                                                                                                                                                                                                                                                                                                                                                   | Padding Indicator byte (for use in Cryptogram for confidentiality DO), |
|                                                                          | PV                                                                                                                                                                                                                                                                                                                                                   | Plain Value,                                                           |
|                                                                          | s                                                                                                                                                                                                                                                                                                                                                    | signature representative, an integer between 0 and n-1,                |
|                                                                          | SSC                                                                                                                                                                                                                                                                                                                                                  | Send Sequence Counter,                                                 |
|                                                                          | SM                                                                                                                                                                                                                                                                                                                                                   | Secure Messaging,                                                      |
| TDEA                                                                     | Triple Data Encryption Algorithm,                                                                                                                                                                                                                                                                                                                    |                                                                        |
| TLV                                                                      | Tag Length Value,                                                                                                                                                                                                                                                                                                                                    |                                                                        |
| VU                                                                       | Vehicle Unit,                                                                                                                                                                                                                                                                                                                                        |                                                                        |
| X.C                                                                      | the certificate of user X issued by a certification<br>authority,                                                                                                                                                                                                                                                                                    |                                                                        |
| X.CA                                                                     | a certification authority of user X,                                                                                                                                                                                                                                                                                                                 |                                                                        |
| X.CA.PK o X.C                                                            | the operation of unwrapping a certificate to extract a<br>public key. It is an infix operator, whose left operand<br>is the public key of a certification authority, and<br>whose right operand is the certificate issued by that<br>certification authority. The outcome is the public key<br>of the user X whose certificate is the right operand, |                                                                        |
| X.PK                                                                     | RSA public key of a user X,                                                                                                                                                                                                                                                                                                                          |                                                                        |
| X.PK[I]                                                                  | RSA encipherment of some information I, using the<br>public key of user X,                                                                                                                                                                                                                                                                           |                                                                        |
| X.SK                                                                     | RSA private key of a user X,                                                                                                                                                                                                                                                                                                                         |                                                                        |
| X.SK[I]                                                                  | RSA encipherment of some information I, using the<br>private key of user X,                                                                                                                                                                                                                                                                          |                                                                        |
| 'xx'                                                                     | an Hexadecimal value,                                                                                                                                                                                                                                                                                                                                |                                                                        |
|                                                                          | concatenation operator.                                                                                                                                                                                                                                                                                                                              |                                                                        |

TCBC TDEA Cipher Block Chaining Mode of Operation

#### 2. CRYPTOGRAPHIC SYSTEMS AND ALGORITHMS

#### 2.1. **Cryptographic systems**

- CSM\_001 Vehicle units and tachograph cards shall use a classical RSA public-key cryptographic system to provide the following security mechanisms:
  - authentication between vehicle units and cards,
  - transport of Triple-DES session keys between vehicle units and tachograph cards,
  - digital signature of data downloaded from vehicle units or tachograph cards to external media.
- CSM\_002 Vehicle units and tachograph cards shall use a Triple DES symmetric cryptographic system to provide a mechanism for data integrity during user data exchange between vehicle units and tachograph cards, and to provide, where applicable, confidentiality of data exchange between vehicle units and tachograph cards.

### 2.2. **Cryptographic algorithms**

2.2.1 *RSA algorithm*

CSM\_003 The RSA algorithm is fully defined by the following relations:

X.SK[*m*] = *s* = *md* mod *n* X.PK[*s*] = *m* = *se* mod *n*

A more comprehensive description of the RSA function can be found in reference [PKCS1]. Public exponent, e, for RSA calculations is an integer between 3 and n-1 satisfying gcd(e, lcm(p-1, q-1))=1.

2.2.2 *Hash algorithm*

- 2.2.3 *Data Encryption Algorithm*
  - CSM\_005 DES based algorithms shall be used in Cipher Block Chaining mode of operation.

#### 3. KEYS AND CERTIFICATES

#### 3.1. **Keys generation and distribution**

3.1.1 *RSA Keys generation and distribution*

CSM\_006 RSA keys shall be generated through three functional hierarchical levels:

— European level,

— Member State level,

- Equipment level.
- CSM\_007 At European level, a single European key pair (EUR.SK and EUR.PK) shall be generated. The European private key shall be used to certify the Member States public keys. Records of all certified keys shall be kept. These tasks shall be handled by a European Certification Authority, under the authority and responsibility of the European Commission.
- CSM\_008 At Member State level, a Member State key pair (MS.SK and MS.PK) shall be generated. Member States public keys shall be certified by the European Certification Authority. The Member State private key shall be used to certify public keys to be inserted in equipment (vehicle unit or tachograph card). Records of all certified public keys shall be kept with the identification of the equipment to which it is intended. These tasks shall be handled by a Member State Certification Authority. A Member State may regularly change its key pair.
- CSM\_009 At equipment level, one single key pair (EQT.SK and EQT.PK) shall be generated and inserted in each equipment. Equipment public keys shall be certified by a Member State Certification Authority. These tasks may be handled by equipment manufacturers, equipment personalisers or Member State authorities. This key pair is used for authentication, digital signature and encipherement services
- CSM\_010 Private keys confidentiality shall be maintained during generation, transport (if any) and storage.

CSM\_004 The digital signature mechanisms shall use the SHA-1 hash algorithm as defined in reference [SHA-1].

![](_page_7_Figure_1.jpeg)

The following picture summarises the data flow of this process:

### 3.1.2 *RSA Test keys*

CSM\_011 For the purpose of equipment testing (including interoperability tests) the European Certification Authority shall generate a different single European test key pair and at least two Member State test key pairs, the public keys of which shall be certified with the European private test key. Manufacturers shall insert, in equipment undergoing type approval tests, test keys certified by one of these Member State test keys.

#### 3.1.3 *Motion sensor keys*

The confidentiality of the three Triple DES keys described below shall be appropriately maintained during generation, transport (if any) and storage.

In order to support tachograph components compliant with ISO 16844, the European Certification Authority and the Member State Certification Authorities shall, in addition, ensure the following:

CSM\_036 The European Certification authority shall generate KmVU and KmWC, two independent and unique Triple DES keys, and generate Km as: Km = KmVU XOR KmWC. The European Certification Authority shall forward these keys, under appropriately secured procedures, to Member States Certification Authorities at their request.

CSM\_037 Member States Certification Authorities shall:

- use Km to encrypt motion sensor data requested by motion sensor manufacturers (data to be encrypted with Km is defined in ISO 16844-3),
- forward KmVU to vehicle unit manufacturers, under appropriately secured procedures, for insertion in vehicle units,
- ensure that KmWC will be inserted in all workshop cards ( in elementary file) during card personalisation.

#### 3.1.4 *T-DES session keys generation and distribution*

- CSM\_012 Vehicle units and tachograph cards shall, as a part of the mutual authentication process, generate and exchange necessary data to elaborate a common Triple DES session key. This exchange of data shall be protected for confidentiality through an RSA crypt-mechanism.
- CSM\_013 This key shall be used for all subsequent cryptographic operations using secure messaging. Its validity shall expire at the end of the session (withdrawal of the card or reset of the card) and/or after 240 use (one use of the key = one command using secure messaging sent to the card and associated response).

#### 3.2. **Keys**

- CSM\_014 RSA keys shall have (whatever the level) the following lengths: modulus *n* 1 024 bits, public exponent *e* 64 bits maximum, private exponent *d* 1 024 bits.
- CSM\_015 Triple DES keys shall have the form (Ka, Kb, Ka) where Ka and Kb are independent 64 bits long keys. No parity error detecting bits shall be set.

### 3.3. **Certificates**

CSM\_016 RSA Public key certificates shall be 'non self-descriptive' 'Card Verifiable' certificates (Ref.: ISO/IEC 7816-8)

### 3.3.1 *Certificates content*

CSM\_017 RSA Public key certificates are built with the following data in the following order:

| Data | Format          | Bytes | Obs                                                       |
|------|-----------------|-------|-----------------------------------------------------------|
| CPI  | INTEGER         | 1     | Certificate Profile Identifier<br>('01' for this version) |
| CAR  | OCTET<br>STRING | 8     | Certification Authority<br>Reference                      |
| CHA  | OCTET<br>STRING | 7     | Certificate Holder Author-<br>isation                     |

| Data | Format          | Bytes | Obs                                                                   |
|------|-----------------|-------|-----------------------------------------------------------------------|
| EOV  | TimeReal        | 4     | Certificate end of validity.<br>Optional, 'FF' padded if<br>not used. |
| CHR  | OCTET<br>STRING | 8     | Certificate Holder Reference                                          |
| n    | OCTET<br>STRING | 128   | Public key (modulus)                                                  |
| e    | OCTET<br>STRING | 8     | Public Key (public exponent)                                          |
|      |                 | 164   |                                                                       |

*Notes:*

The headerlist associated with this certificate content is as follows:

| '4D'                    | '16'                  | '5F<br>29' | '01'       | '42'    | '08'       | '5F<br>4B' | '07'       | '5F<br>24' | '04'       | '5F<br>20' | '08'       | '7F<br>49'                   | '05'                     | '81'        | '81<br>80'     | '82'                | '08'                   |
|-------------------------|-----------------------|------------|------------|---------|------------|------------|------------|------------|------------|------------|------------|------------------------------|--------------------------|-------------|----------------|---------------------|------------------------|
| Extended Headerlist Tag | Length of header list | CPI Tag    | CPI Length | CAR Tag | CAR Length | CHA Tag    | CHA Length | EOV Tag    | EOV Length | CHR Tag    | CHR Length | Public Key Tag (Constructed) | Length of subsequent DOs | modulus Tag | modulus length | public exponent Tag | public exponent length |

- 2. The 'Certification Authority Reference' (CAR) has the purpose of identifying the certificate issuing CA, in such a way that the Data Element can be used at the same time as an Authority Key Identifier to reference the Public Key of the Certification Authority (for coding, see Key Identifier below).
- 3. The 'Certificate Holder Authorisation' (CHA) is used to identify the rights of the certificate holder. It consists of the Tachograph Application ID and of the type of equipment to which the certificate is intended (according to data element, '00' for a Member State).

<sup>1.</sup> The 'Certificate Profile Identifier' (CPI) delineates the exact structure of an authentication certificate. It can be used as an equipment internal identifier of a relevant headerlist which describes the concatenation of Data Elements within the certificate.

- 4. The 'Certificate Holder Reference' (CHR) has the purpose of identifying uniquely the certificate holder, in such a way that the Data Element can be used at the same time as a Subject Key Identifier to reference the Public Key of the certificate holder.
- 5. Key Identifiers uniquely identify certificate holder or certification authorities. They are coded as follows:
  - 5.1 Equipment (VU or Card):

| Data   | Equipment<br>serial<br>number | Date                | Type                     | Manufacturer         |
|--------|-------------------------------|---------------------|--------------------------|----------------------|
| Length | 4 Bytes                       | 2 Bytes             | 1 Byte                   | 1 Byte               |
| Value  | Integer                       | mm yy BCD<br>coding | Manufacturer<br>specific | Manufacturer<br>code |

In the case of a VU, the manufacturer, when requesting certificates, may or may not know the identification of the equipment in which the keys will be inserted.

In the first case, the manufacturer will send the equipment identification with the public key to its Member State authority for certification. The certificate will then contain the equipment identification, and the manufacturer must ensure that keys and certificate are inserted in the intended equipment. The Key identifier has the form shown above.

In the later case, the manufacturer must uniquely identify each certificate request and send this identification with the public key to its Member State authority for certification. The certificate will contain the request identification. The manufacturer must feed back its Member State authority with the assignment of key to equipment (i.e. certificate request identification, equipment identification) after key installation in the equipment. The key identifier has the following form:

| Data   | Certificate<br>request<br>serial<br>number | Date             | Type   | Manufacturer         |
|--------|--------------------------------------------|------------------|--------|----------------------|
| Length | 4 Bytes                                    | 2 Bytes          | 1 Byte | 1 Byte               |
| Value  | Integer                                    | mm yy BCD coding | 'FF'   | Manufacturer<br>code |

5.2 Certification Authority:

| Data   | Authority Identifi<br>cation                                                | Key serial<br>number | Additional<br>info                                              | Identifier |
|--------|-----------------------------------------------------------------------------|----------------------|-----------------------------------------------------------------|------------|
| Length | 4 Bytes                                                                     | 1 Byte               | 2 Bytes                                                         | 1 Byte     |
| Value  | 1 Byte nation<br>numerical code<br>3 Bytes nation<br>alphanumerical<br>code | Integer              | additional<br>coding<br>(CA specific)<br>'FF FF' if not<br>used | '01'       |

The key serial number is used to distinguish the different keys of a Member State, in the case the key is changed.

6. Certificate verifiers shall implicitly know that the public key certified is an RSA key relevant to authentication, digital signature verification and encipherement for confidentiality services (the certificate contains no Object Identifier to specify it).

#### 3.3.2 *Certificates issued*

CSM\_018 The certificate issued is a digital signature with partial recovery of the certificate content in accordance with ISO/IEC 9796-2 (except for its annex A4), with the 'Certification Authority Reference' appended.

X.C = X.CA.SK['6A' || Cr || *Hash* (Cc) || 'BC'] || Cn || X.CAR

With certificate content = Cc = Cr || Cn 106 bytes 58 bytes

*Notes:*

- 1. This certificate is 194 bytes long.
- 2. CAR, being hidden by the signature, is also appended to the signature, such that the Public Key of the Certification Authority may be selected for the verification of the certificate.
- 3. The certificate verifier shall implicitly know the algorithm used by the Certification Authority to sign the certificate.
- 4. The headerlist associated with this issued certificate is as follows:

| '7F 21'                          | '09'                     | '5F 37'       | '81 80'          | '5F 38'       | '3A'             | '42'    | '08'       |
|----------------------------------|--------------------------|---------------|------------------|---------------|------------------|---------|------------|
| CV Certificate Tag (Constructed) | Length of subsequent DOs | Signature Tag | Signature Length | Remainder Tag | Remainder Length | CAR Tag | CAR Length |
|                                  |                          |               |                  |               |                  |         |            |

#### 3.3.3 *Certificate verification and unwrapping*

Certificate verification and unwrapping consists in verifying the signature in accordance with ISO/IEC 9796-2, retrieving the certificate content and the public key contained: X.PK = X.CA.PK <sup>o</sup> X.C, and verifying the validity of the certificate.

CSM\_019 It involves the following steps:

Verify signature and retrieve content:

| — from X.C retrieve Sign, Cn' and X.C =<br>CAR':                                                              | Sign      |  | Cn'       |  | CAR'       |
|---------------------------------------------------------------------------------------------------------------|-----------|--|-----------|--|------------|
|                                                                                                               | 128 Bytes |  | 58 Bytes  |  | 8 Bytes    |
| — from CAR' select appropriate Certification Authority<br>Public Key (if not done before through other means) |           |  |           |  |            |
| — open Sign with CA Public Key: Sr'= X.CA.PK [Sign],                                                          |           |  |           |  |            |
| — check Sr' starts with '6A' and ends with 'BC'                                                               |           |  |           |  |            |
| — compute Cr' and H' from: Sr' =                                                                              | '6A'      |  | Cr'       |  | H'    'BC' |
|                                                                                                               |           |  | 106 Bytes |  | 20 Bytes   |
| — Recover certificate content C' = Cr'    Cn',                                                                |           |  |           |  |            |
| — check Hash (C') = H'                                                                                        |           |  |           |  |            |
| If the checks are OK the certificate is a genuine one, its<br>content is C'.                                  |           |  |           |  |            |
| Verify validity. From C':                                                                                     |           |  |           |  |            |

— if applicable, check End of validity date,

Retrieve and store public key, Key Identifier, Certificate Holder Authorisation and Certificate End of Validity from C':

— X.PK = *n* || *e*

— X.KID = CHR

— X.CHA = CHA

— X.EOV = EOV

### 4. MUTUAL AUTHENTICATION MECHANISM

Mutual authentication between cards and VUs is based on the following principle:

Each party shall demonstrate to the other that it owns a valid key pair, the public key of which has been certified by a Member State certification authority, itself being certified by the European certification authority.

Demonstration is made by signing with the private key a random number sent by the other party, who must recover the random number sent when verifying this signature.

The mechanism is triggered at card insertion by the VU. It starts with the exchange of certificates and unwrapping of public keys, and ends with the setting of a session key.

CSM\_020 The following protocol shall be used (arrows indicate commands and data exchanged (see Appendix 2)):

![](_page_13_Figure_2.jpeg)

▼<u>B</u>

![](_page_14_Figure_2.jpeg)

5. VU-CARDS DATA TRANSFER CONFIDENTIALITY, INTEGRITY AND AUTHENTICATION MECHANISMS

#### 5.1. **Secure Messaging**

- CSM\_021 VU-Cards data transfers integrity shall be protected through Secure Messaging in accordance with references [ISO/IEC 7816-4] and [ISO/IEC 7816-8].
- CSM\_022 When data need to be protected during transfer, a Cryptographic Checksum Data Object shall be appended to the Data Objects sent within the command or the response. The Cryptographic Checksum shall be verified by the receiver.
- CSM\_023 The cryptographic checksum of data sent within a command shall integrate the command header, and all data objects sent (=>CLA = '0C', and all data objects shall be encapsulated with tags in which b1=1).
- CSM\_024 The response status-information bytes shall be protected by a cryptographic checksum when the response contains no data field.

CSM\_025 Cryptographic checksums shall be 4 Bytes long.

The structure of commands and responses when using secure messaging is therefore the following:

The DOs used are a partial set of the Secure Messaging DOs described in ISO/IEC 7816-4:

| Tag  | Mnemonic | Meaning                                                                 |
|------|----------|-------------------------------------------------------------------------|
| '81' | TPV      | Plain Value not BER-TLV coded data (to be protected by CC)              |
| '97' | TLE      | Value of Le in the unsecured command (to be protected by CC)            |
| '99' | TSW      | Status-Info (to be protected by CC)                                     |
| '8E' | TCC      | Cryptographic Checksum                                                  |
| '87' | TPI CG   | Padding Indicator Byte    Cryptogram (Plain Value not coded in BER-TLV) |

Given an unsecured command response pair:

| Command header |                  |     |    | Command body                 |              |            |
|----------------|------------------|-----|----|------------------------------|--------------|------------|
| CLA            | INS              | P1  | P2 | [Lc field]                   | [Data field] | [Le field] |
| four bytes     |                  |     |    | L bytes, denoted as B1 to BL |              |            |
| Response body  | Response trailer |     |    |                              |              |            |
| [Data field]   | SW1              | SW2 |    |                              |              |            |
| Lr data bytes  | two bytes        |     |    |                              |              |            |

The corresponding secured command response pair is:

Secured command:

|      | Command header (CH)         |     |      | Command body |                |                  |      |     |      |      |     |    |                   |
|------|-----------------------------|-----|------|--------------|----------------|------------------|------|-----|------|------|-----|----|-------------------|
|      | CLA                         | INS | P1   | P2           | [New Lc field] | [New Data field] |      |     |      |      |     |    | [New Le<br>field] |
| 'OC' | Length of New<br>Data field |     |      | TPV          | LPV            | PV               | TLE  | LLE | Le   | TCC  | LCC | CC | '00'              |
|      |                             |     | '81' | Lc           | Data field     | '97'             | '01' | Le  | '8E' | '04' | CC  |    |                   |

Data to be integrated in checksum = CH || PB || TPV || LPV || PV || TLE || LLE || Le || PB

PB = Padding Bytes (80 .. 00) in accordance with ISO-IEC 7816-4 and ISO 9797 method 2.

DOs PV and LE are present only when there is some corresponding data in the unsecured command.

Secured response:

1. Case where response data field is not empty and needs not to be protected for confidentiality:

| Response body    |     |            | Response trailer |      |    |
|------------------|-----|------------|------------------|------|----|
| [New Data field] |     |            | new SW1 SW2      |      |    |
| TPV              | LPV | PV         | TCC              | LCC  | CC |
| '81'             | Lr  | Data field | '8E'             | '04' | CC |

Data to be integrated in checksum = TPV || LPV || PV || PB

<sup>2.</sup> Case where response data field is not empty and needs to be protected for confidentiality:

| Response body    |           |          |      |      |    | Response trailer |
|------------------|-----------|----------|------|------|----|------------------|
| [New Data field] |           |          |      |      |    | new SW1 SW2      |
| TPI CG           | LPI<br>CG | PI CG    | TCC  | LCC  | CC |                  |
| '87'             |           | PI    CG | '8E' | '04' | CC |                  |

Data to be carried by CG: non BER-TLV coded data and padding bytes.

Data to be integrated in checksum = TPI CG || LPI CG || PI CG || PB

3. Case where response data field is empty:

| Response body    |      |             |      |      |    | Response trailer |
|------------------|------|-------------|------|------|----|------------------|
| [New Data field] |      |             |      |      |    | new SW1 SW2      |
| TSW              | LSW  | SW          | TCC  | LCC  | CC |                  |
| '99'             | '02' | New SW1 SW2 | '8E' | '04' | CC |                  |

Data to be integrated in checksum = TSW || LSW || SW || PB

#### 5.2. **Treatment of Secure Messaging errors**

- CSM\_026 When the tachograph card recognises an SM error while interpreting a command, then the status bytes must be returned without SM. In accordance with ISO/IEC 7816-4, the following status bytes are defined to indicate SM errors:
  - '66 88': Verification of Cryptographic Checksum failed,
  - '69 87': Expected SM Data Objects missing,
  - '69 88': SM Data Objects incorrect.
- CSM\_027 When the tachograph card returns status bytes without SM DOs or with an erroneous SM DO, the session must be aborted by the VU.

#### 5.3. **Algorithm to compute Cryptographic Checksums**

- CSM\_028 Cryptographic checksums are built using a retail MACs in accordance with ANSI X9.19 with DES:
  - Initial stage: The initial check block y0 is E(Ka, SSC).
  - Sequential stage: The check blocks y1, .., yn are calculated using Ka.
  - Final stage: The cryptographic checksum is calculated from the last check block yn as follows: E(Ka, D(Kb, yn)).

where E() means encryption with DES, and D() means decryption with DES.

The four most significant bytes of the cryptographic checksum are transferred

CSM\_029 The Send Sequence Counter (SSC) shall be initiated during key agreement procedure to:

> Initial SSC: Rnd3 (4 least significant bytes) || Rnd1 (4 least significant bytes).

CSM\_030 The Send Sequence Counter shall be increased by 1 each time before a MAC is calculated (i.e. the SSC for the first command is Initial SSC + 1, the SSC for the first response is Initial SSC + 2).

The following figure shows the calculation of the retail MAC:

![](_page_18_Figure_7.jpeg)

### 5.4. **Algorithm to compute cryptograms for confidentiality DOs**

CSM\_031 Cryptograms are computed using TDEA in TCBC mode of operation in accordance with references [TDES] and [TDES-OP] and with the Null vector as Initial Value block.

The following figure shows the application of keys in TDES:

![](_page_18_Figure_11.jpeg)

#### 6. DATA DOWNLOAD DIGITAL SIGNATURE MECHANISMS

- CSM\_032 The Intelligent Dedicated Equipment (IDE) stores data received from an equipment (VU or card) during one download session within one physical data file. This file must contain the certificates MSi.C and EQT.C. The file contains digital signatures of data blocks as specified in Appendix 7 Data Downloading Protocols.
- CSM\_033 Digital signatures of downloaded data shall use a digital signature scheme with appendix such, that downloaded data may be read without any decipherment if desired.

#### 6.1. **Signature generation**

CSM\_034 Data signature generation by the equipment shall follow the signature scheme with appendix defined in reference [PKCS1] with the SHA-1 hash function:

> Signature = EQT.SK['00' || '01' || *PS* || '00' || DER(SHA-1 (Data))]

> *PS* = Padding string of octets with value 'FF' such that length is 128.

> DER(SHA-1(*M*)) is the encoding of the algorithm ID for the hash function and the hash value into an ASN.1 value of type DigestInfo (distinguished encoding rules):

'30'||'21'||'30'||'09'||'06'||'05'||'2B'||'0E'||'03'||'02'||'1A'|| '05'||'00'||'04'||'14'||Hash Value.

#### 6.2. **Signature verification**

CSM\_035 Data signature verification on downloaded data shall follow the signature scheme with appendix defined in reference [PKCS1] with the SHA-1 hash function.

> The European public key EUR.PK needs to be known independently (and trusted) by the verifier.

> The following table illustrates the protocol an IDE carrying a Control card can follow to verify the integrity of data downloaded and stored on the ESM (External Storage media). The control card is used to perform the decipherement of digital signatures. This function may in this case not be implemented in the IDE.

> The equipment that has downloaded and signed the data to be analysed is denoted EQT.

![](_page_20_Figure_1.jpeg)

#### PART B

### **SECOND-GENERATION TACHOGRAPH SYSTEM**

### 7. INTRODUCTION

### 7.1. **References**

The following references are used in this part of this Appendix.

| AES | National Institute of Standards and Technology (NIST),<br>FIPS PUB 197: Advanced Encryption Standard (AES),<br>November 26, 2001 |
|-----|----------------------------------------------------------------------------------------------------------------------------------|
|-----|----------------------------------------------------------------------------------------------------------------------------------|

- DSS National Institute of Standards and Technology (NIST), FIPS PUB 186-4: Digital Signature Standard (DSS), July 2013
- ISO 7816-4 ISO/IEC 7816-4, Identification cards Integrated circuit cards — Part 4: Organization, security and commands for interchange. Third edition 2013-04-15
- ISO 7816-8 ISO/IEC 7816-8, Identification cards Integrated circuit cards — Part 8: Commands for security operations. Second edition 2004-06-01

- ISO 8825-1 ISO/IEC 8825-1, Information technology ASN.1 encoding rules: Specification of Basic Encoding Rules (BER), Canonical Encoding Rules (CER) and Distinguished Encoding Rules (DER). Fourth edition, 2008-12-15
- ISO 9797-1 ISO/IEC 9797-1, Information technology Security techniques — Message Authentication Codes (MACs) — Part 1: Mechanisms using a block cipher. Second edition, 2011-03-01
- ISO 10116 ISO/IEC 10116, Information technology Security techniques — Modes of operation of an *n*-bit block cipher. Third edition, 2006-02-01
- ISO 16844-3 ISO/IEC 16844-3, Road vehicles Tachograph systems — Part 3: Motion sensor interface. First edition 2004, including Technical Corrigendum 1 2006
- RFC 5480 Elliptic Curve Cryptography Subject Public Key Information, March 2009
- RFC 5639 Elliptic Curve Cryptography (ECC) Brainpool Standard Curves and Curve Generation, 2010
- RFC 5869 HMAC-based Extract-and-Expand Key Derivation Function (HKDF), May 2010
- SHS National Institute of Standards and Technology (NIST), FIPS PUB 180-4: Secure Hash Standard, March 2012
- SP 800-38B National Institute of Standards and Technology (NIST), Special Publication 800-38B: Recommendation for Block Cipher Modes of Operation: The CMAC Mode for Authentication, 2005
- TR-03111 BSI Technical Guideline TR-03111, Elliptic Curve Cryptography, version 2.00, 2012-06-28

### 7.2. **Notations and Abbreviations**

The following notations and abbreviated terms are used in this Appendix:

- AES Advanced Encryption Standard
- CA Certificate Authority
- CAR Certificate Authority Reference
- CBC Cipher Block Chaining (mode of operation)
- CH Command Header
- CHA Certificate Holder Authorisation
- CHR Certificate Holder Reference
- CV Constant Vector
- DER Distinguished Encoding Rules
- DO Data Object
- DSRC Dedicated Short Range Communication
- ECC Elliptic Curve Cryptography
- ECDSA Elliptic Curve Digital Signature Algorithm
- ECDH Elliptic Curve Diffie-Hellman (key agreement algorithm)
- EGF External GNSS Facility
- EQT Equipment

| IDE     | Intelligent Dedicated Equipment                                                                                                       |  |
|---------|---------------------------------------------------------------------------------------------------------------------------------------|--|
| KM      | Motion Sensor Master Key, allowing the pairing of a<br>Vehicle Unit to a Motion Sensor                                                |  |
| KM-VU   | Key inserted in vehicle units, allowing a VU to derive the<br>Motion Sensor Master Key if a workshop card is inserted<br>into the VU  |  |
| KM-WC   | Key inserted in workshop cards, allowing a VU to derive<br>the Motion Sensor Master Key if a workshop card is<br>inserted into the VU |  |
| MAC     | Message Authentication Code                                                                                                           |  |
| MoS     | Motion Sensor                                                                                                                         |  |
| MSB     | Most Significant Bit                                                                                                                  |  |
| PKI     | Public Key Infrastructure                                                                                                             |  |
| RCF     | Remote Communication Facility                                                                                                         |  |
| SSC     | Send Sequence Counter                                                                                                                 |  |
| SM      | Secure Messaging                                                                                                                      |  |
| TDES    | Triple Data Encryption Standard                                                                                                       |  |
| TLV     | Tag Length Value                                                                                                                      |  |
| VU      | Vehicle Unit                                                                                                                          |  |
| X.C     | the public key certificate of user X                                                                                                  |  |
| X.CA    | the certificate authority that issued the certificate of user X                                                                       |  |
| X.CAR   | the certificate authority reference mentioned in the certifi-<br>cate of user X                                                       |  |
| X.CHR   | the certificate holder reference mentioned in the certificate<br>of user X                                                            |  |
| X.PK    | public key of user X                                                                                                                  |  |
| X.SK    | private key of user X                                                                                                                 |  |
| X.PKeph | ephemeral public key of user X                                                                                                        |  |
| X.SKeph | ephemeral private key of user X                                                                                                       |  |
| 'xx'    | a hexadecimal value                                                                                                                   |  |
|         | concatenation operator                                                                                                                |  |

#### 7.3. **Definitions**

The definitions of terms used in this Appendix are included in section I of Annex 1C.

#### 8. CRYPTOGRAPHIC SYSTEMS AND ALGORITHMS

### 8.1. **Cryptographic Systems**

- CSM\_38 Vehicle units and tachograph cards shall use an elliptic curve-based public-key cryptographic system to provide the following security services:
  - mutual authentication between a vehicle unit and a card,

- agreement of AES session keys between a vehicle unit and a card,
- ensuring the authenticity, integrity and non-repudiation of data downloaded from vehicle units or tachograph cards to external media.
- CSM\_39 Vehicle units and external GNSS facilities shall use an elliptic curve-based public-key cryptographic system to provide the following security services:
  - coupling of a vehicle unit and an external GNSS facility,
  - mutual authentication between a vehicle unit and an external GNSS facility,
  - agreement of an AES session key between a vehicle unit and an external GNSS facility.
- CSM\_40 Vehicle units and tachograph cards shall use an AES-based symmetric cryptographic system to provide the following security services:
  - ensuring authenticity and integrity of data exchanged between a vehicle unit and a tachograph card,
  - where applicable, ensuring confidentiality of data exchanged between a vehicle unit and a tachograph card.
- CSM\_41 Vehicle units and external GNSS facilities shall use an AES-based symmetric cryptographic system to provide the following security services:
  - ensuring authenticity and integrity of data exchanged between a vehicle unit and an external GNSS facility.
- CSM\_42 Vehicle units and motion sensors shall use an AES-based symmetric cryptographic system to provide the following security services:
  - pairing of a vehicle unit and a motion sensor,
  - mutual authentication between a vehicle unit and a motion sensor,
  - ensuring confidentiality of data exchanged between a vehicle unit and a motion sensor.
- CSM\_43 Vehicle units and control cards shall use an AES-based symmetric cryptographic system to provide the following security services on the remote communication interface:
  - ensuring confidentiality, authenticity and integrity of data transmitted from a vehicle unit to a control card.

### *Notes:*

— Properly speaking, data is transmitted from a vehicle unit to a remote interrogator under the control of a control officer, using a remote communication facility that may be internal or external to the VU, see Appendix 14. However, the remote interrogator sends

the received data to a control card for decryption and validation of authenticity. From a security point of view, the remote communication facility and the remote interrogator are fully transparent.

— A workshop card offers the same security services for the DSRC interface as a control card does. This allows a workshop to validate the proper functioning of the remote communication interface of a VU, including security. Please refer to section 9.2.2 for more information.

### 8.2. **Cryptographic Algorithms**

8.2.1 *Symmetric Algorithms*

- CSM\_44 Vehicle units, tachograph cards, motion sensors and external GNSS facilities shall support the AES algorithm as defined in [AES], with key lengths of 128, 192 and 256 bits.
- 8.2.2 *Asymmetric Algorithms and Standardized Domain Parameters*
  - CSM\_45 Vehicle units, tachograph cards and external GNSS facilities shall support elliptic curve cryptography with a key size of 256, 384 and 512/521 bits.
  - CSM\_46 Vehicle units, tachograph cards and external GNSS facilities shall support the ECDSA signing algorithm, as specified in [DSS].
  - CSM\_47 Vehicle units, tachograph cards and external GNSS facilities shall support the ECKA-EG key agreement algorithm, as specified in [TR 03111].
  - CSM\_48 Vehicle units, tachograph cards and external GNSS facilities shall support all standardized domain parameters specified in Table 1 below for elliptic curve cryptography.

**Standardized domain parameters**

| Name            | Size (bits) | Reference         | Object identifier |
|-----------------|-------------|-------------------|-------------------|
| NIST P-256      | 256         | [DSS], [RFC 5480] | secp256r1         |
| BrainpoolP256r1 | 256         | [RFC 5639]        | brainpoolP256r1   |
| NIST P-384      | 384         | [DSS], [RFC 5480] | secp384r1         |
| BrainpoolP384r1 | 384         | [RFC 5639]        | brainpoolP384r1   |
| BrainpoolP512r1 | 512         | [RFC 5639]        | brainpoolP512r1   |
| NIST P-521      | 521         | [DSS], [RFC 5480] | secp521r1         |

*Note:* the object identifiers mentioned in the last column of Table 1 are specified in [RFC 5639] for the Brainpool curves and in [RFC 5480] for the NIST curves.

Example 1: the object identifier of the BrainpoolP256r1 curve is

{iso (1)
identified-organization (3) teletrust (36) algorithm (3)
signaturealgorithm (3) ecSign (2) ecStdCurvesAndGeneration (8)
ellipticCurve (1) versionOne (1) 7}.

Or in dot notation: 1.3.36.3.3.2.8.1.1.7.

{iso (1) identified-organization (3) certicom (132) curve (0) 34}.

Or in dot notation: `1.3.132.0.34`.

8.2.3 *Hashing algorithms*

### **M1**

CSM\_49 Vehicle units, tachograph cards and external GNSS facilities shall support the SHA-256, SHA-384 and SHA-512 algorithms specified in [SHS].

### **B**

8.2.4 *Cipher Suites*

CSM\_50 In case a symmetric algorithm, an asymmetric algorithm and/or a hashing algorithm are used together to form a security protocol, their respective key lengths and hash sizes shall be of (roughly) equal strength. Table 2 shows the allowed cipher suites:

### *Table 2*

### **Allowed cipher suites**

| Cipher suite Id | ECC key size (bits) | AES key length (bits) | Hashing<br>algorithm | MAC length<br>(bytes) |
|-----------------|---------------------|-----------------------|----------------------|-----------------------|
| CS#1            | 256                 | 128                   | SHA-256              | 8                     |
| CS#2            | 384                 | 192                   | SHA-384              | 12                    |
| CS#3            | 512/521             | 256                   | SHA-512              | 16                    |

*Note:* ECC keys sizes of 512 bits and 521 bits are considered to be equal in strength for all purposes within this Appendix.

### 9. KEYS AND CERTIFICATES

### 9.1. **Asymmetric Key Pairs and Public Key Certificates**

### 9.1.1 *General*

*Note:* the keys described in this section are used for mutual authentication and secure messaging between vehicle units and tachograph cards and between vehicle units and external GNSS facilities. These processes are described in detail in chapters 10 and 11 of this Appendix.

- CSM\_51 Within the European Smart Tachograph system, ECC key pairs and corresponding certificates shall be generated and managed through three functional hierarchical levels:
  - European level,
  - Member State level,
  - Equipment level.

CSM\_52 Within the entire European Smart Tachograph system, public and private keys and certificates shall be generated, managed and communicated using standardized and secure methods.

#### 9.1.2 *European Level*

- CSM\_53 At European level, a single unique ECC key pair designated as EUR shall be generated. It shall consist of a private key (EUR.SK) and a public key (EUR.PK). This key pair shall form the root key pair of the entire European Smart Tachograph PKI. This task shall be handled by a European Root Certificate Authority (ERCA), under the authority and responsibility of the European Commission.
- CSM\_54 The ERCA shall use the European private key to sign a (self-signed) root certificate of the European public key, and shall communicate this European root certificate to all Member States.
- CSM\_55 The ERCA shall use the European private key to sign the certificates of the Member States public keys upon request. The ERCA shall keep records of all signed Member State public key certificates.
- CSM\_56 As shown in Figure 1 in section 9.1.7, the ERCA shall generate a new European root key pair every 17 years. Whenever the ERCA generates a new European root key pair, it shall create a new self-signed root certificate for the new European public key. The validity period of a European root certificate shall be 34 years plus 3 months.

*Note:* The introduction of a new root key pair also implies that ERCA will generate a new motion sensor master key and a new DSRC master key, see sections 9.2.1.2 and 9.2.2.2.

- CSM\_57 Before generating a new European root key pair, the ERCA shall conduct an analysis of the cryptographic strength that is needed for the new key pair, given it should stay secure for the next 34 years. If found necessary, the ERCA shall switch to a cipher suite that is stronger than the current one, as specified in CSM\_50.
- **▼M1**
- CSM\_58 Whenever it generates a new European root key pair, the ERCA shall create a link certificate for the new European public key and sign it with the previous European private key. The validity period of the link certificate shall be 17 years plus 3 months. This is shown in Figure 1 in section 9.1.7 as well.
- **▼B**

*Note:* Since a link certificate contains the ERCA generation *X* public key and is signed with the ERCA generation *X-1* private key, a link certificate offers equipment issued under generation *X-1* a method to trust equipment issued under generation *X*.

CSM\_59 The ERCA shall not use the private key of a root key pair for any purpose after the moment a new root key certificate becomes valid.

CSM\_60 At any moment in time, the ERCA shall dispose of the following cryptographic keys and certificates:

- The current EUR key pair and corresponding certificate
- All previous EUR certificates to be used for the verification of MSCA certificates that are still valid
- Link certificates for all generations of EUR certificates except the first one
- 9.1.3 *Member State Level*
  - CSM\_61 At Member State level, all Member States required to sign tachograph card certificates shall generate one or more unique ECC key pairs designated as MSCA\_Card. All Member States required to sign certificates for vehicle units or external GNSS facilities shall additionally generate one or more unique ECC key pairs designated as MSCA\_VU-EGF.
  - CSM\_62 The task of generating Member State key pairs shall be handled by a Member State Certificate Authority (MSCA). Whenever a MSCA generates a Member State key pair, it shall send the public key to the ERCA in order to obtain a corresponding Member State certificate signed by the ERCA.
  - CSM\_63 An MSCA shall choose the strength of a Member State key pair equal to the strength of the European root key pair used to sign the corresponding Member State certificate.
  - CSM\_64 An MSCA\_VU-EGF key pair, if present, shall consist of private key MSCA\_VU-EGF.SK and public key MSCA\_VU-EGF.PK. An MSCA shall use the MSCA\_VU-EGF.SK private key exclusively to sign the public key certificates of vehicle units and external GNSS facilities.
  - CSM\_65 An MSCA\_Card key pair shall consist of private key MSCA\_Card.SK and public key MSCA\_Card.PK. An MSCA shall use the MSCA\_Card.SK private key exclusively to sign the public key certificates of tachograph cards.
  - CSM\_66 An MSCA shall keep records of all signed VU certificates, external GNSS facility certificates and card certificates, together with the identification of the equipment for which each certificate is intended.
  - CSM\_67 The validity period of an MSCA\_VU-EGF certificate shall be 17 years plus 3 months. The validity period of an MSCA\_Card certificate shall be 7 years plus 1 month.
  - CSM\_68 As shown in Figure 1 in section 9.1.7, the private key of a MSCA\_VU-EGF key pair and the private key of a MSCA\_Card key pair shall have a key usage period of two years.

| CSM_69 | An MSCA shall not use the private key of an MSCA_VU- EGF key pair for any purpose after the moment its usage period has ended. Neither shall an MSCA use the private key of an MSCA_Card key pair for any purpose after the moment its usage period has ended. |
|--------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
|--------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|

- CSM\_70 At any moment in time, an MSCA shall dispose of the following cryptographic keys and certificates:
  - The current MSCA\_Card key pair and corresponding certificate
  - All previous MSCA\_Card certificates to be used for the verification of the certificates of tachograph cards that are still valid
  - The current EUR certificate necessary for the verification of the current MSCA certificate
  - All previous EUR certificates necessary for the verification of all MSCA certificates that are still valid
- CSM\_71 If an MSCA is required to sign certificates for vehicle units or external GNSS facilities, it shall additionally dispose of the following keys and certificates:
  - The current MSCA\_VU-EGF key pair and corresponding certificate
  - All previous MSCA\_VU-EGF public keys to be used for the verification of the certificates of VUs or external GNSS facilities that are still valid
- 9.1.4 *Equipment Level: Vehicle Units*

# **M1**

CSM\_72 Two unique ECC key pairs shall be generated for each vehicle unit, designated as VU\_MA and VU\_Sign. This task is handled by VU manufacturers. Whenever a VU key pair is generated, the party generating the key shall send the public key to its MSCA, in order to obtain a corresponding VU certificate signed by the MSCA. The private key shall be used only by the vehicle unit.

# **B**

- CSM\_73 The VU\_MA and VU\_Sign certificates of a given vehicle unit shall have the same Certificate Effective Date.
- CSM\_74 A VU manufacturer shall choose the strength of a VU key pair equal to the strength of the MSCA key pair used to sign the corresponding VU certificate.
- CSM\_75 A vehicle unit shall use its VU\_MA key pair, consisting of private key VU\_MA.SK and public key VU\_MA.PK, exclusively to perform VU Authentication towards tachograph cards and external GNSS facilities, as specified in sections 10.3 and 11.4 of this Appendix.
- CSM\_76 A vehicle unit shall be capable of generating ephemeral ECC key pairs and shall use an ephemeral key pair exclusively to perform session key agreement with a tachograph card or external GNSS facility, as specified in sections 10.4 and 11.4 of this Appendix.

- CSM\_77 A vehicle unit shall use the private key VU\_Sign.SK of its VU\_Sign key pair exclusively to sign downloaded data files, as specified in chapter 14 of this Appendix. The corresponding public key VU\_Sign.PK shall be used exclusively to verify signatures created by the vehicle unit.
- CSM\_78 As shown in Figure 1 in section 9.1.7, the validity period of a VU\_MA certificate shall be 15 years and 3 months. The validity period of a VU\_Sign certificate shall also be 15 years and 3 months.

*Notes:*

- The extended validity period of a VU\_Sign certificate allows a Vehicle Unit to create valid signatures over downloaded data during the first three months after it has expired, as required in Regulation (EU) No 581/2010.
- The extended validity period of a VU\_MA certificate is needed to allow the VU to authenticate to a control card or a company card during the first three months after it has expired, such that is it possible to perform a data download.
- CSM\_79 A vehicle unit shall not use the private key of a VU key pair for any purpose after the corresponding certificate has expired.
- CSM\_80 The VU key pairs (except ephemeral keys pairs) and corresponding certificates of a given vehicle unit shall not be replaced or renewed in the field once the vehicle unit has been put in operation.

#### *Notes:*

- Ephemeral key pairs are not included in this requirement, as a new ephemeral key pair is generated by a VU each time Chip Authentication and session key agreement is performed, see section 10.4. Note that ephemeral key pairs do not have corresponding certificates.
- This requirement does not forbid the possibility of replacing static VU key pairs during a refurbishment or repair in a secure environment controlled by the VU manufacturer.
- CSM\_81 When put in operation, vehicle units shall contain the following cryptographic keys and certificates:
  - The VU\_MA private key and corresponding certificate
  - The VU\_Sign private key and corresponding certificate
  - The MSCA\_VU-EGF certificate containing the MSCA\_VU-EGF.PK public key to be used for verification of the VU\_MA certificate and VU\_Sign certificate

- The EUR certificate containing the EUR.PK public key to be used for verification of the MSCA\_VU-EGF certificate
- The EUR certificate whose validity period directly precedes the validity period of the EUR certificate to be used to verify the MSCA\_VU-EGF certificate, if existing
- The link certificate linking these two EUR certificates, if existing
- CSM\_82 In addition to the cryptographic keys and certificates listed in CSM\_81, vehicle units shall also contain the keys and certificates specified in Part A of this Appendix, allowing a vehicle unit to interact with first-generation tachograph cards.
- 9.1.5 *Equipment Level: Tachograph Cards*

### **M1**

CSM\_83 One unique ECC key pair, designated as Card\_MA, shall be generated for each tachograph card. A second unique ECC key pair, designated as Card\_Sign, shall additionally be generated for each driver card and each workshop card. This task may be handled by card manufacturers or card personalisers. Whenever a card key pair is generated, the party generating the key shall send the public key to its MSCA, in order to obtain a corresponding card certificate signed by the MSCA. The private key shall be used only by the tachograph card.

# **B**

- CSM\_84 The Card\_MA and Card\_Sign certificates of a given driver card or workshop card shall have the same Certificate Effective Date.
- CSM\_85 A card manufacturer or card personaliser shall choose the strength of a card key pair equal to the strength of the MSCA key pair used to sign the corresponding card certificate.
- CSM\_86 A tachograph card shall use its Card\_MA key pair, consisting of private key Card\_MA.SK and public key Card\_MA.PK, exclusively to perform mutual authentication and session key agreement towards vehicle units, as specified in sections 10.3 and 10.4 of this Appendix.
- CSM\_87 A driver card or workshop card shall use the private key Card\_Sign.SK of its Card\_Sign key pair exclusively to sign downloaded data files, as specified in chapter 14 of this Appendix. The corresponding public key Card\_Sign.PK shall be used exclusively to verify signatures created by the card.

# **M1**

- CSM\_88 The validity period of a Card\_MA certificate shall be as follows:
  - For driver cards: 5 years
  - For company cards: 5 years
  - For control cards: 2 years
  - For workshop cards: 1 year

CSM\_89 The validity period of a Card\_Sign certificate shall be as follows:

| — For driver cards:   | 5 years and 1 month |
|-----------------------|---------------------|
| — For workshop cards: | 1 year and 1 month  |

*Note:* the extended validity period of a Card\_Sign certificate allows a driver card to create valid signatures over downloaded data during the first month after it has expired. This is necessary in view of Regulation (EU) No 581/2010, which requires that a data download from a driver card must be possible up to 28 days after the last data has been recorded.

- CSM\_90 The key pairs and corresponding certificates of a given tachograph card shall not be replaced or renewed once the card has been issued.
- CSM\_91 When issued, tachograph cards shall contain the following cryptographic keys and certificates:
  - The Card\_MA private key and corresponding certificate
  - For driver cards and workshop cards additionally: the Card\_Sign private key and corresponding certificate
  - The MSCA\_Card certificate containing the MSCA\_Card.PK public key to be used for verification of the Card\_MA certificate and Card\_Sign certificate
  - The EUR certificate containing the EUR.PK public key to be used for verification of the MSCA\_Card certificate.
  - The EUR certificate whose validity period directly precedes the validity period of the EUR certificate to be used to verify the MSCA\_Card certificate, if existing.
  - The link certificate linking these two EUR certificates, if existing.
  - Additionally, for control cards, company cards and workshop cards only, and only if such cards are issued during the first three months of the validity period of a new EUR certificate: the EUR certificate that is two generations older, if existing.

*Note to last bullet:* For example, in the first three months of the ERCA(3) certificate (see Figure 1), the mentioned cards shall contain the ERCA(1) certificate. This is needed to ensure that these cards can be used to perform data downloads from ERCA(1) VUs whose normal 15-year life period plus the 3-months data downloading period expires during these months; see the last bullet of requirement 13) in Annex IC.

CSM\_92 In addition to the cryptographic keys and certificates listed in CSM\_91, tachograph cards shall also contain the keys and certificates specified in Part A of this Appendix, allowing these cards to interact with first-generation VUs.

# **M1**

# **B**

|     | 9.1.6 |        | Equipment Level: External GNSS Facilities                                                                                                                                                                                                                                                                                                                                                                                                      |
|-----|-------|--------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| ▼M1 |       | CSM_93 | One unique ECC key pair shall be generated for each<br>external GNSS facility, designated as EGF_MA. This<br>task is handled by external GNSS facility manufacturers.<br>Whenever an EGF_MA key pair is generated, the party<br>generating th e key shall send the public key to its MSCA<br>in order to obtain a corresponding EGF_MA certificate<br>signed by the MSCA. The private key shall be used<br>only by the external GNSS facility. |
| ▼B  |       | CSM_94 | An EGF manufacturer shall choose the strength of an<br>EGF_MA key pair equal to the strength of the MSCA<br>key pair used to sign the corresponding EGF_MA<br>certificate.                                                                                                                                                                                                                                                                     |
| ▼M1 |       | CSM_95 | An external GNSS facility shall use its EGF_MA key pair,<br>consisting of private key EGF_MA.SK and public key<br>EGF_MA.PK, exclusively to perform mutual authenti-<br>cation and session key agreement towards vehicle units,<br>as specified in section 11.4 of this Appendix.                                                                                                                                                              |
| ▼B  |       | CSM_96 | The validity period of an EGF_MA certificate shall be 15<br>years.                                                                                                                                                                                                                                                                                                                                                                             |
|     |       | CSM_97 | An external GNSS facility shall not use the private key of<br>its EGF_MA key pair for coupling to a vehicle unit after<br>the corresponding certificate has expired.                                                                                                                                                                                                                                                                           |
|     |       |        | Note: as explained in section 11.3.3, an EGF may<br>potentially use its private key for mutual authentication<br>towards the VU it is already coupled to, even after the<br>corresponding certificate has expired.                                                                                                                                                                                                                             |
|     |       | CSM_98 | The EGF_MA key pair and corresponding certificate of a<br>given external GNSS facility shall not be replaced or<br>renewed in the field once the EGF has been put in<br>operation.                                                                                                                                                                                                                                                             |
|     |       |        | Note: This requirement does not forbid the possibility of<br>replacing EGF key pairs during a refurbishment or repair<br>in a secure environment controlled by the EGF manu-<br>facturer.                                                                                                                                                                                                                                                      |
|     |       | CSM_99 | When put in operation, an external GNSS facility shall<br>contain the following cryptographic keys and certificates:                                                                                                                                                                                                                                                                                                                           |
|     |       |        | — The EGF_MA private key and corresponding certifi-<br>cate                                                                                                                                                                                                                                                                                                                                                                                    |
|     |       |        | — The MSCA_VU-EGF certificate containing the<br>MSCA_VU-EGF.PK public key to be used for verifi-<br>cation of the EGF_MA certificate                                                                                                                                                                                                                                                                                                           |
|     |       |        | — The EUR certificate containing the EUR.PK public key<br>to be used for verification of the MSCA_VU-EGF                                                                                                                                                                                                                                                                                                                                       |

certificate

- The EUR certificate whose validity period directly precedes the validity period of the EUR certificate to be used to verify the MSCA\_VU-EGF certificate, if existing
- The link certificate linking these two EUR certificates, if existing

#### 9.1.7 *Overview: Certificate Replacement*

Figure 1 below shows how different generations of ERCA root certificates, ERCA link certificates, MSCA certificates and equipment (VU and card) certificates are issued and used over time:

### **M1**

#### *Figure 1*

#### **Issuance and usage of different generations of ERCA root certificates, ERCA link certificates, MSCA certificates and equipment certificates**

![](_page_33_Figure_8.jpeg)

*Notes to Figure 1:*

- 1. Different generations of the root certificate are indicated by a number in brackets. E.g. ERCA (1) is the first generation of ERCA root certificate; ERCA (2) is the second generation, etc.
- 2. Other certificates are indicated by two numbers in brackets, the first one indicating the root certificate generation under which they are issued, the second one the generation of the certificate itself. E.g. MSCA\_Card (1-1) is the first MSCA\_Card certificate issued under ERCA (1); MSCA\_Card (2-1) is the first MSCA\_Card certificate issued under ERCA (2); MSCA\_Card (2-last) is the last MSCA\_Card certificate issued under ERCA (2); Card\_MA(2-1) is the first Card certificate for mutual authentication that is issued under ERCA (2), etc.
- 3. The MSCA\_Card (2-1) and MSCA\_Card (1-last) certificates are issued at almost but not exactly the same date. MSCA\_Card (2-1) is the first MSCA\_Card certificate issued under ERCA (2) and will be issued slightly later than MSCA\_Card (1-last), the last MSCA\_Card certificate under ERCA (1).
- 4. As shown in the figure, the first VU and Card certificates issued under ERCA (2) will appear almost two years before the last VU and Card certificates issued under ERCA (1) will appear. This is because of the fact that VU and Card certificates are issued under an MSCA certificate, not directly under the ERCA certificate. The MSCA (2-1) certificate will be issued directly after ERCA (2) becomes valid, but the MSCA (1-last) certificate will be issued only slightly before that time, at the last moment the ERCA (1) certificate is still valid. Therefore, these two MSCA certificates will have almost the same validity period, despite the fact that they are of different generations.
- 5. The validity period shown for cards is the one for driver cards (5 years).

### **M1**

6. To save space, the difference in validity period between the Card\_MA and Card\_Sign certificates is shown only for the first generation.

#### **B** 9.2. **Symmetric Keys**

- 9.2.1 *Keys for Securing VU Motion Sensor Communication*
- 9.2.1.1 *General*

*Note:* readers of this section are supposed to be familiar with the contents of [ISO 16844-3] describing the interface between a vehicle unit and a motion sensor. The pairing process between a VU and a motion sensor is described in detail in chapter 12 of this Appendix.

CSM\_100 A number of symmetric keys is needed for pairing vehicle units and motion sensors, for mutual authentication between vehicle units and motion sensors and for encrypting communication between vehicle units and motion sensors, as shown in Table 3. All of these keys shall be AES keys, with a key length equal to the length of the motion sensor master key, which shall be linked to the length of the (foreseen) European root key pair as described in CSM\_50.

| Key                                      | Symbol     | Generated by                                | Generation method                                                      | Stored by                                                                         |
|------------------------------------------|------------|---------------------------------------------|------------------------------------------------------------------------|-----------------------------------------------------------------------------------|
| Motion Sensor Master Key – VU part       | $K_{M-VU}$ | ERCA                                        | Random                                                                 | ERCA, MSCAs involved in issuing VUs certificates, VU manufacturers, vehicle units |
| Motion Sensor Master Key – Workshop part | $K_{M-WC}$ | ERCA                                        | Random                                                                 | ERCA, MSCAs, card manufacturers, workshop cards                                   |
| Motion Sensor Master Key                 | $K_M$      | Not independently generated                 | Calculated as $K_M = K_{M-VU} XOR K_{M-WC}$                            | ERCA, MSCAs involved in issuing motion sensors keys (optionally) (*)              |
| Identification Key                       | $K_{ID}$   | Not independently generated                 | Calculated as $K_{ID} = K_M XOR CV$ , where CV is specified in CSM_106 | ERCA, MSCAs involved in issuing motion sensors keys (optionally) (*)              |
| Pairing Key                              | $K_P$      | Motion sensor manufacturer                  | Random                                                                 | One motion sensor                                                                 |
| Session Key                              | $K_S$      | VU (during pairing of VU and motion sensor) | Random                                                                 | One VU and one motion sensor                                                      |

| Table 3 |  |
|---------|--|
|         |  |

**Keys for securing vehicle unit — motion sensor communication**

(\*) Storage of KM and KID is optional, as these keys can be derived from KM-VU, KM-WC and CV.

- CSM\_101 The European Root Certificate Authority shall generate KM-VU and KM-WC, two random and unique AES keys from which the motion sensor master key KM can be calculated as KM-VU XOR KM-WC. The ERCA shall communicate KM, KM-VU and KM-WC to Member State Certificate Authorities upon their request.
- CSM\_102 The ERCA shall assign to each motion sensor master key KM a unique version number, which shall also be applicable for the constituting keys KM-VU and KM-WC and for the related identification key KID. The ERCA shall inform the MSCAs about the version number when sending KM-VU and KM-WC to them.

*Note:* The version number is used to distinguish different generations of these keys, as explained in detail in section 9.2.1.2.

- CSM\_103 A Member State Certificate Authority shall forward KM-VU, together with its version number, to vehicle unit manufacturers upon their request. The VU manufacturers shall insert KM-VU and its version number in all manufactured VUs.
- CSM\_104 A Member State Certificate Authority shall ensure that KM-WC, together with its version number, is inserted in every workshop card issued under its responsibility.

*Notes:*

<sup>—</sup> See the description of data type in Appendix 2.

- as explained in section 9.2.1.2, in fact multiple generations of KM-WC may have to be inserted in a single workshop card.
- CSM\_105 In addition to the AES key specified in CSM\_104, a MSCA shall ensure that the TDES key KmWC, specified in requirement CSM\_037 in Part A of this Appendix, is inserted in every workshop card issued under its responsibility.

*Notes:*

- This allows a second-generation workshop card to be used for coupling a first-generation VU.
- A second-generation workshop card will contain two different applications, one complying with Part B of this Appendix and one complying with Part A. The latter will contain the TDES key KmWC.
- CSM\_106 An MSCA involved in issuing motion sensors shall derive the identification key from the motion sensor master key by XORing it with a constant vector CV. The value of CV shall be as follows:

**▼M1**

- **▼B**

— For 192-bit motion sensor master keys: CV = '72 AD EA FA 00 BB F4 EE F4 99 15 70 5B 7E EE BB 1C 54 ED 46 8B 0E F8 25'

— For 128-bit motion sensor master keys: CV = 'B6 44 2C 45 0E F8 D3 62 0B 7A 8A 97 91 E4 5D 83'

— For 256-bit motion sensor master keys: CV = '1D 74 DB F0 34 C7 37 2F 65 55 DE D5 DC D1 9A C3 23 D6 A6 25 64 CD BE 2D 42 0D 85 D2 32 63 AD 60'

*Note:* the constant vectors have been generated as follows:

Pi\_10 = first 10 bytes of the decimal portion of the mathematical constant π = '24 3F 6A 88 85 A3 08 D3 13 19'

CV\_128-bits = first 16 bytes of SHA-256(Pi\_10)

CV\_192-bits = first 24 bytes of SHA-384(Pi\_10)

CV\_256-bits = first 32 bytes of SHA-512(Pi\_10)

CSM\_107 **►M1** Each Motion sensor manufacturer shall generate a random and unique pairing key KP for every motion sensor, and shall send each pairing key to its Member State Certificate Authority. The MSCA shall encrypt each pairing key separately with the motion sensor master key KM and shall return the encrypted key to the motion sensor manufacturer. For each encrypted key, the MSCA shall notify the motion sensor manufacturer of the version number of the associated KM. ◄

> *Note:* as explained in section 9.2.1.2, in fact a motion sensor manufacturer may have to generate multiple unique pairing keys for a single motion sensor.

- CSM\_108 Each motion sensor manufacturer shall generate a unique serial number for every motion sensor, and shall send all serial numbers to its Member State Certificate Authority. The MSCA shall encrypt each serial number separately with the identification key KID and shall return the encrypted serial number to the motion sensor manufacturer. For each encrypted serial number, the MSCA shall notify the motion sensor manufacturer of the version number of the associated KID.
- **▼B**
- CSM\_109 For requirements CSM\_107 and CSM\_108, the MSCA shall use the AES algorithm in the Cipher Block Chaining mode of operation, as defined in [ISO 10116], with an interleave parameter *m* = 1 and an initialization vector SV = '00' {16}, i.e. sixteen bytes with binary value 0. When necessary, the MSCA shall use padding method 2 defined in [ISO 9797-1].
- CSM\_110 The motion sensor manufacturer shall store the encrypted pairing key and the encrypted serial number in the intended motion sensor, together with the corresponding plain text values and the version number of KM and KID used for encrypting.

*Note:* as explained in section 9.2.1.2, in fact a motion sensor manufacturer may have to insert multiple encrypted pairing keys and multiple encrypted serial numbers in a single motion sensor.

CSM\_111 In addition to the AES-based cryptographic material specified in CSM\_110, a motion sensor manufacturer may also store in each motion sensor the TDES-based cryptographic material specified in requirement CSM\_037 in Part A of this Appendix.

> *Note:* doing so will allow a second-generation motion sensor to be coupled to a first-generation VU.

- CSM\_112 The length of the session key KS generated by a VU during the pairing to a motion sensor shall be linked to the length of its KM-VU, as described in CSM\_50.
- 9.2.1.2 *Motion Sensor Master Key Replacement in Second-Generation Equipment*
  - CSM\_113 Each motion sensor master key and all related keys (see Table 3) is associated to a particular generation of the ERCA root key pair. These keys shall therefore be replaced every 17 years. The validity period of each motion sensor master key generation shall begin one year before the associated ERCA root key pair becomes valid and shall end when the associated ERCA root key pair expires. This is depicted in Figure 2.

### *Figure 2*

### **Issuance and usage of different generations of the motion sensor master key in vehicle units, motions sensors and workshop cards**

![](_page_38_Figure_3.jpeg)

- CSM\_114 At least one year before generating a new European root key pair, as described in CSM\_56, the ERCA shall generate a new motion sensor master key KM by generating a new KM-VU and KM-WC. The length of the motion sensor master key shall be linked to the foreseen strength of the new European root key pair, according to CSM\_50. The ERCA shall communicate the new KM, KM-VU and KM-WC to the MSCAs upon their request, together with their version number.
- CSM\_115 An MSCA shall ensure that all valid generations of KM-WC are stored in every workshop card issued under its authority, together with their version numbers, as shown in Figure 2.

*Note:* this implies that in the last year of the validity period of an ERCA certificate, workshop cards will be issued with three different generations of KM-WC, as shown in Figure 2.

CSM\_116 In relation to the process described in CSM\_107 and CSM\_108 above: An MSCA shall encrypt each pairing key KP it receives from a motion sensor manufacturer separately with each valid generation of the motion sensor master key KM. An MSCA shall also encrypt each serial number it receives from a motion sensor manufacturer separately with each valid generation of the identification key KID. A motion sensor manufacturer shall store all encryptions of the pairing key and all encryptions of the serial number in the intended motion sensor, together with the corresponding plain text values and the version number(s) of KM and KID used for encrypting.

*Note:* This implies that in the last year of the validity period of an ERCA certificate, motion sensors will be issued with encrypted data based on three different generations of KM, as shown in Figure 2.

CSM\_117 In relation to the process described in CSM\_107 above: Since the length of the pairing key KP shall be linked to the length of KM (see CSM\_100), a motion sensor manufacturer may have to generate up to three different pairing keys (of different lengths) for one motion sensor, in case subsequent generations of KM have different lengths. In such a case, the manufacturer shall send each pairing key to the MSCA. The MSCA shall ensure that each pairing key is encrypted with the correct generation of the motion sensor master key, i.e. the one having the same length.

> *Note:* In case the motion sensor manufacturer chooses to generate a TDES-based pairing key for a second-generation motion sensor (see CSM\_111), the manufacturer shall indicate to the MSCA that the TDES-based motion sensor master key must be used for encrypting this pairing key. This is because the length of a TDES key may be equal to that of an AES key, so the MSCA cannot judge from the key length alone.

CSM\_118 Vehicle unit manufacturers shall insert only one generation of KM-VU in each vehicle unit, together with its version number. This KM-VU generation shall be linked to the ERCA certificate upon which the VU's certificates are based.

*Notes:*

- A vehicle unit based on the generation *X* ERCA certificate shall only contain the generation *X* KM-VU, even if it is issued after the start of the validity period of the generation *X+1* ERCA certificate. This is shown in Figure 2.
- A VU of generation *X* cannot be paired to a motion sensor of generation *X-1*.
- Since workshop cards have a validity period of one year, the result of CSM\_113 — CSM\_118 is that all workshop cards will contain the new KM-WC at the moment the first VU containing the new KM-VU is issued. Therefore, such a VU will always be able to calculate the new KM. Moreover, by that time most new motion sensors will contain encrypted data based on the new KM as well.
- 9.2.2 *Keys for Securing DSRC Communication*
- 9.2.2.1 *General*
  - CSM\_119 The authenticity and confidentiality of data communicated from a vehicle unit to a control authority over a DSRC remote communication channel shall be ensured by means of a set of VU-specific AES keys derived from a single DSRC master key, KMDSRC.

- CSM\_120 The DSRC master key KMDSRC shall be an AES key that is securely generated, stored and distributed by the ERCA. The key length may be 128, 192 or 256 bits and shall be linked to the length of the European root key pair, as described in CSM\_50.
- CSM\_121 The ERCA shall communicate the DSRC master key to Member State Certificate Authorities upon their request in a secure manner, to allow them to derive VU-specific DSRC keys and to ensure that the DSRC master key is inserted in all control cards and workshop cards issued under their responsibility.
- CSM\_122 The ERCA shall assign to each DSRC master key a unique version number. The ERCA shall inform the MSCAs about the version number when sending the DSRC master key to them.

*Note:* The version number is used to distinguish different generations of the DSRC master key, as explained in detail in section 9.2.2.2.

- **▼M1**
- CSM\_123 For every vehicle unit, the vehicle unit manufacturer shall create a unique VU serial number and shall send this number to its Member State Certificate Authority in a request to obtain a set of two VU-specific DSRC keys. The VU serial number shall have data type .

#### *Note:*

- This VU serial number shall be identical to the vuSerialNumber element of VuIdentification, see Appendix 1 and to the Certificate Holder Reference in the VU's certificates.
- The VU serial number may not be known at the moment a vehicle unit manufacturer requests the VU-specific DSRC keys. In this case, the VU manufacturer shall send instead the unique certificate request ID it used when requesting the VU's certificates; see CSM\_153. This certificate request ID shall therefore be equal to the Certificate Holder Reference in the VU's certificates.

- **▼B**
- CSM\_124 Upon receiving a request for VU-specific DSRC keys, the MSCA shall derive two AES keys for the vehicle unit, called K\_VUDSRC\_ENC and K\_VUDSRC\_MAC. These VU-specific keys shall have the same length as the DSRC master key. The MSCA shall use the key derivation function defined in [RFC 5869]. The hash function that is necessary to instantiate the HMAC-Hash function shall be linked to the length of the DSRC master key, as described in CSM\_50. The key derivation function in [RFC 5869] shall be used as follows:

Step 1 (Extract):

<sup>—</sup> *PRK* = HMAC-Hash (*salt, IKM*) where *salt* is an empty string '' and *IKM* is KMDSRC.

Step 2 (Expand):

— *OKM* = *T(1)*, where

*T(1)* = HMAC-Hash (*PRK*, *T(0)* || *info* || '01') with

— *T(0)* = an empty string ('')

— **►M1** *info* = VU serial number or certificate request ID, as specified in CSM\_123 ◄

— K\_VUDSRC\_ENC = first *L* octets of *OKM* and

K\_VUDSRC\_MAC = last *L* octets of *OKM*

where *L* is the required length of K\_VUDSRC\_ENC and K\_VUDSRC\_MAC in octets.

- CSM\_125 The MSCA shall distribute K\_VUDSRC\_ENC and K\_VUDSRC\_MAC to the VU manufacturer in a secure manner for insertion in the intended vehicle unit.
- CSM\_126 When issued, a vehicle unit shall have stored K\_VUDSRC\_ENC and K\_VUDSRC\_MAC in its secure memory, in order to be able to ensure the integrity, authenticity and confidentiality of data sent over the remote communication channel. A vehicle unit shall also store the version number of the DSRC master key used to derive these VU-specific keys.
- CSM\_127 When issued, control cards and workshop cards shall have stored KMDSRC in their secure memory, in order to be able to verify the integrity and authenticity of data sent by a VU over the remote communication channel and to decrypt this data. Control cards and workshop cards shall also store the version number of the DSRC master key.

*Note:* as explained in section 9.2.2.2, in fact multiple generations of KMDSRC may have to be inserted in a single workshop card or control card.

### **M1**

CSM\_128 The MSCA shall keep records of all VU-specific DSRC keys it generated, their version number and the VU serial number or certificate request ID used in deriving them.

# **B**

- 9.2.2.2 *DSRC Master Key Replacement*
  - CSM\_129 Each DSRC master key is associated to a particular generation of the ERCA root key pair. The ERCA shall therefore replace the DSRC master key every 17 years. The validity period of each DSRC master key generation shall begin two years before the associated ERCA root key pair becomes valid and shall end when the associated ERCA root key pair expires. This is depicted in Figure 3.

#### *Figure 3*

### **Issuance and usage of different generations of the DSRC master key in vehicle units, workshop cards and control cards**

![](_page_42_Figure_4.jpeg)

- CSM\_130 At least two years before generating a new European root key pair, as described in CSM\_56, the ERCA shall generate a new DSRC master key. The length of the DSRC key shall be linked to the foreseen strength of the new European root key pair, according to CSM\_50. The ERCA shall communicate the new DSRC master key to the MSCAs upon their request, together with its version number.
- CSM\_131 An MSCA shall ensure that all valid generations of KMDSRC are stored in every control card issued under its authority, together with their version numbers, as shown in Figure 3.

*Note:* this implies that in the last two years of the validity period of an ERCA certificate, control cards will be issued with three different generations of KMDSRC, as shown in Figure 3.

CSM\_132 An MSCA shall ensure that all generations of KMDSRC that have been valid for at least a year and are still valid, are stored in every workshop card issued under its authority, together with their version numbers, as shown in Figure 3.

*Note:* this implies that in the last year of the validity period of an ERCA certificate, workshop cards will be issued with three different generations of KMDSRC, as shown in Figure 3.

CSM\_133 Vehicle unit manufacturers shall insert only one set of VU-specific DSRC keys into each vehicle unit, together with its version number. This set of keys shall be derived from the KMDSRC generation linked to the ERCA certificate upon which the VU's certificates are based.

*Notes:*

- This implies that a vehicle unit based on the generation *X* ERCA certificate shall only contain the generation *X* K\_VUDSRC\_ENC and K\_VUDSRC\_MAC, even if the VU is issued after the start of the validity period of the generation *X+1* ERCA certificate. This is shown in Figure 3.
- Since workshop cards have a validity period of one year and control cards of two years, the result of CSM\_131 — CSM\_133 is that all workshop cards and control cards will contain the new DSRC master key at the moment the first VU containing VU-specific keys based on that master key will be issued.

#### 9.3. **Certificates**

- 9.3.1 *General*
  - CSM\_134 All certificates in the European Smart Tachograph system shall be self-descriptive, card-verifiable (CV) certificates according to [ISO 7816-4] and [ISO 7816-8].
  - CSM\_135 **►M1** The Distinguished Encoding Rules (DER) according to [ISO 8825-1] shall be used to encode the data objects within certificates. Table 4 shows the full certificate encoding, including all tag and length bytes. ◄

*Note:* this encoding results in a Tag-Length-Value (TLV) structure as follows:

- Tag: The tag is encoded in one or two octets and indicates the content.
- Length: The length is encoded as an unsigned integer in one, two, or three octets, resulting in a maximum length of 65 535 octets. The minimum number of octets shall be used.
- Value: The value is encoded in zero or more octets

### 9.3.2 *Certificate Content*

CSM\_136 All certificates shall have the structure shown in the certificate profile in Table 4.

### *Table 4*

### **Certificate Profile version 1**

| Field                | Field ID | Tag     | Length (bytes) | ASN.1 data type<br>(see Appendix 1) |
|----------------------|----------|---------|----------------|-------------------------------------|
| ECC Certificate      | C        | '7F 21' | var            |                                     |
| ECC Certificate Body | B        | '7F 4E' | var            |                                     |

| Field                                  | Field ID | Tag     | Length (bytes) | ASN.1 data type<br>(see Appendix 1) |
|----------------------------------------|----------|---------|----------------|-------------------------------------|
| Certificate Profile<br>Identifier      | CPI      | '5F 29' | '01'           | INTEGER(0..255)                     |
| Certificate Authority<br>Reference     | CAR      | '42'    | '08'           | KeyIdentifier                       |
| Certificate<br>Holder<br>Authorisation | CHA      | '5F 4C' | '07'           | CertificateHolder<br>Authorisation  |
| Public Key                             | PK       | '7F 49' | var            |                                     |
| Domain Parameters                      | DP       | '06'    | var            | OBJECT IDENTIFIER                   |
| Public Point                           | PP       | '86'    | var            | OCTET STRING                        |
| Certificate<br>Holder<br>Reference     | CHR      | '5F 20' | '08'           | KeyIdentifier                       |
| Certificate<br>Effective<br>Date       | CEfD     | '5F 25' | '04'           | TimeReal                            |
| Certificate Expiration<br>Date         | CExD     | '5F 24' | '04'           | TimeReal                            |
| ECC<br>Certificate<br>Signature        | S        | '5F 37' | var            | OCTET STRING                        |

*Note:* the Field ID will be used in later sections of this Appendix to indicate individual fields of a certificate, e.g. X.CAR is the Certificate Authority Reference mentioned in the certificate of user X.

- 9.3.2.1 Certificate Profile Identifier
  - CSM\_137 Certificates shall use a Certificate Profile Identifier to indicate the certificate profile used. Version 1, as specified in Table 4, shall be identified by a value of '00'.
- 9.3.2.2 Certificate Authority Reference
  - CSM\_138 The Certificate Authority Reference shall be used to identify the public key to be used to verify the certificate signature. The Certificate Authority Reference shall therefore be equal to the Certificate Holder Reference in the certificate of the corresponding certificate authority.
  - CSM\_139 An ERCA root certificate shall be self-signed, i.e., the Certificate Authority Reference and the Certificate Holder Reference in the certificate shall be equal.

CSM\_140 For an ERCA link certificate, the Certificate Holder Reference shall be equal to the CHR of the new ERCA root certificate. The Certificate Authority Reference for a link certificate shall be equal to the CHR of the previous ERCA root certificate.

9.3.2.3 Certificate Holder Authorisation

### **M1**

CSM\_141 The Certificate Holder Authorisation shall be used to identify the type of certificate. It consists of the six most significant bytes of the Tachograph Application ID, concatenated with the equipment type, which indicates the type of equipment for which the certificate is intended. In the case of a VU certificate, a driver card certificate or a workshop card certificate, the equipment type is also used to differentiate between a certificate for Mutual Authentication and a certificate for creating digital signatures (see section 9.1 and Appendix 1, data type EquipmentType).

# **B**

#### 9.3.2.4 Public Key

The Public Key nests two data elements: the standardized domain parameters to be used with the public key in the certificate and the value of the public point.

- CSM\_142 The data element Domain Parameters shall contain one of the object identifiers specified in Table 1 to reference a set of standardized domain parameters.
- CSM\_143 The data element Public Point shall contain the public point. Elliptic curve public points shall be converted to octet strings as specified in [TR-03111]. The uncompressed encoding format shall be used. When recovering an elliptic curve point from its encoded format, the validations described in [TR-03111] shall always be carried out.

### 9.3.2.5 Certificate Holder Reference

- CSM\_144 The Certificate Holder Reference is an identifier for the public key provided in the certificate. It shall be used to reference this public key in other certificates.
- CSM\_145 For card certificates and external GNSS facility certificates, the Certificate Holder Reference shall have the data type specified in Appendix 1.
- CSM\_146 For vehicle units, the manufacturer, when requesting a certificate, may or may not know the manufacturer-specific serial number of the VU for which that certificate and the associated private key is intended. In the first case, the Certificate Holder Reference shall have the data type specified in Appendix 1. In the latter case, the Certificate Holder Reference shall have the data type specified in Appendix 1.

*Note:* For a card certificate, the value of the CHR shall be equal to the value of the cardExtendedSerialNumber in EF\_ICC; see Appendix 2. For an EGF certificate, the value of the CHR shall be equal to the value of the sensorGNSSSerialNumber in EF\_ICC; see Appendix 14. For a VU certificate, the value of the CHR shall be equal to the vuSerialNumber element of VuIdentification, see Appendix 1, unless the manufacturer does not know the manufacturer-specific serial number at the time the certificate is requested.

**▼B**

CSM\_147 For ERCA and MSCA certificates, the Certificate Holder Reference shall have the data type specified in Appendix 1.

9.3.2.6 Certificate Effective Date

### **M1**

CSM\_148 The Certificate Effective Date shall indicate the starting date and time of the validity period of the certificate.

**▼B**

- 9.3.2.7 Certificate Expiration Date
  - CSM\_149 The Certificate Expiration Date shall indicate the end date and time of the validity period of the certificate.
- 9.3.2.8 Certificate Signature
  - CSM\_150 The signature on the certificate shall be created over the encoded certificate body, including the certificate body tag and length. The signature algorithm shall be ECDSA, as specified in [DSS], using the hashing algorithm linked to the key size of the signing authority, as specified in CSM\_50. The signature format shall be plain, as specified in [TR-03111].

#### 9.3.3 *Requesting Certificates*

- CSM\_151 **►M1** When requesting a certificate, an MSCA shall send the following data to the ERCA: ◄
  - The Certificate Profile Identifier of the requested certificate
  - The Certificate Authority Reference expected to be used for signing the certificate.
  - The Public Key to be signed
- CSM\_152 In addition to the data in CSM\_151, an MSCA shall send the following data in a certificate request to the ERCA, allowing the ERCA to create the Certificate Holder Reference of the new MSCA certificate:
  - The numerical nation code of the Certification Authority (data type defined in Appendix 1)
  - The alphanumerical nation code of the Certification Authority (data type defined in Appendix 1)
  - The 1-byte serial number to distinguish the different keys of the Certification Authority in the case keys are changed
  - The two-byte field containing Certification Authority specific additional info

- CSM\_153 An equipment manufacturer shall send the following data in a certificate request to an MSCA, allowing the MSCA to create the Certificate Holder Reference of the new equipment certificate:
  - If known (see CSM\_154), a serial number for the equipment, unique for the manufacturer, the equipment's type and the month of manufacturing. Otherwise, a unique certificate request identifier.
  - The month and the year of equipment manufacturing or of the certificate request.

The manufacturer shall ensure that this data is correct and that the certificate returned by the MSCA is inserted in the intended equipment.

### **B**

CSM\_154 In the case of a VU, the manufacturer, when requesting a certificate, may or may not know the manufacturer-specific serial number of the VU for which that certificate and the associated private key is intended. If known, the VU manufacturer shall send the serial number to the MSCA. If not known, the manufacturer shall uniquely identify each certificate request and send this certificate request serial number to the MSCA. The resulting certificate will then contain the certificate request serial number. After inserting the certificate in a specific VU, the manufacturer shall communicate the connection between the certificate request serial number and the VU identification to the MSCA.

### 10. VU- CARD MUTUAL AUTHENTICATION AND SECURE MESSAGING

### 10.1. **General**

- CSM\_155 On a high level, secure communication between a vehicle unit and a tachograph card shall be based on the following steps:
  - First, each party shall demonstrate to the other that it owns a valid public key certificate, signed by a Member State Certificate Authority. In turn, the MSCA public key certificate must be signed by the European root certificate authority. This step is called certificate chain verification and is specified in detail in section 10.2
  - Second, the vehicle unit shall demonstrate to the card that it is in possession of the private key corresponding to the public key in the presented certificate. It does so by signing a random number sent by the card. The card verifies the signature over the random number. If this verification is successful, the VU is authenticated. This step is called VU Authentication and is specified in detail in section 10.3.

- Third, both parties independently calculate two AES session keys using an asymmetric key agreement algorithm. Using one of these session keys, the card creates a message authentication code (MAC) over some data sent by the VU. The VU verifies the MAC. If this verification is successful, the card is authenticated. This step is called Card Authentication and is specified in detail in section 10.4.
- Fourth, the VU and the card shall use the agreed session keys to ensure the confidentiality, integrity and authenticity of all exchanged messages. This is called Secure Messaging and is specified in detail in section 10.5.
- CSM\_156 The mechanism described in CSM\_155 shall be triggered by the vehicle unit whenever a card is inserted into one of its card slots.

#### 10.2. **Mutual Certificate Chain Verification**

- 10.2.1 *Card Certificate Chain Verification by VU*
  - CSM\_157 **►M1** Vehicle units shall use the protocol depicted in Figure 4 for verifying a tachograph card's certificate chain. For every certificate it reads from the card, the VU shall verify that the Certificate Holder Authorisation (CHA) field is correct:
    - The CHA field of the Card certificate shall indicate a card certificate for mutual authentication (see Appendix 1, data type EquipmentType).
    - The CHA of the Card.CA certificate shall indicate an MSCA.
    - The CHA of the Card.Link certificate shall indicate the ERCA. ◄

*Notes to Figure 4:*

- The Card certificates and public keys mentioned in the figure are those for mutual authentication. Section 9.1.5 denotes these as Card\_MA.
- The Card.CA certificates and public keys mentioned in the figure are those for signing card certificates and it is indicated in the CAR of the Card certificate. Section 9.1.3 denotes these as MSCA\_Card.
- The Card.CA.EUR certificate mentioned in the figure is the European root certificate that is indicated in the CAR of the Card.CA certificate.
- The Card.Link certificate mentioned in the figure is the card's link certificate, if present. As specified in section 9.1.2, this is a link certificate for a new European root key pair created by the ERCA and signed by the previous European private key.

- The Card.Link.EUR certificate is the European root certificate that is indicated in the CAR of the Card.Link certificate.
- CSM\_158 As depicted in Figure 4, verification of the card's certificate chain shall begin upon card insertion. The vehicle unit shall read the card holder reference ( ) from EF ICC. The VU shall check if it knows the card, i.e., if it has successfully verified the card's certificate chain in the past and stored it for future reference. If it does, and the card certificate is still valid, the process continues with the verification of the VU certificate chain. Otherwise, the VU shall successively read from the card the MSCA\_Card certificate to be used for verifying the card certificate, the Card.CA. EUR certificate to be used for verifying the MSCA\_Card certificate, and possibly the link certificate, until it finds a certificate it knows or it can verify. If such a certificate is found, the VU shall use that certificate to verify the underlying card certificates it has read from the card. If successful, the process continues with the verification of the VU certificate chain. If not successful, the VU shall ignore the card.

*Note:* There are three ways in which the VU may know the Card.CA.EUR certificate:

- the Card.CA.EUR certificate is the same certificate as the VU's own EUR certificate;
- the Card.CA.EUR certificate precedes the VU's own EUR certificate and the VU contained this certificate already at issuance (see CSM\_81);
- the Card.CA.EUR certificate succeeds the VU's own EUR certificate and the VU received a link certificate in the past from another tachograph card, verified it and stored it for future reference.
- CSM\_159 As indicated in Figure 4, once the VU has verified the authenticity and validity of a previously unknown certificate, it may store this certificate for future reference, such that it does not need to verify that certificate's authenticity again if it is presented to the VU again. Instead of storing the entire certificate, a VU may choose to store only the contents of the Certificate Body, as specified in section 9.3.2. **►M1** Whereas storing of all other types of certificate is optional, it is mandatory for a VU to store a new link certificate presented by a card. ◄
- CSM\_160 The VU shall verify the temporal validity of any certificate read from the card or stored in its memory, and shall reject expired certificates. For verifying the temporal validity of a certificate presented by the card a VU shall use its internal clock.

# *Figure 4* **Protocol for Card Certificate Chain Verification by VU**

![](_page_50_Figure_2.jpeg)

10.2.2 *VU Certificate Chain Verification by Card*

- CSM\_161 **►M1** Tachograph cards shall use the protocol depicted in Figure 5 for verifying a VU's certificate chain. For every certificate presented by the VU, the card shall verify that the Certificate Holder Authorisation (CHA) field is correct:
  - The CHA of the VU.Link certificate shall indicate the ERCA.

- The CHA of the VU.CA certificate shall indicate an MSCA.
- The CHA field of the VU certificate shall indicate a VU certificate for mutual authentication (see Appendix 1, data type EquipmentType). ◄

#### *Figure 5*

![](_page_51_Figure_4.jpeg)

![](_page_51_Figure_5.jpeg)

*Notes to Figure 5:*

<sup>—</sup> The VU certificates and public keys mentioned in the figure are those for mutual authentication. Section 9.1.4 denotes these as VU\_MA.

- The VU.CA certificates and public keys mentioned in the figure are those for signing VU and external GNSS facility certificates. Section 9.1.3 denotes these as MSCA\_VU-EGF.
- The VU.CA.EUR certificate mentioned in the figure is the European root certificate that is indicated in the CAR of the VU.CA certificate.
- The VU.Link certificate mentioned in the figure is the VU's link certificate, if present. As specified in section 9.1.2, this is a link certificate for a new European root key pair created by the ERCA and signed by the previous European private key.
- The VU.Link.EUR certificate is the European root certificate that is indicated in the CAR of the VU.Link certificate.
- CSM\_162 As depicted in Figure 5, verification of the certificate chain of the vehicle unit shall begin with the vehicle unit attempting to set its own public key for use in the tachograph card. If this succeeds, it means that the card successfully verified the VU's certificate chain in the past, and has stored the VU certificate for future reference. In this case, the VU certificate is set for use and the process continues with VU Authentication. If the card does not know the VU certificate, the VU shall successively present the VU.CA certificate to be used for verifying its VU certificate, the VU.CA.EUR certificate to be used for verifying the VU.CA certificate, and possibly the link certificate, in order to find a certificate known or verifiable by the card. If such a certificate is found, the card shall use that certificate to verify the underlying VU certificates presented to it. If successful, the VU shall finally set its public key for use in the tachograph card. If not successful, the VU shall ignore the card.

*Note: There are three ways in which the card may know the VU.CA.EUR certificate:*

- the VU.CA.EUR certificate is the same certificate as the card's own EUR certificate;
- the VU.CA.EUR certificate precedes the card's own EUR certificate and the card contained this certificate already at issuance (see CSM\_91);
- the VU.CA.EUR certificate succeeds the card's own EUR certificate and the card received a link certificate in the past from another vehicle unit, verified it and stored it for future reference.

- CSM\_163 The VU shall use the MSE: Set AT command to set its public key for use in the tachograph card. As specified in Appendix 2, this command contains an indication of the cryptographic mechanism that will be used with the key that is set. This mechanism shall be 'VU Authentication using the ECDSA algorithm, in combination with the hashing algorithm linked to the key size of the VU's VU\_MA key pair, as specified in CSM\_50'.
- CSM\_164 The MSE: Set AT command also contains an indication of the ephemeral key pair which the VU will use during session key agreement (see section 10.4). Therefore, before sending the MSE: Set AT command, the VU shall generate an ephemeral ECC key pair. For generating the ephemeral key pair, the VU shall use the standardized domain parameters indicated in the card certificate. The ephemeral key pair is denoted as (VU.SKeph, VU.PKeph, Card.DP). The VU shall take the x-coordinate of the ECDH ephemeral public point as the key identification; this is called the compressed representation of the public key and denoted as Comp(VU.PKeph).
- **▼M1**
- CSM\_165 If the MSE: Set AT command is successful, the card shall set the indicated VU.PK for subsequent use during Vehicle Authentication, and shall temporarily store Comp(VU.PKeph). In case two or more successful MSE: Set AT commands are sent before session key agreement is performed, the card shall store only the last Comp(VU.PKeph) received. The card shall reset Comp(VU.PKeph) after a successful GENERAL AUTH-ENTICATE command.
- **▼B**
- CSM\_166 The card shall verify the temporal validity of any certificate presented by the VU or referenced by the VU while stored in the card's memory, and shall reject expired certificates.
- CSM\_167 For verifying the temporal validity of a certificate presented by the VU, each tachograph card shall internally store some data representing the current time. This data shall not be directly updatable by a VU. At issuance, the current time of a card shall be set equal to the Effective Date of the card's Card\_MA certificate. A card shall update its current time if the Effective Date of an authentic 'valid source of time' certificate presented by a VU is more recent than the card's current time. In that case, the card shall set its current time to the Effective Date of that certificate. The card shall accept only the following certificates as a valid source of time:

— Second-generation ERCA link certificates

— Second-generation MSCA certificates

— Second-generation VU certificates issued by the same country as the card's own card certificate(s).

*Note:* the last requirement implies that a card shall be able to recognize the CAR of the VU certificate, i.e. the MSCA\_VU-EGF certificate. This will not be the same as the CAR of its own certificate, which is the MSCA\_Card certificate.

CSM\_168 As indicated in Figure 5, once the card has verified the authenticity and validity of a previously unknown certificate, it may store this certificate for future reference, such that it does not need to verify that certificate's authenticity again if it is presented to the card again. Instead xof storing the entire certificate, a card may choose to store only the contents of the Certificate Body, as specified in section 9.3.2.

#### 10.3. **VU Authentication**

- CSM\_169 Vehicle units and cards shall use the VU Authentication protocol depicted in Figure 6 to authenticate the VU towards the card. VU Authentication enables the tachograph card to explicitly verify that the VU is authentic. To do so, the VU shall use its private key to sign a challenge generated by the card.
- CSM\_170 **►M1** Next to the card challenge, the VU shall include in the signature the certificate holder reference taken from the card certificate. ◄

*Note:* This ensures that the card to which the VU authenticates itself is the same card whose certificate chain the VU has verified previously.

CSM\_171 The VU shall also include in the signature the identifier of the ephemeral public key Comp(VU.PKeph) which the VU will use to set up Secure Messaging during the Chip Authentication process specified in section 10.4.

> *Note:* This ensures that the VU with which a card communicates during a Secure Messaging session is the same VU that was authenticated by the card.

### *Figure 6*

### **VU Authentication protocol**

![](_page_55_Figure_3.jpeg)

**▼B**

CSM\_172 If multiple GET CHALLENGE commands are sent by the VU during VU Authentication, the card shall return a new 8-byte random challenge each time, but shall store only the last challenge.

CSM\_173 The signing algorithm used by the VU for VU Authentication shall be ECDSA as specified in [DSS], using the hashing algorithm linked to the key size of the VU's VU\_MA key pair, as specified in CSM\_50. The signature format shall be plain, as specified in [TR-03111]. The VU shall send the resulting signature to the card.

# **M1**

- CSM\_174 Upon receiving the VU's signature in an EXTERNAL AUTHENTICATE command, the card shall
  - Calculate the authentication token by concatenating Card.CHR, the card challenge rcard and the identifier of the VU ephemeral public key Comp(VU.PKeph),
  - Verify the VU's signature using the ECDSA algorithm, using the hashing algorithm linked to the key size of the VU's VU\_MA key pair as specified in CSM\_50, in combination with VU.PK and the calculated authentication token.

#### 10.4. **Chip Authentication and Session Key Agreement**

CSM\_175 Vehicle units and cards shall use the Chip Authentication protocol depicted in **Figure 7** to authenticate the card towards the VU. Chip Authentication enables the vehicle unit to explicitly verify that the card is authentic.

#### *Figure 7*

#### **Chip Authentication and session key agreement**

![](_page_56_Figure_5.jpeg)

CSM\_176 The VU and the card shall take the following steps:

1. The vehicle unit initiates the Chip Authentication process by sending the MSE: Set AT command indicating 'Chip Authentication using the ECDH algorithm resulting in an AES session key length linked to the key size of the card's Card\_MA key pair, as specified in CSM\_50'. The VU shall determine the key size of the card's key pair from the card certificate.

**▼M1**

2. The VU sends the public point VU.PKeph of its ephemeral key pair to the card. The public point shall be converted to an octet string as specified in [TR-03111]. The uncompressed encoding format shall be used. As explained in CSM\_164, the VU generated this ephemeral key pair prior to the verification of the VU certificate chain. The VU sent the identifier of the ephemeral public key Comp(VU.PKeph) to the card, and the card stored it.

- 3. The card computes Comp(VU.PKeph) from VU.PKeph and compares this to the stored value of Comp(VU.PKeph).
- 4. Using the ECDH algorithm in combination with the card's static private key and the VU's ephemeral public key, the card computes a secret K.
- 5. The card chooses a random 8-byte nonce NPICC and uses it to derive two AES session keys KMAC and KENC from K. See CSM\_179.
- 6. Using KMAC, the card computes an authentication token over the VU ephemeral public point: TPICC = CMAC(KMAC, VU.PKeph). The public point shall be in the format used by the VU (see bullet 2 above). The card sends NPICC and TPICC to the vehicle unit.
- 7. Using the ECDH algorithm in combination with the card's static public key and the VU's ephemeral private key, the VU computes the same secret K as the card did in step 4.
- 8. The VU derives session keys KMAC and KENC from K and NPICC; see CSM\_179.
- 9. The VU verifies the authentication token TPICC.
- CSM\_177 In step 3 above, the card shall compute Comp(VU.PKeph) as the x-coordinate of the public point in VU.PKeph.
- CSM\_178 In steps 4 and 7 above, the card and the vehicle unit shall use the ECKA-EG algorithm as defined in [TR-03111].
- CSM\_179 In steps 5 and 8 above, the card and the vehicle unit shall use the key derivation function for AES session keys defined in [TR-03111], with the following precisions and changes:
  - The value of the counter shall be '00 00 00 01' for KENC and '00 00 00 02' for KMAC.
  - The optional nonce *r* shall be used and shall be equal to NPICC.
  - For deriving 128-bits AES keys, the hashing algorithm to be used shall be SHA-256.
  - For deriving 192-bits AES keys, the hashing algorithm to be used shall be SHA-384.
  - For deriving 256-bits AES keys, the hashing algorithm to be used shall be SHA-512.

The length of the session keys (i.e. the length at which the hash is truncated) shall be linked to the size of the Card\_MA key pair, as specified in CSM\_50.

# **B**

**▼M1**

CSM\_180 In steps 6 and 9 above, the card and the vehicle unit shall use the AES algorithm in CMAC mode, as specified in [SP 800-38B]. The length of TPICC shall be linked to the length of the AES session keys, as specified in CSM\_50.

#### 10.5. **Secure Messaging**

- 10.5.1 *General*
  - CSM\_181 All commands and responses exchanged between a vehicle unit and a tachograph card after successful Chip Authentication took place and until the end of the session shall be protected by Secure Messaging.
  - CSM\_182 Except when reading from a file with access condition SM-R-ENC-MAC-G2 (see Appendix 2, section 4), Secure Messaging shall be used in authentication-only mode. In this mode, a cryptographic checksum (a.k.a. MAC) is added to all commands and responses to ensure message authenticity and integrity.
  - CSM\_183 When reading data from a file with access condition SM-R-ENC-MAC-G2, Secure Messaging shall be used in encrypt-then-authenticate mode, i.e. the response data is encrypted first to ensure message confidentiality, and afterwards a MAC over the formatted encrypted data is calculated to ensure authenticity and integrity.
  - CSM\_184 Secure Messaging shall use AES as defined in [AES] with the session keys KMAC and KENC that were agreed during Chip Authentication.
  - CSM\_185 An unsigned integer shall be used as the Send Sequence Counter (SSC) to prevent replay attacks. The size of the SSC shall be equal to the AES block size, i.e. 128 bits. The SSC shall be in MSB-first format. The Send Sequence Counter shall be initialized to zero (i.e. '00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00') when Secure Messaging is started. The SSC shall be increased every time before a command or response APDU is generated, i.e. since the starting value of the SSC in a SM session is 0, in the first command the value of the SSC will be 1. The value of SSC for the first response will be 2.
  - CSM\_186 For message encryption, KENC shall be used with AES in the Cipher Block Chaining (CBC) mode of operation, as defined in [ISO 10116], with an interleave parameter *m* = 1 and an initialization vector SV = E(KENC, SSC), i.e. the current value of the Send Sequence Counter encrypted with KENC.
  - CSM\_187 For message authentication, KMAC shall be used with AES in CMAC mode as specified in [SP 800-38B]. The length of the MAC shall be linked to the length of the AES session keys, as specified in CSM\_50. The Send Sequence Counter shall be included in the MAC by prepending it before the datagram to be authenticated.

10.5.2 *Secure Message Structure*

CSM\_188 Secure Messaging shall make use only of the Secure Messaging data objects (see [ISO 7816-4]) listed in Table 5. In any message, these data objects shall be used in the order specified in this table.

### *Table 5*

### **Secure Messaging Data Objects**

| Data Object Name                                                                           | Tag  | Presence (M) andatory, (C)<br>onditional or (F) orbidden<br>in |           |
|--------------------------------------------------------------------------------------------|------|----------------------------------------------------------------|-----------|
|                                                                                            |      | Commands                                                       | Responses |
| Plain value not encoded in BER-TLV                                                         | '81' | C                                                              | C         |
| Plain value encoded in BER-TLV, but<br>not including SM DOs                                | 'B3' | C                                                              | C         |
| Padding-content indicator followed by<br>cryptogram, plain value not encoded<br>in BER-TLV | '87' | C                                                              | C         |
| Protected Le                                                                               | '97' | C                                                              | F         |
| Processing Status                                                                          | '99' | F                                                              | M         |
| Cryptographic Checksum                                                                     | '8E' | M                                                              | M         |

*Note:* As specified in Appendix 2, tachograph cards may support the READ BINARY and UPDATE BINARY command with an odd INS byte ('B1' resp. 'D7'). These command variants are required to read and update files with more than 32 768 bytes or more. In case such a variant is used, a data object with tag 'B3' shall be used instead of an object with tag '81'. See Appendix 2 for more information.

- CSM\_189 All SM data objects shall be encoded in DER TLV as specified in [ISO 8825-1]. This encoding results in a Tag-Length-Value (TLV) structure as follows:
  - Tag: The tag is encoded in one or two octets and indicates the content.
  - Length: The length is encoded as an unsigned integer in one, two, or three octets, resulting in a maximum length of 65 535 octets. The minimum number of octets shall be used.
  - Value: The value is encoded in zero or more octets
- CSM\_190 APDUs protected by Secure Messaging shall be created as follows:
  - The command header shall be included in the MAC calculation, therefore value '0C'shall be used for the class byte CLA.

- As specified in Appendix 2, all INS bytes shall be even, with the possible exception of odd INS bytes for the READ BINARY and UPDATE BINARY commands.
- The actual value of Lc will be modified to Lc' after application of secure messaging.
- The Data field shall consist of SM data objects.
- In the protected command APDU the new Le byte shall be set to '00'. If required, a data object '97' shall be included in the Data field in order to convey the original value of Le.
- CSM\_191 Any data object to be encrypted shall be padded according to [ISO 7816-4] using padding-content indicator '01'. For the calculation of the MAC, data objects in the APDU shall be padded according to [ISO 7816-4].

*Note:* Padding for Secure Messaging is always performed by the secure messaging layer, not by the CMAC or CBC algorithms.

#### *Summary and Examples*

A command APDU with applied Secure Messaging will have the following structure, depending on the case of the respective unsecured command (DO is data object):

| Case 1:                                                                                                                                      | CLA INS P1 P2    Lc'    DO '8E'    Le                     |  |
|----------------------------------------------------------------------------------------------------------------------------------------------|-----------------------------------------------------------|--|
| Case 2:                                                                                                                                      | CLA INS P1 P2    Lc'    DO '97'    DO'8E'    Le           |  |
| Case 3 (even INS byte):                                                                                                                      | CLA INS P1 P2    Lc'    DO '81'    DO'8E'    Le           |  |
| Case 3 (odd INS byte):                                                                                                                       | CLA INS P1 P2    Lc'    DO 'B3'    DO'8E'    Le           |  |
| Case 4 (even INS byte):                                                                                                                      | CLA INS P1 P2    Lc'    DO '81'    DO'97'    DO'8E'    Le |  |
| Case 4 (odd INS byte):                                                                                                                       | CLA INS P1 P2    Lc'    DO 'B3'    DO'97'    DO'8E'    Le |  |
| where Le = '00' or '00 00' depending on whether short length fields or extended length fields are used; see [ISO 7816-4].                    |                                                           |  |
| A response APDU with applied Secure Messaging will have the following structure, depending on the case of the respective unsecured response: |                                                           |  |
| Case 1 or 3:                                                                                                                                 | DO '99'    DO '8E'    SW1SW2                              |  |
| Case 2 or 4 (even INS byte) without encryption:                                                                                              | DO '81'    DO '99'    DO '8E'    SW1SW2                   |  |
| Case 2 or 4 (even INS byte) with encryption:                                                                                                 | DO '87'    DO '99'    DO '8E'    SW1SW2                   |  |
| Case 2 or 4 (odd INS byte) without encryption:                                                                                               | DO 'B3'    DO '99'    DO '8E'    SW1SW2                   |  |

*Note:* Case 2 or 4 (odd INS byte) with encryption is never used in the communication between a VU and a card.

### **B**

Below are three example APDU transformations for commands with even INS code. Figure 8 shows an authenticated Case 4 command APDU, Figure 9 shows an authenticated Case 1/Case 3 response APDU, and Figure 10 shows an encrypted and authenticated Case 2/Case 4 response APDU.

#### *Figure 8*

#### **Transformation of an authenticated Case 4 Command APDU**

![](_page_61_Figure_4.jpeg)

*Figure 9*

**Transformation of an authenticated Case 1 / Case 3 Response APDU**

![](_page_61_Figure_7.jpeg)

### *Figure 10*

### **Transformation of an encrypted and authenticated Case 2/Case 4 Response APDU**

![](_page_62_Figure_3.jpeg)

### **B**

10.5.3 *Secure Messaging Session Abortion*

CSM\_192 A vehicle unit shall abort an ongoing Secure Messaging session if and only if one of the following conditions occur:

— it receives a plain response APDU,

— it detects a Secure Messaging error in a response APDU:

- An expected Secure Messaging data object is missing, the order of data objects is incorrect, or an unknown data object is included.
- A Secure Messaging data object is incorrect, e.g. the MAC value is incorrect, the TLV structure is incorrect or the padding indicator in tag '87' is not equal to '01'.
- the card sends a status byte indicating it detected an SM error (see CSM\_194),
- the limit for the number of commands and associated responses within the current session is reached. For a given VU, this limit shall be defined by its manufacturer, taking into account the security requirements of the hardware used, with a maximum value of 240 SM commands and associated responses per session.

- CSM\_193 A tachograph card shall abort an ongoing Secure Messaging session if and only if one of the following conditions occur:
  - it receives a plain command APDU,
  - it detects a Secure Messaging error in a command APDU:
    - An expected Secure Messaging data object is missing, the order of data objects is incorrect, or an unknown data object is included.
    - A Secure Messaging data object is incorrect, e.g. the MAC value is incorrect or the TLV structure is incorrect.
  - it is depowered or reset,
  - the VU starts the VU Authentication process,
  - the limit for the number of commands and associated responses within the current session is reached. For a given card, this limit shall be defined by its manufacturer, taking into account the security requirements of the hardware used, with a maximum value of 240 SM commands and associated responses per session.

# **B**

- CSM\_194 Regarding SM error handling by a tachograph card:
  - If in a command APDU some expected Secure Messaging data objects are missing, the order of data objects is incorrect or unknown data objects are included, a tachograph card shall respond with status bytes '69 87'.
  - If a Secure Messaging data object in a command APDU is incorrect, a tachograph card shall respond with status bytes '69 88'.

In such a case, the status bytes shall be returned without using SM.

CSM\_195 If a Secure Messaging session between a VU and a tachograph card is aborted, the VU and the tachograph card shall

— securely destroy the stored session keys

- immediately establish a new Secure Messaging session, as described in sections 10.2 — 10.5.
- CSM\_196 If for any reason the VU decides to restart mutual authentication towards an inserted card, the process shall restart with verification of the card certificate chain, as described in section 10.2, and shall continue as described in sections 10.2 — 10.5.

### 11. VU EXTERNAL GNSS FACILITY COUPLING, MUTUAL AUTH-ENTICATION AND SECURE MESSAGING

- 11.1. **General**
  - CSM\_197 The GNSS facility used by a VU to determine its position may be internal, (i.e. built into the VU casing and not detachable), or it may be an external module. In the first case, there is no need to standardize the internal communication between the GNSS facility and the VU, and the requirements in this chapter do not apply. In the latter case, communication between the VU and the external GNSS facility shall be standardized and protected as described in this chapter.
  - CSM\_198 Secure communication between a vehicle unit and an external GNSS facility shall take place in the same way as secure communication between a vehicle unit and a tachograph card, with the external GNSS facility (EGF) taking the role of the card. All requirements mentioned in chapter 10 for tachograph cards shall be satisfied by an EGF, taking into account the deviations, clarifications and additions mentioned in this chapter. In particular, mutual certificate chain verification, VU Authentication and Chip Authentication shall be performed as described in sections 11.3 and 11.4.
  - CSM\_199 Communication between a vehicle unit and an EGF differs from communication between a vehicle unit and a card in the fact that a vehicle unit and an EGF must be coupled once in a workshop before the VU and the EGF can exchange GNSS-based data during normal operation. The coupling process is described in section 11.2.
  - CSM\_200 For communication between a vehicle unit and an EGF, APDU commands and responses based on [ISO 7816-4] and [ISO 7816-8] shall be used. The exact structure of these APDUs is defined in Appendix 2 of this Annex.

### 11.2. **VU and External GNSS Facility Coupling**

- CSM\_201 A vehicle unit and an EGF in a vehicle shall be coupled by a workshop. Only a coupled vehicle unit and EGF shall be able to communicate during normal operation.
- CSM\_202 Coupling of a vehicle unit and an EGF shall only be possible if the vehicle unit is in calibration mode. The coupling shall be initiated by the vehicle unit.
- CSM\_203 A workshop may re-couple a vehicle unit to another EGF or to the same EGF at any time. During re-coupling, the VU shall securely destroy the existing EGF\_MA certificate in its memory and shall store the EGF\_MA certificate of the EGF to which it is being coupled.
- CSM\_204 A workshop may re-couple an external GNSS facility to another VU or to the same VU at any time. During re-coupling, the EGF shall securely destroy the existing VU\_MA certificate in its memory and shall store the VU\_MA certificate of the VU to which it is being coupled.

#### 11.3. **Mutual Certificate Chain Verification**

- 11.3.1 *General*
  - CSM\_205 Mutual certificate chain verification between a VU and an EGF shall take place only during the coupling of the VU and the EGF by a workshop. During normal operation of a coupled VU and EGF, no certificates shall be verified. Instead, the VU and EGF shall trust the certificates they stored during the coupling, after checking the temporal validity of these certificates. The VU and the EGF shall not trust any other certificates for protecting the VU — EGF communication during normal operation.
- 11.3.2 *During VU EGF Coupling*
  - CSM\_206 During the coupling to an EGF, a vehicle unit shall use the protocol depicted in Figure 4 (section 10.2.1) for verifying the external GNSS facility's certificate chain.

*Notes to Figure 4 within this context:*

- Communication control is out of the scope of this Appendix. However, an EGF is not a smart card and hence the VU will probably not send a Reset to initiate the communication and will not receive an ATR.
- The Card certificates and public keys mentioned in the figure shall be interpreted as the EGF's certificates and public keys for mutual authentication. Section 9.1.6 denotes these as EGF\_MA.
- The Card.CA certificates and public keys mentioned in the figure shall be interpreted as the MSCA's certificates and public keys for signing EGF certificates. Section 9.1.3 denotes these as MSCA\_VU-EGF.
- The Card.CA.EUR certificate mentioned in the figure shall be interpreted as the European root certificate that is indicated in the CAR of the MSCA\_VU-EGF certificate.
- The Card.Link certificate mentioned in the figure shall be interpreted as the EGF's link certificate, if present. As specified in section 9.1.2, this is a link certificate for a new European root key pair created by the ERCA and signed by the previous European private key.
- The Card.Link.EUR certificate is the European root certificate that is indicated in the CAR of the Card.Link certificate.
- Instead of the , the VU shall read the from EF ICC.
- Instead of selecting the Tachograph AID, the VU shall select the EGF AID.

— 'Ignore Card' shall be interpreted as 'Ignore EGF'.

- CSM\_207 Once it has verified the EGF\_MA certificate, the vehicle unit shall store this certificate for use during normal operation; see section 11.3.3.
- CSM\_208 **►M1** During the coupling to a VU, an external GNSS facility shall use the protocol depicted in Figure 5 (section 10.2.2) for verifying the VU's certificate chain. ◄
  - *Notes to Figure 5 within this context:*
  - The VU shall generate a fresh ephemeral key pair using the domain parameters in the EGF certificate.
  - The VU certificates and public keys mentioned in the figure are those for mutual authentication. Section 9.1.4 denotes these as VU\_MA.
  - The VU.CA certificates and public keys mentioned in the figure are those for signing VU and external GNSS facility certificates. Section 9.1.3 denotes these as MSCA\_VU-EGF.
  - The VU.CA.EUR certificate mentioned in the figure is the European root certificate that is indicated in the CAR of the VU.CA certificate.
  - The VU.Link certificate mentioned in the figure is the VU's link certificate, if present. As specified in section 9.1.2, this is a link certificate for a new European root key pair created by the ERCA and signed by the previous European private key.
  - The VU.Link.EUR certificate is the European root certificate that is indicated in the CAR of the VU.Link certificate.
- CSM\_209 In deviation from requirement CSM\_167, an EGF shall use the GNSS time to verify the temporal validity of any certificate presented.

# **M1**

CSM\_210 Once it has verified the VU\_MA certificate, the external GNSS facility shall store this certificate for use during normal operation; see section 11.3.3.

# **B**

### 11.3.3 *During Normal Operation*

CSM\_211 **►M1** During normal operation, a vehicle unit and an EGF shall use the protocol depicted in Figure 11 for verifying the temporal validity of the stored EGF\_MA certificate and for setting the VU\_MA public key for subsequent VU Authentication. No further mutual verification of the certificate chains shall take place during normal operation. ◄

> Note that Figure 11 in essence consists of the first steps shown in Figure 4 and Figure 5. Again, note that since an EGF is not a smart card, the VU will probably not send a Reset to initiate the communication and will not receive an ATR. In any case this is out of the scope of this Appendix.

*Figure 11*

### **Mutual verification of certificate temporal validity during normal VU EGF operation**

![](_page_67_Figure_3.jpeg)

CSM\_212 As shown in Figure 11, the vehicle unit shall log an error if the EGF\_MA certificate is no longer valid. However, mutual authentication, key agreement and subsequent communication via secure messaging shall proceed normally.

### 11.4. **VU Authentication, Chip Authentication and Session Key Agreement**

CSM\_213 VU Authentication, Chip Authentication and session key agreement between a VU and an EGF shall take place during coupling and whenever a Secure Messaging session is re-established during normal operation. The VU and the EGF shall carry out the processes described in sections 10.3 and 10.4. All requirements in these sections shall apply.

### 11.5. **Secure Messaging**

- CSM\_214 All commands and responses exchanged between a vehicle unit and an external GNSS facility after successful Chip Authentication took place and until the end of the session shall be protected by Secure Messaging.in authentication-only mode. All requirements in section 10.5 shall apply.
- CSM\_215 If a Secure Messaging session between a VU and an EGF is aborted, the VU shall immediately establish a new Secure Messaging session, as described in section 11.3.3 and 11.4.
- 12. VU MOTION SENSOR PAIRING AND COMMUNICATION
- 12.1. **General**
  - CSM\_216 A vehicle unit and a motion sensor shall communicate using the interface protocol specified in [ISO 16844-3] during pairing and in normal operation, with the changes described in this chapter and in section 9.2.1.

*Note:* readers of this chapter are supposed to be familiar with the contents of [ISO 16844-3].

#### 12.2. **VU — Motion Sensor Pairing Using Different Key Generations**

As explained in section 9.2.1, the motion sensor master key and all associated keys are regularly replaced. This leads to the presence of up to three motion sensor-related AES keys KM-WC (of consecutive key generations) in workshop cards. Similarly, in motion sensors up to three different AES-based encryptions of data (based on consecutive generations of the motion sensor master key KM) may be present. A vehicle unit contains only one motion sensor-related key KM-VU.

- CSM\_217 A second-generation VU and a second-generation motion sensor shall be paired as follows (compare Table 6 in [ISO 16844-3]):
  - 1. A second-generation workshop card is inserted into the VU and the VU is connected to the motion sensor.
  - 2. The VU reads all available KM-WC keys from the workshop card, inspects their key version numbers and chooses the one matching the version number of the VU's KM-VU key. If the matching KM-WC key is not present on the workshop card, the VU aborts the pairing process and shows an appropriate error message to the workshop card holder.
  - 3. The VU calculates the motion sensor master key KM from KM-VU and KM-WC, and the identification key KID from KM, as specified in section 9.2.1.
  - 4. The VU sends the instruction to initiate the pairing process towards the motion sensor, as described in [ISO 16844-3], and encrypts the serial number it receives from the motion sensor with the identification key KID. The VU sends the encrypted serial number back to the motion sensor.
  - 5. The motion sensor matches the encrypted serial number consecutively with each of the encryptions of the serial number it holds internally. If it finds a match, the VU is authenticated. The motion sensor notes the generation of KID used by the VU and returns the matching encrypted version of its pairing key; i.e. the encryption that was created using the same generation of KM.
  - 6. The VU decrypts the pairing key using KM, generates a session key KS, encrypts it with the pairing key and sends the result to the motion sensor. The motion sensor decrypts KS.
  - 7. The VU assembles the pairing information as defined in [ISO 16844-3], encrypts the information with the pairing key, and sends the result to the motion sensor. The motion sensor decrypts the pairing information.
  - 8. The motion sensor encrypts the received pairing information with the received KS and returns this to the VU. The VU verifies that the pairing information is the same information which the VU sent to the motion sensor

in the previous step. If it is, this proves that the motion sensor used the same KS as the VU and hence in step 5 sent its pairing key encrypted with the correct generation of KM. Hence, the motion sensor is authenticated.

Note that steps 2 and 5 are different from the standard process in [ISO 16844-3]; the other steps are standard.

*Example:* Suppose a pairing takes place in the first year of the validity of the ERCA (3) certificate; see Figure 2 in section 9.2.1.2. Moreover

- Suppose the motion sensor was issued in the last year of the validity of the ERCA (1) certificate. It will therefore contain the following keys and data:
  - Ns[1]: its serial number encrypted with generation 1 of KID,
  - Ns[2]: its serial number encrypted with generation 2 of KID,
  - Ns[3]: its serial number encrypted with generation 3 of KID,
  - KP[1]: its generation-1 pairing key (1), encrypted with generation 1 of KM,
  - KP[2]: its generation-2 pairing key, encrypted with generation 2 of KM,
  - KP[3]: its generation-3 pairing key, encrypted with generation 3 of KM,
- Suppose that the workshop card was issued in the first year of the validity of the ERCA (3) certificate. It will therefore contain the generation 2 and generation 3 of the KM-WC key.
- Suppose the VU is a generation-2 VU, containing the generation 2 of KM-VU.
- In this case, the following will happen in steps 2 5:
- Step 2: The VU reads generation 2 and generation 3 of KM-WC from the workshop card and inspects their version numbers.
- Step 3: The VU combines the generation-2 KM-WC with its KM-VU to compute KM and KID.
- Step 4: The VU encrypts the serial number it receives from the motion sensor with KID.
- Step 5: The motion sensor compares the received data with Ns[1] and doesn't find a match. Next, it compares the data with Ns[2] and finds a match. It concludes that the VU is a generation-2 VU, and therefore sends back KP[2].

<sup>(1)</sup> Note that the generation-1, generation-2 and generation-3 pairing keys may actually be the same key, or may be three different keys having different lengths, as explained in CSM\_117.

#### 12.3. **VU — Motion Sensor Pairing and Communication using AES**

CSM\_218 As specified in Table 3 in section 9.2.1, all keys involved in the pairing of a (second-generation) vehicle unit and a motion sensor and in subsequent communication shall be AES keys, rather than double-length TDES keys as specified in [ISO 16844-3]. These AES keys may have a length of 128, 192 or 256 bits. Since the AES block size is 16 bytes, the length of an encrypted message must be a multiple of 16 bytes, compared to 8 bytes for TDES. Moreover, some of these messages will be used to transport AES keys, the length of which may be 128, 192 or 256 bits. Therefore, the number of data bytes per instruction in Table 5 of [ISO 16844-3] shall be changed as shown in Table 6:

**▼M1**

#### *Table 6*

#### **Number of plaintext and encrypted data bytes per instruction defined in [ISO 16844-3]**

| Instruction | Request / reply | Description of data                    | # of plaintext data<br>bytes according to<br>[ISO 16844-3] | # of plaintext data<br>bytes using AES<br>keys | # of encrypted data bytes when<br>using AES keys of bitlength |         |         |
|-------------|-----------------|----------------------------------------|------------------------------------------------------------|------------------------------------------------|---------------------------------------------------------------|---------|---------|
|             |                 |                                        |                                                            |                                                | 128                                                           | 192     | 256     |
| 10          | request         | Authentication data +<br>file number   | 8                                                          | 8                                              | 16                                                            | 16      | 16      |
| 11          | reply           | Authentication data +<br>file contents | 16 or 32, depend<br>on file                                | 16 or 32, depend on<br>file                    | 32 / 48                                                       | 32 / 48 | 32 / 48 |
| 41          | request         | MoS serial number                      | 8                                                          | 8                                              | 16                                                            | 16      | 16      |
| 41          | reply           | Pairing key                            | 16                                                         | 16 / 24 / 32                                   | 16                                                            | 32      | 32      |
| 42          | request         | Session key                            | 16                                                         | 16 / 24 / 32                                   | 16                                                            | 32      | 32      |
| 43          | request         | Pairing information                    | 24                                                         | 24                                             | 32                                                            | 32      | 32      |
| 50          | reply           | Pairing information                    | 24                                                         | 24                                             | 32                                                            | 32      | 32      |
| 70          | request         | Authentication data                    | 8                                                          | 8                                              | 16                                                            | 16      | 16      |
| 80          | reply           | MoS counter value +<br>auth. data      | 8                                                          | 8                                              | 16                                                            | 16      | 16      |

**▼B**

CSM\_219 The pairing information that is sent in instructions 43 (VU request) and 50 (MoS reply) shall be assembled as specified in section 7.6.10 of [ISO 16844-3], except that the AES algorithm shall be used instead of the TDES algorithm in the pairing data encryption scheme, thus resulting in two AES encryptions, and adopting the padding specified in CSM\_220 to fit with the AES block size. The key K'p used for this encryption shall be generated as follows:

> — In case the pairing key KP is 16 bytes long: K'p = KP XOR (Ns||Ns)

- In case the pairing key KP is 24 bytes long: K'p = KP XOR (Ns||Ns||Ns)
- In case the pairing key KP is 32 bytes long: K'p = KP XOR (Ns||Ns||Ns||Ns)

where Ns is the 8-byte serial number of the motion sensor.

CSM\_220 In case the plaintext data length (using AES keys) is not a multiple of 16 bytes, padding method 2 defined in [ISO 9797- 1] shall be used.

> *Note:* in [ISO 16844-3], the number of plaintext data bytes is always a multiple of 8, such that padding is not necessary when using TDES. The definition of data and messages in [ISO 16844-3] is not changed by this part of this Appendix, thus necessitating the application of padding.

- CSM\_221 For instruction 11 and in case more than one block of data must be encrypted, the Cipher Block Chaining mode of operation shall be used as defined in [ISO 10116], with an interleave parameter *m* = 1. The IV to be used shall be
  - For instruction 11: the 8-byte authentication block specified in section 7.6.3.3 of [ISO 16844-3], padded using padding method 2 defined in [ISO 9797-1]; see also section 7.6.5 and 7.6.6 of [ISO 16844-3].
  - For all other instructions in which more than 16 bytes are transferred, as specified in Table 6: '00' {16}, i.e. sixteen bytes with binary value 0.

*Note:* As shown in section 7.6.5 and 7.6.6 of [ISO 16844-3], when the MoS encrypts data files for inclusion in instruction 11, the authentication block is both

- Used as the initialization vector for the CBC-mode encryption of the data files
- Encrypted and included as the first block in the data that is sent to the VU.

#### 12.4. **VU — Motion Sensor Pairing For Different Equipment Generations**

CSM\_222 As explained in section 9.2.1, a second-generation motion sensor may contain the TDES-based encryption of the pairing data (as defined in Part A of this Appendix), which allows the motion sensor to be paired to a first-generation VU. If this is the case, a first-generation VU and a secondgeneration motion sensor shall be paired as described in Part A of this Appendix and in [ISO 16844-3]. For the pairing process either a first-generation or a second-generation workshop card may be used.

*Notes:*

— It is not possible to pair a second-generation VU to a firstgeneration motion sensor.

— It is not possible to use a first-generation workshop card for coupling a second-generation VU to a motion sensor.

#### 13. SECURITY FOR REMOTE COMMUNICATION OVER DSRC

#### 13.1. **General**

As specified in Appendix 14, a VU regularly generates Remote Tachograph Monitoring (RTM) data and sends this data to the (internal or external) Remote Communication Facility (RCF). The remote communication facility is responsible for sending this data over the DSRC interface described in Appendix 14 to the remote interrogator. Appendix 1 specifies that the RTM data is the concatenation of:

**Encrypted tachograph payload** the encryption of the plaintext tachograph payload

**DSRC security data** described below

The plaintext tachograph payload data format is specified in Appendix 1 and further described in Appendix 14. This section describes the structure of the DSRC security data; the formal specification is in Appendix 1.

- CSM\_223 The plaintext data communicated by a VU to a Remote Communication Facility (if the RCF is external to the VU) or from the VU to a remote interrogator over the DSRC interface (if the RCF is internal in the VU) shall be protected in encrypt-then-authenticate mode, i.e. the tachograph payload data is encrypted first to ensure message confidentiality, and afterwards a MAC is calculated to ensure data authenticity and integrity.
- CSM\_224 The DSRC security data shall consist of the concatenation of the following data elements in the following order; see also Figure 12:

| Current date time              | the current date and time of the VU (data type TimeReal)                                                               |
|--------------------------------|------------------------------------------------------------------------------------------------------------------------|
| Counter                        | a 3-byte counter, see CSM_225                                                                                          |
| VU serial number               | the VU's serial number or certificate request ID (data type VuSerialNumber or CertificateRequestID) – see CSM_123      |
| DSRC master key version number | the 1-byte version number of the DSRC master key from which the VU-specific DSRC keys were derived, see section 9.2.2. |
| MAC                            | the MAC calculated over all previous bytes in the RTM data.                                                            |

CSM\_225 The 3-byte counter in the DSRC security data shall be in MSB-first format. The first time a VU calculates a set of RTM data after it is taken into production, it shall set the value of the counter to 0. The VU shall increase the value of the counter data by 1, each time before it calculates a next set of RTM data.

# **B**

**▼M1**

#### 13.2. **Tachograph Payload Encryption and MAC Generation**

- CSM\_226 Given a plaintext data element with data type as described in Appendix 14, a VU shall encrypt this data as shown in Figure 12: the VU's DSRC key for encryption K\_VUDSRC\_ENC (see section 9.2.2) shall be used with AES in the Cipher Block Chaining (CBC) mode of operation, as defined in [ISO 10116], with an interleave parameter *m* = 1. The initialization vector shall be equal to *IV* = *current date time* || '*00 00 00 00 00 00 00 00 00*' || *counter*, where *current date time* and *counter* are specified in CSM\_224. The data to be encrypted shall be padded using method 2 defined in [ISO 9797-1].
- CSM\_227 A VU shall calculate the MAC in the DSRC security data as shown in Figure 12: the MAC shall be calculated over all preceding bytes in the RTM data, up to and including the DSRC master key version number, and including the tags and lengths of the data objects. The VU shall use its DSRC key for authenticity K\_VUDSRC\_MAC (see section 9.2.2) with the AES algorithm in CMAC mode as specified in [SP 800-38B]. The length of the MAC shall be linked to the length of the VU-specific DSRC keys, as specified in CSM\_50.

#### *Figure 12*

# **Tachograph payload encryption and MAC generation**

![](_page_73_Figure_6.jpeg)

### 13.3. **Verification and Decryption of Tachograph Payload**

CSM\_228 When a remote interrogator receives RTM data from a VU, it shall send the entire RTM data to a control card in the data field of a PROCESS DSRC MESSAGE command, as described in Appendix 2. Then:

- 1. The control card shall inspect the DSRC master key version number in the DSRC security data. If the control card does not know the indicated DSRC master key, it shall return an error specified in Appendix 2 and abort the process.
- 2. The control card shall use the indicated DSRC master key in combination with the VU serial number or the certificate request ID in the DSRC security data to derive the VU-specific DSRC keys K\_VUDSRC\_ENC and K\_VUDSRC\_MAC, as specified in CSM\_124.
- 3. The control card shall use K\_VUDSRC\_MAC to verify the MAC in the DSRC security data, as specified in CSM\_227. If the MAC is incorrect, the control card shall return an error specified in Appendix 2 and abort the process.
- 4. The control card shall use K\_VUDSRC\_ENC to decrypt the encrypted tachograph payload, as specified in CSM\_226. The control card shall remove the padding and shall return the decrypted tachograph payload data to the remote interrogator.
- CSM\_229 In order to prevent replay attacks, the remote interrogator shall verify the freshness of the RTM data by verifying that the *current date time* in the DSRC security data does not deviate too much from the current time of the remote interrogator.

*Notes:*

- This requires the remote interrogator to have an accurate and reliable source of time.
- Since Appendix 14 requires a VU to calculate a new set of RTM data every 60 seconds, and the clock of the VU is allowed to deviate 1 minute from the real time, a lower limit for the freshness of the RTM data is 2 minutes. The actual freshness to be required also depends on the accuracy of the clock of the remote interrogator.
- CSM\_230 When a workshop verifies the correct functioning of the DSRC functionality of a VU, it shall send the entire RTM data received from the VU to a workshop card in the data field of a PROCESS DSRC MESSAGE command, as described in Appendix 2. The workshop card shall perform all checks and actions specified in CSM\_228.

#### 14. SIGNING DATA DOWNLOADS AND VERIFYING SIGNATURES

### 14.1. **General**

CSM\_231 The Intelligent Dedicated Equipment (IDE) shall store data received from a VU or a card during one download session within one physical data file. Data may be stored on an ESM (external storage medium). This file contains digital signatures over data blocks, as specified in Appendix 7. This file shall also contain the following certificates (refer to section 9.1):

# **B**

**▼M1**

— In case of a VU download:

— The VU\_Sign certificate

- The MSCA\_VU-EGF certificate containing the public key to be used for verification of the VU\_Sign certificate
- In case of a Card download:
  - The Card\_Sign certificate
  - The MSCA\_Card certificate containing the public key to be used for verification of the Card\_Sign certificate

CSM\_232 The IDE shall also dispose of.

- In case it uses a control card to verify the signature, as shown in Figure 13: The link certificate linking the latest EUR certificate to the EUR certificate whose validity period directly precedes it, if existing.
- In case it verifies the signature itself: all valid European root certificates.

*Note:* the method the IDE uses to retrieve these certificates is not specified in this Appendix.

#### 14.2. **Signature generation**

CSM\_233 The signing algorithm to create digital signatures over downloaded data shall be ECDSA as specified in [DSS], using the hashing algorithm linked to the key size of the VU or the card, as specified in CSM\_50. The signature format shall be plain, as specified in [TR-03111].

#### 14.3. **Signature verification**

- CSM\_234 **►M1** An IDE may perform verification of a signature over downloaded data itself or it may use a control card for this purpose. In case it uses a control card, signature verification shall take place as shown in Figure 13. For verifying the temporal validity of a certificate presented by the IDE, the control card shall use its internal current time, as specified in CSM\_167. The control card shall update its current time if the Effective Date of an authentic 'valid source of time' certificate is more recent than the card's current time. The card shall accept only the following certificates as a valid source of time:
  - Second-generation ERCA link certificates
  - Second-generation MSCA certificates
  - Second-generation VU\_Sign or Card\_Sign certificates issued by the same country as the control card's own card certificate.

In case it performs signature verification itself, the IDE shall verify the authenticity and validity of all certificates in the certificate chain in the data file, and it shall verify the signature over the data following the signature scheme defined in [DSS]. In both cases, for every certificate read from the data file, it is necessary to verify that the Certificate Holder Authorisation (CHA) field is correct:

— The CHA field of the EQT certificate shall indicate a VU or Card (as applicable) certificate for signing (see Appendix 1, data type EquipmentType).

- The CHA of the EQT.CA certificate shall indicate an MSCA.
- The CHA of the EQT.Link certificate shall indicate the ERCA. ◄

*Notes to Figure 13:*

- The equipment that signed the data to be analysed is denoted EQT.
- The EQT certificates and public keys mentioned in the figure are those for signing, i.e. VU\_Sign or Card\_Sign.
- The EQT.CA certificates and public keys mentioned in the figure are those for signing VU or Card certificates, as applicable.
- The EQT.CA.EUR certificate mentioned in the figure is the European root certificate that is indicated in the CAR of the EQT.CA certificate.
- The EQT.Link certificate mentioned in the figure is the EQT's link certificate, if present. As specified in section 9.1.2, this is a link certificate for a new European root key pair created by the ERCA and signed with the previous European private key.
- The EQT.Link.EUR certificate is the European root certificate that is indicated in the CAR of the EQT.Link certificate.
- CSM\_235 For calculating the hash M sent to the control card in the PSO:Hash command, the IDE shall use the hashing algorithm linked to the key size of the VU or the card from which the data is downloaded, as specified in CSM\_50.
- CSM\_236 For verifying the EQT's signature, the control card shall follow the signature scheme defined in [DSS].

*Note:* This document does not specify any action to undertake if a signature over a downloaded data file cannot be verified or if the verification is unsuccessful.

Figure 13

### Protocol for verification of the signature over a downloaded data file

![](_page_77_Figure_3.jpeg)

▼<u>M1</u>