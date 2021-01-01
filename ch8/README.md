# 8장 고루틴과 채널

Go는 두가지 스타일의 동시성 프로그래밍 스타일을 지원한다.

8장은 CSP를 지원하는 고루틴과 채널을 설명한다.

9장은 좀 더 전통적인 공유 메모리 기반 멀티스레딩 관점을 일부 다룬다.

> CSP 상호 통신하는 연속된 프로세스
>
> 토니 호어가 1978년에 저술한 동시성 기반 기술에 대한 논문에 등장하는 언어
>
> 프로그램이 공유 상태가 없는 프로세스의 병렬 조합이고, 프로세스들은 채널을 통해 통신하고 동기화 한다.

### 8.1 고루틴

스레드와 고루틴의 차이는 근복적으로 질적인 것이 아닌 양적인 것이다.

프로그램이 시작된 두 유일한 고루틴은 main 함수를 호출하는 것으로 이를 메인 고루틴이라고 한다.

#### 예제. spinner

해당 예제에서는 고루틴을 이용한 텍스트 스피너를 보여주기 위해 재귀방법을 소개했지만, 피보나치 수 계산은 동적계획법으로 간단히 계산 할 수 있다.

`-\|/`  으로 문자열을 타이핑했는데, 백슬러시가 포함되어 raw string literal로 표현하고자 한 것이다.

golang에서 문자열은 2가지 종류가 있다.

* Raw String Literal: back quote로 둘러싸인 string
* Interpreted String Literal: 큰 따옴표로 둘러싸인 string

````go
// 재귀 방법
// go run main.go  7.93s user 0.34s system 93% cpu 8.884 total
func fibRecursive(x int) int {
    if x < 2 {
        return x
    }
    return fibRecursive(x-1) + fibRecursive(x-2)
}
// 동적계획법
// go run main.go  0.29s user 0.43s system 49% cpu 1.430 total
var cache [46]int
func fibCache(x int) int {
    if x < 2 {
        return x
    }

    if cache[x] > 0 {
        return cache[x]
    }

    cache[x] = fibCache(x-1) + fibCache(x-2)
    return cache[x]
}
````

### 8.2 예제: 동시 시계 서버
 
`net.Listen()`으로 네트워크 포트로 들어오는 연결을 대기할 수 있다.

`net.Dial()`으로 네트워크 연결을 요청할 수 있다.

> nc (netcat): TCP/UDP 프로토콜을 사용하여 네트워크를 간단하게 읽고 쓰는 유틸리티
>
> cat이 파일에 쓰거나 읽는 것 처럼 nc는 네트워크에 쓰거나 읽는다.

handleConn() 함수가 호출되면 메인 고루틴이 해당 함수를 처리하기 때문에 2개 이상의 연결을 처리할 수 없다.

### 8.3 예제: 동시 에코 서버

netcat2 예제로 표준입력(os.Stdin)의 입력을 네트워크(net.Conn)에 출력한 뒤, 고루틴으로 동시에 네트워크(net.Conn)에서 오는 입력을 표준출력(os.Stdout)에 출력하는 예제이다.

네트워크에 출력하게 되면 reverb1 예제에서는, 네트워크(net.Conn)에서 받은 입력 문자열을 3번에 나누어 다시 네트워크(net.Conn)에 출력해줍니다.

이때 bufio.NewScanner(r) 메소드로 Scanner 객채를 생성해서 사용하는데, Scanner는 주어진 입력에서 다음 토큰, 주로 한줄씩 입력을 처리할때 사용한다.

> 토큰을 나누는 기본 기준은 '\n'

### 8.4 채널

앞서 고루틴은 CSP로부터 온 개념이며, 프로그램간 공유 상태가 없다고 했다. 프로그램간 상태를 공유, 동기화해야할 필요가 있을때는 채널을 사용한다.

채널이란 한 고루틴에서 다른 고루틴으로 값을 보내기 위한 통신 메커니즘이다.

채널의 특징
* 생성
  * ch := make(chan int)
  * 맵과 마찬가지로 make로 생성된 데이터 구조에 대한 `참조`이다. 참조를 함수의 인자로 전달할때의 특성에 주의해야 한다.
  * 채널의 제로값은 nil
* 비교
  * `==`로 비교가능
  * 같은 채널 데이터 구조에 대한 참조의 비교는 참
 * nil과 비교 가능
* 송신/수신
  * ch <- x // 송신
  * x = <-ch // 수신
* 종료
  * close(ch)
  * 닫힌 채널로 송신하면 panic
  
