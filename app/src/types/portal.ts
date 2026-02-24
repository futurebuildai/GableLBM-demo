// Portal TypeScript types — zero `any`

export interface PortalUser {
    id: string;
    customer_id: string;
    email: string;
    name: string;
    role: string;
    created_at: string;
    updated_at: string;
}

export interface PortalConfig {
    id: string;
    dealer_name: string;
    logo_url: string;
    primary_color: string;
    support_email: string;
    support_phone: string;
}

export interface PortalLoginResponse {
    token: string;
    user: PortalUser;
    config: PortalConfig;
}

export interface PortalDashboard {
    balance_due: number;
    credit_limit: number;
    past_due: number;
    recent_orders: PortalOrder[];
}

export interface PortalOrder {
    id: string;
    status: string;
    total_amount: number;
    created_at: string;
    lines: PortalLineItem[];
}

export interface PortalLineItem {
    product_id: string;
    product_sku: string;
    product_name: string;
    quantity: number;
    price_each: number;
}

export interface PortalInvoice {
    id: string;
    order_id: string;
    status: string;
    total_amount: number;
    subtotal: number;
    tax_amount: number;
    payment_terms: string;
    due_date: string | null;
    paid_at: string | null;
    created_at: string;
    lines: PortalLineItem[];
}

export interface PortalDelivery {
    id: string;
    order_id: string;
    status: string;
    pod_proof_url: string | null;
    pod_signed_by: string | null;
    pod_timestamp: string | null;
    created_at: string;
    order_number: string | null;
}

export interface ReorderResponse {
    order_id: string;
    message: string;
}

// --- Catalog Types (Sprint 27) ---

export interface CatalogProduct {
    id: string;
    sku: string;
    name: string;
    category: string;
    species: string;
    grade: string;
    image_url: string;
    uom: string;
    base_price: number;
    customer_price: number;
    price_source: string;
    available: number;
    in_stock: boolean;
}

export interface CatalogDetail extends CatalogProduct {
    weight_lbs: number;
    upc: string;
    vendor: string;
}

// --- Cart Types (Sprint 27) ---

export interface CartItem {
    id: string;
    product_id: string;
    product_sku: string;
    product_name: string;
    image_url: string;
    quantity: number;
    unit_price: number;
    line_total: number;
    available: number;
}

export interface Cart {
    id: string;
    items: CartItem[];
    item_count: number;
    subtotal: number;
}

export interface CheckoutRequest {
    delivery_method: 'DELIVERY' | 'PICKUP';
    delivery_address: string;
    payment_method: 'ACCOUNT' | 'CARD';
    notes: string;
}

export interface CheckoutResponse {
    order_id: string;
    message: string;
}

