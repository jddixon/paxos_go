# paxos_go

A library enabling the achievement of consensus within a cluster of 
independent computational units with unique cryptographic 
identities.  Intra-cluster communications pass over a full mesh of inter-node
connections.  Clients can query the cluster or propose new commands which 
will pass through the consensus process and upon agreement will become part
of the cluster's state.

## Clusters

As the term is used here, a **cluster** is a group of coooperating processors,
**servers**.  The servers are driven by commands issued by **clients**; the
clients are not themselves part of the cluster and generally do not communicate
with one another.  Clients issue **commands**.  Commands are not issued in 
any particular order and may take an arbitrary time to deliver.  While a 
command may never be corrupt, it may be lost or delivered more than once.

## Consensus in the Back End

Each server has at lease one attached **state machine** whose state is 
determined by 
the sequence of commands passed to it.  Servers cooperate to achieve 
**consensus**, which means that each server eventually delivers the same
sequence of commands to its state machine.  (More formally, at any given
time, if *a* is the sequence of commands delivered to state machine *A*
and *b* is the sequence of commands delivered to state machine *B*, then
we are guaranteed that either *a* is a prefix of *b* or *b* is a prefix of
*a* or *a* is identical to *b*.)

A state machine consists of a number of variables.  In paxos_go, one of 
these is *state*, a 64-bit unsigned value.  Each execution of a command
causes the state to be incremented.  Upon successful execution of a command,
the cluster always returns at least the state to the client which requested
execution of the command.

## Applications

A set of state machines as described above is referred to here as an
**application**. An application runs on a single cluster.  Each server in
the cluster has a state machine specific to the application.  Each such
application uses a variant of the Paxos protocol to achieve consensus.  


## Application Front End

In paxos_go, only back-end clients may originate commands which alter the
state of an application's state machine.  Commands to a server contain
enough information to identify the back-end client and the application.
Such commands also contain a hash which can be used to verify the integrity
of the command.  If verification fails, the message (the command) is 
discarded.

External clients to the application generally communicate with servers 
over a different set of listening ports.  External clients need not share
the same model of the data as the consensus servers, except that an 
external client will always be able to obtain the 64-bit unsigned value
identifying the current state of the state machine.

As an example, consider a set of consensus servers which cooperate to appear 
to the
outside world as a group of name servers.  Each name server is responsible
for a set of zone files and will answer standard DNS queries about each of
its zones.  

The back-end cluster can store each such zone file by its content key.  
To the cluster, the state machine maps zone names (the fully qualified
domain name of a zone file) to a hash.  The validity of the hash is 
confirmed by hashing the file.   

## On-line Documentation

More information on the **paxos_go** project can be found [here](https://jddixon.github.io/paxos_go)
