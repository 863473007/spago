// Copyright 2019 spaGO Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"github.com/nlpodyssey/spago/examples/mnist/internal/mnist"
	"github.com/nlpodyssey/spago/examples/mnist/third_party/GoMNIST"
	"github.com/nlpodyssey/spago/pkg/mat/rand"
	"github.com/nlpodyssey/spago/pkg/ml/ag"
	"github.com/nlpodyssey/spago/pkg/ml/nn"
	"github.com/nlpodyssey/spago/pkg/ml/optimizers/gd"
	"github.com/nlpodyssey/spago/pkg/ml/optimizers/gd/adam"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
)

func main() {
	// go tool pprof http://localhost:6060/debug/pprof/profile
	go func() { log.Println(http.ListenAndServe("localhost:6060", nil)) }()

	modelPath := os.Args[1]
	var datasetPath string
	if len(os.Args) > 2 {
		datasetPath = os.Args[2]
	} else {
		// assuming default path
		datasetPath = "examples/mnist/third_party/GoMNIST/data"
	}

	batchSize := 50
	epochs := 20

	// read dataset
	trainSet, testSet, err := GoMNIST.Load(datasetPath)
	if err != nil {
		panic("Error reading MNIST data.")
	}

	// new model initialized with random weights
	model := mnist.NewCNN(
		9,   // kernelSizeX
		9,   // kernelSizeY
		1,   // inputChannels
		10,  // outputChannels
		5,   // maxPoolingRows
		5,   // maxPoolingCols
		160, // hidden
		10,  // out
		ag.OpReLU,
		ag.OpIdentity, // The CrossEntropy loss doesn't require explicit Softmax activation
	)

	rndGen := rand.NewLockedRand(uint64(42))
	mnist.InitCNN(model, rndGen)

	trainer := mnist.NewTrainer(
		model,
		gd.NewOptimizer(adam.New(adam.NewDefaultConfig()), nn.NewDefaultParamsIterator(model)),
		epochs,
		batchSize,
		true,
		mnist.Dataset{Set: trainSet, FeaturesAsVector: false},
		mnist.Dataset{Set: testSet, FeaturesAsVector: false},
		modelPath,
		42,
	)
	trainer.Enjoy() // :)
}
