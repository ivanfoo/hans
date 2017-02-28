# hAns server image
FROM ubuntu:16.04

# Install packages
RUN apt-get -y update && apt-get -y install python-minimal ansible

# Add go binary
ADD build/* /hans/

# Add dasboard files
ADD dashboard/templates /hans/dashboard/templates
ADD dashboard/public /hans/dashboard/public

# Ports
EXPOSE 8080

# Start hAns
WORKDIR /hans
ENTRYPOINT ["/hans/hans"]
