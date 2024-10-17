"use client";
import "leaflet/dist/leaflet.css";
import L, { DivIcon, LatLng } from "leaflet";
import React, { useEffect, useRef, useState } from "react";
import { MapContainer, TileLayer, Marker } from "react-leaflet";

// Import marker icons
import icon from "leaflet/dist/images/marker-icon.png";
import iconShadow from "leaflet/dist/images/marker-shadow.png";

// Set default icon
const DefaultIcon = L.icon({
  iconUrl: icon as unknown as string,
  shadowUrl: iconShadow as unknown as string,
});

L.Marker.prototype.options.icon = DefaultIcon;

import { toast } from "sonner";
import { connectToPersistentLocationWebSocket } from "@/utils/driverWebSocket";

interface DriverMapProps {
  driverId: number | null;
  bookingId: number | null;
  location: { lat: number; lng: number } | null;
  setLocation: React.Dispatch<
    React.SetStateAction<{ lat: number; lng: number } | null>
  >;
}

const createCustomIcon = (
  color: string,
  label: string,
  customClassName: string
): DivIcon => {
  return L.divIcon({
    className: customClassName,
    html: `<div style="background-color: ${color}; width: 3rem; height: 3rem; border-radius: 50%; display: flex; justify-content: center; align-items: center; color: white; font-weight: bold;">${label}</div>`,
  });
};

const driverIcon = createCustomIcon("#FFA500", "Driver", "driver");

const DriverMap: React.FC<DriverMapProps> = ({
  driverId,
  bookingId,
  location,
  setLocation,
}) => {
  const [mapCenter, setMapCenter] = useState<LatLng>(
    new LatLng(19.076, 72.8777)
  );
  const mapRef = useRef<L.Map | null>(null);
  const wsRef = useRef<WebSocket | null>(null);

  useEffect(() => {
    if (!driverId) return;

    let watchId: number;

    if (navigator.geolocation) {
      // Establish WebSocket connection
      wsRef.current = connectToPersistentLocationWebSocket(driverId, bookingId);

      watchId = navigator.geolocation.watchPosition(
        (position) => {
          const { latitude, longitude } = position.coords;
          const newLocation = { lat: latitude, lng: longitude };
          setLocation(newLocation);

          // Update map center
          setMapCenter(new LatLng(newLocation.lat, newLocation.lng));

          // Send location update via WebSocket
          if (wsRef.current && wsRef.current.readyState === WebSocket.OPEN) {
            const locationUpdate = { lat: latitude, lng: longitude };
            wsRef.current.send(JSON.stringify(locationUpdate));
          }
        },
        (error) => {
          console.error("Error getting location:", error);
          toast.error("Error getting location.");
        },
        { enableHighAccuracy: true, maximumAge: 0, timeout: 5000 }
      );
    } else {
      toast.error("Geolocation is not supported by this browser.");
    }

    // Cleanup function
    return () => {
      if (watchId !== undefined) {
        navigator.geolocation.clearWatch(watchId);
      }
      if (wsRef.current) {
        wsRef.current.close();
      }
    };
  }, [driverId, bookingId, setLocation]);

  return (
    <div>
      <MapContainer
        center={mapCenter}
        zoom={13}
        scrollWheelZoom={true}
        style={{ width: "100%", height: "500px" }}
        ref={mapRef}
        attributionControl={false}
        minZoom={5}
        className="rounded-lg border-teal-300 border-2"
      >
        <TileLayer url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png" />

        {location && (
          <Marker position={location} icon={driverIcon}>
            {/* Optionally add a Popup */}
          </Marker>
        )}
      </MapContainer>
    </div>
  );
};

export default DriverMap;
