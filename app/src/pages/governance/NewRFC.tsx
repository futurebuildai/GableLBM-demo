import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { GovernanceService } from '../../services/governance.service';
import { useToast } from '../../components/ui/ToastContext';

export function NewRFC() {
    const navigate = useNavigate();
    const { showToast } = useToast();
    const [loading, setLoading] = useState(false);
    const [formData, setFormData] = useState({
        title: '',
        problem_statement: '',
        proposed_solution: ''
    });

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setLoading(true);
        try {
            await GovernanceService.createRFC(formData);
            navigate('/governance');
        } catch (err) {
            console.error(err);
            showToast('Failed to create RFC', 'error');
        } finally {
            setLoading(false);
        }
    };

    return (
        <div className="p-6 max-w-4xl mx-auto">
            <div className="mb-8">
                <h1 className="text-3xl font-bold text-white mb-2">New RFC Draft</h1>
                <p className="text-slate-400">Describe the problem. The AI Governance Layer will generate the structure.</p>
            </div>

            <form onSubmit={handleSubmit} className="space-y-6">
                <div className="space-y-2">
                    <label className="text-sm font-medium text-slate-300">Title</label>
                    <input
                        type="text"
                        required
                        className="w-full bg-[#0A0B10] border border-white/20 rounded p-3 text-white focus:border-[#00FFA3] outline-none"
                        placeholder="e.g. Implement Zero Trust Auth"
                        value={formData.title}
                        onChange={(e) => setFormData({ ...formData, title: e.target.value })}
                    />
                </div>

                <div className="space-y-2">
                    <label className="text-sm font-medium text-slate-300">Problem Statement</label>
                    <textarea
                        required
                        rows={4}
                        className="w-full bg-[#0A0B10] border border-white/20 rounded p-3 text-white focus:border-[#00FFA3] outline-none"
                        placeholder="What is broken or missing? Why is this important now?"
                        value={formData.problem_statement}
                        onChange={(e) => setFormData({ ...formData, problem_statement: e.target.value })}
                    />
                </div>

                <div className="space-y-2">
                    <label className="text-sm font-medium text-slate-300">Proposed Solution (High Level)</label>
                    <textarea
                        required
                        rows={4}
                        className="w-full bg-[#0A0B10] border border-white/20 rounded p-3 text-white focus:border-[#00FFA3] outline-none"
                        placeholder="Briefly describe the approach..."
                        value={formData.proposed_solution}
                        onChange={(e) => setFormData({ ...formData, proposed_solution: e.target.value })}
                    />
                </div>

                <div className="pt-4 flex items-center space-x-4">
                    <button
                        type="button"
                        onClick={() => navigate('/governance')}
                        className="px-6 py-3 border border-white/10 text-white rounded hover:bg-white/5"
                    >
                        Cancel
                    </button>
                    <button
                        type="submit"
                        disabled={loading}
                        className="px-6 py-3 bg-[#00FFA3] text-black font-bold rounded hover:shadow-[0_0_15px_rgba(0,255,163,0.3)] disabled:opacity-50"
                    >
                        {loading ? 'Generating...' : 'Generate Draft'}
                    </button>
                </div>
            </form>
        </div>
    );
}
