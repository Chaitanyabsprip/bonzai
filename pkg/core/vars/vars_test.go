package vars_test

import (
	"fmt"
	"os"

	"github.com/rwxrob/bonzai/pkg/core/futil"
	"github.com/rwxrob/bonzai/pkg/core/vars"
)

func ExampleSet() {

	file := `testdata/settest.properties`

	defer func() {
		err := os.Remove(file)
		fmt.Println(err)
	}()

	if err := vars.Set(`somekey`, `someval`, file); err != nil {
		fmt.Println(err)
	}

	if err := vars.Set(`otherkey`, ``, file); err != nil {
		fmt.Println(err)
	}

	futil.Cat(file)

	// Output:
	// somekey=someval
	// otherkey=
	// <nil>
}
