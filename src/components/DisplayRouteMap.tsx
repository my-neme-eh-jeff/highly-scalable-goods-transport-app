"use client";
import "leaflet/dist/leaflet.css";
import L, { DivIcon, LatLng } from "leaflet";
import React, { useEffect, useRef } from "react";
import {
  MapContainer,
  TileLayer,
  Marker,
  Polyline,
  Popup,
} from "react-leaflet";

// Import marker icons
import icon from "leaflet/dist/images/marker-icon.png";
import iconShadow from "leaflet/dist/images/marker-shadow.png";

// Set default icon
const DefaultIcon = L.icon({
  iconUrl: icon as unknown as string,
  shadowUrl: iconShadow as unknown as string,
});

L.Marker.prototype.options.icon = DefaultIcon;

type Booking = {
  booking_id: number;
  status: string;
  pickup_location: LatLng;
  dropoff_location: LatLng;
};

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

const pickupIcon = createCustomIcon("#4CAF50", "Pickup", "start");
const dropoffIcon = createCustomIcon("#F44336", "Dropoff", "end");

const BookingRoutesMap = ({
  bookings,
  width = "100%",
  height = "500px",
}: {
  bookings: Booking[];
  width?: string;
  height?: string;
}) => {
  const mapRef = useRef<L.Map | null>(null);

  useEffect(() => {
    // Adjust the map view to include all points when bookings change
    if (mapRef.current && bookings.length > 0) {
      const allPoints: LatLng[] = [];
      bookings.forEach((booking) => {
        allPoints.push(booking.pickup_location);
        allPoints.push(booking.dropoff_location);
      });

      const group = new L.FeatureGroup(
        allPoints.map((point) => L.marker([point.lat, point.lng]))
      );
      mapRef.current.fitBounds(group.getBounds(), { padding: [50, 50] });
    }
  }, [bookings]);

  return (
    <div>
      <MapContainer
        center={[19.076, 72.8777]}
        zoom={13}
        scrollWheelZoom={true}
        style={{ width, height }}
        ref={mapRef}
        attributionControl={false}
        minZoom={5}
        className="rounded-lg border-teal-300 border-2"
      >
        <TileLayer url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png" />

        {bookings.map((booking) => {
          const pickup = booking.pickup_location;
          const dropoff = booking.dropoff_location;

          // Determine color based on status
          const color = booking.status === "COMPLETED" ? "green" : "blue";

          // Path coordinates
          const pathCoordinates = [
            new LatLng(pickup.lat, pickup.lng),
            new LatLng(dropoff.lat, dropoff.lng),
          ];

          return (
            <React.Fragment key={booking.booking_id}>
              {/* Draw line between pickup and dropoff */}
              <Polyline positions={pathCoordinates} color={color} />

              {/* Pickup marker */}
              <Marker position={pickup} icon={pickupIcon}>
                <Popup>
                  Pickup Location for Booking ID: {booking.booking_id}
                </Popup>
              </Marker>

              {/* Dropoff marker */}
              <Marker position={dropoff} icon={dropoffIcon}>
                <Popup>
                  Dropoff Location for Booking ID: {booking.booking_id}
                </Popup>
              </Marker>
            </React.Fragment>
          );
        })}
      </MapContainer>
    </div>
  );
};

export default BookingRoutesMap;
