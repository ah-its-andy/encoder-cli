sudo docker run --rm -it \
  -v /share/caching/transcodes/work:/caching/transcodes/work \
  -v /share/caching/transcodes/output:/caching/transcodes/output \
  -v /share/caching/transcodes/task:/etc/encoder-cli/task \
  standardcore/encoder-cli:alpha-2 /bin/bash