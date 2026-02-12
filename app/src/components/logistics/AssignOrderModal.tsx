import React, { useEffect, useState, useCallback } from 'react';
import { Dialog } from '@headlessui/react';
import { X, Package, MapPin, FileText, AlertTriangle } from 'lucide-react';
import { Button } from '../ui/Button';
import { deliveryService } from '../../services/deliveryService';
import { OrderService } from '../../services/OrderService';
import { useToast } from '../ui/ToastContext';
import type { Order } from '../../types/order';
import type { Vehicle, Delivery } from '../../types/delivery';

interface Props {
    isOpen: boolean;
    onClose: () => void;
    routeId: string;
    vehicleId: string;
    existingDeliveries: Delivery[];
    onAssigned: () => void;
}

export const AssignOrderModal: React.FC<Props> = ({ isOpen, onClose, routeId, vehicleId, existingDeliveries, onAssigned }) => {
    const { showToast } = useToast();
    const [orders, setOrders] = useState<Order[]>([]);
    const [vehicle, setVehicle] = useState<Vehicle | null>(null);
    const [selectedOrderId, setSelectedOrderId] = useState('');
    const [instructions, setInstructions] = useState('');
    const [saving, setSaving] = useState(false);
    const [loading, setLoading] = useState(true);

    const loadData = useCallback(async () => {
        setLoading(true);
        try {
            const [allOrders, vehicles] = await Promise.all([
                OrderService.listOrders(),
                deliveryService.listVehicles(),
            ]);
            // Filter to CONFIRMED orders not already assigned to this route
            const assignedOrderIds = new Set(existingDeliveries.map(d => d.order_id));
            setOrders(allOrders.filter(o => o.status === 'CONFIRMED' && !assignedOrderIds.has(o.id)));
            setVehicle(vehicles.find(v => v.id === vehicleId) || null);
        } catch {
            showToast('Failed to load orders', 'error');
        } finally {
            setLoading(false);
        }
    }, [existingDeliveries, vehicleId, showToast]);

    useEffect(() => {
        if (isOpen) {
            loadData();
        }
    }, [isOpen, loadData]);

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        if (!selectedOrderId) return;

        setSaving(true);
        try {
            const nextSequence = existingDeliveries.length + 1;
            const result = await deliveryService.assignOrder({
                route_id: routeId,
                order_id: selectedOrderId,
                stop_sequence: nextSequence,
                delivery_instructions: instructions || undefined,
            });
            if (result.capacity_warning) {
                const w = result.capacity_warning;
                showToast(
                    `Warning: Vehicle capacity ${w.vehicle_capacity_lbs.toLocaleString()} lbs exceeded — load after: ${w.total_after_lbs.toLocaleString()} lbs`,
                    'error'
                );
            } else {
                showToast('Order assigned to route', 'success');
            }
            onAssigned();
            onClose();
            setSelectedOrderId('');
            setInstructions('');
        } catch {
            showToast('Failed to assign order', 'error');
        } finally {
            setSaving(false);
        }
    };

    const capacityWarning = vehicle?.capacity_weight_lbs
        ? existingDeliveries.length >= 20
            ? 'Route has many stops — verify vehicle can handle the load.'
            : null
        : null;

    return (
        <Dialog open={isOpen} onClose={onClose} className="relative z-50">
            <div className="fixed inset-0 bg-black/80 backdrop-blur-sm" aria-hidden="true" />
            <div className="fixed inset-0 flex items-center justify-center p-4">
                <Dialog.Panel className="w-full max-w-lg transform overflow-hidden rounded-2xl bg-slate-steel border border-white/10 p-6 shadow-xl">
                    <div className="flex items-center justify-between mb-6 border-b border-white/10 pb-4">
                        <Dialog.Title className="text-xl font-bold font-mono text-white flex items-center gap-2">
                            <Package className="text-gable-green" /> Assign Order to Route
                        </Dialog.Title>
                        <button onClick={onClose} className="text-zinc-400 hover:text-white transition-colors">
                            <X size={24} />
                        </button>
                    </div>

                    {vehicle?.capacity_weight_lbs && (
                        <div className="mb-4 p-3 bg-sky-500/10 border border-sky-500/20 rounded-lg text-sm text-sky-300 font-mono">
                            Vehicle capacity: {vehicle.capacity_weight_lbs.toLocaleString()} lbs
                            <span className="text-zinc-500 ml-2">• {existingDeliveries.length} stops assigned</span>
                        </div>
                    )}

                    {capacityWarning && (
                        <div className="mb-4 p-3 bg-amber-500/10 border border-amber-500/20 rounded-lg text-sm text-amber-400 flex items-center gap-2">
                            <AlertTriangle className="w-4 h-4 shrink-0" />
                            {capacityWarning}
                        </div>
                    )}

                    {loading ? (
                        <div className="py-12 text-center text-zinc-500 animate-pulse">Loading orders...</div>
                    ) : (
                        <form onSubmit={handleSubmit} className="space-y-5">
                            {/* Order Selection */}
                            <div>
                                <label className="block text-sm font-medium text-zinc-400 mb-1.5 flex items-center gap-1.5">
                                    <MapPin className="w-3.5 h-3.5" /> Select Order
                                </label>
                                {orders.length === 0 ? (
                                    <p className="text-zinc-500 text-sm py-4 text-center">No confirmed orders available for assignment.</p>
                                ) : (
                                    <div className="max-h-60 overflow-y-auto space-y-2">
                                        {orders.map(order => (
                                            <label
                                                key={order.id}
                                                className={`flex items-center gap-3 p-3 rounded-lg border cursor-pointer transition-all ${selectedOrderId === order.id
                                                        ? 'bg-gable-green/10 border-gable-green/50'
                                                        : 'bg-[#0A0B10] border-white/10 hover:border-white/20'
                                                    }`}
                                            >
                                                <input
                                                    type="radio"
                                                    name="order"
                                                    value={order.id}
                                                    checked={selectedOrderId === order.id}
                                                    onChange={() => setSelectedOrderId(order.id)}
                                                    className="accent-emerald-500"
                                                />
                                                <div className="flex-1 min-w-0">
                                                    <div className="flex justify-between items-center">
                                                        <span className="text-white font-mono text-sm">{order.customer_name}</span>
                                                        <span className="text-xs text-zinc-500 font-mono">${(order.total_amount / 100).toFixed(2)}</span>
                                                    </div>
                                                    <span className="text-xs text-zinc-500">Order #{order.id.slice(0, 8)}</span>
                                                </div>
                                            </label>
                                        ))}
                                    </div>
                                )}
                            </div>

                            {/* Instructions */}
                            <div>
                                <label className="block text-sm font-medium text-zinc-400 mb-1.5 flex items-center gap-1.5">
                                    <FileText className="w-3.5 h-3.5" /> Delivery Instructions (optional)
                                </label>
                                <textarea
                                    value={instructions}
                                    onChange={e => setInstructions(e.target.value)}
                                    rows={2}
                                    placeholder="Gate code, dock preference, contact on arrival..."
                                    className="w-full bg-[#0A0B10] border border-white/10 rounded-lg px-3 py-2.5 text-white text-sm focus:border-gable-green/50 focus:outline-none resize-none"
                                />
                            </div>

                            <div className="flex justify-end gap-3 pt-2">
                                <Button type="button" variant="outline" onClick={onClose}>Cancel</Button>
                                <Button type="submit" disabled={saving || !selectedOrderId || orders.length === 0}>
                                    {saving ? 'Assigning...' : `Assign as Stop #${existingDeliveries.length + 1}`}
                                </Button>
                            </div>
                        </form>
                    )}
                </Dialog.Panel>
            </div>
        </Dialog>
    );
};
