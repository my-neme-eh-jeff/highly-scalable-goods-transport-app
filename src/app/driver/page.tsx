"use client";
import React, { useEffect, useState } from "react";
import { Button } from "@/components/ui/button";
import dynamic from "next/dynamic";
import { type Booking } from "@/types/booking";
import { toast } from "sonner";
import { connectToBookingWebSocket } from "@/utils/driverWebSocket";

const DriverMap = dynamic(
  () => import("@/components/DriverMap"),
  {
    ssr: false,
  }
);

export default function DriverDashboardPage() {
  const [driverId, setDriverId] = useState<number | null>(null);
  const [assignedBooking, setAssignedBooking] = useState<Booking | null>(null);
  const [location, setLocation] = useState<{ lat: number; lng: number } | null>(
    null
  );

  useEffect(() => {
    // Simulate driver login and get driver ID
    const storedDriverId = window.localStorage.getItem("driver_id");
    if (storedDriverId) {
      setDriverId(parseInt(storedDriverId));
    } else {
      const newDriverId = Math.floor(Math.random() * 1000000) + 1;
      window.localStorage.setItem("driver_id", newDriverId.toString());
      setDriverId(newDriverId);
    }
  }, []);

  useEffect(() => {
    if (driverId) {
      // Connect to WebSocket for booking assignments
      connectToBookingWebSocket(driverId, (booking: Booking) => {
        setAssignedBooking(booking);
        toast.success(`New booking assigned: ${booking.booking_id}`);
      });
    }
  }, [driverId]);

  const handleAcceptBooking = async () => {
    if (!assignedBooking || !driverId) return;

    try {
      const response = await fetch(
        "http://localhost:8081/api/driver/accept-booking",
        {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({
            driver_id: driverId,
            booking_id: assignedBooking.booking_id,
          }),
        }
      );
      if (!response.ok) {
        throw new Error("Failed to accept booking.");
      }
      toast.success("Booking accepted.");
    } catch (error) {
      console.error(error);
      toast.error("Failed to accept booking.");
    }
  };

  const handleCompleteRide = async () => {
    if (!assignedBooking || !driverId) return;

    try {
      const response = await fetch(
        "http://localhost:8081/api/driver/complete-ride",
        {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({
            driver_id: driverId,
            booking_id: assignedBooking.booking_id,
          }),
        }
      );
      if (!response.ok) {
        throw new Error("Failed to complete ride.");
      }
      toast.success("Ride completed.");
      setAssignedBooking(null);
    } catch (error) {
      console.error(error);
      toast.error("Failed to complete ride.");
    }
  };

  return (
    <div>
      <h1>Driver Dashboard</h1>
      {assignedBooking ? (
        <div>
          <h2>Assigned Booking ID: {assignedBooking.booking_id}</h2>
          <p>
            Pickup Location:{" "}
            {`${assignedBooking.pickup_location.lat}, ${assignedBooking.pickup_location.lng}`}
          </p>
          <p>
            Drop-off Location:{" "}
            {`${assignedBooking.dropoff_location.lat}, ${assignedBooking.dropoff_location.lng}`}
          </p>
          <Button onClick={handleAcceptBooking}>Accept Booking</Button>
          <Button onClick={handleCompleteRide}>Complete Ride</Button>
        </div>
      ) : (
        <p>No bookings assigned yet.</p>
      )}
      <DriverMap 
        driverId={driverId}
        bookingId={assignedBooking?.booking_id || null}
        location={location}
        setLocation={setLocation}
      />
    </div>
  );
}
