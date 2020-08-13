FROM makeblock/alpine

ARG tmp_api_version=default_value
ENV API_VERSION=$tmp_api_version

ADD ./build /go-gen/

WORKDIR /go-gen

EXPOSE 8080

CMD ["./main"]