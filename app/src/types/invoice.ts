export type InvoiceStatus = 'UNPAID' | 'PARTIAL' | 'PAID' | 'VOID' | 'OVERDUE';

export interface Invoice {
    id: string;
    order_id: string;
    customer_id: string;
    status: InvoiceStatus;
    total_amount: number;
    due_date?: string;
    paid_at?: string;
    created_at: string;
    updated_at: string;

    // Relations
    lines?: InvoiceLine[];
}

export interface InvoiceLine {
    id: string;
    invoice_id: string;
    product_id: string;
    quantity: number;
    price_each: number;
    created_at: string;
}
