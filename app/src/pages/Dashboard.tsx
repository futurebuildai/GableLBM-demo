import { useEffect, useState, useCallback } from 'react';
import { RefreshCw, DollarSign, ShoppingCart, Truck, CreditCard } from 'lucide-react';
import { DashboardService } from '../services/DashboardService';
import { KPICard } from '../components/dashboard/KPICard';
import { RevenueTrendChart } from '../components/dashboard/RevenueTrendChart';
import { OrderStatusChart } from '../components/dashboard/OrderStatusChart';
import { TopCustomersTable } from '../components/dashboard/TopCustomersTable';
import { InventoryAlertsWidget } from '../components/dashboard/InventoryAlertsWidget';
import { RecentOrdersFeed } from '../components/dashboard/RecentOrdersFeed';
import type {
    DashboardSummary,
    InventoryAlert,
    TopCustomer,
    OrderActivity,
    RevenueTrendPoint,
} from '../types/dashboard';

const REFRESH_INTERVAL = 60000; // 60 seconds

export const Dashboard = () => {
    const [summary, setSummary] = useState<DashboardSummary | null>(null);
    const [inventoryAlerts, setInventoryAlerts] = useState<InventoryAlert[]>([]);
    const [topCustomers, setTopCustomers] = useState<TopCustomer[]>([]);
    const [orderActivity, setOrderActivity] = useState<OrderActivity | null>(null);
    const [revenueTrend, setRevenueTrend] = useState<RevenueTrendPoint[]>([]);
    const [loading, setLoading] = useState(true);
    const [lastRefresh, setLastRefresh] = useState<Date>(new Date());
    const [refreshing, setRefreshing] = useState(false);

    const fetchDashboardData = useCallback(async (showSpinner = false) => {
        if (showSpinner) setRefreshing(true);
        try {
            const [summaryData, alertsData, customersData, activityData, trendData] = await Promise.all([
                DashboardService.getSummary(),
                DashboardService.getInventoryAlerts(),
                DashboardService.getTopCustomers(),
                DashboardService.getOrderActivity(),
                DashboardService.getRevenueTrend(),
            ]);
            setSummary(summaryData);
            setInventoryAlerts(alertsData);
            setTopCustomers(customersData);
            setOrderActivity(activityData);
            setRevenueTrend(trendData);
            setLastRefresh(new Date());
        } catch (error) {
            console.error('Failed to fetch dashboard data:', error);
        } finally {
            setLoading(false);
            setRefreshing(false);
        }
    }, []);

    useEffect(() => {
        fetchDashboardData();
        const interval = setInterval(() => fetchDashboardData(), REFRESH_INTERVAL);
        return () => clearInterval(interval);
    }, [fetchDashboardData]);

    const formatCurrency = (cents: number) => {
        return `$${(cents / 100).toLocaleString(undefined, { minimumFractionDigits: 2 })}`;
    };

    return (
        <div className="space-y-6">
            {/* Header */}
            <div className="flex items-center justify-between">
                <div>
                    <h1 className="text-3xl font-bold tracking-tight text-white">Dashboard</h1>
                    <p className="text-zinc-500 text-sm mt-1">
                        Last updated: {lastRefresh.toLocaleTimeString()}
                    </p>
                </div>
                <button
                    onClick={() => fetchDashboardData(true)}
                    disabled={refreshing}
                    className="flex items-center gap-2 px-4 py-2 rounded-lg bg-zinc-800 border border-white/10 text-zinc-300 hover:bg-zinc-700 hover:border-white/20 transition-all duration-200 disabled:opacity-50"
                >
                    <RefreshCw className={`w-4 h-4 ${refreshing ? 'animate-spin' : ''}`} />
                    Refresh
                </button>
            </div>

            {/* KPI Cards */}
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
                <KPICard
                    title="Today's Revenue"
                    value={summary ? formatCurrency(summary.today_revenue) : '$0.00'}
                    trend={summary?.today_revenue_change}
                    icon={<DollarSign className="w-5 h-5" />}
                    loading={loading}
                    valueColor="text-emerald-400"
                />
                <KPICard
                    title="Active Orders"
                    value={summary?.active_orders ?? 0}
                    icon={<ShoppingCart className="w-5 h-5" />}
                    loading={loading}
                />
                <KPICard
                    title="Pending Dispatch"
                    value={summary?.pending_dispatch ?? 0}
                    icon={<Truck className="w-5 h-5" />}
                    loading={loading}
                    valueColor="text-blue-400"
                />
                <KPICard
                    title="Outstanding AR"
                    value={summary ? formatCurrency(summary.outstanding_ar) : '$0.00'}
                    subValue={summary ? `${summary.outstanding_ar_count} invoices` : undefined}
                    icon={<CreditCard className="w-5 h-5" />}
                    loading={loading}
                    valueColor="text-amber-400"
                />
            </div>

            {/* Charts Row */}
            <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
                <div className="lg:col-span-2">
                    <RevenueTrendChart data={revenueTrend} loading={loading} />
                </div>
                <div>
                    <OrderStatusChart
                        statusBreakdown={orderActivity?.status_breakdown ?? {}}
                        loading={loading}
                    />
                </div>
            </div>

            {/* Widgets Row */}
            <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
                <TopCustomersTable customers={topCustomers} loading={loading} />
                <InventoryAlertsWidget alerts={inventoryAlerts} loading={loading} />
                <RecentOrdersFeed orders={orderActivity?.recent_orders ?? []} loading={loading} />
            </div>
        </div>
    );
};
