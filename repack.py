import os
import subprocess
import json
import shutil
import argparse

def main():
    print("Starting repackage...")
    parser = argparse.ArgumentParser(description="Command line argument parser")
    parser.add_argument("--source_dir", dest="source_dir", help="Directory containing source files")
    parser.add_argument("--output_dir", dest="output_dir", help="Directory to write output files to")
    parser.add_argument("--exts", dest="exts",  help="File extension to use (default is \"mp4|mkv\")")
    args = parser.parse_args()
    if len(args.output_dir) == 0:
        print("Please specify a directory to write output files to with the --output_dir argument")
        exit(1)
    if (args.source_dir) == 0:
        print("Please specify a directory containing source files with the --source_dir argument")
        exit(1)
    exts = args.exts
    if exts is None or len(exts) == 0:
        exts = "mp4|mkv"
    
    if os.path.exists(args.output_dir):
        print("Output directory is ok")
    else:
        os.makedirs(args.output_dir)

    run(args.source_dir, args.output_dir, exts)

def run(source_dir, output_dir, file_extensions):
    print(f"Searching video files from {source_dir}")
    
    video_files = find_video_files(source_dir, file_extensions)
    if len(video_files)==0:
        print(f"No video files found")
        exit(0)

    print(f"Found {len(video_files)} video files")
    for file in video_files:
        print(f"Processing file {file}")
        stream_data = ffprobe_show_streams_json(file)
        streams = stream_data["streams"]
        print(f"Found {len(streams)} streams")
        tempDir = create_temp_directory(file)
        print(f"Created temporary directory: {tempDir}")
        extract_streams(file, stream_data, tempDir)
        print("Streams extracted successfully")
        audio_files = find_audio_streams(stream_data)
        aac_converted = False
        flac_converted = False
        for audio_file in audio_files:
            audio_file_full_name = os.path.join(tempDir, audio_file.file)
            if audio_file.file.lower().endswith((".aac")):
                continue
            else:
                if aac_converted:
                    continue
                else:
                    aac_file = os.path.splitext(audio_file.file)[0] + ".aac"
                    aac_file_full_name = os.path.join(tempDir, aac_file)
                    if is_5_1_side(audio_file.stream):
                        convert_audio_to_aac_51(audio_file_full_name, aac_file_full_name)
                    else:
                        convert_audio_to_aac(audio_file_full_name, aac_file_full_name)
                    aac_converted = True
                    print(f"Converted audio file to aac file {aac_file}")

            if audio_file.file.lower().endswith((".flac")):
                continue
            else:
                if flac_converted:
                    continue
                else:
                    flac_file = os.path.splitext(audio_file.file)[0] + ".flac"
                    flac_file_full_name = os.path.join(tempDir, flac_file)
                    convert_audio_to_flac(audio_file_full_name, flac_file_full_name)
                    flac_converted = True
                    print(f"Converted audio file to flac file {flac_file}")

        collect_subtitle_files(file, tempDir)
        for subtitle_file in os.listdir(tempDir):
            subtitle_file_path = os.path.join(tempDir, subtitle_file)
            if os.path.isfile(subtitle_file_path) and subtitle_file.lower().endswith((".ass", ".ssa")):
                convert_subtitle_to_srt(subtitle_file_path)
                print(f"Converted subtitle to srt {subtitle_file}")
        output_file_name = os.path.basename(file)
        output_file_name_without_extension, output_file_name_extension = os.path.splitext(output_file_name)
        output_file = os.path.join(output_dir, output_file_name_without_extension+".mkv")
        print("Writing MKV output file to: " + output_file)
        create_mkv_from_temp_directory(tempDir, output_file)
        shutil.rmtree(tempDir)



def ffprobe_show_streams_json(input_file):
    command = [
        "ffprobe",
        "-v", "quiet",
        "-print_format", "json",
        "-show_format",
        "-show_streams",
        input_file
    ]
    result = subprocess.run(command, capture_output=True, text=True)
    output = result.stdout.strip()
    data = json.loads(output)
    return data

