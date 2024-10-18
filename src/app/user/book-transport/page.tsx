"use client";
import dynamic from "next/dynamic";
const MapComponent = dynamic(() => import("@/components/EditableMap"), {
  ssr: false,
});
import { Button } from "@/components/ui/button";
import { LatLng } from "leaflet";
import { Loader2, MousePointerClick } from "lucide-react";
import { useEffect, useState } from "react";
import { toast } from "sonner";
import { titleVariants } from "@/components/textVariants";
import { cn } from "@/lib/utils";

type FareDetailsResponse = {
  distance_km: number;
  fare_amount: number;
  status: string;
};
type BookingResponse = {
  booking_id: number;
  status: string;
};

export default function UserBookingPage() {
  const [pickupLocation, setPickupLocation] = useState<LatLng | null>(null);
  const [dropoffLocation, setDropoffLocation] = useState<LatLng | null>(null);
  const [fareDetails, setFareDetails] = useState<FareDetailsResponse | null>(
    null
  );
  const [isGetFareLoading, setIsGetFareLoading] = useState(false);
  const [isBookTransportLoading, setIsBookTransportLoading] = useState(false);
  const [user_id, setUserId] = useState(0);
  const bookTransport = async () => {
    if (!fareDetails || user_id == 0) {
      toast.warning("Please get fare details first.");
      return;
    }
    setIsBookTransportLoading(true);
    const id = toast.loading("Booking transport...");
    try {
      const response = await fetch(
        "http://localhost:8081/api/user/book-transport",
        {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({
            user_id: user_id,
            pickup_location: {
              lat: pickupLocation!.lat,
              lng: pickupLocation!.lng,
            },
            dropoff_location: {
              lat: dropoffLocation!.lat,
              lng: dropoffLocation!.lng,
            },
            fare_amt: fareDetails.fare_amount,
          }),
        }
      );
      if (!response.ok) {
        throw new Error("Failed to book transport.");
      }
      const data: BookingResponse = await response.json();
      toast.success(
        "Transport booked successfully. Your booking id is" + data.booking_id
      );
    } catch (error) {
      console.error(error);
      toast.error("Failed to book transport.");
    } finally {
      setIsBookTransportLoading(false);
      toast.dismiss(id);
    }
  };
  const handleGetFare = async () => {
    if (!pickupLocation || !dropoffLocation || user_id == 0) {
      toast.warning("Please select both pickup and dropoff locations.");
      return;
    }
    setIsGetFareLoading(true);
    const id = toast.loading("Fetching fare details...");
    try {
      const response = await fetch("http://localhost:8080/api/user/get-fare", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          user_id: user_id,
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
      toast.success("Fare details fetched successfully.");
    } catch (error) {
      console.error(error);
      toast.error("Failed to fetch fare details.");
    } finally {
      setIsGetFareLoading(false);
      toast.dismiss(id);
    }
  };
  useEffect(() => {
    const generateUserId = () => {
      const existingUserId = window.localStorage.getItem("user_id");

      if (existingUserId) {
        return parseInt(existingUserId);
      }

      const newUserId =
        Math.floor(Math.random() * (Number.MAX_SAFE_INTEGER / 100000)) + 1;
      window.localStorage.setItem("user_id", newUserId.toString());
      return newUserId;
    };
    setUserId(generateUserId());
  }, []);

  return (
    <div className="flex justify-center place-items-center place-content-center flex-col gap-y-10 mt-16">
      <div className="flex flex-col justify-center place-items-center place-content-center">
        {fareDetails ? (
          <>
            <h1 className={cn(titleVariants({ color: "golden", size: "sm" }))}>
              Fare Amount: ₹{Math.floor(fareDetails.fare_amount)}
            </h1>
            <p className="text-sm text-muted-foreground flex mt-1 align-middle place-content-center place-items-center">
              Distance: {Math.floor(fareDetails.distance_km)}km{" "}
              {Math.floor(fareDetails.distance_km * 1000)}m
            </p>
            <Button
              variant={"success"}
              className="mt-4"
              isLoading={isBookTransportLoading}
              onClick={bookTransport}
            >
              {isBookTransportLoading ? (
                <Loader2 className="mr-2 h-4 w-4 animate-spin" />
              ) : null}
              Proceed to Book
            </Button>
          </>
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
        allowReset={true}
        allowTheUseOfUserLocationAsInitial={true}
      />
      <Button
        variant={"success"}
        className="w-40"
        onClick={handleGetFare}
        type="submit"
        isLoading={isGetFareLoading}
      >
        {isGetFareLoading ? (
          <Loader2 className="mr-2 h-4 w-4 animate-spin" />
        ) : null}
        Get Fare
      </Button>
    </div>
  );
}
