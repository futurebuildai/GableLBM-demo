import React from 'react';
import type { Product } from '../../types/product';

interface InventoryTableProps {
    products: Product[];
    onAdjustStock: (product: Product) => void;
}

export const InventoryTable: React.FC<InventoryTableProps> = ({ products, onAdjustStock }) => {
    return (
        <div className="w-full overflow-hidden border border-zinc-800 rounded-lg bg-zinc-900 text-sm">
            <div className="overflow-x-auto">
                <table className="w-full text-left text-zinc-400">
                    <thead className="bg-zinc-950 text-zinc-200 uppercase tracking-wider text-xs font-semibold">
                        <tr>
                            <th className="px-6 py-3 border-b border-zinc-800">SKU</th>
                            <th className="px-6 py-3 border-b border-zinc-800">Description</th>
                            <th className="px-6 py-3 border-b border-zinc-800">UOM</th>
                            <th className="px-6 py-3 border-b border-zinc-800 text-right">On Hand</th>
                            <th className="px-6 py-3 border-b border-zinc-800 text-right">Actions</th>
                        </tr>
                    </thead>
                    <tbody className="divide-y divide-zinc-800">
                        {products.length === 0 ? (
                            <tr>
                                <td colSpan={5} className="px-6 py-8 text-center text-zinc-600">
                                    No products found. Add items to the pile.
                                </td>
                            </tr>
                        ) : (
                            products.map((p) => (
                                <tr key={p.id} className="hover:bg-zinc-800/50 transition-colors">
                                    <td className="px-6 py-3 font-mono text-zinc-100">{p.sku}</td>
                                    <td className="px-6 py-3 text-zinc-300">{p.description}</td>
                                    <td className="px-6 py-3">
                                        <span className="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium bg-zinc-800 text-zinc-300 border border-zinc-700">
                                            {p.uom_primary}
                                        </span>
                                    </td>
                                    <td className="px-6 py-3 text-right font-mono text-emerald-400">
                                        {(p.total_quantity || 0).toFixed(4)}
                                    </td>
                                    <td className="px-6 py-3 text-right">
                                        <button
                                            onClick={() => onAdjustStock(p)}
                                            className="text-emerald-400 hover:text-emerald-300 hover:underline"
                                        >
                                            Adjust
                                        </button>
                                    </td>
                                </tr>
                            ))
                        )}
                    </tbody>
                </table>
            </div>
        </div>
    );
};
