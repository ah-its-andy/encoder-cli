package mkvtoolnix

import (
	"bufio"
	"fmt"
	"io/fs"
	"log"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/ah-its-andy/encoder-cli/dto"
	"github.com/ah-its-andy/encoder-cli/utils"
	"github.com/ah-its-andy/goconf"
)

func ExtractMKV(fileName string, outputDir string) ([]dto.TrackInfo, error) {
	execPath, err := GetExecutableFilePath(utils.ExecutableFile("mkvextract"))
	if err != nil {
		return nil, err
	}

	// tempDir, err := utils.TempDir(fileName)
	// if err != nil {
	// 	return nil, err
	// }
	//tracks input.mkv --all-attachments -d /path/to/output

	cmd := utils.Command(execPath, "tracks", fileName, "--all-attachments", "-d", outputDir)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create StdoutPipe: %v", err)
	}
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start command: %v", err)
	}

	scanner := bufio.NewScanner(stdout)
	trackInfos := make([]dto.TrackInfo, 0)

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "Track ID ") {
			parts := strings.Split(line, ":")
			if len(parts) > 1 {
				id, _ := strconv.Atoi(strings.TrimSpace(parts[0][9:]))
				typeStr := strings.TrimSpace(parts[1])
				codecID := ""
				filePath := ""

				if strings.HasPrefix(typeStr, "video") {
					codecID = strings.TrimSpace(strings.TrimPrefix(typeStr, "video ("))
				} else if strings.HasPrefix(typeStr, "audio") {
					codecID = strings.TrimSpace(strings.TrimPrefix(typeStr, "audio ("))
				} else if strings.HasPrefix(typeStr, "subtitles") {
					codecID = strings.TrimSpace(strings.TrimPrefix(typeStr, "subtitles ("))
				}

				if strings.Contains(line, ": '") && strings.HasSuffix(line, "'") {
					filePath = strings.TrimSuffix(strings.TrimPrefix(line[strings.Index(line, ": '")+3:], "'"), "'")
				}

				trackInfo := dto.TrackInfo{
					TrackID:  id,
					CodecID:  codecID,
					FilePath: filePath,
				}

				trackInfos = append(trackInfos, trackInfo)
			}
		}
	}

	if err := cmd.Wait(); err != nil {
		return nil, fmt.Errorf("command finished with error: %v", err)
	}

	simpleFileName := fileName[len(filepath.Dir(fileName)) : len(fileName)-len(filepath.Ext(fileName))]
	err = filepath.Walk(filepath.Dir(fileName), func(path string, info fs.FileInfo, err error) error {
		fn := info.Name()
		fn = fn[len(filepath.Dir(fn)):]
		if strings.HasPrefix(fn, simpleFileName) {
			log.Printf("copying %s", fn)
			if err := utils.Cp(info.Name(), filepath.Join(outputDir, fn)); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return trackInfos, err
	}

	return trackInfos, nil
}

func GetExecutableFilePath(execFileName ...string) (string, error) {
	execPath, ok := goconf.GetString("tools.mkvtoolnix")
	if !ok {
		return "", fmt.Errorf("mkvtoolnix not found")
	}
	args := make([]string, len(execFileName)+1)
	args[0] = execPath
	copy(args[1:], execFileName)
	ret := filepath.Join(args...)
	return ret, nil
}
