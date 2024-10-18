
export type Booking = {
  booking_id: number;
  user_id: number;
  pickup_location: {
    lat: number;
    lng: number;
  };
  dropoff_location: {
    lat: number;
    lng: number;
  };
  fare_amount: number;
  status: string;
}
