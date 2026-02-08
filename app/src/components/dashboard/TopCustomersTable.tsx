import { Card, CardHeader, CardTitle, CardContent } from '../ui/Card';
import type { TopCustomer } from '../../types/dashboard';

interface TopCustomersTableProps {
    customers: TopCustomer[];
    loading?: boolean;
}

export function TopCustomersTable({ customers, loading = false }: TopCustomersTableProps) {
    if (loading) {
        return (
            <Card variant="glass" className="h-full">
                <CardHeader>
                    <div className="h-6 w-32 bg-white/10 rounded animate-pulse" />
                </CardHeader>
                <CardContent>
                    <div className="space-y-4">
                        {[1, 2, 3, 4, 5].map((i) => (
                            <div key={i} className="flex justify-between items-center">
                                <div className="flex items-center gap-3">
                                    <div className="h-8 w-8 rounded-full bg-white/10 animate-pulse" />
                                    <div className="h-4 w-40 bg-white/10 rounded animate-pulse" />
                                </div>
                                <div className="h-4 w-20 bg-white/10 rounded animate-pulse" />
                            </div>
                        ))}
                    </div>
                </CardContent>
            </Card>
        );
    }

    return (
        <Card variant="glass" className="h-full">
            <CardHeader>
                <CardTitle>Top Customers</CardTitle>
            </CardHeader>
            <CardContent className="p-0">
                <div className="overflow-x-auto">
                    <table className="w-full text-sm text-left">
                        <thead className="text-zinc-400 font-medium border-b border-white/5 bg-white/5">
                            <tr>
                                <th className="px-6 py-3">Customer</th>
                                <th className="px-6 py-3 text-right">Orders</th>
                                <th className="px-6 py-3 text-right">Revenue</th>
                            </tr>
                        </thead>
                        <tbody className="divide-y divide-white/5">
                            {customers.length === 0 ? (
                                <tr>
                                    <td colSpan={3} className="px-6 py-8 text-center text-zinc-500">
                                        No customer data available
                                    </td>
                                </tr>
                            ) : (
                                customers.map((customer) => (
                                    <tr key={customer.customer_id} className="group hover:bg-white/5 transition-colors">
                                        <td className="px-6 py-3 font-medium text-white flex items-center gap-3">
                                            <div className="h-8 w-8 rounded-full bg-gradient-to-br from-zinc-700 to-zinc-800 flex items-center justify-center text-xs font-bold text-white border border-white/10 shadow-sm group-hover:scale-105 transition-transform">
                                                {customer.customer_name.substring(0, 2).toUpperCase()}
                                            </div>
                                            {customer.customer_name}
                                        </td>
                                        <td className="px-6 py-3 text-right font-mono text-zinc-400">
                                            {customer.order_count}
                                        </td>
                                        <td className="px-6 py-3 text-right font-mono font-bold text-gable-green">
                                            ${(customer.total_revenue / 100).toLocaleString(undefined, { minimumFractionDigits: 2 })}
                                        </td>
                                    </tr>
                                ))
                            )}
                        </tbody>
                    </table>
                </div>
            </CardContent>
        </Card>
    );
}
