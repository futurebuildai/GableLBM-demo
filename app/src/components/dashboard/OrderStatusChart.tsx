import { PieChart, Pie, Cell, ResponsiveContainer, Legend, Tooltip } from 'recharts';

interface OrderStatusChartProps {
    statusBreakdown: Record<string, number>;
    loading?: boolean;
}

const STATUS_COLORS: Record<string, string> = {
    PENDING: '#FBBF24',      // Amber
    CONFIRMED: '#38BDF8',    // Blueprint Blue
    PROCESSING: '#818CF8',   // Indigo
    READY: '#34D399',        // Emerald
    ALLOCATED: '#60A5FA',    // Light Blue
    COMPLETED: '#00FFA3',    // Gable Green
    CANCELLED: '#F43F5E',    // Safety Red
};

export function OrderStatusChart({ statusBreakdown, loading }: OrderStatusChartProps) {
    if (loading) {
        return (
            <div className="h-64 bg-zinc-900 rounded-lg border border-white/10 animate-pulse flex items-center justify-center">
                <div className="text-zinc-600">Loading chart...</div>
            </div>
        );
    }

    const data = Object.entries(statusBreakdown).map(([name, value]) => ({
        name,
        value,
        fill: STATUS_COLORS[name] || '#6B7280',
    }));

    if (data.length === 0) {
        return (
            <div className="p-4 bg-zinc-900 rounded-lg border border-white/10">
                <h3 className="text-sm font-medium text-zinc-400 mb-4">Order Status (30 Days)</h3>
                <div className="h-48 flex items-center justify-center text-zinc-500">
                    No order data available
                </div>
            </div>
        );
    }

    return (
        <div className="p-4 bg-zinc-900 rounded-lg border border-white/10">
            <h3 className="text-sm font-medium text-zinc-400 mb-4">Order Status (30 Days)</h3>
            <div className="h-48">
                <ResponsiveContainer width="100%" height="100%">
                    <PieChart>
                        <Pie
                            data={data}
                            cx="50%"
                            cy="50%"
                            innerRadius={50}
                            outerRadius={70}
                            paddingAngle={2}
                            dataKey="value"
                        >
                            {data.map((entry, index) => (
                                <Cell key={`cell-${index}`} fill={entry.fill} />
                            ))}
                        </Pie>
                        <Tooltip
                            contentStyle={{
                                background: '#161821',
                                border: '1px solid rgba(255,255,255,0.1)',
                                borderRadius: '8px',
                            }}
                            labelStyle={{ color: '#fff' }}
                        />
                        <Legend
                            wrapperStyle={{ fontSize: '12px' }}
                            formatter={(value) => <span className="text-zinc-300">{value}</span>}
                        />
                    </PieChart>
                </ResponsiveContainer>
            </div>
        </div>
    );
}
