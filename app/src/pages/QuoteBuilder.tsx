import { useState, useEffect, useCallback } from 'react';
import { CustomerSelect } from '../components/customers/CustomerSelect';
import { LineItemEditor } from '../components/quotes/LineItemEditor';
import { EscalatorToggle } from '../components/quotes/EscalatorToggle';
import { MaterialListUpload } from '../components/quotes/MaterialListUpload';
import { ParsedResultsPanel } from '../components/quotes/ParsedResultsPanel';
import { QuoteService } from '../services/QuoteService';
import { ProductService } from '../services/product.service';
import type { Customer } from '../types/customer';
import type { Product } from '../types/product';
import type { CreateQuoteRequest } from '../types/quote';
import type { QuoteLineEscalator } from '../types/pricing';
import type { ParseResponse, ParsedItem } from '../types/parsing';
import { Save, FileText, Calculator, CreditCard, AlertCircle, TrendingUp } from 'lucide-react';
import { useNavigate } from 'react-router-dom';
import { PageTransition } from '../components/ui/PageTransition';
import { Card, CardContent } from '../components/ui/Card';
import { Button } from '../components/ui/Button';
import { useToast } from '../components/ui/ToastContext';

interface LineWithEscalator {
    product_id: string;
    sku: string;
    description: string;
    quantity: number;
    uom: string;
    unit_price: number;
    escalator: QuoteLineEscalator;
}

const defaultEscalator = (): QuoteLineEscalator => ({
    enabled: false,
    escalation_type: 'PERCENTAGE',
    escalation_rate: 5,
    effective_date: new Date().toISOString().split('T')[0],
    target_date: new Date(Date.now() + 90 * 24 * 60 * 60 * 1000).toISOString().split('T')[0],
});

