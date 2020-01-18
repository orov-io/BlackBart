###########################
# Global build args
###########################
# change service and workdir with your 
ARG SERVICE_NAME
ARG workdir=/${SERVICE_NAME}
############################
# STEP 1 build executable binary
############################

FROM golang:latest

ARG workdir

WORKDIR ${workdir}

RUN ["go", "get", "github.com/pilu/fresh"] 

ENTRYPOINT ["fresh", "go", "run", "main.go"]