import React, { useState } from 'react';
import { RouteList } from '../components/logistics/RouteList';
import { DeliveryList } from '../components/logistics/DeliveryList';

export const DispatchBoard: React.FC = () => {
    const [selectedRouteId, setSelectedRouteId] = useState<string | null>(null);

    return (
        <div className="p-6 h-screen flex flex-col bg-[#0A0B10] text-[#E0E0E0]">
            <div className="flex justify-between items-center mb-6 border-b border-[#FFFFFF10] pb-4">
                <h1 className="text-2xl font-bold text-white tracking-tight">Logistics & Dispatch</h1>
                <div className="text-[#38BDF8] font-mono text-sm">Today: {new Date().toLocaleDateString()}</div>
            </div>

            <div className="flex gap-6 flex-1 min-h-0">
                <div className="w-1/3 overflow-y-auto pr-2">
                    <RouteList onSelectRoute={setSelectedRouteId} />
                </div>
                <div className="w-2/3 overflow-y-auto pl-2 border-l border-[#FFFFFF10]">
                    <DeliveryList routeId={selectedRouteId} />
                </div>
            </div>
        </div>
    );
};
