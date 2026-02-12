import { useState, useEffect, useCallback } from 'react';
import { Card, CardContent } from '../../components/ui/Card';
import { ShoppingCart, RefreshCw, AlertTriangle, ChevronDown, ChevronUp } from 'lucide-react';
import { PortalService } from '../../services/PortalService';
import type { PortalOrder } from '../../types/portal';
import { useToast } from '../../components/ui/ToastContext';

const formatCurrency = (val: number): string =>
    new Intl.NumberFormat('en-US', { style: 'currency', currency: 'USD' }).format(val);

export const PortalOrders = () => {
    const [orders, setOrders] = useState<PortalOrder[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState('');
    const [expandedId, setExpandedId] = useState<string | null>(null);
    const [reorderingId, setReorderingId] = useState<string | null>(null);
    const { showToast } = useToast();

    const fetchOrders = useCallback(() => {
        setLoading(true);
        setError('');
        PortalService.getOrders()
            .then(setOrders)
            .catch(err => setError(err instanceof Error ? err.message : 'Failed to load orders'))
            .finally(() => setLoading(false));
    }, []);

    useEffect(() => { fetchOrders(); }, [fetchOrders]);

    const handleReorder = async (orderId: string) => {
        setReorderingId(orderId);
        try {
            const resp = await PortalService.reorder(orderId);
            showToast(
                `Reorder created! New draft: ${resp.order_id.substring(0, 8).toUpperCase()}`,
                'success',
            );
            fetchOrders();
        } catch (err) {
            showToast(
                err instanceof Error ? err.message : 'Failed to create reorder',
                'error',
            );
        } finally {
            setReorderingId(null);
        }
    };

    if (loading) {
        return (
            <div className="space-y-4">
                {[1, 2, 3, 4].map(i => (
                    <div key={i} className="h-20 bg-white/5 rounded-2xl animate-pulse" />
                ))}
            </div>
        );
    }

    if (error) {
        return (
            <div className="flex flex-col items-center justify-center h-64 text-center">
                <AlertTriangle className="w-12 h-12 text-amber-500 mb-4" />
                <p className="text-zinc-400 mb-4">{error}</p>
                <button
                    onClick={fetchOrders}
                    className="flex items-center gap-2 px-4 py-2 rounded-lg bg-white/5 border border-white/10 text-white hover:bg-white/10 transition-colors"
                >
                    <RefreshCw size={16} /> Retry
                </button>
            </div>
        );
    }

    return (
        <div>
            <div className="mb-6 flex items-center justify-between">
                <div>
                    <h1 className="text-2xl font-bold text-white">Order History</h1>
                    <p className="text-zinc-400 text-sm mt-1">{orders.length} order{orders.length !== 1 ? 's' : ''} found</p>
                </div>
            </div>

            {orders.length === 0 ? (
                <Card variant="glass">
                    <CardContent className="p-12 text-center">
                        <ShoppingCart className="w-12 h-12 text-zinc-600 mx-auto mb-4" />
                        <p className="text-zinc-400">No orders yet.</p>
                    </CardContent>
                </Card>
            ) : (
                <div className="space-y-3">
                    {orders.map(order => (
                        <Card key={order.id} variant="glass" noPadding>
                            {/* Order Row */}
                            <div
                                className="flex items-center justify-between p-4 cursor-pointer hover:bg-white/5 transition-colors"
                                onClick={() => setExpandedId(expandedId === order.id ? null : order.id)}
                            >
                                <div className="flex items-center gap-4">
                                    <div
                                        className="w-10 h-10 rounded-lg flex items-center justify-center"
                                        style={{ backgroundColor: 'rgba(56,189,248,0.1)' }}
                                    >
                                        <ShoppingCart size={18} className="text-blue-400" />
                                    </div>
                                    <div>
                                        <div className="font-mono text-sm font-medium text-white">
                                            {order.id.substring(0, 8).toUpperCase()}
                                        </div>
                                        <div className="text-xs text-zinc-500 mt-0.5">
                                            {new Date(order.created_at).toLocaleDateString()} · {order.lines.length} item{order.lines.length !== 1 ? 's' : ''}
                                        </div>
                                    </div>
                                </div>
                                <div className="flex items-center gap-4">
                                    <div className="text-right">
                                        <div className="font-mono text-sm text-white">{formatCurrency(order.total_amount)}</div>
                                        <StatusBadge status={order.status} />
                                    </div>
                                    <button
                                        onClick={(e) => {
                                            e.stopPropagation();
                                            handleReorder(order.id);
                                        }}
                                        disabled={reorderingId === order.id}
                                        className="px-3 py-1.5 rounded-lg text-xs font-semibold bg-gable-green text-black
                                                   hover:bg-emerald-400 transition-colors disabled:opacity-50 disabled:cursor-not-allowed
                                                   flex items-center gap-1.5"
                                    >
                                        {reorderingId === order.id ? (
                                            <RefreshCw size={12} className="animate-spin" />
                                        ) : (
                                            <ShoppingCart size={12} />
                                        )}
                                        Buy Again
                                    </button>
                                    {expandedId === order.id ? (
                                        <ChevronUp size={16} className="text-zinc-500" />
                                    ) : (
                                        <ChevronDown size={16} className="text-zinc-500" />
                                    )}
                                </div>
                            </div>

                            {/* Expanded Lines */}
                            {expandedId === order.id && order.lines.length > 0 && (
                                <div className="border-t border-white/5 px-4 py-3 bg-white/[0.02]">
                                    <table className="w-full text-sm">
                                        <thead>
                                            <tr className="text-zinc-500 text-xs uppercase tracking-wider">
                                                <th className="text-left py-2 font-medium">Product</th>
                                                <th className="text-left py-2 font-medium">SKU</th>
                                                <th className="text-right py-2 font-medium">Qty</th>
                                                <th className="text-right py-2 font-medium">Price</th>
                                                <th className="text-right py-2 font-medium">Total</th>
                                            </tr>
                                        </thead>
                                        <tbody>
                                            {order.lines.map((line, idx) => (
                                                <tr key={idx} className="border-t border-white/5">
                                                    <td className="py-2 text-white">{line.product_name}</td>
                                                    <td className="py-2 font-mono text-zinc-400 text-xs">{line.product_sku}</td>
                                                    <td className="py-2 text-right font-mono text-zinc-300">{line.quantity}</td>
                                                    <td className="py-2 text-right font-mono text-zinc-300">{formatCurrency(line.price_each)}</td>
                                                    <td className="py-2 text-right font-mono text-white">{formatCurrency(line.quantity * line.price_each)}</td>
                                                </tr>
                                            ))}
                                        </tbody>
                                    </table>
                                </div>
                            )}
                        </Card>
                    ))}
                </div>
            )}
        </div>
    );
};

const StatusBadge = ({ status }: { status: string }) => {
    const colors: Record<string, string> = {
        DRAFT: 'bg-zinc-500/10 text-zinc-400 border-zinc-500/20',
        CONFIRMED: 'bg-blue-500/10 text-blue-400 border-blue-500/20',
        FULFILLED: 'bg-emerald-500/10 text-emerald-400 border-emerald-500/20',
        CANCELLED: 'bg-red-500/10 text-red-400 border-red-500/20',
        ON_HOLD: 'bg-amber-500/10 text-amber-400 border-amber-500/20',
    };
    return (
        <span className={`inline-block px-2 py-0.5 rounded text-[10px] uppercase tracking-wider font-semibold border ${colors[status] || colors.DRAFT}`}>
            {status.replace('_', ' ')}
        </span>
    );
};
