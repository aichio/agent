#include "akya.h"

typedef struct akya_ring_s{
    unsigned int offset;
    unsigned int size;
    unsigned int in;
    unsigned int out;
}akya_ring_t;

#define min(x, y) ({                        \
        typeof(x) _min1 = (x);          \
        typeof(y) _min2 = (y);          \
        (void) (&_min1 == &_min2);      \
        _min1 < _min2 ? _min1 : _min2; })


unsigned int akya_ring_get(akya_ring_t *ring ,unsigned char *data ,unsigned int len)
{
    unsigned int l;
    unsigned char *buffer = (unsigned char *)ring + ring->offset;

    len = min(len, ring->in - ring->out);

    l = min(len ,ring->size - (ring->out & (ring->size - 1)));
    memcpy(data ,buffer + (ring->out & (ring->size - 1)) ,l);

    memcpy(data + l ,buffer,len - l);

    ring->out += len;

    return len;
}

int AkyaRingWait(int fd)
{
    struct pollfd fds;
    int ret;

    fds.fd = fd;
    fds.events = POLLIN;

    return poll(&fds ,1 ,5000);
}

int AkyaRingGet(void *d ,char *data ,int len)
{
	akya_ring_t *ring = (akya_ring_t *)d;
    int ret = 0;

    ret = akya_ring_get(ring ,data ,len);
    return ret;
}