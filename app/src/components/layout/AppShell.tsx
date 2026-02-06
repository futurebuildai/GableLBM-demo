import { LayoutDashboard, Package, Truck, FileText, Settings, Menu } from 'lucide-react';
import { cn } from '../../lib/utils';
import { useState } from 'react';
import { Link, useLocation } from 'react-router-dom';
import { Omnibar } from '../ui/Omnibar';
import { ShortcutsModal } from '../ui/ShortcutsModal';
import { useEffect } from 'react';

export function AppShell({ children }: { children: React.ReactNode }) {
    const [sidebarOpen, setSidebarOpen] = useState(true);
    const [shortcutsOpen, setShortcutsOpen] = useState(false);
    const location = useLocation();

    useEffect(() => {
        const handleKeyDown = (e: KeyboardEvent) => {
            if (e.key === '?' && !e.metaKey && !e.ctrlKey && !['INPUT', 'TEXTAREA'].includes((e.target as HTMLElement).tagName)) {
                e.preventDefault();
                setShortcutsOpen(true);
            }
        };

        window.addEventListener('keydown', handleKeyDown);
        return () => window.removeEventListener('keydown', handleKeyDown);
    }, []);

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
                    <div className="flex-1 font-bold text-gable-green text-xl truncate">
                        {sidebarOpen ? 'GOMORRAH' : 'G'}
                    </div>
                </div>

                <nav className="flex-1 p-2 space-y-1">
                    <NavItem
                        to="/"
                        icon={<LayoutDashboard size={20} />}
                        label="Dashboard"
                        isOpen={sidebarOpen}
                        active={location.pathname === '/'}
                    />
                    <NavItem
                        to="/inventory"
                        icon={<Package size={20} />}
                        label="Inventory"
                        isOpen={sidebarOpen}
                        active={location.pathname === '/inventory'}
                    />
                    <NavItem
                        to="/locations"
                        icon={<LayoutDashboard size={20} />} // Reusing icon or pick MapPin
                        label="Locations"
                        isOpen={sidebarOpen}
                        active={location.pathname.startsWith('/locations')}
                    />
                    <NavItem
                        to="/sales"
                        icon={<FileText size={20} />}
                        label="Sales"
                        isOpen={sidebarOpen}
                        active={location.pathname.startsWith('/sales')}
                    />
                    <NavItem
                        to="/orders"
                        icon={<FileText size={20} />}
                        label="Orders"
                        isOpen={sidebarOpen}
                        active={location.pathname.startsWith('/orders')}
                    />
                    <NavItem
                        to="/invoices"
                        icon={<FileText size={20} />}
                        label="Invoices"
                        isOpen={sidebarOpen}
                        active={location.pathname.startsWith('/invoices')}
                    />
                    <NavItem
                        to="/reports/daily-till"
                        icon={<LayoutDashboard size={20} />}
                        label="Daily Till"
                        isOpen={sidebarOpen}
                        active={location.pathname.startsWith('/reports')}
                    />
                    <NavItem
                        to="/logistics"
                        icon={<Truck size={20} />}
                        label="Logistics"
                        isOpen={sidebarOpen}
                        active={location.pathname.startsWith('/logistics')}
                    />
                </nav>

                <div className="p-2 border-t border-white/10">
                    <NavItem
                        to="/admin"
                        icon={<Settings size={20} />}
                        label="Admin"
                        isOpen={sidebarOpen}
                        active={location.pathname.startsWith('/admin')}
                    />
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
                        <div className="text-xs text-muted-foreground hidden lg:block">
                            Press <span className="font-mono bg-white/10 px-1 rounded">⌘K</span> to search
                        </div>
                        <div className="h-8 w-8 rounded-full bg-gable-green/20 border border-gable-green/50 flex items-center justify-center text-xs font-mono text-gable-green">
                            AD
                        </div>
                    </div>
                </header>

                {/* Page Content */}
                <div className="p-6">
                    {children}
                </div>
            </main>
            <Omnibar />
            <ShortcutsModal isOpen={shortcutsOpen} onClose={() => setShortcutsOpen(false)} />
        </div>
    );
}

function NavItem({ icon, label, isOpen, active = false, to }: { icon: React.ReactNode, label: string, isOpen: boolean, active?: boolean, to: string }) {
    return (
        <Link to={to} className={cn(
            "w-full flex items-center gap-3 px-3 py-2 rounded-md transition-colors text-sm font-medium group",
            active
                ? "bg-gable-green/10 text-gable-green"
                : "text-muted-foreground hover:bg-white/5 hover:text-white"
        )}>
            <span className={cn(active ? "text-gable-green" : "text-muted-foreground group-hover:text-white")}>
                {icon}
            </span>
            {isOpen && <span>{label}</span>}
        </Link>
    )
}
