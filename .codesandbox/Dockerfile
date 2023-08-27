FROM python:3.11.1
RUN apt-get update && apt-get upgrade -y
RUN apt-get install curl git python3 python3-pip tree -y
RUN pip3 install -U pip
RUN pip3 install --upgrade pip
RUN wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
RUN rm -rf /usr/local/go && tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
RUN export PATH=$PATH:/usr/local/go/bin