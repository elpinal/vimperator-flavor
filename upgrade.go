package main

var cmdUpgrade = &Command{
	Run:       runUpgrade,
	UsageLine: "upgrade ",
	Short:     "Upgrade",
	Long: `

	`,
}

func init() {
	// Set your flag here like below.
	// cmdUpgrade.Flag.BoolVar(&flagA, "a", false, "")
}

// runUpgrade executes upgrade command and return exit code.
func runUpgrade(args []string) int {

	return 0
}
