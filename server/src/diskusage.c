#include <stdio.h>
#include <sys/statvfs.h>

int main() {
    struct statvfs stat;

    if (statvfs("/", &stat) != 0) {
        perror("statvfs");
        return 1;
    }

    unsigned long long totalSpace = (unsigned long long) stat.f_blocks * stat.f_frsize;
    unsigned long long freeSpace = (unsigned long long) stat.f_bfree * stat.f_frsize;
    unsigned long long usedSpace = totalSpace - freeSpace;
    double usedPercent = (double) usedSpace / totalSpace * 100;

    printf("Disk usage: %.2lf%%\n", usedPercent);

    return 0;
}
