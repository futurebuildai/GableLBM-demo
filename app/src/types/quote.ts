import type { UOM } from "./product";

export type QuoteState = 'DRAFT' | 'SENT' | 'ACCEPTED' | 'REJECTED' | 'EXPIRED';

export interface QuoteLine {
    id: string;
    quote_id: string;
    product_id: string;
    sku: string;
    description: string;
    quantity: number;
    uom: UOM;
    unit_price: number;
    line_total: number;
    created_at: string;
}

export interface Quote {
    id: string;
    customer_id: string;
    customer_name?: string;
    job_id?: string;
    state: QuoteState;
    total_amount: number;
    expires_at?: string;
    created_at: string;
    updated_at: string;
    lines?: QuoteLine[];
}

// Payload for creating a quote
export interface CreateQuoteRequest {
    customer_id: string;
    job_id?: string;
    lines: Array<{
        product_id: string;
        sku: string;
        description: string;
        quantity: number;
        uom: UOM;
        unit_price: number;
    }>;
}
