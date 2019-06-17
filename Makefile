fmt:
	cd backend && make fmt
	cd frontend && make fmt

stats:
	wc `find backend/src -name \*.go |sort`
	wc `find frontend/src -name \*.js |sort`

docker:
	cd frontend && make docker-build
	
	cd ruby && make docker-build
