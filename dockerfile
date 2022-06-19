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

# Lanuages
RUN apt-get install golang -y
RUN apt-get install python2 -y

# Tools
RUN apt-get install nmap -y
RUN apt-get install braa -y
RUN apt-get install dmitry -y
RUN apt-get install nbtscan -y