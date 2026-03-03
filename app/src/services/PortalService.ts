import type {
    PortalLoginResponse,
    PortalUser,
    PortalConfig,
    PortalDashboard,
    PortalOrder,
    PortalInvoice,
    PortalDelivery,
    ReorderResponse,
    CatalogProduct,
    CatalogDetail,
    Cart,
    CheckoutRequest,
    CheckoutResponse,
    PortalInvite,
    InviteUserRequest,
    UpdateUserRoleRequest,
    UpdateUserStatusRequest,
} from '../types/portal';

const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080';
const REQUEST_TIMEOUT_MS = 10_000;
const MAX_RETRIES = 1;
const RETRY_DELAY_MS = 2_000;
const TOKEN_KEY = 'portal_token';

/**
 * Get stored portal JWT token.
 */
function getToken(): string | null {
    return localStorage.getItem(TOKEN_KEY);
}

/**
 * Store portal JWT token.
 */
export function setToken(token: string): void {
    localStorage.setItem(TOKEN_KEY, token);
}

/**
 * Remove stored portal JWT token.
 */
export function clearToken(): void {
    localStorage.removeItem(TOKEN_KEY);
}

/**
 * Check if user is authenticated.
 */
export function isAuthenticated(): boolean {
    return getToken() !== null;
}

/**
 * Fetch with timeout, retry, and auth header injection.
 */
async function fetchWithRetry<T>(
    url: string,
    options: RequestInit = {},
    retries = MAX_RETRIES,
): Promise<T> {
    let lastError: Error | null = null;

    for (let attempt = 0; attempt <= retries; attempt++) {
        const controller = new AbortController();
        const timeoutId = setTimeout(() => controller.abort(), REQUEST_TIMEOUT_MS);

        try {
            const token = getToken();
            const headers: Record<string, string> = {
                'Content-Type': 'application/json',
                ...(options.headers as Record<string, string> || {}),
            };
            if (token) {
                headers['Authorization'] = `Bearer ${token}`;
            }

            const response = await fetch(url, {
                ...options,
                headers,
                signal: controller.signal,
            });
            clearTimeout(timeoutId);

            if (!response.ok) {
                const text = await response.text().catch(() => response.statusText);
                throw new Error(`API error: ${response.status} ${text}`);
            }

            return await response.json() as T;
        } catch (error) {
            clearTimeout(timeoutId);
            lastError = error instanceof Error ? error : new Error(String(error));

            if (lastError.name === 'AbortError') {
                lastError = new Error('Request timed out');
            }

            if (attempt < retries) {
                await new Promise(resolve => setTimeout(resolve, RETRY_DELAY_MS));
            }
        }
    }

    throw lastError ?? new Error('Request failed');
}

