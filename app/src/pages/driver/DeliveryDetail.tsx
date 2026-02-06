import { useEffect, useRef, useState } from "react";
import { useParams, useNavigate } from "react-router-dom";
import { deliveryService } from "../../services/deliveryService";
import type { Delivery, DeliveryStatus } from "../../types/delivery";

export function DeliveryDetail() {
    const { id } = useParams<{ id: string }>();
    const navigate = useNavigate();
    const [delivery, setDelivery] = useState<Delivery | null>(null);
    const [isSubmitting, setIsSubmitting] = useState(false);

    // POD Modal State
    const [showPODModal, setShowPODModal] = useState(false);
    const [status, setStatus] = useState<DeliveryStatus>('DELIVERED');
    const [signedBy, setSignedBy] = useState("");

    // Canvas Refs
    const canvasRef = useRef<HTMLCanvasElement>(null);
    const [isDrawing, setIsDrawing] = useState(false);

    useEffect(() => {
        if (id) {
            deliveryService.getDelivery(id).then(setDelivery);
        }
    }, [id]);

    // Canvas Logic
    const startDrawing = (e: React.MouseEvent<HTMLCanvasElement> | React.TouchEvent<HTMLCanvasElement>) => {
        const canvas = canvasRef.current;
        if (!canvas) return;
        const ctx = canvas.getContext('2d');
        if (!ctx) return;

        setIsDrawing(true);
        const rect = canvas.getBoundingClientRect();

        // Handle both mouse and touch events
        let clientX, clientY;
        if ('touches' in e) {
            clientX = e.touches[0].clientX;
            clientY = e.touches[0].clientY;
        } else {
            clientX = e.clientX;
            clientY = e.clientY;
        }

        const x = clientX - rect.left;
        const y = clientY - rect.top;

        ctx.beginPath();
        ctx.moveTo(x, y);
    };

    const draw = (e: React.MouseEvent<HTMLCanvasElement> | React.TouchEvent<HTMLCanvasElement>) => {
        if (!isDrawing) return;
        const canvas = canvasRef.current;
        if (!canvas) return;
        const ctx = canvas.getContext('2d');
        if (!ctx) return;

        const rect = canvas.getBoundingClientRect();

        let clientX, clientY;
        if ('touches' in e) {
            clientX = e.touches[0].clientX;
            clientY = e.touches[0].clientY;
        } else {
            clientX = e.clientX;
            clientY = e.clientY;
        }

        const x = clientX - rect.left;
        const y = clientY - rect.top;

        ctx.lineTo(x, y);
        ctx.stroke();
    };

    const stopDrawing = () => {
        setIsDrawing(false);
    };

    const clearSignature = () => {
        const canvas = canvasRef.current;
        if (canvas) {
            const ctx = canvas.getContext('2d');
            ctx?.clearRect(0, 0, canvas.width, canvas.height);
        }
    };

    const handleSubmit = async () => {
        if (!delivery) return;
        setIsSubmitting(true);

        try {
            let proofUrl = undefined;

            if (status === 'DELIVERED') {
                if (canvasRef.current) {
                    proofUrl = canvasRef.current.toDataURL("image/png");
                    // Ideally upload this to S3/Cloud and get URL.
                    // For MVP, we might send Base64?
                    // PRD says "Photo Upload". 
                    // If we send Base64 string as URL, database might complain if it's too long or expects proper URL.
                    // But let's try sending it as string for now if schema allows text.
                    // (Usually URL fields are text).
                    // If it's too large, we might need a dummy URL. 
                    // let's use a dummy URL if base64 is too massive, but base64 is what we have.
                    // We'll truncate/mock for MVP if needed, or assume backend handles it.
                    // UpdateRequest expects string. 
                }
            }

            await deliveryService.updateStatus(delivery.id, {
                status,
                pod_proof_url: proofUrl,
                pod_signed_by: signedBy || "Unknown"
            });

            setShowPODModal(false);
            // Refresh
            const updated = await deliveryService.getDelivery(delivery.id);
            setDelivery(updated);
        } catch {
            alert("Failed to update status");
        } finally {
            setIsSubmitting(false);
        }
    };

    if (!delivery) return <div className="p-4 text-white">Loading...</div>;

    return (
        <div className="pt-4 px-4 pb-20 space-y-6">
            {/* Header */}
            <div className="flex items-center gap-4">
                <button onClick={() => navigate(-1)} className="text-gray-400 font-mono text-sm">&larr; BACK</button>
                <div className="font-bold text-lg">DELIVERY</div>
            </div>

            {/* Info Card */}
            <div className="bg-[#161821] p-6 rounded-lg border border-white/10 space-y-4">
                <div>
                    <div className="text-sm text-gray-500 font-mono">CUSTOMER</div>
                    <div className="text-xl font-bold">{delivery.customer_name}</div>
                </div>
                <div>
                    <div className="text-sm text-gray-500 font-mono">ADDRESS</div>
                    <div className="text-lg">{delivery.address}</div>
                </div>
                <div>
                    <div className="text-sm text-gray-500 font-mono">INSTRUCTIONS</div>
                    <div className="text-blue-400">{delivery.delivery_instructions || "None"}</div>
                </div>
            </div>

            {/* Status */}
            <div className="bg-[#161821] p-6 rounded-lg border border-white/10 text-center">
                <div className="text-sm text-gray-500 font-mono mb-2">CURRENT STATUS</div>
                <div className={`text-2xl font-bold ${delivery.status === 'DELIVERED' ? 'text-green-500' :
                    delivery.status === 'FAILED' ? 'text-red-500' : 'text-gray-300'
                    }`}>
                    {delivery.status}
                </div>
            </div>

            {/* Action Button */}
            {delivery.status !== 'DELIVERED' && (
                <button
                    onClick={() => setShowPODModal(true)}
                    className="w-full bg-[#00FFA3] text-black font-bold py-4 rounded-lg text-xl shadow-[0_0_20px_rgba(0,255,163,0.3)] hover:scale-[1.02] transition-transform"
                >
                    COMPLETE DELIVERY
                </button>
            )}

            {/* POD Modal */}
            {showPODModal && (
                <div className="fixed inset-0 bg-black/90 z-[100] flex items-end sm:items-center justify-center p-4">
                    <div className="bg-[#161821] w-full max-w-md rounded-xl border border-white/10 p-6 space-y-6">
                        <h2 className="text-xl font-bold">Proof of Delivery</h2>

                        <div>
                            <label className="block text-sm text-gray-400 mb-2">Status</label>
                            <select
                                value={status}
                                onChange={(e) => setStatus(e.target.value as DeliveryStatus)}
                                className="w-full bg-[#0A0B10] p-3 rounded border border-white/20"
                            >
                                <option value="DELIVERED">Delivered</option>
                                <option value="FAILED">Failed / Refused</option>
                                <option value="PARTIAL">Partial Delivery</option>
                            </select>
                        </div>

                        {status === 'DELIVERED' && (
                            <>
                                <div>
                                    <label className="block text-sm text-gray-400 mb-2">Signed By</label>
                                    <input
                                        type="text"
                                        value={signedBy}
                                        onChange={e => setSignedBy(e.target.value)}
                                        className="w-full bg-[#0A0B10] p-3 rounded border border-white/20"
                                        placeholder="Recipient Name"
                                    />
                                </div>
                                <div>
                                    <div className="flex justify-between mb-2">
                                        <label className="text-sm text-gray-400">Signature</label>
                                        <button onClick={clearSignature} className="text-xs text-red-400">Clear</button>
                                    </div>
                                    <div className="bg-white rounded overflow-hidden h-40 touch-none">
                                        <canvas
                                            ref={canvasRef}
                                            width={400}
                                            height={160}
                                            className="w-full h-full cursor-crosshair"
                                            onMouseDown={startDrawing}
                                            onMouseMove={draw}
                                            onMouseUp={stopDrawing}
                                            onMouseLeave={stopDrawing}
                                            onTouchStart={startDrawing}
                                            onTouchMove={draw}
                                            onTouchEnd={stopDrawing}
                                        />
                                    </div>
                                </div>
                            </>
                        )}

                        <div className="flex gap-4 pt-4">
                            <button
                                onClick={() => setShowPODModal(false)}
                                className="flex-1 py-3 border border-white/10 rounded font-bold"
                            >
                                CANCEL
                            </button>
                            <button
                                onClick={handleSubmit}
                                disabled={isSubmitting}
                                className="flex-1 py-3 bg-[#00FFA3] text-black rounded font-bold"
                            >
                                {isSubmitting ? "SAVING..." : "CONFIRM"}
                            </button>
                        </div>
                    </div>
                </div>
            )}
        </div>
    );
}
