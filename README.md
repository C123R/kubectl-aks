# kubectl-aks
This is a kubectl plugin to manage Azure Kubernetes Service.

### Usage

```sh
$ kubectl aks list
NAME		VERSION		RESOURCE GROUP
foo-AKS		1.11.5		 fooRG
bar-AKS		1.11.5		 barRG

$ kubectl aks get -n foo-AKS
Merged "foo-AKS" as current context in /Users/*****/.kube/config

```
