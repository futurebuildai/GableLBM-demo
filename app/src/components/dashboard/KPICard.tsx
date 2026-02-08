import { TrendingUp, TrendingDown, Minus } from 'lucide-react';
import { Card, CardContent } from '../ui/Card';
import { Tooltip } from '../ui/Tooltip';
import { motion, useSpring, useTransform, useInView } from 'framer-motion';
import { useEffect, useRef } from 'react';

interface KPICardProps {
    title: string;
    value: string | number;
    subValue?: string;
    trend?: number; // Percentage change
    icon?: React.ReactNode;
    loading?: boolean;
    valueColor?: string;
}

function AnimatedNumber({ value }: { value: number }) {
    const ref = useRef<HTMLSpanElement>(null);
    const inView = useInView(ref);
    const spring = useSpring(0, { mass: 0.8, stiffness: 75, damping: 15 });
    const display = useTransform(spring, (current) =>
        current.toLocaleString(undefined, { minimumFractionDigits: 0, maximumFractionDigits: 0 })
    );

    useEffect(() => {
        if (inView) {
            spring.set(value);
        }
    }, [spring, value, inView]);

    return <motion.span ref={ref}>{display}</motion.span>;
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
            <Card variant="default" className="border-white/5 bg-slate-steel/50">
                <CardContent className="p-6 relative overflow-hidden">
                    <div className="flex justify-between items-start mb-4">
                        <div className="h-4 w-24 bg-white/10 rounded animate-pulse" />
                        <div className="h-8 w-8 bg-white/10 rounded-full animate-pulse" />
                    </div>
                    <div className="h-8 w-32 bg-white/10 rounded mb-2 animate-pulse" />
                    <div className="h-4 w-16 bg-white/10 rounded animate-pulse" />
                    <div className="absolute inset-0 -translate-x-full animate-shimmer bg-gradient-to-r from-transparent via-white/5 to-transparent z-10" />
                </CardContent>
            </Card>
        );
    }

    // Parse numeric value if string starts with $
    const numericValue = typeof value === 'string' && value.startsWith('$')
        ? parseFloat(value.replace(/[^0-9.-]+/g, ""))
        : typeof value === 'number' ? value : 0;

    const isCurrency = typeof value === 'string' && value.startsWith('$');

    return (
        <Card variant="interactive" className="group relative overflow-hidden">
            <div className="absolute top-0 right-0 p-4 opacity-10 group-hover:opacity-20 transition-opacity duration-300 transform group-hover:scale-110">
                {icon && <div className="text-current scale-150">{icon}</div>}
            </div>

            <CardContent className="p-6 relative z-10">
                <div className="flex items-center justify-between mb-2">
                    <h3 className="text-sm font-medium text-zinc-400 font-sans tracking-wide">{title}</h3>
                    {icon && (
                        <Tooltip content={`View details for ${title}`}>
                            <div className={`p-2 rounded-lg bg-white/5 text-zinc-300 group-hover:text-gable-green group-hover:bg-gable-green/10 transition-colors duration-300`}>{icon}</div>
                        </Tooltip>
                    )}
                </div>

                <div className={`text-3xl font-mono font-bold tracking-tight ${valueColor} flex items-baseline gap-1`}>
                    {isCurrency && <span>$</span>}
                    {typeof value === 'number' ? (
                        <AnimatedNumber value={value} />
                    ) : isCurrency ? (
                        // For currency, we might want to handle decimals differently, 
                        // but for now let's just animate the whole number part or display static if standard string
                        <span>{numericValue.toLocaleString(undefined, { minimumFractionDigits: 2, maximumFractionDigits: 2 })}</span>
                    ) : (
                        value
                    )}
                </div>

                <div className="flex items-center gap-2 mt-2 h-6">
                    {trend !== undefined && (
                        <div className={`flex items-center gap-1.5 text-xs font-medium px-2 py-0.5 rounded-full bg-white/5 ${getTrendColor()}`}>
                            {getTrendIcon()}
                            <span>{Math.abs(trend).toFixed(1)}%</span>
                        </div>
                    )}
                    {subValue && <span className="text-xs text-zinc-500 font-mono">{subValue}</span>}
                </div>
            </CardContent>

            {/* Hover Glow Effect */}
            <div className="absolute -bottom-4 -right-4 w-24 h-24 bg-gable-green/20 blur-3xl rounded-full opacity-0 group-hover:opacity-100 transition-opacity duration-500 pointer-events-none" />
        </Card>
    );
}
