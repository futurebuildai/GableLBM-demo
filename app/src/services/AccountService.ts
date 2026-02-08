import type { AccountSummary, CustomerTransaction } from '../types/account';

const API_BASE_URL = 'http://localhost:8080/api';

export const AccountService = {
    getAccountSummary: async (customerId: string): Promise<AccountSummary> => {
        const response = await fetch(`${API_BASE_URL}/accounts/${customerId}`);
        if (!response.ok) {
            throw new Error('Failed to fetch account summary');
        }
        return response.json();
    },

    getTransactions: async (customerId: string): Promise<CustomerTransaction[]> => {
        const response = await fetch(`${API_BASE_URL}/accounts/${customerId}/transactions`);
        if (!response.ok) {
            throw new Error('Failed to fetch transactions');
        }
        return response.json();
    },
};
