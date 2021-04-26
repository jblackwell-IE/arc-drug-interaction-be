FROM alpine:3.12.1

COPY arc-drug-interaction-be /arc-drug-interaction-be

ENTRYPOINT ["/arc-drug-interaction-be"]