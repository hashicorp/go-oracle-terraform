Oracle SDK for Terraform
===========================================

**Note:** This SDK is not meant to be a comprehensive SDK for Oracle Cloud. This is meant to be used solely with Terraform.

Running the SDK Integration Tests
-----------------------------

To authenticate with the Oracle Compute Cloud the following credentails must be set in the following environment variables:

-	`OPC_ENDPOINT` - Endpoint provided by Oracle Public Cloud (e.g. https://api-z13.compute.em2.oraclecloud.com/\)
-	`OPC_USERNAME` - Username for Oracle Public Cloud
-	`OPC_PASSWORD` - Password for Oracle Public Cloud
-	`OPC_IDENTITY_DOMAIN` - Identity domain for Oracle Public Cloud

```sh
$ make testacc
```
