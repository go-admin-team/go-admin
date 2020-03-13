package test

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestMapWithFuncValue(t *testing.T) {
	m := map[int]func(op int) int{}
	m[1] = func(op int) int { return op }
	m[2] = func(op int) int { return op * op }
	m[3] = func(op int) int { return op * op * op }
	t.Log(m[1](2), m[2](2), m[3](2))
}

func TestMapForSet(t *testing.T) {
	mySet := map[int]bool{}
	mySet[1] = true
	n := 1
	if mySet[n] {
		t.Logf("%d is existing", n)
	} else {
		t.Logf("%d is not existing", n)
	}
	mySet[3] = true
	t.Log(len(mySet))
	delete(mySet, 1)
	if mySet[n] {
		t.Logf("%d is existing", n)
	} else {
		t.Logf("%d is not existing", n)
	}
}

func TestStringToRune(t *testing.T) {
	s := "中华人民共和国"
	for _, c := range s {
		t.Logf("%[1]c %[1]d", c)
	}
}

func TestStringFn(t *testing.T) {
	s := "A,B,C"
	parts := strings.Split(s, ",")
	for _, part := range parts {
		t.Log(part)
	}

	t.Log(strings.Join(parts, "-"))
}

func TestConv(t *testing.T) {
	s := strconv.Itoa(10)
	t.Log("str" + s)
	if i, err := strconv.Atoi("10"); err == nil {
		t.Log(10 + i)
	}
}

func returnMultiValues() (int, int) {
	return rand.Intn(10), rand.Intn(20)
}

func timeSpent(inner func(op int) int) func(op int) int {
	return func(n int) int {
		start := time.Now()
		ret := inner(n)
		fmt.Print("time spent:", time.Since(start).Seconds())
		return ret
	}
}

func slowF(op int) int {
	time.Sleep(time.Second * 1)
	return op
}

func TestFn(t *testing.T) {
	a, _ := returnMultiValues()
	t.Log(a)
	tsSf := timeSpent(slowF)
	t.Log(tsSf(10))
}

type Programmer interface {
	WriteHelloWorld() string
}

type JavaProgrammer struct {
}

func (g *JavaProgrammer) WriteHelloWorld() string {
	return "System.IO.Input.Write(\"Hello World!\")"
}

type GoProgrammer struct {
}

func (g *GoProgrammer) WriteHelloWorld() string {
	return "fmt.Print(\"Hello World!\")"
}

func TestClient(t *testing.T) {
	var p Programmer
	var j Programmer
	p = new(GoProgrammer)
	t.Log(p.WriteHelloWorld())
	j = new(JavaProgrammer)
	t.Log(j.WriteHelloWorld())
}

var LessThanTwoError = errors.New("n should be not less than 2")
var LargerThenHuandredError = errors.New("n should be not larger than 100")

func GetFibonacci(n int) ([]int, error) {
	if n < 2 {
		return nil, LessThanTwoError
	}
	if n > 100 {
		return nil, LargerThenHuandredError
	}

	fibList := []int{1, 1}
	for i := 2; i < n; i++ {
		fibList = append(fibList, fibList[i-2]+fibList[i-1])
	}
	return fibList, nil
}

func TestGetFibonacci(t *testing.T) {
	if v, err := GetFibonacci(5); err != nil {
		if err == LessThanTwoError {
			t.Log("It is less.")
		}
		t.Error(err)
	} else {
		t.Log(v)
	}
}

func service() string {
	time.Sleep(time.Millisecond * 50)
	return "Done" + time.Now().String()
}

func otherTask() {
	fmt.Println("working on something else--" + time.Now().String())
	time.Sleep(time.Millisecond * 100)
	fmt.Println("Task is done---------------" + time.Now().String())
}

func TestService(t *testing.T) {
	fmt.Println(service())
	otherTask()
}

func AsyncService() chan string {
	retCh := make(chan string, 1)

	go func() {
		ret := service()
		fmt.Println("returned result.--------" + time.Now().String())
		retCh <- ret
		fmt.Println("service exited.---------" + time.Now().String())
	}()
	return retCh
}

func TestAsynService(t *testing.T) {
	retCh := AsyncService()
	otherTask()
	fmt.Println(<-retCh)
	time.Sleep(time.Second * 1)
}
