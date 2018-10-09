FROM ubuntu:18.04

ADD run /bin/run

ENTRYPOINT ["/bin/run"]
CMD ["-l"]