package state

import (
	"encoding/gob"
	"os"
)

// Persist ...
func (s *State) Persist(path string) error {
	// open the file for writing
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	// encode and write to the file
	// enc := json.NewEncoder(f)
	enc := gob.NewEncoder(f) // I ❤️ Go! What an easy change!
	if err = enc.Encode(s); err != nil {
		return err
	}

	return nil
}

// Restore ...
func (s *State) Restore(path string) error {
	// open the file for reading
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	// decode into the state
	// dec := json.NewDecoder(f)
	dec := gob.NewDecoder(f) // I ❤️ Go! What an easy change!
	if err := dec.Decode(s); err != nil {
		return err
	}

	return nil
}
