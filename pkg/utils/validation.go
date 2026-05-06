// Copyright (c) 2020 Red Hat, Inc.
// Copyright Contributors to the Open Cluster Management project

package utils

import (
	"fmt"

	operatorsv1 "github.com/stolostron/multiclusterhub-operator/api/v1"
)

// ValidationResult contains the outcome of MCH validation
type ValidationResult struct {
	Valid   bool
	Reason  string
	Message string
}

// MchMeetsRuntimePreconditions checks if the MCH CR meets the minimum
// requirements for the operator to proceed with reconciliation.
// Returns false if:
//   - AvailabilityConfig is invalid
//   - Required fields are missing
//   - Deprecated fields are set (if blocking is enabled)
func MchMeetsRuntimePreconditions(mch *operatorsv1.MultiClusterHub) ValidationResult {
	// Check AvailabilityConfig
	if !operatorsv1.AvailabilityConfigIsValid(mch.Spec.AvailabilityConfig) {
		return ValidationResult{
			Valid:   false,
			Reason:  "InvalidAvailabilityConfig",
			Message: fmt.Sprintf("availabilityConfig must be 'High' or 'Basic', got: '%s'", mch.Spec.AvailabilityConfig),
		}
	}

	// Check required fields
	if mch.Spec.LocalClusterName == "" {
		return ValidationResult{
			Valid:   false,
			Reason:  "MissingRequiredField",
			Message: "spec.localClusterName is required",
		}
	}

	// Validate local cluster name length
	if len(mch.Spec.LocalClusterName) >= 35 {
		return ValidationResult{
			Valid:   false,
			Reason:  "InvalidLocalClusterName",
			Message: "spec.localClusterName must be shorter than 35 characters",
		}
	}

	return ValidationResult{Valid: true}
}

// HasDeprecatedFields returns true if any deprecated spec fields are set
func HasDeprecatedFields(mch *operatorsv1.MultiClusterHub) bool {
	if mch.Spec.SeparateCertificateManagement {
		return true
	}
	if mch.Spec.Hive != nil {
		return true
	}
	if mch.Spec.Ingress != nil {
		return true
	}
	if mch.Spec.CustomCAConfigmap != "" {
		return true
	}
	if mch.Spec.EnableClusterBackup {
		return true
	}
	if mch.Spec.EnableClusterProxyAddon {
		return true
	}
	return false
}

// GetDeprecatedFieldNames returns a list of deprecated field names that are currently set
func GetDeprecatedFieldNames(mch *operatorsv1.MultiClusterHub) []string {
	var fields []string

	if mch.Spec.SeparateCertificateManagement {
		fields = append(fields, "spec.separateCertificateManagement")
	}
	if mch.Spec.Hive != nil {
		fields = append(fields, "spec.hive")
	}
	if mch.Spec.Ingress != nil {
		fields = append(fields, "spec.ingress")
	}
	if mch.Spec.CustomCAConfigmap != "" {
		fields = append(fields, "spec.customCAConfigmap")
	}
	if mch.Spec.EnableClusterBackup {
		fields = append(fields, "spec.enableClusterBackup")
	}
	if mch.Spec.EnableClusterProxyAddon {
		fields = append(fields, "spec.enableClusterProxyAddon")
	}

	return fields
}

// ValidateComponentConfigs checks that component configurations are valid
func ValidateComponentConfigs(mch *operatorsv1.MultiClusterHub) ValidationResult {
	if mch.Spec.Overrides == nil || len(mch.Spec.Overrides.Components) == 0 {
		return ValidationResult{Valid: true}
	}

	// Check for duplicate components
	seen := make(map[string]bool)
	for _, cc := range mch.Spec.Overrides.Components {
		if seen[cc.Name] {
			return ValidationResult{
				Valid:   false,
				Reason:  "DuplicateComponent",
				Message: fmt.Sprintf("duplicate component config for: %s", cc.Name),
			}
		}
		seen[cc.Name] = true

		// Check if component is valid
		if !operatorsv1.ValidComponent(cc, operatorsv1.MCHComponents) {
			return ValidationResult{
				Valid:   false,
				Reason:  "InvalidComponent",
				Message: fmt.Sprintf("invalid component: %s is not a known MCH component", cc.Name),
			}
		}
	}

	return ValidationResult{Valid: true}
}
