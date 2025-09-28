# Proto Schema Audit Log - 2025-09-28

This log contains findings from the audit of protobuf schemas against the tachograph regulations.

## File: /home/odsod/github.com/way-platform/tachograph-go/proto/wayplatform/connect/tachograph/vu/v1/overview.proto

### Finding 1: Incomplete `ControlActivity` message
- **Message:** `ControlActivity`
- **Field:** `control_card_number`
- **Issue:** The field is of type `dd.v1.FullCardNumber`. According to Annex IC, section 3.12.12, the stored control activity data must include `card generation`. The `FullCardNumber` data type (Data Dictionary 2.73) does not include the card generation. The `VuControlActivityRecord` definition in the Data Dictionary (2.187) also appears to be missing this field.
- **Recommendation:** To be fully compliant, the `FullCardNumberAndGeneration` data type (Data Dictionary 2.74) should be used instead of `FullCardNumber`, or a separate `generation` field should be added to the `ControlActivity` message.

### Finding 2: Incomplete `CompanyLock` message
- **Message:** `CompanyLock`
- **Field:** `company_card_number`
- **Issue:** The field is of type `dd.v1.FullCardNumber`. According to Annex IC, section 3.12.13, stored company lock data must include `card generation`. The `FullCardNumber` data type (Data Dictionary 2.73) does not include card generation, and the corresponding `VuCompanyLocksRecord` in the Data Dictionary (2.184) is also missing this.
- **Recommendation:** To ensure compliance, consider changing the type of `company_card_number` to `dd.v1.FullCardNumberAndGeneration` (Data Dictionary 2.74) or adding a separate `generation` field.

### Finding 3: Incomplete `DownloadActivity` message
- **Message:** `DownloadActivity`
- **Field:** `full_card_number`
- **Issue:** The field is of type `dd.v1.FullCardNumber`. According to Annex IC, section 3.12.14, the stored download activity data must include `card generation`. The `FullCardNumber` data type (Data Dictionary 2.73) does not include this, and the `VuDownloadActivityData` definition in the Data Dictionary (2.195) also omits it.
- **Recommendation:** To align with the regulation, the type of `full_card_number` should be changed to `dd.v1.FullCardNumberAndGeneration` (Data Dictionary 2.74), or a `generation` field should be added.

### Finding 4: Refactoring Suggestion
- **Messages:** `DownloadActivity`, `CompanyLock`, `ControlActivity`
- **Suggestion:** These messages represent specific, self-contained data records (`VuDownloadActivityData`, `VuCompanyLocksRecord`, `VuControlActivityRecord`) from the vehicle unit's memory. They have stand-alone semantics and are good candidates for being moved to the `wayplatform.connect.tachograph.dd.v1` package. This would improve modularity and allow for their potential reuse in other contexts.

---