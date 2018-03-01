Oracle SDK for Terraform
===========================================

**Note:** This SDK is _not_ meant to be a comprehensive SDK for Oracle Cloud. This is meant to be used solely with Terraform.

OPC Config
----------

To create the Oracle clients, a populated configuration struct is required.
The config struct holds the following fields:

* `Username` - (`*string`) The Username used to authenticate to Oracle Public Cloud.
* `Password` - (`*string`) The Password used to authenticate to Oracle Public Cloud.
* `IdentityDomain` - (`*string`) The identity domain for Oracle Public Cloud.
* `APIEndpoint` - (`*url.URL`) The API Endpoint provided by Oracle Public Cloud.
* `LogLevel` - (`LogLevelType`) Defaults to `opc.LogOff`, can be either `opc.LogOff` or `opc.LogDebug`.
* `Logger` - (`Logger`) Must satisfy the generic `Logger` interface. Defaults to `ioutil.Discard` for the `LogOff` loglevel, and `os.Stderr` for the `LogDebug` loglevel.
* `HTTPClient` - (`*http.Client`) Defaults to generic HTTP Client if unspecified.

Oracle Compute Client
----------------------
The Oracle Compute Client requires an OPC Config object to be populated in order to create the client.

Full example to create an OPC Compute instance:
```go
package main

import (
  "fmt"
  "net/url"
  "github.com/hashicorp/go-oracle-terraform/opc"
  "github.com/hashicorp/go-oracle-terraform/compute"
)

func main() {
  apiEndpoint, err := url.Parse("myAPIEndpoint")
  if err != nil {
    fmt.Errorf("Error parsing API Endpoint: %s", err)
  }

  config := &opc.Config{
    Username: opc.String("myusername"),
    Password: opc.String("mypassword"),
    IdentityDomain: opc.String("myidentitydomain"),
    APIEndpoint: apiEndpoint,
    LogLevel: opc.LogDebug,
    // Logger: # Leave blank to use the default logger, or provide your own
    // HTTPClient: # Leave blank to use default HTTP Client, or provider your own
  }
  // Create the Compute Client
  client, err := compute.NewComputeClient(config)
  if err != nil {
    fmt.Errorf("Error creating OPC Compute Client: %s", err)
  }
  // Create instances client
  instanceClient := client.Instances()

  // Instances Input
  input := &compute.CreateInstanceInput{
    Name: "test-instance",
    Label: "test",
    Shape: "oc3",
    ImageList: "/oracle/public/OL_7.2_UEKR4_x86_64",
    Storage: nil,
    BootOrder: nil,
    SSHKeys: []string{},
    Attributes: map[string]interface{}{},
  }

  // Create the instance
  instance, err := instanceClient.CreateInstance(input)
  if err != nil {
    fmt.Errorf("Error creating instance: %s", err)
  }
  fmt.Printf("Instance Created: %#v", instance)
}
```

Please refer to inline documentation for each resource that the compute client provides.

Running the SDK Integration Tests
-----------------------------

To authenticate with the Oracle Compute Cloud the following credentials must be set in the following environment variables:

-	`OPC_ENDPOINT` - API Endpoint URL for the Compute Classic service (e.g. https://api-z13.compute.em2.oraclecloud.com/\)
- `OPC_STORAGE_ENDPOINT` - API Endpoint URL for the Object Storage Classic service (e.g. https://mydomain.storage.oraclecloud.com/\)
- `ORACLEPAAS_DATABASE_ENDPOINT` - API Endpoint URL for the Database Cloud service (e.g. https://dbaas.oraclecloud.com/\)
- `ORACLEPAAS_JAVA_ENDPOINT` - API Endpoint URL for the Java Cloud service (e.g. https://jaas.oraclecloud.com/\)
- `ORACLEPAAS_APPLICATION_ENDPOINT` - API Endpoint URL for the Application Container Cloud service (e.g. https://apaas.us.oraclecloud.com/\)

-	`OPC_USERNAME` - Username for Oracle Public Cloud
-	`OPC_PASSWORD` - Password for Oracle Public Cloud
-	`OPC_IDENTITY_DOMAIN` - Identity Domain Name (for Traditional Accounts) or Identity Service Id (for Cloud Accounts)


The Integration tests can be ran with the following command:
```sh
$ make testacc
```

Isolating a single SDK package can be done via the `TEST` environment variable
```sh
$ make testacc TEST=./compute
```

Isolating a single test within a package can be done via the `TESTARGS` environment variable
```sh
$ make testacc TEST=./compute TESTARGS='-run=TestAccIPAssociationLifeCycle'
```

Tests are ran with logs being sent to `ioutil.Discard` by default.
Display debug logs inside of tests by setting the `ORACLE_LOG` environment variable to any value.
