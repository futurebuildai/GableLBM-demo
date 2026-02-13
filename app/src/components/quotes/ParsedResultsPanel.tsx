import { useState, useCallback } from 'react';
import { Check, AlertTriangle, ChevronDown, X, Sparkles, Package, ShieldAlert } from 'lucide-react';
import { Card, CardContent } from '../ui/Card';
import { Button } from '../ui/Button';
import type { ParseResponse, ParsedItem, MatchedProduct } from '../../types/parsing';

interface ParsedResultsPanelProps {
    result: ParseResponse;
    onAccept: (items: ParsedItem[]) => void;
    onClose: () => void;
}

/**
 * Side-by-side review panel: Original image vs. Parsed line items.
 * Shows confidence badges, swap SKU dropdowns, and accept controls.
 */
export const ParsedResultsPanel = ({ result, onAccept, onClose }: ParsedResultsPanelProps) => {
    const [items, setItems] = useState<ParsedItem[]>(result.items);
    const [selectedIndices, setSelectedIndices] = useState<Set<number>>(
        new Set(result.items.map((_, i) => i))
    );

    const toggleSelect = useCallback((index: number) => {
        setSelectedIndices(prev => {
            const next = new Set(prev);
            if (next.has(index)) {
                next.delete(index);
            } else {
                next.add(index);
            }
            return next;
        });
    }, []);

    const handleSwapSku = useCallback((itemIndex: number, alt: MatchedProduct) => {
        setItems(prev => {
            const updated = [...prev];
            const item = { ...updated[itemIndex] };
            // Move current match to alternatives, use selected alt as match
            const currentAlts = [...(item.alternatives || [])];
            if (item.matched_product) {
                currentAlts.push(item.matched_product);
            }
            item.matched_product = alt;
            item.alternatives = currentAlts.filter(a => a.product_id !== alt.product_id);
            item.confidence = 0.95; // User-confirmed swap = high confidence
            item.is_special_order = false;
            updated[itemIndex] = item;
            return updated;
        });
    }, []);

    const handleAcceptAll = useCallback(() => {
        onAccept(items);
    }, [items, onAccept]);

    const handleAcceptSelected = useCallback(() => {
        const selected = items.filter((_, i) => selectedIndices.has(i));
        onAccept(selected);
    }, [items, selectedIndices, onAccept]);

    const highConfCount = items.filter(i => i.confidence >= 0.9 && !i.is_special_order).length;
    const lowConfCount = items.filter(i => i.confidence < 0.9 && i.confidence >= 0.5 && !i.is_special_order).length;
    const specialOrderCount = items.filter(i => i.is_special_order).length;

    return (
        <div className="fixed inset-0 bg-black/70 backdrop-blur-sm z-50 flex items-center justify-center p-4" id="parsed-results-overlay">
            <div className="bg-slate-steel rounded-2xl border border-white/10 w-full max-w-6xl max-h-[90vh] flex flex-col shadow-2xl">
                {/* Header */}
                <div className="flex items-center justify-between p-6 border-b border-white/5">
                    <div className="flex items-center gap-3">
                        <div className="w-10 h-10 rounded-xl bg-gable-green/10 border border-gable-green/20 flex items-center justify-center">
                            <Sparkles className="w-5 h-5 text-gable-green" />
                        </div>
                        <div>
                            <h2 className="text-lg font-semibold text-white">AI Parse Results</h2>
                            <p className="text-xs text-zinc-500">
                                {result.item_count} items parsed in {result.parse_time_ms}ms
                                <span className="mx-2">•</span>
                                <span className="text-emerald-400">{highConfCount} auto-matched</span>
                                {lowConfCount > 0 && (
                                    <span className="text-amber-400 ml-2">{lowConfCount} review needed</span>
                                )}
                                {specialOrderCount > 0 && (
                                    <span className="text-rose-400 ml-2">{specialOrderCount} special order</span>
                                )}
                            </p>
                        </div>
                    </div>
                    <button
                        onClick={onClose}
                        className="text-zinc-500 hover:text-white transition-colors p-2 rounded-lg hover:bg-white/5"
                        id="close-parse-results-btn"
                    >
                        <X className="w-5 h-5" />
                    </button>
                </div>

                {/* Content: Side by Side */}
                <div className="flex-1 overflow-hidden grid grid-cols-1 lg:grid-cols-2 gap-0 min-h-0">
                    {/* Left: Original Image */}
                    <div className="border-r border-white/5 p-6 overflow-auto">
                        <h3 className="text-sm font-medium text-zinc-400 uppercase tracking-wider mb-4">Original Document</h3>
                        {result.source_image ? (
                            <div className="rounded-xl overflow-hidden border border-white/10 bg-black/30">
                                <img
                                    src={result.source_image}
                                    alt="Uploaded material list"
                                    className="w-full h-auto object-contain max-h-[60vh]"
                                />
                            </div>
                        ) : (
                            <div className="rounded-xl border border-white/10 bg-black/30 p-12 text-center text-zinc-500">
                                No preview available
                            </div>
                        )}
                    </div>

                    {/* Right: Parsed Items */}
                    <div className="p-6 overflow-auto">
                        <h3 className="text-sm font-medium text-zinc-400 uppercase tracking-wider mb-4">Parsed Line Items</h3>
                        <div className="space-y-3">
                            {items.map((item, idx) => (
                                <ParsedItemCard
                                    key={idx}
                                    item={item}
                                    index={idx}
                                    selected={selectedIndices.has(idx)}
                                    onToggleSelect={() => toggleSelect(idx)}
                                    onSwapSku={(alt) => handleSwapSku(idx, alt)}
                                />
                            ))}
                        </div>
                    </div>
                </div>

                {/* Footer: Actions */}
                <div className="flex items-center justify-between p-6 border-t border-white/5 bg-black/20">
                    <div className="text-sm text-zinc-500">
                        {selectedIndices.size} of {items.length} items selected
                    </div>
                    <div className="flex items-center gap-3">
                        <Button variant="secondary" onClick={onClose}>
                            Cancel
                        </Button>
                        {selectedIndices.size < items.length && (
                            <Button variant="secondary" onClick={handleAcceptSelected} id="accept-selected-btn">
                                <Check className="w-4 h-4 mr-2" />
                                Accept Selected ({selectedIndices.size})
                            </Button>
                        )}
                        <Button onClick={handleAcceptAll} className="shadow-glow" id="accept-all-btn">
                            <Check className="w-4 h-4 mr-2" />
                            Accept All ({items.length})
                        </Button>
                    </div>
                </div>
            </div>
        </div>
    );
};

