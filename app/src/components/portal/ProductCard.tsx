import { Package, Plus, CheckCircle, XCircle } from 'lucide-react';
import type { CatalogProduct } from '../../types/portal';
import { Link } from 'react-router-dom';

const formatCurrency = (val: number): string =>
    new Intl.NumberFormat('en-US', { style: 'currency', currency: 'USD' }).format(val);

interface ProductCardProps {
    product: CatalogProduct;
    onAddToCart: (productId: string, quantity: number) => void;
    adding?: boolean;
}

export const ProductCard = ({ product, onAddToCart, adding }: ProductCardProps) => {
    return (
        <div className="group relative bg-white/[0.03] border border-white/[0.06] rounded-2xl overflow-hidden hover:border-white/10 hover:-translate-y-1 transition-all duration-300">
            {/* Image Area */}
            <Link to={`/portal/catalog/${product.id}`} className="block">
                <div className="aspect-square bg-gradient-to-br from-zinc-800/50 to-zinc-900/50 flex items-center justify-center p-6">
                    {product.image_url ? (
                        <img
                            src={product.image_url}
                            alt={product.name}
                            className="w-full h-full object-contain"
                        />
                    ) : (
                        <Package className="w-16 h-16 text-zinc-600 group-hover:text-zinc-500 transition-colors" />
                    )}
                </div>
            </Link>

            {/* Content */}
            <div className="p-4 space-y-3">
                {/* Category Badge */}
                {product.category && (
                    <span className="inline-block px-2 py-0.5 rounded-full text-[10px] uppercase tracking-wider font-semibold bg-gable-green/10 text-gable-green border border-gable-green/20">
                        {product.category}
                    </span>
                )}

                {/* Product Info */}
                <Link to={`/portal/catalog/${product.id}`} className="block">
                    <h3 className="text-sm font-semibold text-white leading-tight group-hover:text-gable-green transition-colors line-clamp-2">
                        {product.name}
                    </h3>
                    <p className="text-xs text-zinc-500 font-mono mt-1">{product.sku}</p>
                </Link>

                {/* Pricing */}
                <div className="flex items-baseline gap-2">
                    <span className="text-lg font-bold text-white font-mono">
                        {formatCurrency(product.customer_price)}
                    </span>
                    {product.customer_price < product.base_price && (
                        <span className="text-xs text-zinc-500 line-through font-mono">
                            {formatCurrency(product.base_price)}
                        </span>
                    )}
                    <span className="text-[10px] text-zinc-600">/{product.uom}</span>
                </div>

                {/* Availability */}
                <div className="flex items-center gap-1.5">
                    {product.in_stock ? (
                        <>
                            <CheckCircle className="w-3.5 h-3.5 text-emerald-400" />
                            <span className="text-xs text-emerald-400">
                                {product.available > 0
                                    ? `${Math.floor(product.available)} in stock`
                                    : 'In Stock'}
                            </span>
                        </>
                    ) : (
                        <>
                            <XCircle className="w-3.5 h-3.5 text-red-400" />
                            <span className="text-xs text-red-400">Out of Stock</span>
                        </>
                    )}
                </div>

                {/* Add to Cart */}
                <button
                    onClick={() => onAddToCart(product.id, 1)}
                    disabled={adding || !product.in_stock}
                    className="w-full flex items-center justify-center gap-2 py-2.5 rounded-xl text-sm font-semibold transition-all duration-200
                        bg-gable-green/10 text-gable-green border border-gable-green/20
                        hover:bg-gable-green/20 hover:border-gable-green/30
                        disabled:opacity-40 disabled:cursor-not-allowed
                        active:scale-[0.98]"
                >
                    <Plus className="w-4 h-4" />
                    {adding ? 'Adding...' : 'Add to Cart'}
                </button>
            </div>
        </div>
    );
};