#### 8.4.1 버퍼 없는 채널

버퍼 없는 채널에서의 통신은 송신과 수신 고루틴이 동기화되게 되어서 `동기 채널`이라고도 부른다.

송신이 없다면 수신 고루틴은 송신이 올때까지 대기한다.

메인 고루틴이 종료되면 백그라운드 고루틴도 같이 종료가 된다. 백그라운드 고루틴 작업을 기다리게 하기 위해선 두 고루틴을 동기화해야 한다.

채널을 통해 전송된 메시지의 값이 중요하지 않고, 이벤트 자체가 중요할 수 있다. 이런 경우에는 struct{} 타입보다는 bool, int 타입을 사용한다.

#### 8.4.2 파이프라인

3개의 고루틴이 앞선 고루틴이 채널로 생성해주는 값을 받아서 사용하게 할 수 있다.

```
Counter -> Squarer -> Printer
```

이런 파이프라인의 경우 무한대로 값을 생성하면 무한하게 출력이 이루어진다.

유한한 값만을 보내기 위해서는 close(naturals)와 같이 채널을 닫을 수 있다.

채널이 닫히면 값을 모두 소진하고 제로 값을 산출한다.

채널이 닫혔는지를 확인하기 위해서는 수신 시에 ok 요소를 추가로 받을 수 있다.

x, ok := <-naturals

채널을 반드시 닫아야 하는 것은 아니고, 가비지 컬렉터에 의해 참조할 수 없는 채널의 자원은 회수된다.

이미 닫힌 채널을 닫을때도 panic

#### 8.4.3 단방향 채널 타입

```go
func counter(out chan<- int) {}
func squarer(out chan<- int, in <-chan int) {}
func printer(in <-chan int) {}
```

Go의 타입시스템은 보내기 동작이나 받기 동작 중 한 가지만 노출하는 단방향 채널 타입을 제공한다.

chan int 타입은 묵시적으로 chan<- int 혹은 <-chan int 타입들로 변환된다. 양방향 채널은 언제든지 단방향 채널로 변환, 할당될 수 있지만 반대로는 불가능하다.

#### 8.4.4 버퍼 채널

버퍼 채널은 송수신 데이터들이 담기는 큐가 있는 것이며, 이 크기는 make로 만들 때의 용량 인자에 의해 결정된다.

버퍼 채널에서 송신 작업은 큐의 가장 마지막에 요소를 삽입하고, 수신 작업은 큐의 가장 앞쪽 요소를 제거한다.

큐가 가득 차게되면 송신은 대기하게 되고, 큐가 가득 차지 않았을때는 대기 없이 송신이 가능하다.

버퍼 채널의 용량은 cap(), 버퍼된 요소의 개수는 len()으로 확인할 수 있다.

고루틴의 버퍼 채널을 큐로 사용하고 싶은 유횩이 있을 수 있는데, 수신해주는 고루틴이 없는 경우 송신 고루틴이 데드락에 걸리는 위험이 있으니 주의해야 한다.

큐가 필요하다면 슬라이스를 사용하는 것이 좋다.

```go
func mirroredQuery() string {
    responses := make(chan string, 3)
    go func() { responses <- request("asia.gopl.io") }()
    go func() { responses <- request("europe.gopl.io") }()
    go func() { responses <- request("americas.gopl.io") }()
    return <-responses // 가장 빠른 응답 반환
}
```

버퍼되지 않은 채널을 사용했다면 1개의 응답을 반환 후에 나머지 2개의 고루틴은 수신자가 없는 채널로 응답을 송신하는 과정에서 막혔을 것이다.

이런 상황을 고루틴 유출(goroutine leak)이라고 한다.

참조되지 않은 변수와는 달리 유출된 고루틴은 가비지 컬렉터 대상이 안되므로 필요하지 않은 고루틴은 스스로 종료하게 해야 한다.

버퍼 채널의 용량의 프로그램의 정확도에 영향을 주고, 충분하지 않은 버퍼 용량을 할당하면 프로그램이 교착상태에 빠질 수 있다.

송신할 개수의 상한선을 알고 있다면 그만큼의 버퍼 채널을 만들어두기도 한다.

전체 프로그램의 특정 로직에서 부분적인 병목이 있다면 특정 로직을 수행하는 고루틴을 여러개 만듬으로써 성능을 향상 시킬 수 없다.

부분적인 병목 부분을 도와주는 고루틴들을 만들어 성능 향상을 이뤄낼 수 있다.

#### 8.5 병렬 루프

