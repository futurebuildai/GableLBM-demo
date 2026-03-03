import { useState, useEffect, useCallback } from 'react';
import { useNavigate } from 'react-router-dom';
import { Card } from '../../components/ui/Card';
import { Users, UserPlus, RefreshCw, AlertTriangle, Shield, CheckCircle2, XCircle, Clock } from 'lucide-react';
import { PortalService } from '../../services/PortalService';
import type { PortalUser, PortalInvite } from '../../types/portal';
import { useToast } from '../../components/ui/ToastContext';

export const PortalTeam = () => {
    const navigate = useNavigate();
    const { showToast } = useToast();

    const [users, setUsers] = useState<PortalUser[]>([]);
    const [invites, setInvites] = useState<PortalInvite[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState('');

    const fetchData = useCallback(async () => {
        setLoading(true);
        setError('');
        try {
            const [usersData, invitesData] = await Promise.all([
                PortalService.getUsers(),
                PortalService.getInvites()
            ]);
            setUsers(usersData);
            setInvites(invitesData);
        } catch (err) {
            setError(err instanceof Error ? err.message : 'Failed to load team data');
        } finally {
            setLoading(false);
        }
    }, []);

    useEffect(() => { fetchData(); }, [fetchData]);

    const handleRoleChange = async (userId: string, newRole: string) => {
        try {
            await PortalService.updateUserRole(userId, { role: newRole });
            showToast('Role updated successfully', 'success');
            fetchData();
        } catch (err) {
            showToast(err instanceof Error ? err.message : 'Failed to update role', 'error');
        }
    };

    const handleStatusChange = async (userId: string, newStatus: string) => {
        try {
            await PortalService.updateUserStatus(userId, { status: newStatus });
            showToast(`User marked as ${newStatus}`, 'success');
            fetchData();
        } catch (err) {
            showToast(err instanceof Error ? err.message : 'Failed to update status', 'error');
        }
    };

    if (loading) {
        return (
            <div className="space-y-4">
                <div className="h-10 w-1/4 bg-white/5 rounded-lg animate-pulse mb-6" />
                {[1, 2, 3].map(i => (
                    <div key={i} className="h-16 bg-white/5 rounded-2xl animate-pulse" />
                ))}
            </div>
        );
    }

    if (error) {
        return (
            <div className="flex flex-col items-center justify-center h-64 text-center">
                <AlertTriangle className="w-12 h-12 text-amber-500 mb-4" />
                <p className="text-zinc-400 mb-4">{error}</p>
                <button
                    onClick={fetchData}
                    className="flex items-center gap-2 px-4 py-2 rounded-lg bg-white/5 border border-white/10 text-white hover:bg-white/10 transition-colors"
                >
                    <RefreshCw size={16} /> Retry
                </button>
            </div>
        );
    }

    return (
        <div>
            <div className="mb-6 flex items-center justify-between">
                <div>
                    <h1 className="text-2xl font-bold text-white">Team Management</h1>
                    <p className="text-zinc-400 text-sm mt-1">Manage portal access for your company</p>
                </div>
                <button
                    onClick={() => navigate('/portal/team/invite')}
                    className="flex items-center gap-2 px-4 py-2 bg-gable-green text-black font-semibold rounded-lg hover:bg-emerald-400 transition-colors shadow-[0_0_15px_rgba(0,255,163,0.3)]"
                >
                    <UserPlus size={18} /> Invite Member
                </button>
            </div>

            <div className="space-y-8">
                {/* Active Members */}
                <div>
                    <h2 className="text-lg font-medium text-white mb-4 flex items-center gap-2">
                        <Users size={18} className="text-zinc-400" /> Current Members
                    </h2>
                    <Card variant="glass" noPadding>
                        <div className="divide-y divide-white/5">
                            {users.map(user => (
                                <div key={user.id} className="p-4 flex flex-col md:flex-row md:items-center justify-between gap-4 hover:bg-white/[0.02] transition-colors">
                                    <div className="flex items-center gap-4">
                                        <div className="w-10 h-10 rounded-full bg-white/5 border border-white/10 flex items-center justify-center text-white font-medium">
                                            {user.name.charAt(0).toUpperCase()}
                                        </div>
                                        <div>
                                            <div className="font-medium text-white flex items-center gap-2">
                                                {user.name}
                                                {user.status === 'Inactive' && (
                                                    <span className="text-[10px] uppercase px-1.5 py-0.5 rounded bg-red-500/10 text-red-400 border border-red-500/20 font-semibold tracking-wider">
                                                        Inactive
                                                    </span>
                                                )}
                                            </div>
                                            <div className="text-sm text-zinc-500">{user.email}</div>
                                        </div>
                                    </div>
                                    <div className="flex items-center gap-3">
                                        <div className="relative">
                                            <select
                                                value={user.role}
                                                onChange={(e) => handleRoleChange(user.id, e.target.value)}
                                                className="appearance-none bg-white/5 border border-white/10 rounded-lg px-3 py-1.5 pr-8 text-sm text-white focus:outline-none focus:ring-1 focus:ring-gable-green cursor-pointer"
                                            >
                                                <option value="View-Only" className="bg-zinc-900">View-Only</option>
                                                <option value="Buyer" className="bg-zinc-900">Buyer</option>
                                                <option value="Admin" className="bg-zinc-900">Admin</option>
                                            </select>
                                            <Shield size={14} className="absolute right-2 top-1/2 -translate-y-1/2 text-zinc-500 pointer-events-none" />
                                        </div>

                                        {user.status === 'Active' ? (
                                            <button
                                                onClick={() => handleStatusChange(user.id, 'Inactive')}
                                                className="p-2 text-zinc-500 hover:text-red-400 hover:bg-white/5 rounded-lg transition-colors"
                                                title="Deactivate User"
                                            >
                                                <XCircle size={18} />
                                            </button>
                                        ) : (
                                            <button
                                                onClick={() => handleStatusChange(user.id, 'Active')}
                                                className="p-2 text-zinc-500 hover:text-emerald-400 hover:bg-white/5 rounded-lg transition-colors"
                                                title="Reactivate User"
                                            >
                                                <CheckCircle2 size={18} />
                                            </button>
                                        )}
                                    </div>
                                </div>
                            ))}
                        </div>
                    </Card>
                </div>

                {/* Pending Invites */}
                {invites.length > 0 && (
                    <div>
                        <h2 className="text-lg font-medium text-white mb-4 flex items-center gap-2">
                            <Clock size={18} className="text-zinc-400" /> Pending Invites
                        </h2>
                        <Card variant="glass" noPadding>
                            <div className="divide-y divide-white/5">
                                {invites.map(invite => (
                                    <div key={invite.id} className="p-4 flex flex-col md:flex-row md:items-center justify-between gap-4">
                                        <div className="flex items-center gap-4">
                                            <div className="w-10 h-10 rounded-full bg-blue-500/10 border border-blue-500/20 flex items-center justify-center text-blue-400">
                                                <UserPlus size={18} />
                                            </div>
                                            <div>
                                                <div className="font-medium text-white">{invite.email}</div>
                                                <div className="text-sm text-zinc-500">
                                                    Invited on {new Date(invite.created_at).toLocaleDateString()} · Expires {new Date(invite.expires_at).toLocaleDateString()}
                                                </div>
                                            </div>
                                        </div>
                                        <div className="flex items-center gap-3">
                                            <span className="px-2 py-1 rounded bg-white/5 text-zinc-300 text-xs font-medium border border-white/10">
                                                {invite.role}
                                            </span>
                                            <span className="text-[10px] uppercase px-1.5 py-0.5 rounded bg-blue-500/10 text-blue-400 border border-blue-500/20 font-semibold tracking-wider">
                                                Pending
                                            </span>
                                        </div>
                                    </div>
                                ))}
                            </div>
                        </Card>
                    </div>
                )}
            </div>
        </div>
    );
};
