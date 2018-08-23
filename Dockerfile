FROM scratch
COPY  app ./
EXPOSE 3000
ENTRYPOINT ["./app"]
