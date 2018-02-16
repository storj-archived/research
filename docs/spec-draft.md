# Storj Network Specification Draft

This a draft and working document for an informal specification of an improved Storj network. It describes the behaviors of each network component, and the interfaces that each expose and use to interact. This is building on ideas defined in both versions of the Storj whitepaper (v1 in 2014 and v2 in 2016), and the latest Storj network v1.2 implementation – with the goal at being decentralized at the foundation.

There are two main groups of responsibilities; Storj Node and Storj Client. A Storj Node is responsible for storing data geographically distributed with other Storj Nodes. This includes a database layer for file pointer and meta data, and a storage layer for larger file shards. A network of Storj Nodes is responsible for maintaining, replicating and repairing data when it’s lost. Storj Nodes pay each other in cryptocurrency STORJ tokens for storage and bandwidth. The Client is responsible for encrypting, erasure encoding and transferring files to Storj Nodes.

-----------------------------------------

## Storj Node

A Storj Node exposes an interface for Clients to communicate to store and retrieve meta data that describes the locations of the shards that make up files, and for downloading shard data in the network. The client layer API is exposed over a TCP/TLS socket with HTTP with JSON requests and responses most often over TLS. The storage layer interface is a JSON-RPC API over a TCP socket over HTTP.  There are RPC commands for storing and retrieving arbitrary data with contracts from a Node, and endpoints for transferring data to Clients in the network over HTTP. Implements SIP5, SIP6 and SIP9

### Behaviors & Services

- A client API TCP server/cluster → A TCP/TLS socket that will respond to HTTP requests based on various HTTP methods and routes  (see Client Endpoints for details). This API is for the purpose of Clients.
- A storage API TCP host server → A public TCP socket that will respond to HTTP requests. based on JSON-RPC commands and routes. Messages are authenticated via ECDSA between Nodes. The purpose of this API is for storing shards of data on the network to be geographically distributed.
- A meta database API server → A public socket for Node communication. This layer is responsible for geographically distributing and storing user, file pointer/meta, and contact data.
- A payout cron worker → The period of storage events is analyzed and a payout sheet and transaction is created and tokens are sent to other Nodes for services.
- A general purpose cron worker → Runs several different scheduled task for data. finalization and maintenance. This includes analyzing storage events and giving them a final status (see SIP9 for more details).
- A shard replication worker → It’s watching for shards to be lost and making an effort to replicate the shard to receive payments for storing the shard. The action should be able to proved to the shard’s user, to receive payment.
- A shard repair worker → Watches for shards to be completely lost and will recover the missing shards so that it will receive payments for having the shard. The action should be able to proved to the shard’s user, to receive payment.

### Service Dependencies

- Storj Nodes → To be able to store data it needs to join a network.
- Ethereum Node/Wallet → For making STORJ token payments.

### Payout Cron Worker Behaviors

- (todo) Needs to be specified.

### Meta Database Behaviors

- (todo) Needs to be specified.

### Shard Replication Worker Behaviors

- (todo) Needs to be specified.

### Shard Repair Worker Behaviors

- (todo) Needs to be specified.

### Identity

- The identity of a Storj Node is based on a 12-24 word phrase (based on BIP39), keys can be derived from this key for various purposes from signing network requests to deriving private keys for an Ethereum wallet. A nodeID is identified by a RIPEMD160(SHA256) hash of a public key.

### Authentication Methods

- Node ECDSA RPC Authentication
  - Authentication with hierarchically deterministic nodeIDs (Please see SIP32 for more details). Each JSON-RPC message is signed between Nodes.
  - The nodeID is a public key hash of the private key used to sign messages.
- Node ECDSA HTTP Authentication
  - Storj Nodes can authenticate by signing requests and including the signature in the header of the HTTP request.
  - The exact details of this authentication are in SIP6 (expanded for additional use cases)
  - The headers x-node-id, x-node-timestamp, x-node-pubkey, x-node-signature are used for authentication.
