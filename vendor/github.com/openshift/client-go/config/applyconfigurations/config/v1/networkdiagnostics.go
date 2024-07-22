// Code generated by applyconfiguration-gen. DO NOT EDIT.

package v1

import (
	v1 "github.com/openshift/api/config/v1"
)

// NetworkDiagnosticsApplyConfiguration represents an declarative configuration of the NetworkDiagnostics type for use
// with apply.
type NetworkDiagnosticsApplyConfiguration struct {
	Mode            *v1.NetworkDiagnosticsMode                           `json:"mode,omitempty"`
	SourcePlacement *NetworkDiagnosticsSourcePlacementApplyConfiguration `json:"sourcePlacement,omitempty"`
	TargetPlacement *NetworkDiagnosticsTargetPlacementApplyConfiguration `json:"targetPlacement,omitempty"`
}

// NetworkDiagnosticsApplyConfiguration constructs an declarative configuration of the NetworkDiagnostics type for use with
// apply.
func NetworkDiagnostics() *NetworkDiagnosticsApplyConfiguration {
	return &NetworkDiagnosticsApplyConfiguration{}
}

// WithMode sets the Mode field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Mode field is set to the value of the last call.
func (b *NetworkDiagnosticsApplyConfiguration) WithMode(value v1.NetworkDiagnosticsMode) *NetworkDiagnosticsApplyConfiguration {
	b.Mode = &value
	return b
}

// WithSourcePlacement sets the SourcePlacement field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the SourcePlacement field is set to the value of the last call.
func (b *NetworkDiagnosticsApplyConfiguration) WithSourcePlacement(value *NetworkDiagnosticsSourcePlacementApplyConfiguration) *NetworkDiagnosticsApplyConfiguration {
	b.SourcePlacement = value
	return b
}

// WithTargetPlacement sets the TargetPlacement field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the TargetPlacement field is set to the value of the last call.
func (b *NetworkDiagnosticsApplyConfiguration) WithTargetPlacement(value *NetworkDiagnosticsTargetPlacementApplyConfiguration) *NetworkDiagnosticsApplyConfiguration {
	b.TargetPlacement = value
	return b
}
