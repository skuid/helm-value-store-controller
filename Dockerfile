FROM alpine

MAINTAINER Micah Hausler, <micah.hausler@skuid.com>

RUN apk -U add ca-certificates

COPY helm-value-store-controller /bin/helm-value-store-controller

ENV AWS_REGION us-west-2

EXPOSE 3000

ENTRYPOINT ["/bin/helm-value-store-controller"]
