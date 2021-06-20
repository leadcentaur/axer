# Axer 🪓 
A small DNS AXFR testing tool.
* This tool is my attempt at learning GoLang/concurrency.
* Inspired by a pentesting class I took

<p align="center">
  <img src="https://github.com/leadcentaur/axer/blob/5c1bd7dc3e19d2f5d84bb8f1d605a420df58bc70/banner.png">
</p>

## Basic Usage
You pass in a list of domains as input. The tool will then attempt to perform DNS zone transfer's for each domain in the list.
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

