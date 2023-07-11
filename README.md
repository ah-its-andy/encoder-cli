sudo docker run --rm -it \
  -v /caching/transcodes/work:/caching/transcodes/work \
  -v /caching/transcodes/output:/caching/transcodes/output \
  -v /caching/transcodes/task:/etc/encoder-cli/task \
  standardcore/encoder-cli:alpha-2