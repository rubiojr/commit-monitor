/*
 *    Copyright (C) 2020 Sergio Rubio
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
	"github.com/muesli/beehive/bees"
)

// StdoutBeeFactory is a factory for StdoutBees.
type StdoutBeeFactory struct {
	bees.BeeFactory
}

// New returns a new Bee instance configured with the supplied options.
func (factory *StdoutBeeFactory) New(name, description string, options bees.BeeOptions) bees.BeeInterface {
	bee := StdoutBee{
		Bee: bees.NewBee(name, factory.ID(), description, options),
	}
	bee.ReloadOptions(options)

	return &bee
}

// ID returns the ID of this Bee.
func (factory *StdoutBeeFactory) ID() string {
	return "stdoutbee"
}

// Name returns the name of this Bee.
func (factory *StdoutBeeFactory) Name() string {
	return "Stdout"
}

// Description returns the description of this Bee.
func (factory *StdoutBeeFactory) Description() string {
	return "Print events to stdout"
}

// Actions describes the available actions provided by this Bee.
func (factory *StdoutBeeFactory) Actions() []bees.ActionDescriptor {
	actions := []bees.ActionDescriptor{
		{
			Namespace:   factory.Name(),
			Name:        "print",
			Description: "Print event to stdout",
			Options: []bees.PlaceholderDescriptor{
				{
					Name:        "text",
					Type:        "string",
					Description: "The text to print",
					Mandatory:   true,
				},
			},
		},
	}
	return actions
}

func init() {
	f := StdoutBeeFactory{}
	bees.RegisterFactory(&f)
}
