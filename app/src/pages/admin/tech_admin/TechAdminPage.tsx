import { useState, useEffect } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { Key, Globe, Activity, Plus, Trash2, Copy, Check, Eye } from 'lucide-react';
import { Button } from '../../../components/ui/Button';
import { techAdminService, type APIKey } from '../../../services/TechAdminService';
import { cn } from '../../../lib/utils';

// --- Subcomponents (In-file for speed, can break out later) ---

const APIKeyManager = () => {
    const [keys, setKeys] = useState<APIKey[]>([]);
    const [loading, setLoading] = useState(true);
    const [isCreating, setIsCreating] = useState(false);
    const [newKeyName, setNewKeyName] = useState('');
    const [generatedKey, setGeneratedKey] = useState<string | null>(null);
    const [copied, setCopied] = useState(false);

    useEffect(() => {
        loadKeys();
    }, []);

    const loadKeys = async () => {
        try {
            const data = await techAdminService.listKeys();
            setKeys(data);
        } catch (error) {
            console.error(error);
        } finally {
            setLoading(false);
        }
    };

    const handleCreate = async () => {
        if (!newKeyName) return;
        try {
            const res = await techAdminService.createKey(newKeyName, ['read:inventory', 'write:orders']); // Default scopes for now
            setGeneratedKey(res.api_key);
            setNewKeyName(''); // Reset input
            loadKeys();
        } catch (error) {
            console.error(error);
        }
    };

    const handleRevoke = async (id: string) => {
        if (!confirm('Are you sure you want to revoke this key? integrations using it will break immediately.')) return;
        try {
            await techAdminService.revokeKey(id);
            loadKeys();
        } catch (error) {
            console.error(error);
        }
    };

    const copyToClipboard = (text: string) => {
        navigator.clipboard.writeText(text);
        setCopied(true);
        setTimeout(() => setCopied(false), 2000);
    };

    return (
        <div className="space-y-6">
            <div className="flex justify-between items-center">
                <div>
                    <h2 className="text-xl font-bold text-white">API Keys</h2>
                    <p className="text-slate-400 text-sm">Manage access keys for 3rd-party integrations.</p>
                </div>
                <Button onClick={() => setIsCreating(true)} disabled={isCreating || !!generatedKey}>
                    <Plus size={16} className="mr-2" />
                    New API Key
                </Button>
            </div>

            {/* Creation Modal / Inline Form */}
            <AnimatePresence>
                {(isCreating || generatedKey) && (
                    <motion.div
                        initial={{ opacity: 0, height: 0 }}
                        animate={{ opacity: 1, height: 'auto' }}
                        exit={{ opacity: 0, height: 0 }}
                        className="bg-slate-steel/50 border border-white/10 rounded-lg p-6 overflow-hidden"
                    >
                        {!generatedKey ? (
                            <div className="flex gap-4 items-end">
                                <div className="flex-1 space-y-2">
                                    <label className="text-sm font-medium text-slate-300">Key Name (e.g. "Zapier Integration")</label>
                                    <input
                                        type="text"
                                        value={newKeyName}
                                        onChange={(e) => setNewKeyName(e.target.value)}
                                        className="w-full bg-deep-space border border-white/10 rounded px-3 py-2 text-white placeholder-slate-500 focus:outline-none focus:border-gable-green transition-colors"
                                        placeholder="Friendly name..."
                                        autoFocus
                                    />
                                </div>
                                <div className="flex gap-2">
                                    <Button variant="outline" onClick={() => setIsCreating(false)}>Cancel</Button>
                                    <Button onClick={handleCreate} disabled={!newKeyName}>Generate</Button>
                                </div>
                            </div>
                        ) : (
                            <div className="space-y-4">
                                <div className="flex items-center gap-2 text-amber-400 bg-amber-400/10 px-3 py-2 rounded text-sm">
                                    <Eye size={16} />
                                    <span>Copy this key now. You won't be able to see it again!</span>
                                </div>
                                <div className="flex gap-2">
                                    <code className="flex-1 bg-deep-space border border-gable-green/50 rounded px-4 py-3 text-gable-green font-mono text-lg break-all">
                                        {generatedKey}
                                    </code>
                                    <Button onClick={() => copyToClipboard(generatedKey)} variant="outline" className="h-[52px] w-[52px] p-0 flex items-center justify-center">
                                        {copied ? <Check size={20} className="text-gable-green" /> : <Copy size={20} />}
                                    </Button>
                                </div>
                                <div className="flex justify-end">
                                    <Button onClick={() => { setGeneratedKey(null); setIsCreating(false); }}>Done</Button>
                                </div>
                            </div>
                        )}
                    </motion.div>
                )}
            </AnimatePresence>

            {/* Key List */}
            <div className="bg-slate-steel border border-white/5 rounded-lg overflow-hidden">
                <table className="w-full text-left text-sm">
                    <thead className="bg-white/5 text-slate-400 font-medium">
                        <tr>
                            <th className="px-4 py-3">Name</th>
                            <th className="px-4 py-3">Prefix</th>
                            <th className="px-4 py-3">Created</th>
                            <th className="px-4 py-3">Last Used</th>
                            <th className="px-4 py-3">Status</th>
                            <th className="px-4 py-3 text-right">Actions</th>
                        </tr>
                    </thead>
                    <tbody className="divide-y divide-white/5">
                        {keys.map((key) => (
                            <tr key={key.id} className="hover:bg-white/5 transition-colors group">
                                <td className="px-4 py-3 font-medium text-white">{key.name}</td>
                                <td className="px-4 py-3 font-mono text-slate-400">{key.prefix}...</td>
                                <td className="px-4 py-3 text-slate-400">{new Date(key.created_at).toLocaleDateString()}</td>
                                <td className="px-4 py-3 text-slate-400">
                                    {key.last_used_at ? new Date(key.last_used_at).toLocaleDateString() : 'Never'}
                                </td>
                                <td className="px-4 py-3">
                                    {key.revoked_at ? (
                                        <span className="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium bg-red-500/10 text-red-500">Revoked</span>
                                    ) : (
                                        <span className="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium bg-gable-green/10 text-gable-green">Active</span>
                                    )}
                                </td>
                                <td className="px-4 py-3 text-right">
                                    {!key.revoked_at && (
                                        <button
                                            onClick={() => handleRevoke(key.id)}
                                            className="text-slate-500 hover:text-red-400 transition-colors p-1 rounded hover:bg-white/5"
                                            title="Revoke Key"
                                        >
                                            <Trash2 size={16} />
                                        </button>
                                    )}
                                </td>
                            </tr>
                        ))}
                        {keys.length === 0 && !loading && (
                            <tr>
                                <td colSpan={6} className="px-4 py-8 text-center text-slate-500">
                                    No API keys generated yet.
                                </td>
                            </tr>
                        )}
                    </tbody>
                </table>
            </div>
        </div>
    );
};

