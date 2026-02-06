import { LayoutDashboard, FileText, CheckCircle, Menu, LogOut } from 'lucide-react';
import { cn } from '../../lib/utils';
import { useState } from 'react';
import { Link, useLocation, Outlet } from 'react-router-dom';

export function PartnerLayout() {
    const [sidebarOpen, setSidebarOpen] = useState(true);
    const location = useLocation();

    return (
        <div className="min-h-screen bg-deep-space text-foreground flex">
            {/* Sidebar */}
            <aside
                className={cn(
                    "bg-slate-steel border-r border-white/10 transition-all duration-300 flex flex-col fixed inset-y-0 left-0 z-50",
                    sidebarOpen ? "w-64" : "w-16"
                )}
            >
                <div className="h-14 flex items-center px-4 border-b border-white/10">
                    <div className="flex-1 font-bold text-blue-400 text-xl truncate">
                        {sidebarOpen ? 'PARTNER PORTAL' : 'P'}
                    </div>
                </div>

                <nav className="flex-1 p-2 space-y-1">
                    <NavItem
                        to="/partner"
                        icon={<LayoutDashboard size={20} />}
                        label="Dashboard"
                        isOpen={sidebarOpen}
                        active={location.pathname === '/partner'}
                    />
                    <NavItem
                        to="/partner/projects"
                        icon={<CheckCircle size={20} />} // Or Briefcase
                        label="Projects"
                        isOpen={sidebarOpen}
                        active={location.pathname.startsWith('/partner/projects')}
                    />
                    <NavItem
                        to="/partner/invoices"
                        icon={<FileText size={20} />}
                        label="Invoices"
                        isOpen={sidebarOpen}
                        active={location.pathname.startsWith('/partner/invoices')}
                    />
                </nav>

                <div className="p-2 border-t border-white/10">
                    <button className={cn(
                        "w-full flex items-center gap-3 px-3 py-2 rounded-md transition-colors text-sm font-medium text-muted-foreground hover:bg-white/5 hover:text-white"
                    )}>
                        <LogOut size={20} />
                        {sidebarOpen && <span>Sign Out</span>}
                    </button>
                </div>
            </aside>

            {/* Main Content */}
            <main className={cn(
                "flex-1 flex flex-col min-h-screen transition-all duration-300",
                sidebarOpen ? "ml-64" : "ml-16"
            )}>
                {/* Header */}
                <header className="h-14 border-b border-white/10 bg-slate-steel/50 backdrop-blur-md px-4 flex items-center justify-between sticky top-0 z-40">
                    <button
                        onClick={() => setSidebarOpen(!sidebarOpen)}
                        className="p-2 hover:bg-white/5 rounded-md text-muted-foreground hover:text-white"
                    >
                        <Menu size={20} />
                    </button>
                    <div className="flex items-center gap-4">
                        <div className="h-8 w-8 rounded-full bg-blue-500/20 border border-blue-500/50 flex items-center justify-center text-xs font-mono text-blue-400">
                            CTR
                        </div>
                    </div>
                </header>

                {/* Page Content */}
                <div className="p-6">
                    <Outlet />
                </div>
            </main>
        </div>
    );
}

function NavItem({ icon, label, isOpen, active = false, to }: { icon: React.ReactNode, label: string, isOpen: boolean, active?: boolean, to: string }) {
    return (
        <Link to={to} className={cn(
            "w-full flex items-center gap-3 px-3 py-2 rounded-md transition-colors text-sm font-medium group",
            active
                ? "bg-blue-500/10 text-blue-400"
                : "text-muted-foreground hover:bg-white/5 hover:text-white"
        )}>
            <span className={cn(active ? "text-blue-400" : "text-muted-foreground group-hover:text-white")}>
                {icon}
            </span>
            {isOpen && <span>{label}</span>}
        </Link>
    )
}
