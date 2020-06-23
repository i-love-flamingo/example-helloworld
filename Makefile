CONTEXT?=dev
.PHONY: serve serve-jaeger jaeger-docker prepare-start-over

serve:
	DEBUG=1 CONTEXT=$(CONTEXT) go run main.go serve

serve-jaeger:
	DEBUG=1 CONTEXT=$(CONTEXT) go run main.go --flamingo-config 'opencensus.jaeger.enable: true' serve

jaeger-docker:
	docker run -e COLLECTOR_ZIPKIN_HTTP_PORT=9411 -p 5775:5775/udp -p 6831:6831/udp -p 6832:6832/udp -p 5778:5778 -p 16686:16686 -p 14268:14268 -p 9411:9411 jaegertracing/all-in-one:latest

prepare-start-over:
	rm -rf src/helloworld
	touch src/.gitkeep
	rm config/routes.yaml
	rm templates/index.html
	rm templates/hello.html
	sed -i '/helloworld.Module/d' main.go
	sed -i '/example-helloworld/src/helloworld/d' main.go