전체 크기의 이미지에서 섬네일 크기의 이미지를 생성하는 `makeThumbnails()` 함수를 통해서 6가지 일반적인 동시성 패턴을 살펴볼 것이다.

```go
func makeThumbnails(filename []string) {
    for _, f := range filenames {
        if _, err := thumbnail.ImageFile(f); err != nil {
            log.Println(err)
        }
    }
}
```
이 예제에서 각 파일들의 섬네일을 만드는 것은 처리 순서가 중요하지 않다. 서로 완전히 독립적인 하위 작업들로 구성된 이와 같은 문제는 동시성을 구현하고, 병렬처리의 양에 따라 선형적으로 확장되는 성능을 보기에 좋다.

```go
func makeThumbnails2(filename []string) {
    for _, f := range filenames {
        go thumbnail.ImageFile(f)
    }
}
```
이 버전은 작업 전체를 병렬로 수행하기 위해 go 키워드를 추가했다.

이 함수는 모든 고루틴을 하나씩 동시에 시작하지만, 완료될때까지 기다리지는 않는다.

```go
func makeThumbnails3(filename []string) {
    ch := make(chan struct{})
    for _, f := range filenames {
        go func(f string) {
            thumbnail.ImageFile(f)
            ch <- strunc{}{}
        }(f)
    }
    for range filenames {
        <-ch
    }
}
```
이 버전에서는 파일의 개수가 정해진 것을 알기에 그 수만큼 이벤트를 세는 방법으로 전체 고루틴의 완료를 기다린다.

이때 익명함수를 고루틴으로 넘길때 루프변수가 캡쳐되는 문제를 생각해볼 수 있다.

```go
for _, f := range filenames {
    go func() {
        thumbnail.ImageFile(f)
    }
}
```
이런식으로 루프변수를 캡쳐하면 루프의 반복으로 인해 모든 고루틴의 f 값은 슬라이스의 마지막값이 있게된다.

이를 막기위해 버전3에서는 명시적인 파라미터를 추가하였다.

```go
func makeThumbnails4(filename []string) error {
    errors := make(chan error)

    for _, f := range filenames {
        go func(f string) {
            _, err := thumbnail.ImageFile(f)
            errors <- err
        }(f)
    }
    for range filenames {
        if err := <-errors; err != nil {
            return err
        }
    }
    return nil
}
```
이 버전은 에러가 발생할 시에 main으로 값을 반환하는 방법을 다룬 버전이다.

여기에는 앞서얘기한 버그가 존재하는데, 에러가 발생하여 오류를 반환할 시에 수신하는 고루틴이 남지않고, 나머지 고루틴들에서는 해당 채널에 값을 보낼 수 없어 무한히 대기하게 되는 고루틴 유출이 발생한다.
 
```go
func makeThumbnails5(filename []string) (thumbfiles []string, err error) {
    type item struct {
        thumbfile string
        err     error
    }
    ch := make(chan item, len(filenames))
    for _, f := range filenames {
        go func(f string) {
            var it item
            it.thumbfile, it.err := thumbnail.ImageFile(f)
            ch <- it
        }(f)
    }
    for range filenames {
        it := <- ch
        if it.err != nil {
            return nil, it.err
        }
        thumbfiles = append(thumbfiles, it.thumbfile)
    }
    return thumbfiles, nil
}
```
이 버전은 버전4의 고루틴 유출을 막기 위해 충분한 용량의 버퍼 채널을 사용한다.

```go
func makeThumbnails6(filenames <-chan string) int64 {
    size := make(chan int64)
    var wg sync.WaitGroup

    for _, f := range filenames {
        wg.Add(1)
        //worker
        go func(f string) {
            defer wg.Done()
            thumb, err := thumbnail.ImageFile(f)
            if err != nil {
                log.Println(err)
                return
            }
            info, _ := os.Stat(thumb)
            sizes <- info.Size()
        }(f)
    }
    //closer
    go func() {
        wg.Wait()
        close(sizes)
    }()
    var total int64
    for size := range sizes {
        total += size
    }
    return total
}
```
이 버전은 새 파일이 차지하는 전체 바이트 수를 반환하는데, 다른점은 인자로 파일명 슬라이스가 아닌 문자열 채널을 받는다.

이때문에 루프 반복 횟수를 예측할 수 없다.

이런 경우 마지막 고루틴이 완료될 때를 알기 위해서는 `sync.WaitGroup`이라는 카운터 변수를 사용한다.

Add, Done은 서로 비대칭이며, Add(-1)은 Done()과 같다.

