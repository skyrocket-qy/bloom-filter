# bloom-filter

$$
m = -\frac{n \ln p}{(\ln 2)^2}
$$

$$
k = \frac{m}{n} \ln 2
$$

## Definition
m = total number of bits in the Bloom filter’s bit array
n = expected number of distinct items
p = target false positive rate
k = number of hash functions

## Feature
### Formula
- Given n and p, calculate m and k
### Plot
- Given n = 10e6, plot x = p and y = m(MB)
- Given n = 10e6, plot x = p and y = checkTime(ms)
- Given n = 10e4, p = 0.001, plot x = realAmount and y = fpRate