#### *Addendum*

Rules for the computation of daily, weekly and fortnightly driving time

1. Basic computation rules

The VU shall compute the daily driving time, the weekly driving time and the fortnightly driving time using relevant data stored in a driver (or workshop) card inserted in the driver slot (slot 1, card reader #1) of the Vehicle Unit, and selected driver's activities while this card is inserted in the VU.

The driving times shall not be calculated while no driver (or workshop) card is inserted.

UNKNOWN period(s) found during the time period needed for computations shall be assimilated to BREAK/REST.

UNKNOWN periods and activities of negative duration (i.e. start of the activity occurs later than the end of the activity) due to time overlaps between two different VUs or due to time adjustment, are not taken into account.

Activities recorded in the driver card corresponding to 'OUT OF SCOPE' periods in accordance with definition (gg) of Annex IC, shall be interpreted as follows:

- BREAK/REST shall be computed as 'BREAK' or 'REST'
- WORK and DRIVING shall be considered as 'WORK'
- AVAILABILITY shall be considered as 'AVAILABILITY'

In the context of this Addendum, the VU shall assume to have a daily rest period at the beginning of the card activities records.

2. Concepts

The following concepts apply exclusively to this appendix, and are intended to specify the computation of driving times by the VU and its later transmission by the remote communication facility.

(a) 'RTM-shift' is the period between the end of a daily rest period and the end of the directly following daily rest period.

The VU shall start a new RTM-shift after a daily rest period has finished.

The ongoing RTM-shift is the period since the end of last daily rest period;

- (b) 'accumulated driving time' is the sum of the duration of all DRIVING activities of the driver within a period while not in OUT OF SCOPE;
- (c) 'daily driving time' is the accumulated driving time within a RTM-shift;
- (d) 'weekly driving time' is the accumulated driving time for the ongoing week;
- (e) 'continuous rest period' is any uninterrupted period of BREAK/REST;
- (f) 'fortnightly driving time' is the accumulated driving time for the previous and the ongoing week;

(g) 'daily rest period' is a period of BREAK/REST, which can be either

- a regular daily rest period,
- a split daily rest period or
- a reduced daily rest period

In the context of Appendix 14, when a VU is computing weekly rest periods, those weekly rest periods shall be considered as daily rest periods;

(h) 'regular daily rest period' is a continuous rest period of at least 11 hours.

As a matter of exception, when a FERRY/TRAIN CROSSING condition is active the regular daily rest period may be interrupted a maximum of two times by activities other than rest, with a maximal accumulated duration of one hour, i.e. the regular daily rest period containing ferry/train crossing period(s) may be split into two or three parts. The VU shall then compute a regular daily rest period when the accumulated rest time computed according to point 3 is at least 11 hours.

When a regular daily rest period has been interrupted the VU:

- shall not incorporate the driving activity encountered during those interruptions to the computation of the daily driving time, and
- shall start a new RTM-shift at the end of the regular daily rest period that has been interrupted.

##### *Figure 1.*

### **Example of daily rest period interrupted due to ferry/train crossing**

| A<br>0/火/0     | B<br>H | C<br>0/火/0 | D<br>H A | E<br>0/火/0 | F<br>H | G<br>0/火/0 |
|----------------|--------|------------|----------|------------|--------|------------|
| Working Period | 2 h    | 30 min     | 8 h      | 30 min     | 2 h    | New Day    |

- (i) 'reduced daily rest period' is a continuous rest period of at least 9 hours and less than 11 hours;
- (j) 'split daily rest period' is a daily rest period taken in two parts:
  - the first part shall be a continuous rest period of at least 3 hours and less than 9,
  - the second part shall be a continuous rest period of at least 9 hours.

As a matter of exception, when a FERRY/TRAIN CROSSING condition is active during one or both of the parts of a split daily rest period, the split daily rest period may be interrupted a maximum of two times by other activities with the accumulated duration of maximal one hour, i.e.:

- the first part of the split daily rest period may be interrupted one or two times, or
- the second part of the split daily rest period may be interrupted one or two times, or

— the first part of the split daily rest period may be interrupted one time and the second part of the split daily rest period may be interrupted one time.

The VU shall then compute a split daily rest period when the accumulated rest time computed according to point 3 is:

- at least three hours and less than 11 hours for the first rest period and at least 9 hours for the second rest period, when the first rest period has been interrupted by FERRY/TRAIN CROSSING.
- at least three hours and less than 9 hours for the first rest period and at least 9 hours for the second rest period, when the first rest period has not been interrupted by FERRY/TRAIN CROSSING.

##### *Figure 2.*

#### **Example of split daily rest period interrupted due to ferry/train crossing**

| A     | B            | C      | D            | E                  | F            | G      | H            | I       |
|-------|--------------|--------|--------------|--------------------|--------------|--------|--------------|---------|
| 0/火/日 | Image: chair | 0/火/日  | Image: chair | 0/火/日/Image: chair | Image: chair | 0/火/日  | Image: chair | 0/火/日   |
| 4 h   | 1h           | 20 min | 2 h          | 6 h                | 7h           | 20 min | 3h           | New Day |

When the split daily rest period is interrupted, the VU:

- shall not incorporate the driving activity encountered during those interruptions to the computation of the daily driving time, and
- shall start a new RTM-shift at the end of the split daily rest period that has been interrupted;
- (k) 'week' is the period in UTC time between 00:00 hours on Monday and 24:00 hours on Sunday;
- 3. Computation of the rest period when it has been interrupted due to ferry/train crossing

For the computation of the rest period when it has been interrupted due to ferry/train crossing, the VU shall calculate the accumulated rest time according to the following steps:

a) Step 1

