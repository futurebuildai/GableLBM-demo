import { Area, AreaChart, CartesianGrid, ResponsiveContainer, Tooltip, XAxis, YAxis } from 'recharts';
import { Card, CardHeader, CardTitle, CardContent } from '../ui/Card';
import type { RevenueTrendPoint } from '../../types/dashboard';

interface RevenueTrendChartProps {
    data: RevenueTrendPoint[];
    loading?: boolean;
}

interface TooltipPayload {
    value: number;
}

interface CustomTooltipProps {
    active?: boolean;
    payload?: TooltipPayload[];
    label?: string;
}

const CustomTooltip = ({ active, payload, label }: CustomTooltipProps) => {
    if (active && payload && payload.length) {
        return (
            <div className="bg-slate-steel/90 backdrop-blur-md border border-white/10 p-3 rounded-lg shadow-xl">
                <p className="text-zinc-400 text-xs mb-1">{label}</p>
                <p className="text-gable-green font-mono font-bold text-lg">
                    ${(payload[0].value).toLocaleString(undefined, { minimumFractionDigits: 2 })}
                </p>
            </div>
        );
    }
    return null;
};

export function RevenueTrendChart({ data, loading = false }: RevenueTrendChartProps) {
    if (loading) {
        return (
            <Card variant="glass" className="h-[400px]">
                <CardHeader>
                    <div className="h-6 w-48 bg-white/10 rounded animate-pulse" />
                </CardHeader>
                <CardContent className="h-[320px] flex items-end justify-between gap-2 px-6 pb-6">
                    {[...Array(12)].map((_, i) => (
                        <div key={i} className="w-full bg-white/5 rounded-t animate-pulse" style={{ height: `${(i * 13 + 20) % 60 + 20}%` }} />
                    ))}
                </CardContent>
            </Card>
        );
    }

    const formattedData = data.map(point => ({
        ...point,
        revenue: point.revenue / 100 // Convert cents to dollars
    }));

    return (
        <Card variant="glass" className="h-[400px] flex flex-col">
            <CardHeader>
                <CardTitle>Revenue Trend</CardTitle>
            </CardHeader>
            <CardContent className="flex-1 min-h-0 pt-4">
                <ResponsiveContainer width="100%" height="100%">
                    <AreaChart data={formattedData} margin={{ top: 5, right: 0, left: 0, bottom: 0 }}>
                        <defs>
                            <linearGradient id="colorRevenue" x1="0" y1="0" x2="0" y2="1">
                                <stop offset="5%" stopColor="#00FFA3" stopOpacity={0.3} />
                                <stop offset="95%" stopColor="#00FFA3" stopOpacity={0} />
                            </linearGradient>
                        </defs>
                        <CartesianGrid strokeDasharray="3 3" stroke="rgba(255,255,255,0.05)" vertical={false} />
                        <XAxis
                            dataKey="date"
                            stroke="#52525b"
                            fontSize={12}
                            tickLine={false}
                            axisLine={false}
                            tickMargin={10}
                        />
                        <YAxis
                            stroke="#52525b"
                            fontSize={12}
                            tickFormatter={(value) => `$${value.toLocaleString()}`}
                            tickLine={false}
                            axisLine={false}
                            tickMargin={10}
                        />
                        <Tooltip content={<CustomTooltip />} cursor={{ stroke: 'rgba(255,255,255,0.1)', strokeWidth: 1 }} />
                        <Area
                            type="monotone"
                            dataKey="revenue"
                            stroke="#00FFA3"
                            strokeWidth={3}
                            fillOpacity={1}
                            fill="url(#colorRevenue)"
                            activeDot={{ r: 6, strokeWidth: 0, fill: '#fff' }}
                            animationDuration={1500}
                        />
                    </AreaChart>
                </ResponsiveContainer>
            </CardContent>
        </Card>
    );
}
