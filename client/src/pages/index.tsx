import React, { useEffect } from "react";
import { Navigator } from "../components/Nav";
const Index = () => {
  useEffect(() => {
    const ws = new WebSocket("ws://localhost:4000/ws/client?user_id=maotrix");

    ws.onopen = (event) => {
      ws.send(
        JSON.stringify({
          data: "this is test",
          from: "maotrix",
          to: "mcuve",
        })
      );
    };

    ws.onmessage = function (event) {
      console.log(JSON.parse(event.data));
    };
  }, []);

  return (
    <div>
      <Navigator />
    </div>
  );
};

export default Index;
