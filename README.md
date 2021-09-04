# Helmify
[![CI](https://github.com/arttor/helmify/actions/workflows/ci.yml/badge.svg)](https://github.com/arttor/helmify/actions/workflows/ci.yml)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/arttor/helmify)
![GitHub](https://img.shields.io/github/license/arttor/helmify)
![GitHub release (latest by date)](https://img.shields.io/github/v/release/arttor/helmify)
[![Go Report Card](https://goreportcard.com/badge/github.com/arttor/helmify)](https://goreportcard.com/report/github.com/arttor/helmify)
[![GoDoc](https://godoc.org/github.com/arttor/helmify?status.svg)](https://pkg.go.dev/github.com/arttor/helmify?tab=doc)
[![Maintainability](https://api.codeclimate.com/v1/badges/2ee755bb948d363207bb/maintainability)](https://codeclimate.com/github/arttor/helmify/maintainability)
[![Test Coverage](https://api.codeclimate.com/v1/badges/2ee755bb948d363207bb/test_coverage)](https://codeclimate.com/github/arttor/helmify/test_coverage)

CLI that creates [Helm](https://github.com/helm/helm) charts from kubernetes yamls.

Helmify reads a list of [supported k8s objects](#status) from stdin and converts it to a helm chart. Main [use-case](#integrate-to-your-operator-sdkkubebuilder-project) is to generate Helm charts for kubernetes operators build with
[Operator-SDK](https://github.com/operator-framework/operator-sdk) or [Kubebuilder](https://github.com/kubernetes-sigs/kubebuilder). Submit issue if some features missing for your use-case.

TODO:
- clean up
- tests
- regress
- release
- install guide

## Run
Clone repo and execute command: 

```shell
cat test_data/k8s-operator-kustomize.output | go run ./cmd/helmify mychart
```

Will generate `mychart` Helm chart form file `test_data/k8s-operator-kustomize.output` representing typical operator 
[kustomize](https://github.com/kubernetes-sigs/kustomize) output.

## Usage

Example 1: `cat my-app.yaml | helmify mychart`
- will create 'mychart' directory with Helm chart from yaml file with k8s objects.

Example 2: `awk 'FNR==1 && NR!=1  {print "---"}{print}' /<my_directory>/*.yaml | helmify mychart`
- will create 'mychart' directory with Helm chart from all yaml files in `<my_directory> `directory.

Example 3: `kustomize build <kustomize_dir> | helmify mychart`
- will create 'mychart' directory with Helm chart from kustomize output.

### Integrate to your Operator-SDK/Kubebuilder project
Tested with operator-sdk version: "v1.8.0".
1. Open `Makefile` in your operator project generated by 
   [Operator-SDK](https://github.com/operator-framework/operator-sdk) or [Kubebuilder](https://github.com/kubernetes-sigs/kubebuilder).
2. Add these lines to `Makefile`:
```makefile
HELMIFY = $(shell pwd)/bin/helmify
helmify:
	$(call go-get-tool,$(HELMIFY),github.com/arttor/helmify/cmd/helmify@v0.2.2)

helm: manifests kustomize helmify
	$(KUSTOMIZE) build config/default | $(HELMIFY)
```
3. Run `make helm` in project root. It will generate helm chart with name 'chart' in 'chart' directory.

## Available options
Helmify takes a chart name for an argument.
Usage:

```helmify [flags] CHART_NAME```  -  `CHART_NAME` is optional. Default is 'chart'.

| flag | description | sample |
| --- | --- | --- |
| -h -help | Prints help | `helmify -h`|
| -v | Enable verbose output. Prints WARN and INFO. | `helmify -v`|
| -vv | Enable very verbose output. Also prints DEBUG. | `helmify -vv`|
| -version | Print helmify version. | `helmify -version`|

## Status
Supported default operator resources:
- deployment
- service
- RBAC (serviceaccount, (cluster-)role, (cluster-)rolebinding)
- configs (configmap, secret)
- webhooks (cert, issuer, ValidatingWebhookConfiguration)

### Known issues
- Helmify defines application (operator) name as the shortest common prefix of k8s objects names. 
  It is possible because operator-sdk using operator name as prefix by default.
- Helmify will not overwrite `Chart.yaml` file if presented. Done on purpose.
- Helmify will not delete existing template file, only overwrite. So, if you delete CRD, re-run `kustomize | helmify` 
crd file will still be in templates directory. (todo: add option for this)
- Helmify overwrites templates and values files on every run. 
  This meas that all your manual changes in helm template files will be lost on the next run. 
  Use kustomize /config folder as a single source of true and make changes there.
  
## Develop
To support a new type of k8s object template:
1. Implement `helmify.Processor` interface. Place implementation in `pkg/processor`. The package contains 
examples for most k8s objects.
2. Register your processor in the `pkg/app/app.go`
3. Add relevant input sample to `test_data/kustomize.output`.

### Test
For manual testing, run program with debug output:
```shell
cat test_data/k8s-operator-kustomize.output | go run ./cmd/helmify -vv mychart
```
Then inspect logs and generated chart in `./mychart` directory.

To execute tests, run:
```shell
go test ./...
```
Beside unit-tests, project contains e2e test `pkg/app/app_e2e_test.go`.
It's a go test, which uses `test_data/*` to generate a chart in temporary directory. 
Then runs `helm lint --strict` to check if generated chart is valid.
