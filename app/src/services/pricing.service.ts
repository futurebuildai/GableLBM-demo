import type { CalculatedPrice } from '../types/pricing';

const API_URL = 'http://localhost:8080';

export const PricingService = {
    calculatePrice: async (customerId: string, productId: string): Promise<CalculatedPrice> => {
        const response = await fetch(`${API_URL}/pricing/calculate?customer_id=${customerId}&product_id=${productId}`);
        if (!response.ok) {
            throw new Error('Failed to calculate price');
        }
        return response.json();
    }
};
