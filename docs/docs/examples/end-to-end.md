---
sidebar_position: 2
---

# End to end setup
Setup with a http server that takes incoming message, produces message into kafka and processes messages multiple times with failure before passing at the end.

### Enable tracing across entire binary {#enable-tracing}
While kp produces spans, it has no knowledge on where the spans need to be sent and in which format. We expect libary user to set that up.

In [this example](https://github.com/honestbank/kp/blob/c79601fc5a6ccb0543909fa34113a9e700098e2f/v2/examples/full/cmd/main.go#L32) we're setting up jaeger traces to send to tempo.

### Start all dependencies that you need {#starting-dependencies}
Here's a sample [docker-compose.yaml](https://github.com/honestbank/kp/blob/c79601fc5a6ccb0543909fa34113a9e700098e2f/v2/examples/full/docker-compose.yaml#L1) that starts the dependencies that the demo requires.

### Setup tracing such that trace-ids are exposed to the user {#setup-tracing-middleware}
Here's [tracing setup](https://github.com/honestbank/kp/blob/c79601fc5a6ccb0543909fa34113a9e700098e2f/v2/examples/full/server/handlers/tracing.go#L10-L20) that exposes traceparent header via response headers.

### Produce a message with traces {#produce-message}
Producing a message is as simple as [calling `.Produce`](https://github.com/honestbank/kp/blob/c79601fc5a6ccb0543909fa34113a9e700098e2f/v2/examples/full/server/handlers/controller.go#L46) with a request context as a first parameters and the message as second parameter.

### Configuring worker to fail multiple times {#configure-worker}
Configure worker with retries, deadletters and tracing [as shown](https://github.com/honestbank/kp/blob/c79601fc5a6ccb0543909fa34113a9e700098e2f/v2/examples/full/worker/worker.go#L21-L33)

## Execution {#execution}
Simply enter into v2/examples/full directory and simply run `docker-compose up` or `docker compose up` and wait for `docker compose ps sleep` to indicate the container has exited (it'll take about 90 seconds)

Once `docker compose logs web --follow` indicates the server is listening, open the link in the browser and submit an email. Click on the link and refresh after a few seconds for traces to arrive.
