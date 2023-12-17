FROM ubuntu:22.04

EXPOSE 9090

COPY requirements.txt /tmp

RUN apt-get update && apt-get install -y python3-pip
RUN pip install --upgrade pip
RUN pip install -r /tmp/requirements.txt

ENTRYPOINT deepspeech-server --config /opt/deepspeech/config.yaml
