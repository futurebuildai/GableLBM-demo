# Cost Projection: 50 Users @ 24h/wk

**Scenario**: 50 Concurrent Users, 4 hours/day, 6 days/week.
**Architecture**: AWS PaaS (Amplify, App Runner, Aurora v2).

## 1. Workload Estimation
*   **Total Active Hours**: ~1,200 user-hours/week.
*   **Concurrency Peak**: 50 users (assuming all active simultaneously).
*   **Request Rate**: High interaction (Lumber yard counting/sales). Est 10 req/min/user = 500 req/min peak.

## 2. Component Costs

### A. Frontend (AWS Amplify)
*   **Build Minutes**: ~$5/mo (Code commits).
*   **Data Transfer**: 50 users * 100MB/day (optimistic) = 120GB/mo.
*   **Cost**:
    *   Hosting: $0 (First 12 months) or ~$0.023/GB.
    *   **Est**: ~$3 - $5 / month.

### B. Backend (AWS App Runner)
*   **Instance**: 1 vCPU / 2GB RAM ($0.0064/vCPU-hour + $0.0064/GB-hour).
*   **Auto-Scaling**:
    *   Active (10 hours/day for coverage): 2 instances (to handle 50 concurrent w/ low latency).
    *   Idle (14 hours/day): 1 instance (provisioned warmth).
*   **Calculation**:
    *   Active: 2 inst * $0.0128/hr * 260 hrs = $6.65
    *   Idle: 1 inst * $0.0128/hr * 460 hrs = $5.88
    *   Requests: ~1M requests ($0.64/million).
*   **Est**: ~$15 - $20 / month.

### C. Database (Aurora Serverless v2)
*   **Compute (ACU)**:
    *   Min: 0.5 ACU (Idle) = 24/7 baseline.
    *   Burst: Up to 2-4 ACU during 4-hour peak.
*   **Reference**: 1 ACU = ~$0.12/hour.
*   **Calculation**:
    *   Peak (4h * 6d * 4w = 96h): 2 ACU * $0.12 * 96 = $23.04
    *   Off-Peak (624h): 0.5 ACU * $0.12 * 624 = $37.44
*   **Storage**: 10GB ($1.00).
*   **I/O**: ~$10.00 (Standard OLTP).
*   **Est**: ~$75 - $90 / month.

## 3. Total Monthly Projection

| Service | Estimated Cost | Notes |
| :--- | :--- | :--- |
| **Frontend** | $5.00 | Data transfer primarily. |
| **Backend** | $20.00 | App Runner is very efficient. |
| **Database** | $85.00 | Aurora is the premium component. |
| **Total** | **~$110.00 / month** | **$2.20 per user / month** |

## Quick Optimization
To drop this to ~$50/mo:
1.  Switch Database to **RDS T4g.micro** (Standard Provisioned) = ~$15/mo. (Loses auto-scaling, but 50 users is fine on a micro DB).
2.  App Runner is already optimized.

**Verdict**: Extremely affordable. The per-user cost is negligible compared to a SaaS subscription ($50-100/user).
