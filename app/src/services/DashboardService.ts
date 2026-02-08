import type {
    DashboardSummary,
    InventoryAlert,
    TopCustomer,
    OrderActivity,
    RevenueTrendPoint,
} from '../types/dashboard';

const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080';

export const DashboardService = {
    /**
     * Fetches aggregate KPIs for the executive dashboard.
     */
    async getSummary(): Promise<DashboardSummary> {
        const response = await fetch(`${API_URL}/api/v1/dashboard/summary`);
        if (!response.ok) {
            throw new Error('Failed to fetch dashboard summary');
        }
        return response.json();
    },

    /**
     * Fetches products with low or zero stock.
     */
    async getInventoryAlerts(): Promise<InventoryAlert[]> {
        const response = await fetch(`${API_URL}/api/v1/dashboard/inventory-alerts`);
        if (!response.ok) {
            throw new Error('Failed to fetch inventory alerts');
        }
        return response.json();
    },

    /**
     * Fetches top customers by revenue.
     */
    async getTopCustomers(): Promise<TopCustomer[]> {
        const response = await fetch(`${API_URL}/api/v1/dashboard/top-customers`);
        if (!response.ok) {
            throw new Error('Failed to fetch top customers');
        }
        return response.json();
    },

    /**
     * Fetches recent orders and status distribution.
     */
    async getOrderActivity(): Promise<OrderActivity> {
        const response = await fetch(`${API_URL}/api/v1/dashboard/order-activity`);
        if (!response.ok) {
            throw new Error('Failed to fetch order activity');
        }
        return response.json();
    },

    /**
     * Fetches 7-day revenue trend for chart.
     */
    async getRevenueTrend(): Promise<RevenueTrendPoint[]> {
        const response = await fetch(`${API_URL}/api/v1/dashboard/revenue-trend`);
        if (!response.ok) {
            throw new Error('Failed to fetch revenue trend');
        }
        return response.json();
    },
};
