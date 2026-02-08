import type { TopCustomer } from '../../types/dashboard';

interface TopCustomersTableProps {
    customers: TopCustomer[];
    loading?: boolean;
}

export function TopCustomersTable({ customers, loading }: TopCustomersTableProps) {
    if (loading) {
        return (
            <div className="p-4 bg-zinc-900 rounded-lg border border-white/10">
                <h3 className="text-sm font-medium text-zinc-400 mb-4">Top Customers (30 Days)</h3>
                <div className="space-y-3">
                    {[1, 2, 3, 4, 5].map((i) => (
                        <div key={i} className="h-8 bg-zinc-800 rounded animate-pulse" />
                    ))}
                </div>
            </div>
        );
    }

    if (customers.length === 0) {
        return (
            <div className="p-4 bg-zinc-900 rounded-lg border border-white/10">
                <h3 className="text-sm font-medium text-zinc-400 mb-4">Top Customers (30 Days)</h3>
                <div className="text-center py-8 text-zinc-500">No customer data available</div>
            </div>
        );
    }

    return (
        <div className="p-4 bg-zinc-900 rounded-lg border border-white/10">
            <h3 className="text-sm font-medium text-zinc-400 mb-4">Top Customers (30 Days)</h3>
            <table className="w-full">
                <thead>
                    <tr className="text-xs text-zinc-500 uppercase">
                        <th className="text-left pb-2">#</th>
                        <th className="text-left pb-2">Customer</th>
                        <th className="text-right pb-2">Revenue</th>
                        <th className="text-right pb-2">Orders</th>
                    </tr>
                </thead>
                <tbody>
                    {customers.map((customer, index) => (
                        <tr
                            key={customer.customer_id}
                            className="border-t border-white/5 hover:bg-white/5 transition-colors"
                        >
                            <td className="py-2 text-zinc-500 text-sm">{index + 1}</td>
                            <td className="py-2 text-white text-sm truncate max-w-[140px]">
                                {customer.customer_name}
                            </td>
                            <td className="py-2 text-right font-mono text-emerald-400 text-sm">
                                ${(customer.total_revenue / 100).toLocaleString()}
                            </td>
                            <td className="py-2 text-right font-mono text-zinc-400 text-sm">
                                {customer.order_count}
                            </td>
                        </tr>
                    ))}
                </tbody>
            </table>
        </div>
    );
}
