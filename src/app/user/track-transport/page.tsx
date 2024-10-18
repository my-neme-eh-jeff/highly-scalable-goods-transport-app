"use client";

import React, { useEffect, useState } from "react";
import dynamic from "next/dynamic";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { CarTaxiFront, ArrowUpDown } from "lucide-react";
import { LatLng } from "leaflet";

const DisplayRouteMap = dynamic(() => import("@/components/DisplayRouteMap"), {
  ssr: false,
});

type SortField = "ID" | "FareAmount";
type SortOrder = "asc" | "desc";

interface SortConfig {
  field: SortField;
  order: SortOrder;
}

type Booking = {
  ID: number;
  UserID: number;
  DriverID: number;
  PickupLocation: {
    lat: number;
    lng: number;
  };
  DropoffLocation: {
    lat: number;
    lng: number;
  };
  FareAmount: number;
  Status: string;
  CreatedAt: string;
};

export default function TrackTransportPage() {
  const [bookings, setBookings] = useState<Booking[]>([]);
  const [filteredBookings, setFilteredBookings] = useState<Booking[]>([]);
  const [userId, setUserId] = useState<number | null>(null);
  const [statusFilter, setStatusFilter] = useState<string>("all");
  const [sortConfig, setSortConfig] = useState<SortConfig>({
    field: "ID",
    order: "desc",
  });

  // Move localStorage access to a useEffect with proper checks
  useEffect(() => {
    if (typeof window !== "undefined") {
      const storedUserId = localStorage.getItem("user_id");
      if (storedUserId) {
        setUserId(parseInt(storedUserId));
      }
    }
  }, []);

  // Only fetch bookings when userId is available
  useEffect(() => {
    if (userId) {
      fetch(`http://localhost:8081/api/user/bookings?user_id=${userId}`)
        .then((res) => res.json())
        .then((data) => {
          setBookings(data);
          setFilteredBookings(data);
        })
        .catch((err) => console.error(err));
    }
  }, [userId]);

  useEffect(() => {
    if (!bookings) return;

    let sorted = [...bookings];

    // Apply status filter
    if (statusFilter !== "all") {
      sorted = sorted.filter(
        (booking) => booking.Status.toLowerCase() === statusFilter.toLowerCase()
      );
    }

    // Apply sorting
    sorted.sort((a, b) => {
      const multiplier = sortConfig.order === "asc" ? 1 : -1;
      return multiplier * (a[sortConfig.field] - b[sortConfig.field]);
    });

    setFilteredBookings(sorted);
  }, [statusFilter, sortConfig, bookings]);

  const toggleSort = (field: SortField) => {
    setSortConfig((current) => ({
      field,
      order:
        current.field === field && current.order === "asc" ? "desc" : "asc",
    }));
  };

  return (
    <div className="container mx-auto p-6 space-y-6">
      <Card>
        <CardHeader>
          <div className="flex items-center space-x-2">
            <CarTaxiFront className="w-6 h-6" />
            <div>
              <CardTitle className="text-xl">Your Ride History</CardTitle>
              <CardDescription>
                View and track all your past and ongoing rides
              </CardDescription>
            </div>
          </div>
        </CardHeader>
        <CardContent>
          <div className="flex justify-end mb-4">
            <Select value={statusFilter} onValueChange={setStatusFilter}>
              <SelectTrigger className="w-[180px]">
                <SelectValue placeholder="Filter by status" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">All Statuses</SelectItem>
                <SelectItem value="REQUESTED">Requested</SelectItem>
                <SelectItem value="STARTED">Started</SelectItem>
                <SelectItem value="CANCELLED">Cancelled</SelectItem>
              </SelectContent>
            </Select>
          </div>

          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>
                  <Button
                    variant="ghost"
                    onClick={() => toggleSort("ID")}
                    className="flex items-center space-x-1"
                  >
                    Booking ID
                    <ArrowUpDown className="w-4 h-4" />
                  </Button>
                </TableHead>
                <TableHead>Status</TableHead>
                <TableHead>
                  <Button
                    variant="ghost"
                    onClick={() => toggleSort("FareAmount")}
                    className="flex items-center space-x-1"
                  >
                    Fare Amount
                    <ArrowUpDown className="w-4 h-4" />
                  </Button>
                </TableHead>
                <TableHead>Action</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {filteredBookings?.map((booking) => (
                <TableRow key={booking.ID}>
                  <TableCell className="font-medium">
                    #{booking.ID} (
                    {new Date(booking.CreatedAt).toLocaleDateString()})
                  </TableCell>
                  <TableCell>
                    <Badge
                      className={`${
                        booking.Status === "COMPLETED"
                          ? "bg-success"
                          : booking.Status == "REJECT"
                          ? "bg-destructive"
                          : booking.Status === "STARTED"
                          ? "bg-green-100"
                          : booking.Status == "REQUESTED"
                          ? "bg-muted-foreground"
                          : booking.Status === "DRIVER_ASSIGNED"
                          ? "bg-yellow-600"
                          : ""
                      }`}
                    >
                      {booking.Status}
                    </Badge>
                  </TableCell>
                  <TableCell>₹ {booking.FareAmount}</TableCell>
                  <TableCell>
                    <Button
                      variant="outline"
                      size="sm"
                      // onClick={() => handleTrack(booking)}
                    >
                      Track Ride
                    </Button>
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </CardContent>
      </Card>

      {bookings && bookings.length > 0 && (
        <Card>
          <CardHeader>
            <CardTitle>Tracking Booking</CardTitle>
          </CardHeader>
          <CardContent>
            <DisplayRouteMap
              bookings={bookings.map((booking) => ({
                booking_id: booking.ID,
                status: booking.Status,
                pickup_location: new LatLng(
                  booking.PickupLocation.lat,
                  booking.PickupLocation.lng
                ),
                dropoff_location: new LatLng(
                  booking.DropoffLocation.lat,
                  booking.DropoffLocation.lng
                ),
              }))}
            />
          </CardContent>
        </Card>
      )}
    </div>
  );
}
