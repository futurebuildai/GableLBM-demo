import type { DailyTillReport, SalesSummaryReport } from '../types/reporting';
import type { ARAgingReport, CustomerStatement, CreditMemo } from '../types/invoice';

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
    },

    async getARAgingReport(): Promise<ARAgingReport> {
        const response = await fetch(`${API_URL}/api/reports/ar-aging`);
        if (!response.ok) throw new Error('Failed to fetch AR aging report');
        return response.json();
    },

    async getCustomerStatement(customerId: string, start?: string, end?: string): Promise<CustomerStatement> {
        const url = new URL(`${API_URL}/api/reports/customer-statement/${customerId}`);
        if (start) url.searchParams.append('start', start);
        if (end) url.searchParams.append('end', end);
        const response = await fetch(url.toString());
        if (!response.ok) throw new Error('Failed to fetch customer statement');
        return response.json();
    },

    async createCreditMemo(invoiceId: string, amount: number, reason: string): Promise<CreditMemo> {
        const response = await fetch(`${API_URL}/invoices/${invoiceId}/credit-memo`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ amount, reason }),
        });
        if (!response.ok) throw new Error('Failed to create credit memo');
        return response.json();
    }
};
