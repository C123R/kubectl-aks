[![Build Status](https://travis-ci.org/C123R/kubectl-aks.svg?branch=master)](https://travis-ci.org/C123R/kubectl-aks)

# kubectl-aks

This is a kubectl plugin to manage Azure Kubernetes Service. `kubectl-aks` support following operation:

- List AKS cluster from the current Azure Subscrption.
- Get the upgrade versions available for a managed Kubernetes cluster.
- Get available upgardes for the cluster.
- Upgrade a managed Kubernetes cluster to a newer version.
- Upgrade a managed Kubernetes cluster to a newer version.

In order to authenticate against Azure API we need Azure Service Principal. To create a service principal, you can use Azure CLI as shown below:

Note: Make sure you have stored azure.auth file in kubectl home directory as mentioned above or export ENV variable `AZURE_AUTH_LOCATION` with the path of the azure.auth file.

```bash
$ az ad sp create-for-rbac —sdk-auth > ~/.kube/azure.auth

$ cat ~/.kube/azure.auth
{
  "clientId": "XXXXXXXXXXXXXXXXXXXXXXXXXXX",
  "clientSecret": "XXXXXXXXXXXXXXXXXXXXXXXXXXX",
  "subscriptionId": "XXXXXXXXXXXXXXXXXXXXXXXXXXX",
  "tenantId": "XXXXXXXXXXXXXXXXXXXXXXXXXXX",
  "activeDirectoryEndpointUrl": "https://login.microsoftonline.com",
  "resourceManagerEndpointUrl": "https://management.azure.com/"
}
```

You can also create Azure Service Principal using [Azure Portal](https://docs.microsoft.com/en-us/azure-stack/user/azure-stack-create-service-principals#create-service-principal).

## Installation

To use this kubectl-aks plugin, you can follow the official Kubernetes Plugin [documentation](https://kubernetes.io/docs/tasks/extend-kubectl/kubectl-plugins/#using-a-plugin).

- Make it executable:

```bash
chmod u+x kubectl-aks
```

- Place it in your PATH:

```bash
 mv kubectl-aks /usr/local/bin
```

- Now it can be access using `kubectl` command:

```bash
$ kubectl aks
The AKS plugin to manage Azure Kubernetes Service.

Usage:
  aks [command]

Available Commands:
  get          Get Kubernetes cluster configuration.
  get-upgrades Get the upgrade versions available for a managed Kubernetes cluster.
  help         Help about any command
  list         List AKS cluster from the current Azure Subscrption.
  scale        Scale the node pool in a managed Kubernetes cluster.
  upgrade      Upgrade a managed Kubernetes cluster to a newer version.

Flags:
  -h, --help   help for aks

Use "aks [command] --help" for more information about a command.
```

## Usage

- Get list of AKS Clusters from current subscription:

```bash
$ kubectl aks list
NAME           VERSION     NODES   RESOURCE GROUP
foo-AKS        1.13.1        4         fooRG
bar-AKS        1.12.7        5         barRG
```

- Get Kubernetes Credentials for specific cluster and merge with `~/.kube/config` (Default). Credentials can be saved to specific path:

```bash
$ kubectl aks get foo-AKS
Merged "foo-AKS" as current context in /Users/*****/.kube/config

$ kubectl aks get foo-AKS -p /home/foo/config
Merged "foo-AKS" as current context in /home/foo/config
```

- Get list of available upgardes for specific cluster:

``` bash
$ kubectl aks get-upgrades bar-AKS

Current Version: 1.12.7

List of avaliable upgrades for bar-AKS:
1.13.5
```

- Scale up the nodes of AKS cluster:

``` bash
$ kubectl aks scale bar-AKS -c 6
Do you want to continue with this operation? [y|n]: y
Scaling bar-AKS to 6 nodes .. (spinner)
```

- Upgrade kubernetes version of AKS cluster:

``` bash
$ kubectl aks upgrade bar-AKS -v 1.13.5
Kubernetes may be unavailable during cluster upgrades.
Do you want to continue with this operation? [y|n]: y
Upgrading bar-AKS to Kubernetes version 1.13.5 ..(spinner)
```
