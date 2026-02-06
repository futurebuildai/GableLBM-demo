import { motion } from "framer-motion";
import { Shield, Zap, Eye, Cpu, Link } from "lucide-react";

export const ManifestoSection = () => {
    const beliefs = [
        {
            icon: <Shield className="w-6 h-6 text-gable-green" />,
            title: "Total Industry Sovereignty",
            text: "No more renting. Dealers own the code. Co-ops control the protocol. We build for 50-year business continuity, not next quarter's SaaS fees.",
        },
        {
            icon: <Zap className="w-6 h-6 text-gable-green" />,
            title: "Design Over Complexity",
            text: "If a new hire can't sell a stud within 15 minutes, the software has failed. We build for the day-one employee.",
        },
        {
            icon: <Eye className="w-6 h-6 text-gable-green" />,
            title: "Real-Time is the Only Time",
            text: "In a volatile market, 'Nightly Sync' is a death sentence. Information must flow like materials—visibly.",
        },
        {
            icon: <Cpu className="w-6 h-6 text-gable-green" />,
            title: "AI Must Be Alpha",
            text: "We use AI to count windows on blueprints and load trucks in LIFO sequence—not just to summarize emails.",
        },
        {
            icon: <Link className="w-6 h-6 text-gable-green" />,
            title: "Open Integration",
            text: "Connecting your ERP to your website or truck trackers should be standard and free. We've eliminated the 'Integration Tax'.",
        },
    ];

    return (
        <section className="py-24 bg-white/5 relative overflow-hidden">
            <div className="container mx-auto px-4 relative z-10">
                <div className="max-w-4xl mx-auto text-center mb-16">
                    <motion.h2
                        initial={{ opacity: 0, y: 20 }}
                        whileInView={{ opacity: 1, y: 0 }}
                        viewport={{ once: true }}
                        className="text-4xl md:text-5xl font-bold mb-6"
                    >
                        The GableLBM Commitment
                    </motion.h2>
                    <motion.div
                        initial={{ opacity: 0 }}
                        whileInView={{ opacity: 1 }}
                        viewport={{ once: true }}
                        className="h-1 w-20 bg-gable-green mx-auto mb-8"
                    />
                    <motion.p
                        initial={{ opacity: 0 }}
                        whileInView={{ opacity: 1 }}
                        viewport={{ once: true }}
                        className="text-xl text-slate-400 font-light"
                    >
                        GableLBM is the industry's digital commons. A 100% open-source, co-op governed operations core stewarded by FutureBuild AI.
                    </motion.p>
                </div>

                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
                    {beliefs.map((belief, i) => (
                        <motion.div
                            key={i}
                            initial={{ opacity: 0, y: 20 }}
                            whileInView={{ opacity: 1, y: 0 }}
                            viewport={{ once: true }}
                            transition={{ delay: i * 0.1 }}
                            className="glass-card p-8 border-white/10 hover:border-gable-green/30 transition-colors"
                        >
                            <div className="mb-4 p-3 bg-gable-green/10 inline-block rounded-lg">
                                {belief.icon}
                            </div>
                            <h3 className="text-xl font-bold mb-4">{belief.title}</h3>
                            <p className="text-slate-400 leading-relaxed font-light">
                                {belief.text}
                            </p>
                        </motion.div>
                    ))}

                    <motion.div
                        initial={{ opacity: 0, y: 20 }}
                        whileInView={{ opacity: 1, y: 0 }}
                        viewport={{ once: true }}
                        transition={{ delay: 0.5 }}
                        className="p-8 border-2 border-dashed border-white/5 rounded-2xl flex flex-col justify-center items-center text-center group hover:border-gable-green/20 transition-colors"
                    >
                        <h3 className="text-xl font-bold mb-2 text-slate-500 group-hover:text-slate-300 transition-colors">Join the Build</h3>
                        <p className="text-slate-600 group-hover:text-slate-400 transition-colors">Contribute to the open source core of LBM.</p>
                    </motion.div>
                </div>
            </div>
        </section>
    );
};