// --- Sub-components ---

interface ParsedItemCardProps {
    item: ParsedItem;
    index: number;
    selected: boolean;
    onToggleSelect: () => void;
    onSwapSku: (alt: MatchedProduct) => void;
}

const ParsedItemCard = ({ item, index, selected, onToggleSelect, onSwapSku }: ParsedItemCardProps) => {
    const [showAlts, setShowAlts] = useState(false);

    return (
        <Card
            variant="glass"
            className={`transition-all duration-200 ${selected
                    ? 'border-gable-green/30 bg-gable-green/5'
                    : 'border-white/5 opacity-50'
                }`}
        >
            <CardContent className="p-4">
                <div className="flex items-start gap-3">
                    {/* Checkbox */}
                    <button
                        onClick={onToggleSelect}
                        className={`mt-0.5 w-5 h-5 rounded border flex items-center justify-center shrink-0 transition-all ${selected
                                ? 'bg-gable-green border-gable-green text-black'
                                : 'border-white/20 hover:border-white/40'
                            }`}
                        id={`toggle-item-${index}`}
                    >
                        {selected && <Check className="w-3 h-3" />}
                    </button>

                    {/* Content */}
                    <div className="flex-1 min-w-0">
                        {/* Raw text + Confidence badge */}
                        <div className="flex items-center gap-2 mb-2">
                            <ConfidenceBadge confidence={item.confidence} isSpecialOrder={item.is_special_order} />
                            <span className="text-xs text-zinc-500 font-mono truncate">{item.raw_text}</span>
                        </div>

                        {/* Matched Product */}
                        {item.matched_product ? (
                            <div className="flex items-center gap-3">
                                <Package className="w-4 h-4 text-zinc-500 shrink-0" />
                                <div className="flex-1 min-w-0">
                                    <div className="font-mono text-white text-sm">{item.matched_product.sku}</div>
                                    <div className="text-xs text-zinc-400 truncate">{item.matched_product.description}</div>
                                </div>
                                <div className="text-right shrink-0">
                                    <div className="font-mono text-white text-sm">
                                        {item.quantity} <span className="text-zinc-500 text-[10px]">{item.uom}</span>
                                    </div>
                                    <div className="font-mono text-emerald-400 text-xs">
                                        ${item.matched_product.base_price.toFixed(2)}
                                    </div>
                                </div>
                            </div>
                        ) : (
                            <div className="flex items-center gap-3">
                                <ShieldAlert className="w-4 h-4 text-rose-400 shrink-0" />
                                <div className="flex-1">
                                    <div className="text-sm text-rose-300">Special Order</div>
                                    <div className="text-xs text-zinc-500">{item.raw_text}</div>
                                </div>
                                <div className="text-right shrink-0">
                                    <div className="font-mono text-white text-sm">
                                        {item.quantity} <span className="text-zinc-500 text-[10px]">{item.uom}</span>
                                    </div>
                                </div>
                            </div>
                        )}

                        {/* Swap SKU dropdown for low-confidence items */}
                        {item.alternatives && item.alternatives.length > 0 && (
                            <div className="mt-2">
                                <button
                                    onClick={() => setShowAlts(!showAlts)}
                                    className="text-xs text-blueprint-blue hover:text-blue-300 flex items-center gap-1 transition-colors"
                                    id={`swap-sku-toggle-${index}`}
                                >
                                    <ChevronDown className={`w-3 h-3 transition-transform ${showAlts ? 'rotate-180' : ''}`} />
                                    Swap SKU ({item.alternatives.length} alternatives)
                                </button>

                                {showAlts && (
                                    <div className="mt-2 space-y-1.5 pl-4 border-l-2 border-blueprint-blue/20">
                                        {item.alternatives.map((alt) => (
                                            <button
                                                key={alt.product_id}
                                                onClick={() => {
                                                    onSwapSku(alt);
                                                    setShowAlts(false);
                                                }}
                                                className="w-full text-left p-2 rounded-lg bg-white/5 hover:bg-white/10 transition-colors group"
                                            >
                                                <div className="flex items-center justify-between">
                                                    <div>
                                                        <span className="font-mono text-xs text-zinc-300 group-hover:text-white">{alt.sku}</span>
                                                        <span className="text-xs text-zinc-500 ml-2">{alt.description}</span>
                                                    </div>
                                                    <span className="font-mono text-xs text-emerald-400">${alt.base_price.toFixed(2)}</span>
                                                </div>
                                            </button>
                                        ))}
                                    </div>
                                )}
                            </div>
                        )}
                    </div>
                </div>
            </CardContent>
        </Card>
    );
};

// --- Confidence Badge ---

const ConfidenceBadge = ({ confidence, isSpecialOrder }: { confidence: number; isSpecialOrder: boolean }) => {
    if (isSpecialOrder) {
        return (
            <span className="inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-[10px] font-bold uppercase tracking-wider bg-rose-500/15 text-rose-400 border border-rose-500/20">
                <AlertTriangle className="w-2.5 h-2.5" />
                Special
            </span>
        );
    }

    const pct = Math.round(confidence * 100);

    if (confidence >= 0.9) {
        return (
            <span className="inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-[10px] font-bold uppercase tracking-wider bg-emerald-500/15 text-emerald-400 border border-emerald-500/20">
                <Check className="w-2.5 h-2.5" />
                {pct}%
            </span>
        );
    }

    return (
        <span className="inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-[10px] font-bold uppercase tracking-wider bg-amber-500/15 text-amber-400 border border-amber-500/20">
            <AlertTriangle className="w-2.5 h-2.5" />
            {pct}%
        </span>
    );
};
