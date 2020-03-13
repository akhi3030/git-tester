package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

func testBindToPort(executable *Executable, logger *customLogger) error {
	logger.Debugf("Running git init")
	tempDir, err := ioutil.TempDir("", "worktree")
	if err != nil {
		return err
	}

	executable.WorkingDir = tempDir
	_, err = executable.Run("init")
	if err != nil {
		return err
	}

	return fmt.Errorf("Something is wrong")

	for _, dir := range []string{".git", ".git/objects", ".git/refs"} {
		if err = assertDirExistsInDir(tempDir, dir); err != nil {
			logDebugTree(logger, tempDir)
			return err
		}
	}

	for _, file := range []string{".git/HEAD"} {
		if err = assertFileExistsInDir(tempDir, file); err != nil {
			logDebugTree(logger, tempDir)
			return err
		}
	}

	if err = assertFileContents(".git/HEAD", path.Join(tempDir, ".git/HEAD")); err != nil {
		return err
	}

	return nil
}

func assertFileContents(friendlyName string, path string) error {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	actualContents := string(bytes)
	expectedContents := "ref: refs/heads/master\n"
	if actualContents != expectedContents {
		return fmt.Errorf("Expected %s to contain '%s', got '%s'", friendlyName, expectedContents, actualContents)
	}

	return nil
}

func assertDirExistsInDir(parent string, child string) error {
	info, err := os.Stat(path.Join(parent, child))
	if os.IsNotExist(err) {
		return fmt.Errorf("Expected the '%s' directory to be created", child)
	}

	if !info.IsDir() {
		return fmt.Errorf("Expected '%s' to be a directory", child)
	}

	return nil
}

func assertFileExistsInDir(parent string, child string) error {
	info, err := os.Stat(path.Join(parent, child))
	if os.IsNotExist(err) {
		return fmt.Errorf("Expected the '%s' file to be created", child)
	}

	if info.IsDir() {
		return fmt.Errorf("Expected '%s' to be a file", child)
	}

	return nil
}

func logDebugTree(logger *customLogger, dir string) {
	logger.Debugf("Files found in directory: ")
	doLogDebugTree(logger, dir, " ")
	logger.Debugf("")
}

func doLogDebugTree(logger *customLogger, dir string, prefix string) {
	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	for _, info := range entries {
		if info.IsDir() {
			logger.Debugf(prefix + "- " + info.Name() + "/")
			doLogDebugTree(logger, path.Join(dir, info.Name()), prefix+" ")
		} else {
			logger.Debugf(prefix + "- " + info.Name())
		}
	}
	// filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
	// 	if info.IsDir() {
	// 		logger.Debugf(path)
	// 		// doLogDebugTree(logger, path, prefix+"  -")
	// 	} else {
	// 		logger.Debugf(path)
	// 	}

	// 	return nil
	// })
}