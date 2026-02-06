export type PaymentMethod = 'CASH' | 'CARD' | 'CHECK' | 'ACCOUNT';

export interface Payment {
    id: string;
    invoice_id: string;
    amount: number;
    method: PaymentMethod;
    reference: string;
    notes: string;
    created_at: string;
}

export interface CreatePaymentRequest {
    invoice_id: string;
    amount: number;
    method: PaymentMethod;
    reference: string;
    notes: string;
}
