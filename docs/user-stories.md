https://medium.com/@jonatisokon/a-framework-for-user-stories-bc3dc323eca9

## INFO
* Storj Node = Peer in peer to peer network (P2P)
* compensation = payment in STORJ token


# A User should be able to upload a file to the Storj Network
As a user
I want to upload a file to Stroj
Because I want recover from a local storage failure

**GIVEN** I have made the file accessable to my local Storj Node
  **AND** my local Storj Node has encryptd and sharded a file
  **AND** my local Stroj Node has connected to a remote Storj Nodes  
 **WHEN** I push different shards of equal size to different Storj Nodes 
 **THEN** the Storj Node finds other Storj Nodes to satisfy 2f + 1 replicas 
  **AND** waits for acks from it async shard push, within a time limit
  **AND** responds to the origin Storj Node of the acks, before or after a time limit is met
---
# A User should be able to stream a file to the Storj Network
As a user
I want to upload a file to Stroj
Because I want stream a file from the network to a remote server

**GIVEN** I have made the file accessable to my local Storj Node
 **WHEN** my Storj Node shards the file for upload 
 **THEN** my Storj Node encodes equal sized shards of the file with a stream encoding for codata
  **AND** uses a burst length that can satisfy 2f + 1 replicas
---
# A User should be able to access a file from the Storj Network from any Storj Node
As a user
I want to upload a file to Stroj
Because I want copy data from one device to another

**GIVEN** I have made the file accessable to my local Storj Node
  **AND** my local Storj Node has encryptd and sharded a file
  **AND** my local Stroj Node has connected to a remote Storj Nodes  
 **WHEN** I provie my private key or mnemonic code to the remote Storj Node 
 **THEN** the remote Storj Node should be able to query the network based on my public key
  **AND** recover my profile and files
---
# A Storj Node should receive STORJ Tokens
As a Storj Node
I want to receive compensation 
Because I shared my resources

**GIVEN** I have shared resources via my local Storj Node
 **WHEN** a payment interval occurs
 **THEN** I should be able to prove the existences of shards for the interval
---
# A Storj Node should receive STORJ Tokens for finding a shard
As a Storj Node
I want to find data for another Storj Node
Because that will increase my reputation and compensation for compute and bandwith resources

**GIVEN** I have provable meta data on my peers (relative latency, tags, current shards, ...)
 **WHEN** I get a request/ query for a specific shard from another Storj Node peer or not 
 **THEN** I pass the request message to the best peer 
---
# A Storj Node should receive STORJ Tokens for storing shards
As a Storj Node
I want to prove I have data
Because that I will be compensated for what I can prove

**GIVEN** I have a shard 
 **WHEN** get a request for proof of a shard 
 **THEN** I provide the require the data to satisfy the inquiry
---
# A Storj Node should receive STORJ Tokens for proactive replication
As a Storj Node
I want to replicate shards
Because that will increase my bandwith compensation

**GIVEN** I have a shard and list of nodes that claim to store the same shard
 **WHEN** I can prove that it does not satisfy the replica count (remote Storj Node failure/timeout)
 **THEN** I find a suitable remote Storj Node and push the shard
---
# A Storj Node should receive STORJ Tokens for proactive replication
As a Storj Node
I want to prove the data I have is not replcated
Because that will increase my bandwith compensation

**GIVEN** I have a shard and list of nodes that claim to store the same shard
 **WHEN** get responses within at time limit from audit request of proof of a shard
 **THEN** I compare my shard to the responses and report my findings
  **AND** I make sure the audit results satisfy the replication minimum
---
# A Storj Node should receive STORJ Tokens for audits of remote Storj Nodes
As a Storj Node
I want to prove the other Storj Nodes have data
Because that will increase my bandwith and storage compensation

**GIVEN** I have free space
 **WHEN** I request a shard from other Storj Nodes
 **THEN** I receive new a shard
  **AND** I increase the replication of the shard
---
# A Storj Node should pass a guided tour to join the network
As a Storj Node
I want to prove I can store data and provied meta data
Because my reputaion will increase which increase my odds to join the network 
After joining I will be compensated for resources 

**GIVEN** I want to join the Storj Network
  **AND** I have a connection to a remote Storj Node
 **WHEN** I get replication request from a remote Storj Node
 **THEN** I prove the replication happend to the remote Storj Node
  **AND** I connect to a different Storj Node given to me by previous Storj Node 
  **AND** repeat until I have provided enough meta data 



 User Stories
============

## Users paying for storage
1. As a user paying for storage
   I want my data to be available at any time
   So that I can download anything I've uploaded

## Users running farmer nodes
1. As a farmer
   I want to be able to independently track and/or calculate performance and billing-related data
   So that I can trust and act on that information

1. As a farmer
   I want to be paid for what I store
   So that I'm not wasting valuable space

## Third-party developers (using storj)
1. As a third-party developer
   I want (user access control / app ids / permissions)
   So that I don't have to (authentication / registration / key management)

1. As a third-party developer
   I want to be able to use storj in multiple languages
   So that I have options
