package log

import "flag"

// If non-empty, overrides the choice of directory in which to write logs.
// See createLogDirs for the full list of possible destinations.
var logDir = flag.String("log_dir", "/work/logs", "If non-empty, write log files in this directory")


// MaxSize is the maximum size of a log file in bytes.
var MaxSize uint64 = 1024 * 1024 * 1800
