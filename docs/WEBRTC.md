# WebRTC

Selective Forwarding Unit (SFU)

## Articles to read:
- https://gabrieltanner.org/blog/broadcasting-ion-sfu

## How does WebRTC work?
Works via UDP-protocol

It stands for Web Real-Time Communication and allows direct communication between browsers. It does not use a websocket with a central server due to the delay.
With WebRTC the clients can connect to each other and bypass the server.

**Example:**

Given Client X wants to connect with Client Y

1. X asks the **STUN** 'Who am I?'

STUN stands for Session Traveral Utilities for NATS. It can be used by an endpoint to determine the IP address and port allocated to it by a NAT.

2. The STUN returns a **Symmetric NAT** response to X

A symmetric NAT is one where all requests from the same internal IP address and port, to a specific destination IP address and port, are mapped to the same external IP address and port.

3. X then asks the **TURN** for a channel

TURN (Traversal Using Relays around NAT) is a protocol that assists in the traversal of network address translators (NAT) or firewalls for webRTC applications. TURN Server allows clients to send and receive data through an intermediary server. The TURN protocol is the extension to STUN

4. X then sends the **SDP** to the Backend

The Session Description Protocol (SDP) is a format for describing multimedia communication sessions for the purposes of session announcement and session invitation.

5. The backend then send the **Offer SDP** from X to Y
6. Y then sends a **Answer SDP** to the Backend
7. The backend then sends the **Answer SDP** from Y to X
8. X sends its **ICE candidate** to the Backend
9. The backend then sends the **ICE candidate (X)** from X to Y 
10. Y sends its **ICE candidate** to the Backend
11. The backend then sends the **ICE candidate (Y)** from Y to X


## How does WebRTC


