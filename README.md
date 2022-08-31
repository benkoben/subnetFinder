Subnetcalc
===

A simple program that calculcates the next possible subnets inside an existing subnet range.

Why I created this
===

I wanted to fully automate IaC deployments to virtual networks inside Azure. The reason why I wrote it in GOlang is because of flexibility and I wanted to learn more about GOlang concepts.

This program can be run as a standalone binary inside a CI/CD agent or hosted as an API endpoint (for managed CI/CD agents)

Testing
===

Unit testing is done by editing `subnetcalc_tests.go`. Tests are comprised of:
1. A Case variable - Describes the input and expected output
2. A test function - Calls the function or method thats tested with Case variable as argument.

How to use?
===

