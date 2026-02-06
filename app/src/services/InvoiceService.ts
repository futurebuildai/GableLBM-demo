import type { Invoice } from '../types/invoice';

const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8085';

export const InvoiceService = {
    async listInvoices(): Promise<Invoice[]> {
        const response = await fetch(`${API_URL}/invoices`);
        if (!response.ok) {
            throw new Error('Failed to fetch invoices');
        }
        return response.json();
    },

    async getInvoice(id: string): Promise<Invoice> {
        const response = await fetch(`${API_URL}/invoices/${id}`);
        if (!response.ok) {
            throw new Error('Failed to fetch invoice');
        }
        return response.json();
    },

    async emailInvoice(id: string): Promise<void> {
        const response = await fetch(`${API_URL}/api/invoices/${id}/email`, { // Note: /api/ prefix
            method: 'POST'
        });
        if (!response.ok) {
            throw new Error('Failed to email invoice');
        }
    }
};
