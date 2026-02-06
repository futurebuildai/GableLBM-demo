import { useEffect, useState } from "react";
import { deliveryService } from "../../services/deliveryService";
import type { Driver, Route } from "../../types/delivery";
import { useNavigate } from "react-router-dom";

export function RouteList() {
    const [drivers, setDrivers] = useState<Driver[]>([]);
    const [selectedDriver, setSelectedDriver] = useState<string>("");
    const [routes, setRoutes] = useState<Route[]>([]);
    const navigate = useNavigate();

    useEffect(() => {
        deliveryService.listDrivers().then(setDrivers);
    }, []);

    useEffect(() => {
        if (selectedDriver) {
            deliveryService.listRoutes(undefined, selectedDriver).then(setRoutes);
        } else {
            setRoutes([]);
        }
    }, [selectedDriver]);

    return (
        <div className="space-y-4 pt-4">
            <div className="bg-[#161821] p-4 rounded-lg border border-white/10 mx-4">
                <label className="block text-sm text-gray-400 mb-2 font-mono uppercase tracking-wider">Driver Login</label>
                <select
                    value={selectedDriver}
                    onChange={e => setSelectedDriver(e.target.value)}
                    className="w-full bg-[#0A0B10] border border-white/20 p-3 rounded text-white focus:outline-none focus:border-[#00FFA3]"
                >
                    <option value="">Select your name...</option>
                    {drivers && drivers.map(d => (
                        <option key={d.id} value={d.id}>{d.name}</option>
                    ))}
                </select>
            </div>

            <div className="space-y-3 px-4">
                {routes && routes.map(route => (
                    <div
                        key={route.id}
                        onClick={() => navigate(`/driver/routes/${route.id}`)}
                        className="bg-[#161821] p-4 rounded-lg border border-white/10 active:scale-98 transition-transform cursor-pointer hover:border-[#00FFA3]/50"
                    >
                        <div className="flex justify-between items-center mb-2">
                            <span className="font-mono text-sm text-gray-400">
                                {new Date(route.scheduled_date).toLocaleDateString()}
                            </span>
                            <span className={`text-xs px-2 py-1 rounded font-bold ${route.status === 'IN_TRANSIT' ? 'bg-blue-500/20 text-blue-400 border border-blue-500/30' :
                                route.status === 'COMPLETED' ? 'bg-green-500/20 text-green-400 border border-green-500/30' :
                                    'bg-gray-700/50 text-gray-300 border border-white/10'
                                }`}>
                                {route.status}
                            </span>
                        </div>
                        <div className="text-xl font-bold mb-1">{route.vehicle_name}</div>
                        <div className="flex items-center text-sm text-gray-400 gap-2">
                            <span className="bg-white/5 px-2 py-0.5 rounded">{route.stop_count} Stops</span>
                            {route.notes && <span className="italic truncate max-w-[150px]">{route.notes}</span>}
                        </div>
                    </div>
                ))}
                {selectedDriver && routes.length === 0 && (
                    <div className="text-center text-gray-500 py-12 flex flex-col items-center">
                        <div className="text-4xl mb-2">🚚</div>
                        <div>No routes assigned today.</div>
                    </div>
                )}
            </div>
        </div>
    );
}
