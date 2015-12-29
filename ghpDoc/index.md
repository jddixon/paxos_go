<h1 class="libTop">paxos_go</h1>

A library enabling the achievement of consensus within a cluster of
independent computational units with unique cryptographic
identities.  Intra-cluster communications pass over a full mesh of
intra-cluster connections.  **Commanders** can query the cluster or propose
new commands which
will pass through the consensus process and upon agreement will become part
of the cluster's state.

## Clusters

As the term is used here, a **cluster** is an
[xLattice cluster](https://jddixon.github.io/xlCluster_go),
a group of coooperating processors,
**servers**.  Each member of the cluster listens on at least two TCP/IP
connections: one is used for communications with other cluster members
and the second is used for communicating with clients.

Paxos servers are driven by commands issued by their clients, **commanders**.
The commanders are not themselves part of the cluster and generally do not
communicate directly with one another.  Commanders issue **commands**
to the cluster.  Commands are not issued in
any particular order and may take an arbitrary time to deliver.  While a
command may never be corrupt, it may be lost or delivered more than once.

## Consensus: the Back End

Each server has at least one attached **state machine** whose state is
determined by
the sequence of commands passed to it.  Servers cooperate to achieve
**consensus**, which means that each server eventually delivers the same
sequence of commands to its state machine.  (More formally, at any given
time, if *a* is the sequence of commands delivered to state machine *A*
and *b* is the sequence of commands delivered to state machine *B*, then
we are guaranteed that either *a* is a prefix of *b* or *b* is a prefix of
*a* or *a* is identical to *b*.)

A state machine consists of a number of variables.  In paxos_go, one of
these is **stateNdx**, a 64-bit unsigned value.  Each execution of a command
causes `stateNdx` to be incremented.  Upon successful execution of a command,
the cluster always returns at least `stateNdx` to the commander which requested
execution of the command.

## Applications

A set of state machines as described above is referred to here as an
**application**. An application runs on a single cluster.  Each server in
the cluster has a state machine specific to the application.  Each such
application uses a variant of the Paxos protocol to achieve consensus.
Each application running on a cluster has a unique identifier, its **appNdx**.

## pktComms

If there are **N** servers in a cluster, then at boot time each server
opens a connection to each of its **N-1** peers.  The low-level code in
`pktComms/` manages these connections.  At boot, each server sends a **Hello**
to each of its peers, using the protocol specified in `pktComms/p.proto`.
The `Hello` message contains

* a message number (**MsgN**), which is always 1;
* the 32-byte nodeID of the sender;
* the RSA public key used for digital signatures (**SigPubKey**);
* the RSA public key used for encoding communications (**CommsPubKey**);
* the address the server listens on (its ipV4 address and port number in string form);
* optionally a **Salt**; and
* the **digital signature** derived from the 256-bit SHA hash of the fields present.

If the peer accepts the `Hello`, it will reply with an **Ack**, as specified
in the protocol, and leave the connection open.  If there was an error in
the `Hello` message, the peer with reply with an **Error** message, and close
the connection.

Normally there will be no error and the connection is left open.  From
this point the server periodically sends the peer a **KeepAlive** message,
to which the peer will reply with either an `Ack` or an `Error` message.  In
the normal case, the reply is an `Ack` and the connection remains open; if
there was an error, the peer replies with an `Error` message and the connection
is closed.

At any appropriate point, instead of sending a `KeepAlive`, the server can
send a **Bye** message.  If this is well-formed, the peer will reply with
an `Ack` and close the connection.

## Application Front End

In paxos_go, only back-end commanders may originate commands which alter the
state of an application's state machine.  Commands to a server contain
enough information to identify the back-end commander and the application.
Such commands also contain a digital signature which can be used to verify
the integrity
of the command.  If verification fails, the message (the command) is
discarded and an error message returned.

External clients to the application generally communicate with servers
over a different set of listening ports.  External clients need not share
the same model of the data as the consensus servers, except that an
external client will always be able to obtain `stateNdx`, the 64-bit unsigned
value identifying the current state of the state machine.

As an example, consider a set of consensus servers which cooperate to appear
to the
outside world as a group of
[name servers.](https://en.wikipedia.org/wiki/Name_server)
Each name server is responsible
for a set of zone files and will answer standard DNS queries about each of
its zones.

The back-end cluster can store each such zone file by its content key.
To the cluster, the state machine maps zone names (the fully qualified
domain name of a zone file) to a hash.  The validity of the hash is
confirmed by hashing the file.

## References

[Wikipedia](https://en.wikipedia.org/wiki/Paxos (computer science))
has a good description of the Paxos family of algorithms and many good
references.

## Project Status

A rough spec.

