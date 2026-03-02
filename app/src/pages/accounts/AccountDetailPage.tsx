import { useState, useEffect } from 'react';
import { useParams, Link } from 'react-router-dom';
import { CustomerService } from '../../services/CustomerService';
import { AccountService } from '../../services/AccountService';
import type { Customer } from '../../types/customer';
import type { AccountSummary, CustomerTransaction } from '../../types/account';
import { LoadingScreen } from '../../components/ui/LoadingScreen';
import { Card } from '../../components/ui/Card';
import { Button } from '../../components/ui/Button'; // Assuming Button component exists
import { ArrowLeft, CreditCard, Receipt, FileText, Activity, AlertCircle, Users, MessageSquare } from 'lucide-react';
import { cn } from '../../lib/utils';
import { ContactList } from './ContactList';
import { ActivityFeed } from './ActivityFeed';

export function AccountDetailPage() {
    const { id } = useParams<{ id: string }>();
    const [customer, setCustomer] = useState<Customer | null>(null);
    const [summary, setSummary] = useState<AccountSummary | null>(null);
    const [transactions, setTransactions] = useState<CustomerTransaction[]>([]);
    const [loading, setLoading] = useState(true);
    const [activeTab, setActiveTab] = useState<'ledger' | 'invoices' | 'payments' | 'contacts' | 'crm'>('ledger');

    useEffect(() => {
        if (id) {
            loadData(id);
        }
    }, [id]);

    async function loadData(customerId: string) {
        try {
            const [cust, summ, txns] = await Promise.all([
                CustomerService.getCustomer(customerId),
                AccountService.getAccountSummary(customerId),
                AccountService.getTransactions(customerId)
            ]);
            setCustomer(cust);
            setSummary(summ);
            setTransactions(txns);
        } catch (error) {
            console.error('Failed to load account data:', error);
            // Ideally use toast here, but for now just log to avoid build failure on missing Toast import
        } finally {
            setLoading(false);
        }
    }

    if (loading) return <LoadingScreen />;
    if (!customer || !summary) return <div className="p-8 text-center text-zinc-400">Account not found</div>;

    const availablePercentage = summary.credit_limit > 0
        ? (summary.available_credit / summary.credit_limit) * 100
        : 0;

    return (
        <div className="space-y-8 max-w-6xl mx-auto">
            {/* Header */}
            <div>
                <Link to="/accounts" className="inline-flex items-center text-sm text-zinc-400 hover:text-white mb-4 transition-colors">
                    <ArrowLeft size={16} className="mr-1" /> Back to Accounts
                </Link>
                <div className="flex items-start justify-between">
                    <div>
                        <h1 className="text-3xl font-bold text-white">{customer.name}</h1>
                        <div className="flex items-center gap-3 mt-2 text-zinc-400 text-sm">
                            <span className="font-mono bg-white/5 px-2 py-0.5 rounded border border-white/5">#{customer.account_number}</span>
                            <span>{customer.email}</span>
                            <span>•</span>
                            <span>{customer.phone}</span>
                        </div>
                    </div>
                    <div className="flex gap-2">
                        <Button variant="outline" size="sm">Edit Profile</Button>
                        <Button variant="default" size="sm">New Transaction</Button>
                    </div>
                </div>
            </div>

            {/* Financial Overview Cards */}
            <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                <Card className="p-5 border-l-4 border-l-orange-500 bg-gradient-to-br from-slate-steel to-transparent">
                    <div className="flex items-center gap-3 mb-2 text-zinc-400 text-sm font-medium uppercase tracking-wide">
                        <Receipt size={16} className="text-orange-400" />
                        Balance Due
                    </div>
                    <div className="text-3xl font-mono font-bold text-white">
                        ${(summary.balance_due / 100).toLocaleString('en-US', { minimumFractionDigits: 2 })}
                    </div>
                    <div className="mt-2 text-xs text-zinc-500">Current outstanding balance</div>
                </Card>

                <Card className="p-5 border-l-4 border-l-emerald-500 bg-gradient-to-br from-slate-steel to-transparent">
                    <div className="flex items-center gap-3 mb-2 text-zinc-400 text-sm font-medium uppercase tracking-wide">
                        <CreditCard size={16} className="text-emerald-400" />
                        Available Credit
                    </div>
                    <div className={cn("text-3xl font-mono font-bold", availablePercentage < 20 ? "text-red-400" : "text-white")}>
                        ${(summary.available_credit / 100).toLocaleString('en-US', { minimumFractionDigits: 2 })}
                    </div>
                    <div className="mt-2 w-full bg-white/10 h-1.5 rounded-full overflow-hidden">
                        <div
                            className={cn("h-full rounded-full transition-all duration-500", availablePercentage < 20 ? "bg-red-500" : "bg-emerald-500")}
                            style={{ width: `${Math.min(availablePercentage, 100)}%` }}
                        />
                    </div>
                </Card>

                <Card className="p-5 border-l-4 border-l-blue-500 bg-gradient-to-br from-slate-steel to-transparent">
                    <div className="flex items-center gap-3 mb-2 text-zinc-400 text-sm font-medium uppercase tracking-wide">
                        <Activity size={16} className="text-blue-400" />
                        Credit Limit
                    </div>
                    <div className="text-3xl font-mono font-bold text-white">
                        ${(summary.credit_limit / 100).toLocaleString('en-US', { minimumFractionDigits: 2 })}
                    </div>
                    <div className="mt-2 text-xs text-zinc-500">Total approved credit line</div>
                </Card>
            </div>

            {/* Tabs & Content */}
            <div className="space-y-4">
                <div className="flex items-center gap-1 border-b border-white/10 pb-1">
                    <Tab active={activeTab === 'ledger'} onClick={() => setActiveTab('ledger')} label="Activity Ledger" icon={<Activity size={16} />} />
                    <Tab active={activeTab === 'invoices'} onClick={() => setActiveTab('invoices')} label="Invoices" icon={<FileText size={16} />} />
                    <Tab active={activeTab === 'payments'} onClick={() => setActiveTab('payments')} label="Payments" icon={<CreditCard size={16} />} />
                    <Tab active={activeTab === 'contacts'} onClick={() => setActiveTab('contacts')} label="Contacts" icon={<Users size={16} />} />
                    <Tab active={activeTab === 'crm'} onClick={() => setActiveTab('crm')} label="CRM Activity" icon={<MessageSquare size={16} />} />
                </div>

                <div className="min-h-[400px]">
                    {activeTab === 'ledger' && (
                        <div className="border border-white/5 rounded-lg overflow-hidden bg-slate-steel/20">
                            <table className="w-full text-sm">
                                <thead className="bg-white/5 text-zinc-400 font-medium border-b border-white/5">
                                    <tr>
                                        <th className="px-4 py-3 text-left">Date</th>
                                        <th className="px-4 py-3 text-left">Type</th>
                                        <th className="px-4 py-3 text-left">Description</th>
                                        <th className="px-4 py-3 text-right">Amount</th>
                                        <th className="px-4 py-3 text-right">Running Balance</th>
                                    </tr>
                                </thead>
                                <tbody className="divide-y divide-white/5">
                                    {transactions.length === 0 ? (
                                        <tr>
                                            <td colSpan={5} className="px-4 py-8 text-center text-zinc-500">No transactions found</td>
                                        </tr>
                                    ) : (
                                        transactions.map(txn => (
                                            <tr key={txn.id} className="group hover:bg-white/5 transition-colors">
                                                <td className="px-4 py-3 font-mono text-zinc-400">
                                                    {new Date(txn.created_at).toLocaleDateString()}
                                                </td>
                                                <td className="px-4 py-3">
                                                    <span className={cn(
                                                        "px-2 py-0.5 rounded text-xs font-bold uppercase tracking-wider",
                                                        txn.type === 'INVOICE' ? "bg-orange-500/10 text-orange-400" :
                                                            txn.type === 'PAYMENT' ? "bg-emerald-500/10 text-emerald-400" :
                                                                "bg-blue-500/10 text-blue-400"
                                                    )}>
                                                        {txn.type}
                                                    </span>
                                                </td>
                                                <td className="px-4 py-3 text-white">{txn.description}</td>
                                                <td className={cn("px-4 py-3 text-right font-mono font-medium", txn.amount > 0 ? "text-white" : "text-emerald-400")}>
                                                    {txn.amount > 0 ? '+' : ''}{(txn.amount / 100).toFixed(2)}
                                                </td>
                                                <td className="px-4 py-3 text-right font-mono text-zinc-300">
                                                    {(txn.balance_after / 100).toFixed(2)}
                                                </td>
                                            </tr>
                                        ))
                                    )}
                                </tbody>
                            </table>
                        </div>
                    )}
                    {activeTab === 'contacts' && (
                        <div className="mt-4">
                            <ContactList customerId={customer.id} />
                        </div>
                    )}
                    {activeTab === 'crm' && (
                        <div className="mt-4">
                            <ActivityFeed customerId={customer.id} />
                        </div>
                    )}
                    {(activeTab !== 'ledger' && activeTab !== 'contacts' && activeTab !== 'crm') && (
                        <div className="flex flex-col items-center justify-center h-64 border border-dashed border-white/10 rounded-lg text-zinc-500">
                            <AlertCircle size={32} className="mb-2 opacity-50" />
                            <p>This view is under construction.</p>
                        </div>
                    )}
                </div>
            </div>
        </div>
    );
}

function Tab({ active, onClick, label, icon }: { active: boolean, onClick: () => void, label: string, icon: React.ReactNode }) {
    return (
        <button
            onClick={onClick}
            className={cn(
                "flex items-center gap-2 px-4 py-2.5 text-sm font-medium border-b-2 transition-colors",
                active
                    ? "text-gable-green border-gable-green bg-gable-green/5"
                    : "text-zinc-400 border-transparent hover:text-white hover:border-white/20"
            )}
        >
            {icon}
            {label}
        </button>
    )
}
