import { AlertTriangle, AlertCircle } from 'lucide-react';
import type { InventoryAlert } from '../../types/dashboard';

interface InventoryAlertsWidgetProps {
    alerts: InventoryAlert[];
    loading?: boolean;
}

export function InventoryAlertsWidget({ alerts, loading }: InventoryAlertsWidgetProps) {
    if (loading) {
        return (
            <div className="p-4 bg-zinc-900 rounded-lg border border-white/10">
                <h3 className="text-sm font-medium text-zinc-400 mb-4">Inventory Alerts</h3>
                <div className="space-y-3">
                    {[1, 2, 3].map((i) => (
                        <div key={i} className="h-12 bg-zinc-800 rounded animate-pulse" />
                    ))}
                </div>
            </div>
        );
    }

    if (alerts.length === 0) {
        return (
            <div className="p-4 bg-zinc-900 rounded-lg border border-emerald-900/30">
                <h3 className="text-sm font-medium text-zinc-400 mb-4">Inventory Alerts</h3>
                <div className="flex items-center gap-2 text-emerald-400 py-4">
                    <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                    </svg>
                    <span>All stock levels healthy</span>
                </div>
            </div>
        );
    }

    return (
        <div className="p-4 bg-zinc-900 rounded-lg border border-white/10">
            <h3 className="text-sm font-medium text-zinc-400 mb-4">
                Inventory Alerts <span className="text-rose-400 ml-1">({alerts.length})</span>
            </h3>
            <div className="space-y-2 max-h-64 overflow-y-auto">
                {alerts.map((alert) => (
                    <div
                        key={alert.product_id}
                        className={`p-3 rounded-lg border ${alert.alert_type === 'OUT_OF_STOCK'
                                ? 'border-rose-900/50 bg-rose-950/20'
                                : 'border-amber-900/50 bg-amber-950/20'
                            }`}
                    >
                        <div className="flex items-start gap-2">
                            {alert.alert_type === 'OUT_OF_STOCK' ? (
                                <AlertCircle className="w-4 h-4 text-rose-400 mt-0.5 shrink-0" />
                            ) : (
                                <AlertTriangle className="w-4 h-4 text-amber-400 mt-0.5 shrink-0" />
                            )}
                            <div className="flex-1 min-w-0">
                                <div className="text-sm text-white truncate">{alert.name}</div>
                                <div className="text-xs text-zinc-500 font-mono">{alert.sku}</div>
                            </div>
                            <div className="text-right shrink-0">
                                <div
                                    className={`text-sm font-mono font-bold ${alert.alert_type === 'OUT_OF_STOCK' ? 'text-rose-400' : 'text-amber-400'
                                        }`}
                                >
                                    {alert.current_qty}
                                </div>
                                <div className="text-xs text-zinc-500">/ {alert.reorder_qty}</div>
                            </div>
                        </div>
                    </div>
                ))}
            </div>
        </div>
    );
}
