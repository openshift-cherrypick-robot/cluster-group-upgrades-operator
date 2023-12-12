package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSafeResourceNames(t *testing.T) {

	name := "cnfdf18-new-common-cnfdf18-looooong-subscriptions-policy"
	namespace := "ztp-install"
	safeName := NewSafeResourceName(name, "", MaxPolicyNameLength, len(namespace)+1)
	assert.Equal(t, MaxPolicyNameLength-len(namespace)-1, len(safeName))
	assert.Equal(t, name[:MaxPolicyNameLength-len(namespace)-7]+"-", safeName[:MaxPolicyNameLength-len(namespace)-6])

	name = "cnfdf18-new-common-cnfdf18-looooong-subscriptions-policy-config"
	safeName = NewSafeResourceName(name, "", MaxPolicyNameLength, 0)
	assert.Equal(t, MaxPolicyNameLength, len(safeName))
	assert.Equal(t, name[:MaxPolicyNameLength-7]+"-", safeName[:MaxPolicyNameLength-6])

	name = "cnfdf18-new-common-cnfdf18-looooong-subscriptions-policy-placement"
	safeName = NewSafeResourceName(name, "", MaxObjectNameLength, 0)
	assert.Equal(t, len(name)+6, len(safeName))
	assert.Equal(t, name+"-", safeName[:len(name)+1])

	name = "cnfdf18-new-common-cnfdf18-looooong-subscriptions-policy"
	safeName = NewSafeResourceName(name, "kuttl", MaxPolicyNameLength, len(namespace)+1)
	assert.Equal(t, MaxPolicyNameLength-len(namespace)-1, len(safeName))
	assert.Equal(t, name[:MaxPolicyNameLength-len(namespace)-7]+"-kuttl", safeName)

	name = "cnfdf18-new-common-cnfdf18-looooong-subscriptions-policy-config"
	safeName = NewSafeResourceName(name, "kuttl", MaxPolicyNameLength, 0)
	assert.Equal(t, MaxPolicyNameLength, len(safeName))
	assert.Equal(t, name[:MaxPolicyNameLength-6]+"-kuttl", safeName)

	name = "cnfdf18-new-common-cnfdf18-looooong-subscriptions-policy-placement"
	safeName = NewSafeResourceName(name, "kuttl", MaxObjectNameLength, 0)
	assert.Equal(t, len(name)+6, len(safeName))
	assert.Equal(t, name+"-kuttl", safeName)

	name = "cgu-sriov-cloudransno-site9-spree-lb-du-cvslcm-4.14.0-rc.4-config"
	safeName = NewSafeResourceName(name, "", MaxPolicyNameLength, 0)
	assert.Equal(t, MaxPolicyNameLength-1, len(safeName))
	assert.Equal(t, name[:MaxPolicyNameLength-7]+"-", safeName[:MaxPolicyNameLength-6])

	name = "cgu-sriov-cloudransno-site9-spree-lb-du-cvslcm-4.14.0-rc.4"
	safeName = NewSafeResourceName(name, "", MaxPolicyNameLength, 8)
	assert.Equal(t, MaxPolicyNameLength-9, len(safeName))
	assert.Equal(t, name[:MaxPolicyNameLength-15]+"-", safeName[:MaxPolicyNameLength-14])
}
