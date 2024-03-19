package main

import (
	"fmt"
	"os"
	"testing"
	"time"
)

// go test -run=TestFooer

func TestFooer(t *testing.T) {
	result := Fooer(10)
	if result != "foo" {
		t.Errorf("Result was incorrect, got: %s, want: %s", result, "foo")
	}
	// t.Errorf 不会让test执行停止，所以这行会打印
	fmt.Println("continue to run this")
	// t.Fatal会中断测试进程
	t.Fatal("t.Fatal...")
	fmt.Println("do not continue to run this")
	// t.Error不会中断测试进程
	t.Error("This won't be executed")

	t.Helper()

	/*
		TestMain..
		continue to run this
		--- FAIL: TestFooer (0.00s)
		    mian_test.go:15: Result was incorrect, got: 10, want: foo
		    mian_test.go:20: t.Fatal...
		TestFooer2-------
		FAIL
		exit status 1
		FAIL    testing_demo    1.038s

	*/

}

func TestFooer2(t *testing.T) {
	fmt.Println("TestFooer2-------")
}

//func TestFooerB(t *testing.B) {
//	fmt.Println("TestFooer2-------")
//}

// M is a type passed to a TestMain function to run the actual tests
// 在所有test执行之前都会执行这个 TestMain
func TestMain(t *testing.M) {
	fmt.Println("TestMain..")
	os.Exit(t.Run())
}

// 并发执行
func TestParallel(t *testing.T) {
	t.Run("P1", func(t *testing.T) {
		t.Parallel() //
		fmt.Println("[TestParallel] - P1 - ", time.Time{}.UnixNano())
	})
	t.Run("P2", func(t *testing.T) {
		t.Parallel()
		fmt.Println("[TestParallel] - P2 - ", time.Time{}.UnixNano())
	})
}

// skip() 方法允许你区分 单元测试和集成测试
// t.Cleanup() 会执行cleanup方法，无所谓放在代码的什么位置，都会在最后执行，类似于 defer func()
func TestSkip(t *testing.T) {

	t.Cleanup(cleanup)

	// go test -run=TestSkip -v -test.short
	//用这个标志跳过，是将整个方法全部返回了，下面的代码也不会执行了，相当于一个return
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	fmt.Println("[Not Short] not skip")
}

func cleanup() {
	fmt.Println("[cleanup]...")
}

func TestTmpDir(t *testing.T) {
	tempDir := t.TempDir()
	fmt.Println("[TestTmpDir] tempDir=", tempDir)

	// [TestTmpDir] tempDir= /var/folders/x_/6nsxqsk919x_jln07lf91s240000gn/T/TestTmpDir3394366779/001
}
