import { LocationManager } from '../components/location/LocationManager';

export const Locations = () => {
    return (
        <div className="p-8 max-w-7xl mx-auto">
            <div className="mb-8">
                <h1 className="text-3xl font-bold text-zinc-100 tracking-tight">Locations</h1>
                <p className="text-zinc-500 mt-1">Manage Warehouse Zones, Aisles, and Bins</p>
            </div>
            <LocationManager />
        </div>
    );
};
