FROM ubuntu

RUN apt-get update && apt-get install ffmpeg mkvtoolnix -y && apt-get clean 

ADD encoder-cli /etc/encoder-cli/encoder-cli

ADD conf /etc/encoder-cli/conf

CMD ["/etc/encoder-cli/encoder-cli", "-c", "/etc/encoder-cli/task.yaml"]