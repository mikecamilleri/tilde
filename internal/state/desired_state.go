package state

// TODO: wrap maps in types as in auths

// newDesiredState is used to avoid nil pointer problems
func newDesiredState() desiredState {
	return desiredState{
		// TODO: ...
	}
}

// desiredState ...
type desiredState struct {
	Features map[FeatureID]*Feature
	// TODO: ...
}
