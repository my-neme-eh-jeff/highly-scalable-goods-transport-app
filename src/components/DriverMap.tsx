// components/DriverMap.tsx

"use client";
import "leaflet/dist/leaflet.css";
import L, { DivIcon, LatLng } from "leaflet";
import React, { useEffect, useRef, useState } from "react";
import { MapContainer, TileLayer, Marker, Popup } from "react-leaflet";

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
import { type Booking } from "@/types/booking";

interface DriverMapProps {
  driverId: number | null;
  assignedBooking: Booking | null;
}

const createCustomIcon = (
  color: string,
  label: string,
  customClassName: string
): DivIcon => {
  return L.divIcon({
    className: customClassName,
    html: `<div style="background-color: ${color}; width: 2.5rem; height: 2.5rem; border-radius: 50%; display: flex; justify-content: center; align-items: center; color: white; font-weight: bold;">${label}</div>`,
  });
};

const driverIcon = createCustomIcon("#FFA500", "YOU", "driver-icon");
const pickupIcon = createCustomIcon("red", "Pick", "pickup-icon");
const dropoffIcon = createCustomIcon("#28A745", "Drop", "dropoff-icon");

const DriverMap: React.FC<DriverMapProps> = ({ driverId, assignedBooking }) => {
  const [driverLocation, setDriverLocation] = useState<LatLng | null>(null);
  const [pickupLocation, setPickupLocation] = useState<LatLng | null>(null);
  const [dropoffLocation, setDropoffLocation] = useState<LatLng | null>(null);
  const [mapCenter, setMapCenter] = useState<LatLng>(
    new LatLng(19.076, 72.8777)
  ); // Default to Mumbai coordinates
  const mapRef = useRef<L.Map | null>(null);
  const wsRef = useRef<WebSocket | null>(null);

  useEffect(() => {
    if (assignedBooking) {
      setPickupLocation(
        new LatLng(
          assignedBooking.pickup_location.lat,
          assignedBooking.pickup_location.lng
        )
      );
      setDropoffLocation(
        new LatLng(
          assignedBooking.dropoff_location.lat,
          assignedBooking.dropoff_location.lng
        )
      );
    } else {
      setPickupLocation(null);
      setDropoffLocation(null);
    }
  }, [assignedBooking]);

  useEffect(() => {
    if (!driverId) return;

    let watchId: number;

    if (navigator.geolocation) {
      // Establish WebSocket connection
      wsRef.current = connectToPersistentLocationWebSocket(
        driverId,
        assignedBooking?.booking_id || null
      );

      watchId = navigator.geolocation.watchPosition(
        (position) => {
          const { latitude, longitude } = position.coords;
          const newLocation = new LatLng(latitude, longitude);
          setDriverLocation(newLocation);

          // Update map center to driver's current location
          setMapCenter(newLocation);

          // Send location update via WebSocket
          if (wsRef.current && wsRef.current.readyState === WebSocket.OPEN) {
            console.log("Sending location update:", {
              lat: latitude,
              lng: longitude,
            });
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
  }, [driverId, assignedBooking]);

  useEffect(() => {
    // Fit map bounds to include all markers
    if (mapRef.current && driverLocation) {
      const bounds = L.latLngBounds([driverLocation]);

      if (pickupLocation) bounds.extend(pickupLocation);
      if (dropoffLocation) bounds.extend(dropoffLocation);

      mapRef.current.fitBounds(bounds, { padding: [50, 50] });
    }
  }, [driverLocation, pickupLocation, dropoffLocation]);

  return (
    <div>
      <MapContainer
        center={mapCenter}
        zoom={10}
        scrollWheelZoom={true}
        style={{
          width: "800px",
          height: "500px",
          margin: "auto auto",
          zIndex: 10,
        }}
        ref={mapRef}
        attributionControl={false}
        minZoom={5}
        className="rounded-lg border-teal-300 border-2"
      >
        <TileLayer url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png" />

        {driverLocation && (
          <Marker position={driverLocation} icon={driverIcon}>
            <Popup>Your Location</Popup>
          </Marker>
        )}
        {pickupLocation && (
          <Marker position={pickupLocation} icon={pickupIcon}>
            <Popup>Pickup Location</Popup>
          </Marker>
        )}
        {dropoffLocation && (
          <Marker position={dropoffLocation} icon={dropoffIcon}>
            <Popup>Drop-off Location</Popup>
          </Marker>
        )}
      </MapContainer>
    </div>
  );
};

export default DriverMap;
