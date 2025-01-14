#ifdef __cplusplus
extern "C"
{
#endif

    __global__ void medianFilter(uint8_t* imgIn, uint8_t* imgOut, int width, int height) {
        const int windowSize = 3;
        int x = blockIdx.x * blockDim.x + threadIdx.x;
        int y = blockIdx.y * blockDim.y + threadIdx.y;
        
        if (x >= width || y >= height) return;

        int halfWindow = windowSize / 2;
        uint8_t window[windowSize * windowSize];

        int count = 0;
        for (int wy = -halfWindow; wy <= halfWindow; ++wy) {
            for (int wx = -halfWindow; wx <= halfWindow; ++wx) {
                int nx = min(max(x + wx, 0), width - 1);
                int ny = min(max(y + wy, 0), height - 1);
                window[count++] = imgIn[ny * width + nx];
            }
        }

        for (int i = 0; i < count - 1; ++i) {
            for (int j = i + 1; j < count; ++j) {
                if (window[i] > window[j]) {
                    uint8_t temp = window[i];
                    window[i] = window[j];
                    window[j] = temp;
                }
            }
        }

        imgOut[y * width + x] = window[count / 2];
    }


#ifdef __cplusplus
}
#endif