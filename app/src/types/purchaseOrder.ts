export type POStatus = 'DRAFT' | 'SENT' | 'PARTIAL' | 'RECEIVED' | 'CANCELLED';

export interface PurchaseOrder {
    id: string;
    vendor_id?: string;
    vendor_name?: string;
    status: POStatus;
    created_at: string;
    updated_at: string;
    lines?: PurchaseOrderLine[];
    line_count?: number;
    total_cost?: number;
}

export interface PurchaseOrderLine {
    id: string;
    po_id: string;
    product_id?: string;
    description: string;
    quantity: number;
    qty_received: number;
    cost: number;
    linked_so_line_id?: string;
}

export interface CreatePORequest {
    vendor_id: string;
    lines: CreatePOLine[];
}

export interface CreatePOLine {
    product_id: string;
    description: string;
    quantity: number;
    cost: number;
}

export interface ReceivePORequest {
    lines: ReceiveLine[];
}

export interface ReceiveLine {
    line_id: string;
    qty_received: number;
    location_id: string;
}
