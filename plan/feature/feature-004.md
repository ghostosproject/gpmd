## Feature 004: Add Wasm Module Request support

A Node should be able to request a WASM module in order to run.
WASM modules are stored on device and managed by GPMD.


**Tasks**
- [x] ([T003](task/task-003.md)) Create a repository for storing Wasm Module
- [x] ([T004](task/task-004.md)) Handle request from GVM for Wasm Module

**Dependencies**
Ghost Message Protocol [[gmp/feature-004](ghost-network-common/feature-004.md)  :green_circle:

**Commits**
task(task-003): add wasm upload functionality & commands (255fbf1df62ffc059975539420a438d7546f0424)
task(task-004): add wasm module request support (eeadaec8dbe056f5eba47e98c829b87167c5c3b2)