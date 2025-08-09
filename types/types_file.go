package types

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type XFile struct {
	path string
}

// File 创建 XFile 实例
func File(path string) XFile {
	return XFile{path: path}
}

// Path 获取文件路径
func (f XFile) Path() string {
	return f.path
}

// Exists 判断文件是否存在
func (f XFile) Exists() bool {
	_, err := os.Stat(f.path)
	return !os.IsNotExist(err)
}

// IsFile 判断是否为文件
func (f XFile) IsFile() bool {
	info, err := os.Stat(f.path)
	return err == nil && !info.IsDir()
}

// IsDir 判断是否为目录
func (f XFile) IsDir() bool {
	info, err := os.Stat(f.path)
	return err == nil && info.IsDir()
}

// Size 获取文件大小（字节）
func (f XFile) Size() int64 {
	info, err := os.Stat(f.path)
	if err != nil {
		return 0
	}
	return info.Size()
}

// ModTime 获取修改时间
func (f XFile) ModTime() time.Time {
	info, err := os.Stat(f.path)
	if err != nil {
		return time.Time{}
	}
	return info.ModTime()
}

// Permission 获取文件权限
func (f XFile) Permission() os.FileMode {
	info, err := os.Stat(f.path)
	if err != nil {
		return 0
	}
	return info.Mode()
}

// Name 获取文件名（包含扩展名）
func (f XFile) Name() string {
	return filepath.Base(f.path)
}

// BaseName 获取文件名（不包含扩展名）
func (f XFile) BaseName() string {
	name := f.Name()
	ext := f.Ext()
	if ext != "" {
		return strings.TrimSuffix(name, ext)
	}
	return name
}

// Ext 获取文件扩展名
func (f XFile) Ext() string {
	return filepath.Ext(f.path)
}

// Dir 获取文件所在目录
func (f XFile) Dir() string {
	return filepath.Dir(f.path)
}

// DirFile 获取目录的 XFile 实例
func (f XFile) DirFile() XFile {
	return File(f.Dir())
}

// AbsPath 获取绝对路径
func (f XFile) AbsPath() (string, error) {
	return filepath.Abs(f.path)
}

// Read 读取文件内容为字节数组
func (f XFile) Read() ([]byte, error) {
	return os.ReadFile(f.path)
}

