# hAns standalone image
FROM ubuntu:16.04

# Mongo repository
ENV MONGO_GPG_KEY 0C49F3730359A14518585931BC711F9BA15703C6
ENV MONGO_DEB_REPO deb [ arch=amd64,arm64 ] http://repo.mongodb.org/apt/ubuntu xenial/mongodb-org/3.4 multiverse

# Install required packages
RUN apt-key adv --keyserver hkp://keyserver.ubuntu.com:80 --recv ${MONGO_GPG_KEY} \
	&& echo ${MONGO_DEB_REPO} > /etc/apt/sources.list.d/mongodb-org-3.4.list \
	&& apt-get -y update \
	&& apt-get -y install git python-minimal ansible mongodb-org beanstalkd supervisor

# Clean apt cache
RUN rm -rf /var/cache/apt

# Minimal mongo setup
RUN mkdir -p /data/db && chown mongodb:mongodb /data/db

# Create directory for queue persistence
RUN mkdir -p /hans/queue

# Add go binary
ADD build/* /hans/

# Add dashboard files
ADD dashboard/templates /hans/dashboard/templates
ADD dashboard/public /hans/dashboard/public

# Add supervisor config
ADD supervisord.conf /etc/supervisor/conf.d/supervisord.conf

# Ports
EXPOSE 8080

# Start all services
WORKDIR /hans
CMD ["/usr/bin/supervisord"]
