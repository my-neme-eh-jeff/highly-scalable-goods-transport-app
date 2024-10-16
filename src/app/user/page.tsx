"use client";
import MapComponent from "@/components/Map";
import { Button } from "@/components/ui/button";
import { LatLng } from "leaflet";
import { MousePointerClick } from "lucide-react";
import { useState } from "react";
import { toast } from "sonner";

type FareDetailsResponse = {
  distance_km: number;
  fareAmount: number;
  status: string;
};

export default function UserPage() {
  const [pickupLocation, setPickupLocation] = useState<LatLng | null>(null);
  const [dropoffLocation, setDropoffLocation] = useState<LatLng | null>(null);
  const [fareDetails, setFareDetails] = useState<FareDetailsResponse | null>(
    null
  );

  const bookTransport = async () => {
    
  };
  const handleGetFare = async () => {
    if (!pickupLocation || !dropoffLocation) {
      toast.warning("Please select both pickup and dropoff locations.");
      return;
    }

    const id = toast.loading("Fetching fare details...");
    try {
      const response = await fetch("http://localhost:8080/api/user/get-fare", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          user_id: 1,
          pickup_location: {
            lat: pickupLocation.lat,
            lng: pickupLocation.lng,
          },
          dropoff_location: {
            lat: dropoffLocation.lat,
            lng: dropoffLocation.lng,
          },
        }),
      });

      if (!response.ok) {
        throw new Error("Failed to fetch fare details.");
      }
      const data: FareDetailsResponse = await response.json();
      setFareDetails(data);
      toast.success("Fare details fetched successfully.", { id });
    } catch (error) {
      console.error(error);
      toast.error("Failed to fetch fare details.", { id });
    } finally {
      toast.dismiss(id);
    }
  };

  return (
    <div className="flex justify-center place-items-center place-content-center flex-col gap-y-10 mt-16">
      <div className="flex flex-col justify-center place-items-center place-content-center">
        {fareDetails ? (
          <div>
            <h1 className="scroll-m-20 border-b pb-2 text-3xl font-semibold tracking-tight first:mt-0">
              Fare Amount: {fareDetails.fareAmount}
            </h1>
            <p className="text-sm text-muted-foreground flex mt-1 align-middle place-content-center place-items-center">
              Distance: {fareDetails.distance_km} km
            </p>
            <Button onClick={bookTransport}>Proceed to Book</Button>
          </div>
        ) : (
          <>
            <h1 className="scroll-m-20 border-b pb-2 text-3xl font-semibold tracking-tight first:mt-0">
              Get fare
            </h1>
            <p className="text-sm text-muted-foreground flex mt-1 align-middle place-content-center place-items-center">
              Click <MousePointerClick className="mx-2" /> on the map to drop
              pins
            </p>
          </>
        )}
      </div>
      <MapComponent
        pickupLocation={pickupLocation}
        setPickupLocation={setPickupLocation}
        dropoffLocation={dropoffLocation}
        setDropoffLocation={setDropoffLocation}
        width="750px"
        height="500px"
        allowTheUseOfUserLocationAsInitial={true}
        allowReset={true}
      />
      <Button
        variant={"default"}
        className="text-black bg-green-400 active:bg-green-600 hover:bg-green-600 transition-all w-40 "
        onClick={handleGetFare}
        type="submit"
      >
        Get Fare
      </Button>
    </div>
  );
}
