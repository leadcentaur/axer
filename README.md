# Axer 🪓 
A concurrent DNS AXFR testing tool.
* This tool is my attempt at learning GoLang/concurrency.
* Inspired by a pentesting class I took

## Basic Usage
You pass in a list of domains as input. The tool will then attempt to perform a DNS zone transfer on each domain in the list.
If sucessful, it will dump to the screen.

```shell
go run axer.go -f domains.txt
```

Or

```shell
go run axer.go domains.txt
```
## Install

```shell
go get -u github.com/leadcentaur/axer
```

