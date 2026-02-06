import { motion } from "framer-motion";
import { Menu, X } from "lucide-react";
import { useState } from "react";

export const Navbar = () => {
    const [isOpen, setIsOpen] = useState(false);

    return (
        <nav className="fixed top-0 left-0 right-0 z-50 bg-gable-bg/80 backdrop-blur-lg border-b border-white/5">
            <div className="container mx-auto px-4 h-20 flex items-center justify-between">
                <div className="flex items-center space-x-2">
                    <img src="/logo.png" alt="Gable Logo" className="h-8 w-auto" />
                    <span className="font-bold text-xl tracking-tight uppercase tracking-widest">GableLBM</span>
                </div>

                <div className="hidden md:flex items-center space-x-10 text-sm font-medium text-slate-400">
                    <a href="#" className="hover:text-gable-green transition-colors">The Protocol</a>
                    <a href="#" className="hover:text-gable-green transition-colors">The Foundation</a>
                    <a href="#" className="hover:text-gable-green transition-colors">Governance</a>
                </div>

                <div className="hidden md:block">
                    <button className="industrial-button glass-card hover:bg-white/5 px-6 py-2 rounded-lg text-sm border-white/10 font-bold">
                        Get Started
                    </button>
                </div>

                <button className="md:hidden text-white" onClick={() => setIsOpen(!isOpen)}>
                    {isOpen ? <X /> : <Menu />}
                </button>
            </div>

            {/* Mobile Menu */}
            {isOpen && (
                <motion.div
                    initial={{ opacity: 0, y: -20 }}
                    animate={{ opacity: 1, y: 0 }}
                    className="md:hidden bg-gable-surface p-6 border-b border-white/5 space-y-1"
                >
                    <a href="#" className="block py-3 text-slate-300 border-b border-white/5 last:border-0 hover:text-gable-green transition-colors">Platform</a>
                    <a href="#" className="block py-3 text-slate-300 border-b border-white/5 last:border-0 hover:text-gable-green transition-colors">Manifesto</a>
                    <a href="#" className="block py-3 text-slate-300 border-b border-white/5 last:border-0 hover:text-gable-green transition-colors">Governance</a>
                    <a href="#" className="block py-3 text-gable-green font-bold last:border-0">Get Started</a>
                </motion.div>
            )}
        </nav>
    );
};