// ReadString 读取文件内容为字符串
func (f XFile) ReadString() (string, error) {
	data, err := f.Read()
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// ReadLines 读取文件内容为行数组
func (f XFile) ReadLines() ([]string, error) {
	file, err := os.Open(f.path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines, scanner.Err()
}

// ReadJSON 读取 JSON 文件并解析
func (f XFile) ReadJSON(v interface{}) error {
	data, err := f.Read()
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}

// Write 写入字节数据到文件
func (f XFile) Write(data []byte) error {
	return f.WriteWithPerm(data, 0644)
}

// WriteWithPerm 以指定权限写入字节数据到文件
func (f XFile) WriteWithPerm(data []byte, perm os.FileMode) error {
	// 确保目录存在
	if err := f.DirFile().MkdirAll(); err != nil {
		return err
	}
	return os.WriteFile(f.path, data, perm)
}

// WriteString 写入字符串到文件
func (f XFile) WriteString(content string) error {
	return f.Write([]byte(content))
}

// WriteLines 写入字符串数组到文件（每行一个元素）
func (f XFile) WriteLines(lines []string) error {
	content := strings.Join(lines, "\n")
	return f.WriteString(content)
}

// WriteJSON 将对象序列化为 JSON 写入文件
func (f XFile) WriteJSON(v interface{}) error {
	return f.WriteJSONWithIndent(v, "  ")
}

// WriteJSONWithIndent 将对象序列化为带缩进的 JSON 写入文件
func (f XFile) WriteJSONWithIndent(v interface{}, indent string) error {
	var data []byte
	var err error

	if indent == "" {
		data, err = json.Marshal(v)
	} else {
		data, err = json.MarshalIndent(v, "", indent)
	}

	if err != nil {
		return err
	}
	return f.Write(data)
}

// Append 追加字节数据到文件
func (f XFile) Append(data []byte) error {
	file, err := os.OpenFile(f.path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(data)
	return err
}

// AppendString 追加字符串到文件
func (f XFile) AppendString(content string) error {
	return f.Append([]byte(content))
}

// AppendLine 追加一行到文件
func (f XFile) AppendLine(line string) error {
	return f.AppendString(line + "\n")
}

// Copy 复制文件到目标路径
func (f XFile) Copy(dst string) error {
	return f.CopyWithPerm(dst, 0644)
}

// CopyWithPerm 以指定权限复制文件到目标路径
func (f XFile) CopyWithPerm(dst string, perm os.FileMode) error {
	src, err := os.Open(f.path)
	if err != nil {
		return err
	}
	defer src.Close()

	// 确保目标目录存在
	dstFile := File(dst)
	if err := dstFile.DirFile().MkdirAll(); err != nil {
		return err
	}

	dest, err := os.OpenFile(dst, os.O_RDWR|os.O_CREATE|os.O_TRUNC, perm)
	if err != nil {
		return err
	}
	defer dest.Close()

	_, err = io.Copy(dest, src)
	return err
}

// Move 移动文件到目标路径
func (f XFile) Move(dst string) error {
	// 先尝试重命名（同一文件系统）
	err := os.Rename(f.path, dst)
	if err == nil {
		f.path = dst
		return nil
	}

	// 如果重命名失败，则复制后删除
	if err := f.Copy(dst); err != nil {
		return err
	}

	if err := f.Delete(); err != nil {
		return err
	}

	f.path = dst
	return nil
}

// Delete 删除文件
func (f XFile) Delete() error {
	return os.Remove(f.path)
}

// DeleteAll 递归删除文件或目录
func (f XFile) DeleteAll() error {
	return os.RemoveAll(f.path)
}

// Mkdir 创建目录
func (f XFile) Mkdir() error {
	return f.MkdirWithPerm(0755)
}

// MkdirWithPerm 以指定权限创建目录
func (f XFile) MkdirWithPerm(perm os.FileMode) error {
	return os.Mkdir(f.path, perm)
}

// MkdirAll 递归创建目录
func (f XFile) MkdirAll() error {
	return f.MkdirAllWithPerm(0755)
}

// MkdirAllWithPerm 以指定权限递归创建目录
func (f XFile) MkdirAllWithPerm(perm os.FileMode) error {
	return os.MkdirAll(f.path, perm)
}

// List 列出目录下的文件和子目录
func (f XFile) List() ([]XFile, error) {
	entries, err := os.ReadDir(f.path)
	if err != nil {
		return nil, err
	}

	files := make([]XFile, len(entries))
	for i, entry := range entries {
		files[i] = File(filepath.Join(f.path, entry.Name()))
	}

	return files, nil
}

// ListFiles 列出目录下的文件（不包括子目录）
func (f XFile) ListFiles() ([]XFile, error) {
	all, err := f.List()
	if err != nil {
		return nil, err
	}

	var files []XFile
	for _, file := range all {
		if file.IsFile() {
			files = append(files, file)
		}
	}

	return files, nil
}

// ListDirs 列出目录下的子目录（不包括文件）
func (f XFile) ListDirs() ([]XFile, error) {
	all, err := f.List()
	if err != nil {
		return nil, err
	}

	var dirs []XFile
	for _, file := range all {
		if file.IsDir() {
			dirs = append(dirs, file)
		}
	}

	return dirs, nil
}

// Walk 递归遍历目录
func (f XFile) Walk(fn func(XFile) error) error {
	return filepath.Walk(f.path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		return fn(File(path))
	})
}

// Find 在目录中查找匹配的文件
func (f XFile) Find(pattern string) ([]XFile, error) {
	var matches []XFile

	err := f.Walk(func(file XFile) error {
		matched, err := filepath.Match(pattern, file.Name())
		if err != nil {
			return err
		}
		if matched {
			matches = append(matches, file)
		}
		return nil
	})

	return matches, err
}

// FindFiles 在目录中查找匹配的文件（不包括目录）
func (f XFile) FindFiles(pattern string) ([]XFile, error) {
	all, err := f.Find(pattern)
	if err != nil {
		return nil, err
	}

	var files []XFile
	for _, file := range all {
		if file.IsFile() {
			files = append(files, file)
		}
	}

	return files, nil
}

// Join 连接路径
func (f XFile) Join(elem ...string) XFile {
	elements := append([]string{f.path}, elem...)
	return File(filepath.Join(elements...))
}

// Chmod 修改文件权限
func (f XFile) Chmod(mode os.FileMode) error {
	return os.Chmod(f.path, mode)
}

// Chown 修改文件所有者（仅Unix系统）
func (f XFile) Chown(uid, gid int) error {
	return os.Chown(f.path, uid, gid)
}

// Touch 创建空文件或更新修改时间
func (f XFile) Touch() error {
	if !f.Exists() {
		return f.WriteString("")
	}

	now := time.Now()
	return os.Chtimes(f.path, now, now)
}

// IsEmpty 判断文件是否为空
func (f XFile) IsEmpty() bool {
	return f.Size() == 0
}

// IsHidden 判断是否为隐藏文件
func (f XFile) IsHidden() bool {
	name := f.Name()
	return len(name) > 0 && name[0] == '.'
}

// MimeType 获取文件 MIME 类型（简单实现）
func (f XFile) MimeType() string {
	ext := strings.ToLower(f.Ext())
	mimeTypes := map[string]string{
		".txt":  "text/plain",
		".html": "text/html",
		".css":  "text/css",
		".js":   "application/javascript",
		".json": "application/json",
		".xml":  "application/xml",
		".pdf":  "application/pdf",
		".zip":  "application/zip",
		".jpg":  "image/jpeg",
		".jpeg": "image/jpeg",
		".png":  "image/png",
		".gif":  "image/gif",
		".svg":  "image/svg+xml",
		".mp4":  "video/mp4",
		".mp3":  "audio/mpeg",
		".wav":  "audio/wav",
	}

	if mime, exists := mimeTypes[ext]; exists {
		return mime
	}
	return "application/octet-stream"
}

// ReadStream 创建读取流
func (f XFile) ReadStream() (*os.File, error) {
	return os.Open(f.path)
}

// WriteStream 创建写入流
func (f XFile) WriteStream() (*os.File, error) {
	// 确保目录存在
	if err := f.DirFile().MkdirAll(); err != nil {
		return nil, err
	}
	return os.Create(f.path)
}

// AppendStream 创建追加流
func (f XFile) AppendStream() (*os.File, error) {
	// 确保目录存在
	if err := f.DirFile().MkdirAll(); err != nil {
		return nil, err
	}
	return os.OpenFile(f.path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
}

// String 实现 Stringer 接口
func (f XFile) String() string {
	return f.path
}

// Equal 判断两个文件路径是否相等
func (f XFile) Equal(other XFile) bool {
	abs1, err1 := f.AbsPath()
	abs2, err2 := other.AbsPath()

	if err1 != nil || err2 != nil {
		return f.path == other.path
	}

	return abs1 == abs2
}

// Backup 备份文件（添加时间戳后缀）
func (f XFile) Backup() error {
	timestamp := time.Now().Format("20060102_150405")
	backupPath := fmt.Sprintf("%s.backup_%s", f.path, timestamp)
	return f.Copy(backupPath)
}

// BackupTo 备份文件到指定路径
func (f XFile) BackupTo(backupPath string) error {
	return f.Copy(backupPath)
}

// TempFile 在系统临时目录创建临时文件
func TempFile(pattern string) (XFile, error) {
	file, err := os.CreateTemp("", pattern)
	if err != nil {
		return XFile{}, err
	}
	defer file.Close()

	return File(file.Name()), nil
}

// TempDir 在系统临时目录创建临时目录
func TempDir(pattern string) (XFile, error) {
	dir, err := os.MkdirTemp("", pattern)
	if err != nil {
		return XFile{}, err
	}

	return File(dir), nil
}
