"use client";

import "leaflet/dist/leaflet.css";
import L, { DivIcon, LatLng } from "leaflet";
import React, { useEffect, useRef, useState } from "react";
import {
  MapContainer,
  Marker,
  Polyline,
  Popup,
  TileLayer,
} from "react-leaflet";

import icon from "leaflet/dist/images/marker-icon.png";
import iconShadow from "leaflet/dist/images/marker-shadow.png";

const DefaultIcon = L.icon({
  iconUrl: icon as unknown as string,
  shadowUrl: iconShadow as unknown as string,
});
L.Marker.prototype.options.icon = DefaultIcon;

import { Button } from "./ui/button";

type Booking = {
  booking_id: number;
  status: string;
  pickup_location: LatLng;
  dropoff_location: LatLng;
};

type Props = {
  bookings: Booking[];
};

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

const DisplayRouteMap: React.FC<Props> = ({ bookings }) => {
  const [activeBooking, setActiveBooking] = useState<Booking | null>(null);
  const [positions, setPositions] = useState<LatLng[]>([]);
  const eventSourceRef = useRef<EventSource | null>(null);
  const polylineRef = useRef<L.Polyline | null>(null);
  const [driverPosition, setDriverPosition] = useState<LatLng | null>(null);

  useEffect(() => {
    // Clean up on unmount
    return () => {
      if (eventSourceRef.current) {
        eventSourceRef.current.close();
      }
    };
  }, []);

  const handleTrack = (booking: Booking) => {
    setActiveBooking(booking);
    setPositions([]);

    if (eventSourceRef.current) {
      eventSourceRef.current.close();
    }

    // Connect to SSE endpoint
    eventSourceRef.current = new EventSource(
      `http://localhost:8082/api/user/track-transport?booking_id=${booking.booking_id}`
    );

    eventSourceRef.current.onmessage = (e) => {
      const data = JSON.parse(e.data);
      const newLatLng = new LatLng(data.lat, data.lng);
      setPositions((prevPositions) => [...prevPositions, newLatLng]);
      setDriverPosition(newLatLng);
    };

    eventSourceRef.current.addEventListener("end", () => {
      console.log("Tracking ended for booking:", booking.booking_id);
      if (eventSourceRef.current) {
        eventSourceRef.current.close();
      }
    });
  };

  return (
    <div>
      <div className="flex space-x-4 mb-4">
        {bookings.map((booking) => (
          <Button
            key={booking.booking_id}
            onClick={() => handleTrack(booking)}
            variant={
              activeBooking?.booking_id === booking.booking_id
                ? "default"
                : "outline"
            }
          >
            Track Booking #{booking.booking_id}
          </Button>
        ))}
      </div>

      {activeBooking && (
        <MapContainer
          center={activeBooking.pickup_location}
          zoom={13}
          style={{ height: "500px", width: "100%" }}
        >
          <TileLayer
            attribution="&copy; OpenStreetMap contributors"
            url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
          />
          <Marker position={activeBooking.pickup_location} icon={pickupIcon}>
            <Popup>Pickup Location</Popup>
          </Marker>
          <Marker position={activeBooking.dropoff_location} icon={dropoffIcon}>
            <Popup>Dropoff Location</Popup>
          </Marker>
          {positions.length > 0 && (
            <>
              <Polyline
                positions={positions}
                color="blue"
                ref={(ref) => {
                  polylineRef.current = ref;
                }}
              />
              {driverPosition && (
                <Marker position={driverPosition} icon={driverIcon}>
                  <Popup>Driver&apos;s Current Location</Popup>
                </Marker>
              )}
            </>
          )}
        </MapContainer>
      )}
    </div>
  );
};

export default DisplayRouteMap;
