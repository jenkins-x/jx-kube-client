FROM scratch
EXPOSE 8080
ENTRYPOINT ["/jx-kube-client"]
COPY ./bin/ /