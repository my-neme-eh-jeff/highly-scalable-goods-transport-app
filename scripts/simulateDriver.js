import WebSocket from "ws";

const driverId = 2;
const bookingId = 123; 

const ws = new WebSocket(
  `ws://localhost:8083/ws/driver/update-location?driver_id=${driverId}&booking_id=${bookingId}`
);

ws.on("open", function open() {
  console.log("Connected to update driver location service.");

  let lat = 19.076;
  let lng = 72.8777;
  const interval = setInterval(() => {
    lat += 0.0001;
    lng += 0.0001;
    const locationUpdate = { lat, lng };
    ws.send(JSON.stringify(locationUpdate));
  }, 2000);

  setTimeout(() => {
    clearInterval(interval);
    ws.close();

    fetch("http://localhost:8081/api/driver/complete-ride", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ driver_id: driverId, booking_id: bookingId }),
    })
      .then((response) => response.json())
      .then((data) => console.log("Ride completed:", data))
      .catch((error) => console.error("Error completing ride:", error));
  }, 20000); 
});
