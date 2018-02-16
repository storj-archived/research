# Storj Network Specification Draft

This a draft and working document for an informal specification of an improved Storj network. It describes the behaviors of each network component, and the interfaces that each expose and use to interact. This is building on ideas defined in both versions of the Storj whitepaper (v1 in 2014 and v2 in 2016), and the latest Storj network v1.2 implementation – with the goal at being decentralized at the foundation.

There are two main groups of responsibilities; Storj Node and Storj Client. A Storj Node is responsible for storing data geographically distributed with other Storj Nodes. This includes a database layer for file pointer and meta data, and a storage layer for larger file shards. A network of Storj Nodes is responsible for maintaining, replicating and repairing data when it’s lost. Storj Nodes pay each other in cryptocurrency STORJ tokens for storage and bandwidth. The Client is responsible for encrypting, erasure encoding and transferring files to Storj Nodes.

## Storj Node

A Storj Node exposes an interface for Clients to communicate to store and retrieve meta data that describes the locations of the shards that make up files, and for downloading shard data in the network. The client layer API is exposed over a TCP/TLS socket with HTTP with JSON requests and responses most often over TLS. The storage layer interface is a JSON-RPC API over a TCP socket over HTTP.  There are RPC commands for storing and retrieving arbitrary data with contracts from a Node, and endpoints for transferring data to Clients in the network over HTTP. Implements SIP5, SIP6 and SIP9

### Behaviors & Services

### Service Dependencies

### Payout Cron Worker Behaviors

### Meta Database Behaviors

### Shard Replication Worker Behaviors

### Shard Repair Worker Behaviors

### Identity

### Authentication Methods

### Client API Error Response Codes

### Client Endpoints

### Storage API

- `POST /` with JSON-RPC commands:
  - `PING` → Request to send back a PONG
    - Authorization
      - None
    - Params
      - contact → The requesters contact object
    - Response
      - contact → The responders contact object
  - `MIRROR` → Request to download shard from another Node. The request is finished   nearly immediately, however the actual shard transfer is then run in the background. Once the transfer is finished an exchange report is sent back to the Bridge for the contract, to confirm that the transfer is finished.
    - Authorization
      - Must be part of “trusted” bridges configuration
      - Must already have a contract for the data from a previous ALLOC command
      - Must have a token from a previous RETRIEVE command
    - Params
      - contact → The requesters contact object
      - node → A farmer contact object to download the data
      - port → The port of the Node
      - address → The IP address of the Node
      - nodeID → The nodeID of the Node
      - token → The token to use to download the shard
      - signature → ECDSA signature data for the request
    - Response
      - (empty object)
  - `RETRIEVE` → Request to download a shard. The shard’s contract describes which identity will pay for the shard, and therefore which requests to authorize to receive a token to download.
    - Authorization
      - Must be authenticated with identity on the shard contract
    - Params
      - contact → The requesters contact object
      - data_hash → The hash of the requested data
      - nonce → Used for signature verification
      - signature → ECDSA signature data for the request
    Response
      - token → The token for downloading the shard
  - `ALLOC` → Request to upload a shard and responds with a token
    - Authorization
      - An algorithm for determining if the request is within acceptable thresholds for expected payout, and if the asking Node has a positive reputation for making payments. In the case that the request is reject a list of reasons should be given as to reason for the rejection.
    - Params
      - contact → The requesters contact object
      - dataHash → The hash of the data
      - dataSize → The size of the data in bytes
      - storeBegin → When the data is initially to be stored
      - storeEnd → Desired length of storage
      - nonce → Used for signature verification
      - signature → ECDSA signature data for the request
    - Response
      - contract → A farmer signed contract object
      - token → A token to upload data for the shard
  - `RENEW` → Request to renew length of shard storage
    - Authorization
      - Must be authenticated with identity on the shard contract
    - Params
      - dataHash → The data that should be renewed
      - storeEnd → The new end date
      - nonce →Used for signature verification
      - signature → ECDSA signature data for the request
    - Response
      - dataHash → The data hash that was renewed
      - storeEnd → The new store end date

### Storage Transfer API

- `GET /shards/<shard-hash>?token=<token>` → Download a shard
  - Authorization
    - The token must be found and associated with the shard hash
  - Responses
    - 200 → Success
  - Behaviors
    - Once the shard has been successfully transferred an exchange report is sent to the Bridge for the shard to confirm the transfer.
- `POST /shards/<shard-hash>?token=<token>` → Upload a shard
  - Authorization
    - The token must be found and associated with the shard hash
  - Responses
    - 200 → Success
  - Behaviors
    - Once the shard has been successfully transferred an exchange report is sent to the Bridge for the shard to confirm the transfer.

## Storj Client

A Client will transfer files to and from the Storj network. It will encrypt and erasure encode the file and coordinate with a Storj Node to store file pointer and meta data to later retrieve the file. The Client transfers shards directly to other Storj Nodes over HTTP concurrently. Implements SIP5 and SIP9.

### Identity

- The identity of a Storj Client is based on a 12-24 word phrase (based on BIP39), keys can be derived from this key for various purposes from signing requests to deriving private keys for file encryption/decryption.

### Behaviors

- Lightweight client for concurrent transfer of files on the network with encryption and erasure encoding.

### Service Dependencies

- Storj Node → A Client needs to communicate with a Storj Node to be able to retrieve file pointer and meta data to download data from the network. A network of Storj Nodes are necessary to be able to store data in the network, and the Client’s Node needs to make payments for services.

### Upload Behavior

- File is encrypted with AES-256-CTR with a key derived from the Encryption Key seed (this is the twelve to twenty-four words sometimes called a Mnemonic, see BIP39) and an index.
- The file is then encoded with Reed Solomon erasure encoding, expanding the total size to 1 and 2/3 of it’s original size (technically this ratio is adjustable). A shard size is determined at this point, as a multiple of 2MiB (e.g. 2, 4, 8, 16…).
- Each shard is then hashed, with SHA256 and RIPEMD160, and the Client will ask for a location, a Node, to store the shard from the Client’s Node.
- The Client’s Node then selects other Nodes based on reputation and asks many of them concurrently if they are willing to store the data. The Nodes respond and are put into a cache of available mirrors for the data.
- The Client receives in response a Node contact, which includes a nodeID, and IP address, port and token. The token is used to authorize the upload from the Client.
- Each shard for the file repeats steps 4 and 5 concurrently, and data is uploaded to each Node at the same time for faster transfers.
- When each shard is complete an Exchange Report is sent to the Client’s Node with a success or fail status.
- An HMAC is generated from the hash of each shard, this is later used to verify the integrity of the file and that it hasn’t been modified.
- Once complete, the file meta data is sent to the Client’s Node and finalizes the upload.

### Download Behavior

- The Client requests the file meta data from the Client’s Node. This includes the HMAC, the decryption index, the size, and a list of all of the hashes. The integrity of the file can be verified before decryption, using the Encryption Key, HMAC and shard hashes to avoid any   issues as detailed in The Cryptographic Doom Principle.
- The Client then requests the locations of the shards from the Client’s Node.
- The Client’s Node then reaches out to about six of the known Nodes storing that shard, and asks for a retrieval token. Of those that receive a response, a token and the contact details for the Node are then sent to the Client.
- The Client then downloads each encrypted shard directly from each Node to disk at the position it will exist in the file and verifies the hash of the shard.
- For each shard an Exchange Report is sent to Client’s Node to report its success or failure. This information can later be used to improve the ability to retrieve the files.
- If there are any shards that didn’t receive location information, the Client will recover those shards from the Reed-Solomon encoding.
- The file will now be decrypted and the file returned to its original size.
