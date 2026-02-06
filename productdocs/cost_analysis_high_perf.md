# Cost Projection: "Maximum Performance" Tier

**Budget**: $500/mo ($10/user/mo).
**Goal**: Zero latency, no cold starts, instant reporting.

## 1. Architecture Upgrades

### A. Database: Aurora Serverless v2 (Turbo Mode)
*   **Strategy**: Instead of scaling down to 0.5 ACU (2GB RAM) when idle, we keep the database "hot" at **4 ACU** (approx 8GB RAM + 2 vCPU dedicated) 24/7.
*   **Benefit**: Complex reports (Daily Till, aggregated sales) will run instantly. No "warm up" lag in the morning.
*   **Cost**:
    *   4 ACU * $0.12/hour * 730 hours = **$350.40 / month**.

### B. Backend: App Runner (HPC Config)
*   **Strategy**: Vertical Scaling + redundancy.
*   **Config**: 2 vCPU / 4 GB RAM per container.
*   **Provisioning**: Keep **2 instances** active 24/7 (High Availability).
*   **Benefit**: CPU-intensive tasks (PDF generation, encryption) are 2x faster. Zero cold starts.
*   **Cost**:
    *   Instance: ~$0.026/hour * 2 instances * 730 hours = **$37.96 / month**.
    *   Requests: ~$10.00 (buffer).
    *   Total: **~$48.00 / month**.

### C. Frontend: AWS Amplify
*   **Strategy**: Standard Edge Cache.
*   **Cost**: **~$5.00 / month**.

## 2. Total Monthly Projection

| Service | Cost | Performance Impact |
| :--- | :--- | :--- |
| **Database (4 ACU)** | $350.50 | **Lightning fast reports**, zero-lag concurrency. |
| **Backend (2x vCPU)** | $48.00 | **Sub-second PDF generation**, instant API response. |
| **Frontend** | $5.00 | Global CDN distribution. |
| **Total** | **~$403.50 / month** | **$8.07 per user / month** |

## verdict
This architecture fits perfectly within your **$10/user** limit. By allocating ~85% of the budget to the Database (the heart of an ERP), you ensure that the system feels "instant" even as you load it with years of data.
