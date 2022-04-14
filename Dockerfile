FROM golang:latest AS gobuilder
WORKDIR /go/src/gradio
COPY . .
RUN go get -v -d
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o gradio

FROM scratch
COPY --from=gobuilder /go/src/gradio/gradio /
EXPOSE 3000
CMD [ "/gradio" ]
