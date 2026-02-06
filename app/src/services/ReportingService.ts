import type { DailyTillReport, SalesSummaryReport } from '../types/reporting';

const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080';

export const ReportingService = {
    async getDailyTill(date?: string): Promise<DailyTillReport> {
        const url = new URL(`${API_URL}/api/reports/daily-till`);
        if (date) url.searchParams.append('date', date);
        const response = await fetch(url.toString());
        if (!response.ok) throw new Error('Failed to fetch daily till');
        return response.json();
    },

    async getSalesSummary(start?: string, end?: string): Promise<SalesSummaryReport> {
        const url = new URL(`${API_URL}/api/reports/sales-summary`);
        if (start) url.searchParams.append('start', start);
        if (end) url.searchParams.append('end', end);
        const response = await fetch(url.toString());
        if (!response.ok) throw new Error('Failed to fetch sales summary');
        return response.json();
    }
};
