package cmd

func FlatpakInstall(args []string) error {
	if err := validateNames(args, "usage: kk flatpak-inst <app-id> [app-id...]"); err != nil {
		return err
	}
	commandArgs := append([]string{"install", "--noninteractive", "flathub"}, args...)
	return runCommand("flatpak", commandArgs...)
}
