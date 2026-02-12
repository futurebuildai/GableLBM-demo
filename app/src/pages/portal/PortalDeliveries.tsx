import { useState, useEffect, useCallback } from 'react';
import { Card, CardContent } from '../../components/ui/Card';
import { Truck, Camera, RefreshCw, AlertTriangle, X, User, Clock } from 'lucide-react';
import { PortalService } from '../../services/PortalService';
import type { PortalDelivery } from '../../types/portal';

export const PortalDeliveries = () => {
    const [deliveries, setDeliveries] = useState<PortalDelivery[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState('');
    const [lightboxUrl, setLightboxUrl] = useState<string | null>(null);

    const fetchDeliveries = useCallback(() => {
        setLoading(true);
        setError('');
        PortalService.getDeliveries()
            .then(setDeliveries)
            .catch(err => setError(err instanceof Error ? err.message : 'Failed to load deliveries'))
            .finally(() => setLoading(false));
    }, []);

    // eslint-disable-next-line react-hooks/set-state-in-effect -- async fetch→setState is standard pattern
    useEffect(() => { fetchDeliveries(); }, [fetchDeliveries]);

    // Close lightbox on Escape key
    useEffect(() => {
        if (!lightboxUrl) return;
        const handleEsc = (e: KeyboardEvent) => {
            if (e.key === 'Escape') setLightboxUrl(null);
        };
        window.addEventListener('keydown', handleEsc);
        return () => window.removeEventListener('keydown', handleEsc);
    }, [lightboxUrl]);

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

    return (
        <div>
            <div className="mb-6">
                <h1 className="text-2xl font-bold text-white">Deliveries</h1>
                <p className="text-zinc-400 text-sm mt-1">{deliveries.length} deliver{deliveries.length !== 1 ? 'ies' : 'y'} found</p>
            </div>

            {deliveries.length === 0 ? (
                <Card variant="glass">
                    <CardContent className="p-12 text-center">
                        <Truck className="w-12 h-12 text-zinc-600 mx-auto mb-4" />
                        <p className="text-zinc-400">No deliveries yet.</p>
                    </CardContent>
                </Card>
            ) : (
                <div className="space-y-3">
                    {deliveries.map(del => (
                        <Card key={del.id} variant="glass" noPadding>
                            <div className="flex items-center justify-between p-4 hover:bg-white/5 transition-colors">
                                <div className="flex items-center gap-4">
                                    <div
                                        className="w-10 h-10 rounded-lg flex items-center justify-center"
                                        style={{ backgroundColor: deliveryStatusColor(del.status).bg }}
                                    >
                                        <Truck size={18} style={{ color: deliveryStatusColor(del.status).fg }} />
                                    </div>
                                    <div>
                                        <div className="font-mono text-sm font-medium text-white">
                                            DEL-{del.id.substring(0, 8).toUpperCase()}
                                        </div>
                                        <div className="text-xs text-zinc-500 mt-0.5">
                                            Order: {del.order_number?.substring(0, 8).toUpperCase() || del.order_id.substring(0, 8).toUpperCase()}
                                            {' · '}
                                            {new Date(del.created_at).toLocaleDateString()}
                                        </div>
                                    </div>
                                </div>

                                <div className="flex items-center gap-4">
                                    {/* POD Info */}
                                    <div className="flex items-center gap-3">
                                        {del.pod_signed_by && (
                                            <div className="flex items-center gap-1 text-xs text-zinc-400">
                                                <User size={12} />
                                                <span>{del.pod_signed_by}</span>
                                            </div>
                                        )}
                                        {del.pod_timestamp && (
                                            <div className="flex items-center gap-1 text-xs text-zinc-500">
                                                <Clock size={12} />
                                                <span>{new Date(del.pod_timestamp).toLocaleString()}</span>
                                            </div>
                                        )}
                                    </div>

                                    <DeliveryStatusBadge status={del.status} />

                                    {/* POD Photo */}
                                    {del.pod_proof_url ? (
                                        <button
                                            onClick={() => setLightboxUrl(del.pod_proof_url)}
                                            className="w-12 h-12 rounded-lg overflow-hidden border border-white/10 hover:border-gable-green/50 transition-colors relative group"
                                        >
                                            <img
                                                src={del.pod_proof_url}
                                                alt="Proof of Delivery"
                                                className="w-full h-full object-cover"
                                            />
                                            <div className="absolute inset-0 bg-black/40 flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity">
                                                <Camera size={14} className="text-white" />
                                            </div>
                                        </button>
                                    ) : (
                                        <div className="w-12 h-12 rounded-lg border border-white/5 flex items-center justify-center bg-white/5">
                                            <Camera size={14} className="text-zinc-600" />
                                        </div>
                                    )}
                                </div>
                            </div>
                        </Card>
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
