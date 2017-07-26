package ggprov

import "log"

// AptUpdate run apt update
func AptUpdate() error {
	log.Println("Running apt update")
	return RunCommand("apt", []string{"update"})
}

// AptInstall run apt install and install the specified package(s)
func AptInstall(packages []string) error {
	log.Println("Running apt install for:", packages)
	return RunCommand("apt", append([]string{"install", "-y"}, packages...))
}
