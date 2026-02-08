import { TrendingUp, TrendingDown, Minus } from 'lucide-react';

interface KPICardProps {
    title: string;
    value: string | number;
    subValue?: string;
    trend?: number; // Percentage change
    icon?: React.ReactNode;
    loading?: boolean;
    valueColor?: string;
}

export function KPICard({
    title,
    value,
    subValue,
    trend,
    icon,
    loading = false,
    valueColor = 'text-white',
}: KPICardProps) {
    const getTrendIcon = () => {
        if (trend === undefined || trend === null) return null;
        if (trend > 0) return <TrendingUp className="w-4 h-4 text-emerald-400" />;
        if (trend < 0) return <TrendingDown className="w-4 h-4 text-rose-400" />;
        return <Minus className="w-4 h-4 text-zinc-500" />;
    };

    const getTrendColor = () => {
        if (trend === undefined || trend === null) return 'text-zinc-500';
        if (trend > 0) return 'text-emerald-400';
        if (trend < 0) return 'text-rose-400';
        return 'text-zinc-500';
    };

    if (loading) {
        return (
            <div className="p-6 rounded-lg bg-zinc-900 border border-white/10 animate-pulse">
                <div className="h-4 w-24 bg-zinc-800 rounded mb-3" />
                <div className="h-8 w-32 bg-zinc-800 rounded" />
            </div>
        );
    }

    return (
        <div className="p-6 rounded-lg bg-zinc-900 border border-white/10 hover:border-white/20 transition-all duration-200 hover:-translate-y-0.5">
            <div className="flex items-center justify-between mb-1">
                <h3 className="text-sm font-medium text-zinc-400">{title}</h3>
                {icon && <span className="text-zinc-500">{icon}</span>}
            </div>
            <div className={`text-2xl font-mono font-bold ${valueColor}`}>{value}</div>
            <div className="flex items-center gap-2 mt-1">
                {trend !== undefined && (
                    <div className={`flex items-center gap-1 text-sm ${getTrendColor()}`}>
                        {getTrendIcon()}
                        <span>{Math.abs(trend).toFixed(1)}%</span>
                    </div>
                )}
                {subValue && <span className="text-xs text-zinc-500">{subValue}</span>}
            </div>
        </div>
    );
}
