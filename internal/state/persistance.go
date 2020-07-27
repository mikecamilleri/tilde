package state

import (
	"encoding/gob"
	"os"
)

// Persist ...
// TODO: open and write to copies, then move
func (s *State) Persist(path string) error {
	// open the files for writing and defer close
	af, err := os.Create(path + "/auths.gob")
	if err != nil {
		return err
	}
	defer af.Close()

	cf, err := os.Create(path + "/current.gob")
	if err != nil {
		return err
	}
	defer cf.Close()

	// lock the state and defer unlock
	s.Lock()
	defer s.Unlock()

	// encode and write to the file
	afEnc := gob.NewEncoder(af)
	if err := afEnc.Encode(&s.current); err != nil {
		return err
	}

	cfEnc := gob.NewEncoder(cf)
	if err := cfEnc.Encode(&s.current); err != nil {
		return err
	}

	return nil
}

// Restore ...
func (s *State) Restore(path string) error {
	// open the files for reading and defer close
	af, err := os.Open(path + "/auths.gob")
	if err != nil {
		return err
	}
	defer af.Close()

	cf, err := os.Open(path + "/current.gob")
	if err != nil {
		return err
	}
	defer cf.Close()

	// lock the state and defer unlock
	s.Lock()
	defer s.Unlock()

	// decode into the state
	// dec := json.NewDecoder(f)
	afDec := gob.NewDecoder(af) // I ❤️ Go! What an easy change!
	if err := afDec.Decode(&s.auths); err != nil {
		return err
	}

	cfDec := gob.NewDecoder(cf) // I ❤️ Go! What an easy change!
	if err := cfDec.Decode(&s.current); err != nil {
		return err
	}

	return nil
}
