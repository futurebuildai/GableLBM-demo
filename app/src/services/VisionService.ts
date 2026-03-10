import type { BlueprintScanResponse } from '../types/configurator';

const API_URL = 'http://localhost:8080';

export const VisionService = {
    async scanBlueprint(
        blueprintText: string,
        configSelections: Record<string, string>
    ): Promise<BlueprintScanResponse> {
        const response = await fetch(`${API_URL}/api/vision/scan`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
                blueprint_text: blueprintText,
                config_selections: configSelections,
            }),
        });
        if (!response.ok) throw new Error('Blueprint scan failed');
        return response.json();
    },
};
