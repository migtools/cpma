FROM golang:latest AS builder
RUN mkdir /app 
ADD . /app/ 
WORKDIR /app 
RUN make

FROM registry.access.redhat.com/ubi8/ubi
COPY --from=builder /app/bin/cpma /usr/local/bin/cpma
WORKDIR /mnt 
ENTRYPOINT ["cpma"]
LABEL RUN docker run -it --rm -v \${PWD}:/mnt:z -v \$HOME/.kube:/.kube:z -v \$HOME/.ssh:/.ssh:z -u \${UID} \${IMAGE}
