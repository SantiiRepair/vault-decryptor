FROM python:3.11.1
RUN apt-get update && apt-get upgrade -y
RUN wget https://dot.net/v1/dotnet-install.sh -O dotnet-install.sh
RUN chmod +x ./dotnet-install.sh
RUN ./dotnet-install.sh --version latest
RUN apt-get install curl git mono-complete nuget python3 python3-pip tree -y
RUN pip3 install -U pip
RUN pip3 install --upgrade pip
RUN export DOTNET_ROOT=$HOME/.dotnet
RUN export PATH=$PATH:$HOME/.dotnet:$HOME/.dotnet/tools