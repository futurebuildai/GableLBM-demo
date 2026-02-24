import { Outlet, NavLink } from "react-router-dom";
import { ClipboardList, Package, ScanBarcode } from "lucide-react";

const navItems = [
    { icon: ClipboardList, label: "Pick", path: "/yard" },
    { icon: Package, label: "Inventory", path: "/yard/inventory" },
    { icon: ScanBarcode, label: "Receiving", path: "/yard/receiving" },
];

export function YardLayout() {
    return (
        <div className="min-h-screen bg-[#0A0B10] text-white font-sans md:max-w-md md:mx-auto md:border-x md:border-white/10 relative shadow-2xl flex flex-col">
            <header className="h-16 flex items-center justify-between px-4 border-b border-white/10 bg-[#161821]/80 backdrop-blur-md sticky top-0 z-50">
                <div className="font-bold text-lg tracking-wider font-mono">
                    GABLE<span className="text-amber-400">YARD</span>
                </div>
                <div className="h-8 w-8 rounded-full bg-amber-400/10 border border-amber-400/20 flex items-center justify-center text-xs font-mono text-amber-400">
                    Y1
                </div>
            </header>

            <main className="flex-1 pb-20 overflow-y-auto">
                <Outlet />
            </main>

            {/* Bottom Tab Navigation */}
            <nav className="fixed bottom-0 left-0 right-0 md:max-w-md md:mx-auto h-16 bg-[#161821]/95 backdrop-blur-md border-t border-white/10 flex items-center justify-around z-50">
                {navItems.map(item => (
                    <NavLink
                        key={item.path}
                        to={item.path}
                        end={item.path === "/yard"}
                        className={({ isActive }) =>
                            `flex flex-col items-center gap-1 px-4 py-2 rounded-lg transition-colors ${isActive
                                ? "text-amber-400"
                                : "text-zinc-500 hover:text-zinc-300"
                            }`
                        }
                    >
                        <item.icon className="w-5 h-5" />
                        <span className="text-[10px] font-mono uppercase tracking-wider">{item.label}</span>
                    </NavLink>
                ))}
            </nav>
        </div>
    );
}
