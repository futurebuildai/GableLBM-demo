import { useEffect, useState } from "react";
import { useParams, useNavigate } from "react-router-dom";
import { deliveryService } from "../../services/deliveryService";
import type { Delivery } from "../../types/delivery";

export function StopList() {
    const { id } = useParams<{ id: string }>();
    const navigate = useNavigate();
    const [deliveries, setDeliveries] = useState<Delivery[]>([]);

    useEffect(() => {
        if (id) {
            // Fetch Route Details (we might need a getRoute method or reuse list?)
            // deliveryService.getRoute(id) doesn't exist yet in service, wait, listRoutes returns simplified. 
            // We probably need getRoute for details or just pass state. 
            // For MVP, just fetch deliveries. 

            deliveryService.listDeliveries(id).then(setDeliveries);

            // Hack: fetch route info via listRoutes filter logic if simplified
            // Or adding getRoute to service is better.
        }
    }, [id]);

    if (!id) return <div>Invalid Route</div>;

    return (
        <div className="space-y-4 pt-4 px-4">
            <div className="flex items-center justify-between mb-4">
                <button onClick={() => navigate('/driver')} className="text-gray-400 font-mono text-sm">&larr; BACK</button>
                <div className="font-bold">STOP LIST</div>
                <div className="w-8"></div>
            </div>

            <div className="space-y-3">
                {deliveries.length === 0 && <div className="text-gray-500 text-center py-8">No stops on this route.</div>}

                {deliveries.map((d) => (
                    <div
                        key={d.id}
                        onClick={() => navigate(`/driver/deliveries/${d.id}`)}
                        className={`p-4 rounded-lg border flex items-center gap-4 transition-all active:scale-98 relative overflow-hidden ${d.status === 'DELIVERED' ? 'bg-[#161821]/50 border-green-500/20 opacity-75' :
                            d.status === 'FAILED' ? 'bg-red-500/10 border-red-500/20' :
                                d.stop_sequence === 1 && d.status === 'PENDING' ? 'bg-[#161821] border-[#00FFA3] shadow-[0_0_15px_rgba(0,255,163,0.1)]' :
                                    'bg-[#161821] border-white/10'
                            }`}
                    >
                        {/* Sequence Badge */}
                        <div className={`flex-shrink-0 w-8 h-8 rounded-full flex items-center justify-center font-bold font-mono ${d.status === 'DELIVERED' ? 'bg-green-500/20 text-green-500' :
                            'bg-white/10 text-white'
                            }`}>
                            {d.stop_sequence}
                        </div>

                        <div className="flex-1 min-w-0">
                            <div className="font-bold truncate text-white">{d.customer_name}</div>
                            <div className="text-sm text-gray-400 truncate">{d.address}</div>
                            <div className="text-xs font-mono text-gray-500 mt-1">{d.order_number}</div>
                        </div>

                        {d.status === 'DELIVERED' && <div className="text-green-500 text-xl">✓</div>}
                    </div>
                ))}
            </div>
        </div>
    );
}
