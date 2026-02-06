import { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { GovernanceService } from '../../services/governance.service';
import type { RFC } from '../../types/governance';

export function RFCDetail() {
    const { id } = useParams<{ id: string }>();
    const navigate = useNavigate();
    const [rfc, setRfc] = useState<RFC | null>(null);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        if (id) loadRFC(id);
    }, [id]);

    const loadRFC = async (rfcId: string) => {
        try {
            const data = await GovernanceService.getRFC(rfcId);
            setRfc(data);
        } catch (e) {
            console.error(e);
            alert('Failed to load RFC');
        } finally {
            setLoading(false);
        }
    };

    if (loading) return <div className="p-8 text-slate-400">Loading protocol...</div>;
    if (!rfc) return <div className="p-8 text-red-500">RFC Not Found</div>;

    return (
        <div className="flex h-full">
            {/* Sidebar / Meta */}
            <div className="w-80 border-r border-white/10 p-6 bg-[#0A0B10] flex flex-col space-y-6">
                <div>
                    <button
                        onClick={() => navigate('/governance')}
                        className="text-slate-400 hover:text-white flex items-center mb-6 text-sm"
                    >
                        ← Back to Governance
                    </button>
                    <h1 className="text-xl font-bold text-white mb-2">{rfc.title}</h1>
                    <div className="inline-block px-2 py-1 rounded text-xs font-mono uppercase bg-slate-800 text-slate-300 border border-slate-700">
                        {rfc.status}
                    </div>
                </div>

                <div className="space-y-4">
                    <div>
                        <h3 className="text-xs uppercase text-slate-500 font-bold mb-1">Author</h3>
                        <p className="text-slate-300 font-mono text-sm">Owner Bob</p>
                    </div>
                    <div>
                        <h3 className="text-xs uppercase text-slate-500 font-bold mb-1">Created</h3>
                        <p className="text-slate-300 font-mono text-sm">{new Date(rfc.created_at).toLocaleString()}</p>
                    </div>
                    <div>
                        <h3 className="text-xs uppercase text-slate-500 font-bold mb-1">Last Updated</h3>
                        <p className="text-slate-300 font-mono text-sm">{new Date(rfc.updated_at).toLocaleString()}</p>
                    </div>
                </div>

                <div className="pt-6 border-t border-white/10">
                    <button className="w-full industrial-button bg-[#00FFA3] text-black px-4 py-2 font-medium mb-2">
                        Edit RFC
                    </button>
                    <button className="w-full border border-white/10 text-white px-4 py-2 font-medium hover:bg-white/5">
                        Export to PDF
                    </button>
                </div>
            </div>

            {/* Main Content (Document) */}
            <div className="flex-1 overflow-auto bg-[#161821] p-8">
                <div className="max-w-4xl mx-auto bg-[#0A0B10] border border-white/5 p-12 min-h-screen shadow-2xl">
                    {/* Simple Markdown-like rendering */}
                    <pre className="font-mono text-slate-300 whitespace-pre-wrap leading-relaxed">
                        {rfc.content}
                    </pre>
                </div>
            </div>
        </div>
    );
}
