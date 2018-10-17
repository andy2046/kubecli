# kubecli

`kubectl config` provides no subcommand to list user records or batch delete context / cluster records.

`kubecli config` provides subcommands to batch delete and list user / context / cluster records.

```sh
kubecli config current-context      # Display the current-context
kubecli config delete-cluster NAME  # Delete the specified cluster NAME from the kubeconfig
kubecli config delete-context NAME  # Delete the specified context NAME from the kubeconfig
kubecli config delete-user NAME     # Delete the specified user NAME from the kubeconfig
kubecli config get-clusters         # Display clusters defined in the kubeconfig
kubecli config get-contexts         # Display contexts defined in the kubeconfig
kubecli config get-users            # Display users defined in the kubeconfig
```
