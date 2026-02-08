import React, { useState } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { cn } from '../../lib/utils';

interface TooltipProps {
    content: string;
    children: React.ReactNode;
    position?: 'top' | 'bottom' | 'left' | 'right';
    className?: string;
}

export const Tooltip: React.FC<TooltipProps> = ({
    content,
    children,
    position = 'top',
    className
}) => {
    const [isVisible, setIsVisible] = useState(false);

    const positions = {
        top: { bottom: '100%', left: '50%', x: '-50%', y: -8 },
        bottom: { top: '100%', left: '50%', x: '-50%', y: 8 },
        left: { right: '100%', top: '50%', y: '-50%', x: -8 },
        right: { left: '100%', top: '50%', y: '-50%', x: 8 }
    };

    return (
        <div
            className="relative inline-block"
            onMouseEnter={() => setIsVisible(true)}
            onMouseLeave={() => setIsVisible(false)}
        >
            {children}
            <AnimatePresence>
                {isVisible && (
                    <motion.div
                        initial={{ opacity: 0, scale: 0.95, ...positions[position] }}
                        animate={{ opacity: 1, scale: 1, ...positions[position] }}
                        exit={{ opacity: 0, scale: 0.95 }}
                        transition={{ duration: 0.15 }}
                        className={cn(
                            "absolute z-50 px-3 py-1.5 text-xs font-medium text-white bg-zinc-900 border border-white/10 rounded shadow-xl whitespace-nowrap pointer-events-none backdrop-blur-md",
                            className
                        )}
                        style={{
                            ...positions[position]
                        }}
                    >
                        {content}
                        {/* Arrow */}
                        <div
                            className={cn(
                                "absolute w-2 h-2 bg-zinc-900 border-white/10 rotate-45",
                                position === 'top' && "bottom-[-5px] left-1/2 -translate-x-1/2 border-b border-r",
                                position === 'bottom' && "top-[-5px] left-1/2 -translate-x-1/2 border-t border-l",
                                position === 'left' && "right-[-5px] top-1/2 -translate-y-1/2 border-t border-r",
                                position === 'right' && "left-[-5px] top-1/2 -translate-y-1/2 border-b border-l"
                            )}
                        />
                    </motion.div>
                )}
            </AnimatePresence>
        </div>
    );
};
