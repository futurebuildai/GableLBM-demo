import {
    AreaChart,
    Area,
    XAxis,
    YAxis,
    Tooltip,
    ResponsiveContainer,
} from 'recharts';
import type { RevenueTrendPoint } from '../../types/dashboard';

interface RevenueTrendChartProps {
    data: RevenueTrendPoint[];
    loading?: boolean;
}

export function RevenueTrendChart({ data, loading }: RevenueTrendChartProps) {
    if (loading) {
        return (
            <div className="h-64 bg-zinc-900 rounded-lg border border-white/10 animate-pulse flex items-center justify-center">
                <div className="text-zinc-600">Loading chart...</div>
            </div>
        );
    }

    // Format data for chart - convert cents to dollars
    const chartData = data.map((point) => ({
        date: point.date.slice(5), // Show MM-DD
        revenue: point.revenue / 100,
    }));

    return (
        <div className="p-4 bg-zinc-900 rounded-lg border border-white/10">
            <h3 className="text-sm font-medium text-zinc-400 mb-4">Revenue Trend (7 Days)</h3>
            <div className="h-56">
                <ResponsiveContainer width="100%" height="100%">
                    <AreaChart data={chartData} margin={{ top: 5, right: 20, left: 10, bottom: 5 }}>
                        <defs>
                            <linearGradient id="revenueGradient" x1="0" y1="0" x2="0" y2="1">
                                <stop offset="5%" stopColor="#00FFA3" stopOpacity={0.3} />
                                <stop offset="95%" stopColor="#00FFA3" stopOpacity={0} />
                            </linearGradient>
                        </defs>
                        <XAxis
                            dataKey="date"
                            axisLine={false}
                            tickLine={false}
                            tick={{ fill: '#71717a', fontSize: 12 }}
                        />
                        <YAxis
                            axisLine={false}
                            tickLine={false}
                            tick={{ fill: '#71717a', fontSize: 12 }}
                            tickFormatter={(value) => `$${value.toLocaleString()}`}
                        />
                        <Tooltip
                            contentStyle={{
                                background: '#161821',
                                border: '1px solid rgba(255,255,255,0.1)',
                                borderRadius: '8px',
                            }}
                            labelStyle={{ color: '#fff' }}
                            formatter={(value) => [`$${(value as number)?.toLocaleString() ?? 0}`, 'Revenue']}
                        />
                        <Area
                            type="monotone"
                            dataKey="revenue"
                            stroke="#00FFA3"
                            strokeWidth={2}
                            fill="url(#revenueGradient)"
                        />
                    </AreaChart>
                </ResponsiveContainer>
            </div>
        </div>
    );
}
