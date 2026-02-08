import React, { useState } from 'react';
import { RouteList } from '../components/logistics/RouteList';
import { DeliveryList } from '../components/logistics/DeliveryList';
import { PageTransition } from '../components/ui/PageTransition';
import { Truck, Map, Calendar } from 'lucide-react';
import { Card, CardContent } from '../components/ui/Card';

export const DispatchBoard: React.FC = () => {
    const [selectedRouteId, setSelectedRouteId] = useState<string | null>(null);

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
                            <RouteList onSelectRoute={setSelectedRouteId} selectedRouteId={selectedRouteId} />
                        </CardContent>
                    </Card>

                    {/* Right Panel: Delivery Manifest & Map Placeholder */}
                    <div className="w-2/3 flex flex-col gap-6">
                        <Card variant="glass" className="flex-1 flex flex-col overflow-hidden">
                            <CardContent className="p-0 flex-1 overflow-hidden flex flex-col">
                                <DeliveryList routeId={selectedRouteId} />
                            </CardContent>
                        </Card>

                        {/* Map Placeholder - Future Integration */}
                        <Card variant="glass" className="h-[300px] relative overflow-hidden group">
                            <div className="absolute inset-0 bg-[#161821] flex items-center justify-center">
                                <div className="absolute inset-0 opacity-20 bg-[radial-gradient(#2d3748_1px,transparent_1px)] [background-size:16px_16px]"></div>
                                <div className="text-center z-10">
                                    <div className="w-16 h-16 rounded-full bg-gable-green/10 flex items-center justify-center mx-auto mb-4 group-hover:scale-110 transition-transform duration-500">
                                        <Map className="w-8 h-8 text-gable-green" />
                                    </div>
                                    <h3 className="text-white font-medium">Live Fleet Map</h3>
                                    <p className="text-zinc-500 text-sm mt-1">Real-time telematics integration coming soon</p>
                                </div>
                            </div>
                        </Card>
                    </div>
                </div>
            </div>
        </PageTransition>
    );
};
