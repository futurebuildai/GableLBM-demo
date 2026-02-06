import type { MillworkOption, CreateOptionRequest } from '../types/millwork';

const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080';

export const MillworkService = {
    async getOptionsByCategory(category: string): Promise<MillworkOption[]> {
        const response = await fetch(`${API_URL}/api/millwork/options?category=${category}`);
        if (!response.ok) {
            throw new Error('Failed to fetch millwork options');
        }
        return response.json();
    },

    async createOption(option: CreateOptionRequest): Promise<MillworkOption> {
        const response = await fetch(`${API_URL}/api/millwork/options`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(option),
        });
        if (!response.ok) {
            throw new Error('Failed to create millwork option');
        }
        return response.json();
    },
};
