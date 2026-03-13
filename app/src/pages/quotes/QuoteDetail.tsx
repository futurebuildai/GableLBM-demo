import { useEffect, useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { FileText, Download, ArrowLeft, ShoppingCart, Send, Check, X, Sparkles, Eye, Map, Package, AlertTriangle, ShieldAlert } from 'lucide-react';
import { QuoteService } from '../../services/QuoteService';
import type { Quote, QuoteState, ParseMapItem } from '../../types/quote';
import { useToast } from '../../components/ui/ToastContext';
import { OrderService } from '../../services/OrderService';

type Tab = 'details' | 'original' | 'mapping';

export default function QuoteDetail() {
    const { id } = useParams();
    const navigate = useNavigate();
    const { showToast } = useToast();
    const [quote, setQuote] = useState<Quote | null>(null);
    const [loading, setLoading] = useState(true);
    const [processing, setProcessing] = useState(false);
    const [activeTab, setActiveTab] = useState<Tab>('details');

    useEffect(() => {
        if (id) loadQuote(id);
    }, [id]);

    async function loadQuote(quoteId: string) {
        try {
            const data = await QuoteService.getQuote(quoteId);
            setQuote(data);
        } catch (error) {
            console.error(error);
            showToast('Failed to load quote', 'error');
        } finally {
            setLoading(false);
        }
    }

    async function handleStateChange(state: QuoteState) {
        if (!quote) return;
        setProcessing(true);
        try {
            const updated = await QuoteService.updateQuoteState(quote.id, state);
            setQuote(updated);
            showToast(`Quote marked as ${state.toLowerCase()}`, 'success');
        } catch (error) {
            showToast(`Failed: ${error instanceof Error ? error.message : 'Unknown error'}`, 'error');
        } finally {
            setProcessing(false);
        }
    }

    async function handleConvert() {
        if (!quote) return;
        setProcessing(true);
        try {
            const orderPayload = await QuoteService.convertToOrder(quote.id);
            const order = await OrderService.createOrder(orderPayload);
            showToast('Quote converted to order', 'success');
            navigate(`/erp/orders/${order.id}`);
        } catch (error) {
            showToast(`Failed: ${error instanceof Error ? error.message : 'Unknown error'}`, 'error');
        } finally {
            setProcessing(false);
        }
    }

    if (loading || !quote) {
        return <div className="text-white p-8">Loading quote...</div>;
    }

    const stateColors: Record<string, string> = {
        DRAFT: 'bg-zinc-500/20 text-zinc-400 border-zinc-500/30',
        SENT: 'bg-blue-500/20 text-blue-400 border-blue-500/30',
        ACCEPTED: 'bg-emerald-500/20 text-emerald-400 border-emerald-500/30',
        REJECTED: 'bg-red-500/20 text-red-400 border-red-500/30',
        EXPIRED: 'bg-amber-500/20 text-amber-400 border-amber-500/30',
    };

    const tabs: { id: Tab; label: string; icon: typeof FileText; show: boolean }[] = [
        { id: 'details', label: 'Details', icon: FileText, show: true },
        { id: 'original', label: 'Original Upload', icon: Eye, show: quote.source === 'ai' && !!quote.original_filename },
        { id: 'mapping', label: 'AI Mapping', icon: Map, show: quote.source === 'ai' && !!quote.parse_map?.length },
    ];

    return (
        <div className="space-y-6 max-w-6xl mx-auto">
            {/* Header */}
            <div className="flex items-center justify-between pb-6 border-b border-white/10">
                <div>
                    <button onClick={() => navigate('/erp/quotes')} className="text-zinc-500 hover:text-white text-sm flex items-center gap-1 mb-3 transition-colors">
                        <ArrowLeft size={14} /> Back to Quotes
                    </button>
                    <div className="flex items-center gap-4 mb-2">
                        <h1 className="text-3xl font-bold font-mono text-white">Quote #{quote.id.slice(0, 8)}</h1>
                        <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium border ${stateColors[quote.state] || ''}`}>
                            {quote.state}
                        </span>
                        {quote.source === 'ai' && (
                            <span className="inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-[10px] font-bold uppercase tracking-wider bg-violet-500/15 text-violet-400 border border-violet-500/20">
                                <Sparkles className="w-3 h-3" /> AI Parsed
                            </span>
                        )}
                    </div>
                    <p className="text-muted-foreground">
                        {quote.customer_name || quote.customer_id.slice(0, 8)} &middot; Created {new Date(quote.created_at).toLocaleDateString()}
                    </p>
                </div>
                <div className="flex gap-2">
                    {quote.state === 'DRAFT' && (
                        <button onClick={() => handleStateChange('SENT')} disabled={processing}
                            className="bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-500 transition-colors flex items-center gap-2 text-sm font-medium disabled:opacity-50">
                            <Send size={14} /> Mark Sent
                        </button>
                    )}
                    {(quote.state === 'DRAFT' || quote.state === 'SENT') && (
                        <>
                            <button onClick={() => handleStateChange('ACCEPTED')} disabled={processing}
                                className="bg-emerald-600 text-white px-4 py-2 rounded hover:bg-emerald-500 transition-colors flex items-center gap-2 text-sm font-medium disabled:opacity-50">
                                <Check size={14} /> Accept
                            </button>
                            <button onClick={() => handleStateChange('REJECTED')} disabled={processing}
                                className="bg-red-600/80 text-white px-4 py-2 rounded hover:bg-red-500 transition-colors flex items-center gap-2 text-sm font-medium disabled:opacity-50">
                                <X size={14} /> Reject
                            </button>
                        </>
                    )}
                    {(quote.state === 'DRAFT' || quote.state === 'SENT' || quote.state === 'ACCEPTED') && (
                        <button onClick={handleConvert} disabled={processing}
                            className="bg-gable-green text-black px-4 py-2 rounded hover:bg-gable-green/90 transition-colors flex items-center gap-2 text-sm font-bold disabled:opacity-50">
                            <ShoppingCart size={14} /> Convert to Order
                        </button>
                    )}
                </div>
            </div>

            {/* Tabs */}
            <div className="flex gap-1 border-b border-white/10">
                {tabs.filter(t => t.show).map(tab => (
                    <button
                        key={tab.id}
                        onClick={() => setActiveTab(tab.id)}
                        className={`flex items-center gap-2 px-4 py-3 text-sm font-medium transition-colors border-b-2 -mb-px ${
                            activeTab === tab.id
                                ? 'text-gable-green border-gable-green'
                                : 'text-zinc-500 border-transparent hover:text-white hover:border-white/20'
                        }`}
                    >
                        <tab.icon size={16} />
                        {tab.label}
                    </button>
                ))}
            </div>

            {/* Tab Content */}
            {activeTab === 'details' && <DetailsTab quote={quote} />}
            {activeTab === 'original' && <OriginalUploadTab quote={quote} />}
            {activeTab === 'mapping' && <MappingTab parseMap={quote.parse_map || []} />}
        </div>
    );
}

// --- Details Tab ---
function DetailsTab({ quote }: { quote: Quote }) {
    return (
        <div className="space-y-6">
            {/* Summary Cards */}
            <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
                <SummaryCard label="Total Amount" value={`$${quote.total_amount.toFixed(2)}`} accent="text-gable-green" />
                <SummaryCard label="Lines" value={String(quote.lines?.length || 0)} />
                <SummaryCard label="Source" value={quote.source === 'ai' ? 'AI Parsed' : 'Manual'} />
                <SummaryCard label="State" value={quote.state} />
            </div>

            {/* Timeline */}
            <div className="bg-slate-steel border border-white/10 rounded-lg p-6">
                <h3 className="text-sm font-medium text-zinc-400 uppercase tracking-wider mb-4">Timeline</h3>
                <div className="space-y-3 text-sm">
                    <TimelineEntry label="Created" date={quote.created_at} />
                    {quote.sent_at && <TimelineEntry label="Sent" date={quote.sent_at} />}
                    {quote.accepted_at && <TimelineEntry label="Accepted" date={quote.accepted_at} />}
                    {quote.rejected_at && <TimelineEntry label="Rejected" date={quote.rejected_at} />}
                </div>
            </div>

            {/* Line Items */}
            <div className="bg-slate-steel border border-white/10 rounded-lg overflow-hidden">
                <table className="w-full text-left text-sm">
                    <thead>
                        <tr className="border-b border-white/10 bg-white/5">
                            <th className="p-4 font-medium text-muted-foreground">SKU</th>
                            <th className="p-4 font-medium text-muted-foreground">Description</th>
                            <th className="p-4 font-medium text-muted-foreground text-right">Qty</th>
                            <th className="p-4 font-medium text-muted-foreground text-right">Unit Price</th>
                            <th className="p-4 font-medium text-muted-foreground text-right">Total</th>
                        </tr>
                    </thead>
                    <tbody className="divide-y divide-white/5">
                        {(quote.lines || []).map(line => (
                            <tr key={line.id} className="hover:bg-white/5 transition-colors">
                                <td className="p-4 font-mono text-white">{line.sku}</td>
                                <td className="p-4 text-zinc-300">{line.description}</td>
                                <td className="p-4 text-right font-mono text-zinc-300">
                                    {line.quantity} <span className="text-zinc-600 text-xs">{line.uom}</span>
                                </td>
                                <td className="p-4 text-right font-mono text-zinc-300">${line.unit_price.toFixed(2)}</td>
                                <td className="p-4 text-right font-mono text-gable-green font-medium">${line.line_total.toFixed(2)}</td>
                            </tr>
                        ))}
                    </tbody>
                    {(quote.lines?.length || 0) > 0 && (
                        <tfoot className="bg-white/5 border-t border-white/10">
                            <tr>
                                <td colSpan={4} className="p-4 text-right font-medium text-zinc-400 uppercase tracking-wider text-xs">Total</td>
                                <td className="p-4 text-right font-mono text-xl font-bold text-gable-green">${quote.total_amount.toFixed(2)}</td>
                            </tr>
                        </tfoot>
                    )}
                </table>
            </div>
        </div>
    );
}

// --- Original Upload Tab ---
function OriginalUploadTab({ quote }: { quote: Quote }) {
    const fileUrl = QuoteService.getOriginalFileUrl(quote.id);
    const isImage = quote.original_content_type?.startsWith('image/');
    const isPdf = quote.original_content_type === 'application/pdf';

    return (
        <div className="space-y-4">
            <div className="flex items-center justify-between">
                <h3 className="text-sm font-medium text-zinc-400 uppercase tracking-wider">Original Material List Upload</h3>
                <a href={fileUrl} download={quote.original_filename || 'original-upload'}
                    className="text-gable-green hover:text-gable-green/80 text-sm flex items-center gap-2 transition-colors">
                    <Download size={14} /> Download Original
                </a>
            </div>

            <div className="bg-slate-steel border border-white/10 rounded-lg overflow-hidden">
                {isImage && (
                    <img src={fileUrl} alt="Original material list" className="w-full max-h-[70vh] object-contain bg-black/30 p-4" />
                )}
                {isPdf && (
                    <iframe src={fileUrl} className="w-full h-[70vh]" title="Original PDF" />
                )}
                {!isImage && !isPdf && (
                    <div className="p-12 text-center text-zinc-500">
                        <FileText className="w-12 h-12 mx-auto mb-4 opacity-50" />
                        <p>Preview not available for {quote.original_content_type}</p>
                        <a href={fileUrl} download className="text-gable-green hover:underline text-sm mt-2 inline-block">
                            Download to view
                        </a>
                    </div>
                )}
            </div>
        </div>
    );
}

// --- Mapping Tab ---
function MappingTab({ parseMap }: { parseMap: ParseMapItem[] }) {
    return (
        <div className="space-y-4">
            <div className="flex items-center justify-between">
                <h3 className="text-sm font-medium text-zinc-400 uppercase tracking-wider">AI Extraction Mapping</h3>
                <div className="text-xs text-zinc-500">
                    {parseMap.length} items extracted &middot;
                    {' '}{parseMap.filter(i => i.confidence >= 0.9 && !i.is_special_order).length} high confidence &middot;
                    {' '}{parseMap.filter(i => i.is_special_order).length} special order
                </div>
            </div>

            <div className="space-y-3">
                {parseMap.map((item, idx) => (
                    <div key={idx} className="bg-slate-steel border border-white/10 rounded-lg p-4">
                        {/* Raw text source */}
                        <div className="flex items-center gap-3 mb-3 pb-3 border-b border-white/5">
                            <span className="text-[10px] font-bold uppercase tracking-wider text-zinc-600 bg-white/5 px-2 py-0.5 rounded">
                                Raw Input
                            </span>
                            <span className="font-mono text-sm text-zinc-300">{item.raw_text}</span>
                            <ConfidenceBadge confidence={item.confidence} isSpecialOrder={item.is_special_order} />
                        </div>

                        {/* Mapping arrow */}
                        <div className="flex items-start gap-4">
                            <div className="flex-1">
                                {item.matched_product ? (
                                    <div className="flex items-center gap-3">
                                        <Package className="w-4 h-4 text-emerald-400 shrink-0" />
                                        <div>
                                            <div className="font-mono text-white text-sm">{item.matched_product.sku}</div>
                                            <div className="text-xs text-zinc-400">{item.matched_product.description}</div>
                                        </div>
                                        <div className="ml-auto text-right">
                                            <div className="font-mono text-sm text-white">
                                                {item.quantity} <span className="text-zinc-500 text-xs">{item.uom}</span>
                                            </div>
                                            <div className="font-mono text-xs text-emerald-400">${item.matched_product.base_price.toFixed(2)}/ea</div>
                                        </div>
                                    </div>
                                ) : (
                                    <div className="flex items-center gap-3">
                                        <ShieldAlert className="w-4 h-4 text-rose-400 shrink-0" />
                                        <div>
                                            <div className="text-sm text-rose-300">No catalog match — Special Order</div>
                                            <div className="text-xs text-zinc-500">{item.raw_text}</div>
                                        </div>
                                        <div className="ml-auto font-mono text-sm text-white">
                                            {item.quantity} <span className="text-zinc-500 text-xs">{item.uom}</span>
                                        </div>
                                    </div>
                                )}
                            </div>
                        </div>

                        {/* Alternatives considered */}
                        {item.alternatives && item.alternatives.length > 0 && (
                            <div className="mt-3 pt-3 border-t border-white/5">
                                <span className="text-[10px] font-bold uppercase tracking-wider text-zinc-600">
                                    Alternatives Considered
                                </span>
                                <div className="mt-2 space-y-1.5 pl-4 border-l-2 border-zinc-700/50">
                                    {item.alternatives.map((alt, aidx) => (
                                        <div key={aidx} className="flex items-center justify-between text-xs text-zinc-500">
                                            <span><span className="font-mono text-zinc-400">{alt.sku}</span> — {alt.description}</span>
                                            <span className="font-mono text-zinc-400">${alt.base_price.toFixed(2)}</span>
                                        </div>
                                    ))}
                                </div>
                            </div>
                        )}
                    </div>
                ))}
            </div>
        </div>
    );
}

// --- Shared Components ---
function SummaryCard({ label, value, accent }: { label: string; value: string; accent?: string }) {
    return (
        <div className="bg-slate-steel border border-white/10 rounded-lg p-4">
            <div className="text-xs text-zinc-500 uppercase tracking-wider mb-1">{label}</div>
            <div className={`text-lg font-mono font-bold ${accent || 'text-white'}`}>{value}</div>
        </div>
    );
}

function TimelineEntry({ label, date }: { label: string; date: string }) {
    return (
        <div className="flex items-center gap-3">
            <div className="w-2 h-2 rounded-full bg-gable-green/60" />
            <span className="text-zinc-400 w-20">{label}</span>
            <span className="text-white font-mono text-xs">{new Date(date).toLocaleString()}</span>
        </div>
    );
}

function ConfidenceBadge({ confidence, isSpecialOrder }: { confidence: number; isSpecialOrder: boolean }) {
    if (isSpecialOrder) {
        return (
            <span className="inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-[10px] font-bold uppercase tracking-wider bg-rose-500/15 text-rose-400 border border-rose-500/20">
                <AlertTriangle className="w-2.5 h-2.5" /> Special
            </span>
        );
    }
    const pct = Math.round(confidence * 100);
    if (confidence >= 0.9) {
        return (
            <span className="inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-[10px] font-bold uppercase tracking-wider bg-emerald-500/15 text-emerald-400 border border-emerald-500/20">
                <Check className="w-2.5 h-2.5" /> {pct}%
            </span>
        );
    }
    return (
        <span className="inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-[10px] font-bold uppercase tracking-wider bg-amber-500/15 text-amber-400 border border-amber-500/20">
            <AlertTriangle className="w-2.5 h-2.5" /> {pct}%
        </span>
    );
}
