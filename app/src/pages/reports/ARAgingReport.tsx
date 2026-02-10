import { useEffect, useState } from 'react';
import { ReportingService } from '../../services/ReportingService';
import type { ARAgingReport } from '../../types/invoice';
import { useToast } from '../../components/ui/ToastContext';
import { PageTransition } from '../../components/ui/PageTransition';
import { Card, CardContent } from '../../components/ui/Card';
import { DollarSign } from 'lucide-react';

export function ARAgingReportPage() {
    const { showToast } = useToast();
    const [report, setReport] = useState<ARAgingReport | null>(null);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        loadReport();
    }, []);

    const loadReport = async () => {
        try {
            const data = await ReportingService.getARAgingReport();
            setReport(data);
        } catch (err) {
            console.error(err);
            showToast('Failed to load AR aging report', 'error');
        } finally {
            setLoading(false);
        }
    };

    const fmt = (v: number) => `$${v.toFixed(2)}`;

    if (loading) {
        return (
            <div className="p-12 flex justify-center">
                <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-gable-green"></div>
            </div>
        );
    }

    return (
        <PageTransition>
            <div className="mb-8">
                <h1 className="text-3xl font-bold text-white flex items-center gap-3">
                    <DollarSign className="w-8 h-8 text-gable-green" />
                    AR Aging Report
                </h1>
                <p className="text-zinc-500 mt-1">
                    Accounts receivable aging as of {report?.as_of_date || 'today'}
                </p>
            </div>

            {/* Summary Cards */}
            {report && (
                <div className="grid grid-cols-2 md:grid-cols-5 gap-4 mb-6">
                    {[
                        { label: 'Current (0-30)', value: report.total_current, color: 'text-emerald-400' },
                        { label: '31-60 Days', value: report.total_31_60, color: 'text-amber-400' },
                        { label: '61-90 Days', value: report.total_61_90, color: 'text-orange-400' },
                        { label: '90+ Days', value: report.total_over_90, color: 'text-rose-400' },
                        { label: 'Grand Total', value: report.grand_total, color: 'text-white' },
                    ].map((item) => (
                        <Card key={item.label} variant="glass">
                            <CardContent className="p-4 text-center">
                                <p className="text-xs text-zinc-500 uppercase tracking-wider mb-1">{item.label}</p>
                                <p className={`text-xl font-mono font-bold ${item.color}`}>{fmt(item.value)}</p>
                            </CardContent>
                        </Card>
                    ))}
                </div>
            )}

            {/* Detail Table */}
            <Card variant="glass">
                <CardContent className="p-0">
                    {!report || report.buckets.length === 0 ? (
                        <div className="p-12 text-center text-zinc-500">No outstanding receivables</div>
                    ) : (
                        <table className="w-full text-sm text-left">
                            <thead className="bg-white/5 text-zinc-400 uppercase tracking-wider text-xs font-semibold">
                                <tr>
                                    <th className="px-6 py-4">Customer</th>
                                    <th className="px-6 py-4 text-right">Current</th>
                                    <th className="px-6 py-4 text-right">31-60</th>
                                    <th className="px-6 py-4 text-right">61-90</th>
                                    <th className="px-6 py-4 text-right">90+</th>
                                    <th className="px-6 py-4 text-right">Total</th>
                                </tr>
                            </thead>
                            <tbody className="divide-y divide-white/5">
                                {report.buckets.map((bucket) => (
                                    <tr key={bucket.customer_id} className="hover:bg-white/5 transition-colors">
                                        <td className="px-6 py-4 text-white font-medium">
                                            {bucket.customer_name}
                                            <span className="text-zinc-500 text-xs ml-2">({bucket.customer_id.slice(0, 8)})</span>
                                        </td>
                                        <td className="px-6 py-4 text-right font-mono text-emerald-400">{fmt(bucket.current)}</td>
                                        <td className="px-6 py-4 text-right font-mono text-amber-400">{fmt(bucket.days_31_60)}</td>
                                        <td className="px-6 py-4 text-right font-mono text-orange-400">{fmt(bucket.days_61_90)}</td>
                                        <td className="px-6 py-4 text-right font-mono text-rose-400">{fmt(bucket.over_90)}</td>
                                        <td className="px-6 py-4 text-right font-mono text-white font-bold">{fmt(bucket.total)}</td>
                                    </tr>
                                ))}
                            </tbody>
                            <tfoot className="bg-white/5 border-t border-white/10">
                                <tr className="font-bold">
                                    <td className="px-6 py-4 text-zinc-400 uppercase text-xs">Totals</td>
                                    <td className="px-6 py-4 text-right font-mono text-emerald-400">{fmt(report.total_current)}</td>
                                    <td className="px-6 py-4 text-right font-mono text-amber-400">{fmt(report.total_31_60)}</td>
                                    <td className="px-6 py-4 text-right font-mono text-orange-400">{fmt(report.total_61_90)}</td>
                                    <td className="px-6 py-4 text-right font-mono text-rose-400">{fmt(report.total_over_90)}</td>
                                    <td className="px-6 py-4 text-right font-mono text-white">{fmt(report.grand_total)}</td>
                                </tr>
                            </tfoot>
                        </table>
                    )}
                </CardContent>
            </Card>
        </PageTransition>
    );
}
