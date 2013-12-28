package ros


type Node interface {
    NewPublisher(topic string, msgType MessageType) Publisher
    NewSubscriber(topic string, msgType MessageType, callback interface{}) Subscriber
    NewServiceClient(service string, srvType ServiceType) ServiceClient
    NewServiceServer(service string, srvType ServiceType, callback interface{}) ServiceServer

    OK() bool
    SpinOnce()
    Spin()
    Shutdown()

    GetParam(name string) (interface{}, error)
    SetParam(name string, value interface{}) error
    HasParam(name string) (bool, error)
    SearchParam(name string) (string, error)
    DeleteParam(name string) error

    Logger() Logger
}


func NewNode(name string) Node {
    return newDefaultNode(name)
}


type Publisher interface {
    Publish(msg Message)
    Shutdown()
}


type Subscriber interface {
    Shutdown()
}


type ServiceHandler interface {}


type ServiceFactory interface {
    Name() string
    MD5Sum() string
}


type ServiceServer interface {
    Shutdown()
}


type ServiceClient interface {
    Call(srv Service) error
    Shutdown()
}


type TimeBase interface {
    Sec() uint32
    SetSec(sec uint32)
    NSec() uint32
    SetNSec(nsec uint32)
    Set(sec uint32, nsec uint32)
    IsZero() bool
    ToSec() float64
    ToNSec() uint64
    FromSec(sec float64)
    FromNSec(nsec uint64)
    Add(d Duration) Time
    Sub(d Duration) Time
    Diff(t Time) Duration
    Cmp(t Time) int
}


type Time interface {
    TimeBase
}


func Now() Time {
    return now()
}


func NewTime() Time {
    return new(_Time)
}


type Duration interface {
    TimeBase
    Sleep() error
}


func NewDuration() Duration {
    return new(duration)
}


type Rate interface {
    CycleTime() Duration
    ExpectedCycleTime() Duration
    Reset()
    Sleep() error
}


func NewRate(frequency float64) Rate {
    return newRate(frequency)
}


func NewRateFromCycleTime(d Duration) Rate {
    return newRateFromCycleTime(d)
}





