import { Card, CardHeader, CardTitle, CardContent } from '../ui/Card';
import type { RecentOrder } from '../../types/dashboard';

interface RecentOrdersFeedProps {
    orders: RecentOrder[];
    loading?: boolean;
}

const getStatusColor = (status: string) => {
    switch (status.toLowerCase()) {
        case 'submitted': return 'bg-amber-500/10 text-amber-500 border-amber-500/20';
        case 'confirmed': return 'bg-blue-500/10 text-blue-500 border-blue-500/20';
        case 'processing': return 'bg-indigo-500/10 text-indigo-500 border-indigo-500/20';
        case 'ready': return 'bg-emerald-500/10 text-emerald-500 border-emerald-500/20';
        case 'completed': return 'bg-gable-green/10 text-gable-green border-gable-green/20';
        case 'cancelled': return 'bg-rose-500/10 text-rose-500 border-rose-500/20';
        default: return 'bg-zinc-500/10 text-zinc-500 border-zinc-500/20';
    }
};

export function RecentOrdersFeed({ orders, loading = false }: RecentOrdersFeedProps) {
    if (loading) {
        return (
            <Card variant="glass" className="h-full">
                <CardHeader>
                    <div className="h-6 w-32 bg-white/10 rounded animate-pulse" />
                </CardHeader>
                <CardContent>
                    <div className="space-y-4">
                        {[1, 2, 3, 4, 5].map((i) => (
                            <div key={i} className="flex justify-between items-center">
                                <div className="space-y-2">
                                    <div className="h-4 w-24 bg-white/10 rounded animate-pulse" />
                                    <div className="h-3 w-16 bg-white/10 rounded animate-pulse" />
                                </div>
                                <div className="h-6 w-20 bg-white/10 rounded animate-pulse" />
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
                <CardTitle>Recent Activity</CardTitle>
            </CardHeader>
            <CardContent className="p-0">
                <div className="divide-y divide-white/5">
                    {orders.length === 0 ? (
                        <div className="p-6 text-center text-zinc-500">No recent orders</div>
                    ) : (
                        orders.map((order) => (
                            <div key={order.order_id} className="p-4 hover:bg-white/5 transition-colors group">
                                <div className="flex justify-between items-start mb-1">
                                    <div className="font-medium text-white group-hover:text-gable-green transition-colors">
                                        {order.customer_name}
                                    </div>
                                    <span className={`text-[10px] px-2 py-0.5 rounded-full border uppercase tracking-wider font-semibold ${getStatusColor(order.status)}`}>
                                        {order.status}
                                    </span>
                                </div>
                                <div className="flex justify-between items-center text-xs text-zinc-500">
                                    <div className="font-mono">{order.order_id.substring(0, 8)}...</div>
                                    <div className="font-mono text-zinc-400">
                                        ${(order.total_amount / 100).toLocaleString(undefined, { minimumFractionDigits: 2 })}
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
