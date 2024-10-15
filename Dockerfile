FROM alpine

WORKDIR /app
ADD ./demoApp /app

ENTRYPOINT ["./demoApp"]