The VU shall detect interruptions to the rest time occurring before the activation of the FERRY/TRAIN CROSSING (BEGIN) flag, according to figure 3 and in its case figure 4, and shall evaluate for each interruption detected if the following conditions are met:

- the interruption makes the total duration of the interruptions detected, including in its case interruptions occurring during the first part of a split daily rest period due to ferry/train crossing, to exceed more than one hour in total,
- the interruption makes the total number of interruptions detected, including in its case interruptions occurring during the first part of a split daily rest period due to ferry/train crossing, to be bigger than two,
- there is an 'Entry of place where daily work periods end' stored after the interruption ended.

If none of the above conditions are met, the continuous rest period immediately preceding the interruption shall be added to the accumulated rest time.

If at least one of the above conditions is met, the VU shall either stop the computation of the accumulated rest time according to step 2 or detect interruptions to the rest time occurring after the FERRY/TRAIN CROSSING (BEGIN) flag according to step 3.

b) Step 2

For each interruption detected according to step 1, the VU shall evaluate whether the computation of the accumulated rest time should stop. The VU shall stop the computation process when two continuous rest periods occurring before the activation of the FERRY/TRAIN CROSSING (BEGIN) flag have been added to the accumulated rest time, including in its case rest periods added in the first part of a split daily rest period also interrupted by ferry/train crossing. Otherwise, the VU shall proceed according to step 3.

c) Step 3

If after performance of step 2 the VU continues the computation of the accumulated rest time, the VU shall detect interruptions occurring after the deactivation of the FERRY/TRAIN CROSSING condition according to figure 3 and in its case figure 4.

For each interruption found, the VU shall evaluate if the interruption makes the accumulated time of all the interruptions detected to exceed more than one hour in total, in which case the computation of the accumulated rest period shall finish at the end of the continuous rest period previous to the interruption. Otherwise, the continuous rest periods occurring after the respective interruptions shall be added to the computation of the daily rest period until the condition in step 4 is fulfilled.

d) Step 4

The computation of the accumulated rest time shall stop when the VU has added, as result of steps 1 and 3, a maximum of two continuous rest periods to the rest period for which the FERRY/TRAIN CROSSING condition is activated, including in its case interruptions occurring during the first part of a split daily rest period due to ferry/train crossing.

##### *Figure 3.*

**Processing of rest times by the VU in order to determine whether an interrupted rest period shall compute as regular daily rest period or as the first part of a split daily rest period**

![](_page_3_Figure_11.jpeg)

*Figure 4.*

**Processing of rest times by the VU in order to determine whether an interrupted rest period shall compute as the second part of a split daily rest period**

![](_page_4_Figure_3.jpeg)

*Figure 5.*

**Example of a daily rest period interrupted more than twice causing rest period H not to be included in the computation**

![](_page_4_Figure_6.jpeg)

![](_page_4_Figure_7.jpeg)

**Example of a daily rest period where Ferry/Train Calculation period is commenced at end of work period**

![](_page_4_Figure_9.jpeg)

![](_page_4_Figure_10.jpeg)

**Example of a daily rest period interrupted more than twice causing rest period B not to be included in the computation**

| A                                                | B        | C                                         | D        | E                                         | F        | G                                         | H             | I                                                |
|--------------------------------------------------|----------|-------------------------------------------|----------|-------------------------------------------|----------|-------------------------------------------|---------------|--------------------------------------------------|
| $\oslash/\text{\textasciicircum}/\oslash/\vdash$ | $\vdash$ | $\oslash/\text{\textasciicircum}/\oslash$ | $\vdash$ | $\oslash/\text{\textasciicircum}/\oslash$ | $\vdash$ | $\oslash/\text{\textasciicircum}/\oslash$ | $\vdash$      | $\oslash/\text{\textasciicircum}/\oslash/\vdash$ |
| 4,5h                                             | 1h       | 10 min                                    | 1h       | 10 min                                    | 1h       | 10 min                                    | 9h            |                                                  |
| Working                                          | Rest     | Movement                                  | Rest     | Movement                                  | Rest     | Embarking                                 | Rest on ferry |                                                  |
|                                                  |          |                                           |          |                                           |          |                                           |               | Start of New Shift                               |

*Figure 8.*

### **Example of a split daily rest period interrupted once during the first rest period and once during 2nd rest period**

| A         | B    | C         | D             | E         | F    | G         | H             | I                  |
|-----------|------|-----------|---------------|-----------|------|-----------|---------------|--------------------|
| $Θ/∗/□/h$ | h    | $Θ/∗/□$   | h             | $Θ/∗/□/h$ | h    | $Θ/∗/□$   | h             | $Θ/∗/□/h$          |
| 3h        | 1h   | 10 min    | 2h            | 6h        | 2h   | 10 min    | 7h            |                    |
| Working   | Rest | Embarking | Rest on ferry | Working   | Rest | Embarking | Rest on ferry | Start of New Shift |

## 4. Computation of daily, weekly and fortnightly- driving times

The VU shall compute the daily driving time(s) for the ongoing and previous RTM-shifts. The driving time occurring during the interruptions of the daily rest periods shall not be added to the computation of the daily driving time, when such interruptions are due to ferry/train crossing and the requirements provided for in paragraphs (h) and (j) of point 2 and in point 3 have been fulfilled. Nevertheless, insofar as a complete regular or split daily rest period has not been computed by the VU according to point 3, the driving times occurring during the interruptions shall be added to the daily driving time for the ongoing RTM-shift.

The VU shall also compute the weekly and the fortnightly driving times. The driving time occurring during the interruptions of the daily rest periods due to ferry/train crossing shall be added to the computation of the weekly and the fortnightly driving times.