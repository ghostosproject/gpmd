## Task 003: Create a repository for storing Wasm Module


**Description**

GPMD is responsible for holding all of the WASM modules for the GVMs on the local device.
Modules will be stored as wasm binaries in a specific folder with the name of the module and version in the {module-version.wasm} format. (ex. /tmp/ghost-network/wasm/node-v0.0.1.wasm)

- Create a model for holding the Module information
- A wasm.modules file may be implemented to keep track of the modules on machine

```go
type Module struct {
    Name        string                      `json:"name"`
    Versions    map[string]ModuleVersion    `json:"versions"`
}

type ModuleVersion struct {
    Version string `json:"version"`
    File    string `json:"file"`
    Hash    string `json:"hash"`
}
```

Module Get Call
- get's the binary for the wasm module
- compares with the binary hash with the saved hash value
- returns the file as base64 encoded string along with the name and version being sent


Module Upload/Download for GPMD
- GPMD needs to be able to download modules.
- Initially we can add functionality to GPMD commands to upload wasm modules from a local source
- In the future, we can have functionality to download from github or even DOS (Distributed Object Storage) when it is implemented.

```
$: gpmd wasm upload --name=node --version=v0.0.1 --file=/Users/ghost/node-v0.0.1.wasm
$: gpmd wasm get node@v0.0.1 --location=dos.ghost                       ***not_implemented***
$: gpmd wasm get node@v0.0.1 --location=github.com/ghost/node@v0.0.1    ***not_implemented***
```

**Commits**
task(task-003): add wasm upload functionality & commands (255fbf1df62ffc059975539420a438d7546f0424)