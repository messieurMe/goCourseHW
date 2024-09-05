package fact

import (
	"fmt"
	"golang.org/x/sync/errgroup"
	"io"
	"strconv"
	"strings"
	"sync"
)

type Input struct {
	NumsOfGoroutine int   // n - число горутин
	Numbers         []int // слайс чисел, которые необходимо факторизовать
}

type Factorization interface {
	Work(Input, io.Writer) error
}

func NewFactorization() *NewFactorizationImpl {
	return &NewFactorizationImpl{}
}

type NewFactorizationImpl struct {
}

type ConcurrentWriter struct {
	mutex   *sync.Mutex
	writer  io.Writer
	counter int
}

func (c *ConcurrentWriter) Write(toWrite string) error {
	c.mutex.Lock()

	curr := c.counter
	_, err := c.writer.Write([]byte(fmt.Sprintf("line %d, %s\n", curr, toWrite)))

	c.counter += 1
	c.mutex.Unlock()
	return err
}

func toStr(x int) string {
	return strconv.Itoa(x)
}

func factorCore(rawNumber int) *[]int {
	var factorization []int
	var number = rawNumber
	if rawNumber < 0 {
		factorization = append(factorization, -1)
		number = -1 * rawNumber
	}

	var nextPrime int

	for number != 1 {
		nextPrime = 2
		for number%nextPrime != 0 {
			nextPrime++
		}
		number /= nextPrime
		factorization = append(factorization, nextPrime)
	}

	if len(factorization) == 0 || rawNumber == -1 {
		factorization = append(factorization, number)
	}
	return &factorization
}

func factor(rawNumber int, counter *ConcurrentWriter) error {

	var factorization = *factorCore(rawNumber)

	resultStringBuilder := strings.Builder{}
	resultStringBuilder.WriteString(fmt.Sprintf("%d = %d", rawNumber, factorization[0]))
	for i := 1; i < len(factorization); i++ {
		resultStringBuilder.WriteString(" * " + toStr(factorization[i]))
	}

	// without lock we cannot guarantee order of lines
	// for example
	// goroutine 1 read atomicInt and increased it
	// goroutine 2 read atomicInt and increased it
	// goroutine 2 writes line
	// goroutine 1 writes line
	// So we write in different goroutines, but with lock
	return counter.Write(resultStringBuilder.String())
}

type Finished struct{}

func (f *NewFactorizationImpl) Work(input Input, writer io.Writer) error {
	var onii = make(chan Finished, input.NumsOfGoroutine)

	group := errgroup.Group{}
	group.SetLimit(input.NumsOfGoroutine)

	writeMutex := sync.Mutex{}
	var counter = ConcurrentWriter{
		mutex:   &writeMutex,
		writer:  writer,
		counter: 1,
	}

	for _, i := range input.Numbers {
		localI := i
		group.Go(func() error {
			return factor(localI, &counter)
		})
	}

	foundError := group.Wait()
	close(onii)

	return foundError
}
