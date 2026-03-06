# 🚀 exercism_watcher

![Exercism Watcher](watcher.png?raw=true "Exercism Watcher")

Watcher of Exercisms

This is a collection of helper utilities for exercism

## 📦 Build

```
cmake .
make
sudo make install
```

## Overview of `exercism_watcher.go`

`exercism_watcher.go` is a small utility written in Go that watches a directory tree for any file changes and automatically runs the Go test suite for the affected package. It is useful when working on Exercism exercises, where you often want immediate feedback on whether your solution passes the tests.

### Key Components

1. **File System Watcher** – Uses the [fsnotify](https://github.com/fsnotify/fsnotify) package to monitor file system events. The watcher is created with `fsnotify.NewWatcher()` and recursively adds every sub‑directory (excluding hidden ones) to the watch list.
2. **Recursive Directory Registration** – `filepath.WalkDir` walks the current directory (`"."`). For each directory that does **not** start with a dot, it calls `watcher.Add(path)` to start watching that directory.
3. **Test Runner** – The `runTests(dir string)` helper runs `go test ./...` inside the directory where the change occurred. It captures combined stdout/stderr, logs the output, and distinguishes between success and failure.
4. **Event Loop** – A goroutine reads from `watcher.Events` and `watcher.Errors`. For each event, it extracts the directory (`filepath.Dir(event.Name)`) and invokes `runTests` for that directory, providing immediate feedback.
5. **Blocking Main** – The program blocks forever with `<-make(chan struct{})` so the watcher continues running until the process is terminated.

### How It Works

- When the program starts, it sets up the watcher and registers all visible directories.
- As soon as a file is created, modified, renamed, or removed, `fsnotify` emits an event.
- The event handler determines the directory of the changed file and runs `go test` there. If the tests pass, a success log is printed; otherwise the failure and test output are logged.

### Typical Usage

```bash
# Run the watcher in the root of an Exercism Go exercise repository
exercism_watcher
```

The watcher will keep the terminal output updated with the results of each test run, allowing you to iterate quickly on your solution.

---

Feel free to modify the watcher (e.g., filter events, change the test command, or integrate with other languages) to suit your workflow.
