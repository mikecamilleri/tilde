package state

import (
	"encoding/json"
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
	enc := json.NewEncoder(f)
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
	dec := json.NewDecoder(f)
	if err := dec.Decode(s); err != nil {
		return err
	}

	return nil
}