sizes 채널은 각 파일의 크기를 메인 코루틴으로 돌려주며, 메인 고루틴에서는 range 루프로 받은 파일의 크기를 합한다.

이때 closer 고루틴은 전체 작업자가 완료되면 sizes 채널을 닫아준다.

대기(range 루프)와 닫기(close 고루틴)의 순서가 바뀌면 대기(range 루프)가 종료되지 않기 때문에 순서에 주의해야 한다.

#### 8.6 예제: 동시 웹 크롤러

url을 커맨라인 인자로 받으면 해당 url에 접근하고 거기서 다음 url list를 추출한 뒤 다시 접근하는 (BFS 방식) 웹 크롤러이다.

for loop 내부에서 고루틴을 생성할때 루프변수 캡쳐시 명시적인 파라미터로 넘겨주는 것을 볼 수 있다.

```go
// #1
go func() { worklist <- os.Args[1:] }()
// #2
worklist <- os.Args[1:]
```
crawl1에서는 #1의 방법을 사용하고 있고, #2는 채널을 수신하는쪽 (for loop)가 같은 메인 고루틴에 존재해 데드락이 발생한다

데드락을 해결하기 위해서는 별도의 고루틴에서 채널에 송신하거나, 버퍼 채널을 사용해야 한다.

> crawl1 예제를 실행시에 file open이 과도하게 발생해 컴퓨터가 셧다운되는 현상이 발생했다. 절대 오래 실행하지 말것

crawl1 예제는 무제한적으로 병렬적이라 컴퓨터 리소스를 과도하게 점유해버리는 문제가 생긴다.

이를 제한하기 위해서 용량이 n인 버퍼 채널로 카운팅 세마포어를 구현할 수 있다.

채널 버퍼의 빈 슬롯 n개를 두고 각각을 토큰처럼 사용하는 것이다.

토큰을 획득하면 작업을 수행하고, 획득하지 못했으면 토큰을 획득할때까지 대기한다. 작업을 마치면 토큰을 반납한다.

````go
var tokens = make(chan struct{}, 20)
...
    tokens <- struct{}{}
    list, err := links.Extract(url)
    <-tokens
````
crawl2 예제에서는 용량 20의 버퍼 채널을 생성했고, 제한하고 싶은 IO작업(list.Extract의 net.Dial)과 최대한 가깝게 채널의 token을 획득하고 반납하는 로직을 추가한다.

도달 가능한 링크를 모두 찾은 뒤에도 프로그램이 종료되지 않는 문제는 카운터 n 변수를 활용함으로써 해결했다.

crawl3 예제에서 보여주는 동시성 제어하는 또다른 방법은, 작업을 수행하는 고루틴을 미리 필요한 만큼 생성해두는 것이다.

고루틴 내부에서 채널을 수신하고 있다가 메인 고루틴에서 채널로 송신하면 미리 대기하고 있던 고루틴들이 채널에서 수신하여 작업을 수행하는 것이다.

#### 8.7 select를 통한 다중화

`select`는 다른 언어의 `switch`와 유사하지만 case문에 채널을 받는다는게 특징이다.

case문에 있는 채널 중 하나로부터 값이 들어오면 수행하며, 보통은 for loop로 반복시킵니다. (반복시키지 않으면 1번만 실행하고 끝)

일정시간 동안 아무 채널에도 값이 들어오지 않으면 특정 로직을 수행시킬 수도 있고, `case <-time.After(10 * time.Second)`와 같이 구현할 수 있다.

countdown3 예제에서 time.Tick() 함수를 호출해서 일정 주기마다 이벤트를 수신합니다.

##### tick(), 고루틴 leak
Tick() 함수는 유용하지만 for loop 종료 후에 아무도 수신하지 않는 tick 채널로 계속해서 이벤트를 보내기 때문에 고루틴 leak이 발생할 수 있다.

전체 어플리케이션에 걸쳐 tick이 필요할때만 사용하거나, 명시적으로 tick.Stop()을 호출해 이벤트 보내는 것을 정지시키는 것이 필요하다.

##### select, case, default
select, case 문을 사용하게 되면 case 문의 채널에 값이 들어올때까지 select 문에서 블록된다.

이때 대기하지 않고 바로 어떤 작업을 수행해야 할 수도 있다. 이때 default 문을 작성해주면 대기없이 바로 어떤 작업을 수행시킬 수 있다.

주기적으로 이를 반복하면서 이벤트를 기다리는 것을 채널 풀링(channel polling)이라 한다.