- Client ECDSA HTTP Authentication
  - Users can authenticate by signing HTTP requests with included HTTP headers, as similar with Node ECDSA Authentication.
  - Uses a different set of HTTP headers: x-client-id, x-client-timestamp, x-client-pubkey, x-client-signature.

### Client API Error Response Codes

- 400 → Bad request. The request supplied invalid arguments.
- 401 → Not authorized. The resource requires authentication.
- 404 → Not found. The resource is not found.
- 420 → Transfer rate limit. This rate limit code is given for uploads and downloads when the rate of bytes transferred exceeds configurable thresholds.
- 429 → Request rate limit. This rate limit code is given when the number of requests per minute exceeds the configurable thresholds per each endpoint.
- 500 → Internal service error. The server experienced an internal error when processing the request.
- 503 → Service unavailable. The service is currently unavailable.

### Client Endpoints

#### Users
- `POST /users` → Creates a user in the network, a user is required to authenticate with many of the other endpoints.
  - Authorization
    Request must be signed by a user
  - Body params
    - pubkey (required, ECDSA 256-bit compressed public key)
    - referralPartner (optional, partner code string, example “filezilla”)
    - email (optional, valid email address)
  - Responses
    - 201 → Successfully created user
    - 400 → Bad request params or user already exists
  - Response body (user object)
    - referralPartner    (partner string)
    - created    (ISO 8601 date string)
    - id (hash of public key)
    - email (email string)
    - deleted (boolean)
- `DELETE /users/:id` → Prepares user to be schedule to be deleted. There is a grace period before complete removal. All data associated with the user should be closed out and no longer kept around, repaired or replicated.
  - Authorization
    - Requires user authentication
  - Request body params
    - redirect (A redirect <url> for the email)
  - Responses
    - 200 → Email has been dispatched
  - Response body
    - (user object)
- `GET /users/confirmations/:token?redirect=<url>` → Confirm different actions such as confirming an email address.
  - Authorization
    - No authentication (aside from knowing the activation token)
  - Responses
    - 200 → Action has been confirmed
    - 302 → Action has been confirmed and redirecting to <url>
  - Response body
    - (user object)

#### Files
- `GET /files?startDate=<timestamp>?tag=<tag>` → List user files sorted by timestamp or by tag.
  - Authorization
    - Requires user authentication
  - Responses
    - 200 → Successfully found results
  - Response body
    - An array of “bucket entry” objects
      - bucket → bucket id
      - mimetype (unused)
      - filename → The encrypted filename
      - frame → frame id that describes all the shards of the file
      - size → The client reported file size
      - id → The file id
      - created → The date the file was created (ISO 8601 date string)
      - hmac → An object with properties for “value” and “type”
- `GET /file-ids/:name` → Get file id by encrypted name
  - Authorization
    - Requires user authentication
  - Request body params
  - Responses
    - 404 → Not found
    - 200 → Successfully found file
  - Response body
    - id → Matching file id
- `GET /files/:file?skip=<number>&limit=<number>&exclude=<node-ids>` → Get pointers for a file to retrieve from Nodes. This request reaches out the the network of Storj Nodes and asks several of them for a shard download token. Multiple requests to this endpoint is necessary to retrieve all of the shards for a file. A best effort is made to retrieve tokens to download the shard however there is a possibility that the shard will no longer be accessible. Clients are designed to handle this case and can recover the shard with erasure encoding. The shard hashes will be returned in the request even if a download token can not be retrieved.
  - Authentication
    - Requires user authentication
    - Optional token based authorization (unused)
  - Responses
    - 404 → Not found
    - 420 → Transfer rate error
  - Response body
    - An array of shard objects
      - token → The download token used to download the shard from the Node
      - hash → The hash of the shard
      - index → The index of the shard in the file
      - size → The size of the shard in bytes
      - node → An object that describes the Node
      - address → The hostname that the Node is reachable
      - port → The port the Node is reachable
      - nodeID → The id of the Node
