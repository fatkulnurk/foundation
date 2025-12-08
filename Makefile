
gen-mocks:
	mockgen --package mocks --source=cache/cache.go --destination=cache/mocks/mock_cache.go
	mockgen --package mocks --source=container/container.go --destination=container/mocks/mock_container.go
	mockgen --package mocks --source=httpclient/httpclient.go --destination=httpclient/mocks/mock_httpclient.go
	mockgen --package mocks --source=httprouter/router.go --destination=httprouter/mocks/mock_router.go
	mockgen --package mocks --source=logging/logging.go --destination=logging/mocks/mock_logging.go
	mockgen --package mocks --source=mailer/mailer.go --destination=mailer/mocks/mock_mailer.go
	mockgen --package mocks --source=queue/queue.go --destination=queue/mocks/mock_queue.go
	mockgen --package mocks --source=storage/storage.go --destination=storage/mocks/mock_storage.go
	mockgen --package mocks --source=view/view.go --destination=view/mocks/mock_view.go
