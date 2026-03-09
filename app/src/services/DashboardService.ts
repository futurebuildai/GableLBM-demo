import type {
    DashboardSummary,
    InventoryAlert,
    TopCustomer,
    OrderActivity,
    RevenueTrendPoint,
} from '../types/dashboard';

const API_URL = 'https://backend-production-bdf8.up.railway.app';
const REQUEST_TIMEOUT_MS = 10_000; // 10 seconds
const MAX_RETRIES = 1;
const RETRY_DELAY_MS = 2_000;

/**
 * Fetch with timeout and retry logic.
 * Uses AbortController for request cancellation and retries on network errors.
 */
async function fetchWithRetry<T>(url: string, retries = MAX_RETRIES): Promise<T> {
    let lastError: Error | null = null;

    for (let attempt = 0; attempt <= retries; attempt++) {
        const controller = new AbortController();
        const timeoutId = setTimeout(() => controller.abort(), REQUEST_TIMEOUT_MS);

        try {
            const response = await fetch(url, { signal: controller.signal });
            clearTimeout(timeoutId);

            if (!response.ok) {
                throw new Error(`API error: ${response.status} ${response.statusText}`);
            }

            return await response.json() as T;
        } catch (error) {
            clearTimeout(timeoutId);
            lastError = error instanceof Error ? error : new Error(String(error));

            // Don't retry on non-network errors (4xx responses parsed above throw before reaching catch)
            if (lastError.name === 'AbortError') {
                lastError = new Error('Request timed out');
            }

            // If we have retries left, wait then retry
            if (attempt < retries) {
                await new Promise(resolve => setTimeout(resolve, RETRY_DELAY_MS));
            }
        }
    }

    throw lastError ?? new Error('Request failed');
}

export const DashboardService = {
    /**
     * Fetches aggregate KPIs for the executive dashboard.
     */
    async getSummary(): Promise<DashboardSummary> {
        return fetchWithRetry<DashboardSummary>(`${API_URL}/api/v1/dashboard/summary`);
    },

    /**
     * Fetches products with low or zero stock.
     */
    async getInventoryAlerts(): Promise<InventoryAlert[]> {
        return fetchWithRetry<InventoryAlert[]>(`${API_URL}/api/v1/dashboard/inventory-alerts`);
    },

    /**
     * Fetches top customers by revenue.
     */
    async getTopCustomers(): Promise<TopCustomer[]> {
        return fetchWithRetry<TopCustomer[]>(`${API_URL}/api/v1/dashboard/top-customers`);
    },

    /**
     * Fetches recent orders and status distribution.
     */
    async getOrderActivity(): Promise<OrderActivity> {
        return fetchWithRetry<OrderActivity>(`${API_URL}/api/v1/dashboard/order-activity`);
    },

    /**
     * Fetches 7-day revenue trend for chart.
     */
    async getRevenueTrend(): Promise<RevenueTrendPoint[]> {
        return fetchWithRetry<RevenueTrendPoint[]>(`${API_URL}/api/v1/dashboard/revenue-trend`);
    },
};
