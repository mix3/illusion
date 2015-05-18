FROM ubuntu

MAINTAINER mix3

RUN apt-get update
RUN apt-get -y install wget jq curl unzip

WORKDIR /opt

RUN export ILLUSION_VERSION=`curl https://api.github.com/repos/mix3/illusion/releases | jq -r ".[0].tag_name"` && wget https://github.com/mix3/illusion/releases/download/${ILLUSION_VERSION}/illusion-${ILLUSION_VERSION}-linux-amd64.zip -O /opt/illusion.zip
RUN unzip illusion.zip
WORKDIR /opt/illusion

ADD config.toml /opt/illusion/config.toml

CMD ./illusion
