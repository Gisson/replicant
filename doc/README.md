## Replicant base notions
### Transaction
Replicant has the notion of transaction. This is an object where (by a certain configuration) makes some action which will produce a result. This result will be fed into the emitters
 to produce significate results to be later consumed by whoever wishes to consume it.

There are normally two types of transactions, scheduled transactions, and ad-hoc transactions. In the case of a scheduled transaction, Replicant will
 keep a pool of registered transactions and will schedule them appropriately in order to make them run orderly and smoothly. In the latter, Replicant
 will execute the transaction at the moment of the call, this can be a one one-shot trasaction (meaning it will only exist for the moment of the run),
 or this can be a transaction that already exists but will be forced to run at a certain moment.

To create a transaction we need to create a yaml (it can be JSON but we will use all examples as yaml) following a strict DSL described in the [transaction DSL doc](transaction_dsl.md).

### Driver
Replicant has the notion of driver. Driver is an object which will be dinamically created and will be the way replicant will use to fetch metrics (make the requests and get results).

### Emitter
Repliocant has the notion of emitter. Emitter(s) is the target of the results. It is where the Replicant results will be published. For example publishing the results as
 prometheus metrics.

### Data Store
TODO

## Replicant architecture

Replicant has the following main components:
 - Webserver (mandatory);
 - Manager (mandatory);
 - Drivers (mandatory at least one);
 - Emitters (mandatory at least one);
 - Data store (mandatory);
 - Callbacks (optional).


## Webserver
Responsible for receiving connections from the clients and handle the correct endpoints forwarding them wherever they belong


## Manager
The core part of Replicant. This component handles the scheduled transactions registered by the client. It will also orchestrate
 everything to fetch and place/publish the results as it is fit.

## Drivers
To put it simple, it is which type of scripts does Replicant allow.

## Emitters
Objects which will be used to publish the results.

## Data store
TODO

## Callbacks
In the case of endpoints that returhn a callback, Replicant will use this callback in order to fetch the results in a real manner, not a fake one.
