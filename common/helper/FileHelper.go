package Server_Helper

import (
	"io"
	"mime/multipart"
	"os"
)

func FileToBytes(file *os.File) ([]byte, error) {
	// 将文件指针移动到文件的开头
	_, err := file.Seek(0, io.SeekStart)
	if err != nil {
		return nil, err
	}

	// 读取整个文件内容
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func SaveMultipartFile(mf multipart.File) (*os.File, error) {
	// 创建一个临时文件
	tempFile, err := os.CreateTemp("", "uploaded-*")
	if err != nil {
		return nil, err
	}

	// 确保关闭临时文件
	defer func() {
		if err != nil {
			tempFile.Close()
			os.Remove(tempFile.Name())
		}
	}()

	// 将 multipart.File 的内容拷贝到临时文件
	_, err = io.Copy(tempFile, mf)
	if err != nil {
		return nil, err
	}

	// 重置文件指针到文件开头
	_, err = tempFile.Seek(0, io.SeekStart)
	if err != nil {
		return nil, err
	}

	return tempFile, nil
}

func BytesToFile(data []byte) (*os.File, error) {
	tempFile, err := os.CreateTemp("", "uploaded-*")
	if err != nil {
		return nil, err
	}
	if _, err = tempFile.Write(data); err != nil {
		os.Remove(tempFile.Name())
		tempFile.Close()
		return nil, err
	}
	if _, err = tempFile.Seek(0, 0); err != nil {
		os.Remove(tempFile.Name())
		tempFile.Close()
		return nil, err
	}
	return tempFile, nil
}
