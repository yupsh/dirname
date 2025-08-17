package command

type ZeroFlag bool

const (
	Zero   ZeroFlag = true
	NoZero ZeroFlag = false
)

type flags struct {
	Zero ZeroFlag
}

func (z ZeroFlag) Configure(flags *flags) { flags.Zero = z }
