# FROM python:3.11.1
FROM golang:bullseye
RUN apt-get update && apt-get upgrade -y
RUN apt-get install curl snap git python3 python3-pip tree -y
RUN pip3 install -U pip
RUN pip3 install --upgrade pip