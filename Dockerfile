FROM ubuntu:22.04 AS deps

RUN apt update && apt install -y ffmpeg mkvtoolnix python3 python-is-python3 && apt clean

FROM deps AS runner

WORKDIR /data

ENV ENABLE_ACC_CONVERT 1
ENV ENABLE_FLAC_CONVERT 0
ENV REMOVE_AAC 0
ENV ENABLE_SRT_CONVERT 1
ENV LANGUAGES "eng,chi,zho"

ADD repack.py repack.py

CMD ["python", "/data/repack.py", "--source_dir", "/data/source", "--output_dir", "/data/output", "--exts", "mp4|mkv"]