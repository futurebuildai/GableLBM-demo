import React from 'react';
import { motion } from 'framer-motion';
import { BrandLogo } from './BrandLogo';

export const LoadingScreen: React.FC = () => {
    return (
        <div className="fixed inset-0 bg-[#0A0B10] flex flex-col items-center justify-center z-[100]">
            <motion.div
                initial={{ opacity: 0, scale: 0.8 }}
                animate={{ opacity: 1, scale: 1 }}
                transition={{ duration: 0.5 }}
                className="relative"
            >
                <div className="w-24 h-24 rounded-full border-4 border-white/5 border-t-gable-green animate-spin absolute inset-0" />
                <div className="w-24 h-24 flex items-center justify-center">
                    <BrandLogo variant="mark" size="lg" className="text-white" />
                </div>
            </motion.div>

            <motion.div
                initial={{ opacity: 0, y: 10 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ delay: 0.2 }}
                className="mt-8"
            >
                <BrandLogo variant="text" size="xl" />
            </motion.div>

            <motion.div
                initial={{ width: 0 }}
                animate={{ width: 100 }}
                transition={{ delay: 0.4, duration: 1.5, repeat: Infinity }}
                className="mt-4 h-1 bg-gable-green/50 rounded-full overflow-hidden"
            >
                <div className="h-full bg-gable-green w-full origin-left animate-progress" />
            </motion.div>

            <p className="mt-4 text-zinc-500 font-mono text-xs uppercase tracking-widest animate-pulse">
                Initializing System...
            </p>
        </div>
    );
};
