import { motion } from "framer-motion";
import { ArrowRight } from "lucide-react";

export const CTA = () => {
    return (
        <section className="py-24 px-4 relative overflow-hidden">
            {/* Background decoration */}
            <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-[800px] h-[800px] bg-gable-green/5 blur-[160px] rounded-full z-0" />

            <div className="container mx-auto relative z-10">
                <div className="grid grid-cols-1 lg:grid-cols-2 gap-8 items-stretch pt-12">
                    {/* Dealer Track */}
                    <motion.div
                        initial={{ opacity: 0, x: -20 }}
                        whileInView={{ opacity: 1, x: 0 }}
                        viewport={{ once: true }}
                        className="glass-card p-12 flex flex-col items-center text-center space-y-6 border-gable-green/10 bg-gable-green/5"
                    >
                        <h3 className="text-3xl font-bold uppercase tracking-widest text-gable-green">LBM Dealers</h3>
                        <p className="text-lg text-slate-400 font-light">
                            Own your data. Control your roadmap. Eliminate the legacy vendor tax and reclaim your yard's sovereignty.
                        </p>
                        <button className="industrial-button bg-white text-gable-bg rounded-lg hover-depth font-bold group px-10 py-4 uppercase tracking-tighter mt-auto">
                            Reclaim Your Yard
                            <ArrowRight className="ml-2 h-4 w-4 transition-transform group-hover:translate-x-1" />
                        </button>
                    </motion.div>

                    {/* Co-op Track */}
                    <motion.div
                        initial={{ opacity: 0, x: 20 }}
                        whileInView={{ opacity: 1, x: 0 }}
                        viewport={{ once: true }}
                        className="glass-card p-12 flex flex-col items-center text-center space-y-6 border-gable-blue/10 bg-gable-blue/5"
                    >
                        <h3 className="text-3xl font-bold uppercase tracking-widest text-gable-blue">Co-Op Alliances</h3>
                        <p className="text-lg text-slate-400 font-light">
                            We are seeking visionary co-op partners to join the **GableLBM Advisory Board**. Lead the movement to standardize the industry core.
                        </p>
                        <button className="industrial-button border-2 border-gable-blue/30 text-white rounded-lg hover:border-gable-blue/60 hover-depth font-bold group px-10 py-4 uppercase tracking-tighter mt-auto">
                            Join Advisory Board
                            <ArrowRight className="ml-2 h-4 w-4 transition-transform group-hover:translate-x-1" />
                        </button>
                    </motion.div>
                </div>
            </div>
        </section>
    );
};