- `DELETE /files/:file` → Delete a file from a bucket. Will cascade deletion of all  associated data, and marks data as no longer used using storage events stored locally at the Node.
  - Authorization
    - Requires user authentication
  - Responses
    - 204 → Successfully deleted file
    - 404 → File not found
- `GET /files/:file/info` → Gets all information about a file.
  - Authorization
    - Requires user authentication
  - Responses
    - 404 → File not found
    - 200 → Success
  - Response body
    - user → The user id
    - name → Encrypted filename
    - shards → An array of objects
    - hash → The hash of the shard
    - parity → Boolean if the shard is parity
    - bytes → The bytes of the shard
    - nodes → An array of Todes
    - begin → A timestamp in milliseconds
    - end → A timestamp in milliseconds
    - id → The farmer id
    - bytes → Client reported size of file
    - id → file id
    - created → The created timestamp in milliseconds
    - hmac → An object with “type” and “value” properties
    - erasure → An string that describes the erasure encoding used
    - index → A key derivation index for encryption/decryption
    - tags → An array of tags
- `POST /files` → Create a new empty file that will get updated with additional shards.
  - Authorization
    - Requires user authentication
  - Request body params
    - filename → The encrypted filename
    - hmac → An object with “type” and “value” properties. Valid types include “sha512”, “sha256” and “ripemd160”, the value should be a hex string.
    - erasure → A string that indicates the erasure encoding. Valid values include “reedsolomon”.
    - index → A 256-bit number as hex string. Used for derivation of a key for file encryption/decryption.
    - tags → An array of ids or words
  - Responses
    - 400 → Bad request error, some of the body params may have an issue, such as exceeding maximum lengths.
    - 201 → Successfully created file
  - Response body (file object)
- `PUT /files/:file/shards/:index` → Requests to add a shard to a file. This will make a network call to several Storj Nodes to see which Nodes are available to store data. Nodes are selected based on reputation currently using the responseTime metric of a contact. See SIP6 for more details.
  - Authorization
    - Requires user authentication
    - Request body params
      - exclude → A list of nodeIDs to exclude
      - size → The size of the shard in bytes
      - hash → The hash of the shard
      - parity → Boolean if the shard is parity
  - Responses
    - 200 → Success
    - 503 → Unable to get an upload offer from a Node in time
    - 400 → Invalid requests
    - 420 → Transfer rate limit exceeded
  - Response body
    - hash → The hash of the shard
    - token → The authorization token for uploading to the Node
    - node → An object describing where to find the Node
    - address → The publicly reachable IP address of the Node
    - port → The publicly reachable port of the Node
    - nodeID → The nodeID
- `POST /reports` → Creates a new exchange report describing a data transfer between Client and Node. Reports of the message MIRROR_SUCCESS, MIRROR_FAILED, SHARD_UPLOADED and DOWNLOAD_ERROR will all trigger a shard to be replicated. Note: Please see SIP9 for further details on how exchange reports are handled.
  - Authorization
    - Requires user or Node authentication
  - Request body params
    - token → The token used for the shard transfer used to lookup the storage event
    - exchangeResultCode → 1000 for success, and 1100 for failure
    - exchangeResultMessage → The associated message for the report
    - exchangeStart → A unix timestamp in milliseconds when the transfer started
    - exchangeEnd → A unix timestamp in milliseconds when the transfer ended
  - Responses
    - 200 → The report has already been submitted
    - 201 → The storage event was updated with the exchange report
    - 404 → The token does not match any existing storage event
    - 400 → Invalid exchange report
  - Response body
    - (empty object)

#### Node Contacts

