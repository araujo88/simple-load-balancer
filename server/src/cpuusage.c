// In this program:

// We read CPU statistics from /proc/stat.
// We then sleep for 1 second and read the statistics again.
// We calculate the CPU utilization during that second.
// This C program should be executed in a Linux environment.

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>

struct CpuStat {
    long user;
    long nice;
    long system;
    long idle;
};

int get_cpu_stat(struct CpuStat *stat) {
    FILE *file = fopen("/proc/stat", "r");
    if (!file) {
        perror("fopen");
        return -1;
    }

    int ret = fscanf(file, "cpu %ld %ld %ld %ld", &stat->user, &stat->nice, &stat->system, &stat->idle);
    fclose(file);

    return ret == 4 ? 0 : -1;
}

double compute_cpu_utilization(const struct CpuStat *prev, const struct CpuStat *curr) {
    long prev_total = prev->user + prev->nice + prev->system + prev->idle;
    long curr_total = curr->user + curr->nice + curr->system + curr->idle;
    long idle_diff = curr->idle - prev->idle;
    long total_diff = curr_total - prev_total;
    return 100*(1.0 - ((double) idle_diff / total_diff));
}

int main() {
    struct CpuStat prev_stat, curr_stat;

    if (get_cpu_stat(&prev_stat) < 0) {
        fprintf(stderr, "Error retrieving CPU stat\n");
        return 1;
    }

    sleep(1);

    if (get_cpu_stat(&curr_stat) < 0) {
        fprintf(stderr, "Error retrieving CPU stat\n");
        return 1;
    }

    double cpu_utilization = compute_cpu_utilization(&prev_stat, &curr_stat);
    printf("CPU Utilization: %.2lf%%\n", cpu_utilization);
    return 0;
}
