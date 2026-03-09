import type {
    Vehicle, Driver, Route, Delivery, CapacityWarning,
    CreateVehicleRequest, CreateDriverRequest, CreateRouteRequest,
    AssignOrderRequest, UpdateDeliveryStatusRequest
} from '../types/delivery';

const API_BASE = 'https://backend-production-bdf8.up.railway.app';

export const deliveryService = {
    // Fleet
    listVehicles: async (): Promise<Vehicle[]> => {
        const res = await fetch(`${API_BASE}/api/v1/delivery/vehicles`);
        if (!res.ok) throw new Error('Failed to fetch vehicles');
        return res.json();
    },

    createVehicle: async (req: CreateVehicleRequest): Promise<Vehicle> => {
        const res = await fetch(`${API_BASE}/api/v1/delivery/vehicles`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(req),
        });
        if (!res.ok) throw new Error('Failed to create vehicle');
        return res.json();
    },

    listDrivers: async (): Promise<Driver[]> => {
        const res = await fetch(`${API_BASE}/api/v1/delivery/drivers`);
        if (!res.ok) throw new Error('Failed to fetch drivers');
        return res.json();
    },

    createDriver: async (req: CreateDriverRequest): Promise<Driver> => {
        const res = await fetch(`${API_BASE}/api/v1/delivery/drivers`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(req),
        });
        if (!res.ok) throw new Error('Failed to create driver');
        return res.json();
    },

    // Routes
    listRoutes: async (date?: string, driverId?: string): Promise<Route[]> => {
        const params = new URLSearchParams();
        if (date) params.append('date', date);
        if (driverId) params.append('driver_id', driverId);

        const res = await fetch(`${API_BASE}/api/v1/delivery/routes?${params.toString()}`);
        if (!res.ok) throw new Error('Failed to fetch routes');
        return res.json();
    },

    createRoute: async (req: CreateRouteRequest): Promise<Route> => {
        const res = await fetch(`${API_BASE}/api/v1/delivery/routes`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(req),
        });
        if (!res.ok) throw new Error('Failed to create route');
        return res.json();
    },

    reorderStops: async (routeId: string, orderedDeliveryIds: string[]): Promise<void> => {
        const res = await fetch(`${API_BASE}/api/v1/delivery/routes/${routeId}/reorder`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ ordered_delivery_ids: orderedDeliveryIds })
        });
        if (!res.ok) throw new Error('Failed to reorder stops');
    },

    dispatchRoute: async (id: string): Promise<void> => {
        const res = await fetch(`${API_BASE}/api/v1/delivery/routes/${id}/dispatch`, {
            method: 'POST'
        });
        if (!res.ok) throw new Error('Failed to dispatch route');
    },

    // Deliveries
    listDeliveries: async (routeId: string): Promise<Delivery[]> => {
        const res = await fetch(`${API_BASE}/api/v1/delivery/routes/${routeId}/deliveries`);
        if (!res.ok) throw new Error('Failed to fetch deliveries');
        return res.json();
    },

    getDelivery: async (id: string): Promise<Delivery> => {
        const res = await fetch(`${API_BASE}/api/v1/delivery/deliveries/${id}`);
        if (!res.ok) throw new Error('Failed to fetch delivery');
        return res.json();
    },

    assignOrder: async (req: AssignOrderRequest): Promise<{ delivery: Delivery; capacity_warning?: CapacityWarning }> => {
        const res = await fetch(`${API_BASE}/api/v1/delivery/deliveries`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(req),
        });
        if (!res.ok) throw new Error('Failed to assign order');
        return res.json();
    },

    updateStatus: async (id: string, req: UpdateDeliveryStatusRequest): Promise<void> => {
        const res = await fetch(`${API_BASE}/api/v1/delivery/deliveries/${id}/status`, {
            method: 'PUT',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(req),
        });
        if (!res.ok) throw new Error('Failed to update delivery status');
    }
};
