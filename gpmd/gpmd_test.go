package gpmd

import (
	"os"
	"testing"
)

func TestCreateGPMD(t *testing.T) {
	t.Run("No P2P", func(t *testing.T) {
		pmd := CreateGPMD(false, "", "")
		if pmd.P2PHost != nil {
			t.Errorf("Host created when non is expected")
		}
	})

	t.Run("P2P", func(t *testing.T) {
		pmd := CreateGPMD(true, "8282", "./private.key")
		if pmd.P2PHost == nil {
			t.Errorf("Host was not created when one was expected")
		}
		if len(pmd.Networks) != 1 {
			t.Errorf("One network was expected but non were given")
		}
		os.Remove("./private.key")
	})
}

func TestSetupP2P(t *testing.T) {
	id := ""
	t.Run("No Private Key", func(t *testing.T) {
		hst := setupP2P("8282", "")
		if hst == nil {
			t.Errorf("Bad Host")
		}
		id = hst.ID().String()
		hst.Close()
	})

	t.Run("Private Key", func(t *testing.T) {
		hst := setupP2P("8282", "./private.key")
		if hst == nil {
			t.Errorf("Bad Host")
		}
		if hst.ID().String() != id {
			t.Errorf("Id not correct")
		}
		hst.Close()
		os.Remove("./private.key")
	})
}
