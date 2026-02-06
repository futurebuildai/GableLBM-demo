export type LocationType = 'ZONE' | 'AISLE' | 'RACK' | 'SHELF' | 'BIN' | 'YARD';

export interface Location {
    id: string;
    parent_id?: string;
    path: string;
    type: LocationType;
    code: string;
    description?: string;
    created_at: string;
    updated_at: string;
    children?: Location[];
}

export interface CreateLocationRequest {
    parent_id?: string;
    type: LocationType;
    code: string;
    description?: string;
}
