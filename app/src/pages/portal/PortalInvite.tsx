import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { Card, CardContent } from '../../components/ui/Card';
import { Mail, Shield, UserPlus, ArrowLeft, Loader2 } from 'lucide-react';
import { PortalService } from '../../services/PortalService';
import { useToast } from '../../components/ui/ToastContext';

export const PortalInvite = () => {
    const navigate = useNavigate();
    const { showToast } = useToast();

    const [email, setEmail] = useState('');
    const [role, setRole] = useState('Buyer');
    const [submitting, setSubmitting] = useState(false);

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        if (!email) return;

        setSubmitting(true);
        try {
            await PortalService.inviteUser({ email, role });
            showToast(`Invitation sent to ${email}`, 'success');
            navigate('/portal/team');
        } catch (err) {
            showToast(err instanceof Error ? err.message : 'Failed to send invite', 'error');
            setSubmitting(false);
        }
    };

    return (
        <div className="max-w-2xl mx-auto">
            <div className="mb-6 flex items-center gap-4">
                <button
                    onClick={() => navigate('/portal/team')}
                    className="p-2 text-zinc-400 hover:text-white hover:bg-white/5 rounded-lg transition-colors"
                >
                    <ArrowLeft size={20} />
                </button>
                <div>
                    <h1 className="text-2xl font-bold text-white">Invite Team Member</h1>
                    <p className="text-zinc-400 text-sm mt-1">Send an invitation to join your company's portal</p>
                </div>
            </div>

            <Card variant="glass" className="p-1">
                <CardContent className="p-6">
                    <form onSubmit={handleSubmit} className="space-y-6">
                        <div>
                            <label className="block text-sm font-medium text-zinc-300 mb-2">Email Address</label>
                            <div className="relative">
                                <Mail className="absolute left-3 top-1/2 -translate-y-1/2 text-zinc-500" size={18} />
                                <input
                                    type="email"
                                    required
                                    value={email}
                                    onChange={(e) => setEmail(e.target.value)}
                                    placeholder="colleague@company.com"
                                    className="w-full bg-black/20 border border-white/10 rounded-lg py-2.5 pl-10 pr-4 text-white placeholder-zinc-600 focus:outline-none focus:ring-1 focus:ring-gable-green transition-all"
                                />
                            </div>
                        </div>

                        <div>
                            <label className="block text-sm font-medium text-zinc-300 mb-2">Role & Permissions</label>
                            <div className="grid grid-cols-1 md:grid-cols-3 gap-3">
                                {['View-Only', 'Buyer', 'Admin'].map((r) => (
                                    <label
                                        key={r}
                                        className={`relative flex flex-col p-4 cursor-pointer rounded-xl border transition-all ${role === r
                                                ? 'bg-[#00FFA3]/10 border-[#00FFA3] shadow-[0_0_15px_rgba(0,255,163,0.1)]'
                                                : 'bg-black/20 border-white/5 hover:bg-white/5 hover:border-white/20'
                                            }`}
                                    >
                                        <input
                                            type="radio"
                                            name="role"
                                            value={r}
                                            checked={role === r}
                                            onChange={() => setRole(r)}
                                            className="sr-only"
                                        />
                                        <div className="flex items-center justify-between mb-2">
                                            <span className={`font-semibold ${role === r ? 'text-[#00FFA3]' : 'text-white'}`}>{r}</span>
                                            {role === r && <Shield size={16} className="text-[#00FFA3]" />}
                                        </div>
                                        <span className="text-xs text-zinc-500 leading-relaxed">
                                            {r === 'View-Only' && 'Can view orders, invoices, and catalog. Cannot place orders.'}
                                            {r === 'Buyer' && 'Can view items and place new orders on account.'}
                                            {r === 'Admin' && 'Full access. Can place orders and manage company team members.'}
                                        </span>
                                    </label>
                                ))}
                            </div>
                        </div>

                        <div className="pt-4 border-t border-white/5 flex items-center justify-end gap-3">
                            <button
                                type="button"
                                onClick={() => navigate('/portal/team')}
                                className="px-4 py-2 text-sm font-medium text-zinc-400 hover:text-white hover:bg-white/5 rounded-lg transition-colors"
                            >
                                Cancel
                            </button>
                            <button
                                type="submit"
                                disabled={submitting || !email}
                                className="flex items-center gap-2 px-6 py-2 bg-gable-green text-black font-semibold rounded-lg hover:bg-emerald-400 transition-all shadow-[0_0_15px_rgba(0,255,163,0.3)] disabled:opacity-50 disabled:shadow-none disabled:cursor-not-allowed"
                            >
                                {submitting ? <Loader2 size={18} className="animate-spin" /> : <UserPlus size={18} />}
                                Send Invitation
                            </button>
                        </div>
                    </form>
                </CardContent>
            </Card>
        </div>
    );
};
