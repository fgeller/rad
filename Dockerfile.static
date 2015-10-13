FROM tianon/true
MAINTAINER Felix Geller "fgeller@gmail.com"
ADD sap /
VOLUME /packs
EXPOSE 3025
CMD ["/sap", "-packdir", "/packs", "-addr", "0.0.0.0:3025"]
