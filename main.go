package main

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"reflect"
	"time"
)

func main() {

	fmt.Println("hello world!")

	name := "rick"
	fmt.Println(name)
	fmt.Println(&name)

	//调用这个函数，传入一个匿名函数，匿名函数没有方法名，所以直接是 func(x int, y int)
	funcWithFuncParams(func(x int, y int) int {
		return x + y
	})

	//传入一个有名称的函数
	funcWithFuncParams(multi)

	/**
	plugin是一个package
	ServeOpts是上面package下的一个struct，这个Serve方法需要传入 *ServeOpts struct的*
	这个地方是创建这个ServeOpts，ServeOpts中有一个属性 ProviderFunc，这个属性是一个函数：type ProviderFunc func() *schema.Provider
	所以需要传入一个函数，而func() *schema.Provider 就是创建了一个匿名函数，这个函数由自己的provider提供
	这个schema.Provider包含了所有的这个provider的信息
	*/
	//plugin.Serve(&plugin.ServeOpts{
	//	ProviderFunc: func() *schema.Provider {
	//		return demo_provider.Provider()
	//	},
	//})

	// go中的 * 和 &
	var varA int
	var pointVarA *int

	varA = 1
	fmt.Println("varA原值", varA)
	fmt.Println("varA指针值", &varA)
	pointVarA = &varA
	fmt.Println("pointVarA原指针值", pointVarA)
	fmt.Println("pointVarA原指针值使用*后值", *pointVarA)

	// 通过函数返回值修改值
	var name1 string = "rick"
	name1 = addPrefixStr(name1)
	fmt.Println("利用函数返回值修改对象值的方法：name=", name1)

	//不通过函数返回值，通过传入函数指针来修改。
	// 减少内存开销： 有时候在函数中需要传递一个大对象（如大型结构体）给另一个函数处理，直接传递对象值会拷贝整个对象，造成内存开销较大。使用指针传递对象的地址可以避免拷贝整个对象，减少内存开销
	addPrefixStrByPoint(&name1)
	fmt.Println("利用函数返回值修改对象值的方法：name=", name1)

	// channel
	var c = make(chan int)
	go process(c)
	for i := range c {
		fmt.Println(i)
	}

	// strings 中的方法

	// 读文件 os.Open(caCert) io.Reader

	// 范型
	var m1 MyMap[int, string] = map[int]string{
		1: "2",
	}
	fmt.Println("m1=", m1)
	//结构体范型
	var car Car[string] = Car[string]{
		Data: "data.",
	}
	fmt.Println("car=", car)
	// channel generic
	var ch MyChannel[string] = make(MyChannel[string])
	fmt.Println("ch=", ch)
	//
	qu := Queue[string]{
		elements: []string{"one", "two"},
	}
	qu.Put("three")
	fmt.Println("qu=:", qu)

	//interface{} 就等于 any
	//var i1 interface{} = 123
	var i1 any = "123"
	switch i1.(type) {
	case int:
		fmt.Println("i1 int")
	case string:
		fmt.Println("i1 str")
	default:
		fmt.Println("i1 default")
	}
	// 通过反射获取类型
	r1 := reflect.ValueOf(i1)
	fmt.Printf("reflect r1=%v, kind=%v\n", r1, r1.Kind())

	//内存地址的地址
	str11 := "str"
	str11B := &str11
	fmt.Println("str11B = ", str11B)
	str11BB := &str11B
	fmt.Println("str11BB = ", str11BB)
	fmt.Println("str11BB1 = ", *str11BB)

	// go中，在对象中的函数叫方法，否则都叫函数，方法是绑定对象的

	//struct default
	structDefaultValue := People{}.Init()
	structDefaultValue1 := new(People).Init()
	fmt.Println("structDefaultValue=", structDefaultValue)
	fmt.Println("structDefaultValue=", structDefaultValue1)

	// context
	ctx := context.Background()
	anotherCtx := context.WithValue(ctx, "key", "value")
	showVFromCtx(anotherCtx, "key")

	// struct没有变量名，只有变量类型
	ot := OnlyType{
		structDefaultValue1,
	}
	ot.People.Init()
	fmt.Println("\not=", ot.People.firstName)

	//go对struct和interface的继承如下，将struct和interface不用加变量，直接写在struct或者interface中

	// 新建对象 - new 方式
	p1 := new(People).Init()
	fmt.Println("p1 new people: ", p1)
	//直接初始化
	p2 := People{
		"xie1", "zhao1",
	}
	fmt.Println("p2 new people: ", p2)
	//指定属性
	p3 := People{
		firstName: "xie3",
		lastName:  "zhao3",
	}
	fmt.Println("p3 new people: ", p3)

	// struct中的tag
	// func Unmarshal(data []byte, v any) error {
	xmlLogResp := XmlLogResponse{}
	xmlResp := `
<?xml version="1.0" encoding="UTF-8" ?>
<Context>
	<pageNum>1</pageNum>
	<pageSize>10</pageSize>
	<code>9999</code>
	<message>this is a message</message>
	<detail>this is the detai</detail>
	<cost>121</cost>
</Context>
`
	//因为Unmarshal方法是没有返回的，所以要把内存地址传进去让修改
	fmt.Println("Unmarshal 前 xmlLogResp=", xmlLogResp)
	if err := xml.Unmarshal([]byte(xmlResp), &xmlLogResp); err != nil {
		fmt.Errorf("could not parse job XML: %w", err)
	}
	fmt.Printf("\nUnmarshal 后 jsonLogResp: \npageNum=%v,\npageSize=%v,\ncode=%v,\nmessage=%v,\nOther=%v,\nDescription=%v",
		xmlLogResp.PageNum,
		xmlLogResp.PageSize,
		xmlLogResp.Code,
		xmlLogResp.Message,
		xmlLogResp.Other,
		xmlLogResp.Description,
	)

	//生成xml
	xmlLogResp = XmlLogResponse{
		PageSize: "123",
		PageNum:  "12",
	}
	xmlByte, _ := xml.Marshal(&xmlLogResp)
	fmt.Println("\nxmlByte = \n", string(xmlByte))

	// 反序列化json
	jsonStr := `
		{
			"code": "9999",
			"message": "this is a message"
		}
		`
	jsonResp := new(JsonLogResponse)
	json.Unmarshal([]byte(jsonStr), &jsonResp)
	fmt.Printf("\njsonResp=\n code=%v \n message=%v", jsonResp.Code, jsonResp.Message)

	// 序列化
	jsonObj := JsonLogResponse{
		Code:    "12",
		Message: "this is 12",
	}

	bjson, _ := json.Marshal(jsonObj)
	fmt.Println("\njsonObj序列化后为：", string(bjson))

	// try catch

	//
	err := fmt.Errorf("fmt Errorf: %w", errors.New("test error fmt"))
	fmt.Println(err)
}