def extract_streams(input_file, streams_data, output_dir):
    for stream in streams_data["streams"]:
        stream_index = stream["index"]
        codec_name = stream["codec_name"]
        codec_type = stream["codec_type"]

        if codec_name.lower() == "subrip":
            codec_name = "srt"
            
        output_file = os.path.join(output_dir, f"{codec_type}_{stream_index}.{codec_name}")

        # 构建提取命令
        ffmpeg_command = [
            "ffmpeg",
            "-hide_banner",
            "-i", input_file,
            "-map", f"0:{stream_index}",
            "-c", "copy",
            output_file
        ]

        # 执行提取命令
        subprocess.run(ffmpeg_command)

def is_5_1_side(audio_stream):
    channels = audio_stream.get('channels')
    layout = audio_stream.get('channel_layout')
    if channels == 6 and layout == '5.1(side)':
            return True
    return False

def find_audio_streams(streams_data):    
    audio_streams = []
    
    for stream in streams_data["streams"]:
        stream_index = stream["index"]
        codec_type = stream["codec_type"]
        codec_name = stream["codec_name"]

        if codec_type == "audio":
            output_file = f"{codec_type}_{stream_index}.{codec_name}"
            audio_streams.append(AudioStream(stream, output_file))
    
    return audio_streams

def find_video_streams(streams_data):
    output_files = []
    for stream in streams_data["streams"]:
        stream_index = stream["index"]
        codec_type = stream["codec_type"]
        codec_name = stream["codec_name"]

        if codec_type == "video":
            output_file = f"{codec_type}_{stream_index}.{codec_name}"
            output_files.append(output_file)
    
    return output_files

def convert_audio_to_aac_51(input_file, output_file):
    ffmpeg_command = [
        "ffmpeg",
        "-hide_banner",
        "-i", input_file,
        "-c:a", "aac",
        "-ac", "6",
        "-strict", "-2",
        output_file
    ]
    subprocess.run(ffmpeg_command)

def convert_audio_to_aac(input_file, output_file):
    ffmpeg_command = [
        "ffmpeg",
        "-hide_banner",
        "-i", input_file,
        "-c:a", "aac",
        "-strict", "-2",
        output_file
    ]
    subprocess.run(ffmpeg_command)

def convert_audio_to_flac(input_file, output_file)  :
    ffmpeg_command = [
        "ffmpeg",
        "-hide_banner",
        "-i", input_file,
        "-c:a", "flac",
        output_file
    ]

    # 执行转码命令
    subprocess.run(ffmpeg_command)

def collect_subtitle_files(input_file, output_dir):
    source_dir = os.path.dirname(input_file)
    source_filename = os.path.splitext(os.path.basename(input_file))[0]
    for file in os.listdir(source_dir):
        file_path = os.path.join(source_dir, file)
        if os.path.isfile(file_path) and file.startswith(source_filename) and file.lower().endswith((".srt", ".ass", ".ssa")):
            output_file = os.path.join(output_dir, file)
            shutil.copy2(file_path, output_file)

def convert_subtitle_to_srt(input_file):
    if os.path.isfile(input_file) and input_file.lower().endswith((".ass", ".ssa")):
        output_file = os.path.splitext(input_file)[0] + ".srt"

        ffmpeg_command = [
            "ffmpeg",
            "-hide_banner",
            "-i", input_file,
            output_file
        ]

        subprocess.run(ffmpeg_command)

def create_temp_directory(input_file):
    source_dir = os.path.dirname(input_file)
    source_filename = os.path.splitext(os.path.basename(input_file))[0]
    temp_dir = os.path.join(source_dir, source_filename+"_temp")
    if os.path.exists(temp_dir):
        shutil.rmtree(temp_dir)

    os.makedirs(temp_dir)

    return temp_dir

def find_video_files(root_directory, file_extensions):
    video_files = []

    for root, dirs, files in os.walk(root_directory):
        for file in files:
            file_path = os.path.join(root, file)
            if file_path.find("@Recycle") != -1:
                continue
            if file_path.find("/.") != -1:
                continue
            if file.lower().endswith(tuple(file_extensions.split("|"))):
                video_files.append(file_path)

    return sorted(video_files)

def create_mkv_from_temp_directory(temp_directory, output_file):
    command = f'mkvmerge -o "{output_file}" "{temp_directory}"/*'
    # 执行命令
    subprocess.run(command, shell=True)


class AudioStream:
    def __init__(self, stream, file):
        self.stream = stream
        self.file = file

if __name__ == "__main__":
    main()