import React from 'react';
import type { ProductDetail, PIMMedia } from '../../../types/pim';
import { Package, Weight, BarChart3, DollarSign, Layers, Tag } from 'lucide-react';

interface Props {
    product: ProductDetail;
}

export const ProductOverviewTab: React.FC<Props> = ({ product }) => {
    const available = (product.total_quantity || 0) - (product.total_allocated || 0);
    const primaryImage = product.media?.find((m: PIMMedia) => m.is_primary) || product.media?.[0];
    const visiblePrice = product.base_price || 0;
    const margin = product.target_margin || 0;

    return (
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
            {/* Primary Image */}
            <div className="lg:col-span-1">
                <div className="bg-zinc-900 border border-white/10 rounded-xl overflow-hidden aspect-square flex items-center justify-center">
                    {primaryImage ? (
                        <img src={primaryImage.url} alt={primaryImage.alt_text || product.description} className="w-full h-full object-cover" />
                    ) : (
                        <div className="flex flex-col items-center gap-3 text-zinc-500">
                            <Package className="w-16 h-16" />
                            <span className="text-sm">No image</span>
                        </div>
                    )}
                </div>
            </div>

            {/* Product Info */}
            <div className="lg:col-span-2 space-y-6">
                {/* Info Grid */}
                <div className="grid grid-cols-2 sm:grid-cols-3 gap-4">
                    <InfoCard icon={<Tag className="w-4 h-4" />} label="SKU" value={product.sku} />
                    <InfoCard icon={<Layers className="w-4 h-4" />} label="UOM" value={product.uom_primary} />
                    <InfoCard icon={<Package className="w-4 h-4" />} label="Vendor" value={product.vendor || 'N/A'} />
                    <InfoCard icon={<Weight className="w-4 h-4" />} label="Weight" value={`${(product.weight_lbs || 0).toFixed(1)} lbs`} />
                    <InfoCard icon={<DollarSign className="w-4 h-4" />} label="Avg Cost" value={`$${(product.average_unit_cost || 0).toFixed(2)}`} accent="emerald" />
                    <InfoCard icon={<DollarSign className="w-4 h-4" />} label="Base Price" value={`$${visiblePrice.toFixed(2)}`} accent="green" />
                    <InfoCard icon={<BarChart3 className="w-4 h-4" />} label="Target Margin" value={`${margin.toFixed(1)}%`} />
                    <InfoCard icon={<BarChart3 className="w-4 h-4" />} label="Commission" value={`${(product.commission_rate || 0).toFixed(1)}%`} />
                    {product.upc && <InfoCard icon={<Tag className="w-4 h-4" />} label="UPC" value={product.upc} />}
                </div>

                {/* Stock Summary */}
                <div>
                    <h3 className="text-sm font-medium text-zinc-400 uppercase tracking-wider mb-3">Stock Summary</h3>
                    <div className="grid grid-cols-3 gap-4">
                        <StockCard label="On Hand" value={product.total_quantity || 0} />
                        <StockCard label="Allocated" value={product.total_allocated || 0} color="amber" />
                        <StockCard label="Available" value={available} color={available < 100 ? 'rose' : 'emerald'} />
                    </div>
                </div>

                {/* Reorder Info */}
                {(product.reorder_point || 0) > 0 && (
                    <div className="bg-zinc-900 border border-white/10 rounded-xl p-4">
                        <h3 className="text-sm font-medium text-zinc-400 uppercase tracking-wider mb-2">Reorder Settings</h3>
                        <div className="flex gap-6 text-sm">
                            <div>
                                <span className="text-zinc-500">Reorder Point: </span>
                                <span className="text-white font-mono">{(product.reorder_point || 0).toLocaleString()}</span>
                            </div>
                            <div>
                                <span className="text-zinc-500">Reorder Qty: </span>
                                <span className="text-white font-mono">{(product.reorder_qty || 0).toLocaleString()}</span>
                            </div>
                        </div>
                    </div>
                )}

                {/* PIM Content Preview */}
                {product.content?.short_description && (
                    <div className="bg-zinc-900 border border-white/10 rounded-xl p-4">
                        <h3 className="text-sm font-medium text-zinc-400 uppercase tracking-wider mb-2">Description</h3>
                        <p className="text-zinc-300 text-sm">{product.content.short_description}</p>
                    </div>
                )}
            </div>
        </div>
    );
};

const InfoCard: React.FC<{ icon: React.ReactNode; label: string; value: string; accent?: string }> = ({ icon, label, value, accent }) => (
    <div className="bg-zinc-900 border border-white/10 rounded-lg p-3">
        <div className="flex items-center gap-1.5 text-zinc-500 text-xs mb-1">
            {icon}
            {label}
        </div>
        <div className={`font-mono text-sm font-medium ${accent === 'emerald' ? 'text-emerald-400' : accent === 'green' ? 'text-gable-green' : 'text-white'}`}>
            {value}
        </div>
    </div>
);

const StockCard: React.FC<{ label: string; value: number; color?: string }> = ({ label, value, color = 'white' }) => (
    <div className="bg-zinc-900 border border-white/10 rounded-lg p-4 text-center">
        <div className="text-xs text-zinc-500 mb-1">{label}</div>
        <div className={`text-2xl font-mono font-bold ${
            color === 'emerald' ? 'text-emerald-400' :
            color === 'amber' ? 'text-amber-400' :
            color === 'rose' ? 'text-rose-500' :
            'text-white'
        }`}>
            {value.toLocaleString()}
        </div>
    </div>
);
