# Feature Specification: Tech Admin Panel

**Goal**: Empower non-technical administrators to manage the technical aspects of their ERP, specifically API access and 3rd-party integrations, without needing an engineer.

## 1. Overview
The Tech Admin Panel is a specialized section of the Admin Dashboard (`/admin/tech`). It abstracts complex backend configurations (API Gateway, Webhooks, Cors) into a user-friendly UI.

## 2. Core Features

### A. API Key Management
*   **Generate Keys**: Create new API keys with specific scopes (e.g., `read:inventory`, `write:orders`).
*   **Revoke/Rotate**: One-click revocation of compromised keys.
*   **Usage Monitoring**: Simple charts showing API usage/limits per key.

### B. Integration Marketplace ("One-Click Config")
*   **Pre-built Integrations**:
    *   **Run Payments**: Toggle to enable, enter credentials.
    *   **QuickBooks Online**: OAuth2 flow to connect GL.
    *   **Avalara**: Tax calculation connection.
*   **Custom Webhooks**:
    *   UI to define "Event Triggers" (e.g., `Order Created`, `Inventory Low`).
    *   Enter a Destination URL.
    *   Test button to send a sample payload.

### C. Database & Data Export
*   **"Your Data" Export**: One-click full SQL dump download (for self-hosting backups).
*   **Schema Viewer**: Visual representation of the data model for analysts.

### D. System Health (For Self-Hosted)
*   **Logs**: View recent system logs/errors.
*   **Environment Variables**: safely view non-secret env vars.
*   **Update Checker**: Check GitHub for newer releases of GableLBM.

## 3. UI/UX Design
*   **Theme**: "Dark Console" aesthetic (consistent with the rest of the app, but perhaps with more "technical" monospace fonts).
*   **Safety**: "Danger Zone" red accents for destructive actions (Revoking keys, etc.).

## 4. Technical Implementation
*   **Backend**:
    *   New API endpoints: `POST /api/admin/keys`, `GET /api/admin/integrations`.
    *   Middleware to enforce scopes on generated keys.
*   **Frontend**:
    *   New Route: `/admin/tech`.
    *   Components: `KeyGenerator`, `WebhookConfig`, `IntegrationCard`.

## 5. Success Metrics
*   **Time to Hello World**: A non-engineer should be able to generate a key and make a successful cURL request within 5 minutes.
*   **Integration Speed**: reduced time to connect a 3rd party tool (e.g., Zapier) from days to minutes.
