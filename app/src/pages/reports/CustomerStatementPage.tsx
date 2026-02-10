import { useState } from 'react';
import { ReportingService } from '../../services/ReportingService';
import type { CustomerStatement } from '../../types/invoice';
import { useToast } from '../../components/ui/ToastContext';
import { PageTransition } from '../../components/ui/PageTransition';
import { Card, CardContent } from '../../components/ui/Card';
import { Button } from '../../components/ui/Button';
import { FileText, Search } from 'lucide-react';

export function CustomerStatementPage() {
    const { showToast } = useToast();
    const [customerId, setCustomerId] = useState('');
    const [startDate, setStartDate] = useState('');
    const [endDate, setEndDate] = useState('');
    const [statement, setStatement] = useState<CustomerStatement | null>(null);
    const [loading, setLoading] = useState(false);

    const loadStatement = async () => {
        if (!customerId.trim()) {
            showToast('Enter a customer ID', 'error');
            return;
        }
        setLoading(true);
        try {
            const data = await ReportingService.getCustomerStatement(customerId, startDate || undefined, endDate || undefined);
            setStatement(data);
        } catch (err) {
            console.error(err);
            showToast('Failed to load statement', 'error');
        } finally {
            setLoading(false);
        }
    };

    const fmt = (v: number) => `$${v.toFixed(2)}`;

    return (
        <PageTransition>
            <div className="mb-8">
                <h1 className="text-3xl font-bold text-white flex items-center gap-3">
                    <FileText className="w-8 h-8 text-gable-green" />
                    Customer Statement
                </h1>
                <p className="text-zinc-500 mt-1">Generate account activity statements</p>
            </div>

            {/* Search Controls */}
            <Card variant="glass" className="mb-6 no-print">
                <CardContent className="p-6">
                    <div className="flex flex-wrap gap-4 items-end">
                        <div className="flex-1 min-w-[200px]">
                            <label className="text-xs text-zinc-500 uppercase tracking-wider block mb-1">Customer ID</label>
                            <input
                                type="text"
                                value={customerId}
                                onChange={(e) => setCustomerId(e.target.value)}
                                placeholder="UUID..."
                                className="w-full bg-black/20 border border-white/10 rounded px-3 py-2 text-white font-mono focus:border-[#00FFA3] outline-none"
                            />
                        </div>
                        <div>
                            <label className="text-xs text-zinc-500 uppercase tracking-wider block mb-1">Start Date</label>
                            <input
                                type="date"
                                value={startDate}
                                onChange={(e) => setStartDate(e.target.value)}
                                className="bg-black/20 border border-white/10 rounded px-3 py-2 text-white focus:border-[#00FFA3] outline-none"
                            />
                        </div>
                        <div>
                            <label className="text-xs text-zinc-500 uppercase tracking-wider block mb-1">End Date</label>
                            <input
                                type="date"
                                value={endDate}
                                onChange={(e) => setEndDate(e.target.value)}
                                className="bg-black/20 border border-white/10 rounded px-3 py-2 text-white focus:border-[#00FFA3] outline-none"
                            />
                        </div>
                        <Button onClick={loadStatement} disabled={loading} isLoading={loading}>
                            <Search className="w-4 h-4 mr-2" />
                            Generate
                        </Button>
                        {statement && (
                            <Button variant="outline" onClick={() => window.print()}>
                                <FileText className="w-4 h-4 mr-2" />
                                Print
                            </Button>
                        )}
                    </div>
                </CardContent>
            </Card>

            <style type="text/css" media="print">
                {`
                    @page { size: portrait; margin: 2cm; }
                    body { background: white !important; color: black !important; }
                    .no-print, header, aside, .omnibar { display: none !important; }
                    .print-only { display: block !important; }
                    * { border-color: #ddd !important; }
                    .text-white { color: black !important; }
                    .text-zinc-400, .text-zinc-500 { color: #666 !important; }
                    .bg-white\\/5 { background: transparent !important; }
                    .bg-black\\/20 { background: transparent !important;border: 1px solid #ccc; }
                `}
            </style>


            {/* Statement */}
            {
                statement && (
                    <Card variant="glass">
                        <CardContent className="p-6">
                            <div className="flex justify-between items-start mb-6">
                                <div>
                                    <h2 className="text-xl font-bold text-white">{statement.customer_name}</h2>
                                    <p className="text-sm text-zinc-400">
                                        Period: {statement.start_date} to {statement.end_date}
                                    </p>
                                </div>
                                <div className="text-right">
                                    <p className="text-xs text-zinc-500">Opening Balance</p>
                                    <p className="text-lg font-mono font-bold text-white">{fmt(statement.open_balance)}</p>
                                </div>
                            </div>

                            <table className="w-full text-sm text-left">
                                <thead className="bg-white/5 text-zinc-400 uppercase tracking-wider text-xs font-semibold">
                                    <tr>
                                        <th className="px-4 py-3">Date</th>
                                        <th className="px-4 py-3">Type</th>
                                        <th className="px-4 py-3">Description</th>
                                        <th className="px-4 py-3 text-right">Debit</th>
                                        <th className="px-4 py-3 text-right">Credit</th>
                                        <th className="px-4 py-3 text-right">Balance</th>
                                    </tr>
                                </thead>
                                <tbody className="divide-y divide-white/5">
                                    {(!statement.lines || statement.lines.length === 0) ? (
                                        <tr>
                                            <td colSpan={6} className="px-4 py-8 text-center text-zinc-500 italic">
                                                No transactions in this period
                                            </td>
                                        </tr>
                                    ) : (
                                        statement.lines.map((line, idx) => (
                                            <tr key={idx} className="hover:bg-white/5 transition-colors">
                                                <td className="px-4 py-3 font-mono text-zinc-300">{line.date}</td>
                                                <td className="px-4 py-3">
                                                    <span className={`px-2 py-0.5 rounded text-xs font-bold uppercase ${line.type === 'INVOICE' ? 'text-blue-400 bg-blue-500/10' :
                                                        line.type === 'PAYMENT' ? 'text-emerald-400 bg-emerald-500/10' :
                                                            line.type === 'REFUND' ? 'text-amber-400 bg-amber-500/10' :
                                                                'text-zinc-400 bg-zinc-500/10'
                                                        }`}>
                                                        {line.type}
                                                    </span>
                                                </td>
                                                <td className="px-4 py-3 text-zinc-300">{line.description}</td>
                                                <td className="px-4 py-3 text-right font-mono text-rose-400">
                                                    {line.debit > 0 ? fmt(line.debit) : ''}
                                                </td>
                                                <td className="px-4 py-3 text-right font-mono text-emerald-400">
                                                    {line.credit > 0 ? fmt(line.credit) : ''}
                                                </td>
                                                <td className="px-4 py-3 text-right font-mono text-white font-bold">
                                                    {fmt(line.balance)}
                                                </td>
                                            </tr>
                                        ))
                                    )}
                                </tbody>
                            </table>

                            <div className="flex justify-end mt-4 pt-4 border-t border-white/10">
                                <div className="text-right">
                                    <p className="text-xs text-zinc-500">Closing Balance</p>
                                    <p className="text-2xl font-mono font-bold text-white">{fmt(statement.close_balance)}</p>
                                </div>
                            </div>
                        </CardContent>
                    </Card>
                )
            }
        </PageTransition >
    );
}
