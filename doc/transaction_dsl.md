## Replicant transaction DSL
Replicant's transaction DSL is as follows:

| attribute    | required  | type          | description  |   |
|--------------|-----------|---------------|--------------|---|
| name         | true      | string        | Name of the transaction.  |   |
| driver       | true      | string        | Driver to be used to execute the script. Can be web (FerretQL), javascript, or go  |   |
| schedule     | true      | string        | When to schedule the job (TODO possible values)  |   |
| timeout      | true      | string        | Value until the request times out (TODO possible values)  |   |
| retry_count  | true      | int           | Number of retris for the request.  |   |
| inputs       | true      | list(object)  | Variables to be used in the template.  |   |
| metadata     | false     | list(object)  | Metadata to be added to the transaction. |   |
| script       | true      | string        | Script to be run depending on the driver.  |   |


