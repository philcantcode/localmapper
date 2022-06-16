# OS Setup
FROM ubuntu:latest
RUN apt-get update
RUN mkdir /src

# Basic tools
RUN apt-get install nano -y

# Network tools
RUN apt-get install net-tools -y
RUN apt-get install iputils-ping -y
RUN apt-get install ca-certificates -y

# Golang
RUN apt-get install golang -y

# Tools
RUN apt-get install nmap -y