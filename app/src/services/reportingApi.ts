import axios from 'axios';

const API_BASE = '/api/reporting';

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

export const reportingApi = {
  // Ad-hoc query preview
  previewReport: async (entityType: string, definition: ReportDefinition) => {
    const response = await axios.post(`${API_BASE}/builder/preview`, {
      entity_type: entityType,
      definition,
    });
    return response.data;
  },

  // Export
  exportReport: async (entityType: string, format: 'csv' | 'xlsx', definition: ReportDefinition) => {
    const response = await axios.post(`${API_BASE}/builder/export`, {
      entity_type: entityType,
      format,
      definition,
    }, {
      responseType: 'blob', // Important for file downloads
    });
    return response.data;
  },

  // Saved Reports CRUD
  listSavedReports: async (): Promise<SavedReport[]> => {
    const response = await axios.get(`${API_BASE}/saved`);
    return response.data;
  },

  getSavedReport: async (id: string): Promise<SavedReport> => {
    const response = await axios.get(`${API_BASE}/saved/${id}`);
    return response.data;
  },

  saveReport: async (report: Partial<SavedReport>): Promise<SavedReport> => {
    const response = await axios.post(`${API_BASE}/save`, report);
    return response.data;
  },
  
  updateSavedReport: async (id: string, report: Partial<SavedReport>) => {
    const response = await axios.put(`${API_BASE}/saved/${id}`, report);
    return response.data;
  },

  deleteSavedReport: async (id: string) => {
    await axios.delete(`${API_BASE}/saved/${id}`);
  }
};
