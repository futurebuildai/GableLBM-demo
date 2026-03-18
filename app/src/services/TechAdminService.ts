const API_URL = import.meta.env.VITE_API_URL || '';

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

export interface AISettings {
    configured: boolean;
    source: 'admin' | 'env' | 'none';
    key_hint?: string;
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
    },

    // --- AI Settings ---

    async getAISettings(): Promise<AISettings> {
        const response = await fetch(`${API_URL}/api/admin/settings/ai`);
        if (!response.ok) throw new Error('Failed to fetch AI settings');
        return response.json();
    },

    async saveAIKey(apiKey: string): Promise<void> {
        const response = await fetch(`${API_URL}/api/admin/settings/ai`, {
            method: 'PUT',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ api_key: apiKey }),
        });
        if (!response.ok) {
            const text = await response.text();
            throw new Error(text || 'Failed to save API key');
        }
    },

    async deleteAIKey(): Promise<void> {
        const response = await fetch(`${API_URL}/api/admin/settings/ai`, {
            method: 'DELETE',
        });
        if (!response.ok) throw new Error('Failed to delete API key');
    },

    // --- Gemini Settings ---

    async getGeminiSettings(): Promise<AISettings> {
        const response = await fetch(`${API_URL}/api/admin/settings/gemini`);
        if (!response.ok) throw new Error('Failed to fetch Gemini settings');
        return response.json();
    },

    async saveGeminiKey(apiKey: string): Promise<void> {
        const response = await fetch(`${API_URL}/api/admin/settings/gemini`, {
            method: 'PUT',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ api_key: apiKey }),
        });
        if (!response.ok) {
            const text = await response.text();
            throw new Error(text || 'Failed to save Gemini API key');
        }
    },

    async deleteGeminiKey(): Promise<void> {
        const response = await fetch(`${API_URL}/api/admin/settings/gemini`, {
            method: 'DELETE',
        });
        if (!response.ok) throw new Error('Failed to delete Gemini API key');
    },
};
