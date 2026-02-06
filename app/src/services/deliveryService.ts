import type {
    Vehicle, Driver, Route, Delivery,
    CreateVehicleRequest, CreateDriverRequest, CreateRouteRequest,
    AssignOrderRequest, UpdateDeliveryStatusRequest
} from '../types/delivery';

const API_BASE = '/api/v1/delivery';

export const deliveryService = {
    // Fleet
    listVehicles: async (): Promise<Vehicle[]> => {
        const res = await fetch(`${API_BASE}/vehicles`);
        if (!res.ok) throw new Error('Failed to fetch vehicles');
        return res.json();
    },

    createVehicle: async (req: CreateVehicleRequest): Promise<Vehicle> => {
        const res = await fetch(`${API_BASE}/vehicles`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(req),
        });
        if (!res.ok) throw new Error('Failed to create vehicle');
        return res.json();
    },

    listDrivers: async (): Promise<Driver[]> => {
        const res = await fetch(`${API_BASE}/drivers`);
        if (!res.ok) throw new Error('Failed to fetch drivers');
        return res.json();
    },

    createDriver: async (req: CreateDriverRequest): Promise<Driver> => {
        const res = await fetch(`${API_BASE}/drivers`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(req),
        });
        if (!res.ok) throw new Error('Failed to create driver');
        return res.json();
    },

    // Routes
    listRoutes: async (date?: string): Promise<Route[]> => {
        const params = new URLSearchParams();
        if (date) params.append('date', date);

        const res = await fetch(`${API_BASE}/routes?${params.toString()}`);
        if (!res.ok) throw new Error('Failed to fetch routes');
        return res.json();
    },

    createRoute: async (req: CreateRouteRequest): Promise<Route> => {
        const res = await fetch(`${API_BASE}/routes`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(req),
        });
        if (!res.ok) throw new Error('Failed to create route');
        return res.json();
    },

    dispatchRoute: async (id: string): Promise<void> => {
        const res = await fetch(`${API_BASE}/routes/${id}/dispatch`, {
            method: 'POST'
        });
        if (!res.ok) throw new Error('Failed to dispatch route');
    },

    // Deliveries
    listDeliveries: async (routeId: string): Promise<Delivery[]> => {
        const res = await fetch(`${API_BASE}/routes/${routeId}/deliveries`);
        if (!res.ok) throw new Error('Failed to fetch deliveries');
        return res.json();
    },

    assignOrder: async (req: AssignOrderRequest): Promise<Delivery> => {
        const res = await fetch(`${API_BASE}/deliveries`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(req),
        });
        if (!res.ok) throw new Error('Failed to assign order');
        return res.json();
    },

    updateStatus: async (id: string, req: UpdateDeliveryStatusRequest): Promise<void> => {
        const res = await fetch(`${API_BASE}/deliveries/${id}/status`, {
            method: 'PUT',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(req),
        });
        if (!res.ok) throw new Error('Failed to update delivery status');
    }
};
