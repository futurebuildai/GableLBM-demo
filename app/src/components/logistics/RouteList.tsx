import React, { useEffect, useState } from 'react';
import type { Route, RouteStatus } from '../../types/delivery';
import { deliveryService } from '../../services/deliveryService';

// Simple helper for dates
const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString();
};

const StatusBadge: React.FC<{ status: RouteStatus }> = ({ status }) => {
    let color = 'bg-[#FFFFFF10] text-gray-400';
    switch (status) {
        case 'SCHEDULED': color = 'bg-[#38BDF820] text-[#38BDF8] border border-[#38BDF850]'; break;
        case 'IN_TRANSIT': color = 'bg-[#EAB30820] text-[#EAB308] border border-[#EAB30850]'; break;
        case 'COMPLETED': color = 'bg-[#00FFA320] text-[#00FFA3] border border-[#00FFA350]'; break;
        case 'CANCELLED': color = 'bg-[#F43F5E20] text-[#F43F5E] border border-[#F43F5E50]'; break;
    }
    return <span className={`px-2 py-0.5 rounded text-[10px] font-mono uppercase tracking-wider ${color}`}>{status}</span>;
};

interface RouteListProps {
    onSelectRoute: (routeId: string) => void;
}

export const RouteList: React.FC<RouteListProps> = ({ onSelectRoute }) => {
    const [routes, setRoutes] = useState<Route[]>([]);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        loadRoutes();
    }, []);

    const loadRoutes = async () => {
        try {
            const data = await deliveryService.listRoutes();
            setRoutes(data);
        } catch (err) {
            console.error(err);
        } finally {
            setLoading(false);
        }
    };

    if (loading) return <div className="text-[#FFFFFF40] animate-pulse">Loading Routes...</div>;

    return (
        <div className="space-y-4">
            <h2 className="text-lg font-bold mb-4 text-[#00FFA3]">Routes</h2>
            <div className="space-y-3">
                {routes.map(route => (
                    <div
                        key={route.id}
                        onClick={() => onSelectRoute(route.id)}
                        className="bg-[#161821] border border-[#FFFFFF10] p-4 rounded hover:border-[#00FFA3] hover:shadow-[0_0_15px_rgba(0,255,163,0.1)] cursor-pointer transition duration-200 group"
                    >
                        <div className="flex justify-between items-start mb-2">
                            <div className="font-bold text-white group-hover:text-[#00FFA3] transition-colors">{route.vehicle_name}</div>
                            <StatusBadge status={route.status} />
                        </div>
                        <div className="text-sm text-[#9CA3AF] mb-1">{route.driver_name}</div>
                        <div className="flex justify-between items-end mt-2">
                            <div className="text-xs text-[#6B7280] font-mono">{formatDate(route.scheduled_date)}</div>
                            <div className="text-xs font-mono bg-[#FFFFFF05] px-2 py-1 rounded text-[#E0E0E0]">{route.stop_count} Stops</div>
                        </div>
                    </div>
                ))}
                {routes.length === 0 && (
                    <div className="text-[#FFFFFF40] text-center py-4 border border-dashed border-[#FFFFFF20] rounded">No routes found.</div>
                )}
            </div>
            <button className="mt-4 w-full bg-[#00FFA3] text-black font-bold py-3 rounded hover:bg-[#00E090] hover:shadow-[0_0_20px_rgba(0,255,163,0.3)] transition-all uppercase tracking-wide text-sm">
                + Create New Route
            </button>
        </div>
    );
};
