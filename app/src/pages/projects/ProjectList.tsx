import { useState, useEffect, useCallback } from 'react';
import { useNavigate } from 'react-router-dom';
import { Card, CardContent } from '../../components/ui/Card';
import { FolderGit2, Plus, ArrowRight, RefreshCw, AlertTriangle, Briefcase } from 'lucide-react';
import { ProjectService } from '../../services/ProjectService';
import type { Project } from '../../types/project';
import { useToast } from '../../components/ui/ToastContext';
import { format } from 'date-fns';

export const ProjectList = () => {
    const navigate = useNavigate();
    const { showToast } = useToast();

    const [projects, setProjects] = useState<Project[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState('');
    const [isCreating, setIsCreating] = useState(false);
    const [newProjectName, setNewProjectName] = useState('');

    const fetchProjects = useCallback(async () => {
        setLoading(true);
        setError('');
        try {
            const data = await ProjectService.listProjects();
            setProjects(data);
        } catch (err) {
            setError(err instanceof Error ? err.message : 'Failed to load projects');
        } finally {
            setLoading(false);
        }
    }, []);

    useEffect(() => { fetchProjects(); }, [fetchProjects]);

    const handleCreateProject = async (e: React.FormEvent) => {
        e.preventDefault();
        if (!newProjectName.trim()) return;

        try {
            const newProject = await ProjectService.createProject({ name: newProjectName });
            showToast('Project created successfully', 'success');
            setIsCreating(false);
            setNewProjectName('');
            setProjects([newProject, ...projects]);
            navigate(`/portal/projects/${newProject.id}`);
        } catch (err) {
            showToast(err instanceof Error ? err.message : 'Failed to create project', 'error');
        }
    };

    if (loading) {
        return (
            <div className="space-y-4">
                <div className="h-10 w-1/4 bg-white/5 rounded-lg animate-pulse mb-6" />
                {[1, 2, 3].map(i => (
                    <div key={i} className="h-20 bg-white/5 rounded-2xl animate-pulse" />
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
                    onClick={fetchProjects}
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
                    <h1 className="text-2xl font-bold text-white flex items-center gap-3">
                        <FolderGit2 className="text-gable-green" /> Projects
                    </h1>
                    <p className="text-zinc-400 text-sm mt-1">Organize orders, deliveries, and invoices by job site or project</p>
                </div>
                {!isCreating && (
                    <button
                        onClick={() => setIsCreating(true)}
                        className="flex items-center gap-2 px-4 py-2 bg-gable-green text-black font-semibold rounded-lg hover:bg-emerald-400 transition-colors shadow-[0_0_15px_rgba(0,255,163,0.3)]"
                    >
                        <Plus size={18} /> New Project
                    </button>
                )}
            </div>

            {isCreating && (
                <Card variant="glass" className="mb-6 border-gable-green/30">
                    <CardContent className="p-4">
                        <form onSubmit={handleCreateProject} className="flex gap-3 items-end">
                            <div className="flex-1">
                                <label className="block text-sm font-medium text-zinc-300 mb-1">Project Name</label>
                                <input
                                    type="text"
                                    autoFocus
                                    required
                                    value={newProjectName}
                                    onChange={(e) => setNewProjectName(e.target.value)}
                                    placeholder="e.g. 123 Main St Subivision"
                                    className="w-full bg-black/20 border border-white/10 rounded-lg py-2 px-3 text-white placeholder-zinc-600 focus:outline-none focus:ring-1 focus:ring-gable-green"
                                />
                            </div>
                            <button
                                type="button"
                                onClick={() => setIsCreating(false)}
                                className="px-4 py-2 text-sm font-medium text-zinc-400 hover:text-white transition-colors"
                            >
                                Cancel
                            </button>
                            <button
                                type="submit"
                                disabled={!newProjectName.trim()}
                                className="px-4 py-2 bg-white/10 text-white font-semibold rounded-lg hover:bg-white/20 transition-colors disabled:opacity-50"
                            >
                                Create
                            </button>
                        </form>
                    </CardContent>
                </Card>
            )}

            {projects.length === 0 ? (
                <Card variant="glass">
                    <CardContent className="p-12 text-center">
                        <Briefcase className="w-12 h-12 text-zinc-600 mx-auto mb-4" />
                        <h3 className="text-lg font-medium text-white mb-2">No projects found</h3>
                        <p className="text-zinc-400 mb-6">Group your orders and deliveries by creating your first project.</p>
                        <button
                            onClick={() => setIsCreating(true)}
                            className="inline-flex items-center gap-2 px-4 py-2 bg-gable-green text-black font-semibold rounded-lg hover:bg-emerald-400 transition-colors"
                        >
                            <Plus size={18} /> New Project
                        </button>
                    </CardContent>
                </Card>
            ) : (
                <div className="space-y-3">
                    {projects.map(project => (
                        <Card key={project.id} variant="glass" noPadding>
                            <div
                                onClick={() => navigate(`/portal/projects/${project.id}`)}
                                className="p-4 flex items-center justify-between cursor-pointer hover:bg-white/5 transition-colors group"
                            >
                                <div className="flex items-center gap-4">
                                    <div className="w-10 h-10 rounded-lg bg-emerald-500/10 border border-emerald-500/20 flex items-center justify-center text-emerald-400 group-hover:bg-emerald-500/20 transition-colors">
                                        <FolderGit2 size={20} />
                                    </div>
                                    <div>
                                        <h3 className="font-medium text-white">{project.name}</h3>
                                        <div className="text-sm text-zinc-500 mt-1 flex items-center gap-2">
                                            <span>Created {format(new Date(project.created_at), 'MMM d, yyyy')}</span>
                                            <span>·</span>
                                            <span className={`text-[10px] uppercase font-semibold tracking-wider px-1.5 py-0.5 rounded border ${project.status === 'Active'
                                                    ? 'bg-blue-500/10 text-blue-400 border-blue-500/20'
                                                    : 'bg-zinc-500/10 text-zinc-400 border-zinc-500/20'
                                                }`}>
                                                {project.status}
                                            </span>
                                        </div>
                                    </div>
                                </div>
                                <ArrowRight className="text-zinc-600 group-hover:text-emerald-400 transition-colors" />
                            </div>
                        </Card>
                    ))}
                </div>
            )}
        </div>
    );
};
