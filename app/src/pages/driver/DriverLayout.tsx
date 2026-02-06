import { Outlet } from "react-router-dom";

export function DriverLayout() {
    return (
        <div className="min-h-screen bg-[#0A0B10] text-white font-sans md:max-w-md md:mx-auto md:border-x md:border-white/10 relative shadow-2xl">
            <header className="h-16 flex items-center justify-between px-4 border-b border-white/10 bg-[#161821]/80 backdrop-blur-md sticky top-0 z-50">
                <div className="font-bold text-lg tracking-wider font-mono">
                    GABLE<span className="text-[#00FFA3]">DRIVER</span>
                </div>
                <div className="h-8 w-8 rounded-full bg-white/10 flex items-center justify-center text-xs font-mono">
                    D1
                </div>
            </header>

            <main className="pb-8">
                <Outlet />
            </main>
        </div>
    );
}
