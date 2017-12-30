# validate

Easily validate struct fields and values.

# Usage

```go
import (
	"fmt"
	"time"

	"github.com/olivoil/pkg/validate"
)

type Input struct {
	Name string `validate:"required,max(15)"`
	User struct {
		Balance float64 `validate:"-"`
	} `validate:"required"`
	Amount float64 `validate:"lte($.User.Balance)`
	Duration time.Duration `validate:"lte(5m)"`
	Contacts []string `validate:"required,each(required,email|phone)"`
}

input := Input{
	Name: "morethan15characters",
	User: struct {
		Balance: 10.0,
	},
	Amount: 11.0,
	Duration: 6 * time.Minute,
	Contacts: []string{"notavalid@email", "", "123.555.789"},
}

if err := validate.Struct(); err != nil {
	fmt.Println(err.Error())
}
```