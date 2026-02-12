import { useState, useEffect, useCallback } from 'react';
import { Card, CardHeader, CardTitle, CardContent } from '../../components/ui/Card';
import { DollarSign, CreditCard, AlertTriangle, ShoppingCart, ArrowRight, RefreshCw } from 'lucide-react';
import { Link } from 'react-router-dom';
import { PortalService } from '../../services/PortalService';
import type { PortalDashboard as DashboardData } from '../../types/portal';

const formatCurrency = (val: number): string =>
    new Intl.NumberFormat('en-US', { style: 'currency', currency: 'USD' }).format(val);

export const PortalDashboard = () => {
    const [data, setData] = useState<DashboardData | null>(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState('');

    const fetchData = useCallback(() => {
        setLoading(true);
        setError('');
        PortalService.getDashboard()
            .then(setData)
            .catch(err => setError(err instanceof Error ? err.message : 'Failed to load dashboard'))
            .finally(() => setLoading(false));
    }, []);

    // eslint-disable-next-line react-hooks/set-state-in-effect -- async fetch→setState is standard pattern
    useEffect(() => { fetchData(); }, [fetchData]);

    const userName = (() => {
        try {
            const stored = localStorage.getItem('portal_user');
            if (stored) {
                const u = JSON.parse(stored) as { name: string };
                return u.name?.split(' ')[0] || 'Contractor';
            }
        } catch { /* ignore */ }
        return 'Contractor';
    })();

    if (loading) {
        return (
            <div className="space-y-6">
                <div className="h-10 w-72 bg-white/5 rounded-lg animate-pulse" />
                <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
                    {[1, 2, 3].map(i => (
                        <div key={i} className="h-32 bg-white/5 rounded-2xl animate-pulse" />
                    ))}
                </div>
                <div className="h-64 bg-white/5 rounded-2xl animate-pulse" />
            </div>
        );
    }

    if (error) {
        return (
            <div className="flex flex-col items-center justify-center h-64 text-center">
                <AlertTriangle className="w-12 h-12 text-amber-500 mb-4" />
                <p className="text-zinc-400 mb-4">{error}</p>
                <button
                    onClick={fetchData}
                    className="flex items-center gap-2 px-4 py-2 rounded-lg bg-white/5 border border-white/10 text-white hover:bg-white/10 transition-colors"
                >
                    <RefreshCw size={16} /> Retry
                </button>
            </div>
        );
    }

    const stats = [
        {
            title: 'Current Balance',
            value: formatCurrency(data?.balance_due || 0),
            icon: DollarSign,
            color: '#00FFA3',
        },
        {
            title: 'Credit Limit',
            value: formatCurrency(data?.credit_limit || 0),
            icon: CreditCard,
            color: '#38BDF8',
        },
        {
            title: 'Past Due',
            value: formatCurrency(data?.past_due || 0),
            icon: AlertTriangle,
            color: data?.past_due && data.past_due > 0 ? '#F43F5E' : '#00FFA3',
        },
    ];

    return (
        <div>
            <div className="mb-8">
                <h1 className="text-display-large text-white">Welcome back, {userName}</h1>
                <p className="text-zinc-400 mt-2 text-lg">Here's your account overview.</p>
            </div>

            {/* Stats Row */}
            <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
                {stats.map((stat) => (
                    <Card key={stat.title} variant="glass" className="group hover:-translate-y-1 transition-transform duration-300">
                        <CardContent className="p-6">
                            <div className="flex justify-between items-start mb-4">
                                <div
                                    className="p-3 rounded-lg border"
                                    style={{
                                        backgroundColor: `${stat.color}10`,
                                        borderColor: `${stat.color}20`,
                                    }}
                                >
                                    <stat.icon className="w-6 h-6" style={{ color: stat.color }} />
                                </div>
                            </div>
                            <p className="text-zinc-400 text-sm font-medium mb-1">{stat.title}</p>
                            <h3 className="text-3xl font-bold text-white font-mono tracking-tight">{stat.value}</h3>
                        </CardContent>
                    </Card>
                ))}
            </div>

            {/* Recent Orders */}
            <Card variant="glass">
                <CardHeader className="flex flex-row items-center justify-between pb-2">
                    <CardTitle>Recent Orders</CardTitle>
                    <Link
                        to="/portal/orders"
                        className="text-sm text-gable-green hover:underline flex items-center gap-1"
                    >
                        View All <ArrowRight className="w-4 h-4" />
                    </Link>
                </CardHeader>
                <CardContent>
                    {data?.recent_orders && data.recent_orders.length > 0 ? (
                        <div className="space-y-3">
                            {data.recent_orders.map((order) => (
                                <div
                                    key={order.id}
                                    className="flex justify-between items-center p-3 rounded-lg hover:bg-white/5 transition-colors border border-transparent hover:border-white/5"
                                >
                                    <div>
                                        <div className="font-medium text-white font-mono text-sm">
                                            {order.id.substring(0, 8).toUpperCase()}
                                        </div>
                                        <div className="text-xs text-zinc-500 mt-0.5">
                                            {new Date(order.created_at).toLocaleDateString()} · {order.lines.length} item{order.lines.length !== 1 ? 's' : ''}
                                        </div>
                                    </div>
                                    <div className="text-right flex items-center gap-3">
                                        <div>
                                            <div className="font-mono text-zinc-300 text-sm">{formatCurrency(order.total_amount)}</div>
                                            <StatusBadge status={order.status} />
                                        </div>
                                        <Link
                                            to="/portal/orders"
                                            className="p-1.5 rounded-lg hover:bg-white/5 text-zinc-500 hover:text-white transition-colors"
                                        >
                                            <ShoppingCart size={14} />
                                        </Link>
                                    </div>
                                </div>
                            ))}
                        </div>
                    ) : (
                        <p className="text-zinc-500 text-sm py-4">No orders yet.</p>
                    )}
                </CardContent>
            </Card>
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
