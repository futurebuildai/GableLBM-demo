import { useState, useEffect, useCallback } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { GovernanceService } from '../../services/governance.service';
import type { RFC } from '../../types/governance';
import { PageTransition } from '../../components/ui/PageTransition';
import { Button } from '../../components/ui/Button';
import { Card, CardContent } from '../../components/ui/Card';
import { ArrowLeft, Edit2, FileText, Calendar, User, Clock, CheckCircle, AlertCircle } from 'lucide-react';
import { useToast } from '../../components/ui/ToastContext';

export function RFCDetail() {
    const { id } = useParams<{ id: string }>();
    const navigate = useNavigate();
    const { showToast } = useToast();
    const [rfc, setRfc] = useState<RFC | null>(null);
    const [loading, setLoading] = useState(true);

    const loadRFC = useCallback(async (rfcId: string) => {
        try {
            const data = await GovernanceService.getRFC(rfcId);
            setRfc(data);
        } catch (e) {
            console.error(e);
            showToast('Failed to load RFC', 'error');
        } finally {
            setLoading(false);
        }
    }, [showToast]);

    useEffect(() => {
        if (id) loadRFC(id);
    }, [id, loadRFC]);

    if (loading) {
        return (
            <div className="p-12 flex justify-center">
                <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-gable-green"></div>
            </div>
        );
    }

    if (!rfc) return <div className="p-8 text-rose-500">RFC Not Found</div>;

    const statusConfig = (status: string) => {
        switch (status.toLowerCase()) {
            case 'approved': return { color: 'text-emerald-400', bg: 'bg-emerald-500/10', border: 'border-emerald-500/20', icon: CheckCircle };
            case 'rejected': return { color: 'text-rose-400', bg: 'bg-rose-500/10', border: 'border-rose-500/20', icon: AlertCircle };
            case 'review': return { color: 'text-amber-400', bg: 'bg-amber-500/10', border: 'border-amber-500/20', icon: Clock };
            default: return { color: 'text-zinc-400', bg: 'bg-zinc-500/10', border: 'border-zinc-500/20', icon: FileText };
        }
    };

    const status = statusConfig(rfc.status);
    const StatusIcon = status.icon;

    return (
        <PageTransition>
            <div className="flex flex-col lg:flex-row h-[calc(100vh-6rem)] gap-6">
                {/* Sidebar / Meta */}
                <div className="lg:w-80 shrink-0">
                    <Card variant="glass" className="h-full bg-slate-steel/50 border-r border-white/5 lg:border-none">
                        <CardContent className="p-6 flex flex-col h-full">
                            <button
                                onClick={() => navigate('/governance')}
                                className="text-zinc-400 hover:text-white flex items-center mb-8 text-sm transition-colors group"
                            >
                                <ArrowLeft className="w-4 h-4 mr-2 group-hover:-translate-x-1 transition-transform" />
                                Back to Governance
                            </button>

                            <h1 className="text-xl font-bold text-white mb-4 leading-tight">{rfc.title}</h1>

                            <div className={`inline-flex items-center gap-2 px-3 py-1.5 rounded-lg text-sm font-medium border w-fit mb-8 ${status.bg} ${status.color} ${status.border}`}>
                                <StatusIcon className="w-4 h-4" />
                                <span className="uppercase tracking-wide text-xs">{rfc.status}</span>
                            </div>

                            <div className="space-y-6 flex-1">
                                <div className="flex items-start gap-3">
                                    <User className="w-5 h-5 text-zinc-500 mt-0.5" />
                                    <div>
                                        <h3 className="text-xs uppercase text-zinc-500 font-bold mb-0.5">Author</h3>
                                        <p className="text-zinc-300 text-sm">Owner Bob</p>
                                    </div>
                                </div>
                                <div className="flex items-start gap-3">
                                    <Calendar className="w-5 h-5 text-zinc-500 mt-0.5" />
                                    <div>
                                        <h3 className="text-xs uppercase text-zinc-500 font-bold mb-0.5">Created</h3>
                                        <p className="text-zinc-300 text-sm font-mono">{new Date(rfc.created_at).toLocaleDateString()}</p>
                                    </div>
                                </div>
                                <div className="flex items-start gap-3">
                                    <Clock className="w-5 h-5 text-zinc-500 mt-0.5" />
                                    <div>
                                        <h3 className="text-xs uppercase text-zinc-500 font-bold mb-0.5">Last Updated</h3>
                                        <p className="text-zinc-300 text-sm font-mono">{new Date(rfc.updated_at).toLocaleDateString()}</p>
                                    </div>
                                </div>
                            </div>

                            <div className="pt-6 border-t border-white/5 space-y-3 mt-auto">
                                <Button className="w-full justify-center">
                                    <Edit2 className="w-4 h-4 mr-2" />
                                    Edit RFC
                                </Button>
                                <Button variant="outline" className="w-full justify-center border-white/10 hover:bg-white/5 text-zinc-300">
                                    <FileText className="w-4 h-4 mr-2" />
                                    Export PDF
                                </Button>
                            </div>
                        </CardContent>
                    </Card>
                </div>

                {/* Main Content (Document) */}
                <div className="flex-1 overflow-auto rounded-xl border border-white/10 bg-[#0A0B10] shadow-2xl relative">
                    {/* Paper texture overlay could go here */}
                    <div className="max-w-4xl mx-auto p-12 min-h-full">
                        <div className="prose prose-invert max-w-none">
                            <pre className="font-mono text-zinc-300 whitespace-pre-wrap leading-relaxed text-sm">
                                {rfc.content}
                            </pre>
                        </div>
                    </div>
                </div>
            </div>
        </PageTransition>
    );
}
