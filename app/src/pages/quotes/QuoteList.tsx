import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { Plus, ArrowRight, ShoppingCart } from 'lucide-react';
import { QuoteService } from '../../services/QuoteService';
import { OrderService } from '../../services/OrderService';
import type { Quote } from '../../types/quote';
import { useToast } from '../../components/ui/ToastContext';

export default function QuoteList() {
    const navigate = useNavigate();
    const { showToast } = useToast();
    const [quotes, setQuotes] = useState<Quote[]>([]);
    const [loading, setLoading] = useState(true);
    const [converting, setConverting] = useState<string | null>(null);

    useEffect(() => {
        loadQuotes();
    }, []);

    async function loadQuotes() {
        try {
            const data = await QuoteService.listQuotes();
            setQuotes(data || []);
        } catch (error) {
            console.error('Failed to load quotes:', error);
        } finally {
            setLoading(false);
        }
    }

    async function handleConvert(quoteId: string) {
        setConverting(quoteId);
        try {
            const orderPayload = await QuoteService.convertToOrder(quoteId);
            const order = await OrderService.createOrder(orderPayload);
            showToast('Quote converted to order successfully', 'success');
            navigate(`/erp/orders/${order.id}`);
        } catch (error) {
            showToast(`Failed to convert: ${error instanceof Error ? error.message : 'Unknown error'}`, 'error');
        } finally {
            setConverting(null);
        }
    }

    const stateColors: Record<string, string> = {
        DRAFT: 'bg-zinc-500/20 text-zinc-400 border-zinc-500/30',
        SENT: 'bg-blue-500/20 text-blue-400 border-blue-500/30',
        ACCEPTED: 'bg-emerald-500/20 text-emerald-400 border-emerald-500/30',
        REJECTED: 'bg-red-500/20 text-red-400 border-red-500/30',
        EXPIRED: 'bg-amber-500/20 text-amber-400 border-amber-500/30',
    };

    return (
        <div className="space-y-6">
            <div className="flex items-center justify-between">
                <div>
                    <h1 className="text-3xl font-bold tracking-tight text-white font-mono">Quotes</h1>
                    <p className="text-muted-foreground mt-2">Manage sales quotes and convert to orders.</p>
                </div>
                <button
                    onClick={() => navigate('/quotes/new')}
                    className="bg-gable-green text-black font-bold px-4 py-2 rounded hover:bg-gable-green/90 transition-colors flex items-center gap-2"
                >
                    <Plus size={16} /> New Quote
                </button>
            </div>

            <div className="bg-slate-steel border border-white/10 rounded-lg overflow-hidden">
                <table className="w-full text-left text-sm">
                    <thead>
                        <tr className="border-b border-white/10 bg-white/5">
                            <th className="p-4 font-medium text-muted-foreground">Quote ID</th>
                            <th className="p-4 font-medium text-muted-foreground">Date</th>
                            <th className="p-4 font-medium text-muted-foreground">Customer</th>
                            <th className="p-4 font-medium text-muted-foreground">State</th>
                            <th className="p-4 font-medium text-muted-foreground text-right">Total</th>
                            <th className="p-4 font-medium text-muted-foreground text-right">Actions</th>
                        </tr>
                    </thead>
                    <tbody className="divide-y divide-white/5">
                        {loading ? (
                            <tr>
                                <td colSpan={6} className="p-8 text-center text-muted-foreground">Loading quotes...</td>
                            </tr>
                        ) : quotes.length === 0 ? (
                            <tr>
                                <td colSpan={6} className="p-8 text-center text-muted-foreground">
                                    No quotes found. Create your first quote to get started.
                                </td>
                            </tr>
                        ) : (
                            quotes.map((quote) => (
                                <tr key={quote.id} className="hover:bg-white/5 transition-colors">
                                    <td className="p-4 font-mono text-white/80">#{quote.id.slice(0, 8)}</td>
                                    <td className="p-4 text-white/80">{new Date(quote.created_at).toLocaleDateString()}</td>
                                    <td className="p-4 text-white font-medium">{quote.customer_name || quote.customer_id.slice(0, 8)}</td>
                                    <td className="p-4">
                                        <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium border ${stateColors[quote.state] || ''}`}>
                                            {quote.state}
                                        </span>
                                    </td>
                                    <td className="p-4 font-mono text-right text-gable-green">
                                        ${quote.total_amount.toFixed(2)}
                                    </td>
                                    <td className="p-4 text-right">
                                        <div className="flex items-center justify-end gap-2">
                                            {(quote.state === 'DRAFT' || quote.state === 'SENT' || quote.state === 'ACCEPTED') && (
                                                <button
                                                    onClick={() => handleConvert(quote.id)}
                                                    disabled={converting === quote.id}
                                                    className="text-gable-green hover:text-gable-green/80 transition-colors flex items-center gap-1 text-xs font-medium disabled:opacity-50"
                                                    title="Convert to Order"
                                                >
                                                    <ShoppingCart size={14} />
                                                    {converting === quote.id ? 'Converting...' : 'Convert'}
                                                </button>
                                            )}
                                            <button
                                                onClick={() => navigate(`/erp/quotes/${quote.id}`)}
                                                className="text-white/50 hover:text-white transition-colors"
                                            >
                                                <ArrowRight size={18} />
                                            </button>
                                        </div>
                                    </td>
                                </tr>
                            ))
                        )}
                    </tbody>
                </table>
            </div>
        </div>
    );
}
