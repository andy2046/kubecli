package cmd

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path"
	"regexp"

	"github.com/andy2046/kubecli/internal/types"
)

const (
	kubectl   = "kubectl"
	usagePath = "path of kube config file, default to `$HOME/.kube/config`"
)

var (
	defaultKubeConfigPath = path.Join(userHomeDir(), "./.kube/config")
	kubectlCmd            = []string{"config", "unset"}
	// configCommand is the config subcommand.
	configCommand  *flag.FlagSet
	logger         = log.New(os.Stdout, "kubecli 🍻  ", 0)
	kubeConfigPath string
	notFound       = "-"
)

func init() {
	configCommand = flag.NewFlagSet("config", flag.ExitOnError)
}

// Parse validate flags.
func Parse() {
	// get kube config path from flag or env.
	configCommand.StringVar(&kubeConfigPath, "path", notFound, usagePath)
	if kubeConfigPath == notFound {
		kubeConfigPath = getEnv("kube-config-path", notFound)
		if kubeConfigPath == notFound {
			kubeConfigPath = defaultKubeConfigPath
		}
	}

	// verify that subcommand has been provided
	// os.Arg[0] is main command, os.Arg[1] is the subcommand
	if len(os.Args) < 2 {
		printExit("subcommand required")
	}

	// switch on the subcommand
	switch os.Args[1] {
	case "config":
		configCommand.Parse(os.Args[2:])
	case "-h", "-help", "--help":
		PrintUsage()
		return
	default:
		printExit("unknown subcommand")
	}

	// extract vars
	if configCommand.Parsed() {
		args := configCommand.Args()
		n := configCommand.NArg()

		// verify that SUBCOMMAND for config command has been provided
		// args[0] should be {current-context|delete-cluster|delete-context|delete-user|get-clusters|get-contexts|get-users}
		if n < 1 {
			printExit("SUBCOMMAND for config command required")
		}

		switch args[0] {
		case "current-context":
			currentContext()
		case "get-users":
			getUsers()
		case "get-clusters":
			getClusters()
		case "get-contexts":
			getContexts()
		case "delete-cluster":
			if n < 2 {
				printExit("cluster NAME required")
			}
			deleteCluster(args[1:])
		case "delete-context":
			if n < 2 {
				printExit("context NAME required")
			}
			deleteContext(args[1:])
		case "delete-user":
			if n < 2 {
				printExit("user NAME required")
			}
			deleteUser(args[1:])
		default:
			printExit("unknown command")
		}
	}
}

func printExit(str string) {
	logger.Println(str)
	PrintUsage()
	os.Exit(1)
}

// PrintUsage prints usage.
func PrintUsage() {
	logger.Printf(`Available Commands:
  current-context      Display the current-context
  delete-cluster NAME  Delete the specified cluster NAME from the kubeconfig
  delete-context NAME  Delete the specified context NAME from the kubeconfig
  delete-user NAME     Delete the specified user NAME from the kubeconfig
  get-clusters         Display clusters defined in the kubeconfig
  get-contexts         Display contexts defined in the kubeconfig
  get-users            Display users defined in the kubeconfig

Usage:
  kubecli config [-path] SUBCOMMAND [options]
  -path for %v

Use "kubecli {-h|--help}" for more information.`, usagePath)
}

func deleteCluster(clusters []string) {
	data, err := ioutil.ReadFile(kubeConfigPath)
	if err != nil {
		logger.Fatalln(err)
	}

	k := &types.KubeConfig{}
	err = k.Parse(data)
	if err != nil {
		logger.Fatalln(err)
	}

	for _, c := range clusters {
		r := regexp.MustCompile(c)
		for _, kubeC := range k.Clusters {
			if r.MatchString(kubeC.Name) {
				cmd := exec.Command(kubectl, append(kubectlCmd, "clusters."+kubeC.Name)...)
				out, err := cmd.CombinedOutput()
				if err != nil {
					logger.Fatalln(err, string(out))
				}
			}
		}
	}
}

func deleteContext(contexts []string) {
	data, err := ioutil.ReadFile(kubeConfigPath)
	if err != nil {
		logger.Fatalln(err)
	}

	k := &types.KubeConfig{}
	err = k.Parse(data)
	if err != nil {
		logger.Fatalln(err)
	}

	for _, c := range contexts {
		r := regexp.MustCompile(c)
		for _, kubeC := range k.Contexts {
			if r.MatchString(kubeC.Name) {
				cmd := exec.Command(kubectl, append(kubectlCmd, "contexts."+kubeC.Name)...)
				out, err := cmd.CombinedOutput()
				if err != nil {
					logger.Fatalln(err, string(out))
				}
			}
		}
	}
}

func deleteUser(users []string) {
	data, err := ioutil.ReadFile(kubeConfigPath)
	if err != nil {
		logger.Fatalln(err)
	}

	k := &types.KubeConfig{}
	err = k.Parse(data)
	if err != nil {
		logger.Fatalln(err)
	}

	for _, c := range users {
		r := regexp.MustCompile(c)
		for _, kubeC := range k.Users {
			if r.MatchString(kubeC.Name) {
				cmd := exec.Command(kubectl, append(kubectlCmd, "users."+kubeC.Name)...)
				out, err := cmd.CombinedOutput()
				if err != nil {
					logger.Fatalln(err, string(out))
				}
			}
		}
	}
}

func currentContext() {
	data, err := ioutil.ReadFile(kubeConfigPath)
	if err != nil {
		logger.Fatalln(err)
	}

	k := &types.KubeConfig{}
	err = k.Parse(data)
	if err != nil {
		logger.Fatalln(err)
	}
	logger.Println(k.CurrentContext)
}

func getUsers() {
	data, err := ioutil.ReadFile(kubeConfigPath)
	if err != nil {
		logger.Fatalln(err)
	}

	k := &types.KubeConfig{}
	err = k.Parse(data)
	if err != nil {
		logger.Fatalln(err)
	}
	for _, kubeC := range k.Users {
		logger.Println(kubeC.Name)
	}
}

func getContexts() {
	data, err := ioutil.ReadFile(kubeConfigPath)
	if err != nil {
		logger.Fatalln(err)
	}

	k := &types.KubeConfig{}
	err = k.Parse(data)
	if err != nil {
		logger.Fatalln(err)
	}
	for _, kubeC := range k.Contexts {
		logger.Println(kubeC.Name)
	}
}

func getClusters() {
	data, err := ioutil.ReadFile(kubeConfigPath)
	if err != nil {
		logger.Fatalln(err)
	}

	k := &types.KubeConfig{}
	err = k.Parse(data)
	if err != nil {
		logger.Fatalln(err)
	}
	for _, kubeC := range k.Clusters {
		logger.Println(kubeC.Name)
	}
}

// getEnv lookup key in env, or return fallback if not found.
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func userHomeDir() string {
	usr, err := user.Current()
	if err != nil {
		logger.Fatalln(err)
	}
	return usr.HomeDir
}
