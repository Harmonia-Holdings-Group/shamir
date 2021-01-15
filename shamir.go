package shamir

import "fmt"

type Shamir struct{}

func (s *Shamir) SayHi() string {
	return fmt.Sprintf("Hi!")
}