// 在struct中定义序列化和反序列化字段
// 后面的内容是struct tag，标签，是用来辅助反射的 `json:"code"`
type JsonLogResponse struct {
	Code    any    `json:"code"`
	Message string `json:"message"`
}

type XmlLogResponse struct {
	XMLName     xml.Name `xml:"Context"`
	Code        string   `xml:"code"`    //将元素名写入该字段
	Message     string   `xml:"message"` //将Message该属性的值写入该字段
	PageNum     string   `xml:"pageNum"`
	PageSize    string   `xml:"pageSize"`
	Other       []string `xml:",any"`      //any指的是剩余的，这个要用一个字符数组去接收
	Description string   `xml:",innerxml"` //Unmarshal函数直接将对应原始XML文本写入该字段
}

// struct 中引用 struct，直接创创建struct，而不是引入已经创建好的，适用于不会被其他struct引用，也就是不那么通用
type LogRecord struct {
	logPath string
	logMeta struct {
		CreateTime time.Time
		FileSize   float32
	}
}

// 对接口的继承
type engine interface {
	run()
}
type Bicycle interface {
	engine
	drive()
}

// struct 对象的 继承
// People 是 OnlyType的一个属性，也可以直接调用people中的字段
type OnlyType struct {
	People
}

// context
func showVFromCtx(ctx context.Context, key any) {
	fmt.Printf("ctx key=%v, value=%v", key, ctx.Value(key))
}

// struct default values
type People struct {
	firstName string
	lastName  string
}

func (p People) Init() People {
	p.lastName = "xie"
	p.firstName = "zhao"
	return p
}

// 自定义范型的数据结构
type MyMap[KEY int | string, VALUE int | string] map[KEY]VALUE
type MyChannel[T int | string] chan T

// 自定义struct中的结构
type Car[T int | string] struct {
	Data T
}

type Queue[T any] struct {
	elements []T
}

func (q *Queue[T]) Put(value T) {
	q.elements = append(q.elements, value)
}

func process(c chan int) {
	//defer close(c)
	for i := 0; i < 5; i++ {
		c <- i
	}
	// fatal error: all goroutines are asleep - deadlock!
	// 不关闭的话，主线程等待打印，go线程等待写入，所以需要把 chan 关掉
	// 也可以用上面的 defer close(c) ，告诉这个函数执行完 close channels
	close(c)
}

func addPrefixStrByPoint(name *string) {
	fmt.Println("name *string = ", name)
	//*name 则解引用name的值
	//把处理后的值再复制给name的指针
	*name = "pre_" + *name
}

// defer
func addPrefixStr(name string) string {
	defer fmt.Println("addPrefixStr done.. defer..")
	return "pre_" + name
}

// 这个函数的参数是传入一个函数，这个函数传入两个int变量，返回一个int值
func funcWithFuncParams(fn func(int, int) int) {
	result := fn(3, 5)
	fmt.Println(result)
}

func multi(x int, y int) int {
	return x * y
}
