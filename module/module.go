package module

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
)

type Module struct {
	Name     string                   `json:"name"`
	Versions map[string]ModuleVersion `json:"versions"`
}

type ModuleVersion struct {
	Version string `json:"version"`
	File    string `json:"file"`
	Hash    string `json:"hash"`
}

type ModuleService struct {
	WorkingDir string            `json:"dir"`
	Modules    map[string]Module `json:"modules"`
}

type ModuleReturn struct {
	Name string `json:"name"`
	Hash string `json:"hash"`
	Body string `json:"body"`
}

func (mod ModuleService) Print() {
	fmt.Println(mod.Modules)
}

func (mod ModuleService) GetModule(name, version string) (ModuleReturn, error) {
	module := mod.Modules[name].Versions[version]
	fmt.Println(module.File)
	bts, err := os.ReadFile(mod.WorkingDir + "/" + module.File)
	if err != nil {
		fmt.Println("could not read file: ", err)
		return ModuleReturn{}, err
	}
	base := base64.StdEncoding.EncodeToString(bts)
	ret := ModuleReturn{
		Name: name,
		Hash: module.Hash,
		Body: base,
	}
	return ret, nil

	// search ghost.modules.json for the file
	// if version is specified, use that version, if not, get the latest tagged version...
	// check if the file exists
	// if not return error
	// if so get the bytes
}

func (mod ModuleService) SetModule(name, version, file string) error {
	// check if the file exists
	//
	bts, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	// load the ghost.modules.json
	//

	// check if the module exists
	module := mod.Modules[name]
	if module.Name == "" {
		fmt.Println("Nod found...")
		// create the module
		mod.Modules[name] = Module{
			Name:     name,
			Versions: map[string]ModuleVersion{},
		}
		module = mod.Modules[name]
	}
	if module.Versions[version].Version != "" {
		return fmt.Errorf("Version already exists")
	}
	// create file and get hash
	fileName := name + "-" + version + ".wasm"
	f, err := os.Create(mod.WorkingDir + "/" + fileName)
	if err != nil {
		return err
	}
	_, err = f.Write(bts)
	if err != nil {
		return err
	}
	hashBytes := sha256.Sum256(bts)

	// 3. Encode the byte slice to a human-readable hexadecimal string
	hashString := hex.EncodeToString(hashBytes[:])
	mod.Modules[name].Versions[version] = ModuleVersion{
		Version: version,
		File:    fileName,
		Hash:    hashString,
	}

	mds, err := json.MarshalIndent(mod.Modules, "", "    ")
	err = os.WriteFile(mod.WorkingDir+"/ghost.modules.json", mds, 0666)
	if err != nil {
		return err
	}

	return nil

}

func CreateModuleService() ModuleService {
	// create the folder if not exists
	// dirPath := "/tmp/.ghost-network/modules"
	dirPath := ".ghost-network/modules"
	// dirPath := os.TempDir() + ".ghost-network/modules"
	// if it does exist read the
	_, err := os.Stat(dirPath)

	if err != nil {
		err = os.MkdirAll(dirPath, 0777)
		if err != nil {
			panic(err)
		}
	}

	mods := map[string]Module{}

	mds, err := json.MarshalIndent(mods, "", "    ")

	bts, err := os.ReadFile(dirPath + "/ghost.modules.json")
	if err != nil {
		os.WriteFile(dirPath+"/ghost.modules.json", mds, 0666)
		bts, err = os.ReadFile(dirPath + "/ghost.modules.json")
	}

	err = json.Unmarshal(bts, &mods)
	if err != nil {
		panic(err)
	}
	return ModuleService{
		Modules:    mods,
		WorkingDir: dirPath,
	}
}
