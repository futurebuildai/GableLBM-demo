import { motion } from "framer-motion";
import { ArrowRight } from "lucide-react";

export const Hero = () => {
    return (
        <section className="relative min-h-screen flex items-center justify-center overflow-hidden py-20 px-4">
            {/* Background elements */}
            <div className="absolute inset-0 z-0">
                <div className="absolute top-1/4 left-1/4 w-96 h-96 bg-gable-green/10 rounded-full blur-[128px] animate-pulse-slow" />
                <div className="absolute bottom-1/4 right-1/4 w-96 h-96 bg-gable-blue/10 rounded-full blur-[128px] animate-pulse-slow" />
                <div className="absolute inset-0 bg-[url('https://www.transparenttextures.com/patterns/carbon-fibre.png')] opacity-5" />
            </div>

            <div className="container mx-auto relative z-10">
                <div className="flex flex-col items-center text-center space-y-8">
                    <motion.div
                        initial={{ opacity: 0, y: 20 }}
                        animate={{ opacity: 1, y: 0 }}
                        transition={{ duration: 0.5 }}
                        className="flex items-center space-x-3 mb-4"
                    >
                        <img src="/logo.png" alt="Gable Logo" className="h-12 w-auto" />
                        <div className="h-8 w-px bg-white/10 hidden sm:block" />
                        <span className="text-gable-green font-mono text-sm tracking-widest uppercase hidden sm:block">
                            GableLBM: Open Source & Co-Op Governed
                        </span>
                    </motion.div>

                    <motion.h1
                        initial={{ opacity: 0, y: 20 }}
                        animate={{ opacity: 1, y: 0 }}
                        transition={{ duration: 0.5, delay: 0.1 }}
                        className="text-5xl md:text-7xl lg:text-9xl font-bold max-w-7xl leading-tight"
                    >
                        Own the Core. <br />
                        <span className="text-gable-green italic">Sovereign LBM Operations.</span>
                    </motion.h1>

                    <motion.p
                        initial={{ opacity: 0, y: 20 }}
                        animate={{ opacity: 1, y: 0 }}
                        transition={{ duration: 0.5, delay: 0.2 }}
                        className="text-xl md:text-2xl text-slate-400 max-w-4xl font-light"
                    >
                        Stop renting your business foundation. GableLBM is the first 100% open-source operations core for the independent yard—stewarded by <strong>FutureBuild AI</strong> and governed by industry co-ops.
                        No more ransom. No more silos. Just your yard, your way.
                    </motion.p>

                    <motion.div
                        initial={{ opacity: 0, y: 20 }}
                        animate={{ opacity: 1, y: 0 }}
                        transition={{ duration: 0.5, delay: 0.3 }}
                        className="flex flex-col sm:flex-row items-center space-y-4 sm:space-y-0 sm:space-x-6 pt-8 w-full justify-center"
                    >
                        <button className="industrial-button electric-glow bg-gable-green text-gable-bg rounded-lg hover-depth font-bold group px-10 py-4 uppercase tracking-tighter">
                            Secure Your Yard
                            <ArrowRight className="ml-2 h-5 w-5 transition-transform group-hover:translate-x-1" />
                        </button>
                        <button className="industrial-button border-2 border-white/10 text-white rounded-lg hover:border-white/30 hover-depth font-bold group px-10 py-4 uppercase tracking-tighter transition-colors">
                            Apply for Advisory Board
                            <ArrowRight className="ml-2 h-5 w-5 transition-transform group-hover:translate-x-1" />
                        </button>
                    </motion.div>

                    {/* Technical Stat Cards (Glassmorphism) */}
                    <motion.div
                        initial={{ opacity: 0, scale: 0.95 }}
                        animate={{ opacity: 1, scale: 1 }}
                        transition={{ duration: 0.5, delay: 0.4 }}
                        className="grid grid-cols-1 md:grid-cols-3 gap-6 pt-20 w-full max-w-5xl"
                    >
                        {[
                            { label: "Core Status", value: "STABLE", detail: "Member-tested core" },
                            { label: "Ownership", value: "MEMBER OWNED", detail: "Zero-tax infrastructure" },
                            { label: "Governance", value: "ROADMAP CONTROL", detail: "Co-op driven development" },
                        ].map((stat, i) => (
                            <div key={i} className="glass-card p-6 text-left hover-depth border border-white/5">
                                <p className="text-xs font-mono text-gable-green tracking-tighter uppercase">{stat.label}</p>
                                <h3 className="text-3xl font-bold mt-2 font-mono text-white tracking-tighter">{stat.value}</h3>
                                <p className="text-sm text-slate-500 mt-1">{stat.detail}</p>
                            </div>
                        ))}
                    </motion.div>
                </div>
            </div>
        </section>
    );
};
