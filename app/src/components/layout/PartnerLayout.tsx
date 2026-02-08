import { useState } from 'react';
import { Outlet, Link, useLocation } from "react-router-dom";
import { LayoutDashboard, FileText, ClipboardList, Bell, ChevronLeft, ChevronRight } from "lucide-react";
import { cn } from "../../lib/utils";

import { BrandLogo } from '../ui/BrandLogo';

export const PartnerLayout = () => {
    const [sidebarOpen, setSidebarOpen] = useState(true);
    const location = useLocation();

    const navItems = [
        { icon: <LayoutDashboard size={20} />, label: "Overview", path: "/partner" },
        { icon: <ClipboardList size={20} />, label: "Projects", path: "/partner/projects" },
        { icon: <FileText size={20} />, label: "Invoices", path: "/partner/invoices" },
    ];

    return (
        <div className="min-h-screen bg-deep-space text-foreground flex font-sans selection:bg-gable-green/30">
            {/* Sidebar */}
            <aside
                className={cn(
                    "bg-slate-steel border-r border-white/10 transition-all duration-300 flex flex-col fixed inset-y-0 left-0 z-50 shadow-2xl",
                    sidebarOpen ? "w-64" : "w-20"
                )}
            >
                <div className="h-16 flex items-center justify-between px-6 border-b border-white/5 bg-white/5 backdrop-blur-sm">
                    <div className={cn("flex items-center gap-3 transition-opacity duration-200", !sidebarOpen && "opacity-0 hidden")}>
                        <BrandLogo variant="full" size="md" />
                    </div>
                    {!sidebarOpen && (
                        <div className="mx-auto flex items-center justify-center">
                            <BrandLogo variant="mark" size="md" />
                        </div>
                    )}
                </div>

                <nav className="flex-1 py-6 px-3 space-y-1">
                    {navItems.map((item) => {
                        const isActive = location.pathname === item.path;
                        return (
                            <Link
                                key={item.path}
                                to={item.path}
                                className={cn(
                                    "flex items-center gap-3 px-3 py-3 rounded-lg transition-all duration-200 group relative overflow-hidden",
                                    isActive
                                        ? "bg-gable-green/10 text-gable-green shadow-[0_0_20px_rgba(0,255,163,0.1)]"
                                        : "text-zinc-400 hover:text-zinc-100 hover:bg-white/5"
                                )}
                            >
                                {isActive && (
                                    <div className="absolute inset-y-0 left-0 w-1 bg-gable-green rounded-r-full" />
                                )}
                                <span className={cn("transition-transform duration-200 group-hover:scale-110", isActive && "scale-110")}>
                                    {item.icon}
                                </span>
                                <span className={cn(
                                    "font-medium transition-all duration-300 origin-left",
                                    sidebarOpen ? "opacity-100 translate-x-0" : "opacity-0 -translate-x-4 absolute"
                                )}>
                                    {item.label}
                                </span>
                            </Link>
                        );
                    })}
                </nav>

                <div className="p-4 border-t border-white/5 bg-white/5">
                    <button
                        onClick={() => setSidebarOpen(!sidebarOpen)}
                        className="w-full flex items-center justify-center p-2 rounded-lg hover:bg-white/5 text-zinc-500 hover:text-white transition-colors"
                    >
                        {sidebarOpen ? <ChevronLeft size={20} /> : <ChevronRight size={20} />}
                    </button>

                    <div className={cn("mt-4 flex items-center gap-3 transition-all duration-300", !sidebarOpen && "justify-center")}>
                        <div className="h-8 w-8 rounded-full bg-gradient-to-br from-indigo-500 to-purple-600 p-[1px]">
                            <div className="h-full w-full rounded-full bg-slate-steel flex items-center justify-center">
                                <span className="font-bold text-xs text-white">JD</span>
                            </div>
                        </div>
                        {sidebarOpen && (
                            <div className="overflow-hidden">
                                <div className="text-sm font-medium text-white truncate">John Doe</div>
                                <div className="text-xs text-zinc-500 truncate">Acme Builders</div>
                            </div>
                        )}
                    </div>
                </div>
            </aside>

            {/* Main Content */}
            <main className={cn(
                "flex-1 flex flex-col min-h-screen transition-all duration-300",
                sidebarOpen ? "ml-64" : "ml-20"
            )}>
                {/* Header */}
                <header className="h-16 border-b border-white/10 bg-slate-steel/80 backdrop-blur-md sticky top-0 z-40 px-8 flex items-center justify-between shadow-sm">
                    <div className="flex items-center gap-4">
                        <h2 className="text-xl font-semibold text-white tracking-tight">
                            {navItems.find(i => i.path === location.pathname)?.label || 'Portal'}
                        </h2>
                    </div>
                    <div className="flex items-center gap-4">
                        <button className="relative p-2 text-zinc-400 hover:text-white transition-colors rounded-full hover:bg-white/5">
                            <Bell size={20} />
                            <span className="absolute top-1.5 right-1.5 w-2 h-2 bg-rose-500 rounded-full border border-slate-steel"></span>
                        </button>
                    </div>
                </header>

                <div className="flex-1 p-8 overflow-auto">
                    <Outlet />
                </div>
            </main>
        </div>
    );
};
