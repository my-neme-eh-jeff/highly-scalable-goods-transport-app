import { type Booking } from "@/types/booking";

export function connectToBookingWebSocket(
    driverId: number,
    onBookingAssigned: (booking: Booking) => void
) {
    const ws = new WebSocket(`ws://localhost:8084/ws/driver/assign?driver_id=${driverId}`);

    ws.onopen = () => {
        console.log("Connected to booking assignment WebSocket.");
    };

    ws.onmessage = (event) => {
        console.log("event", event)
        const booking = JSON.parse(event.data);
        console.log("Booking assigned:", booking);
        onBookingAssigned(booking);
    };

    ws.onerror = (error) => {
        console.log(error)
    };

    ws.onclose = () => {
        console.log("Booking assignment WebSocket connection closed.");
    };
}

export function connectToPersistentLocationWebSocket(
    driverId: number,
    bookingId: number | null
): WebSocket {
    const wsUrl = bookingId
        ? `ws://localhost:8083/ws/driver/update-location?driver_id=${driverId}&booking_id=${bookingId}`
        : `ws://localhost:8083/ws/driver/update-location?driver_id=${driverId}`;

    const ws = new WebSocket(wsUrl);

    ws.onopen = () => {
        console.log("Connected to location WebSocket.");
    };

    ws.onerror = (error) => {
        console.log("location vale me error h")
        console.log("WebSocket error:", error);
    };

    ws.onclose = () => {
        console.log("Location WebSocket connection closed.");
    };

    return ws;
}
