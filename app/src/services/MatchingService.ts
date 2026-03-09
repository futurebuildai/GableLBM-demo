import type { MatchResult, MatchException, MatchConfig, UpdateMatchConfigRequest } from '../types/matching';

const API = 'https://backend-production-bdf8.up.railway.app';

export async function runMatch(poId: string): Promise<MatchResult> {
    const res = await fetch(`${API}/api/matching/run/${poId}`, { method: 'POST' });
    if (!res.ok) throw new Error(await res.text());
    return res.json();
}

export async function getMatchResult(poId: string): Promise<MatchResult> {
    const res = await fetch(`${API}/api/matching/results/${poId}`);
    if (!res.ok) throw new Error('Failed to fetch match result');
    return res.json();
}

export async function listExceptions(): Promise<MatchException[]> {
    const res = await fetch(`${API}/api/matching/exceptions`);
    if (!res.ok) throw new Error('Failed to fetch exceptions');
    return res.json();
}

export async function getMatchConfig(): Promise<MatchConfig> {
    const res = await fetch(`${API}/api/matching/config`);
    if (!res.ok) throw new Error('Failed to fetch match config');
    return res.json();
}

export async function updateMatchConfig(data: UpdateMatchConfigRequest): Promise<MatchConfig> {
    const res = await fetch(`${API}/api/matching/config`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(data),
    });
    if (!res.ok) throw new Error(await res.text());
    return res.json();
}
