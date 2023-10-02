#include <stdio.h>
#include <stdlib.h>

int main() {
    FILE *file = fopen("/proc/meminfo", "r");
    if (file == NULL) {
        perror("Unable to open /proc/meminfo");
        return 1;
    }

    unsigned long long totalMemory;
    unsigned long long freeMemory;
    unsigned long long availableMemory;
    unsigned long long bufferMemory;
    unsigned long long cachedMemory;

    char line[256];
    while (fgets(line, sizeof(line), file)) {
        if (sscanf(line, "MemTotal: %llu kB", &totalMemory) == 1)
            continue;
        if (sscanf(line, "MemFree: %llu kB", &freeMemory) == 1)
            continue;
        if (sscanf(line, "MemAvailable: %llu kB", &availableMemory) == 1)
            continue;
        if (sscanf(line, "Buffers: %llu kB", &bufferMemory) == 1)
            continue;
        if (sscanf(line, "Cached: %llu kB", &cachedMemory) == 1)
            continue;
    }

    fclose(file);

    unsigned long long usedMemory = totalMemory - availableMemory;
    double usedMemoryPercentage = ((double)usedMemory / totalMemory) * 100.0;

    printf("Memory Usage: %.2f%%\n", usedMemoryPercentage);

    return 0;
}
