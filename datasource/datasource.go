package datasource

type DatasourceI interface {
	// connect()
	insert()
	insertBatch()
	update()
	delete()
	close()
}
