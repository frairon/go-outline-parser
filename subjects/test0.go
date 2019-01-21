package subjects

type Foo struct {
	Public    int
	nonPublic int
}
type Iface interface {
	FuncA(int, string)
}

func helloWorld() {

}

func (f *Foo) fooPointerMember2(arg1 int, args2 string, xx chan int32) {

}

func (f Foo) fooValueMember() {

}
