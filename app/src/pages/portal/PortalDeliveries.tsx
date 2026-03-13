import { useState, useEffect, useCallback, useMemo } from 'react';
import { Card, CardContent } from '../../components/ui/Card';
import { Truck, Camera, RefreshCw, AlertTriangle, X, User, Clock, Phone, MapPin, FileText, ChevronDown, ChevronUp } from 'lucide-react';
import { PortalService } from '../../services/PortalService';
import type { PortalDelivery } from '../../types/portal';

type FilterTab = 'all' | 'active' | 'upcoming' | 'completed';

export const PortalDeliveries = () => {
    const [deliveries, setDeliveries] = useState<PortalDelivery[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState('');
    const [lightboxUrl, setLightboxUrl] = useState<string | null>(null);
    const [filter, setFilter] = useState<FilterTab>('all');
    const [expandedId, setExpandedId] = useState<string | null>(null);

    const fetchDeliveries = useCallback(() => {
        setLoading(true);
        setError('');
        PortalService.getDeliveries()
            .then(setDeliveries)
            .catch(err => setError(err instanceof Error ? err.message : 'Failed to load deliveries'))
            .finally(() => setLoading(false));
    }, []);

    useEffect(() => { fetchDeliveries(); }, [fetchDeliveries]);

    useEffect(() => {
        if (!lightboxUrl) return;
        const handleEsc = (e: KeyboardEvent) => {
            if (e.key === 'Escape') setLightboxUrl(null);
        };
        window.addEventListener('keydown', handleEsc);
        return () => window.removeEventListener('keydown', handleEsc);
    }, [lightboxUrl]);

    const { active, upcoming, completed } = useMemo(() => {
        const active: PortalDelivery[] = [];
        const upcoming: PortalDelivery[] = [];
        const completed: PortalDelivery[] = [];
        for (const d of deliveries) {
            if (d.status === 'OUT_FOR_DELIVERY') active.push(d);
            else if (d.status === 'PENDING') upcoming.push(d);
            else completed.push(d);
        }
        return { active, upcoming, completed };
    }, [deliveries]);

    const filtered = useMemo(() => {
        if (filter === 'active') return active;
        if (filter === 'upcoming') return upcoming;
        if (filter === 'completed') return completed;
        return deliveries;
    }, [filter, active, upcoming, completed, deliveries]);

    if (loading) {
        return (
            <div className="space-y-4">
                {[1, 2, 3].map(i => (
                    <div key={i} className="h-24 bg-white/5 rounded-2xl animate-pulse" />
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
                    onClick={fetchDeliveries}
                    className="flex items-center gap-2 px-4 py-2 rounded-lg bg-white/5 border border-white/10 text-white hover:bg-white/10 transition-colors"
                >
                    <RefreshCw size={16} /> Retry
                </button>
            </div>
        );
    }

    const tabs: { key: FilterTab; label: string; count: number }[] = [
        { key: 'all', label: 'All', count: deliveries.length },
        { key: 'active', label: 'Active', count: active.length },
        { key: 'upcoming', label: 'Upcoming', count: upcoming.length },
        { key: 'completed', label: 'Completed', count: completed.length },
    ];

    return (
        <div>
            <div className="mb-6">
                <h1 className="text-2xl font-bold text-white">Deliveries</h1>
                <p className="text-zinc-400 text-sm mt-1">Track your orders from warehouse to jobsite</p>
            </div>

            {/* Filter Tabs */}
            <div className="flex gap-1 mb-6 bg-white/5 rounded-lg p-1 w-fit">
                {tabs.map(t => (
                    <button
                        key={t.key}
                        onClick={() => setFilter(t.key)}
                        className={`px-3 py-1.5 rounded-md text-xs font-medium transition-colors ${filter === t.key ? 'bg-gable-green/20 text-gable-green' : 'text-zinc-400 hover:text-white'}`}
                    >
                        {t.label} <span className="ml-1 text-[10px] opacity-70">{t.count}</span>
                    </button>
                ))}
            </div>

            {deliveries.length === 0 ? (
                <Card variant="glass">
                    <CardContent className="p-12 text-center">
                        <Truck className="w-12 h-12 text-zinc-600 mx-auto mb-4" />
                        <p className="text-zinc-400">No deliveries yet.</p>
                    </CardContent>
                </Card>
            ) : filtered.length === 0 ? (
                <Card variant="glass">
                    <CardContent className="p-8 text-center">
                        <p className="text-zinc-500 text-sm">No deliveries match this filter.</p>
                    </CardContent>
                </Card>
            ) : (
                <div className="space-y-3">
                    {filtered.map(del => (
                        <DeliveryCard
                            key={del.id}
                            delivery={del}
                            expanded={expandedId === del.id}
                            onToggle={() => setExpandedId(expandedId === del.id ? null : del.id)}
                            onLightbox={setLightboxUrl}
                        />
                    ))}
                </div>
            )}

            {/* Lightbox */}
            {lightboxUrl && (
                <div
                    className="fixed inset-0 z-[100] bg-black/80 backdrop-blur-md flex items-center justify-center p-8"
                    onClick={() => setLightboxUrl(null)}
                >
                    <div className="relative max-w-3xl w-full">
                        <button
                            onClick={() => setLightboxUrl(null)}
                            className="absolute -top-12 right-0 p-2 rounded-lg bg-white/10 text-white hover:bg-white/20 transition-colors"
                        >
                            <X size={20} />
                        </button>
                        <img
                            src={lightboxUrl}
                            alt="Proof of Delivery"
                            className="w-full h-auto rounded-2xl border border-white/10 shadow-2xl"
                        />
                        <p className="text-center text-zinc-400 text-sm mt-4">Proof of Delivery Photo</p>
                    </div>
                </div>
            )}
        </div>
    );
};

function etaLabel(eta: string | null): string | null {
    if (!eta) return null;
    const arrival = new Date(eta);
    const now = new Date();
    const diffMs = arrival.getTime() - now.getTime();
    if (diffMs <= 0) return 'Arriving now';
    const diffMin = Math.round(diffMs / 60000);
    if (diffMin < 120) return `Arriving in ${diffMin} min`;
    return `ETA ${arrival.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })}`;
}

const STATUS_STEPS = ['PENDING', 'OUT_FOR_DELIVERY', 'DELIVERED'] as const;

function StatusTimeline({ status }: { status: string }) {
    const currentIdx = STATUS_STEPS.indexOf(status as typeof STATUS_STEPS[number]);
    const isFailed = status === 'FAILED' || status === 'PARTIAL';

    return (
        <div className="flex items-center gap-1 w-full max-w-xs">
            {STATUS_STEPS.map((step, i) => {
                const isActive = i <= currentIdx;
                const isCurrent = i === currentIdx;
                return (
                    <div key={step} className="flex items-center flex-1">
                        <div className={`w-3 h-3 rounded-full border-2 shrink-0 transition-colors ${
                            isFailed && isCurrent ? 'border-red-400 bg-red-400/30' :
                            isActive ? 'border-gable-green bg-gable-green/30' :
                            'border-zinc-600 bg-transparent'
                        }`} />
                        {i < STATUS_STEPS.length - 1 && (
                            <div className={`flex-1 h-0.5 mx-1 transition-colors ${
                                i < currentIdx ? 'bg-gable-green/50' : 'bg-zinc-700'
                            }`} />
                        )}
                    </div>
                );
            })}
            <span className="text-[9px] text-zinc-500 uppercase ml-2 whitespace-nowrap">
                {STATUS_STEPS[Math.max(0, currentIdx)] ?? status}
            </span>
        </div>
    );
}

function DeliveryCard({ delivery: del, expanded, onToggle, onLightbox }: {
    delivery: PortalDelivery;
    expanded: boolean;
    onToggle: () => void;
    onLightbox: (url: string) => void;
}) {
    const isActive = del.status === 'OUT_FOR_DELIVERY';
    const isCompleted = del.status === 'DELIVERED' || del.status === 'FAILED' || del.status === 'PARTIAL';
    const eta = etaLabel(del.estimated_arrival);

    return (
        <Card variant="glass" noPadding>
            {/* Header row */}
            <button onClick={onToggle} className="w-full text-left p-4 hover:bg-white/5 transition-colors">
                <div className="flex items-center justify-between">
                    <div className="flex items-center gap-4">
                        <div
                            className="w-10 h-10 rounded-lg flex items-center justify-center shrink-0"
                            style={{ backgroundColor: deliveryStatusColor(del.status).bg }}
                        >
                            <Truck size={18} style={{ color: deliveryStatusColor(del.status).fg }} />
                        </div>
                        <div>
                            <div className="flex items-center gap-2">
                                <span className="font-mono text-sm font-medium text-white">
                                    DEL-{del.id.substring(0, 8).toUpperCase()}
                                </span>
                                <DeliveryStatusBadge status={del.status} />
                            </div>
                            <div className="text-xs text-zinc-500 mt-0.5 flex items-center gap-2">
                                <span className="flex items-center gap-1">
                                    <FileText size={10} />
                                    Order {del.order_number?.substring(0, 8).toUpperCase() || del.order_id.substring(0, 8).toUpperCase()}
                                </span>
                                {del.scheduled_date && (
                                    <span className="flex items-center gap-1">
                                        <Clock size={10} />
                                        {new Date(del.scheduled_date).toLocaleDateString()}
                                    </span>
                                )}
                                {!del.scheduled_date && (
                                    <span>{new Date(del.created_at).toLocaleDateString()}</span>
                                )}
                            </div>
                        </div>
                    </div>

                    <div className="flex items-center gap-3">
                        {/* ETA for active deliveries */}
                        {isActive && eta && (
                            <span className="text-xs font-semibold text-sky-400 bg-sky-500/10 px-2 py-1 rounded border border-sky-500/20">
                                {eta}
                            </span>
                        )}
                        {/* Stop progress */}
                        {del.stop_sequence != null && del.total_stops != null && (
                            <span className="text-[10px] text-zinc-500 font-mono">
                                Stop {del.stop_sequence} of {del.total_stops}
                            </span>
                        )}
                        {expanded ? <ChevronUp size={16} className="text-zinc-500" /> : <ChevronDown size={16} className="text-zinc-500" />}
                    </div>
                </div>

                {/* Status timeline for active deliveries */}
                {isActive && (
                    <div className="mt-3 pl-14">
                        <StatusTimeline status={del.status} />
                    </div>
                )}
            </button>

            {/* Expanded details */}
            {expanded && (
                <div className="border-t border-white/5 p-4 pl-[4.5rem] space-y-3">
                    {/* Driver & vehicle info */}
                    {(del.driver_name || del.vehicle_name) && (
                        <div className="flex flex-wrap gap-4 text-sm">
                            {del.driver_name && (
                                <div className="flex items-center gap-2 text-zinc-300">
                                    <User size={14} className="text-zinc-500" />
                                    <span>{del.driver_name}</span>
                                    {del.driver_phone && (
                                        <a href={`tel:${del.driver_phone}`} className="flex items-center gap-1 text-sky-400 hover:text-sky-300">
                                            <Phone size={12} />
                                            <span className="text-xs">{del.driver_phone}</span>
                                        </a>
                                    )}
                                </div>
                            )}
                            {del.vehicle_name && (
                                <div className="flex items-center gap-2 text-zinc-400">
                                    <Truck size={14} className="text-zinc-500" />
                                    <span>{del.vehicle_name}</span>
                                </div>
                            )}
                        </div>
                    )}

                    {/* Address */}
                    {del.delivery_address && (
                        <div className="flex items-start gap-2 text-sm text-zinc-400">
                            <MapPin size={14} className="text-zinc-500 mt-0.5 shrink-0" />
                            <span>{del.delivery_address}</span>
                        </div>
                    )}

                    {/* Delivery instructions */}
                    {del.delivery_instructions && (
                        <div className="text-sm bg-amber-500/5 text-amber-500/90 p-3 rounded-lg border border-amber-500/10 flex gap-2">
                            <span className="font-bold text-[10px] px-1.5 py-0.5 bg-amber-500/20 rounded h-fit shrink-0">NOTE</span>
                            <span className="italic">{del.delivery_instructions}</span>
                        </div>
                    )}

                    {/* POD info for completed deliveries */}
                    {isCompleted && (del.pod_signed_by || del.pod_proof_url) && (
                        <div className="flex items-center gap-4 pt-2 border-t border-white/5">
                            {del.pod_signed_by && (
                                <div className="flex items-center gap-1 text-xs text-zinc-400">
                                    <User size={12} />
                                    <span>Signed by {del.pod_signed_by}</span>
                                </div>
                            )}
                            {del.pod_timestamp && (
                                <div className="flex items-center gap-1 text-xs text-zinc-500">
                                    <Clock size={12} />
                                    <span>{new Date(del.pod_timestamp).toLocaleString()}</span>
                                </div>
                            )}
                            {del.pod_proof_url && (
                                <button
                                    onClick={(e) => { e.stopPropagation(); onLightbox(del.pod_proof_url!); }}
                                    className="w-10 h-10 rounded-lg overflow-hidden border border-white/10 hover:border-gable-green/50 transition-colors relative group"
                                >
                                    <img src={del.pod_proof_url} alt="POD" className="w-full h-full object-cover" />
                                    <div className="absolute inset-0 bg-black/40 flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity">
                                        <Camera size={12} className="text-white" />
                                    </div>
                                </button>
                            )}
                        </div>
                    )}
                </div>
            )}
        </Card>
    );
}

const deliveryStatusColor = (status: string): { fg: string; bg: string } => {
    const map: Record<string, { fg: string; bg: string }> = {
        PENDING: { fg: '#F59E0B', bg: 'rgba(245,158,11,0.1)' },
        OUT_FOR_DELIVERY: { fg: '#38BDF8', bg: 'rgba(56,189,248,0.1)' },
        DELIVERED: { fg: '#00FFA3', bg: 'rgba(0,255,163,0.1)' },
        FAILED: { fg: '#F43F5E', bg: 'rgba(244,63,94,0.1)' },
        PARTIAL: { fg: '#A78BFA', bg: 'rgba(167,139,250,0.1)' },
    };
    return map[status] || map.PENDING;
};

const DeliveryStatusBadge = ({ status }: { status: string }) => {
    const colors: Record<string, string> = {
        PENDING: 'bg-amber-500/10 text-amber-400 border-amber-500/20',
        OUT_FOR_DELIVERY: 'bg-blue-500/10 text-blue-400 border-blue-500/20',
        DELIVERED: 'bg-emerald-500/10 text-emerald-400 border-emerald-500/20',
        FAILED: 'bg-red-500/10 text-red-400 border-red-500/20',
        PARTIAL: 'bg-purple-500/10 text-purple-400 border-purple-500/20',
    };
    return (
        <span className={`inline-block px-2 py-0.5 rounded text-[10px] uppercase tracking-wider font-semibold border whitespace-nowrap ${colors[status] || colors.PENDING}`}>
            {status.replace(/_/g, ' ')}
        </span>
    );
};
