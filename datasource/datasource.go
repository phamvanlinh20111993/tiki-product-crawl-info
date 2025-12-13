package datasource

type datasourceI interface {
	// connect()
	insert()
	insertBatch()
	update()
	delete()
	close()
}
