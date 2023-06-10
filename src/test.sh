#ginkgo -r --randomizeAllSpecs --randomizeSuites --failOnPending --cover --trace --race --progress
ginkgo -r -tags unit --randomizeAllSpecs --randomizeSuites --failOnPending --cover --trace --race --progress -v
ginkgo -r -tags integration --randomizeAllSpecs --randomizeSuites --failOnPending --cover --trace --race --progress -v
