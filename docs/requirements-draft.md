# Storj Network Requirements

## Primary Queries
  - A user with only the knowledge of an encryption key (e.g. “uncover earn cigar humor stable tuition obvious car weird learn genuine million”) should be able to access all file data from any node. This data should only be available to modify by the user with that encryption key.
  - A node should be able to gain knowledge of shards that are in need of replication or repair, and without central coordination, and be able to replicate and repair autonomously and verified by other nodes in the network. The current state of replication should be available for all nodes to discovery.
  - A node should be able to get the current contact information for another node, for example the IP address and port.

## Fault Tolerances
 - The network should be able to tolerate less than 51% network failure, for example:
   - Geographic based failures, such as natural disasters; floods, fire and attacks.
   - Social related failures; social engineering, betrayal, sabotage, sybil, eclipse and other attacks.
   - Engineering related failures; bugs, malware and viruses.

## Scaling
 - Adding more nodes to the network should increase the overall ability to store data, from the capacity of actual storage available, to bandwidth, and the maximum rate that new files can be stored.
- The work that every node must participate should be minimal.

## Nodes Joining
- New nodes can join the network autonomously and without central coordination as to encourage broad replication of data across diverse resources.
- The cost to join the network should be equivalent to the work that would be necessary when they abruptly leave the network. Alternatively, leaving the network and data unexpectedly should come with a cost that would be equivalent to the resources necessary to recover the data that was lost.
- Cryptocurrency payout incentives can encourage people to contribute resources to the network and share computing and storage resources.

## Data Structures
- File Descriptor
  - Authorization
    - Can only be written and modified by authenticating as the user
  - Data
    - user (able to be cryptographically authenticated, e.g. a public key hash)
    - filename (encrypted)
    - created (timestamp)
    - index (used for key derivation and encryption)
    - hmac (used for authentication, depends the type/algo and the value)
    - bytes (the size of the file in bytes)
    - data (encrypted data, less than 2MB)
    - shards (describes a sequence of shard hashes that make up a larger file, greater than 2MB)
    - tags (a list of tags for searching files)
- Shard Description
  - Authorization
    - Requires network consensus with auditing to add and remove the list of nodes that are currently storing a particular shard. The shard descriptor can be read by all nodes so that the network can automatically replicate and repair shards, and cryptocurrency payments can be made for storing data.
  - Data
    - hash (the cryptographic secure hash of the shard)
    - bytes (the size of the shard in bytes)
    - erasure (information necessary to recovery this shard, such as the shards necessary)
    - nodes (a description of the nodes that currently hold this data)
    - begin (when the data was originally stored for a farmer)
    - end (when the data was no longer stored for a farmer)
    - id (an identify for the farmer node that can be used to find the location in the network)
- Node Descriptor
  - Authorization
    - Can only be created or modified by authenticating as the node, and can be read by anyone. There should be some cost associated with creating this data, for example posting a cryptocurrency bond or performing proof-of-work hashing, such as defined in hashcash.
  - Data
    - id (a hash of the extended public key)
    - addresses (the IP and port addresses for the node, including IPv4 and IPv6)
    - spaceAvailable (boolean if the node is accepting data)
    - protocol (the current protocol version)
