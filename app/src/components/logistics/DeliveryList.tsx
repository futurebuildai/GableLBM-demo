import React, { useEffect, useState } from 'react';
import type { Delivery } from '../../types/delivery';
import { deliveryService } from '../../services/deliveryService';
import { MapPin, Box, FileText, ArrowRight } from 'lucide-react';
import { Button } from '../ui/Button';

interface DeliveryListProps {
    routeId: string | null;
}

export const DeliveryList: React.FC<DeliveryListProps> = ({ routeId }) => {
    const [deliveries, setDeliveries] = useState<Delivery[]>([]);
    const [loading, setLoading] = useState(false);

    useEffect(() => {
        if (routeId) {
            loadDeliveries(routeId);
        } else {
            setDeliveries([]);
        }
    }, [routeId]);

    const loadDeliveries = async (id: string) => {
        setLoading(true);
        try {
            const data = await deliveryService.listDeliveries(id);
            setDeliveries(data);
        } catch (err) {
            console.error(err);
        } finally {
            setLoading(false);
        }
    };

    if (!routeId) {
        return (
            <div className="flex flex-col items-center justify-center h-full text-zinc-500 gap-4 p-12">
                <div className="w-16 h-16 rounded-full bg-white/5 flex items-center justify-center">
                    <MapPin className="w-8 h-8 opacity-50" />
                </div>
                <p>Select a route from the left to view its delivery manifest.</p>
            </div>
        );
    }

    if (loading) return (
        <div className="flex flex-col items-center justify-center h-full text-zinc-500 gap-4">
            <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-gable-green"></div>
            <p>Loading Manifest...</p>
        </div>
    );

    return (
        <div className="flex flex-col h-full">
            <div className="p-4 border-b border-white/5 bg-white/5 flex justify-between items-center">
                <h2 className="text-lg font-bold text-white flex items-center gap-2">
                    <FileText className="w-5 h-5 text-sky-400" />
                    Delivery Manifest
                </h2>
                <div className="text-xs text-zinc-400 font-mono">
                    {deliveries.length} DROPS
                </div>
            </div>

            <div className="flex-1 overflow-y-auto p-6 space-y-6 relative">
                {/* Timeline Line */}
                {deliveries.length > 0 && (
                    <div className="absolute left-[2.25rem] top-6 bottom-6 w-px bg-gradient-to-b from-gable-green/50 via-white/10 to-transparent"></div>
                )}

                {deliveries.map((delivery, index) => (
                    <div key={delivery.id} className="relative pl-12 group">
                        {/* Timeline Node */}
                        <div className="absolute left-6 top-6 -translate-x-1/2 w-6 h-6 rounded-full bg-[#0A0B10] border-2 border-gable-green flex items-center justify-center shadow-[0_0_10px_rgba(0,255,163,0.3)] z-10 text-[10px] font-bold text-white">
                            {index + 1}
                        </div>

                        <div className="bg-[#161821] border border-white/5 p-5 rounded-xl hover:border-gable-green/30 hover:bg-white/5 transition-all duration-300 group-hover:translate-x-1">
                            <div className="flex justify-between items-start mb-2">
                                <span className="font-bold text-lg text-white group-hover:text-gable-green transition-colors">{delivery.customer_name}</span>
                                <span className="text-[10px] font-mono uppercase bg-white/5 px-2 py-1 rounded text-zinc-400 border border-white/5">
                                    {delivery.status}
                                </span>
                            </div>

                            <div className="flex items-start gap-2 text-zinc-400 text-sm mb-4">
                                <MapPin className="w-4 h-4 shrink-0 mt-0.5 text-zinc-600" />
                                {delivery.address}
                            </div>

                            <div className="flex items-center gap-4 text-xs font-mono text-zinc-500 pl-6 border-l-2 border-white/5">
                                <span className="flex items-center gap-1.5">
                                    <Box className="w-3 h-3" />
                                    Order #{delivery.order_number}
                                </span>
                            </div>

                            {delivery.delivery_instructions && (
                                <div className="mt-4 text-sm bg-amber-500/5 text-amber-500/90 p-3 rounded-lg border border-amber-500/10 italic flex gap-2">
                                    <span className="not-italic font-bold text-[10px] px-1.5 py-0.5 bg-amber-500/20 rounded h-fit">NOTE</span>
                                    {delivery.delivery_instructions}
                                </div>
                            )}
                        </div>
                    </div>
                ))}

                {deliveries.length === 0 && (
                    <div className="text-zinc-500 text-center py-12">No deliveries assigned to this route.</div>
                )}
            </div>

            <div className="p-4 border-t border-white/5 bg-white/5">
                <Button variant="outline" className="w-full border-dashed border-white/20 hover:border-gable-green/50 text-gable-green hover:bg-gable-green/5">
                    <ArrowRight className="w-4 h-4 mr-2" />
                    Assign Order to Route
                </Button>
            </div>
        </div>
    );
};
