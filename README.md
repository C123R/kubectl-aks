[![Build Status](https://travis-ci.org/C123R/kubectl-aks.svg?branch=master)](https://travis-ci.org/C123R/kubectl-aks)


# kubectl-aks
This is a kubectl plugin to manage Azure Kubernetes Service.

### Usage

```sh
$ kubectl aks list
NAME		VERSION		RESOURCE GROUP
foo-AKS		1.11.5		 fooRG
bar-AKS		1.11.5		 barRG

$ kubectl aks get foo-AKS
Merged "foo-AKS" as current context in /Users/*****/.kube/config

$ kubectl aks get foo-AKS -p /home/foo/config
Merged "foo-AKS" as current context in /home/foo/config

```
