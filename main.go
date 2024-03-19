package main

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
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

	err := fmt.Errorf("fmt Errorf: %w", errors.New("test error fmt"))
	fmt.Println(err, "\n")

	// json json.Unmarshal 反序列化时接收的是 []byte 字符数组
	// json json.NewDecoder 从流中进行解码
	jsonResp1 := new(JsonLogResponse)
	json.NewDecoder(strings.NewReader(jsonStr)).Decode(jsonResp1)
	fmt.Println("json.NewDecoder -> jsonResp1=", jsonResp1)

	// defer : defer是go中一种延迟调用机制，defer后面的函数只有在当前函数执行完毕后才能执行，将延迟的语句按defer的逆序进行执行，
	//也就是说先被defer的语句最后被执行，最后被defer的语句，最先被执行，通常用于释放资源
	// 多个defer多话，从后往前执行
	// 先输出“echo return”后输出“defer echo1”，意味着先return后执行defer函数
	echo()
	// defer 最大的功能是 panic 后依然有效,所以defer可以保证你的一些资源一定会被关闭，从而避免一些异常出现的问题
	//deferPanic()

	// 对于指针数组是要求数组里面都是int类似变量的指针，用string类型的不行，它不仅要存指针还要对变量类型做校验的
	//var ptr []*int
	//i := "str"
	//append(ptr, &i)

	// interface
	//在 golang中，只要interface中任意一个方法，被interface继承或者被struct的func实现，都是这个interface的实现类
	// golang 完全是通过函数名称是不是一样来决定是不是一个接口的实现

	//关于值传递和引用传递
	// 不用指针的话就是值传递，一般在方法上，或者函数输入输出上会看到加*号的对象
	// 指针就是类似引用传递. 一般要修改对象的话就让传指针，否则就值传递就行

	// errors.New(strconv.Itoa(status))
	// strconv 主要用于字符串和基本数据类型之间的相互转换
	str12 := "123"
	atoi, err := strconv.Atoi(str12) // 字符串转int
	fmt.Println(atoi)
	atoi += 1
	str122 := strconv.Itoa(atoi) // int转字符串
	fmt.Println(str122)

	var bytel = []byte("1")
	strconv.AppendInt(bytel, 3, 10)
	fmt.Println(bytel)
	// strconv 中 base的解释。
	// 根据base，基数的不同输出的不同。base 的选择（例如，二进制，十进制、十六进制等）
	var strconvD int64 = 32
	strconvDS2 := strconv.FormatInt(strconvD, 2)
	strconvDS10 := strconv.FormatInt(strconvD, 10)
	strconvDS16 := strconv.FormatInt(strconvD, 16)
	fmt.Println("int64转string-2进制：", strconvDS2)  //1100100
	fmt.Println("int64转string-10进制", strconvDS10) //100
	fmt.Println("int64转string-16进制", strconvDS16) //20

	strconvDS2 = "10"
	//把字符串转成int，base表明字符串是多少进制的，bitSize说明你转成这个int的范围
	//the result must fit into. Bit sizes 0, 8, 16, 32, and 64 correspond to int, int8, int16, int32, and int64
	parseInt, _ := strconv.ParseInt(strconvDS2, 2, 0)                             // 说明我要给你转换的字符串是一个2进制的
	fmt.Printf("\nstring转int,2进制. 原字符串=%v, parseInt后=%v\n", strconvDS2, parseInt) //用2进制表示的10，则值为2
	strconvDS2 = "17"                                                             //这个值根据下面base，下面base是8，则是8进制，则这个地方每个位上的数字不能大于8，大于8则会返回0
	parseInt1, _ := strconv.ParseInt(strconvDS2, 8, 32)
	fmt.Printf("\nstring转int,8进制. 原字符串=%v, parseInt后=%v\n", strconvDS2, parseInt1) //用8进制表示的8

	// map
	map1 := map[string]string{
		"name": "rick",
		"age":  "18",
	}
	s, ok := map1["name"]
	fmt.Println(ok)
	if ok {
		fmt.Println("map s = ", s)
	}

	// 创建固定长度的数组.
	x1 := [2]int{}
	fmt.Println("x1=", x1)
	x2 := make([]int, 2)
	fmt.Println("x2=", x2)

	//数组切片 [:] 切片后要重新赋值回原数组
	x21 := []int{1, 2, 3, 4, 5}
	fmt.Printf("\nx21=%v", x21)
	x21 = x21[:3] // [1 2 3]
	fmt.Printf("\n从开始取到第三位 x21[:3]=%v", x21)
	x21 = []int{1, 2, 3, 4, 5}
	x21 = x21[1:] // [2 3 4 5]
	fmt.Printf("\n从第一个开始往后取 x21[1:]=%v", x21)
	x21 = []int{1, 2, 3, 4, 5}
	x21 = x21[2:4] // [3 4]
	fmt.Printf("\n从第二个往后选，选到第四位 x21[2:4]=%v", x21)

	// range带 index.
	// append 往上面 x1和x2插入时，因为长度是固定的，而且有默认值，所以没有地方去插入了，所以不能插入
	x3 := []int{1, 2, 3, 4} // append不能为空
	x3 = append(x3, 5)
	for i, _ := range x1 {
		fmt.Printf("\nx1 range: index=%v, value=%v", i, x3[i])
	}
	fmt.Println("\nx3=", x3)
	// 利用 append 去删除元素
	// 原理：拿到要删除的值的前面的分片，然后用这个值后面到末尾的一个分片去append进去，就把要删除的那个位置覆盖掉，从而达到删除的目的
	x3 = append(x3[:2], x3[3:]...)
	fmt.Println("利用 append delete 后x3=", x3)

	//对象类型 - 只针对 any 或者 interface{} 类型的变量，才能获取type
	var objType anyType = 123
	switch t := objType.(type) {
	case string:
		fmt.Println("objType is string.t=", t)
	case int:
		fmt.Println("objType is int.t=", t)
	}
	// 在 Go 中，类型断言使用 x.(T) 的语法进行，其中 x 是接口类型的变量，也可以是any类型的，T 是具体的类型
	if _, ok := objType.(int); ok {
		fmt.Println("objType is int... by objType.(int)")
	}

	// 别的包里的属性，用 type 去创建一个属性，不同于 struct 中的属性，这个是go文件里面的属性
	values := Values{
		"key": []string{"a", "b"},
	}
	fmt.Println("别的包中定义的属性. values=", values)

}

// 用来验证 switch t := objType.(type)
type anyType any

// interface
type Animal interface {
	engine
	name() string
	run()
	eat()
}
type Dog struct{}

func (d Dog) run() {
	fmt.Println("dog run")
}

// 接口用作依赖注入
type Database interface {
	Connect() error
	Disconnect() error
	Query(string) ([]byte, error)
}
type User struct{}
type UserRepository struct {
	db Database
}

func (r UserRepository) GetUser(id int) (*User, error) {
	data, err := r.db.Query(fmt.Sprintf("SELECT * FROM users WHERE id = %d", id))
	if err != nil {
		return nil, err
	}
	// parse data and return User object
	fmt.Println(data)
	return nil, nil
}

// defer
func echo() string {
	defer echo1()
	defer echo2()
	defer echo3()
	return EchoReturn()
}
func EchoReturn() string {
	fmt.Println("echo return")
	return "echo return"
}
func echo1() string {
	fmt.Println("defer echo1")
	return "echo1"
}
func echo2() {
	fmt.Println("echo2")
}
func echo3() {
	fmt.Println("echo3")
}
func deferPanic() {
	defer echo1()
	panic("panic ....")
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

// struct 对象的 “继承”，其实在golang中这个是 组合
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
