import type { CalculatedPrice } from '../types/pricing';

const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080';

export const PricingService = {
    calculatePrice: async (customerId: string, productId: string, quantity?: number, jobId?: string): Promise<CalculatedPrice> => {
        const params = new URLSearchParams({
            customer_id: customerId,
            product_id: productId,
        });
        if (quantity && quantity > 0) {
            params.set('quantity', quantity.toString());
        }
        if (jobId) {
            params.set('job_id', jobId);
        }
        const response = await fetch(`${API_URL}/pricing/calculate?${params.toString()}`);
        if (!response.ok) {
            throw new Error('Failed to calculate price');
        }
        return response.json();
    }
};
