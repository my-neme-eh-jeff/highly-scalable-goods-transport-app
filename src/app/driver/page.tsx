"use client";

import React, { useEffect, useState } from "react";
import dynamic from "next/dynamic";
import { toast } from "sonner";
import { connectToBookingWebSocket } from "@/utils/driverWebSocket";
import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogFooter,
  DialogTitle,
  DialogDescription,
} from "@/components/ui/dialog";
import { type Booking } from "@/types/booking";
import { titleVariants } from "@/components/textVariants";
import { cn } from "@/lib/utils";

const DriverMap = dynamic(() => import("@/components/DriverMap"), {
  ssr: false,
});

export default function DriverDashboardPage() {
  const [driverId, setDriverId] = useState<number | null>(null);
  const [assignedBooking, setAssignedBooking] = useState<Booking | null>(null);
  const [showDialog, setShowDialog] = useState(false);

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
        setShowDialog(true); // Show the dialog when a new booking is assigned
        toast.success(`New booking assigned: ${booking.booking_id}`);
      });
    }
  }, [driverId]);

  const handleAcceptBooking = async () => {
    if (!assignedBooking || !driverId) return;

    try {
      const response = await fetch(
        "http://localhost:8081/api/driver/respond-booking",
        {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({
            driver_id: driverId,
            booking_id: assignedBooking.booking_id,
            response: "ACCEPT",
          }),
        }
      );
      if (!response.ok) {
        throw new Error("Failed to accept booking.");
      }
      toast.success("Booking accepted.");
      setShowDialog(false);
    } catch (error) {
      console.error(error);
      toast.error("Failed to accept booking.");
    }
  };

  const handleRejectBooking = async () => {
    if (!assignedBooking || !driverId) return;

    try {
      const response = await fetch(
        "http://localhost:8081/api/driver/respond-booking",
        {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({
            driver_id: driverId,
            booking_id: assignedBooking.booking_id,
            response: "REJECT",
          }),
        }
      );
      if (!response.ok) {
        throw new Error("Failed to reject booking.");
      }
      toast.success("Booking rejected.");
      setAssignedBooking(null);
      setShowDialog(false);
    } catch (error) {
      console.error(error);
      toast.error("Failed to reject booking.");
    }
  };

  const DialogComponent = () => (
    <Dialog open={showDialog} onOpenChange={setShowDialog}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>New Booking Assigned</DialogTitle>
          <DialogDescription>
            A new booking has been assigned to you. Do you want to accept it?
          </DialogDescription>
        </DialogHeader>
        {assignedBooking && (
          <div className="space-y-2">
            <p>
              <strong>Booking ID:</strong> {assignedBooking.booking_id}
            </p>
            <p>
              <strong>Pickup Location:</strong>{" "}
              {`${assignedBooking.pickup_location.lat}, ${assignedBooking.pickup_location.lng}`}
            </p>
            <p>
              <strong>Drop-off Location:</strong>{" "}
              {`${assignedBooking.dropoff_location.lat}, ${assignedBooking.dropoff_location.lng}`}
            </p>
            <p>
              <strong>Fare Amount:</strong> ₹{assignedBooking.fare_amount}
            </p>
          </div>
        )}
        <DialogFooter>
          <Button variant="secondary" onClick={handleRejectBooking}>
            Reject
          </Button>
          <Button onClick={handleAcceptBooking}>Accept</Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );

  return (
    <div className="container mx-auto mt-8">
      <h1
        className={cn(
          titleVariants({ color: "blue", size: "sm" }),
          "text-center justify-center place-items-center place-content-center flex mb-4"
        )}
      >
        Driver Dashboard
      </h1>
      <DriverMap driverId={driverId} assignedBooking={assignedBooking} />
      <DialogComponent />
    </div>
  );
}
