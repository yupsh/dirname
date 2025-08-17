package opt

// Boolean flag types with constants
type ZeroFlag bool
const (
	Zero   ZeroFlag = true
	NoZero ZeroFlag = false
)

// Flags represents the configuration options for the dirname command
type Flags struct {
	Zero ZeroFlag // End output with NUL character instead of newline
}

// Configure methods for the opt system
func (z ZeroFlag) Configure(flags *Flags) { flags.Zero = z }
