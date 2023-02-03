FROM ubuntu
RUN apt-get -y update && apt-get install -y wget unzip curl

WORKDIR /mc

RUN wget {{.}} -O bedrock-server.zip 

RUN unzip bedrock-server.zip && rm bedrock-server.zip

ENV LD_LIBRARY_PATH=/mc

CMD ["./bedrock_server"]