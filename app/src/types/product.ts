export type UOM =
    | 'PCS'
    | 'EA'
    | 'LF'
    | 'SF'
    | 'BF'
    | 'MBF'
    | 'SQ'
    | 'BOX'
    | 'CTN'
    | 'RL'
    | 'GAL'
    | 'LBS'
    | 'BAG'
    | 'BUNDLE'
    | 'PAIR'
    | 'SET';

export interface Product {
    id: string;
    sku: string;
    description: string;
    uom_primary: UOM;
    total_quantity?: number;
    total_allocated?: number;
    created_at: string;
    updated_at: string;
}

export interface Inventory {
    id: string;
    product_id: string;
    location: string;
    quantity: number;
    allocated?: number;
    updated_at: string;
}
