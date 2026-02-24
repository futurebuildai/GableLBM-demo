import { useNavigate } from "react-router-dom";
import { Monitor, Truck, Warehouse, ShoppingBag, ArrowRight, ExternalLink } from "lucide-react";

const surfaces = [
    {
        id: "erp",
        title: "ERP Desktop",
        subtitle: "Back-Office Operations",
        description: "Full lumber yard management: inventory, quotes, orders, invoices, purchasing, accounting, dispatch, and millwork configurator.",
        icon: Monitor,
        color: "#00FFA3",
        path: "/erp",
        stats: "31 views",
        badge: "Primary",
    },
    {
        id: "portal",
        title: "Customer Portal",
        subtitle: "Contractor Self-Service",
        description: "White-labeled B2B portal where contractors check orders, pay invoices, and track deliveries in real-time.",
        icon: ShoppingBag,
        color: "#00FFA3",
        path: "/portal",
        stats: "4 views",
        badge: "B2B",
    },
    {
        id: "driver",
        title: "Driver Mobile",
        subtitle: "Delivery Management",
        description: "Mobile-first app for delivery drivers: route assignments, stop-by-stop navigation, and proof-of-delivery capture.",
        icon: Truck,
        color: "#00FFA3",
        path: "/driver",
        stats: "3 views",
        badge: "Mobile",
    },
    {
        id: "yard",
        title: "Yard Mobile",
        subtitle: "Picks & Inventory",
        description: "Handheld interface for yard pickers: pick queues, inventory lookup, cycle counts, and PO receiving.",
        icon: Warehouse,
        color: "#F59E0B",
        path: "/yard",
        stats: "5 views",
        badge: "Mobile",
    },
];

export function DemoLanding() {
    const navigate = useNavigate();

    return (
        <div className="min-h-screen font-sans" style={{ backgroundColor: "#0A0B10" }}>
            {/* Ambient glow */}
            <div
                className="fixed inset-0 pointer-events-none"
                style={{
                    background:
                        "radial-gradient(ellipse 80% 50% at 50% -10%, rgba(0,255,163,0.08) 0%, transparent 60%), radial-gradient(ellipse 60% 40% at 80% 90%, rgba(245,158,11,0.05) 0%, transparent 50%)",
                }}
            />

            <div className="relative max-w-5xl mx-auto px-6 py-16 md:py-24">
                {/* Header */}
                <div className="text-center mb-16">
                    <div className="inline-flex items-center gap-2 mb-6 px-4 py-1.5 rounded-full border border-white/10 bg-white/5 text-xs font-mono text-zinc-400 uppercase tracking-wider">
                        <span className="w-1.5 h-1.5 rounded-full bg-emerald-400 animate-pulse" />
                        Live Demo Environment
                    </div>
                    <h1 className="text-4xl md:text-6xl font-bold text-white tracking-tight mb-4">
                        Gable<span style={{ color: "#00FFA3" }}>LBM</span>
                    </h1>
                    <p className="text-lg md:text-xl text-zinc-400 max-w-2xl mx-auto leading-relaxed">
                        Modern ERP for lumber &amp; building materials dealers.
                        <br className="hidden md:block" />
                        Explore four purpose-built app surfaces, fully seeded with demo data.
                    </p>
                </div>

                {/* Surface Cards */}
                <div className="grid grid-cols-1 md:grid-cols-2 gap-5">
                    {surfaces.map((s) => (
                        <button
                            key={s.id}
                            onClick={() => navigate(s.path)}
                            className="group relative text-left rounded-2xl border border-white/[0.06] bg-white/[0.02] p-6 md:p-8 transition-all duration-300 hover:border-white/15 hover:bg-white/[0.04] hover:-translate-y-1 hover:shadow-2xl active:scale-[0.98] cursor-pointer overflow-hidden"
                        >
                            {/* Hover glow */}
                            <div
                                className="absolute inset-0 opacity-0 group-hover:opacity-100 transition-opacity duration-500 pointer-events-none"
                                style={{
                                    background: `radial-gradient(circle at 50% 100%, ${s.color}08 0%, transparent 60%)`,
                                }}
                            />

                            <div className="relative">
                                {/* Top row: icon + badges */}
                                <div className="flex items-start justify-between mb-5">
                                    <div
                                        className="p-3 rounded-xl border transition-colors duration-300"
                                        style={{
                                            backgroundColor: `${s.color}08`,
                                            borderColor: `${s.color}15`,
                                        }}
                                    >
                                        <s.icon
                                            className="w-6 h-6 transition-transform duration-300 group-hover:scale-110"
                                            style={{ color: s.color }}
                                        />
                                    </div>
                                    <div className="flex items-center gap-2">
                                        <span className="text-[10px] font-mono px-2 py-0.5 rounded border border-white/10 text-zinc-500 uppercase tracking-wider">
                                            {s.badge}
                                        </span>
                                        <span className="text-[10px] font-mono px-2 py-0.5 rounded border border-white/10 text-zinc-600">
                                            {s.stats}
                                        </span>
                                    </div>
                                </div>

                                {/* Title */}
                                <h2 className="text-xl md:text-2xl font-bold text-white mb-1 tracking-tight">
                                    {s.title}
                                </h2>
                                <p
                                    className="text-sm font-medium mb-3"
                                    style={{ color: s.color }}
                                >
                                    {s.subtitle}
                                </p>

                                {/* Description */}
                                <p className="text-sm text-zinc-500 leading-relaxed mb-6">
                                    {s.description}
                                </p>

                                {/* CTA */}
                                <div className="flex items-center gap-2 text-sm font-medium text-zinc-400 group-hover:text-white transition-colors">
                                    Launch Demo
                                    <ArrowRight className="w-4 h-4 transition-transform duration-300 group-hover:translate-x-1" />
                                </div>
                            </div>
                        </button>
                    ))}
                </div>

                {/* Footer */}
                <div className="mt-16 text-center space-y-4">
                    <div className="flex items-center justify-center gap-6 text-xs text-zinc-600">
                        <span className="flex items-center gap-1.5">
                            <span className="w-1.5 h-1.5 rounded-full bg-emerald-500" />
                            42 products seeded
                        </span>
                        <span className="flex items-center gap-1.5">
                            <span className="w-1.5 h-1.5 rounded-full bg-blue-500" />
                            12 customers
                        </span>
                        <span className="flex items-center gap-1.5">
                            <span className="w-1.5 h-1.5 rounded-full bg-amber-500" />
                            132 orders
                        </span>
                    </div>
                    <p className="text-xs text-zinc-700">
                        All data is synthetic and resets periodically.{" "}
                        <a
                            href="https://gablelbm.com"
                            target="_blank"
                            rel="noopener noreferrer"
                            className="inline-flex items-center gap-1 text-zinc-500 hover:text-white transition-colors"
                        >
                            gablelbm.com <ExternalLink className="w-3 h-3" />
                        </a>
                    </p>
                </div>
            </div>
        </div>
    );
}
