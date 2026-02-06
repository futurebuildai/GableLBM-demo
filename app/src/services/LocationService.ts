import type { Location } from '../types/location';

const API_URL = 'http://localhost:8080'; // TODO: use environment variable

export const LocationService = {
    async listLocations(): Promise<Location[]> {
        const response = await fetch(`${API_URL}/locations`);
        if (!response.ok) {
            throw new Error('Failed to fetch locations');
        }
        return response.json();
    },

    async createLocation(data: Omit<Location, 'id' | 'created_at' | 'updated_at'>): Promise<Location> {
        const response = await fetch(`${API_URL}/locations`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(data),
        });
        if (!response.ok) {
            throw new Error('Failed to create location');
        }
        return response.json();
    }
};
