## Feature 002: Add ability for GPMD to facilitate connection between two nodes.

GPMD accepts a connection request from a node (gvm1).
Looks for the node(gvm2) gvm1 is trying to connect to.
If the node is alive and can be accessed, GPMD will send connection information for gvm2 to gvm1

Connection Information
- address
- name
- id
- port

**Tasks**

**Dependencies**
Ghost Message Protocol [[gmp/feature-001](ghost-network-common/feature-001.md)  :green_circle:

**Commits**

feat(feature-002): add node connectivity support (d2cc39416ec4b8891b602e7a36569e16ef492513)