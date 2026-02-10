import React, { useEffect } from 'react';
import { MapContainer, TileLayer, Marker, Popup, useMap } from 'react-leaflet';
import 'leaflet/dist/leaflet.css';
import type { Delivery } from '../../types/delivery';
import L from 'leaflet';

// Fix for default marker icon in Leaflet with webpack/vite
// We'll use a simple divIcon with Lucide for better styling anyway, but for now standard markers
// or custom SVG markers are best.
import icon from 'leaflet/dist/images/marker-icon.png';
import iconShadow from 'leaflet/dist/images/marker-shadow.png';

const DefaultIcon = L.icon({
    iconUrl: icon,
    shadowUrl: iconShadow,
    iconSize: [25, 41],
    iconAnchor: [12, 41]
});

L.Marker.prototype.options.icon = DefaultIcon;

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
    // Default center (San Francisco)
    const center: [number, number] = [37.7749, -122.4194];

    // Filter valid locations
    const validDeliveries = deliveries.filter(d => d.latitude && d.longitude);

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
                className="opacity-60 saturate-0 invert" // Dark mode style hack
            />

            <MapBounds deliveries={validDeliveries} />

            {validDeliveries.map((delivery, idx) => (
                <Marker
                    key={delivery.id}
                    position={[delivery.latitude!, delivery.longitude!]}
                >
                    <Popup className="text-black">
                        <div className="font-bold flex items-center gap-2">
                            <span className="bg-gable-green text-black rounded-full w-5 h-5 flex items-center justify-center text-xs">
                                {idx + 1}
                            </span>
                            {delivery.customer_name}
                        </div>
                        <div className="text-xs mt-1">{delivery.address}</div>
                        <div className="text-xs text-gray-500 mt-1">Order #{delivery.order_number}</div>
                    </Popup>
                </Marker>
            ))}
        </MapContainer>
    );
};
