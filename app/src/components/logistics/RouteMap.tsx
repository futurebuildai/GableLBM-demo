import React, { useEffect, useMemo } from 'react';
import { MapContainer, TileLayer, Marker, Popup, Polyline, useMap } from 'react-leaflet';
import 'leaflet/dist/leaflet.css';
import type { Delivery } from '../../types/delivery';
import L from 'leaflet';

function makeStopIcon(index: number, status: string) {
    const isDelivered = status === 'DELIVERED';
    const isFailed = status === 'FAILED' || status === 'PARTIAL';
    const bg = isDelivered ? '#10b981' : isFailed ? '#ef4444' : '#00FFA3';
    const text = isDelivered || isFailed ? '#fff' : '#000';
    const shadow = isDelivered ? 'rgba(16,185,129,0.5)' : isFailed ? 'rgba(239,68,68,0.4)' : 'rgba(0,255,163,0.5)';

    return L.divIcon({
        className: '',
        iconSize: [28, 28],
        iconAnchor: [14, 14],
        popupAnchor: [0, -16],
        html: `<div style="
            width:28px;height:28px;border-radius:50%;
            background:${bg};color:${text};
            display:flex;align-items:center;justify-content:center;
            font-weight:700;font-size:12px;font-family:monospace;
            border:2px solid rgba(255,255,255,0.3);
            box-shadow:0 0 12px ${shadow}, 0 2px 6px rgba(0,0,0,0.4);
        ">${index + 1}</div>`,
    });
}

interface RouteMapProps {
    deliveries: Delivery[];
}

// Component to auto-fit bounds
function MapBounds({ deliveries }: { deliveries: Delivery[] }) {
    const map = useMap();

    useEffect(() => {
        if (deliveries.length === 0) return;

        const bounds = L.latLngBounds(deliveries
            .filter(d => d.latitude && d.longitude)
            .map(d => [d.latitude!, d.longitude!]));

        if (bounds.isValid()) {
            map.fitBounds(bounds, { padding: [50, 50] });
        }
    }, [deliveries, map]);

    return null;
}

export const RouteMap: React.FC<RouteMapProps> = ({ deliveries }) => {
    // Default center (Portland, OR area — matches fleet)
    const center: [number, number] = [45.5152, -122.6784];

    // Filter valid locations
    const validDeliveries = deliveries.filter(d => d.latitude && d.longitude);

    // Build polyline path from stops in order
    const routePath: [number, number][] = useMemo(
        () => validDeliveries.map(d => [d.latitude!, d.longitude!]),
        [validDeliveries]
    );

    return (
        <MapContainer
            center={center}
            zoom={11}
            style={{ height: '100%', width: '100%', background: '#161821' }}
            className="z-0"
        >
            <TileLayer
                attribution='&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
                url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
                className="opacity-60 saturate-0 invert"
            />

            <MapBounds deliveries={validDeliveries} />

            {/* Route line connecting stops in order */}
            {routePath.length >= 2 && (
                <>
                    {/* Glow / shadow line */}
                    <Polyline
                        positions={routePath}
                        pathOptions={{
                            color: '#00FFA3',
                            weight: 6,
                            opacity: 0.2,
                            lineCap: 'round',
                            lineJoin: 'round',
                        }}
                    />
                    {/* Main route line */}
                    <Polyline
                        positions={routePath}
                        pathOptions={{
                            color: '#00FFA3',
                            weight: 3,
                            opacity: 0.8,
                            dashArray: '8, 12',
                            lineCap: 'round',
                            lineJoin: 'round',
                        }}
                    />
                </>
            )}

            {/* Stop markers */}
            {validDeliveries.map((delivery, idx) => (
                <Marker
                    key={delivery.id}
                    position={[delivery.latitude!, delivery.longitude!]}
                    icon={makeStopIcon(idx, delivery.status)}
                >
                    <Popup className="text-black">
                        <div className="font-bold flex items-center gap-2">
                            <span style={{
                                background: '#00FFA3', color: '#000',
                                borderRadius: '50%', width: 20, height: 20,
                                display: 'inline-flex', alignItems: 'center', justifyContent: 'center',
                                fontSize: 11, fontWeight: 700
                            }}>
                                {idx + 1}
                            </span>
                            {delivery.customer_name}
                        </div>
                        <div className="text-xs mt-1">{delivery.address}</div>
                        <div className="text-xs text-gray-500 mt-1">
                            Order #{delivery.order_number} &middot; {delivery.status}
                        </div>
                    </Popup>
                </Marker>
            ))}

            {/* Empty state message */}
            {validDeliveries.length === 0 && <EmptyMapOverlay />}
        </MapContainer>
    );
};

function EmptyMapOverlay() {
    return (
        <div style={{
            position: 'absolute', top: 0, left: 0, right: 0, bottom: 0,
            display: 'flex', alignItems: 'center', justifyContent: 'center',
            zIndex: 1000, pointerEvents: 'none',
        }}>
            <div style={{
                background: 'rgba(10,11,16,0.8)', padding: '12px 20px',
                borderRadius: 8, color: '#71717a', fontSize: 13,
                border: '1px solid rgba(255,255,255,0.05)',
            }}>
                Select a route to view stops on the map
            </div>
        </div>
    );
}
