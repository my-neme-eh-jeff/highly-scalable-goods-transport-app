export function connectToBookingWebSocket(
    driverId: number,
    onBookingAssigned: (booking: any) => void
) {
    const ws = new WebSocket(`ws://localhost:8084/ws/driver/assign?driver_id=${driverId}`);

    ws.onopen = () => {
        console.log("Connected to booking assignment WebSocket.");
    };

    ws.onmessage = (event) => {
        const booking = JSON.parse(event.data);
        onBookingAssigned(booking);
    };

    ws.onerror = (error) => {
        console.error("WebSocket error:", error);
    };

    ws.onclose = () => {
        console.log("WebSocket connection closed.");
        // Optionally implement reconnection logic
    };
}

export function connectToLocationWebSocket(
    driverId: number,
    bookingId: number,
    lat: number,
    lng: number
) {
    const ws = new WebSocket(
        `ws://localhost:8083/ws/driver/update-location?driver_id=${driverId}&booking_id=${bookingId}`
    );

    ws.onopen = () => {
        const locationUpdate = { lat, lng };
        ws.send(JSON.stringify(locationUpdate));
    };

    ws.onerror = (error) => {
        console.error("WebSocket error:", error);
    };

    ws.onclose = () => {
        console.log("Location WebSocket connection closed.");
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
        console.error("WebSocket error:", error);
    };

    ws.onclose = () => {
        console.log("Location WebSocket connection closed.");
    };

    return ws;
}
