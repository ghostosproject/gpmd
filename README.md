### Ghost Port Mapper Daemon (GPMD)

The main role of GPMD in Phase One is to orchestrate communication between local GVM Nodes. When a GVM is created, it will reach out to GPMD and register itself. If the registration is successfull, the GVM becomes a node on the local network. The Node is now able to be discover by other nodes on the network as well as be discovered. 
GPMD will orchestrate the connection between these two nodes. GPMD only orchestrates the connection. Once a conection is established, all communication is directly between the two nodes.


All of the features, tasks and bugs are documented in the [handle](https://github.com/ghostosproject/gpmd/tree/handle) branch of this repository.

For a deep dive into the internals of GVM as well as other pieces of The Ghost Network, visit [ghostosproject.github.io](https://ghostosproject.github.io/ghost-network).