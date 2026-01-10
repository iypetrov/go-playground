package main

import (
	"errors"
	"fmt"
)

// func readFile(filename string) (string, error) {
// 	file, err := os.Open(filename)
// 	if err != nil {
// 		return "", err
// 	}
//
// 	defer file.Close()
//
// 	buffer := make([]byte, 100)
// 	_, err = file.Read(buffer)
// 	if err != nil {
// 		return "", err
// 	}
//
// 	return string(buffer), nil
// }
//
// func lifo() {
// 	fmt.Println("LIFO")
// 	defer fmt.Println("First Deferred")
// 	defer fmt.Println("Second Deferred")
// 	defer fmt.Println("Third Deferred")
// }
//
// func wrongErrorHandling() {
// 	file, err := os.Open("nonexistent.txt")
// 	defer file.Close()
// 	if err != nil {
// 		fmt.Println("Error opening file:", err)
// 		return
// 	}
// }
//
// type App struct {
// 	Name    string `json:"name"`
// 	Version string `json:"version"`
// }
//
// func readJsonFileManual(filename string) (App, error) {
// 	file, err := os.Open(filename)
// 	if err != nil {
// 		return App{}, err
// 	}
//
// 	data, err := io.ReadAll(file)
// 	if err != nil {
// 		file.Close()
// 		return App{}, err
// 	}
//
// 	var app App
// 	err = json.Unmarshal(data, &app)
// 	if err != nil {
// 		file.Close()
// 		return App{}, err
// 	}
//
// 	return app, nil
// }
//
// func readJsonFile(filename string) (App, error) {
// 	file, err := os.Open(filename)
// 	if err != nil {
// 		return App{}, err
// 	}
// 	defer file.Close()
//
// 	data, err := io.ReadAll(file)
// 	if err != nil {
// 		return App{}, err
// 	}
//
// 	var app App
// 	err = json.Unmarshal(data, &app)
// 	if err != nil {
// 		return App{}, err
// 	}
//
// 	return app, nil
// }
//
// func printFileContent(filename *string) {
// 	file, _ := os.Open(*filename)
// 	defer file.Close()
//
// 	fmt.Printf("Content of file %s:\n", *filename)
//
// 	buffer := make([]byte, 100)
// 	file.Read(buffer)
// 	fmt.Println(string(buffer))
// }
//
// // func captureByValue()  {
// // 	file := "file.txt"
// // 	defer printFileContent(file)
// //
// // 	file = "app.json"
// // }
// //
// // func captureByReferenceClosure()  {
// // 	file := "file.txt"
// // 	defer func() {
// // 		printFileContent(file)
// // 	}()
// //
// // 	file = "app.json"
// // }
//
// func captureByReferencePointer()  {
// 	file := "file.txt"
// 	defer printFileContent(&file)
//
// 	file = "app.json"
// }
//
type Foo struct{
	Name string
}

func OpenFoo(name string) (Foo, error) {
	foo := Foo{Name: name}
	return foo, nil
}

func (*Foo) Close() error {
	return fmt.Errorf("Closing Foo failed")
}
//
// func NotHandlingCloseError() error {
// 	foo, err := OpenFoo()
// 	if err != nil {
// 		return err
// 	}
// 	defer foo.Close()
//
// 	fmt.Println("Using", foo.Name)
// 	return fmt.Errorf("Some error occurred while using Foo")
// }
//
// func HandlingCloseError() error {
// 	foo, err := OpenFoo()
// 	if err != nil {
// 		return err
// 	}
// 	defer func() {
// 		err = foo.Close()
// 		if err != nil {
// 			fmt.Println("Error during closing Foo:", err)
// 		}
// 	}()
//
// 	return nil
// }
//
// // func wrongReadFiles(ch <-chan string) error {
// // 	for path := range ch {
// // 		file, err := os.Open(path)
// // 		if err != nil {
// // 			return err
// // 		}
// //
// // 		defer file.Close() // Close will not be called till the end of the function
// // 	}
// // 	return nil
// // }
// //
// // func readFile(path string) {
// // 		file, err := os.Open(path)
// // 		if err != nil {
// // 			return err
// // 		}
// //
// // 		defer file.Close()
// // }
// //
// // func wrongReadFiles(ch <-chan string) error {
// // 	for path := range ch {
// // 		readFile(path)
// // 	}
// // 	return nil
// // }

type File struct {
	Name string
}

func Open(name string) (*File, error) {
	file := &File{Name: name}
	return file, nil
}

func (f *File) Close() error {
	return fmt.Errorf("Error closing file: %s", f.Name)
}

// func Solve(filename string) {
// 	file, err := OpenFile(filename)
// 	if err != nil {
// 		fmt.Println("Error opening file:", err)
// 		return
// 	}
// 	defer file.Close()
// 
// 	// Simulate file processing
// 	fmt.Println("Processing file:", file.Name)
// }
// 
// func VoidLoopErrorHandling() {
// 	fileNames := []string{"file1.txt", "file2.txt", "file3.txt"}
// 	for _, name := range fileNames {
// 		err := func() error {
// 			file, err := OpenFile(name)
// 			if err != nil {
// 				return err
// 			}
// 			defer file.Close()
// 			// Simulate file processing
// 			fmt.Println("Processing file:", file.Name)
// 			return nil
// 		}()
// 		if err != nil {
// 			fmt.Println(err)
// 		}
// 	}
// }
// 
// func correctReadFilesClosure(files []string) error {
// 	for _, file := range files {
// 		err := func() error {
// 			f, err := OpenFoo(file)
// 			if err != nil {
// 				return err
// 			}
// 
// 			defer f.Close()
// 
// 			return nil
// 		}()
// 		if err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }

func correctErrorHandling() (err error) {
	file, err := Open("nonexistent.txt")
	if err != nil {
		return err
	}
	defer func() {
		err = errors.Join(err, file.Close())
	}()

	return fmt.Errorf("Random error")
}

func foo() (err error) {
	defer func() {
        if r := recover(); r != nil {
			err = fmt.Errorf("no worries, recovered from panic: %v", r)
        }
    }()
	panic("trigger panic for demonstration")
}

func main() {
	err := foo()
	if err != nil {
		fmt.Println(err.Error())
	}
}
