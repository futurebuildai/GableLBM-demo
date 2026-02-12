import React, { useEffect, useState, useCallback } from 'react';
import { Dialog } from '@headlessui/react';
import { X, Truck, User, Calendar, FileText } from 'lucide-react';
import { Button } from '../ui/Button';
import { deliveryService } from '../../services/deliveryService';
import { useToast } from '../ui/ToastContext';
import type { Vehicle, Driver, CreateRouteRequest } from '../../types/delivery';

interface Props {
    isOpen: boolean;
    onClose: () => void;
    onCreated: () => void;
}

export const CreateRouteModal: React.FC<Props> = ({ isOpen, onClose, onCreated }) => {
    const { showToast } = useToast();
    const [vehicles, setVehicles] = useState<Vehicle[]>([]);
    const [drivers, setDrivers] = useState<Driver[]>([]);
    const [vehicleId, setVehicleId] = useState('');
    const [driverId, setDriverId] = useState('');
    const [scheduledDate, setScheduledDate] = useState(() => new Date().toISOString().split('T')[0]);
    const [notes, setNotes] = useState('');
    const [saving, setSaving] = useState(false);

    const loadFleet = useCallback(async () => {
        try {
            const [v, d] = await Promise.all([
                deliveryService.listVehicles(),
                deliveryService.listDrivers(),
            ]);
            setVehicles(v);
            setDrivers(d.filter(dr => dr.status === 'ACTIVE'));
        } catch {
            showToast('Failed to load fleet data', 'error');
        }
    }, [showToast]);

    useEffect(() => {
        if (isOpen) {
            loadFleet();
        }
    }, [isOpen, loadFleet]);

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        if (!vehicleId || !driverId || !scheduledDate) return;

        setSaving(true);
        try {
            const req: CreateRouteRequest = {
                vehicle_id: vehicleId,
                driver_id: driverId,
                scheduled_date: scheduledDate,
                notes: notes || undefined,
            };
            await deliveryService.createRoute(req);
            showToast('Route created successfully', 'success');
            onCreated();
            onClose();
            setVehicleId('');
            setDriverId('');
            setNotes('');
        } catch {
            showToast('Failed to create route', 'error');
        } finally {
            setSaving(false);
        }
    };

    const selectedVehicle = vehicles.find(v => v.id === vehicleId);

    return (
        <Dialog open={isOpen} onClose={onClose} className="relative z-50">
            <div className="fixed inset-0 bg-black/80 backdrop-blur-sm" aria-hidden="true" />
            <div className="fixed inset-0 flex items-center justify-center p-4">
                <Dialog.Panel className="w-full max-w-lg transform overflow-hidden rounded-2xl bg-slate-steel border border-white/10 p-6 shadow-xl">
                    <div className="flex items-center justify-between mb-6 border-b border-white/10 pb-4">
                        <Dialog.Title className="text-xl font-bold font-mono text-white flex items-center gap-2">
                            <Truck className="text-gable-green" /> Create Route
                        </Dialog.Title>
                        <button onClick={onClose} className="text-zinc-400 hover:text-white transition-colors">
                            <X size={24} />
                        </button>
                    </div>

                    <form onSubmit={handleSubmit} className="space-y-5">
                        {/* Vehicle */}
                        <div>
                            <label className="block text-sm font-medium text-zinc-400 mb-1.5 flex items-center gap-1.5">
                                <Truck className="w-3.5 h-3.5" /> Vehicle
                            </label>
                            <select
                                value={vehicleId}
                                onChange={e => setVehicleId(e.target.value)}
                                required
                                className="w-full bg-[#0A0B10] border border-white/10 rounded-lg px-3 py-2.5 text-white text-sm focus:border-gable-green/50 focus:outline-none"
                            >
                                <option value="">Select vehicle...</option>
                                {vehicles.map(v => (
                                    <option key={v.id} value={v.id}>
                                        {v.name} — {v.vehicle_type.replace('_', ' ')} ({v.license_plate})
                                        {v.capacity_weight_lbs ? ` • ${v.capacity_weight_lbs.toLocaleString()} lbs` : ''}
                                    </option>
                                ))}
                            </select>
                            {selectedVehicle?.capacity_weight_lbs && (
                                <p className="text-xs text-zinc-500 mt-1 font-mono">
                                    Capacity: {selectedVehicle.capacity_weight_lbs.toLocaleString()} lbs
                                </p>
                            )}
                        </div>

                        {/* Driver */}
                        <div>
                            <label className="block text-sm font-medium text-zinc-400 mb-1.5 flex items-center gap-1.5">
                                <User className="w-3.5 h-3.5" /> Driver
                            </label>
                            <select
                                value={driverId}
                                onChange={e => setDriverId(e.target.value)}
                                required
                                className="w-full bg-[#0A0B10] border border-white/10 rounded-lg px-3 py-2.5 text-white text-sm focus:border-gable-green/50 focus:outline-none"
                            >
                                <option value="">Select driver...</option>
                                {drivers.map(d => (
                                    <option key={d.id} value={d.id}>
                                        {d.name} {d.phone_number ? `(${d.phone_number})` : ''}
                                    </option>
                                ))}
                            </select>
                        </div>

                        {/* Date */}
                        <div>
                            <label className="block text-sm font-medium text-zinc-400 mb-1.5 flex items-center gap-1.5">
                                <Calendar className="w-3.5 h-3.5" /> Scheduled Date
                            </label>
                            <input
                                type="date"
                                value={scheduledDate}
                                onChange={e => setScheduledDate(e.target.value)}
                                required
                                className="w-full bg-[#0A0B10] border border-white/10 rounded-lg px-3 py-2.5 text-white text-sm focus:border-gable-green/50 focus:outline-none"
                            />
                        </div>

                        {/* Notes */}
                        <div>
                            <label className="block text-sm font-medium text-zinc-400 mb-1.5 flex items-center gap-1.5">
                                <FileText className="w-3.5 h-3.5" /> Notes (optional)
                            </label>
                            <textarea
                                value={notes}
                                onChange={e => setNotes(e.target.value)}
                                rows={2}
                                placeholder="Special instructions, area focus, etc."
                                className="w-full bg-[#0A0B10] border border-white/10 rounded-lg px-3 py-2.5 text-white text-sm focus:border-gable-green/50 focus:outline-none resize-none"
                            />
                        </div>

                        <div className="flex justify-end gap-3 pt-2">
                            <Button type="button" variant="outline" onClick={onClose}>Cancel</Button>
                            <Button type="submit" disabled={saving || !vehicleId || !driverId}>
                                {saving ? 'Creating...' : 'Create Route'}
                            </Button>
                        </div>
                    </form>
                </Dialog.Panel>
            </div>
        </Dialog>
    );
};