- `GET /contacts?address=<address>&skip=<number>&limit=<number>` → Get a list of contacts matching the query. Note: Reputation metrics and other metrics are not included here as this will be data that is only for a specific Node.
  - Authorization
    - No authentication
  - Responses
    - 200 → Successfully found contacts
  - Response body
    - An array of contact objects
      - spaceAvailable → If space is reported to be available
      - port → The port that can be used when connecting to the Node
      - address → The IP address for connecting to the Node
      - protocol → The protocol version of the contact
      - nodeID → The id of the Node
- `GET /contacts/:nodeID` → Get a specific contact by nodeID.
  - Authorization
    - No authentication
  - Responses
    - 404 → Not found
    - 200 → Found contact
  - Response body
    - (contact object)
- `PATCH /contacts/:nodeID` → Update a contact.
  - Authentication
    - Requires Node authentication
  - Request body params
    - address → The UP address that the Node can be reached
    - port → The port the farmer can be reached
    - protocol → The protocol version of the Node
    - spaceAvailable → If the Node has any space available
  - Responses
    - 201 → Contact was updated (bug: should likely be 200 since a new contact wasn’t created)
  - Response body (bug: should be whole contact object)
    - nodeID → The updated nodeID
    - address → The updated address
    - port → The updated port
    - spaceAvailable → The updated spaceAvailable value
- `POST /contacts` → Create a new Node contact
  - Authorization
    - Requires completing proof-of-work challenge
  - Request body params
    - address → The IP address of the Node
    - port → The port of the Node
    - spaceAvailable → A boolean of the space available
    - protocol → The protocol version
  - Responses
    - 200 → The contact was successfully created
  - Response body
    - nodeID →The nodeID of the Node
    - address → The address of the Node
    - port → The port of the Node
- `POST /contacts/challenges` → Gets a proof-of-work challenge to create a new contact. Please see SIP6 for further details on the challenge.
  - Authorization
    - No authentication
  - Responses
    - 201 → Successfully created challenge
  - Response body
    - challenge → A 256-bit hex string to be hashed
    - target → A 256-bit hex string which the resulting hash must be less than

#### Health
- `GET /health` → Get the health of the server
  - Authorization
    - No authentication
  - Responses
    - 503 → Service unavailable
    - 200 → Success
  - Response body
    - (empty object)

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

-----------------------------------------

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

-----------------------------------------

## References

- “Storj Whitepaper v1” https://storj.io/storj2014.pdf
- “Storj Whitepaper v2” https://storj.io/storj.pdf
- “File Encryption and Erasure Encoding Standard (SIP5)” https://github.com/Storj/sips/blob/master/sip-0005.md
- “Farmer Load Balancing Based on Reputation (SIP6)” https://github.com/Storj/sips/blob/master/sip-0006.md
- "Bandwidth Reputation and Accounting (SIP9)” https://github.com/Storj/sips/blob/master/sip-0009.md
- “Basic access authentication” https://en.wikipedia.org/wiki/Basic_access_authentication
- “Elliptic Curve Digital Signature Algorithm (ECDSA)” https://en.wikipedia.org/wiki/Elliptic_Curve_Digital_Signature_Algorithm
- “The Cryptographic Doom Principle” https://www.moxie.org/blog/the-cryptographic-doom-principle/
- “Advanced Encryption Standard (AES)” https://en.wikipedia.org/wiki/Advanced_Encryption_Standard
- “CTR Block cipher mode of operation” https://en.wikipedia.org/wiki/Block_cipher_mode_of_operation#Counter_(CTR)
- “Hash-based message authentication code (HMAC)” https://en.wikipedia.org/wiki/Hash-based_message_authentication_code
- “Hash-based message authentication code” https://en.wikipedia.org/wiki/Galois/Counter_Mode
- “Mnemonic code for generating deterministic keys (BIP39)” https://github.com/bitcoin/bips/blob/master/bip-0039.mediawiki
- “Reed–Solomon error correction” https://en.wikipedia.org/wiki/Reed%E2%80%93Solomon_error_correction
