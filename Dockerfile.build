FROM golang
ADD sap/*.go /sap/
ADD sap/Makefile /sap/
ADD shared/*.go /shared/
WORKDIR /sap
RUN make buildgo
CMD ["/bin/bash"]
