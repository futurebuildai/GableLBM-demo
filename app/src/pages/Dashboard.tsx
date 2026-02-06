
export const Dashboard = () => {
    return (
        <div className="space-y-6">
            <h1 className="text-3xl font-bold tracking-tight text-white">Dashboard</h1>

            <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
                {/* Stats Cards */}
                <div className="p-6 rounded-lg bg-zinc-900 border border-white/10">
                    <h3 className="text-sm font-medium text-zinc-400">Total Revenue</h3>
                    <div className="mt-2 text-2xl font-mono font-bold text-white">$45,231.89</div>
                </div>
                <div className="p-6 rounded-lg bg-zinc-900 border border-white/10">
                    <h3 className="text-sm font-medium text-zinc-400">Active Orders</h3>
                    <div className="mt-2 text-2xl font-mono font-bold text-emerald-400">12</div>
                </div>
                <div className="p-6 rounded-lg bg-zinc-900 border border-white/10">
                    <h3 className="text-sm font-medium text-zinc-400">Pending Dispatch</h3>
                    <div className="mt-2 text-2xl font-mono font-bold text-blue-400">4</div>
                </div>
            </div>
        </div>
    );
};
