// Copyright 2020 spaGO Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"os"

	"github.com/nlpodyssey/spago/cmd/bert/internal/app"
)

func main() {
	app.NewBertApp().Run(os.Args)
}
