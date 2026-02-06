export interface StockAdjustmentRequest {
    product_id: string;
    location_id?: string;
    quantity: number;
    reason: string;
    is_delta: boolean;
}

const API_URL = 'http://localhost:8080';

export const InventoryService = {
    async adjustStock(data: StockAdjustmentRequest): Promise<void> {
        const response = await fetch(`${API_URL}/inventory/adjust`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(data),
        });
        if (!response.ok) {
            throw new Error('Failed to adjust stock');
        }
    },

    // Helper to get inventory for a product if we want to show it in the modal
    async getInventoryByProduct(_productId: string): Promise<any[]> {
        // TODO: Implement list endpoint in backend if not exists, 
        // or rely on the main product list to contain total if needed.
        // For Sprint 03, we rely on the main list or assume 0 start.
        return [];
    }
};
