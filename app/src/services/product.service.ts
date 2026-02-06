import type { Product } from '../types/product';

const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8085';

export const ProductService = {
    async getProducts(): Promise<Product[]> {
        const response = await fetch(`${API_URL}/products`);
        if (!response.ok) {
            throw new Error('Failed to fetch products');
        }
        return response.json();
    },

    async createProduct(product: Omit<Product, 'id' | 'created_at' | 'updated_at'>): Promise<Product> {
        const response = await fetch(`${API_URL}/products`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(product),
        });

        if (!response.ok) {
            const errorText = await response.text();
            throw new Error(errorText || 'Failed to create product');
        }

        return response.json();
    },
};