##### 채널의 제로값, nil
채널의 네가지 특성을 다시 살펴보자
1. nil 채널에게 값을 넣으려하면 영구 블록 (데드락)
2. nil 채널에게서 값을 받으려하면 영구 블록 (데드락)
3. closed 채널에게 값을 넣으려하면 panic 발생
4. closed 채널에게서 값을 읽으려하면 zero 리턴. 버퍼 채널이라면 차있는 만큼은 읽은 뒤 zero 리턴

채널의 제로값은 nil인데, 송신,수신 작업이 데드락을 발생시키는 nil 채널이 유용한 경우가 있다.

1. select 문에서 채널이 nil 값이면 선택되지 않는다.
이 특성을 이용해서 특성 채널을 활성화하거나, 비활성화하는데 nil 채널을 사용할 수 있다.

2. closed 채널인지를 확인
closed 채널인지를 확인하기 위해선 두번째 반환값. ok를 받고 분기해야 한다.
이때 채널을 close 되었고 읽을 값이 없다면, nil값을 넣어줌으로써 채널이 현재 비활성화상태인지 마킹할 수 있다.

#### 8.8 예제: 동시 디렉토리 탐색

linux du 커맨드와 비슷하게 동작하는 프로그램을 만든다.

du1 예제는 프로그램이 모두 수행된 뒤 총 파일 개수와 크기 합을 출력한다.

du2 예제는 -v 플래그를 받으면 정해진 주기마다 그때까지 집계된 총 파일 개수와 크기 합을 출력한다.

> 이 예제에서는 break "레이블"을 사용하여 for, select 문을 모두 빠져나온다.

du3 에제는 재귀적으로 순차적으로 디렉토리를 순회하는게 아니라 sync.WaitGroup, 세마포어를 이용해 병렬적으로 수행한다.

> 카운팅 세마포어를 사용할때는 io작업이 일어나는 작업과 가장 근접한 부분이 로직을 작성하는게 좋다. 여기선 ioutil.ReadDir()

#### 8.9 취소

때로는 고루틴이 현재 수행 중인 작업을 중지하게 지시할 필요가 있다. 예를 들면 서버 작업 중에 클라이언트 연결이 끊기는 경우이다.

한 고루틴에서 다른 고루틴을 직접 종료할 수는 없다. countdown 에제에서는 abort라는 채널에 값을 하나 보내서 고루틴이 스스로 종료할 수 있게 했다.

종료시켜야하는 고루틴이 여러개가 될 경우 abort 채널에 값을 보내는 것으로는 고루틴 종료가 힘들 수 있다.

고루틴 중 일부가 이미 종료되어서 abort 채널로 더 이상 송신할 수 없는 상태가 될 수도 있고,

고루틴이 또다른 고루틴을 생성하는 경우에은 abort 채널로 보내는 이벤트의 개수가 부족할 수도 있다.

채널에 이벤트를 브로드캐스트해서 여러 고루틴에서 취소 이벤트가 발생했음을 알 수 있고, 나중에도 알 수 있는 메커니즘이 필요하다.

채널이 닫히고 송신된 모든 값을 소진한 이후에는 수신 작업이 즉시 수행되고 제로값을 반환하는 특성을 이용하여 브로드캐스트를 구현할 수 있다. (du4)

done 채널을 하나 만들고 고루틴들을 종료시켜야하는 경우 close(done)을 호출해 해당 채널을 수신하고 있는 곳에 즉시 반환값을 주도록 한다.

로직들을 빠르게 취소시키려면 더 많이 로직에 관여해야한다. du4 예제에서는 고루틴으로 동작하는 walkDir 함수와 세마포어 경합이 일어날 수 있는 dirents 함수에 취소여부를 확인하는 로직이 추가됬다.

취소시에 고루틴이 모두 잘 종료되었는지 확인하기 위해서 취소 이벤트 발생시에 select 문에서 메인고루틴으로 반환하는게 아니라 panic을 발생시키면 고루틴 스택을 볼 수 있다.

```
panic: canceled

goroutine 1 [running]:
main.main()
        /Users/1111252/workspace/project/book-The-Go-Programming-Language/ch8/du4/main.go:64 +0x425
exit status 2
```

#### 8.10 예제: 채팅 서버

마지막 채팅 서버 예제에는 4가지 고루틴이 있다. main, broadcaster, handleConn, clientWriter

각 클라이언트 접속마다 클라이언트의 입력을 받아주고, 연결, 종료를 관리하는 handleConn 한개와

서버(브로드캐스터)에서 메시지를 받으면 클라이언트에 보내주는 clientWriter 한개씩이 있다.




