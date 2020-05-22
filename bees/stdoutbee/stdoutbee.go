/*
 *    Copyright (C) 2020      Sergio Rubio
 *
 *    This program is free software: you can redistribute it and/or modify
 *    it under the terms of the GNU Affero General Public License as published
 *    by the Free Software Foundation, either version 3 of the License, or
 *    (at your option) any later version.
 *
 *    This program is distributed in the hope that it will be useful,
 *    but WITHOUT ANY WARRANTY; without even the implied warranty of
 *    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *    GNU Affero General Public License for more details.
 *
 *    You should have received a copy of the GNU Affero General Public License
 *    along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 *    Authors:
 *      Sergio Rubio <sergio@rubio.im>
 */

package stdoutbee

import (
	"fmt"

	"github.com/muesli/beehive/bees"
	log "github.com/sirupsen/logrus"
)

// StdoutBee is an example for a Bee skeleton, designed to help you get started
// with writing your own Bees.
type StdoutBee struct {
	bees.Bee
}

// Run executes the Bee's event loop.
func (mod *StdoutBee) Run(eventChan chan bees.Event) {
}

// Action triggers the action passed to it.
func (mod *StdoutBee) Action(action bees.Action) []bees.Placeholder {
	outs := []bees.Placeholder{}

	switch action.Name {
	case "print":
		text := action.Options.Value("text")
		fmt.Println(text)
	default:
		log.Debug("ignoring non-supported action " + action.Name)
	}

	return outs
}

// ReloadOptions parses the config options and initializes the Bee.
func (mod *StdoutBee) ReloadOptions(options bees.BeeOptions) {
	mod.SetOptions(options)
}
