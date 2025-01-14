#ifdef __cplusplus
extern "C" {
#endif

__global__ void medianFilter(unsigned char* imgIn, unsigned char* imgOut, int width, int height) {
    const int windowSize = 3;
    const int halfWindow = windowSize / 2;

    // Thread coordinates
    int x = blockIdx.x * blockDim.x + threadIdx.x;
    int y = blockIdx.y * blockDim.y + threadIdx.y;

    // Threads outside image boundary
    if (x >= width || y >= height) return;

    unsigned char window[windowSize * windowSize];
    int count = 0;

    for (int wy = -halfWindow; wy <= halfWindow; ++wy) {
        for (int wx = -halfWindow; wx <= halfWindow; ++wx) {
            int nx = min(max(x + wx, 0), width - 1);
            int ny = min(max(y + wy, 0), height - 1);
            window[count++] = imgIn[ny * width + nx];
        }
    }

    for (int i = 1; i < count; ++i) {
        unsigned char key = window[i];
        int j = i - 1;
        while (j >= 0 && window[j] > key) {
            window[j + 1] = window[j];
            j--;
        }
        window[j + 1] = key;
    }

    imgOut[y * width + x] = window[count / 2];
}

#ifdef __cplusplus
}
#endif
