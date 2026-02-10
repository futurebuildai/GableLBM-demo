import type { Quote, CreateQuoteRequest } from '../types/quote';

const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080';

export const QuoteService = {
    async createQuote(request: CreateQuoteRequest): Promise<Quote> {
        const response = await fetch(`${API_URL}/quotes`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(request),
        });

        if (!response.ok) {
            const errorText = await response.text();
            throw new Error(errorText || 'Failed to create quote');
        }

        return response.json();
    },

    async getQuote(id: string): Promise<Quote> {
        const response = await fetch(`${API_URL}/quotes/${id}`);
        if (!response.ok) {
            throw new Error('Failed to fetch quote');
        }
        return response.json();
    },

    async listQuotes(): Promise<Quote[]> {
        const response = await fetch(`${API_URL}/quotes`);
        if (!response.ok) {
            throw new Error('Failed to fetch quotes');
        }
        return response.json();
    },

    async convertToOrder(quoteId: string): Promise<{ customer_id: string; quote_id: string; lines: { product_id: string; quantity: number; price_each: number }[] }> {
        const response = await fetch(`${API_URL}/quotes/${quoteId}/convert`, {
            method: 'POST',
        });
        if (!response.ok) {
            const text = await response.text();
            throw new Error(text || 'Failed to convert quote to order');
        }
        return response.json();
    }
};
