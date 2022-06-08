# StaticBackend
The static backend is a backend that is filled with data in startup and further provides DHCP Packets on that basis.

The static backend is setup to read entries from a json file provided via the `YNDD_ZTP_STATIC_DATASTORE_SOURCE` environment variable.

The file can be set relative to the CWD or as an absolute Path.

An example of the file is provided in this repository "examples/staticBackend/staticBackendEntries.json"