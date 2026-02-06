
import { Link } from 'react-router-dom';

export function ProjectList() {
    // Mock Data
    const projects = [
        { id: '1', name: 'Smith Residence', active: true, quotes: 2, orders: 5 },
        { id: '2', name: 'Downtown Lofts', active: true, quotes: 0, orders: 12 },
        { id: '3', name: 'Lake House Reno', active: false, quotes: 0, orders: 8 },
    ];

    return (
        <div className="space-y-6">
            <div className="flex items-center justify-between">
                <h1 className="text-2xl font-bold text-white">Projects</h1>
            </div>

            <div className="bg-slate-steel rounded-lg border border-white/10 overflow-hidden">
                <table className="w-full text-left text-sm text-muted-foreground">
                    <thead className="bg-white/5 text-white font-medium uppercase text-xs">
                        <tr>
                            <th className="px-4 py-3">Project Name</th>
                            <th className="px-4 py-3">Status</th>
                            <th className="px-4 py-3">Open Quotes</th>
                            <th className="px-4 py-3">Orders</th>
                            <th className="px-4 py-3">Action</th>
                        </tr>
                    </thead>
                    <tbody className="divide-y divide-white/5">
                        {projects.map((p) => (
                            <tr key={p.id} className="hover:bg-white/5 transition-colors">
                                <td className="px-4 py-3 font-medium text-white">{p.name}</td>
                                <td className="px-4 py-3">
                                    <span className={`px-2 py-0.5 rounded text-xs ${p.active ? 'bg-emerald-500/20 text-emerald-400' : 'bg-white/10 text-muted-foreground'}`}>
                                        {p.active ? 'Active' : 'Archived'}
                                    </span>
                                </td>
                                <td className="px-4 py-3 text-white">{p.quotes}</td>
                                <td className="px-4 py-3">{p.orders}</td>
                                <td className="px-4 py-3">
                                    <Link to={`/partner/projects/${p.id}`} className="text-blue-400 hover:underline">View</Link>
                                </td>
                            </tr>
                        ))}
                    </tbody>
                </table>
            </div>
        </div>
    );
}
