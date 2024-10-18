"use client";
import "leaflet/dist/leaflet.css";
import L, {
  DivIcon,
  LatLng,
  Map as LeafletMap,
  type LatLngLiteral,
} from "leaflet";
import React, { useEffect, useRef, useState } from "react";
import {
  MapContainer,
  Marker,
  Popup,
  TileLayer,
  useMapEvents,
} from "react-leaflet";

//!import order matters
import icon from "leaflet/dist/images/marker-icon.png";

import iconShadow from "leaflet/dist/images/marker-shadow.png";

import "leaflet/dist/leaflet.css";

const DefaultIcon = L.icon({
  iconUrl: icon as unknown as string,
  shadowUrl: iconShadow as unknown as string,
});

L.Marker.prototype.options.icon = DefaultIcon;

import { toast } from "sonner";
import { Button } from "./ui/button";

type MapComponentProps = {
  pickupLocation: LatLng | null;
  setPickupLocation: React.Dispatch<React.SetStateAction<LatLng | null>>;
  dropoffLocation: LatLng | null;
  setDropoffLocation: React.Dispatch<React.SetStateAction<LatLng | null>>;
  width: string;
  height: string;
  allowReset: boolean;
  onLocationSelect?: (latlng: LatLng) => void;
  allowTheUseOfUserLocationAsInitial: boolean;
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

const startIcon = createCustomIcon("#4CAF50", "Pickup", "start");
const endIcon = createCustomIcon("#F44336", "Dropoff", "end");
const initialIcon = createCustomIcon("#2196F3", "You", "initial");

const MapComponent: React.FC<MapComponentProps> = ({
  pickupLocation,
  setPickupLocation,
  dropoffLocation,
  setDropoffLocation,
  width,
  height,
  allowReset,
  onLocationSelect,
  allowTheUseOfUserLocationAsInitial,
}) => {
  const [userPosition, setUserPosition] = useState<LatLngLiteral>({
    lat: 19.076,
    lng: 72.8777,
  });
  const mapRef = useRef<LeafletMap | null>(null);
  const effectRan = useRef(false);
  const [isDefaultLocation, setIsDefaultLocation] = useState(true);
  const [toggleMarker, setToggleMarker] = useState<"pickup" | "drop-off">(
    "pickup"
  );
  const [zoom, setZoom] = useState(13);
  useEffect(() => {
    if (effectRan.current) return;

    if (navigator.geolocation) {
      navigator.geolocation.getCurrentPosition(
        (position) => {
          const { latitude, longitude } = position.coords;
          setUserPosition({ lat: latitude, lng: longitude });
          setIsDefaultLocation(false);
          setZoom(15);
        },
        (err) => {
          if (err.code === 1) {
            toast.error("Please accept the browser request for Geolocation.");
          } else if (err.code === 2) {
            toast.error(
              "Location information is unavailable. Please turn on GPS"
            );
            toast.warning("Set default location to Mumbai, India");
          } else if (err.code === 3) {
            toast.error("The request to get user location timed out.");
          }
        }
      );
    } else {
      toast.error("Geolocation is not supported by this browser.");
    }

    effectRan.current = true;
  }, []);

  const useUserLocationAsInitial = () => {
    setPickupLocation(new LatLng(userPosition.lat, userPosition.lng));
  };

  const LocationMarker: React.FC = () => {
    useMapEvents({
      click(e) {
        if (toggleMarker == "pickup") {
          setPickupLocation(e.latlng);
          setToggleMarker("drop-off");
        } else if (toggleMarker == "drop-off") {
          setDropoffLocation(e.latlng);
          setToggleMarker("pickup");
        }
        if (onLocationSelect) {
          onLocationSelect(e.latlng);
        }
      },
    });
    return null;
  };

  const handleReset = () => {
    setPickupLocation(null);
    setDropoffLocation(null);
    setToggleMarker("pickup");
    if (mapRef.current) {
      mapRef.current.setView(userPosition!, zoom);
    }
  };

  return (
    <div>
      <MapContainer
        center={userPosition}
        zoom={zoom}
        scrollWheelZoom={true}
        style={{ width, height }}
        ref={mapRef}
        attributionControl={false}
        minZoom={5}
        className="rounded-lg border-teal-300 border-2"
      >
        <TileLayer url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png" />
        <LocationMarker />
        {!isDefaultLocation && (
          <Marker position={userPosition} icon={initialIcon}>
            <Popup>Your Location</Popup>
          </Marker>
        )}
        {pickupLocation && (
          <Marker draggable position={pickupLocation} icon={startIcon}>
            <Popup>Pickup Location</Popup>
          </Marker>
        )}
        {dropoffLocation && (
          <Marker draggable position={dropoffLocation} icon={endIcon}>
            <Popup>Drop-off Location</Popup>
          </Marker>
        )}
      </MapContainer>
      <div className="flex gap-x-10 justify-between mt-2">
        {allowReset && (
          <Button variant={"destructive"} onClick={handleReset}>
            Reset Locations
          </Button>
        )}
        {!isDefaultLocation && allowTheUseOfUserLocationAsInitial && (
          <Button onClick={useUserLocationAsInitial}>
            Use current Location
          </Button>
        )}
      </div>
    </div>
  );
};

export default MapComponent;
