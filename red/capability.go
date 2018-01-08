package red

// Authentication capabilities
const (
	CapabilityAuthSpice uint8 = 1
	CapabilityAuthSASL  uint8 = 2
)

// Capability is a bitwise capability set
type Capability uint32

// Test whether bit i is set.
func (c *Capability) Test(i uint32) bool {
	if i >= 32 {
		return false
	}
	return *c&(1<<i) > 0
}

// Set bit i to 1
func (c *Capability) Set(i uint32) *Capability {
	if i >= 32 {
		return c
	}
	*c |= 1 << i
	return c
}

// Clear bit i to 0
func (c *Capability) Clear(i uint32) *Capability {
	if i >= 32 {
		return c
	}
	*c &^= 1 << i
	return c
}

// SetTo sets bit i to value
func (c *Capability) SetTo(i uint32, value bool) *Capability {
	if value {
		return c.Set(i)
	}
	return c.Clear(i)
}

// Flip bit at i
func (c *Capability) Flip(i uint32) *Capability {
	if i >= 32 {
		return c
	}
	*c ^= 1 << i
	return c
}
