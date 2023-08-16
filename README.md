# pow jabbar WIP
## Anti DDOS via Proof of Work

PoW Jabbar is API library, implementing a security protocol that requires clients to perform a proof-of-work computation, like Adam Back's Hashcash system (the same as used in bitcoin), to gain access to the API. 
By ensuring requests are computationally costly for the client but lightweight for the server, the system effectively counters denial-of-service and spam attacks.

## Features

### Asynchronous Proof-of-Work Challenge
Utilizing the SHA-256 hashing algorithm, the system require the computational proof-of-work to be executed solely on the client side, making it resource-intensive for the requester but cheap for the server.

### Stateless Challenge Design
Challenges are self-contained and do not require any external storage. This design allows high level of scalability across API nodes.


### Timestamp Detection

Each request is timestamped to ensure timely submissions, limiting the window of access and preventing outdated challenges.

### HMAC Signature

To safeguard against replay attacks, each request incorporates an HMAC signature, ensuring that the request is both untampered and unique.

### Difficulty Adjustment

The system dynamically adjusts the computational challenge's difficulty, ensuring that it remains effective regardless of evolving hardware or network conditions.

### Client-Side Implementation

With SHA-256 as the foundation, the system is designed for easy client-side implementation using JavaScript.

### Inspiration
Heavely insipred by Islam Bekbuzarov and his https://github.com/blkmlk/ddos-pow

### Links
https://en.wikipedia.org/wiki/Hashcash


