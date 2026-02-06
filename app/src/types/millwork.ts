export interface MillworkOption {
    id: string;
    category: string;
    name: string;
    price_adjustment: number;
    attributes: Record<string, any>;
    created_at: string;
    updated_at: string;
}

export interface CreateOptionRequest {
    category: string;
    name: string;
    price_adjustment: number;
    attributes: Record<string, any>;
}

export interface MillworkConfiguration {
    doorType: MillworkOption | null;
    material: MillworkOption | null;
    glass: MillworkOption | null;
    width: number;
    height: number;
}
