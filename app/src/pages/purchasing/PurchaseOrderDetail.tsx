import { useEffect, useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { PurchaseOrderService } from '../../services/PurchaseOrderService';
import { LocationService } from '../../services/LocationService';
import type { PurchaseOrder, PurchaseOrderLine } from '../../types/purchaseOrder';
import type { Location } from '../../types/location';
import { useToast } from '../../components/ui/ToastContext';
import { PageTransition } from '../../components/ui/PageTransition';
import { Card, CardContent } from '../../components/ui/Card';
import { Button } from '../../components/ui/Button';
import { ArrowLeft, Send, PackageCheck } from 'lucide-react';

export function PurchaseOrderDetail() {
    const { id } = useParams<{ id: string }>();
    const navigate = useNavigate();
    const { showToast } = useToast();
    const [po, setPO] = useState<PurchaseOrder | null>(null);
    const [locations, setLocations] = useState<Location[]>([]);
    const [receiving, setReceiving] = useState(false);
    const [receiveData, setReceiveData] = useState<Record<string, { qty: number; locationId: string }>>({});
    const [isSubmitting, setIsSubmitting] = useState(false);

    useEffect(() => {
        if (id) {
            loadPO(id);
            LocationService.listLocations().then(setLocations);
        }
    }, [id]);

    const loadPO = async (poId: string) => {
        try {
            const data = await PurchaseOrderService.getPO(poId);
            setPO(data);
            // Initialize receive data
            const initial: Record<string, { qty: number; locationId: string }> = {};
            (data.lines || []).forEach((line: PurchaseOrderLine) => {
                initial[line.id] = { qty: line.quantity - line.qty_received, locationId: '' };
            });
            setReceiveData(initial);
        } catch (err) {
            console.error(err);
            showToast('Failed to load purchase order', 'error');
        }
    };

    const handleSubmitPO = async () => {
        if (!po) return;
        setIsSubmitting(true);
        try {
            await PurchaseOrderService.submitPO(po.id);
            showToast('Purchase order submitted to vendor', 'success');
            loadPO(po.id);
        } catch (err) {
            console.error(err);
            showToast('Failed to submit PO', 'error');
        } finally {
            setIsSubmitting(false);
        }
    };

    const handleReceive = async () => {
        if (!po) return;
        setIsSubmitting(true);
        try {
            const lines = Object.entries(receiveData)
                .filter(([, v]) => v.qty > 0 && v.locationId)
                .map(([lineId, v]) => ({
                    line_id: lineId,
                    qty_received: v.qty,
                    location_id: v.locationId,
                }));

            if (lines.length === 0) {
                showToast('Enter quantities and select locations to receive', 'error');
                setIsSubmitting(false);
                return;
            }

            await PurchaseOrderService.receivePO(po.id, { lines });
            showToast('Items received into inventory', 'success');
            setReceiving(false);
            loadPO(po.id);
        } catch (err) {
            console.error(err);
            showToast('Failed to receive items', 'error');
        } finally {
            setIsSubmitting(false);
        }
    };

    if (!po) {
        return (
            <div className="p-12 flex justify-center">
                <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-gable-green"></div>
            </div>
        );
    }

    const canSubmit = po.status === 'DRAFT';
    const canReceive = po.status === 'SENT' || po.status === 'PARTIAL';

    return (
        <PageTransition>
            <div className="flex items-center gap-4 mb-6">
                <button onClick={() => navigate('/purchasing')} className="p-2 rounded-full bg-white/5 hover:bg-white/10 text-zinc-400 transition-colors">
                    <ArrowLeft className="w-5 h-5" />
                </button>
                <div className="flex-1">
                    <h1 className="text-2xl font-bold text-white">PO #{po.id.slice(0, 8)}</h1>
                    <p className="text-sm text-zinc-400">Status: <span className="font-bold uppercase">{po.status}</span></p>
                </div>
                <div className="flex gap-3">
                    {canSubmit && (
                        <Button onClick={handleSubmitPO} disabled={isSubmitting} isLoading={isSubmitting}>
                            <Send className="w-4 h-4 mr-2" />
                            Submit to Vendor
                        </Button>
                    )}
                    {canReceive && !receiving && (
                        <Button onClick={() => setReceiving(true)}>
                            <PackageCheck className="w-4 h-4 mr-2" />
                            Receive Items
                        </Button>
                    )}
                </div>
            </div>

            <Card variant="glass">
                <CardContent className="p-0">
                    <table className="w-full text-sm text-left">
                        <thead className="bg-white/5 text-zinc-400 uppercase tracking-wider text-xs font-semibold">
                            <tr>
                                <th className="px-6 py-4">Description</th>
                                <th className="px-6 py-4 text-right">Ordered</th>
                                <th className="px-6 py-4 text-right">Received</th>
                                <th className="px-6 py-4 text-right">Unit Cost</th>
                                <th className="px-6 py-4 text-right">Line Total</th>
                                {receiving && <th className="px-6 py-4 text-right">Receive Qty</th>}
                                {receiving && <th className="px-6 py-4">Location</th>}
                            </tr>
                        </thead>
                        <tbody className="divide-y divide-white/5">
                            {(po.lines || []).map((line) => {
                                const remaining = line.quantity - line.qty_received;
                                return (
                                    <tr key={line.id} className="hover:bg-white/5 transition-colors">
                                        <td className="px-6 py-4">
                                            <span className="text-white">{line.description}</span>
                                            {line.product_id && (
                                                <span className="text-zinc-500 text-xs ml-2">({line.product_id.slice(0, 8)})</span>
                                            )}
                                        </td>
                                        <td className="px-6 py-4 text-right font-mono text-zinc-300">{line.quantity}</td>
                                        <td className="px-6 py-4 text-right font-mono">
                                            <span className={line.qty_received >= line.quantity ? 'text-emerald-400' : 'text-amber-400'}>
                                                {line.qty_received}
                                            </span>
                                        </td>
                                        <td className="px-6 py-4 text-right font-mono text-zinc-300">${line.cost.toFixed(2)}</td>
                                        <td className="px-6 py-4 text-right font-mono text-emerald-400 font-bold">
                                            ${(line.quantity * line.cost).toFixed(2)}
                                        </td>
                                        {receiving && (
                                            <td className="px-6 py-4 text-right">
                                                <input
                                                    type="number"
                                                    min="0"
                                                    max={remaining}
                                                    step="any"
                                                    value={receiveData[line.id]?.qty || 0}
                                                    onChange={(e) => setReceiveData(prev => ({
                                                        ...prev,
                                                        [line.id]: { ...prev[line.id], qty: Number(e.target.value) },
                                                    }))}
                                                    className="w-24 bg-black/20 border border-white/10 rounded px-2 py-1 text-white font-mono text-right focus:border-[#00FFA3] outline-none"
                                                    disabled={remaining <= 0}
                                                />
                                            </td>
                                        )}
                                        {receiving && (
                                            <td className="px-6 py-4">
                                                <select
                                                    value={receiveData[line.id]?.locationId || ''}
                                                    onChange={(e) => setReceiveData(prev => ({
                                                        ...prev,
                                                        [line.id]: { ...prev[line.id], locationId: e.target.value },
                                                    }))}
                                                    className="w-40 bg-black/20 border border-white/10 rounded px-2 py-1 text-white focus:border-[#00FFA3] outline-none"
                                                    disabled={remaining <= 0}
                                                >
                                                    <option value="">Select...</option>
                                                    {locations.map(loc => (
                                                        <option key={loc.id} value={loc.id}>
                                                            {loc.path || loc.code}
                                                        </option>
                                                    ))}
                                                </select>
                                            </td>
                                        )}
                                    </tr>
                                );
                            })}
                        </tbody>
                    </table>
                </CardContent>
            </Card>

            {receiving && (
                <div className="flex justify-end gap-3 mt-4">
                    <button
                        onClick={() => setReceiving(false)}
                        className="px-4 py-2 text-gray-400 hover:text-white transition-colors"
                    >
                        Cancel
                    </button>
                    <Button onClick={handleReceive} disabled={isSubmitting} isLoading={isSubmitting} className="shadow-glow">
                        <PackageCheck className="w-4 h-4 mr-2" />
                        Confirm Receipt
                    </Button>
                </div>
            )}
        </PageTransition>
    );
}
