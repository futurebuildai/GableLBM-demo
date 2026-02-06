import { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';
import { InvoiceService } from '../../services/InvoiceService';
import { paymentService } from '../../services/paymentService';
import { PaymentModal } from '../../components/invoices/PaymentModal';
import type { Invoice } from '../../types/invoice';
import type { Payment, CreatePaymentRequest } from '../../types/payment';
import { Download, CreditCard, Mail } from 'lucide-react';

const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080';

export default function InvoiceDetail() {
    const { id } = useParams();
    const [invoice, setInvoice] = useState<Invoice | null>(null);
    const [payments, setPayments] = useState<Payment[]>([]);
    const [loading, setLoading] = useState(true);
    const [isPaymentModalOpen, setIsPaymentModalOpen] = useState(false);

    useEffect(() => {
        if (id) {
            loadInvoice(id);
            loadPayments(id);
        }
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

    async function loadPayments(id: string) {
        try {
            const data = await paymentService.getHistory(id);
            setPayments(data);
        } catch (error) {
            console.error("Failed to load payments", error);
        }
    }

    const handlePayment = async (input: CreatePaymentRequest) => {
        await paymentService.createPayment(input);
        if (id) {
            await loadInvoice(id);
            await loadPayments(id);
        }
    };

    if (loading || !invoice) return <div className="text-white">Loading invoice...</div>;

    const totalPaid = payments.reduce((sum, p) => sum + p.amount, 0);
    const amountDue = invoice.total_amount - totalPaid;

    return (
        <div className="space-y-8 max-w-4xl mx-auto pb-20">
            <div className="flex items-center justify-between pb-6 border-b border-white/10">
                <div>
                    <h1 className="text-3xl font-bold font-mono text-white">Invoice #{invoice.id.slice(0, 8)}</h1>
                    <p className="text-muted-foreground mt-1">Order Ref: <span className="font-mono text-zinc-400">{invoice.order_id.slice(0, 8)}</span></p>
                </div>
                <div className="flex gap-3">
                    <button
                        onClick={async () => {
                            if (!invoice.id) return;
                            try {
                                await InvoiceService.emailInvoice(invoice.id);
                                alert('Invoice emailed successfully!');
                            } catch {
                                alert('Failed to email invoice');
                            }
                        }}
                        className="bg-white/10 text-white hover:bg-white/20 px-4 py-2 rounded flex items-center gap-2 transition-colors border border-white/10"
                    >
                        <Mail size={18} /> Email
                    </button>
                    <button
                        onClick={() => window.open(`${API_URL}/documents/print/invoice/${invoice.id}`, '_blank')}
                        className="bg-white/10 text-white hover:bg-white/20 px-4 py-2 rounded flex items-center gap-2 transition-colors border border-white/10"
                    >
                        <Download size={18} /> Download
                    </button>

                    {invoice.status !== 'PAID' && (
                        <button
                            onClick={() => setIsPaymentModalOpen(true)}
                            className="bg-emerald-600 hover:bg-emerald-500 text-white px-4 py-2 rounded flex items-center gap-2 transition-colors font-medium shadow-lg shadow-emerald-900/20"
                        >
                            <CreditCard size={18} /> Pay
                        </button>
                    )}
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
                                ${invoice.status === 'PARTIAL' ? 'bg-blue-500/10 text-blue-500 border-blue-500/20' : ''}
                                ${invoice.status === 'PAID' ? 'bg-emerald-500/10 text-emerald-500 border-emerald-500/20' : ''}
                            `}>
                                {invoice.status}
                            </span>
                        </div>
                    </div>
                </div>
            </div>

            <div className="bg-zinc-900 rounded-lg border border-zinc-800 overflow-hidden">
                <div className="px-6 py-4 border-b border-zinc-800">
                    <h3 className="text-zinc-100 font-bold">Line Items</h3>
                </div>
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
                                <td className="px-6 py-4 text-white font-medium">{line.product_id.slice(0, 8)}</td>
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

            {/* Payment History Section */}
            {
                payments.length > 0 && (
                    <div className="bg-zinc-900 rounded-lg border border-zinc-800 overflow-hidden">
                        <div className="px-6 py-4 border-b border-zinc-800 flex justify-between items-center">
                            <h3 className="text-zinc-100 font-bold">Payment History</h3>
                            <span className="text-zinc-400 text-sm">Paid: <span className="text-green-400 font-mono">${totalPaid.toFixed(2)}</span></span>
                        </div>
                        <table className="w-full text-left text-sm">
                            <thead className="bg-zinc-950 text-zinc-400 uppercase text-xs">
                                <tr>
                                    <th className="px-6 py-4">Date</th>
                                    <th className="px-6 py-4">Method</th>
                                    <th className="px-6 py-4">Reference</th>
                                    <th className="px-6 py-4 text-right">Amount</th>
                                </tr>
                            </thead>
                            <tbody className="divide-y divide-zinc-800">
                                {payments.map(p => (
                                    <tr key={p.id}>
                                        <td className="px-6 py-4 text-zinc-300">{new Date(p.created_at).toLocaleString()}</td>
                                        <td className="px-6 py-4 text-zinc-300 font-bold">{p.method}</td>
                                        <td className="px-6 py-4 text-zinc-400 font-mono text-xs">{p.reference || '-'}</td>
                                        <td className="px-6 py-4 text-right text-white font-mono font-bold">${p.amount.toFixed(2)}</td>
                                    </tr>
                                ))}
                            </tbody>
                        </table>
                    </div>
                )
            }

            {
                invoice.id && (
                    <PaymentModal
                        isOpen={isPaymentModalOpen}
                        onClose={() => setIsPaymentModalOpen(false)}
                        onSave={handlePayment}
                        invoiceId={invoice.id}
                        amountDue={amountDue > 0 ? amountDue : 0}
                    />
                )
            }
        </div >
    );
}
