import React, { useEffect, useState } from 'react';
import type { Delivery } from '../../types/delivery';
import { deliveryService } from '../../services/deliveryService';

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
        return <div className="text-[#FFFFFF40] p-12 text-center border border-dashed border-[#FFFFFF20] rounded-lg bg-[#FFFFFF05]">Select a route to view manifest</div>;
    }

    if (loading) return <div className="text-[#FFFFFF40] animate-pulse">Loading Manifest...</div>;

    return (
        <div className="space-y-4">
            <h2 className="text-lg font-bold mb-4 text-[#38BDF8]">Manifest</h2>
            <div className="space-y-4">
                {deliveries.map((delivery, index) => (
                    <div key={delivery.id} className="bg-[#161821] border border-[#FFFFFF10] p-4 rounded flex items-start gap-4 hover:bg-[#1C1F2B] transition-colors">
                        <div className="bg-[#FFFFFF05] rounded text-[#38BDF8] w-8 h-8 flex items-center justify-center font-mono font-bold shrink-0 border border-[#FFFFFF05]">
                            {index + 1}
                        </div>
                        <div className="flex-1">
                            <div className="flex justify-between mb-1">
                                <span className="font-bold text-white text-lg">{delivery.customer_name}</span>
                                <span className="text-xs font-mono bg-[#FFFFFF10] px-2 py-1 rounded text-[#9CA3AF] border border-[#FFFFFF05]">{delivery.status}</span>
                            </div>
                            <div className="text-[#9CA3AF] mb-2">{delivery.address}</div>
                            <div className="flex items-center gap-4 text-xs font-mono text-[#6B7280]">
                                <span>Order #{delivery.order_number}</span>
                            </div>
                            {delivery.delivery_instructions && (
                                <div className="mt-3 text-sm bg-[#EAB30810] text-[#EAB308] p-3 rounded border border-[#EAB30820]">
                                    <span className="font-bold mr-2">NOTE:</span>
                                    {delivery.delivery_instructions}
                                </div>
                            )}
                        </div>
                    </div>
                ))}
                {deliveries.length === 0 && (
                    <div className="text-[#FFFFFF40] text-center py-8">No deliveries assigned to this route.</div>
                )}
            </div>
            <button className="mt-4 w-full border border-dashed border-[#FFFFFF20] py-3 rounded text-[#00FFA3] hover:bg-[#00FFA310] hover:border-[#00FFA350] transition-colors uppercase tracking-wide text-sm font-bold">
                + Assign Order to Route
            </button>
        </div>
    );
};
