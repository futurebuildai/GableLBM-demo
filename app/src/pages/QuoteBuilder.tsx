import { useState, useEffect } from 'react';
import { CustomerSelect } from '../components/customers/CustomerSelect';
import { LineItemEditor } from '../components/quotes/LineItemEditor';
import { QuoteService } from '../services/QuoteService';
import { ProductService } from '../services/product.service';
import type { Customer } from '../types/customer';
import type { Product } from '../types/product';
import type { CreateQuoteRequest } from '../types/quote';
import { Save, FileText } from 'lucide-react';
import { useNavigate } from 'react-router-dom';

export const QuoteBuilder = () => {
    const navigate = useNavigate();
    const [customer, setCustomer] = useState<Customer | null>(null);
    const [products, setProducts] = useState<Product[]>([]);
    const [lines, setLines] = useState<CreateQuoteRequest['lines']>([]);
    const [loading, setLoading] = useState(false);

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
            unit_price: unitPrice
        }]);
    };

    const handleSave = async () => {
        if (!customer) return;
        setLoading(true);
        try {
            await QuoteService.createQuote({
                customer_id: customer.id,
                lines
            });
            // Show success or navigate?
            alert('Quote created successfully!');
            navigate('/'); // Or to quote list
        } catch (err) {
            console.error(err);
            alert('Failed to save quote');
        } finally {
            setLoading(false);
        }
    };

    const totalAmount = lines.reduce((sum, line) => sum + (line.quantity * line.unit_price), 0);

    return (
        <div className="p-6 max-w-7xl mx-auto text-white">
            <header className="flex justify-between items-center mb-8">
                <div className="flex items-center">
                    <div className="p-3 bg-[#00FFA3]/10 rounded-lg mr-4">
                        <FileText className="w-8 h-8 text-[#00FFA3]" />
                    </div>
                    <div>
                        <h1 className="text-2xl font-bold tracking-tight">New Quote</h1>
                        <p className="text-gray-500">Create a quick quote for a customer.</p>
                    </div>
                </div>
                <div className="flex gap-2">
                    <button
                        onClick={handleSave}
                        disabled={!customer || lines.length === 0 || loading}
                        className="flex items-center bg-[#00FFA3] text-black px-6 py-2 rounded font-medium hover:bg-[#00FFA3]/80 disabled:opacity-50 transition-all shadow-[0_0_20px_rgba(0,255,163,0.3)]"
                    >
                        <Save className="w-4 h-4 mr-2" />
                        {loading ? 'Saving...' : 'Save Quote'}
                    </button>
                </div>
            </header>

            <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
                {/* Left Column: Customer & Details */}
                <div className="lg:col-span-1 space-y-6">
                    <section className="bg-[#161821] border border-white/10 rounded-xl p-6">
                        <h2 className="text-lg font-medium mb-4">Customer Details</h2>
                        <CustomerSelect
                            onSelect={setCustomer}
                            selectedCustomerId={customer?.id}
                        />

                        {customer && (
                            <div className="mt-6 space-y-3 text-sm border-t border-white/10 pt-4">
                                <div className="flex justify-between">
                                    <span className="text-gray-500">Account #</span>
                                    <span className="font-mono">{customer.account_number}</span>
                                </div>
                                <div className="flex justify-between">
                                    <span className="text-gray-500">Price Level</span>
                                    <span className="text-[#00FFA3]">{customer.price_level?.name || 'Retail'}</span>
                                </div>
                                <div className="flex justify-between">
                                    <span className="text-gray-500">Credit Limit</span>
                                    <span className="font-mono">${customer.credit_limit?.toLocaleString() || '0.00'}</span>
                                </div>
                            </div>
                        )}
                    </section>
                </div>

                {/* Right Column: Lines */}
                <div className="lg:col-span-2 space-y-6">
                    <section className="bg-[#161821] border border-white/10 rounded-xl p-6 min-h-[600px]">
                        <h2 className="text-lg font-medium mb-4">Line Items</h2>

                        <LineItemEditor products={products} onAddLine={handleAddLine} />

                        {/* Lines Table */}
                        <div className="mt-6 overflow-hidden border border-white/10 rounded-lg">
                            <table className="w-full text-sm text-left">
                                <thead className="bg-white/5 text-gray-400">
                                    <tr>
                                        <th className="px-4 py-3 font-medium">SKU</th>
                                        <th className="px-4 py-3 font-medium">Description</th>
                                        <th className="px-4 py-3 font-medium text-right">Qty</th>
                                        <th className="px-4 py-3 font-medium text-right">Price</th>
                                        <th className="px-4 py-3 font-medium text-right">Total</th>
                                    </tr>
                                </thead>
                                <tbody className="divide-y divide-white/5">
                                    {lines.length === 0 && (
                                        <tr>
                                            <td colSpan={5} className="px-4 py-8 text-center text-gray-600 italic">
                                                No items added yet.
                                            </td>
                                        </tr>
                                    )}
                                    {lines.map((line, idx) => (
                                        <tr key={idx} className="hover:bg-white/5 transition-colors">
                                            <td className="px-4 py-3 font-mono text-[#00FFA3]">{line.sku}</td>
                                            <td className="px-4 py-3">{line.description}</td>
                                            <td className="px-4 py-3 text-right font-mono">
                                                {line.quantity} <span className="text-gray-600 text-xs">{line.uom}</span>
                                            </td>
                                            <td className="px-4 py-3 text-right font-mono">${line.unit_price.toFixed(2)}</td>
                                            <td className="px-4 py-3 text-right font-mono font-bold">${(line.quantity * line.unit_price).toFixed(2)}</td>
                                        </tr>
                                    ))}
                                </tbody>
                                <tfoot className="bg-white/5 border-t border-white/10">
                                    <tr>
                                        <td colSpan={4} className="px-4 py-4 text-right font-medium text-gray-400 uppercase tracking-wider">Total Amount</td>
                                        <td className="px-4 py-4 text-right font-mono text-xl font-bold text-[#00FFA3]">${totalAmount.toFixed(2)}</td>
                                    </tr>
                                </tfoot>
                            </table>
                        </div>
                    </section>
                </div>
            </div>
        </div>
    );
};
