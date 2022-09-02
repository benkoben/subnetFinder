Subnetcalc
===

A simple program that calculcates the next possible subnets inside an existing subnet range.

Why I created this
===

In some cases when landingzones in the cloud start growing. I created this tool in order to calculate subnets prefixes within an existing virtual network in order to make infrastrucutre deployments fully automated. Think "project vending machine" and "network provisioning". 

Testing
===

Unit testing is done by editing `subnetcalc_tests.go`. Tests are comprised of:
1. A Case variable - Describes the input and expected output
2. A test function - Calls the function or method thats tested with Case variable as argument.

How to use?
==

```
az network vnet show -n hub-vnet-weeu-dev-001 -g connectivity-rg-weeu-dev-001 -o json | go run main.go -new-subnets '[{"aks":24}, {"dbxPriv": 28}, {"dbsPub": 22}]'

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

