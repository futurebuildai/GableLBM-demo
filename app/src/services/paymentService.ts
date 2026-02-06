import type { Payment, CreatePaymentRequest } from '../types/payment';

const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080';

export const paymentService = {
    createPayment: async (req: CreatePaymentRequest): Promise<Payment> => {
        const response = await fetch(`${API_URL}/api/payments`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(req),
        });
        if (!response.ok) throw new Error('Failed to create payment');
        return response.json();
    },

    getHistory: async (invoiceId: string): Promise<Payment[]> => {
        const response = await fetch(`${API_URL}/api/invoices/${invoiceId}/payments`);
        if (!response.ok) throw new Error('Failed to fetch payments');
        return response.json();
    }
};
