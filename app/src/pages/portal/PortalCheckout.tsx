import { useState, useEffect, useCallback } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { ArrowLeft, MapPin, Truck, Store, CreditCard, Building2, CheckCircle } from 'lucide-react';
import { Card, CardContent } from '../../components/ui/Card';
import { PortalService } from '../../services/PortalService';
import type { Cart, CheckoutRequest } from '../../types/portal';

const formatCurrency = (val: number): string =>
    new Intl.NumberFormat('en-US', { style: 'currency', currency: 'USD' }).format(val);

export const PortalCheckout = () => {
    const navigate = useNavigate();
    const [cart, setCart] = useState<Cart | null>(null);
    const [loading, setLoading] = useState(true);
    const [submitting, setSubmitting] = useState(false);
    const [error, setError] = useState('');
    const [success, setSuccess] = useState<string | null>(null);

    // Form state
    const [deliveryMethod, setDeliveryMethod] = useState<'DELIVERY' | 'PICKUP'>('DELIVERY');
    const [deliveryAddress, setDeliveryAddress] = useState('');
    const [paymentMethod, setPaymentMethod] = useState<'ACCOUNT' | 'CARD'>('ACCOUNT');
    const [notes, setNotes] = useState('');

    const fetchCart = useCallback(() => {
        setLoading(true);
        PortalService.getCart()
            .then(setCart)
            .catch(err => setError(err instanceof Error ? err.message : 'Failed to load cart'))
            .finally(() => setLoading(false));
    }, []);

    useEffect(() => { fetchCart(); }, [fetchCart]);

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        if (!cart || cart.items.length === 0) return;

        if (deliveryMethod === 'DELIVERY' && !deliveryAddress.trim()) {
            setError('Delivery address is required');
            return;
        }

        setSubmitting(true);
        setError('');

        const req: CheckoutRequest = {
            delivery_method: deliveryMethod,
            delivery_address: deliveryAddress,
            payment_method: paymentMethod,
            notes,
        };

        try {
            const resp = await PortalService.checkout(req);
            setSuccess(resp.order_id);
        } catch (err) {
            setError(err instanceof Error ? err.message : 'Checkout failed');
        } finally {
            setSubmitting(false);
        }
    };

    if (loading) {
        return (
            <div className="space-y-4 max-w-3xl mx-auto">
                <div className="h-10 w-48 bg-white/5 rounded-lg animate-pulse" />
                <div className="h-64 bg-white/5 rounded-2xl animate-pulse" />
                <div className="h-48 bg-white/5 rounded-2xl animate-pulse" />
            </div>
        );
    }

    if (success) {
        return (
            <div className="flex flex-col items-center justify-center py-20 text-center max-w-lg mx-auto">
                <div className="w-20 h-20 rounded-full bg-emerald-500/10 flex items-center justify-center mb-6">
                    <CheckCircle className="w-10 h-10 text-emerald-400" />
                </div>
                <h2 className="text-2xl font-bold text-white mb-2">Order Placed!</h2>
                <p className="text-zinc-400 mb-2">Your order has been submitted successfully.</p>
                <p className="text-sm font-mono text-zinc-500 mb-8">
                    Order ID: {success.substring(0, 8).toUpperCase()}
                </p>
                <div className="flex gap-3">
                    <Link
                        to="/portal/orders"
                        className="px-6 py-3 rounded-xl bg-gable-green text-black font-semibold hover:bg-gable-green/90 transition-colors"
                    >
                        View Orders
                    </Link>
                    <Link
                        to="/portal/catalog"
                        className="px-6 py-3 rounded-xl bg-white/5 border border-white/10 text-white font-semibold hover:bg-white/10 transition-colors"
                    >
                        Continue Shopping
                    </Link>
                </div>
            </div>
        );
    }

    if (!cart || cart.items.length === 0) {
        return (
            <div className="text-center py-16">
                <p className="text-zinc-400">Your cart is empty.</p>
                <Link to="/portal/catalog" className="text-gable-green hover:underline mt-2 inline-block">
                    Browse Catalog
                </Link>
            </div>
        );
    }

    return (
        <div className="max-w-3xl mx-auto">
            {/* Header */}
            <Link
                to="/portal/cart"
                className="inline-flex items-center gap-1.5 text-sm text-zinc-400 hover:text-white transition-colors mb-6"
            >
                <ArrowLeft className="w-4 h-4" /> Back to Cart
            </Link>

            <h1 className="text-display-large text-white mb-8">Checkout</h1>

            <form onSubmit={handleSubmit} className="space-y-6">
                {/* Order Summary */}
                <Card variant="glass">
                    <CardContent className="p-5">
                        <h3 className="text-sm font-semibold text-white uppercase tracking-wider mb-3">Order Summary</h3>
                        <div className="space-y-2 text-sm">
                            {cart.items.map(item => (
                                <div key={item.id} className="flex justify-between text-zinc-400">
                                    <span className="truncate mr-4">
                                        {item.product_name} x{item.quantity}
                                    </span>
                                    <span className="font-mono text-white shrink-0">{formatCurrency(item.line_total)}</span>
                                </div>
                            ))}
                        </div>
                        <div className="border-t border-white/10 mt-3 pt-3 flex justify-between">
                            <span className="font-medium text-white">Total</span>
                            <span className="text-xl font-bold text-white font-mono">{formatCurrency(cart.subtotal)}</span>
                        </div>
                    </CardContent>
                </Card>

                {/* Delivery Method */}
                <Card variant="glass">
                    <CardContent className="p-5">
                        <h3 className="text-sm font-semibold text-white uppercase tracking-wider mb-4">Delivery Method</h3>
                        <div className="grid grid-cols-2 gap-3">
                            <button
                                type="button"
                                onClick={() => setDeliveryMethod('DELIVERY')}
                                className={`flex flex-col items-center gap-2 p-4 rounded-xl border transition-all ${deliveryMethod === 'DELIVERY'
                                        ? 'border-gable-green bg-gable-green/5 text-gable-green'
                                        : 'border-white/10 text-zinc-400 hover:border-white/20'
                                    }`}
                            >
                                <Truck className="w-6 h-6" />
                                <span className="text-sm font-semibold">Delivery</span>
                            </button>
                            <button
                                type="button"
                                onClick={() => setDeliveryMethod('PICKUP')}
                                className={`flex flex-col items-center gap-2 p-4 rounded-xl border transition-all ${deliveryMethod === 'PICKUP'
                                        ? 'border-gable-green bg-gable-green/5 text-gable-green'
                                        : 'border-white/10 text-zinc-400 hover:border-white/20'
                                    }`}
                            >
                                <Store className="w-6 h-6" />
                                <span className="text-sm font-semibold">Will Call</span>
                            </button>
                        </div>

                        {deliveryMethod === 'DELIVERY' && (
                            <div className="mt-4">
                                <label className="flex items-center gap-1.5 text-sm text-zinc-400 mb-2">
                                    <MapPin className="w-4 h-4" /> Delivery Address
                                </label>
                                <textarea
                                    value={deliveryAddress}
                                    onChange={e => setDeliveryAddress(e.target.value)}
                                    placeholder="Enter delivery address..."
                                    rows={3}
                                    className="w-full px-4 py-3 rounded-xl bg-white/5 border border-white/10 text-white placeholder-zinc-500 focus:outline-none focus:border-gable-green/50 transition-colors text-sm resize-none"
                                />
                            </div>
                        )}
                    </CardContent>
                </Card>

                {/* Payment Method */}
                <Card variant="glass">
                    <CardContent className="p-5">
                        <h3 className="text-sm font-semibold text-white uppercase tracking-wider mb-4">Payment Method</h3>
                        <div className="grid grid-cols-2 gap-3">
                            <button
                                type="button"
                                onClick={() => setPaymentMethod('ACCOUNT')}
                                className={`flex flex-col items-center gap-2 p-4 rounded-xl border transition-all ${paymentMethod === 'ACCOUNT'
                                        ? 'border-gable-green bg-gable-green/5 text-gable-green'
                                        : 'border-white/10 text-zinc-400 hover:border-white/20'
                                    }`}
                            >
                                <Building2 className="w-6 h-6" />
                                <span className="text-sm font-semibold">Charge to Account</span>
                            </button>
                            <button
                                type="button"
                                onClick={() => setPaymentMethod('CARD')}
                                className={`flex flex-col items-center gap-2 p-4 rounded-xl border transition-all ${paymentMethod === 'CARD'
                                        ? 'border-gable-green bg-gable-green/5 text-gable-green'
                                        : 'border-white/10 text-zinc-400 hover:border-white/20'
                                    }`}
                            >
                                <CreditCard className="w-6 h-6" />
                                <span className="text-sm font-semibold">Pay by Card</span>
                            </button>
                        </div>
                    </CardContent>
                </Card>

                {/* Order Notes */}
                <Card variant="glass">
                    <CardContent className="p-5">
                        <h3 className="text-sm font-semibold text-white uppercase tracking-wider mb-3">Order Notes</h3>
                        <textarea
                            value={notes}
                            onChange={e => setNotes(e.target.value)}
                            placeholder="Special instructions, PO number, job reference..."
                            rows={3}
                            className="w-full px-4 py-3 rounded-xl bg-white/5 border border-white/10 text-white placeholder-zinc-500 focus:outline-none focus:border-gable-green/50 transition-colors text-sm resize-none"
                        />
                    </CardContent>
                </Card>

                {/* Error */}
                {error && (
                    <div className="px-4 py-3 rounded-xl bg-red-500/10 border border-red-500/20 text-red-400 text-sm">
                        {error}
                    </div>
                )}

                {/* Submit */}
                <button
                    type="submit"
                    disabled={submitting}
                    className="w-full py-4 rounded-xl text-lg font-bold bg-gable-green text-black hover:bg-gable-green/90 transition-all disabled:opacity-50 disabled:cursor-not-allowed active:scale-[0.99]"
                >
                    {submitting ? 'Placing Order...' : `Place Order - ${formatCurrency(cart.subtotal)}`}
                </button>
            </form>
        </div>
    );
};
