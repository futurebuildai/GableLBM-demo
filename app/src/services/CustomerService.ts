import type { Customer, PriceLevel } from '../types/customer';

const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080';

export const CustomerService = {
    async listCustomers(): Promise<Customer[]> {
        const response = await fetch(`${API_URL}/customers`);
        if (!response.ok) {
            throw new Error('Failed to fetch customers');
        }
        return response.json();
    },

    async createCustomer(customer: Omit<Customer, 'id' | 'created_at' | 'updated_at' | 'balance_due'>): Promise<Customer> {
        const response = await fetch(`${API_URL}/customers`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(customer),
        });

        if (!response.ok) {
            const errorText = await response.text();
            throw new Error(errorText || 'Failed to create customer');
        }

        return response.json();
    },

    async listPriceLevels(): Promise<PriceLevel[]> {
        const response = await fetch(`${API_URL}/price_levels`);
        if (!response.ok) {
            throw new Error('Failed to fetch price levels');
        }
        return response.json();
    }
};
