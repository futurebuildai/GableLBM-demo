import type {
    Project,
    ProjectDashboardDTO,
    CreateProjectRequest,
    UpdateProjectRequest
} from '../types/project';

const API_URL = '';
const TOKEN_KEY = 'portal_token';

async function fetchWithAuth<T>(url: string, options: RequestInit = {}): Promise<T> {
    const token = localStorage.getItem(TOKEN_KEY);
    const headers: Record<string, string> = {
        'Content-Type': 'application/json',
        ...(options.headers as Record<string, string> || {}),
    };
    if (token) {
        headers['Authorization'] = `Bearer ${token}`;
    }

    const response = await fetch(url, { ...options, headers });
    if (!response.ok) {
        const text = await response.text().catch(() => response.statusText);
        throw new Error(`API error: ${response.status} ${text}`);
    }

    return await response.json() as T;
}

export const ProjectService = {
    async listProjects(): Promise<Project[]> {
        return fetchWithAuth<Project[]>(`${API_URL}/api/portal/v1/projects`);
    },

    async getProjectDashboard(id: string): Promise<ProjectDashboardDTO> {
        return fetchWithAuth<ProjectDashboardDTO>(`${API_URL}/api/portal/v1/projects/${id}`);
    },

    async createProject(req: CreateProjectRequest): Promise<Project> {
        return fetchWithAuth<Project>(`${API_URL}/api/portal/v1/projects`, {
            method: 'POST',
            body: JSON.stringify(req),
        });
    },

    async updateProject(id: string, req: UpdateProjectRequest): Promise<Project> {
        return fetchWithAuth<Project>(`${API_URL}/api/portal/v1/projects/${id}`, {
            method: 'PUT',
            body: JSON.stringify(req),
        });
    },
};
