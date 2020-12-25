# 인터페이스

인터페이스 타입은 다른 타입의 동작을 일반화하거나 추상화해서 표현한다.

인터페이스의 일반화를 통해 함수를 특정 구현의 세부 사항에 구애받지 않고 더 유연하고 융통성 있게  작성할 수 있다.

## 7.1 인터페이스 규약

### 구상(concrete) vs 추상(abstract)

구상 타입이라는 말이 어색하게 느껴질 수 있다. 이는 Golang에서만 등장하는 단어는 아니고 객체 지향 관련 서적으로 자주 나오는 단어이다.

`concrete class` 구상 클래스란 모든 operation의 구현을 제공하는 클래스이다.

`abstract class` 추상 클래스란 abstract operation 을 포함하고 있는 클래스이고, 구현은 제공하지 않고 시그니쳐만 제공한다.

`추상`과 쌍을 이루는 단어가 `구상`이고 `구체적인`과 같은 단어로 생각해보면 이해가 된다.

### 연습문제 7.2

````go
type ByteCounter int

func (c *ByteCounter) Write(p []byte) (int, error) {
  *c += ByteCounter(len(p))
  return len(p), nil
}

func CountingWriter(w io.Writer) (io.Writer, *int64) {
    c := &WrapperWriter{w, 0}
    return c, &c.written
}
...
    var c ByteCounter
    c4, n := CountingWriter(&c) // O
    c4, n := CountingWriter(c) // X
````

위와 같이 코드가 작성되었을때 CountingWriter 함수의 인자로 io.Writer를 줄때 ByteCounter 타입을 그대로 넘겨주면 컴파일 에러가 발생한다.

ByteCounter가 io.Writer를 구현하고 있긴하지만 pointer receiver 형태로 구현하고 있어서 ByteCounter 타입 그대로 넘겨주면 구현하고 있지 않은걸로 인식된다.

따라서 pointer receiver 형태로, 주소값을 넘겨주어야 한다.

### 7.2 인터페이스 타입

기존 인터페이스의 조합으로 이뤄진 새 인터페이스 선언을 할 수 있다. 구조체 내장과 유사한 문법으로 인터페이스 내장(embedding) 문법도 가능하다.

```go
type ReadWriter interface {
    Reader
    Writer
}

// 혼합해서 사용 가능
type ReadWriter interface {
    Read(p []byte) (n int, err error)
    Writer
}
```


