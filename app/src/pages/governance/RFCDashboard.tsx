import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { GovernanceService } from '../../services/governance.service';
import type { RFC } from '../../types/governance';

export function RFCDashboard() {
    const [rfcs, setRfcs] = useState<RFC[]>([]);
    const navigate = useNavigate();

    useEffect(() => {
        loadRFCs();
    }, []);

    const loadRFCs = async () => {
        try {
            const data = await GovernanceService.listRFCs();
            setRfcs(data);
        } catch (e) {
            console.error(e);
        }
    };

    const statusColor = (status: string) => {
        switch (status) {
            case 'approved': return 'text-green-500 bg-green-500/10 border-green-500/20';
            case 'rejected': return 'text-red-500 bg-red-500/10 border-red-500/20';
            case 'review': return 'text-yellow-500 bg-yellow-500/10 border-yellow-500/20';
            default: return 'text-slate-400 bg-slate-500/10 border-slate-500/20';
        }
    };

    return (
        <div className="h-full flex flex-col p-6 space-y-6">
            <div className="flex justify-between items-center">
                <div>
                    <h1 className="text-3xl font-bold text-white tracking-tight">Governance</h1>
                    <p className="text-slate-400 mt-1">Manage architectural decisions and RFCs.</p>
                </div>
                <button
                    onClick={() => navigate('/governance/new')}
                    className="industrial-button bg-[#00FFA3] text-black px-4 py-2 font-medium hover:shadow-[0_0_15px_rgba(0,255,163,0.3)] transition-all"
                >
                    Draft New RFC
                </button>
            </div>

            <div className="flex-1 overflow-auto border border-white/10 rounded-lg bg-[#0A0B10]">
                <table className="w-full text-left">
                    <thead className="bg-[#161821] text-xs uppercase text-slate-500 font-mono sticky top-0">
                        <tr>
                            <th className="px-6 py-3">Status</th>
                            <th className="px-6 py-3">Title</th>
                            <th className="px-6 py-3">Problem</th>
                            <th className="px-6 py-3 text-right">Created</th>
                        </tr>
                    </thead>
                    <tbody className="divide-y divide-white/5">
                        {rfcs.map((rfc) => (
                            <tr
                                key={rfc.id}
                                onClick={() => navigate(`/governance/${rfc.id}`)}
                                className="hover:bg-white/5 cursor-pointer transition-colors"
                            >
                                <td className="px-6 py-4">
                                    <span className={`px-2 py-1 rounded text-xs font-mono uppercase border ${statusColor(rfc.status)}`}>
                                        {rfc.status}
                                    </span>
                                </td>
                                <td className="px-6 py-4 text-white font-medium">{rfc.title}</td>
                                <td className="px-6 py-4 text-slate-400 truncate max-w-md">{rfc.problem_statement}</td>
                                <td className="px-6 py-4 text-right text-slate-500 font-mono text-sm">
                                    {new Date(rfc.created_at).toLocaleDateString()}
                                </td>
                            </tr>
                        ))}
                        {rfcs.length === 0 && (
                            <tr>
                                <td colSpan={4} className="px-6 py-12 text-center text-slate-500">
                                    No RFCs found. Create one to get started.
                                </td>
                            </tr>
                        )}
                    </tbody>
                </table>
            </div>
        </div>
    );
}