const IntegrationCard = ({ name, description, icon: Icon, connected = false }: { name: string, description: string, icon: any, connected?: boolean }) => (
    <div className="bg-slate-steel border border-white/5 p-6 rounded-lg flex items-start justify-between hover:border-white/10 transition-colors group">
        <div className="flex gap-4">
            <div className={cn("w-12 h-12 rounded-lg flex items-center justify-center shrink-0 transition-colors", connected ? "bg-gable-green/20 text-gable-green" : "bg-white/5 text-slate-400 group-hover:text-white")}>
                <Icon size={24} />
            </div>
            <div>
                <h3 className="text-lg font-semibold text-white mb-1">{name}</h3>
                <p className="text-slate-400 text-sm leading-relaxed">{description}</p>
            </div>
        </div>
        <Button variant={connected ? "outline" : "default"} className={cn(connected && "border-gable-green/50 text-gable-green hover:bg-gable-green/10")}>
            {connected ? "Configure" : "Connect"}
        </Button>
    </div>
);

const Integrations = () => {
    return (
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-4">
            <IntegrationCard
                name="Run Payments"
                description="Secure payment processing for card-present and online transactions. Preferred Partner."
                icon={Activity} // Placeholder icon
            />
            <IntegrationCard
                name="QuickBooks Online"
                description="Automatically sync invoices, payments, and customers with your General Ledger."
                icon={Globe}
            />
            <IntegrationCard
                name="Avalara AvaTax"
                description="Real-time tax calculation and compliance for all 50 states."
                icon={Globe}
            />
            <IntegrationCard
                name="Zapier / Make"
                description="Connect GableLBM to 5,000+ other apps via webhooks."
                icon={Activity}
            />
        </div>
    );
};

// --- Main Page Component ---

export const TechAdminPage = () => {
    const [activeTab, setActiveTab] = useState<'keys' | 'integrations' | 'health'>('keys');

    const tabs = [
        { id: 'keys', label: 'API Keys', icon: Key },
        { id: 'integrations', label: 'Integrations', icon: Globe },
        { id: 'health', label: 'System Health', icon: Activity },
    ];

    return (
        <div className="min-h-screen bg-deep-space p-8 space-y-8">
            {/* Header */}
            <div>
                <h1 className="text-3xl font-bold text-white tracking-tight mb-2">Tech Admin</h1>
                <p className="text-slate-400">Manage integrations, API access, and system configuration.</p>
            </div>

            {/* Tabs */}
            <div className="flex gap-1 border-b border-white/10">
                {tabs.map((tab) => (
                    <button
                        key={tab.id}
                        onClick={() => setActiveTab(tab.id as any)}
                        className={cn(
                            "flex items-center gap-2 px-6 py-3 text-sm font-medium transition-colors relative",
                            activeTab === tab.id ? "text-gable-green" : "text-slate-400 hover:text-white"
                        )}
                    >
                        <tab.icon size={16} />
                        {tab.label}
                        {activeTab === tab.id && (
                            <motion.div
                                layoutId="activeTab"
                                className="absolute bottom-0 left-0 right-0 h-0.5 bg-gable-green"
                            />
                        )}
                    </button>
                ))}
            </div>

            {/* Content Area */}
            <div className="max-w-5xl">
                {activeTab === 'keys' && <APIKeyManager />}
                {activeTab === 'integrations' && <Integrations />}
                {activeTab === 'health' && (
                    <div className="bg-slate-steel border border-white/5 rounded-lg p-12 text-center">
                        <Activity className="w-12 h-12 text-slate-500 mx-auto mb-4" />
                        <h3 className="text-lg font-medium text-white">System Health ok</h3>
                        <p className="text-slate-400 mt-2">All services are running normally. Logs coming soon.</p>
                    </div>
                )}
            </div>
        </div>
    );
};

export default TechAdminPage;
