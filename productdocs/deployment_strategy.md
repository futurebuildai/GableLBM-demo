# Deployment Architecture: AWS PaaS Strategy

**Objective**: Deploy GableLBM using high-abstraction AWS Managed Services for minimal operational overhead.
**Tooling**: Terraform (Infrastructure as Code).

## 1. Architecture Diagram
```mermaid
graph TD
    User([User]) --> Amplify[AWS Amplify (Frontend / CDN)]
    User --> AppRunner[AWS App Runner (Backend API)]
    
    subgraph AWS_VPC [VPC]
        AppRunner --> |Private Link| Aurora[Aurora Serverless v2 (Postgres)]
        AppRunner --> |Egress| InternetGateway
    end
    
    GitHub((GitHub Repo)) --> |Push Master| Amplify
    GitHub --> |Action: Build Image| ECR[Elastic Container Registry]
    ECR --> |New Image Trigger| AppRunner
```

## 2. Terraform Resource Plan

This plan is divided into 4 core modules to keep state manageable.

### Module A: Networking (The Plumbing)
*   **Resource**: `aws_vpc`
    *   CIDR: `10.0.0.0/16`
*   **Resource**: `aws_subnet` (Public & Private)
    *   2 Public Subnets (for Load Balancers/NAT if needed)
    *   2 Private Subnets (for Database & App Runner VPC Connector)
*   **Resource**: `aws_security_group`
    *   `db_sg`: Allow Inbound 5432 only from `app_runnner_sg`.

### Module B: Database (The State)
*   **Resource**: `aws_rds_cluster` (Aurora Postgres)
    *   Engine: `aurora-postgresql`
    *   Mode: `provisioned` (Serverless v2 uses provisioned instances heavily modified)
    *   `serverlessv2_scaling_configuration`:
        *   Min Capacity: 0.5 ACU
        *   Max Capacity: 2.0 ACU
*   **Resource**: `aws_rds_cluster_instance`
    *   Instance Class: `db.serverless`
*   **Output**: `db_endpoint`, `db_secret_arn` (Secrets Manager)

### Module C: Backend Compute (The API)
*   **Resource**: `aws_ecr_repository`
    *   Name: `gable-backend`
*   **Resource**: `aws_apprunner_service`
    *   Source: Image Repository (ECR)
    *   Instance Config: 1 vCPU, 2 GB Mem
    *   **Network Configuration**:
        *   `vpc_connector_arn`: Points to Private Subnets.
    *   **Environment Variables**:
        *   `DATABASE_URL`: Injected from Secrets Manager.
        *   `JWKS_URL`: Clerk/Auth provider URL.

### Module D: Frontend (The UI)
*   **Resource**: `aws_amplify_app`
    *   Repository: `github.com/futurebuild/GableLBM`
    *   Build Spec: `generic_build.yml` (npm run build)
*   **Resource**: `aws_amplify_branch`
    *   Branch: `master`
    *   Stage: `PRODUCTION`
    *   Enable Auto Build: `true`

## 3. Deployment Strategy (GitHub Actions)

We will need a `.github/workflows/deploy-prod.yml` to glue the backend build to the infrastructure.

1.  **Terraform Apply**: Runs on PR merge to provision changes.
2.  **Backend Build**:
    *   `docker build -t gable-backend .`
    *   `aws ecr login`
    *   `docker push $ECR_REGISTRY/gable-backend:latest`
    *   (App Runner auto-deploys upon seeing the new tag).

## 4. Estimated Costs (Low Traffic / Dev)
*   **Aurora v2**: ~$40/mo (0.5 ACU baseline).
*   **App Runner**: ~$15/mo (per active instance, scales to 0).
*   **Amplify**: Free tier eligible / Pay per GB transfer (negligible).
*   **Total**: ~$60-$80/month for a production-grade, auto-scaling environment.

## 5. Next Steps
1.  Setup `aws-cli` and `terraform` locally.
2.  Create S3 Backend for Terraform State (`gable-tf-state`).
3.  Write `main.tf` implementing Module A (Networking) first.