export const QuoteBuilder = () => {
    const navigate = useNavigate();
    const { showToast } = useToast();
    const [customer, setCustomer] = useState<Customer | null>(null);
    const [products, setProducts] = useState<Product[]>([]);
    const [lines, setLines] = useState<LineWithEscalator[]>([]);
    const [loading, setLoading] = useState(false);

    // AI Parsing state
    const [parseResult, setParseResult] = useState<ParseResponse | null>(null);
    const [showParsePanel, setShowParsePanel] = useState(false);
    const [aiSource, setAiSource] = useState(false);
    const [lastParseResult, setLastParseResult] = useState<ParseResponse | null>(null);

    useEffect(() => {
        const loadProducts = async () => {
            try {
                const data = await ProductService.getProducts();
                setProducts(data);
            } catch (err) {
                console.error("Failed to load products", err);
            }
        };
        loadProducts();
    }, []);

    const handleAddLine = (product: Product, quantity: number, unitPrice: number) => {
        setLines([...lines, {
            product_id: product.id,
            sku: product.sku,
            description: product.description,
            uom: product.uom_primary,
            quantity,
            unit_price: unitPrice,
            escalator: defaultEscalator(),
        }]);
    };

    const handleEscalatorChange = (idx: number, escalator: QuoteLineEscalator) => {
        const updated = [...lines];
        updated[idx] = { ...updated[idx], escalator };
        setLines(updated);
    };

    // --- AI Parsing Handlers ---
    const handleParseComplete = useCallback((result: ParseResponse) => {
        setParseResult(result);
        setShowParsePanel(true);
    }, []);

    const handleAcceptParsed = useCallback((parsedItems: ParsedItem[]) => {
        const newLines: LineWithEscalator[] = parsedItems.map(item => ({
            product_id: item.matched_product?.product_id || '',
            sku: item.matched_product?.sku || 'SPECIAL-ORDER',
            description: item.matched_product?.description || item.raw_text,
            quantity: item.quantity,
            uom: item.matched_product?.uom || item.uom,
            unit_price: item.matched_product?.base_price || 0,
            escalator: defaultEscalator(),
        }));
        setLines(prev => [...prev, ...newLines]);
        setAiSource(true);
        setLastParseResult(parseResult);
        setShowParsePanel(false);
        setParseResult(null);
        showToast(`${parsedItems.length} items added from material list`, 'success');
    }, [showToast, parseResult]);

    const handleSave = async () => {
        if (!customer) return;
        setLoading(true);
        try {
            const payload: CreateQuoteRequest = {
                customer_id: customer.id,
                source: aiSource ? 'ai' : 'manual',
                lines: lines.map(l => ({
                    product_id: l.product_id,
                    sku: l.sku,
                    description: l.description,
                    quantity: l.quantity,
                    uom: l.uom as import('../types/product').UOM,
                    unit_price: l.unit_price,
                })),
            };

            // Attach AI parse data if available
            if (aiSource && lastParseResult) {
                payload.parse_map = lastParseResult.items;
                // Store original file as base64 (source_image is already a data URI)
                if (lastParseResult.source_image) {
                    const [header, data] = lastParseResult.source_image.split(',');
                    const contentType = header?.match(/data:([^;]+)/)?.[1] || 'application/octet-stream';
                    payload.original_file = data;
                    payload.original_content_type = contentType;
                    payload.original_filename = 'material-list-upload';
                }
            }

            const quote = await QuoteService.createQuote(payload);
            showToast('Draft quote created', 'success');
            navigate(`/erp/quotes/${quote.id}`);
        } catch (err) {
            console.error(err);
            showToast('Failed to save quote', 'error');
        } finally {
            setLoading(false);
        }
    };

    const totalAmount = lines.reduce((sum, line) => sum + (line.quantity * line.unit_price), 0);
    const escalatedTotal = lines.reduce((sum, line) => {
        if (line.escalator.enabled && line.escalator.result) {
            return sum + (line.quantity * line.escalator.result.future_price);
        }
        return sum + (line.quantity * line.unit_price);
    }, 0);
    const hasEscalators = lines.some(l => l.escalator.enabled && l.escalator.result);
    const hasStaleLines = lines.some(l => l.escalator.result?.is_stale);
    const isOverLimit = customer ? (customer.balance_due + totalAmount) > customer.credit_limit : false;

    return (
        <PageTransition>
            <div className="flex flex-col md:flex-row md:items-center justify-between gap-4 mb-8">
                <div>
                    <h1 className="text-display-large text-white flex items-center gap-3">
                        <FileText className="w-10 h-10 text-gable-green" />
                        New Quote
                    </h1>
                    <p className="text-zinc-500 mt-1 max-w-2xl text-lg">
                        Draft a new pricing proposal.
                    </p>
                </div>
                <Button
                    onClick={handleSave}
                    disabled={!customer || lines.length === 0 || loading}
                    isLoading={loading}
                    className="shadow-glow"
                >
                    <Save className="w-4 h-4 mr-2" />
                    Create Quote
                </Button>
            </div>

            <div className="grid grid-cols-1 lg:grid-cols-12 gap-8">
                {/* Left Column: Customer & Details */}
                <div className="lg:col-span-4 space-y-6">
                    <Card variant="glass">
                        <CardContent className="p-6">
                            <h2 className="text-lg font-medium text-white mb-4 flex items-center gap-2">
                                <CreditCard className="w-5 h-5 text-zinc-400" />
                                Customer Details
                            </h2>
                            <CustomerSelect
                                onSelect={setCustomer}
                                selectedCustomerId={customer?.id}
                            />

                            {customer && (
                                <div className="mt-6 space-y-4 text-sm border-t border-white/5 pt-6">
                                    <div className="flex justify-between items-center bg-white/5 p-3 rounded-lg">
                                        <span className="text-zinc-400">Account #</span>
                                        <span className="font-mono text-white font-bold">{customer.account_number}</span>
                                    </div>
                                    <div className="flex justify-between items-center">
                                        <span className="text-zinc-400">Price Level</span>
                                        <span className="text-gable-green font-medium px-2 py-0.5 rounded bg-gable-green/10 border border-gable-green/20">
                                            {customer.price_level?.name || 'Retail'}
                                        </span>
                                    </div>
                                    <div className="space-y-2 pt-2">
                                        <div className="flex justify-between">
                                            <span className="text-zinc-400">Credit Limit</span>
                                            <span className="font-mono text-zinc-200">${customer.credit_limit?.toLocaleString() || '0.00'}</span>
                                        </div>
                                        <div className="flex justify-between">
                                            <span className="text-zinc-400">Balance Due</span>
                                            <span className={`font-mono ${customer.balance_due > customer.credit_limit ? 'text-rose-500 font-bold' : 'text-zinc-200'}`}>
                                                ${customer.balance_due.toLocaleString()}
                                            </span>
                                        </div>
                                        <div className="flex justify-between border-t border-white/5 pt-2">
                                            <span className="text-zinc-400">Available</span>
                                            <span className={`font-mono font-bold ${(customer.credit_limit - customer.balance_due) < 0 ? 'text-rose-500' : 'text-emerald-400'}`}>
                                                ${(customer.credit_limit - customer.balance_due).toLocaleString()}
                                            </span>
                                        </div>
                                    </div>
                                    {isOverLimit && (
                                        <div className="flex items-start gap-3 bg-rose-500/10 border border-rose-500/20 text-rose-400 text-xs p-3 rounded-lg">
                                            <AlertCircle className="w-4 h-4 shrink-0 mt-0.5" />
                                            <p>This quote exceeds the customer's credit limit. Approval will be required.</p>
                                        </div>
                                    )}
                                </div>
                            )}
                        </CardContent>
                    </Card>

                    <Card variant="glass" className="bg-gradient-to-br from-gable-green/5 to-emerald-900/5 border-gable-green/20">
                        <CardContent className="p-6">
                            <h2 className="text-lg font-medium text-white mb-4 flex items-center gap-2">
                                <Calculator className="w-5 h-5 text-gable-green" />
                                Quote Summary
                            </h2>
                            <div className="flex items-baseline justify-between">
                                <span className="text-zinc-400">Subtotal</span>
                                <span className="text-2xl font-mono font-bold text-white">${totalAmount.toFixed(2)}</span>
                            </div>

                            {/* Escalated Total */}
                            {hasEscalators && (
                                <div className="mt-3 pt-3 border-t border-white/5">
                                    <div className="flex items-baseline justify-between">
                                        <span className="text-zinc-400 flex items-center gap-1.5 text-sm">
                                            <TrendingUp className="w-3.5 h-3.5 text-gable-green" />
                                            Escalated Total
                                        </span>
                                        <span className="text-xl font-mono font-bold text-emerald-400">
                                            ${escalatedTotal.toFixed(2)}
                                        </span>
                                    </div>
                                    <div className="text-[10px] text-zinc-500 text-right mt-1">
                                        +${(escalatedTotal - totalAmount).toFixed(2)} from escalators
                                    </div>
                                </div>
                            )}

                            {/* Stale Lines Warning */}
                            {hasStaleLines && (
                                <div className="mt-3 flex items-center gap-2 bg-amber-500/10 border border-amber-500/20 text-amber-400 text-xs p-2.5 rounded-lg">
                                    <AlertCircle className="w-3.5 h-3.5 shrink-0" />
                                    Some lines have stale pricing
                                </div>
                            )}

                            <div className="text-xs text-zinc-500 text-right mt-1">Tax calculated at invoicing</div>
                        </CardContent>
                    </Card>
                </div>

                {/* Right Column: Lines */}
                <div className="lg:col-span-8 space-y-6">
                    <Card variant="glass" className="h-full">
                        <CardContent className="p-6">
                            <div className="flex items-center justify-between mb-6">
                                <h2 className="text-lg font-medium text-white">Line Items</h2>
                                <MaterialListUpload
                                    onParseComplete={handleParseComplete}
                                    disabled={loading}
                                />
                            </div>

                            <LineItemEditor products={products} customerId={customer?.id} onAddLine={handleAddLine} />

                            {/* Lines Table */}
                            <div className="mt-8 rounded-lg overflow-hidden border border-white/5 bg-black/20">
                                <table className="w-full text-sm text-left">
                                    <thead className="bg-white/5 text-zinc-400 uppercase tracking-wider text-xs font-semibold">
                                        <tr>
                                            <th className="px-6 py-4">SKU / Description</th>
                                            <th className="px-6 py-4 text-right">Qty</th>
                                            <th className="px-6 py-4 text-right">Unit Price</th>
                                            <th className="px-6 py-4 text-right">Total</th>
                                        </tr>
                                    </thead>
                                    <tbody className="divide-y divide-white/5">
                                        {lines.length === 0 && (
                                            <tr>
                                                <td colSpan={4} className="px-6 py-12 text-center text-zinc-500 italic">
                                                    No items added yet. Start building the quote above.
                                                </td>
                                            </tr>
                                        )}
                                        {lines.map((line, idx) => (
                                            <tr key={idx} className="group hover:bg-white/5 transition-colors">
                                                <td className="px-6 py-4">
                                                    <div className="font-mono text-white mb-0.5 group-hover:text-gable-green transition-colors">{line.sku}</div>
                                                    <div className="text-zinc-400 text-xs">{line.description}</div>

                                                    {/* Escalator Toggle */}
                                                    <EscalatorToggle
                                                        basePrice={line.unit_price}
                                                        escalator={line.escalator}
                                                        onChange={(esc) => handleEscalatorChange(idx, esc)}
                                                    />
                                                </td>
                                                <td className="px-6 py-4 text-right font-mono text-zinc-300 align-top">
                                                    {line.quantity} <span className="text-zinc-600 text-[10px] ml-1">{line.uom}</span>
                                                </td>
                                                <td className="px-6 py-4 text-right font-mono text-zinc-300 align-top">
                                                    ${line.unit_price.toFixed(2)}
                                                    {line.escalator.result && (
                                                        <div className="text-xs text-emerald-400 mt-1">
                                                            → ${line.escalator.result.future_price.toFixed(2)}
                                                        </div>
                                                    )}
                                                </td>
                                                <td className="px-6 py-4 text-right font-mono font-bold text-emerald-400 align-top">
                                                    ${(line.quantity * line.unit_price).toFixed(2)}
                                                    {line.escalator.result && (
                                                        <div className="text-xs text-emerald-300/70 mt-1">
                                                            → ${(line.quantity * line.escalator.result.future_price).toFixed(2)}
                                                        </div>
                                                    )}
                                                </td>
                                            </tr>
                                        ))}
                                    </tbody>
                                    {lines.length > 0 && (
                                        <tfoot className="bg-white/5 border-t border-white/10">
                                            <tr>
                                                <td colSpan={3} className="px-6 py-4 text-right font-medium text-zinc-400 uppercase tracking-wider text-xs">Total Amount</td>
                                                <td className="px-6 py-4 text-right font-mono text-xl font-bold text-gable-green">${totalAmount.toFixed(2)}</td>
                                            </tr>
                                        </tfoot>
                                    )}
                                </table>
                            </div>
                        </CardContent>
                    </Card>
                </div>
            </div>

            {/* AI Parse Results Overlay */}
            {showParsePanel && parseResult && (
                <ParsedResultsPanel
                    result={parseResult}
                    onAccept={handleAcceptParsed}
                    onClose={() => {
                        setShowParsePanel(false);
                        setParseResult(null);
                    }}
                />
            )}
        </PageTransition>
    );
};
