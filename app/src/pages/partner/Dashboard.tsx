import { Card, CardHeader, CardTitle, CardContent } from '../../components/ui/Card';
import { PageTransition } from '../../components/ui/PageTransition';
import { DollarSign, Package, Clock, ArrowRight } from 'lucide-react';
import { Link } from 'react-router-dom';

interface StatCardProps {
    title: string;
    value: string;
    icon: React.ElementType;
    color: string;
}

const StatCard = ({ title, value, icon: Icon, color }: StatCardProps) => (
    <Card variant="glass" className="group hover:-translate-y-1 transition-transform duration-300">
        <CardContent className="p-6">
            <div className="flex justify-between items-start mb-4">
                <div className={`p-3 rounded-lg bg-${color}-500/10 border border-${color}-500/20`}>
                    <Icon className={`w-6 h-6 text-${color}-500`} />
                </div>
            </div>
            <div>
                <p className="text-zinc-400 text-sm font-medium mb-1">{title}</p>
                <h3 className="text-3xl font-bold text-white font-mono tracking-tight">{value}</h3>
            </div>
        </CardContent>
    </Card>
);

export const PartnerDashboard = () => {
    return (
        <PageTransition>
            <div className="mb-8">
                <h1 className="text-display-large text-white">Welcome back, John</h1>
                <p className="text-zinc-400 mt-2 text-lg">Here's what's happening with your projects today.</p>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
                <StatCard
                    title="Active Projects"
                    value="12"
                    icon={Package}
                    color="amber"
                />
                <StatCard
                    title="Open Invoices"
                    value="$4,250.00"
                    icon={DollarSign}
                    color="emerald"
                />
                <StatCard
                    title="Pending Quotes"
                    value="3"
                    icon={Clock}
                    color="blue"
                />
            </div>

            <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
                <Card variant="glass" className="h-full">
                    <CardHeader className="flex flex-row items-center justify-between pb-2">
                        <CardTitle>Recent Orders</CardTitle>
                        <Link to="/partner/orders" className="text-sm text-gable-green hover:underline flex items-center gap-1">
                            View All <ArrowRight className="w-4 h-4" />
                        </Link>
                    </CardHeader>
                    <CardContent>
                        <div className="space-y-4">
                            {[1, 2, 3].map((i) => (
                                <div key={i} className="flex justify-between items-center p-3 rounded-lg hover:bg-white/5 transition-colors border border-transparent hover:border-white/5">
                                    <div>
                                        <div className="font-medium text-white">Order #100{i}</div>
                                        <div className="text-sm text-zinc-500">Project: Maple Ave Renovation</div>
                                    </div>
                                    <div className="text-right">
                                        <div className="font-mono text-zinc-300">$1,2{i}0.00</div>
                                        <span className="inline-block px-2 py-0.5 rounded text-[10px] bg-blue-500/10 text-blue-400 border border-blue-500/20 uppercase tracking-wider font-semibold">
                                            Processing
                                        </span>
                                    </div>
                                </div>
                            ))}
                        </div>
                    </CardContent>
                </Card>

                <Card variant="glass" className="h-full">
                    <CardHeader className="flex flex-row items-center justify-between pb-2">
                        <CardTitle>Approvals Needed</CardTitle>
                    </CardHeader>
                    <CardContent>
                        <div className="space-y-4">
                            <div className="p-4 rounded-lg bg-amber-500/5 border border-amber-500/10 hover:bg-amber-500/10 transition-colors">
                                <div className="flex justify-between items-start mb-2">
                                    <h4 className="font-medium text-white">Quote #Q-2024-001</h4>
                                    <span className="text-xs text-amber-500 font-medium bg-amber-500/10 px-2 py-1 rounded">Pending Approval</span>
                                </div>
                                <p className="text-sm text-zinc-400 mb-3">Lumber package for 123 Oak St framing.</p>
                                <div className="flex gap-2">
                                    <button className="text-xs bg-gable-green text-black font-semibold px-3 py-1.5 rounded hover:bg-emerald-400 transition-colors">
                                        Approve Quote
                                    </button>
                                    <button className="text-xs bg-white/5 text-white px-3 py-1.5 rounded hover:bg-white/10 transition-colors border border-white/10">
                                        Review Details
                                    </button>
                                </div>
                            </div>
                        </div>
                    </CardContent>
                </Card>
            </div>
        </PageTransition>
    );
};
