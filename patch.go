package patch

import (
	"github.com/pkg/errors"
	"google.golang.org/genproto/protobuf/field_mask"
)

type GetMasker interface {
	GetMask() *field_mask.FieldMask
}

// Patch encompasses the patching logic for some resource. This type is
// intended to be used with composition.
type Patch struct {
	masker GetMasker
	err    error
}

// Err is a getter for the Patch.err field.
func (p Patch) Err() error {
	return p.err
}

// SetErr is a setter for the Patch.err field.
func (p *Patch) SetErr(err error) {
	p.err = err
}

// NewPatch initializes a Patch with a masker. The masker is typically a grpc
// request object that utilizes the FieldMask type.
func NewPatch(masker GetMasker) *Patch {
	return &Patch{
		masker: masker,
		err:    nil,
	}
}

// ValidatePaths ensures that a set of paths exist within p.masker.
func (p Patch) ValidatePaths(paths ...string) error {
	var m = make(map[string]bool)
	for _, path := range paths {
		m[path] = false
	}
	for _, mask := range p.masker.GetMask().GetPaths() {
		state, ok := m[mask]
		if !ok {
			continue
		}
		if state {
			continue
		}
		m[mask] = true
	}
	for path, exists := range m {
		if exists == false {
			return errors.Errorf("failed to ValidatePaths\tpath \"%s\" not in mask", path)
		}
	}
	return nil
}
