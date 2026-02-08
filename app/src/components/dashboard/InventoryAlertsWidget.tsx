import { AlertTriangle, Package } from 'lucide-react';
import { Card, CardHeader, CardTitle, CardContent } from '../ui/Card';
import type { InventoryAlert } from '../../types/dashboard';

interface InventoryAlertsWidgetProps {
    alerts: InventoryAlert[];
    loading?: boolean;
}

export function InventoryAlertsWidget({ alerts, loading = false }: InventoryAlertsWidgetProps) {
    if (loading) {
        return (
            <Card variant="glass" className="h-full">
                <CardHeader>
                    <div className="h-6 w-32 bg-white/10 rounded animate-pulse" />
                </CardHeader>
                <CardContent>
                    <div className="space-y-4">
                        {[1, 2, 3].map((i) => (
                            <div key={i} className="flex gap-4">
                                <div className="h-10 w-10 bg-white/10 rounded animate-pulse shrink-0" />
                                <div className="space-y-2 flex-1">
                                    <div className="h-4 w-3/4 bg-white/10 rounded animate-pulse" />
                                    <div className="h-3 w-1/2 bg-white/10 rounded animate-pulse" />
                                </div>
                            </div>
                        ))}
                    </div>
                </CardContent>
            </Card>
        );
    }

    return (
        <Card variant="glass" className="h-full">
            <CardHeader>
                <CardTitle className="flex items-center gap-2">
                    <AlertTriangle className="w-5 h-5 text-amber-500" />
                    Inventory Alerts
                </CardTitle>
            </CardHeader>
            <CardContent className="p-0">
                <div className="divide-y divide-white/5">
                    {alerts.length === 0 ? (
                        <div className="p-8 text-center flex flex-col items-center gap-3">
                            <div className="h-12 w-12 rounded-full bg-emerald-500/10 flex items-center justify-center">
                                <Package className="w-6 h-6 text-emerald-500" />
                            </div>
                            <p className="text-zinc-500">All inventory levels healthy</p>
                        </div>
                    ) : (
                        alerts.map((alert) => (
                            <div key={alert.product_id} className="p-4 hover:bg-white/5 transition-colors group">
                                <h4 className="text-sm font-medium text-white mb-1 group-hover:text-amber-400 transition-colors">
                                    {alert.name}
                                </h4>
                                <div className="flex items-center justify-between text-xs">
                                    <span className="text-zinc-500">
                                        SKU: <span className="font-mono">{alert.sku}</span>
                                    </span>
                                    <div className="flex items-center gap-3">
                                        <span className="text-zinc-400">
                                            Current: <span className="text-white font-mono font-bold">{alert.current_qty}</span>
                                        </span>
                                        <span className={`font-medium ${alert.alert_type === 'OUT_OF_STOCK' ? 'text-rose-500' : 'text-amber-500'}`}>
                                            {alert.alert_type === 'OUT_OF_STOCK' ? 'Out of Stock' : 'Low Stock'}
                                        </span>
                                    </div>
                                </div>
                            </div>
                        ))
                    )}
                </div>
            </CardContent>
        </Card>
    );
}
