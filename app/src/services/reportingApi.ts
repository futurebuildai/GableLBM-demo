const API_BASE = 'https://backend-production-bdf8.up.railway.app';

export interface ReportDefinition {
  columns: ReportColumn[];
  filters: ReportFilter[];
  groupings: ReportGrouping[];
}

export interface ReportColumn {
  field: string;
  label: string;
  aggregation?: string;
}

export interface ReportFilter {
  field: string;
  operator: string;
  value: any;
}

export interface ReportGrouping {
  field: string;
}

export interface SavedReport {
  id: string;
  name: string;
  description: string;
  entity_type: string;
  definition_json: ReportDefinition;
  created_at: string;
}

export interface ReportSchedule {
  id: string;
  report_id: string;
  cron_expression: string;
  recipients: string[];
  status: string;
}

const fetchJson = async <T>(url: string, options?: RequestInit): Promise<T> => {
  const token = localStorage.getItem('token');
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
    ...(options?.headers as Record<string, string> || {}),
  };
  if (token) headers['Authorization'] = `Bearer ${token}`;

  const response = await fetch(url, { ...options, headers });
  if (!response.ok) throw new Error(`API Error: ${response.status}`);
  return response.json();
};

export const reportingApi = {
  // Ad-hoc query preview
  previewReport: async (entityType: string, definition: ReportDefinition) => {
    return fetchJson<any>(`${API_BASE}/builder/preview`, {
      method: 'POST',
      body: JSON.stringify({ entity_type: entityType, definition })
    });
  },

  // Export
  exportReport: async (entityType: string, format: 'csv' | 'xlsx', definition: ReportDefinition) => {
    const token = localStorage.getItem('token');
    const headers: Record<string, string> = {
      'Content-Type': 'application/json'
    };
    if (token) headers['Authorization'] = `Bearer ${token}`;

    const response = await fetch(`${API_BASE}/builder/export`, {
      method: 'POST',
      headers,
      body: JSON.stringify({ entity_type: entityType, format, definition })
    });
    if (!response.ok) throw new Error(`Export failed`);
    return response.blob();
  },

  // Saved Reports CRUD
  listSavedReports: async (): Promise<SavedReport[]> => {
    return fetchJson<SavedReport[]>(`${API_BASE}/saved`);
  },

  getSavedReport: async (id: string): Promise<SavedReport> => {
    return fetchJson<SavedReport>(`${API_BASE}/saved/${id}`);
  },

  saveReport: async (report: Partial<SavedReport>): Promise<SavedReport> => {
    return fetchJson<SavedReport>(`${API_BASE}/save`, {
      method: 'POST',
      body: JSON.stringify(report)
    });
  },

  updateSavedReport: async (id: string, report: Partial<SavedReport>) => {
    return fetchJson<SavedReport>(`${API_BASE}/saved/${id}`, {
      method: 'PUT',
      body: JSON.stringify(report)
    });
  },

  deleteSavedReport: async (id: string) => {
    const token = localStorage.getItem('token');
    const headers: Record<string, string> = {};
    if (token) headers['Authorization'] = `Bearer ${token}`;

    const response = await fetch(`${API_BASE}/saved/${id}`, {
      method: 'DELETE',
      headers
    });
    if (!response.ok) throw new Error(`Delete failed`);
  }
};
