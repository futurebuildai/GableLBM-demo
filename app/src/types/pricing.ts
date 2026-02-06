export type PricingSource = "CONTRACT" | "TIER" | "RETAIL";

export interface CalculatedPrice {
    product_id: string;
    original_price: number;
    final_price: number;
    discount_pct: number;
    source: PricingSource;
    details: string;
}
