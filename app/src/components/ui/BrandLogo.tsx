import { cn } from "../../lib/utils";

interface BrandLogoProps {
    variant?: 'full' | 'mark' | 'text';
    className?: string;
    size?: 'sm' | 'md' | 'lg' | 'xl';
}

export function BrandLogo({ variant = 'full', className, size = 'md' }: BrandLogoProps) {
    const sizeClasses = {
        sm: "h-6",
        md: "h-8",
        lg: "h-12",
        xl: "h-16"
    };

    if (variant === 'mark') {
        return (
            <svg
                viewBox="0 0 64 64"
                fill="none"
                xmlns="http://www.w3.org/2000/svg"
                className={cn(sizeClasses[size], "w-auto", className)}
            >
                {/* Minimalist Gable Roof */}
                <path d="M4 36 L32 8 L60 36" stroke="currentColor" strokeWidth="6" strokeLinecap="round" strokeLinejoin="round" />

                {/* Lightbulb (Filament style) */}
                <path d="M32 16 V28" stroke="#00FFA3" strokeWidth="4" strokeLinecap="round" />
                <path d="M22 28 C22 28 22 46 32 46 C42 46 42 28 42 28" stroke="#00FFA3" strokeWidth="4" strokeLinecap="round" fill="none" />

                {/* Bulb Base */}
                <rect x="28" y="46" width="8" height="6" rx="1" fill="#00FFA3" />

                {/* Glow Effect (Subtle) */}
                <circle cx="32" cy="32" r="16" fill="#00FFA3" fillOpacity="0.1" />
            </svg>
        );
    }

    if (variant === 'text') {
        return (
            <div className={cn("flex items-center gap-0.5", className)}>
                <span className="font-bold tracking-tight text-white">Gable</span>
                <span className="font-light tracking-widest text-gable-green">LBM</span>
            </div>
        );
    }

    // Full Variant
    return (
        <div className={cn("flex items-center gap-3 select-none", className)}>
            <div className={cn("relative shrink-0 flex items-center justify-center",
                size === 'sm' && "w-6 h-6",
                size === 'md' && "w-8 h-8",
                size === 'lg' && "w-12 h-12",
                size === 'xl' && "w-16 h-16"
            )}>
                <BrandLogo variant="mark" className="w-full h-full text-white" />
            </div>
            <div className={cn("flex items-baseline leading-none",
                size === 'sm' && "text-lg",
                size === 'md' && "text-xl",
                size === 'lg' && "text-3xl",
                size === 'xl' && "text-4xl"
            )}>
                <span className="font-bold tracking-tight text-white">Gable</span>
                <span className="font-light tracking-widest text-gable-green ml-0.5">LBM</span>
            </div>
        </div>
    );
}
