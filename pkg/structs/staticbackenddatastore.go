package structs

type StaticBackendDatastoreEntry struct {
	ClientIdentifier  *ClientIdentifier  `json:"clientIdentifier"`
	DeviceInformation *DeviceInformation `json:"deviceInformation"`
}

type StaticBackendDatastore struct {
	Datastore []*StaticBackendDatastoreEntry `json:"staticBackendDatastore"`
}