export const PortalService = {
    /** Authenticate contractor and return JWT + user + config. */
    async login(email: string, password: string): Promise<PortalLoginResponse> {
        return fetchWithRetry<PortalLoginResponse>(
            `${API_URL}/api/portal/v1/login`,
            { method: 'POST', body: JSON.stringify({ email, password }) },
            0,
        );
    },

    /** Get portal branding config (public). */
    async getConfig(): Promise<PortalConfig> {
        return fetchWithRetry<PortalConfig>(`${API_URL}/api/portal/v1/config`);
    },

    /** Get contractor dashboard data. */
    async getDashboard(): Promise<PortalDashboard> {
        return fetchWithRetry<PortalDashboard>(`${API_URL}/api/portal/v1/dashboard`);
    },

    /** List order history. */
    async getOrders(): Promise<PortalOrder[]> {
        return fetchWithRetry<PortalOrder[]>(`${API_URL}/api/portal/v1/orders`);
    },

    /** Get single order with lines. */
    async getOrder(id: string): Promise<PortalOrder> {
        return fetchWithRetry<PortalOrder>(`${API_URL}/api/portal/v1/orders/${id}`);
    },

    /** Create a reorder from historical order. */
    async reorder(orderId: string): Promise<ReorderResponse> {
        return fetchWithRetry<ReorderResponse>(
            `${API_URL}/api/portal/v1/orders/reorder`,
            { method: 'POST', body: JSON.stringify({ order_id: orderId }) },
        );
    },

    /** List invoices. */
    async getInvoices(): Promise<PortalInvoice[]> {
        return fetchWithRetry<PortalInvoice[]>(`${API_URL}/api/portal/v1/invoices`);
    },

    /** Get single invoice with lines. */
    async getInvoice(id: string): Promise<PortalInvoice> {
        return fetchWithRetry<PortalInvoice>(`${API_URL}/api/portal/v1/invoices/${id}`);
    },

    /** List deliveries with POD info. */
    async getDeliveries(): Promise<PortalDelivery[]> {
        return fetchWithRetry<PortalDelivery[]>(`${API_URL}/api/portal/v1/deliveries`);
    },

    /** Get single delivery with POD info. */
    async getDelivery(id: string): Promise<PortalDelivery> {
        return fetchWithRetry<PortalDelivery>(`${API_URL}/api/portal/v1/deliveries/${id}`);
    },

    // --- Catalog Methods (Sprint 27) ---

    /** Browse product catalog with optional filters. */
    async getCatalog(params?: {
        q?: string;
        category?: string;
        species?: string;
        grade?: string;
    }): Promise<CatalogProduct[]> {
        const searchParams = new URLSearchParams();
        if (params?.q) searchParams.set('q', params.q);
        if (params?.category) searchParams.set('category', params.category);
        if (params?.species) searchParams.set('species', params.species);
        if (params?.grade) searchParams.set('grade', params.grade);
        const qs = searchParams.toString();
        return fetchWithRetry<CatalogProduct[]>(
            `${API_URL}/api/portal/v1/catalog${qs ? `?${qs}` : ''}`,
        );
    },

    /** Get single catalog product detail. */
    async getCatalogProduct(id: string): Promise<CatalogDetail> {
        return fetchWithRetry<CatalogDetail>(`${API_URL}/api/portal/v1/catalog/${id}`);
    },

    // --- Cart Methods (Sprint 27) ---

    /** Get current shopping cart. */
    async getCart(): Promise<Cart> {
        return fetchWithRetry<Cart>(`${API_URL}/api/portal/v1/cart`);
    },

    /** Add item to cart. */
    async addToCart(productId: string, quantity: number): Promise<Cart> {
        return fetchWithRetry<Cart>(
            `${API_URL}/api/portal/v1/cart/items`,
            { method: 'POST', body: JSON.stringify({ product_id: productId, quantity }) },
        );
    },

    /** Update cart item quantity. */
    async updateCartItem(itemId: string, quantity: number): Promise<Cart> {
        return fetchWithRetry<Cart>(
            `${API_URL}/api/portal/v1/cart/items/${itemId}`,
            { method: 'PUT', body: JSON.stringify({ quantity }) },
        );
    },

    /** Remove item from cart. */
    async removeCartItem(itemId: string): Promise<Cart> {
        return fetchWithRetry<Cart>(
            `${API_URL}/api/portal/v1/cart/items/${itemId}`,
            { method: 'DELETE' },
        );
    },

    /** Place order from cart. */
    async checkout(req: CheckoutRequest): Promise<CheckoutResponse> {
        return fetchWithRetry<CheckoutResponse>(
            `${API_URL}/api/portal/v1/checkout`,
            { method: 'POST', body: JSON.stringify(req) },
        );
    },

    // --- User Management Methods (Sprint 34) ---

    /** Get portal users. */
    async getUsers(): Promise<PortalUser[]> {
        return fetchWithRetry<PortalUser[]>(`${API_URL}/api/portal/v1/users`);
    },

    /** Get active invites. */
    async getInvites(): Promise<PortalInvite[]> {
        return fetchWithRetry<PortalInvite[]>(`${API_URL}/api/portal/v1/invites`);
    },

    /** Invite a new user. */
    async inviteUser(req: InviteUserRequest): Promise<PortalInvite> {
        return fetchWithRetry<PortalInvite>(
            `${API_URL}/api/portal/v1/invites`,
            { method: 'POST', body: JSON.stringify(req) },
        );
    },

    /** Update a user's role. */
    async updateUserRole(id: string, req: UpdateUserRoleRequest): Promise<void> {
        return fetchWithRetry<void>(
            `${API_URL}/api/portal/v1/users/${id}/role`,
            { method: 'PUT', body: JSON.stringify(req) },
        );
    },

    /** Update a user's status. */
    async updateUserStatus(id: string, req: UpdateUserStatusRequest): Promise<void> {
        return fetchWithRetry<void>(
            `${API_URL}/api/portal/v1/users/${id}/status`,
            { method: 'PUT', body: JSON.stringify(req) },
        );
    },
};
