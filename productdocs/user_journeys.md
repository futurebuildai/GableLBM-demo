# User Journeys: Critical LBM Flows

## 1. The "Contractor Rush" Counter Sale
**Actor:** Counter Carl (Sales Rep) & Bob (Contractor)
**Context:** Bob is in a hurry, needs 50 studs and 2 boxes of screws for a job starting in 1 hour. It's 6:30 AM.

### Flow
1.  **Selection:** Bob dumps 2 boxes of screws on the counter. "I need these and 50 2x4 8-footers."
2.  **Order Entry:**
    *   Carl hits `Ctrl+N` (New Order).
    *   Scans screw boxes. System beeps.
    *   Types `248` (Speed code for 2x4-8'). System auto-completes "Hem-Fir Stud 2x4 92-5/8".
    *   Inputs `50`.
3.  **Cross-Sell (AI):**
    *   System toast notification: "Bob usually buys 'Simpson A35 Clips' with studs. Ask?"
    *   Carl: "Need clips?" Bob: "Nah, got 'em."
4.  **Checkout:**
    *   Carl asks: "Job?"
    *   Bob: "The Smith house."
    *   Carl Types `Sm...`, selects "Smith Residence (Job #1042)".
    *   System applies "Job Pricing" (Level 3 - 5%).
5.  **Payment:** "Put it on my account."
    *   Carl spins signature pad. Bob signs.
6.  **Fulfillment:**
    *   System prints "Pick Ticket" at the Yard Printer (Zone B) immediately.
    *   Carl hands Bob recipe/invoice. "Go see Yvonne in the yard."

## 2. The "Yard Layout" Pick & Load
**Actor:** Yardmaster Yvonne
**Context:** Ticket #9022 (Bob's studs) just printed in the shed.

### Flow
1.  **Notification:** Yvonne's rugged tablet buzzes. "New Urgent Pick: Bob - 50 Studs - Zone B".
2.  **Locating:**
    *   Yvonne opens the ticket on tablet.
    *   Map highlights "Row 4, Bin 2".
3.  **Picking:**
    *   Yvonne drives forklift to Row 4.
    *   Visual inspection: "Bin 2 looks messy."
    *   She grabs a blend of 50 studs.
4.  **Verification:**
    *   She taps "Confirm Pick" on tablet.
    *   Input: "Did you pull from a new bunk?" -> "No."
5.  **Handoff:**
    *   Bob pulls his truck around.
    *   Yvonne loads the bunk.
    *   She snaps a photo of the load in the truck bed with the tablet (Proof of condition).
    *   Bob thumbs-up. Yvonne taps "Dispatched".
6.  **Inventory Update:** System decrements "On Hand" in Zone B immediately.

## 3. The "Custom Millwork" Special Order
**Actor:** Contractor Ken & Counter Carl
**Context:** Ken needs a custom arched window that takes 4 weeks to make.

### Flow
1.  **Configuration:**
    *   Carl opens "Millwork Configurator" module.
    *   Selects Vendor: "Andersen".
    *   Selects Type: "400 Series Arch".
    *   Inputs dimensions: `48" x 60"`.
2.  **Quoting:**
    *   System calls Vendor API (EDI/Rest) for real-time cost & lead time.
    *   Returns: "$850 Cost, Est Ship: Feb 24th".
    *   Carl applies margin (30%). Quote: $1,215.
3.  **Deposit:**
    *   System flags: "Special Order > $1,000 requires 50% deposit."
    *   Ken pays $600 card.
4.  **Procurement:**
    *   System *automatically* generates PO #5501 to Andersen.
    *   Status set to "Ordered - Awaiting Vendor Conf".
5.  **Tracking:**
    *   Feb 24th: Truck arrives. Yvonne scans barcode.
    *   System matches PO #5501 -> Links to Customer Order #9099.
    *   **Trigger:** SMS sent to Ken: "Your window has arrived."
