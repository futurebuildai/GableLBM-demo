import { useEffect, useState } from 'react';
import { ProductService } from '../services/product.service';
import type { Product } from '../types/product';
import { InventoryTable } from '../components/inventory/InventoryTable';
import { AddProductModal } from '../components/inventory/AddProductModal';
import { StockAdjustmentModal } from '../components/inventory/StockAdjustmentModal';
import { InventoryTransferModal } from '../components/inventory/InventoryTransferModal';
import { Plus, Search } from 'lucide-react';

export const Inventory = () => {
    const [products, setProducts] = useState<Product[]>([]);
    const [isLoading, setIsLoading] = useState(true);
    const [isModalOpen, setIsModalOpen] = useState(false);
    const [error, setError] = useState('');

    // Stock Adjustment State
    const [isStockModalOpen, setIsStockModalOpen] = useState(false);
    const [isTransferModalOpen, setIsTransferModalOpen] = useState(false);
    const [selectedProduct, setSelectedProduct] = useState<Product | null>(null);

    const loadProducts = async () => {
        try {
            setIsLoading(true);
            const data = await ProductService.getProducts();
            setProducts(data);
            setError('');
        } catch (err) {
            setError('Failed to load products');
            console.error(err);
        } finally {
            setIsLoading(false);
        }
    };

    useEffect(() => {
        loadProducts();
    }, []);

    const handleSaveProduct = async (productData: Omit<Product, 'id' | 'created_at' | 'updated_at'>) => {
        await ProductService.createProduct(productData);
        await loadProducts(); // Refresh list
    };

    const handleAdjustStock = (product: Product) => {
        setSelectedProduct(product);
        setIsStockModalOpen(true);
    };

    const handleTransferStock = (product: Product) => {
        setSelectedProduct(product);
        setIsTransferModalOpen(true);
    };

    return (
        <div className="p-8 max-w-7xl mx-auto">
            <div className="flex justify-between items-center mb-8">
                <div>
                    <h1 className="text-3xl font-bold text-zinc-100 tracking-tight">The Pile</h1>
                    <p className="text-zinc-500 mt-1">Master Inventory & SKU Management</p>
                </div>
                <button
                    onClick={() => setIsModalOpen(true)}
                    className="flex items-center gap-2 bg-amber-600 hover:bg-amber-500 text-white px-4 py-2 rounded font-medium transition-colors shadow-lg shadow-amber-900/20"
                >
                    <Plus className="w-4 h-4" />
                    Add Product
                </button>
            </div>

            <div className="mb-6 flex gap-4">
                {/* Search Bar - Placeholder for now */}
                <div className="relative flex-1 max-w-md">
                    <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-zinc-500" />
                    <input
                        type="text"
                        placeholder="Search SKUs..."
                        className="w-full bg-zinc-900 border border-zinc-800 rounded pl-10 pr-4 py-2 text-zinc-300 focus:outline-none focus:ring-1 focus:ring-amber-500"
                    />
                </div>
            </div>

            {error && (
                <div className="mb-6 p-4 bg-red-900/20 border border-red-900/50 text-red-400 rounded">
                    {error}
                </div>
            )}

            {isLoading ? (
                <div className="text-zinc-500 text-center py-12 animate-pulse">Loading core inventory...</div>
            ) : (
                <InventoryTable
                    products={products}
                    onAdjustStock={handleAdjustStock}
                    onTransferStock={handleTransferStock}
                />
            )}

            <AddProductModal
                isOpen={isModalOpen}
                onClose={() => setIsModalOpen(false)}
                onSave={handleSaveProduct}
            />

            <StockAdjustmentModal
                isOpen={isStockModalOpen}
                onClose={() => setIsStockModalOpen(false)}
                product={selectedProduct}
                onSuccess={() => {
                    loadProducts(); // Reload to update "On Hand" if we calculate it (currently hardcoded 0.0000, need backend update to show real sum)
                }}
            />

            <InventoryTransferModal
                isOpen={isTransferModalOpen}
                onClose={() => setIsTransferModalOpen(false)}
                product={selectedProduct}
                onSuccess={() => {
                    loadProducts();
                }}
            />
        </div>
    );
};
