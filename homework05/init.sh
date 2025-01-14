#!/bin/bash
#SBATCH --nodes=1
#SBATCH --reservation=fri

module load Go
module load CUDA

export CGO_CFLAGS=$(pkg-config --cflags cudart-12.6)
export CGO_LDFLAGS=$(pkg-config --libs cudart-12.6)
export PATH="~/go/bin/:$PATH"

go install github.com/InternatBlackhole/cudago/CudaGo@latest