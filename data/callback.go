package data

type ProgressCallback func(current, total int, message string)

type undeterminedFileProgressCallback func(filePath, fileName string)
