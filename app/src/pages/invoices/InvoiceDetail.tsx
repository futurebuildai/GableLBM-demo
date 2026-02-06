import { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';
import { InvoiceService } from '../../services/InvoiceService';
import type { Invoice } from '../../types/invoice';
import { Download } from 'lucide-react';

const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080';

export default function InvoiceDetail() {
    const { id } = useParams();
    const [invoice, setInvoice] = useState<Invoice | null>(null);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        if (id) loadInvoice(id);
    }, [id]);

    async function loadInvoice(id: string) {
        try {
            const data = await InvoiceService.getInvoice(id);
            setInvoice(data);
        } catch (error) {
            console.error(error);
        } finally {
            setLoading(false);
        }
    }

    if (loading || !invoice) return <div className="text-white">Loading invoice...</div>;

    return (
        <div className="space-y-8 max-w-4xl mx-auto">
            <div className="flex items-center justify-between pb-6 border-b border-white/10">
                <div>
                    <h1 className="text-3xl font-bold font-mono text-white">Invoice #{invoice.id.slice(0, 8)}</h1>
                    <p className="text-muted-foreground mt-1">Order Ref: <span className="font-mono text-zinc-400">{invoice.order_id.slice(0, 8)}</span></p>
                </div>
                <div className="flex gap-3">
                    <button
                        onClick={() => window.open(`${API_URL}/documents/print/invoice/${invoice.id}`, '_blank')}
                        className="bg-white/10 text-white hover:bg-white/20 px-4 py-2 rounded flex items-center gap-2 transition-colors"
                    >
                        <Download size={18} /> Download PDF
                    </button>
                    {/* Placeholder for Pay Now */}
                </div>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
                <div className="bg-zinc-900 p-6 rounded-lg border border-zinc-800">
                    <h3 className="text-zinc-500 uppercase text-xs font-bold mb-4">Bill To</h3>
                    <div className="text-zinc-300">
                        <p className="font-mono text-white mb-2">{invoice.customer_id.slice(0, 8)}</p>
                        <p>123 Construction Way</p>
                        <p>Builder Town, ST 12345</p>
                    </div>
                </div>
                <div className="bg-zinc-900 p-6 rounded-lg border border-zinc-800 text-right">
                    <h3 className="text-zinc-500 uppercase text-xs font-bold mb-4">Invoice Details</h3>
                    <div className="space-y-2">
                        <div className="flex justify-between">
                            <span className="text-zinc-400">Issue Date</span>
                            <span className="text-zinc-200">{new Date(invoice.created_at).toLocaleDateString()}</span>
                        </div>
                        <div className="flex justify-between">
                            <span className="text-zinc-400">Due Date</span>
                            <span className="text-zinc-200">{invoice.due_date ? new Date(invoice.due_date).toLocaleDateString() : 'Net 30'}</span>
                        </div>
                        <div className="flex justify-between items-center mt-4 pt-4 border-t border-zinc-800">
                            <span className="text-zinc-400">Status</span>
                            <span className={`inline-flex items-center px-2 py-0.5 rounded text-xs font-medium border
                                ${invoice.status === 'UNPAID' ? 'bg-amber-500/10 text-amber-500 border-amber-500/20' : ''}
                                ${invoice.status === 'PAID' ? 'bg-emerald-500/10 text-emerald-500 border-emerald-500/20' : ''}
                            `}>
                                {invoice.status}
                            </span>
                        </div>
                    </div>
                </div>
            </div>

            <div className="bg-zinc-900 rounded-lg border border-zinc-800 overflow-hidden">
                <table className="w-full text-left text-sm">
                    <thead className="bg-zinc-950 text-zinc-400 uppercase text-xs">
                        <tr>
                            <th className="px-6 py-4">Item</th>
                            <th className="px-6 py-4 text-right">Qty</th>
                            <th className="px-6 py-4 text-right">Rate</th>
                            <th className="px-6 py-4 text-right">Amount</th>
                        </tr>
                    </thead>
                    <tbody className="divide-y divide-zinc-800">
                        {invoice.lines?.map(line => (
                            <tr key={line.id}>
                                <td className="px-6 py-4 text-white font-medium">{line.product_id.slice(0, 8)}</td>{/* Resolve Name */}
                                <td className="px-6 py-4 text-right text-zinc-300 font-mono">{line.quantity}</td>
                                <td className="px-6 py-4 text-right text-zinc-300 font-mono">${line.price_each.toFixed(2)}</td>
                                <td className="px-6 py-4 text-right text-white font-mono font-bold">${(line.quantity * line.price_each).toFixed(2)}</td>
                            </tr>
                        ))}
                    </tbody>
                    <tfoot className="bg-zinc-950">
                        <tr>
                            <td colSpan={3} className="px-6 py-4 text-right text-zinc-400 font-bold uppercase">Total Due</td>
                            <td className="px-6 py-4 text-right text-emerald-500 font-bold font-mono text-xl">${invoice.total_amount.toFixed(2)}</td>
                        </tr>
                    </tfoot>
                </table>
            </div>

        </div>
    );
}
