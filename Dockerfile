FROM ubuntu
RUN apt-get -y update && apt-get install -y wget unzip

WORKDIR /mc

RUN wget https://minecraft.azureedge.net/bin-linux/bedrock-server-1.18.12.01.zip -O bedrock-server.zip 

RUN unzip bedrock-server.zip && rm bedrock-server.zip

ENV LD_LIBRARY_PATH=/mc

CMD ["./bedrock_server"]