import React, { useState } from 'react';
import { RouteList } from '../components/logistics/RouteList';
import { DeliveryList } from '../components/logistics/DeliveryList';
import { RouteMap } from '../components/logistics/RouteMap';
import { PageTransition } from '../components/ui/PageTransition';
import { Truck, Calendar } from 'lucide-react';
import { Card, CardContent } from '../components/ui/Card';
import type { Delivery } from '../types/delivery';

export const DispatchBoard: React.FC = () => {
    const [selectedRouteId, setSelectedRouteId] = useState<string | null>(null);
    const [selectedVehicleId, setSelectedVehicleId] = useState<string | undefined>(undefined);
    const [currentDeliveries, setCurrentDeliveries] = useState<Delivery[]>([]);

    const handleSelectRoute = (routeId: string, vehicleId?: string) => {
        setSelectedRouteId(routeId);
        setSelectedVehicleId(vehicleId);
        // Reset deliveries when route changes (optional, but cleaner)
        if (routeId !== selectedRouteId) {
            setCurrentDeliveries([]);
        }
    };

    return (
        <PageTransition>
            <div className="h-[calc(100vh-2rem)] flex flex-col">
                <div className="flex justify-between items-center mb-6">
                    <div>
                        <h1 className="text-display-large text-white flex items-center gap-3">
                            <Truck className="w-10 h-10 text-gable-green" />
                            Logistics & Dispatch
                        </h1>
                        <p className="text-zinc-500 mt-1 text-lg">
                            Manage fleet routing and delivery schedules.
                        </p>
                    </div>
                    <div className="flex items-center gap-2 px-4 py-2 rounded-lg bg-white/5 border border-white/10 text-zinc-300 font-mono text-sm">
                        <Calendar className="w-4 h-4 text-gable-green" />
                        Today: {new Date().toLocaleDateString(undefined, { weekday: 'short', month: 'long', day: 'numeric' })}
                    </div>
                </div>

                <div className="flex gap-6 flex-1 min-h-0">
                    {/* Left Panel: Route List */}
                    <Card variant="glass" className="w-1/3 flex flex-col overflow-hidden">
                        <CardContent className="p-0 flex-1 overflow-hidden flex flex-col">
                            <RouteList onSelectRoute={handleSelectRoute} selectedRouteId={selectedRouteId} />
                        </CardContent>
                    </Card>

                    {/* Right Panel: Delivery Manifest & Map */}
                    <div className="w-2/3 flex flex-col gap-6">
                        <Card variant="glass" className="flex-1 flex flex-col overflow-hidden">
                            <CardContent className="p-0 flex-1 overflow-hidden flex flex-col">
                                <DeliveryList
                                    routeId={selectedRouteId}
                                    vehicleId={selectedVehicleId}
                                    onDeliveriesChange={setCurrentDeliveries}
                                />
                            </CardContent>
                        </Card>

                        <Card variant="glass" className="h-[350px] relative overflow-hidden">
                            <RouteMap deliveries={currentDeliveries} />
                        </Card>
                    </div>
                </div>
            </div>
        </PageTransition>
    );
};
