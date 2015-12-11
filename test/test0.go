package sometest

type Foo struct {
	Public    int
	nonPublic int
}

func helloWorld() {

}

func (f *Foo) fooPointerMember() {

}

func (f Foo) fooValueMember() {

}
