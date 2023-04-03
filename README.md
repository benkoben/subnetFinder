[![Build master](https://github.com/RedeployAB/azsubnetfinder/actions/workflows/build.yml/badge.svg)](https://github.com/RedeployAB/azsubnetfinder/actions/workflows/build.yml)


Subnetcalc
===

A simple program that calculates the next possible subnets inside an existing subnet range. 

Features
===
* Supports multiple address spaces within a single VNET object
* No authentication necessary. The program relies on passing information from tools like az-cli, echo or shell variables
* Runs as a standalone lightweight binary file instead of a script

Why I created this
===

I created this tool in order to calculate subnets prefixes within an existing virtual network. The idea is to help CD pipelines that deploy "vending machine" infrastrucutre. It works by generating network parameters based on existing virtual networks.

How to use?
==

VNET object is read form STDIN or from the `-vnet` flag. The input string must represent a JSON structure that has the following keys:

* `addressSpace.addressPrefixes: []`
* `subnets: []`

Desired subnets are read from the `-new-subnets` flag

**Method 1:**
```
$> az network vnet show -n hub-vnet-weeu-dev-001 -g connectivity-rg-weeu-dev-001 -o json | go run main.go -new-subnets '[{"aks":24}, {"dbxPriv": 28}, {"dbsPub": 22}]'

{
  "parameters": [
    {
      "name": "aks",
      "prefix": "10.100.0.0/24"
    },
    {
      "name": "dbxPriv",
      "prefix": "10.100.1.0/28"
    },
    {
      "name": "dbsPub",
      "prefix": "10.100.4.0/22"
    }
  ]
}}
```

**Method 2:**
```
$> VNET=$(az network vnet show -n hub-vnet-weeu-dev-001 -g connectivity-rg-weeu-dev-001 -o json)
$> DESIRED_SUBNETS='[{"aks":24}, {"dbxPriv": 28}, {"dbsPub": 22}]'
$> go run main.go -new-subnets "${DESIRED_SUBNETS}" -vnet "${VNET}"

{
  "parameters": [
    {
      "name": "aks",
      "prefix": "10.100.0.0/24"
    },
    {
      "name": "dbxPriv",
      "prefix": "10.100.1.0/28"
    },
    {
      "name": "dbsPub",
      "prefix": "10.100.4.0/22"
    }
  ]
}}
```
