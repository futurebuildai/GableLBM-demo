import { useState } from 'react';
import type { BlueprintScanResponse } from '../../types/configurator';
import { VisionService } from '../../services/VisionService';
import { motion, AnimatePresence } from 'framer-motion';
import { Upload, AlertTriangle, CheckCircle, Eye, FileText, Zap } from 'lucide-react';

const SAMPLE_BLUEPRINT = `STRUCTURAL FRAMING PLAN — LOT 42
Wall framing: 2x6 SYP #2 studs
Stud height: 10' stud walls
Spacing: 16" OC
Headers: 4x12 Douglas Fir Select Structural
Treated sill plate: 2x6 PT SYP
Roof rafters: 2x8 SPF #2 @ 24" OC`;

export const BlueprintVerifier = () => {
    const [blueprintText, setBlueprintText] = useState('');
    const [configSelections, setConfigSelections] = useState<Record<string, string>>({
        Species: 'Douglas Fir',
        Grade: '#2',
        Treatment: 'None',
        Dimensions: '2x4-8',
    });
    const [scanResult, setScanResult] = useState<BlueprintScanResponse | null>(null);
    const [scanning, setScanning] = useState(false);
    const [scanError, setScanError] = useState<string | null>(null);
    const [dragOver, setDragOver] = useState(false);

    const handleScan = async () => {
        if (!blueprintText.trim()) return;
        setScanning(true);
        setScanError(null);
        try {
            const result = await VisionService.scanBlueprint(blueprintText, configSelections);
            setScanResult(result);
        } catch (err) {
            const message = err instanceof Error ? err.message : 'Blueprint scan failed';
            setScanError(message);
        } finally {
            setScanning(false);
        }
    };

    const loadSample = () => {
        setBlueprintText(SAMPLE_BLUEPRINT);
    };

    return (
        <div className="min-h-[calc(100vh-6rem)]">
            {/* Header */}
            <div className="mb-8">
                <div className="flex items-center gap-3 mb-2">
                    <div className="w-10 h-10 rounded-lg bg-blue-500/20 border border-blue-500/30 flex items-center justify-center">
                        <Eye size={20} className="text-blue-400" />
                    </div>
                    <div>
                        <h1 className="text-3xl font-bold text-white">Blueprint Verifier</h1>
                        <span className="text-xs font-mono bg-amber-500/20 text-amber-400 px-2 py-0.5 rounded border border-amber-500/30">
                            AI PROTOTYPE
                        </span>
                    </div>
                </div>
                <p className="text-gray-400 mt-1">
                    Upload blueprint specs and compare against your configurator selections to identify mismatches
                </p>
            </div>

            <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
                {/* Left: Input Panel */}
                <div className="space-y-6">
                    {/* Blueprint Text Input */}
                    <div className="bg-[#161821] border border-white/10 rounded-xl p-6">
                        <div className="flex items-center justify-between mb-4">
                            <h3 className="font-semibold text-white flex items-center gap-2">
                                <FileText size={18} className="text-blue-400" />
                                Blueprint Specifications
                            </h3>
                            <button
                                onClick={loadSample}
                                className="text-xs text-[#00FFA3] hover:text-[#00FFA3]/80 font-medium transition-colors"
                            >
                                Load Sample →
                            </button>
                        </div>

                        {/* Drag & Drop Zone */}
                        <div
                            onDragOver={(e) => { e.preventDefault(); setDragOver(true); }}
                            onDragLeave={() => setDragOver(false)}
                            onDrop={(e) => {
                                e.preventDefault();
                                setDragOver(false);
                                // Simulated: just show the drop area acknowledgment
                            }}
                            className={`border-2 border-dashed rounded-xl p-4 mb-4 text-center transition-all ${dragOver
                                ? 'border-[#00FFA3] bg-[#00FFA3]/5'
                                : 'border-white/10 hover:border-white/20'
                                }`}
                        >
                            <Upload size={24} className="mx-auto text-gray-500 mb-2" />
                            <div className="text-sm text-gray-400">
                                Drop a PDF or paste blueprint text below
                            </div>
                            <div className="text-xs text-gray-600 mt-1">
                                AI extraction from PDFs is a prototype — paste text for best results
                            </div>
                        </div>

                        <textarea
                            value={blueprintText}
                            onChange={(e) => setBlueprintText(e.target.value)}
                            placeholder="Paste blueprint specifications here..."
                            rows={10}
                            className="w-full bg-[#0A0B10] border border-white/10 rounded-lg p-4 text-sm font-mono text-gray-300 
                                       focus:border-blue-500/50 focus:ring-1 focus:ring-blue-500/20 outline-none resize-none"
                        />
                    </div>

                    {/* Current Config Selections */}
                    <div className="bg-[#161821] border border-white/10 rounded-xl p-6">
                        <h3 className="font-semibold text-white mb-4 flex items-center gap-2">
                            <Zap size={18} className="text-amber-400" />
                            Configurator Selections (to compare against)
                        </h3>
                        <div className="grid grid-cols-2 gap-3">
                            {Object.entries(configSelections).map(([key, value]) => (
                                <div key={key}>
                                    <label className="text-xs text-gray-500 block mb-1">{key}</label>
                                    <input
                                        value={value}
                                        onChange={(e) => setConfigSelections(prev => ({ ...prev, [key]: e.target.value }))}
                                        className="w-full bg-[#0A0B10] border border-white/10 rounded-lg px-3 py-2 text-sm text-white
                                                   focus:border-[#00FFA3]/50 outline-none"
                                    />
                                </div>
                            ))}
                        </div>
                    </div>

                    {/* Scan Button */}
                    <button
                        onClick={handleScan}
                        disabled={!blueprintText.trim() || scanning}
                        className={`w-full py-4 rounded-xl font-bold text-lg flex items-center justify-center gap-3 transition-all ${blueprintText.trim() && !scanning
                            ? 'bg-gradient-to-r from-blue-500 to-indigo-500 text-white hover:shadow-[0_0_30px_rgba(59,130,246,0.3)]'
                            : 'bg-gray-800 text-gray-600 cursor-not-allowed'
                            }`}
                    >
                        {scanning ? (
                            <>
                                <div className="w-5 h-5 border-2 border-white border-t-transparent rounded-full animate-spin" />
                                Analyzing Blueprint...
                            </>
                        ) : (
                            <>
                                <Eye size={20} />
                                Scan & Compare
                            </>
                        )}
                    </button>

                    {scanError && (
                        <div className="bg-red-500/10 border border-red-500/30 rounded-xl p-4 flex items-center gap-3">
                            <AlertTriangle size={18} className="text-red-400 shrink-0" />
                            <span className="text-sm text-red-300">{scanError}</span>
                        </div>
                    )}
                </div>

                {/* Right: Results Panel */}
                <div className="space-y-6">
                    <AnimatePresence mode="wait">
                        {scanResult ? (
                            <motion.div
                                initial={{ opacity: 0, y: 20 }}
                                animate={{ opacity: 1, y: 0 }}
                                exit={{ opacity: 0 }}
                                className="space-y-6"
                            >
                                {/* Summary */}
                                <div className={`p-6 rounded-xl border ${scanResult.mismatches.length === 0
                                    ? 'bg-emerald-500/5 border-emerald-500/30'
                                    : 'bg-amber-500/5 border-amber-500/30'
                                    }`}>
                                    <div className="flex items-center gap-3 mb-2">
                                        {scanResult.mismatches.length === 0 ? (
                                            <CheckCircle size={24} className="text-emerald-400" />
                                        ) : (
                                            <AlertTriangle size={24} className="text-amber-400" />
                                        )}
                                        <span className="font-semibold text-white">{scanResult.summary}</span>
                                    </div>
                                </div>

                                {/* Extracted Dimensions */}
                                <div className="bg-[#161821] border border-white/10 rounded-xl p-6">
                                    <h3 className="font-semibold text-white mb-4">Extracted Dimensions</h3>
                                    <div className="space-y-2">
                                        {Object.entries(scanResult.extracted_dimensions).map(([key, value]) => (
                                            <div key={key} className="flex justify-between items-center py-2 border-b border-white/5 last:border-0">
                                                <span className="text-gray-400 text-sm capitalize">
                                                    {key.replace(/_/g, ' ')}
                                                </span>
                                                <span className="font-mono text-white bg-blue-500/10 px-3 py-1 rounded border border-blue-500/20">
                                                    {value}
                                                </span>
                                            </div>
                                        ))}
                                    </div>
                                </div>

                                {/* Mismatches */}
                                {scanResult.mismatches.length > 0 && (
                                    <div className="bg-[#161821] border border-white/10 rounded-xl p-6">
                                        <h3 className="font-semibold text-white mb-4 flex items-center gap-2">
                                            <AlertTriangle size={18} className="text-amber-400" />
                                            Mismatches ({scanResult.mismatches.length})
                                        </h3>
                                        <div className="space-y-3">
                                            {scanResult.mismatches.map((m, i) => (
                                                <motion.div
                                                    key={i}
                                                    initial={{ opacity: 0, x: 20 }}
                                                    animate={{ opacity: 1, x: 0 }}
                                                    transition={{ delay: i * 0.1 }}
                                                    className={`p-4 rounded-lg border ${m.severity === 'error'
                                                        ? 'bg-red-500/5 border-red-500/30'
                                                        : 'bg-amber-500/5 border-amber-500/30'
                                                        }`}
                                                >
                                                    <div className="flex items-center gap-2 mb-2">
                                                        <span className={`text-xs font-bold uppercase px-2 py-0.5 rounded ${m.severity === 'error'
                                                            ? 'bg-red-500/20 text-red-400'
                                                            : 'bg-amber-500/20 text-amber-400'
                                                            }`}>
                                                            {m.severity}
                                                        </span>
                                                        <span className="font-medium text-white">{m.field}</span>
                                                    </div>
                                                    <div className="text-sm text-gray-300 mb-3">{m.message}</div>
                                                    <div className="grid grid-cols-2 gap-3 text-sm">
                                                        <div className="bg-[#0A0B10] rounded p-2">
                                                            <div className="text-xs text-gray-500 mb-1">Blueprint</div>
                                                            <div className="font-mono text-blue-400">{m.blueprint_value}</div>
                                                        </div>
                                                        <div className="bg-[#0A0B10] rounded p-2">
                                                            <div className="text-xs text-gray-500 mb-1">Configurator</div>
                                                            <div className="font-mono text-amber-400">{m.config_value}</div>
                                                        </div>
                                                    </div>
                                                </motion.div>
                                            ))}
                                        </div>
                                    </div>
                                )}
                            </motion.div>
                        ) : (
                            <motion.div
                                initial={{ opacity: 0 }}
                                animate={{ opacity: 1 }}
                                className="bg-[#161821] border border-white/10 rounded-xl p-12 text-center"
                            >
                                <div className="w-16 h-16 rounded-full bg-blue-500/10 border border-blue-500/20 flex items-center justify-center mx-auto mb-4">
                                    <Eye size={32} className="text-blue-400/50" />
                                </div>
                                <h3 className="text-lg font-semibold text-gray-400 mb-2">No Scan Results</h3>
                                <p className="text-sm text-gray-600">
                                    Paste blueprint text and click "Scan & Compare" to see mismatch analysis
                                </p>
                            </motion.div>
                        )}
                    </AnimatePresence>
                </div>
            </div>
        </div>
    );
};
