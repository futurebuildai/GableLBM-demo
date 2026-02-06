export const Footer = () => {
    return (
        <footer className="py-12 border-t border-white/5 px-4">
            <div className="container mx-auto">
                <div className="flex flex-col md:flex-row justify-between items-center gap-8">
                    <div className="flex items-center space-x-2">
                        <img src="/logo.png" alt="FutureBuild Logo" className="h-6 w-auto grayscale" />
                        <span className="font-bold text-lg tracking-tight text-white/60">GableLBM by FutureBuild AI</span>
                    </div>

                    <div className="flex space-x-8 text-sm text-slate-500 font-mono uppercase tracking-tighter">
                        <a href="#" className="hover:text-gable-green transition-colors">Discord</a>
                        <a href="#" className="hover:text-gable-green transition-colors">Docs</a>
                        <a href="#" className="hover:text-gable-green transition-colors">Legal</a>
                    </div>

                    <div className="text-xs text-slate-600 font-mono">
                        © 2026 FUTUREBUILD AI. A DELAWARE C-CORP. OPEN SOURCE REVOLUTION.
                    </div>
                </div>
            </div>
        </footer>
    );
};
