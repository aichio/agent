#include <unistd.h>
#include <stdio.h>
#include <string.h>
#include <sys/mman.h>
#include <sys/types.h>
#include <fcntl.h>
#include <poll.h>

int AkyaRingGet(void *d ,char *data ,int len);
int AkyaRingWait(int fd);