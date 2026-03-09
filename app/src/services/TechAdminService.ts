const API_URL = 'https://backend-production-bdf8.up.railway.app';

export interface APIKey {
    id: string;
    name: string;
    prefix: string;
    scopes: string[];
    created_at: string;
    last_used_at?: string;
    revoked_at?: string;
}

export interface CreateKeyResponse {
    api_key: string;
    key: APIKey;
}

export const techAdminService = {
    async listKeys(): Promise<APIKey[]> {
        const response = await fetch(`${API_URL}/api/admin/keys`);
        if (!response.ok) {
            throw new Error('Failed to fetch API keys');
        }
        const data = await response.json();
        return data || [];
    },

    async createKey(name: string, scopes: string[]): Promise<CreateKeyResponse> {
        const response = await fetch(`${API_URL}/api/admin/keys`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ name, scopes }),
        });
        if (!response.ok) {
            throw new Error('Failed to create API key');
        }
        return response.json();
    },

    async revokeKey(id: string): Promise<void> {
        const response = await fetch(`${API_URL}/api/admin/keys/${id}`, {
            method: 'DELETE',
        });
        if (!response.ok) {
            throw new Error('Failed to revoke API key');
        }
    }
};
