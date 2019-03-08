package main

import "os"

func removeFile(fileName string){
	//Remove file at end
	if _, err := os.Stat(fileName); err == nil {
		err := os.Remove(fileName)
		if err != nil {
			onError(err)
		}
	}
}
