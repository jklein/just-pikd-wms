// Copyright G2G Market Inc, 2015

/*
Package daos contains data access objects which interact with
data storage layers such as databases and caches.

Its primary method of data interchange is with models from the "models" package.
Methods that write should accept models as parameters, and methods that read should return
models as well.

The daos should generally only be invoked from controllers, and especially not from models since
that would create circular dependencies.
*/
package daos
