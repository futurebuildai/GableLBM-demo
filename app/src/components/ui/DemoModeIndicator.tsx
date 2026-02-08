import React from 'react';
import { Beaker } from 'lucide-react';

export const DemoModeIndicator: React.FC = () => {
    return (
        <div className="fixed bottom-4 right-4 z-50 flex items-center gap-2 px-3 py-1.5 rounded-full bg-gable-green/10 border border-gable-green/20 backdrop-blur-md shadow-lg pointer-events-none">
            <Beaker className="w-3 h-3 text-gable-green animate-pulse" />
            <span className="text-[10px] font-mono font-bold uppercase tracking-wider text-gable-green">
                Demo Mode With Sample Data
            </span>
        </div>
    );
};
