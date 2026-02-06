
export function PartnerDashboard() {
    return (
        <div className="space-y-6">
            <h1 className="text-2xl font-bold text-white">Dashboard</h1>

            <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
                <div className="bg-slate-steel p-6 rounded-lg border border-white/10">
                    <div className="text-muted-foreground text-sm">Current Balance</div>
                    <div className="text-3xl font-bold text-white mt-1">$12,450.00</div>
                    <div className="text-xs text-emerald-400 mt-2">Within Credit Limit</div>
                </div>

                <div className="bg-slate-steel p-6 rounded-lg border border-white/10">
                    <div className="text-muted-foreground text-sm">Active Jobs</div>
                    <div className="text-3xl font-bold text-white mt-1">4</div>
                    <div className="text-xs text-muted-foreground mt-2">2 Pending Quotes</div>
                </div>

                <div className="bg-slate-steel p-6 rounded-lg border border-white/10">
                    <div className="text-muted-foreground text-sm">Open Invoices</div>
                    <div className="text-3xl font-bold text-white mt-1">3</div>
                    <div className="text-xs text-amber-500 mt-2">1 Overdue</div>
                </div>
            </div>

            {/* Recent Activity Mock */}
            <div className="bg-slate-steel rounded-lg border border-white/10">
                <div className="p-4 border-b border-white/10">
                    <h2 className="font-semibold text-white">Recent Activity</h2>
                </div>
                <div className="p-4 text-muted-foreground text-sm">
                    No recent activity to show.
                </div>
            </div>
        </div>
    );
}
