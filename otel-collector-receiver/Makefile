run-receiver:
	@docker run -d --name jaeger \
  		-e COLLECTOR_OTLP_ENABLED=true \
  		-p 16686:16686 \
  		-p 14317:4317 \
  		-p 14318:4318 \
  		jaegertracing/all-in-one:1.41

init-collctor:
	@ocb --config builder-config.yaml

run-collector:
	@./otelcol-dev/otelcol-dev --config config.yaml

invoke:
	@telemetrygen traces --otlp-insecure --traces 1
