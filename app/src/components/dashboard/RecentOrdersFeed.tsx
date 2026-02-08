import { Package, Clock } from 'lucide-react';
import type { RecentOrder } from '../../types/dashboard';

interface RecentOrdersFeedProps {
    orders: RecentOrder[];
    loading?: boolean;
}

const STATUS_STYLES: Record<string, string> = {
    PENDING: 'bg-amber-500/20 text-amber-400',
    CONFIRMED: 'bg-blue-500/20 text-blue-400',
    PROCESSING: 'bg-indigo-500/20 text-indigo-400',
    READY: 'bg-emerald-500/20 text-emerald-400',
    COMPLETED: 'bg-emerald-500/20 text-emerald-400',
    CANCELLED: 'bg-rose-500/20 text-rose-400',
};

function formatTimeAgo(dateStr: string): string {
    const date = new Date(dateStr);
    const now = new Date();
    const diffMs = now.getTime() - date.getTime();
    const diffMins = Math.floor(diffMs / 60000);
    const diffHours = Math.floor(diffMins / 60);
    const diffDays = Math.floor(diffHours / 24);

    if (diffMins < 1) return 'Just now';
    if (diffMins < 60) return `${diffMins}m ago`;
    if (diffHours < 24) return `${diffHours}h ago`;
    return `${diffDays}d ago`;
}

export function RecentOrdersFeed({ orders, loading }: RecentOrdersFeedProps) {
    if (loading) {
        return (
            <div className="p-4 bg-zinc-900 rounded-lg border border-white/10">
                <h3 className="text-sm font-medium text-zinc-400 mb-4">Recent Orders</h3>
                <div className="space-y-3">
                    {[1, 2, 3, 4, 5].map((i) => (
                        <div key={i} className="h-14 bg-zinc-800 rounded animate-pulse" />
                    ))}
                </div>
            </div>
        );
    }

    if (orders.length === 0) {
        return (
            <div className="p-4 bg-zinc-900 rounded-lg border border-white/10">
                <h3 className="text-sm font-medium text-zinc-400 mb-4">Recent Orders</h3>
                <div className="text-center py-8 text-zinc-500">No recent orders</div>
            </div>
        );
    }

    return (
        <div className="p-4 bg-zinc-900 rounded-lg border border-white/10">
            <h3 className="text-sm font-medium text-zinc-400 mb-4">Recent Orders</h3>
            <div className="space-y-2">
                {orders.map((order) => (
                    <div
                        key={order.order_id}
                        className="flex items-center gap-3 p-3 rounded-lg bg-zinc-800/50 hover:bg-zinc-800 transition-colors"
                    >
                        <div className="p-2 rounded-lg bg-zinc-700">
                            <Package className="w-4 h-4 text-zinc-400" />
                        </div>
                        <div className="flex-1 min-w-0">
                            <div className="flex items-center gap-2">
                                <span className="text-sm text-white truncate">{order.customer_name}</span>
                                <span
                                    className={`px-2 py-0.5 text-xs rounded-full ${STATUS_STYLES[order.status] || 'bg-zinc-600 text-zinc-300'
                                        }`}
                                >
                                    {order.status}
                                </span>
                            </div>
                            <div className="flex items-center gap-1 text-xs text-zinc-500">
                                <Clock className="w-3 h-3" />
                                {formatTimeAgo(order.created_at)}
                            </div>
                        </div>
                        <div className="text-right shrink-0">
                            <div className="text-sm font-mono font-bold text-emerald-400">
                                ${(order.total_amount / 100).toLocaleString()}
                            </div>
                        </div>
                    </div>
                ))}
            </div>
        </div>
    );
}
