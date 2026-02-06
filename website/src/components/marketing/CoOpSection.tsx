import { motion } from "framer-motion";
import { Users, Database, Cpu } from "lucide-react";

export const CoOpSection = () => {
    return (
        <section className="py-24 bg-gable-surface/30 px-4">
            <div className="container mx-auto">
                <div className="flex flex-col lg:flex-row items-center gap-16">
                    <div className="lg:w-1/2">
                        <motion.div
                            initial={{ opacity: 0, x: -20 }}
                            whileInView={{ opacity: 1, x: 0 }}
                            viewport={{ once: true }}
                            className="space-y-6"
                        >
                            <div className="inline-flex items-center space-x-2 text-gable-green font-mono text-sm tracking-widest uppercase mb-4">
                                <Users className="w-4 h-4" />
                                <span>Industry Partners</span>
                            </div>
                            <h2 className="text-4xl md:text-5xl font-bold leading-tight">
                                Sovereignty for <span className="whitespace-nowrap">the Co-Op Alliances.</span>
                            </h2>
                            <p className="text-xl text-slate-400 font-light leading-relaxed">
                                We're proposing a new industrial contract. The <strong>GableLBM Industry Board</strong> provides the sovereign core for independent member alliances. Decisions don't come from Silicon Valley; they come from the yards that build America.
                            </p>

                            <div className="grid grid-cols-1 sm:grid-cols-2 gap-6 pt-8">
                                <div className="space-y-2 pb-4 sm:pb-0">
                                    <h4 className="flex items-center space-x-2 font-bold">
                                        <Database className="w-4 h-4 text-gable-blue" />
                                        <span>Member Editions</span>
                                    </h4>
                                    <p className="text-sm text-slate-500">Launch exclusive, branded versions of the platform for your members.</p>
                                </div>
                                <div className="space-y-2">
                                    <h4 className="flex items-center space-x-2 font-bold">
                                        <Cpu className="w-4 h-4 text-gable-blue" />
                                        <span>Governed Roadmap</span>
                                    </h4>
                                    <p className="text-sm text-slate-500">Every co-op gets a vote on the shared industry development roadmap.</p>
                                </div>
                            </div>
                        </motion.div>
                    </div>

                    <div className="lg:w-1/2 w-full">
                        <motion.div
                            initial={{ opacity: 0, scale: 0.95 }}
                            whileInView={{ opacity: 1, scale: 1 }}
                            viewport={{ once: true }}
                            className="glass-card p-2 md:p-8"
                        >
                            <div className="bg-gable-bg rounded-lg p-6 border border-white/5 font-mono text-xs sm:text-sm text-slate-500 overflow-hidden relative">
                                <div className="absolute top-0 right-0 p-4">
                                    <Database className="w-4 h-4 opacity-20" />
                                </div>
                                <pre className="whitespace-pre-wrap">
                                    {`// Gable Federated Context
{
  "tenant": "LBM_COOP_NORTH",
  "action": "SYNC_MASTER_CATALOG",
  "security": "FEDERATED_TRUST_L3",
  "modules": [
    "inventory.v1",
    "pricing.v1",
    "logistics.v1"
  ],
  "status": "PROPAGATING..."
}`}
                                </pre>
                                <div className="mt-8 space-y-3">
                                    <div className="h-1.5 w-full bg-white/5 rounded-full overflow-hidden">
                                        <motion.div
                                            className="h-full bg-gable-green"
                                            initial={{ width: 0 }}
                                            whileInView={{ width: "65%" }}
                                            transition={{ duration: 2, repeat: Infinity }}
                                        />
                                    </div>
                                    <div className="flex justify-between text-[10px] uppercase tracking-widest">
                                        <span>Alpha Node 01</span>
                                        <span className="text-gable-green text-opacity-80">Connected</span>
                                    </div>
                                </div>
                            </div>
                        </motion.div>
                    </div>
                </div>
            </div>
        </section>
    );
};
