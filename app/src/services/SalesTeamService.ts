import type { SalesPerson } from '../types/salesteam';

const API_URL = '';

export const SalesTeamService = {
    async listSalesTeam(): Promise<SalesPerson[]> {
        const response = await fetch(`${API_URL}/sales-team`);
        if (!response.ok) {
            throw new Error('Failed to fetch sales team');
        }
        return response.json();
    },

    async getSalesPerson(id: string): Promise<SalesPerson> {
        const response = await fetch(`${API_URL}/sales-team/${id}`);
        if (!response.ok) {
            throw new Error('Failed to fetch salesperson');
        }
        return response.json();
    },
};
