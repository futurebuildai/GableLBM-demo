import { Cell, Legend, Pie, PieChart, ResponsiveContainer, Tooltip } from 'recharts';
import { Card, CardHeader, CardTitle, CardContent } from '../ui/Card';

interface OrderStatusChartProps {
    statusBreakdown: Record<string, number>;
    loading?: boolean;
}

const COLORS = ['#00FFA3', '#38BDF8', '#818cf8', '#fbbf24', '#f43f5e', '#a1a1aa'];

// Helper to map status to readable names
const formatStatus = (status: string) => {
    return status.split('_').map(word => word.charAt(0).toUpperCase() + word.slice(1)).join(' ');
};

interface TooltipPayload {
    name: string;
    value: number | string;
}

const CustomTooltip = ({ active, payload }: { active?: boolean; payload?: TooltipPayload[] }) => {
    if (active && payload && payload.length) {
        return (
            <div className="bg-slate-steel/90 backdrop-blur-md border border-white/10 p-2 rounded-lg shadow-xl text-xs">
                <span className="text-zinc-400">{payload[0].name}: </span>
                <span className="font-mono font-bold text-white">{payload[0].value}</span>
            </div>
        );
    }
    return null;
};

export function OrderStatusChart({ statusBreakdown, loading = false }: OrderStatusChartProps) {
    if (loading) {
        return (
            <Card variant="glass" className="h-[400px]">
                <CardHeader>
                    <div className="h-6 w-32 bg-white/10 rounded animate-pulse" />
                </CardHeader>
                <CardContent className="flex items-center justify-center h-[320px]">
                    <div className="h-48 w-48 rounded-full border-4 border-white/10 border-t-gable-green/50 animate-spin" />
                </CardContent>
            </Card>
        );
    }

    const data = Object.entries(statusBreakdown).map(([status, count]) => ({
        name: formatStatus(status),
        value: count,
        statusKey: status
    })).filter(item => item.value > 0); // Only show statuses with orders

    const totalOrders = data.reduce((acc, curr) => acc + curr.value, 0);

    return (
        <Card variant="glass" className="h-[400px] flex flex-col">
            <CardHeader>
                <CardTitle>Order Status</CardTitle>
            </CardHeader>
            <CardContent className="flex-1 min-h-0 relative">
                <div className="absolute inset-0 flex items-center justify-center flex-col pointer-events-none">
                    <span className="text-4xl font-bold text-white font-mono">{totalOrders}</span>
                    <span className="text-xs text-zinc-500 uppercase tracking-widest">Active</span>
                </div>
                <ResponsiveContainer width="100%" height="100%">
                    <PieChart>
                        <Pie
                            data={data}
                            innerRadius={80}
                            outerRadius={110}
                            paddingAngle={5}
                            dataKey="value"
                            stroke="none"
                        >
                            {data.map((_, index) => (
                                <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
                            ))}
                        </Pie>
                        <Tooltip content={<CustomTooltip />} />
                        <Legend
                            verticalAlign="bottom"
                            height={36}
                            iconType="circle"
                            formatter={(value) => <span className="text-zinc-400 text-xs ml-1">{value}</span>}
                        />
                    </PieChart>
                </ResponsiveContainer>
            </CardContent>
        </Card>
    );
}
