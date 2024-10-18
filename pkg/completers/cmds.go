package completers

import (
	bonzai "github.com/rwxrob/bonzai/pkg"
	"github.com/rwxrob/bonzai/pkg/fn/filt"
	"github.com/rwxrob/bonzai/pkg/set"
)

type cmds struct{}

var Cmds = new(cmds)

// Complete resolves completion as follows:
//
//  1. If leaf has Comp function, delegate to it
//
//  2. If leaf has no arguments, return all Cmds and Params
//
//  3. If first argument is the name of a Command return it only even
//     if in the Hidden list
//
//  4. Otherwise, return every Command or Param that is not in the
//     Hidden list and HasPrefix matching the first arg
//
// See bonzai.Completer.
func (cmds) Complete(x bonzai.Command, args ...string) []string {

	// if has completer, delegate
	if c := x.GetComp(); c != nil {
		return c.Complete(x, args...)
	}

	// not sure we've completed the command name itself yet
	if len(args) == 0 {
		return []string{x.GetName()}
	}

	// build list of visible commands and params
	list := []string{}
	list = append(list, x.GetCommandNames()...)
	list = append(list, x.GetParams()...)
	list = append(list, x.GetShortcuts()...)
	list = set.MinusAsString[string, string](list, x.GetHidden())

	if len(args) == 0 {
		return list
	}

	return filt.HasPrefix(list, args[0])
}
