## 0.1.5 (July 5, 2017

 * storage: User must pass in Storage URL to CRUD resources [GH-74]

## 0.1.4 (June 30, 2017)

 * opc: Fix infinite loop around auth token exceeding it's 25 minute duration. [GH-73]

## 0.1.3 (June 30, 2017)

  * opc: Add additional logs instance logs [GH-72]
  
  * opc: Increase instance creation and deletion timeout [GH-72]

## 0.1.2 (June 30, 2017)


FEATURES:

  * opc: Add image snapshots [GH-67]
  
  * storage: Storage containers have been added [GH-70]


IMPROVEMENTS: 
  
  * opc: Refactored client to be generic for multiple Oracle api endpoints [GH-68]
  
  * opc: Instance creation retries when an instance enters a deleted state [GH-71]
  
## 0.1.1 (May 31, 2017)

IMPROVEMENTS:

 * opc: Add max_retries capabilities [GH-66]
 
## 0.1.0 (May 25, 2017)

BACKWARDS INCOMPATIBILITIES / NOTES:

 * Initial Release of OPC SDK
