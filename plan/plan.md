## Ghost Port Manager Daemon
The Ghost Port Manager Daemon (GPMD) is the network manager for the Ghost Network. It is responsible for managing the connections between virtual machines as well as access to global resources.

### Phase One - Local Ghost Network ([P001](../../.project/project.md#phase-one-local-ghost-network (P001)))

**Purpose**: GPMD should be able to register newly created Virtual Machines, and connect two machines together.

**Features**
- [x] ([F001](feature/feature-001.md)) Create the structure for keeping track of Virtual Machines on the device.
- [x] ([F002](feature/feature-002.md)) Add ability for GPMD to facilitate connection between two nodes.
- [x] ([F003](feature/feature-003.md)) Add Node Discovery Response
- [x] ([F004](feature/feature-004.md)) Add Wasm Module Request support
- [x] ([F005](feature/feature-005.md)) Add detached mode for GMPD & other flags

**Bugs**
- [x] ([B001](bug/bug-001.md)) gpmd wasm upload command not exiting connection when complete

**Dependencies**
- [x] GMP - Ghost Network Common (Ghost Message Protocol) v0.0.5